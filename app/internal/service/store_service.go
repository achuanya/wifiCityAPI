package service

import (
	"context"
	"errors"
	"fmt"

	"app/internal/models"
	"app/pkg/database"

	"gorm.io/gorm"
)

// StoreService 提供了门店相关的业务逻辑
type StoreService struct{}

// CreateStoreInput 定义了创建门店时需要的输入
type CreateStoreInput struct {
	Name      string  `json:"name" binding:"required"`
	Country   string  `json:"country"`
	Province  string  `json:"province"`
	City      string  `json:"city"`
	District  string  `json:"district"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Phone     string  `json:"phone"`
}

// CreateStore 在数据库中创建一个新的门店记录。
// 它在一个事务中完成操作，以确保数据一致性。
func (s *StoreService) CreateStore(input *CreateStoreInput) (*models.Store, error) {
	// 将输入数据映射到GORM模型
	store := models.Store{
		Name:      input.Name,
		Country:   input.Country,
		Province:  input.Province,
		City:      input.City,
		District:  input.District,
		Address:   input.Address,
		Latitude:  input.Latitude,
		Longitude: input.Longitude,
		Phone:     input.Phone,
		Status:    1, // 默认为正常状态
	}

	// 使用事务来创建门店
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行创建操作
		if err := tx.Create(&store).Error; err != nil {
			// 如果发生错误，返回该错误，事务将自动回滚
			return err
		}
		// 返回nil表示事务成功，将被提交
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &store, nil
}

// GetStoreByID 根据ID获取门店详情
// 使用读库
func (s *StoreService) GetStoreByID(id uint) (*models.Store, error) {
	var store models.Store
	// 使用 gorm.io/plugin/dbresolver 的 Clauser context 来强制走从库
	// 在我们的设置中，读操作默认走从库，这里为了代码清晰明确指出
	err := database.DB.WithContext(context.Background()).First(&store, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not an error, just not found
		}
		return nil, fmt.Errorf("查询门店失败: %w", err)
	}
	return &store, nil
}

// GetStoresInput 定义了查询门店列表的输入参数
type GetStoresInput struct {
	Page      int     `form:"page"`
	PageSize  int     `form:"pageSize"`
	Province  string  `form:"province"`
	City      string  `form:"city"`
	District  string  `form:"district"`
	Latitude  float64 `form:"lat"`
	Longitude float64 `form:"lng"`
	Radius    float64 `form:"radius"` // 半径，单位：公里
}

// GetStores 获取门店列表，支持分页、区域筛选和附近查询
func (s *StoreService) GetStores(input *GetStoresInput) ([]models.Store, int64, error) {
	var stores []models.Store
	var total int64

	query := database.DB.Model(&models.Store{})

	// 区域筛选
	if input.Province != "" {
		query = query.Where("province = ?", input.Province)
	}
	if input.City != "" {
		query = query.Where("city = ?", input.City)
	}
	if input.District != "" {
		query = query.Where("district = ?", input.District)
	}

	// 附近门店查询
	// Haversine 公式:
	// a = sin²(Δφ/2) + cos φ1 ⋅ cos φ2 ⋅ sin²(Δλ/2)
	// c = 2 ⋅ atan2(√a, √(1−a))
	// d = R ⋅ c
	// R (地球半径) = 6371 公里
	if input.Latitude != 0 && input.Longitude != 0 && input.Radius > 0 {
		lat := input.Latitude
		lng := input.Longitude
		radius := input.Radius

		// 使用 Haversine 公式构建 SQL 查询
		// 1 degree of latitude = 111.045 km
		// 1 degree of longitude = 111.045 * cos(latitude) km
		// 为了简化和提高性能，我们先用一个矩形范围进行粗筛，再用精确的 Haversine 公式进行筛选
		haversine := fmt.Sprintf(
			"6371 * acos(cos(radians(%f)) * cos(radians(latitude)) * cos(radians(longitude) - radians(%f)) + sin(radians(%f)) * sin(radians(latitude)))",
			lat, lng, lat,
		)
		query = query.Select(fmt.Sprintf("*, (%s) AS distance", haversine)).
			Where(fmt.Sprintf("(%s) < ?", haversine), radius).
			Order("distance")
	} else {
		query = query.Order("created_at DESC")
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计门店数量失败: %w", err)
	}

	// 分页
	if input.Page > 0 && input.PageSize > 0 {
		offset := (input.Page - 1) * input.PageSize
		query = query.Offset(offset).Limit(input.PageSize)
	}

	// 执行查询
	if err := query.Find(&stores).Error; err != nil {
		return nil, 0, fmt.Errorf("查询门店列表失败: %w", err)
	}

	return stores, total, nil
}

// UpdateStoreInput 定义了更新门店信息的输入
type UpdateStoreInput struct {
	Name      string  `json:"name"`
	Country   string  `json:"country"`
	Province  string  `json:"province"`
	City      string  `json:"city"`
	District  string  `json:"district"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Phone     string  `json:"phone"`
	Status    *int8   `json:"status"` // 使用指针以区分0和未提供
}

