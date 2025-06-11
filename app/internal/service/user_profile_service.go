package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"app/internal/models"
	"app/pkg/database"

	"gorm.io/gorm"
)

// UserProfileService 提供了用户相关的业务逻辑
type UserProfileService struct{}

// CreateOrUpdateUserInput 定义了创建或更新用户时的输入
type CreateOrUpdateUserInput struct {
	UserUnionID     string `json:"user_union_id" binding:"required"`
	OpenID          string `json:"open_id"`
	WechatNickname  string `json:"wechat_nickname"`
	WechatAvatarURL string `json:"wechat_avatar_url"`
	PhoneNumber     string `json:"phone_number"`
	Gender          *int8  `json:"gender"`
	Language        string `json:"language"`
	Country         string `json:"country"`
	Province        string `json:"province"`
	City            string `json:"city"`
}

// CreateOrUpdateUserProfile 创建或更新一个用户的信息。
// 它会根据提供的UnionID判断是创建新用户还是更新现有用户。
// 整个操作在一个事务中完成。
func (s *UserProfileService) CreateOrUpdateUserProfile(input *CreateOrUpdateUserInput) (*models.UserProfile, error) {
	var user models.UserProfile
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 使用 FirstOrInit 查找或初始化用户。
		// 如果用户存在，记录将被加载到 user 变量中。
		// 如果不存在，user 将根据查询条件（UnionID）和输入（input）被初始化。
		if err := tx.Where(&models.UserProfile{UserUnionID: input.UserUnionID}).
			Attrs(&models.UserProfile{
				OpenID:          input.OpenID,
				WechatNickname:  input.WechatNickname,
				WechatAvatarURL: input.WechatAvatarURL,
				PhoneNumber:     input.PhoneNumber,
				Language:        input.Language,
				Country:         input.Country,
				Province:        input.Province,
				City:            input.City,
			}).FirstOrInit(&user).Error; err != nil {
			return err
		}

		// 更新字段
		user.OpenID = input.OpenID
		user.WechatNickname = input.WechatNickname
		user.WechatAvatarURL = input.WechatAvatarURL
		user.PhoneNumber = input.PhoneNumber
		user.Language = input.Language
		user.Country = input.Country
		user.Province = input.Province
		user.City = input.City
		if input.Gender != nil {
			user.Gender = *input.Gender
		}

		// Save 会自动处理创建或更新
		return tx.Save(&user).Error
	})

	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUnionID 根据 UnionID 获取用户详情
func (s *UserProfileService) GetUserByUnionID(unionID string) (*models.UserProfile, error) {
	var user models.UserProfile
	err := database.DB.WithContext(context.Background()).Where("user_union_id = ?", unionID).First(&user).Error
	return &user, err
}

// GetUserByOpenID 根据 OpenID 获取用户详情
func (s *UserProfileService) GetUserByOpenID(openID string) (*models.UserProfile, error) {
	var user models.UserProfile
	err := database.DB.WithContext(context.Background()).Where("open_id = ?", openID).First(&user).Error
	return &user, err
}

// GetUserByPhone 根据手机号获取用户详情
func (s *UserProfileService) GetUserByPhone(phone string) (*models.UserProfile, error) {
	var user models.UserProfile
	err := database.DB.WithContext(context.Background()).Where("phone_number = ?", phone).First(&user).Error
	return &user, err
}

// BindPhoneNumber 用户绑定手机号
func (s *UserProfileService) BindPhoneNumber(unionID string, phoneNumber string, countryCode string) (*models.UserProfile, error) {
	var user models.UserProfile

	// 检查用户是否存在
	result := database.DB.Where("user_union_id = ?", unionID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, result.Error
	}

	// 检查该手机号是否已被其他用户绑定
	var existingUser models.UserProfile
	result = database.DB.Where("phone_number = ? AND user_union_id != ?", phoneNumber, unionID).First(&existingUser)
	if result.Error == nil {
		return nil, fmt.Errorf("该手机号已被其他用户绑定")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	// 更新用户手机号
	user.PhoneNumber = phoneNumber
	user.PhoneCountryCode = countryCode

	if err := database.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// UnbindPhoneNumber 用户解绑手机号
func (s *UserProfileService) UnbindPhoneNumber(unionID string) (*models.UserProfile, error) {
	var user models.UserProfile

	// 检查用户是否存在
	result := database.DB.Where("user_union_id = ?", unionID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, result.Error
	}

	// 清空用户手机号
	user.PhoneNumber = ""
	user.PhoneCountryCode = ""

	if err := database.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserScanHistoryInput 定义了获取用户扫码历史的输入参数
type GetUserScanHistoryInput struct {
	UserUnionID string `form:"user_union_id" binding:"required"`
	Page        int    `form:"page"`
	PageSize    int    `form:"pageSize"`
	StartDate   string `form:"start_date"` // 格式: YYYY-MM-DD
	EndDate     string `form:"end_date"`   // 格式: YYYY-MM-DD
}

// UserScanHistoryItem 用户扫码历史的结构
type UserScanHistoryItem struct {
	StoreID     uint      `json:"store_id"`
	StoreName   string    `json:"store_name"`
	ScanTime    time.Time `json:"scan_time"`
	SuccessFlag bool      `json:"success_flag"`
	WifiSSID    string    `json:"wifi_ssid"`
	DeviceInfo  string    `json:"device_info"`
	LocationLat float64   `json:"location_lat"`
	LocationLng float64   `json:"location_lng"`
	NetworkType string    `json:"network_type"`
	FailReason  string    `json:"fail_reason,omitempty"`
}

// GetUserScanHistory 获取用户扫码门店历史
func (s *UserProfileService) GetUserScanHistory(input *GetUserScanHistoryInput) ([]UserScanHistoryItem, int64, error) {
	// 验证用户是否存在
	var user models.UserProfile
	if err := database.DB.Where("user_union_id = ?", input.UserUnionID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, fmt.Errorf("用户不存在")
		}
		return nil, 0, err
	}

	query := database.DB.Table("scan_log AS sl").
		Select("sl.log_id, sl.store_id, s.name AS store_name, sl.scan_time, sl.success_flag, sl.wifi_ssid, "+
			"sl.device_info, sl.location_lat, sl.location_lng, sl.network_type, sl.fail_reason_message AS fail_reason").
		Joins("LEFT JOIN store AS s ON sl.store_id = s.store_id").
		Where("sl.user_union_id = ?", input.UserUnionID)

	// 日期范围筛选
	if input.StartDate != "" {
		query = query.Where("sl.scan_time >= ?", input.StartDate)
	}
	if input.EndDate != "" {
		query = query.Where("sl.scan_time <= ?", input.EndDate+" 23:59:59")
	}

	// 查询总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计扫码历史失败: %w", err)
	}

	// 分页
	if input.Page > 0 && input.PageSize > 0 {
		offset := (input.Page - 1) * input.PageSize
		query = query.Offset(offset).Limit(input.PageSize)
	}

	// 排序
	query = query.Order("sl.scan_time DESC")

	// 执行查询
	var results []UserScanHistoryItem
	if err := query.Find(&results).Error; err != nil {
		return nil, 0, fmt.Errorf("查询扫码历史失败: %w", err)
	}

	return results, total, nil
}
