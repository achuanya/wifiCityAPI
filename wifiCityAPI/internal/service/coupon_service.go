package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin/wifiCityAPI/internal/models"
	"github.com/gin-gonic/gin/wifiCityAPI/pkg/database"
	"gorm.io/gorm"
)

// CouponService 提供了优惠券相关的业务逻辑
type CouponService struct{}

// CreateCouponInput 定义了创建优惠券的输入
type CreateCouponInput struct {
	CouponName        string  `json:"coupon_name" binding:"required"`
	CouponCode        string  `json:"coupon_code"`
	CouponType        string  `json:"coupon_type" binding:"required"`
	Value             float64 `json:"value" binding:"required"`
	MinPurchaseAmount float64 `json:"min_purchase_amount"`
	UsageLimitPerUser int     `json:"usage_limit_per_user"`
	TotalQuantity     int     `json:"total_quantity"`
	StartTime         string  `json:"start_time" binding:"required"` // "2006-01-02 15:04:05"
	EndTime           string  `json:"end_time" binding:"required"`
	ValidityDays      int     `json:"validity_days"`
	StoreID           *uint   `json:"store_id"`
	Description       string  `json:"description"`
}

// CreateCoupon 创建一个新的优惠券。
// 在事务中执行。
func (s *CouponService) CreateCoupon(input *CreateCouponInput) (*models.Coupon, error) {
	// 解析时间字符串
	startTime, err := time.Parse("2006-01-02 15:04:05", input.StartTime)
	if err != nil {
		return nil, fmt.Errorf("无效的开始时间格式: %w", err)
	}
	endTime, err := time.Parse("2006-01-02 15:04:05", input.EndTime)
	if err != nil {
		return nil, fmt.Errorf("无效的结束时间格式: %w", err)
	}

	coupon := models.Coupon{
		CouponName:        input.CouponName,
		CouponCode:        input.CouponCode,
		CouponType:        input.CouponType,
		Value:             input.Value,
		MinPurchaseAmount: input.MinPurchaseAmount,
		UsageLimitPerUser: input.UsageLimitPerUser,
		TotalQuantity:     input.TotalQuantity,
		StartTime:         startTime,
		EndTime:           endTime,
		ValidityDays:      input.ValidityDays,
		StoreID:           input.StoreID,
		Description:       input.Description,
		Status:            1, // 默认启用
	}

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&coupon).Error
	})

	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

// GetCouponByID 根据ID获取优惠券详情
func (s *CouponService) GetCouponByID(id uint) (*models.Coupon, error) {
	var coupon models.Coupon
	err := database.DB.WithContext(context.Background()).First(&coupon, id).Error
	return &coupon, err
}

// GetCouponsInput 定义了查询优惠券的输入
type GetCouponsInput struct {
	StoreID  *uint `form:"store_id"`
	Status   *int8 `form:"status"`
	Page     int   `form:"page"`
	PageSize int   `form:"pageSize"`
}

// GetCoupons 查询优惠券列表
func (s *CouponService) GetCoupons(input *GetCouponsInput) ([]models.Coupon, int64, error) {
	var coupons []models.Coupon
	var total int64

	db := database.DB.WithContext(context.Background()).Model(&models.Coupon{})
	if input.StoreID != nil {
		db = db.Where("store_id = ? OR store_id IS NULL", *input.StoreID) // 门店券或平台通用券
	}
	if input.Status != nil {
		db = db.Where("status = ?", *input.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 10
	}
	offset := (input.Page - 1) * input.PageSize

	err := db.Order("created_at DESC").Offset(offset).Limit(input.PageSize).Find(&coupons).Error
	return coupons, total, err
}

// UpdateCouponInput 定义了更新优惠券的输入
type UpdateCouponInput struct {
	CouponName        string   `json:"coupon_name"`
	Description       string   `json:"description"`
	MinPurchaseAmount *float64 `json:"min_purchase_amount"`
	TotalQuantity     *int     `json:"total_quantity"`
	Status            *int8    `json:"status"`
}

// UpdateCoupon 更新一个已存在的优惠券。
// 在事务中执行"先读后写"。
func (s *CouponService) UpdateCoupon(id uint, input *UpdateCouponInput) (*models.Coupon, error) {
	var coupon models.Coupon
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 查找记录
		if err := tx.First(&coupon, id).Error; err != nil {
			return err
		}

		// 2. 更新字段
		if input.CouponName != "" {
			coupon.CouponName = input.CouponName
		}
		if input.Description != "" {
			coupon.Description = input.Description
		}
		if input.MinPurchaseAmount != nil {
			coupon.MinPurchaseAmount = *input.MinPurchaseAmount
		}
		if input.TotalQuantity != nil {
			coupon.TotalQuantity = *input.TotalQuantity
		}
		if input.Status != nil {
			coupon.Status = *input.Status
		}

		// 3. 保存更新
		return tx.Save(&coupon).Error
	})

	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

// DeleteCoupon 删除一个优惠券。
// 在事务中执行。
func (s *CouponService) DeleteCoupon(id uint) error {
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&models.Coupon{}, id)
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
