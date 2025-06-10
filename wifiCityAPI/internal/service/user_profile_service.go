package service

import (
	"context"

	"github.com/gin-gonic/gin/wifiCityAPI/internal/models"
	"github.com/gin-gonic/gin/wifiCityAPI/pkg/database"
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
