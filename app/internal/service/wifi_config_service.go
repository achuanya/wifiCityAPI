package service

import (
	"context"
	"fmt"

	"app/internal/models"
	"app/pkg/database"

	"gorm.io/gorm"
)

// WifiConfigService 提供了 WIFI 配置相关的业务逻辑
type WifiConfigService struct{}

// CreateWifiConfigInput 定义了创建 WIFI 配置的输入
type CreateWifiConfigInput struct {
	StoreID           uint   `json:"store_id" binding:"required"`
	SSID              string `json:"ssid" binding:"required"`
	PasswordEncrypted string `json:"password_encrypted" binding:"required"`
	EncryptionType    string `json:"encryption_type"`
	WifiType          string `json:"wifi_type"`
	MaxConnections    int    `json:"max_connections"`
}

// CreateWifiConfig 创建一个新的 WIFI 配置
// 它在一个事务中完成此操作。
func (s *WifiConfigService) CreateWifiConfig(input *CreateWifiConfigInput) (*models.WifiConfig, error) {
	wifiConfig := models.WifiConfig{
		StoreID:           input.StoreID,
		SSID:              input.SSID,
		PasswordEncrypted: input.PasswordEncrypted,
		EncryptionType:    input.EncryptionType,
		WifiType:          input.WifiType,
		MaxConnections:    input.MaxConnections,
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&wifiConfig).Error; err != nil {
			return err
		}
		// 可以在此事务中添加其他相关操作，例如更新门店的wifi_count字段
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &wifiConfig, nil
}

// CreateBatchWifiConfigs 批量创建WIFI配置
func (s *WifiConfigService) CreateBatchWifiConfigs(inputs []*CreateWifiConfigInput) ([]models.WifiConfig, error) {
	var createdConfigs []models.WifiConfig

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		for _, input := range inputs {
			config := models.WifiConfig{
				StoreID:           input.StoreID,
				SSID:              input.SSID,
				PasswordEncrypted: input.PasswordEncrypted,
				EncryptionType:    input.EncryptionType,
				WifiType:          input.WifiType,
				MaxConnections:    input.MaxConnections,
			}
			if err := tx.Create(&config).Error; err != nil {
				// 如果任何一个创建失败，则回滚整个事务
				return fmt.Errorf("创建WIFI配置 '%s' 失败: %w", input.SSID, err)
			}
			createdConfigs = append(createdConfigs, config)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdConfigs, nil
}

// GetWifiConfigByID 根据ID获取WIFI配置详情
// 只读操作，将走从库。
func (s *WifiConfigService) GetWifiConfigByID(id uint) (*models.WifiConfig, error) {
	var wifiConfig models.WifiConfig
	err := database.DB.WithContext(context.Background()).First(&wifiConfig, id).Error
	return &wifiConfig, err
}

// GetWifiConfigsByStoreID 根据门店ID获取所有WIFI配置
func (s *WifiConfigService) GetWifiConfigsByStoreID(storeID uint) ([]models.WifiConfig, error) {
	var wifiConfigs []models.WifiConfig
	err := database.DB.WithContext(context.Background()).Where("store_id = ?", storeID).Find(&wifiConfigs).Error
	return wifiConfigs, err
}

// GetWifiConfigsByStoreAndTypeInput 定义获取门店特定类型WIFI配置的输入参数
type GetWifiConfigsByStoreAndTypeInput struct {
	StoreID  uint   `form:"store_id" binding:"required"`
	WifiType string `form:"wifi_type" binding:"required"`
}

// GetWifiConfigsByStoreAndType 获取门店特定类型的WIFI配置
func (s *WifiConfigService) GetWifiConfigsByStoreAndType(input *GetWifiConfigsByStoreAndTypeInput) ([]models.WifiConfig, error) {
	var wifiConfigs []models.WifiConfig

	// 查询指定门店下指定类型的WIFI配置
	err := database.DB.WithContext(context.Background()).
		Where("store_id = ? AND wifi_type = ?", input.StoreID, input.WifiType).
		Find(&wifiConfigs).Error

	if err != nil {
		return nil, fmt.Errorf("查询门店特定类型WIFI配置失败: %w", err)
	}

	return wifiConfigs, nil
}

// UpdateWifiConfigInput 定义了更新WIFI配置的输入
type UpdateWifiConfigInput struct {
	SSID              string `json:"ssid"`
	PasswordEncrypted string `json:"password_encrypted"`
	EncryptionType    string `json:"encryption_type"`
	WifiType          string `json:"wifi_type"`
	MaxConnections    *int   `json:"max_connections"`
}

// UpdateWifiConfig 更新一个已存在的WIFI配置
// 在事务中执行"先读后写"。
func (s *WifiConfigService) UpdateWifiConfig(id uint, input *UpdateWifiConfigInput) (*models.WifiConfig, error) {
	var wifiConfig models.WifiConfig
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 在事务中查找记录
		if err := tx.First(&wifiConfig, id).Error; err != nil {
			return err
		}

		// 2. 更新字段
		if input.SSID != "" {
			wifiConfig.SSID = input.SSID
		}
		if input.PasswordEncrypted != "" {
			wifiConfig.PasswordEncrypted = input.PasswordEncrypted
		}
		if input.EncryptionType != "" {
			wifiConfig.EncryptionType = input.EncryptionType
		}
		if input.WifiType != "" {
			wifiConfig.WifiType = input.WifiType
		}
		if input.MaxConnections != nil {
			wifiConfig.MaxConnections = *input.MaxConnections
		}

		// 3. 在事务中保存
		if err := tx.Save(&wifiConfig).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &wifiConfig, nil
}

// DeleteWifiConfig 删除一个WIFI配置
// 在事务中执行。
func (s *WifiConfigService) DeleteWifiConfig(id uint) error {
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&models.WifiConfig{}, id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
	return err
}

// DeleteBatchWifiConfigs 批量删除多个WIFI配置
func (s *WifiConfigService) DeleteBatchWifiConfigs(ids []uint) error {
	// 使用事务确保数据一致性
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 查询这些WIFI配置关联的门店信息
		var wifiConfigs []models.WifiConfig
		if err := tx.Where("wifi_id IN ?", ids).Find(&wifiConfigs).Error; err != nil {
			return err
		}

		// 创建一个map来存储每个门店的WIFI数量变更
		storeWifiCount := make(map[uint]int)
		for _, wifi := range wifiConfigs {
			storeWifiCount[wifi.StoreID]--
		}

		// 删除WIFI配置
		if err := tx.Where("wifi_id IN ?", ids).Delete(&models.WifiConfig{}).Error; err != nil {
			return err
		}

		// 更新每个门店的WIFI数量
		for storeID, count := range storeWifiCount {
			if count != 0 {
				if err := tx.Exec("UPDATE store SET wifi_count = wifi_count + ? WHERE store_id = ?", count, storeID).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}