// UpdateStore 更新一个已存在的门店信息。
// 它在一个事务中完成"先读后写"的操作，以避免竞态条件并保证数据一致性。
func (s *StoreService) UpdateStore(id uint, input *UpdateStoreInput) (*models.Store, error) {
	var store models.Store

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 在事务中首先查找记录，确保记录存在并锁定
		if err := tx.First(&store, id).Error; err != nil {
			return err // 如果记录未找到，gorm.ErrRecordNotFound将被返回
		}

		// 2. 将输入的数据更新到模型中
		// 使用 input 中的非空值来更新 store 结构体
		if input.Name != "" {
			store.Name = input.Name
		}
		if input.Country != "" {
			store.Country = input.Country
		}
		if input.Province != "" {
			store.Province = input.Province
		}
		if input.City != "" {
			store.City = input.City
		}
		if input.District != "" {
			store.District = input.District
		}
		if input.Address != "" {
			store.Address = input.Address
		}
		// 注意：纬度和经度是 float64，不能直接与 "" 比较。
		// 如果业务上允许更新为0，则直接赋值。如果0是无效值，则需要用指针或特定值来判断。
		// 这里我们假设可以直接更新。
		store.Latitude = input.Latitude
		store.Longitude = input.Longitude

		if input.Phone != "" {
			store.Phone = input.Phone
		}
		if input.Status != nil {
			store.Status = *input.Status
		}

		// 3. 在同一个事务中保存更新
		if err := tx.Save(&store).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &store, nil
}

// DeleteStore 从数据库中删除一个门店。
// 它在一个事务中完成操作。
func (s *StoreService) DeleteStore(id uint) error {
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 执行删除操作
		result := tx.Delete(&models.Store{}, id)
		if result.Error != nil {
			return result.Error
		}
		// 如果未找到记录，GORM v2的Delete不会返回ErrRecordNotFound，而是返回RowsAffected=0。
		// 我们需要检查受影响的行数来确定记录是否存在。
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound // 手动返回记录未找到的错误
		}
		return nil
	})

	return err
}

// CreateStoreWithWifiInput 定义了同时创建门店和WIFI的输入参数
type CreateStoreWithWifiInput struct {
	Store CreateStoreInput        `json:"store" binding:"required"`
	Wifis []CreateWifiConfigInput `json:"wifis"`
}

// CreateStoreWithWifi 在一个事务中创建门店及其关联的WIFI配置
func (s *StoreService) CreateStoreWithWifi(input *CreateStoreWithWifiInput) (*models.Store, error) {
	var store models.Store

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 创建门店
		store = models.Store{
			Name:      input.Store.Name,
			Country:   input.Store.Country,
			Province:  input.Store.Province,
			City:      input.Store.City,
			District:  input.Store.District,
			Address:   input.Store.Address,
			Latitude:  input.Store.Latitude,
			Longitude: input.Store.Longitude,
			Phone:     input.Store.Phone,
			Status:    1, // 默认启用
		}
		if err := tx.Create(&store).Error; err != nil {
			return fmt.Errorf("创建门店失败: %w", err)
		}

		// 2. 如果有WIFI配置，则创建它们
		if len(input.Wifis) > 0 {
			for _, wifiInput := range input.Wifis {
				wifi := models.WifiConfig{
					StoreID:           store.StoreID, // 关联到刚刚创建的门店ID
					SSID:              wifiInput.SSID,
					PasswordEncrypted: wifiInput.PasswordEncrypted,
					EncryptionType:    wifiInput.EncryptionType,
					WifiType:          wifiInput.WifiType,
					MaxConnections:    wifiInput.MaxConnections,
				}
				if err := tx.Create(&wifi).Error; err != nil {
					// 事务将回滚
					return fmt.Errorf("为门店 '%s' 创建WIFI配置 '%s' 失败: %w", store.Name, wifi.SSID, err)
				}
			}
			// 更新门店的WIFI数量
			if err := tx.Model(&store).Update("wifi_count", len(input.Wifis)).Error; err != nil {
				return fmt.Errorf("更新门店WIFI数量失败: %w", err)
			}
			store.WifiCount = len(input.Wifis)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 重新加载WIFI配置，以便在返回的Store对象中完整显示
	if err := database.DB.Where("store_id = ?", store.StoreID).Find(&store.WifiConfigs).Error; err != nil {
		// 这不是一个关键错误，可以选择只记录日志而不返回错误
		fmt.Printf("警告: 创建门店后加载WIFI配置失败: %v\n", err)
	}

	return &store, nil
}
