package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin/wifiCityAPI/internal/models"
	"github.com/gin-gonic/gin/wifiCityAPI/pkg/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

// CreateBatchCoupons 批量创建优惠券
func (s *CouponService) CreateBatchCoupons(inputs []*CreateCouponInput) ([]models.Coupon, error) {
	var createdCoupons []models.Coupon
	var couponsToCreate []models.Coupon

	// 先将输入转换为模型对象，并进行基本校验
	for _, input := range inputs {
		startTime, err := time.ParseInLocation("2006-01-02 15:04:05", input.StartTime, time.Local)
		if err != nil {
			return nil, fmt.Errorf("优惠券 '%s' 的开始时间格式无效: %w", input.CouponName, err)
		}
		endTime, err := time.ParseInLocation("2006-01-02 15:04:05", input.EndTime, time.Local)
		if err != nil {
			return nil, fmt.Errorf("优惠券 '%s' 的结束时间格式无效: %w", input.CouponName, err)
		}

		coupon := models.Coupon{
			CouponName:        input.CouponName,
			CouponCode:        input.CouponCode,
			CouponType:        input.CouponType,
			Value:             input.Value,
			MinPurchaseAmount: input.MinPurchaseAmount,
			UsageLimitPerUser: input.UsageLimitPerUser,
			TotalQuantity:     input.TotalQuantity,
			IssuedQuantity:    0,
			StartTime:         startTime,
			EndTime:           endTime,
			ValidityDays:      input.ValidityDays,
			StoreID:           input.StoreID,
			Description:       input.Description,
			Status:            1, // 默认为启用
		}
		couponsToCreate = append(couponsToCreate, coupon)
	}

	// 在一个事务中批量创建
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&couponsToCreate).Error; err != nil {
			return fmt.Errorf("批量创建优惠券失败: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	createdCoupons = couponsToCreate
	return createdCoupons, nil
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
	StartTime         *string  `json:"start_time,omitempty"` // "2006-01-02 15:04:05"
	EndTime           *string  `json:"end_time,omitempty"`
	UsageLimitPerUser *int     `json:"usage_limit_per_user"`
	StoreID           *uint    `json:"store_id"`
}

// UpdateCoupon 更新指定的优惠券信息
func (s *CouponService) UpdateCoupon(id uint, input *UpdateCouponInput) (*models.Coupon, error) {
	// 使用事务确保更新的原子性
	tx := database.DB.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("开启事务失败: %w", tx.Error)
	}

	// 先查找优惠券
	var coupon models.Coupon
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&coupon, id).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("优惠券不存在")
		}
		return nil, fmt.Errorf("查找优惠券失败: %w", err)
	}

	// 使用 map 来构建需要更新的字段，避免零值问题
	updates := make(map[string]interface{})

	if input.CouponName != "" {
		updates["coupon_name"] = input.CouponName
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}
	if input.MinPurchaseAmount != nil {
		updates["min_purchase_amount"] = *input.MinPurchaseAmount
	}
	if input.TotalQuantity != nil {
		updates["total_quantity"] = *input.TotalQuantity
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}
	if input.UsageLimitPerUser != nil {
		updates["usage_limit_per_user"] = *input.UsageLimitPerUser
	}
	// 特别处理 store_id，因为它可能被设为 null
	if input.StoreID != nil {
		updates["store_id"] = input.StoreID
	}

	// 处理时间格式
	if input.StartTime != nil {
		startTime, err := time.ParseInLocation("2006-01-02 15:04:05", *input.StartTime, time.Local)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("无效的开始时间格式: %w", err)
		}
		updates["start_time"] = startTime
	}
	if input.EndTime != nil {
		endTime, err := time.ParseInLocation("2006-01-02 15:04:05", *input.EndTime, time.Local)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("无效的结束时间格式: %w", err)
		}
		updates["end_time"] = endTime
	}

	// 如果没有提供任何更新字段，则直接返回
	if len(updates) == 0 {
		tx.Rollback() // 虽然没有操作，但保持良好实践
		return &coupon, nil
	}

	// 执行更新
	if err := tx.Model(&coupon).Updates(updates).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("更新优惠券失败: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("提交事务失败: %w", err)
	}

	return &coupon, nil
}

// GetAvailableCouponsForUserInput 定义了查询用户可领取优惠券的输入
type GetAvailableCouponsForUserInput struct {
	UserID   string `form:"user_id" binding:"required"`
	StoreID  *uint  `form:"store_id"` // 可选，用于筛选特定门店的优惠券
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
}

// GetAvailableCouponsForUser 查询指定用户可领取的优惠券列表
func (s *CouponService) GetAvailableCouponsForUser(input *GetAvailableCouponsForUserInput) ([]models.Coupon, int64, error) {
	var availableCoupons []models.Coupon
	var total int64

	now := time.Now()

	// 基础查询条件：启用、在有效期内、有库存
	baseQuery := database.DB.Model(&models.Coupon{}).
		Where("status = 1 AND ? BETWEEN start_time AND end_time AND (total_quantity = 0 OR issued_quantity < total_quantity)", now)

	// 门店筛选条件：全平台通用券 或 特定门店券
	if input.StoreID != nil {
		baseQuery = baseQuery.Where("store_id IS NULL OR store_id = ?", *input.StoreID)
	} else {
		// 如果不提供 store_id，通常只返回全平台通用券
		baseQuery = baseQuery.Where("store_id IS NULL")
	}

	// 核心逻辑：过滤掉用户已达领取上限的优惠券
	// 使用子查询来计算用户已领取的数量
	subQuery := "usage_limit_per_user = 0 OR (SELECT count(*) FROM coupon_log WHERE coupon_log.coupon_id = coupon.coupon_id AND coupon_log.user_union_id = ? AND coupon_log.action_type = 'RECEIVE' AND coupon_log.status = 1) < coupon.usage_limit_per_user"
	finalQuery := baseQuery.Where(subQuery, input.UserID)

	// 计算总数
	if err := finalQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计可领取优惠券数量失败: %w", err)
	}

	// 分页
	if input.Page > 0 && input.PageSize > 0 {
		offset := (input.Page - 1) * input.PageSize
		finalQuery = finalQuery.Offset(offset).Limit(input.PageSize)
	}

	// 执行查询
	if err := finalQuery.Find(&availableCoupons).Error; err != nil {
		return nil, 0, fmt.Errorf("查询可领取优惠券列表失败: %w", err)
	}

	return availableCoupons, total, nil
}

// DeleteCoupon 软删除一张优惠券
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
