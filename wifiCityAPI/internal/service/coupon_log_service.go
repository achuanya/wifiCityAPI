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

// CouponLogService 提供了优惠券日志相关的业务逻辑
type CouponLogService struct{}

// LogActionInput 定义了记录优惠券日志的通用输入
type LogActionInput struct {
	CouponID       uint     `json:"coupon_id" binding:"required"`
	UserUnionID    string   `json:"user_union_id" binding:"required"`
	StoreID        *uint    `json:"store_id"`
	ActionType     string   `json:"action_type" binding:"required,oneof=ISSUE RECEIVE USE EXPIRE REFUND"`
	OrderID        *string  `json:"order_id"`
	AmountDeducted *float64 `json:"amount_deducted"`
	Remark         string   `json:"remark"`
}

// CreateCouponLog 创建优惠券日志，并根据操作类型执行特定逻辑
func (s *CouponLogService) CreateCouponLog(input *LogActionInput) (*models.CouponLog, error) {
	// 对于"领取"操作，需要执行特殊逻辑并使用事务
	if input.ActionType == "RECEIVE" {
		return s.receiveCoupon(input)
	}

	// 对于其他操作类型 (ISSUE, USE, EXPIRE, REFUND)，暂时只记录日志
	log := models.CouponLog{
		CouponID:       input.CouponID,
		UserUnionID:    input.UserUnionID,
		StoreID:        input.StoreID,
		ActionType:     input.ActionType,
		ActionTime:     time.Now(), // 记录当前时间
		OrderID:        "",         // 可根据需要从 input 赋值
		AmountDeducted: 0,          // 可根据需要从 input 赋值
		Status:         1,          // 默认为成功
		Remark:         input.Remark,
	}
	if input.OrderID != nil {
		log.OrderID = *input.OrderID
	}
	if input.AmountDeducted != nil {
		log.AmountDeducted = *input.AmountDeducted
	}

	if err := database.DB.Create(&log).Error; err != nil {
		return nil, fmt.Errorf("创建优惠券日志失败: %w", err)
	}

	return &log, nil
}

// receiveCoupon 处理用户领取优惠券的逻辑
func (s *CouponLogService) receiveCoupon(input *LogActionInput) (*models.CouponLog, error) {
	var log *models.CouponLog

	// 使用 GORM 的事务来确保数据一致性
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var coupon models.Coupon
		var userLogCount int64

		// 1. 锁定并查找优惠券信息
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&coupon, input.CouponID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("优惠券不存在")
			}
			return fmt.Errorf("查询优惠券失败: %w", err)
		}

		// 2. 校验优惠券状态和有效期
		if coupon.Status != 1 {
			return errors.New("优惠券已禁用")
		}
		now := time.Now()
		if now.Before(coupon.StartTime) || now.After(coupon.EndTime) {
			return errors.New("优惠券不在可用时间内")
		}

		// 3. 校验库存
		if coupon.TotalQuantity > 0 && coupon.IssuedQuantity >= coupon.TotalQuantity {
			return errors.New("优惠券已领完")
		}

		// 4. 校验用户领取限制
		if coupon.UsageLimitPerUser > 0 {
			tx.Model(&models.CouponLog{}).
				Where("user_union_id = ? AND coupon_id = ? AND action_type = 'RECEIVE' AND status = 1", input.UserUnionID, input.CouponID).
				Count(&userLogCount)
			if userLogCount >= int64(coupon.UsageLimitPerUser) {
				return errors.New("已达到该优惠券的领取上限")
			}
		}

		// 5. 创建领取日志
		log = &models.CouponLog{
			CouponID:    input.CouponID,
			UserUnionID: input.UserUnionID,
			StoreID:     input.StoreID,
			ActionType:  "RECEIVE",
			ActionTime:  now,
			Status:      1,
			Remark:      "用户成功领取",
		}
		if err := tx.Create(log).Error; err != nil {
			return fmt.Errorf("创建领取日志失败: %w", err)
		}

		// 6. 更新优惠券已发行数量
		coupon.IssuedQuantity++
		if err := tx.Save(&coupon).Error; err != nil {
			return fmt.Errorf("更新优惠券数量失败: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return log, nil
}

// GetCouponLogsInput 定义了查询优惠券日志的输入
type GetCouponLogsInput struct {
	UserUnionID *string `form:"user_union_id"`
	CouponID    *uint   `form:"coupon_id"`
	StoreID     *uint   `form:"store_id"`
	ActionType  *string `form:"action_type"`
	Page        int     `form:"page"`
	PageSize    int     `form:"pageSize"`
}

// GetCouponLogs 根据条件查询优惠券日志
func (s *CouponLogService) GetCouponLogs(input *GetCouponLogsInput) ([]models.CouponLog, int64, error) {
	db := database.DB.WithContext(context.Background())
	query := db.Model(&models.CouponLog{})

	if input.UserUnionID != nil && *input.UserUnionID != "" {
		query = query.Where("user_union_id = ?", *input.UserUnionID)
	}
	if input.CouponID != nil {
		query = query.Where("coupon_id = ?", *input.CouponID)
	}
	if input.StoreID != nil {
		query = query.Where("store_id = ?", *input.StoreID)
	}
	if input.ActionType != nil && *input.ActionType != "" {
		query = query.Where("action_type = ?", *input.ActionType)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var logs []models.CouponLog
	if input.Page > 0 && input.PageSize > 0 {
		offset := (input.Page - 1) * input.PageSize
		query = query.Offset(offset).Limit(input.PageSize)
	}

	if err := query.Order("action_time DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
