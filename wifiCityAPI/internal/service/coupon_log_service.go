package service

import (
	"context"

	"github.com/gin-gonic/gin/wifiCityAPI/internal/models"
	"github.com/gin-gonic/gin/wifiCityAPI/pkg/database"
	"gorm.io/gorm"
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

// CreateCouponLog 在数据库中创建一条新的优惠券日志。
// 它在一个事务中完成此操作，并在某些情况下（如"使用"和"退款"）更新优惠券的已发行数量。
func (s *CouponLogService) CreateCouponLog(input *LogActionInput) (*models.CouponLog, error) {
	log := models.CouponLog{
		CouponID:    input.CouponID,
		UserUnionID: input.UserUnionID,
		StoreID:     input.StoreID,
		ActionType:  input.ActionType,
		Remark:      input.Remark,
		Status:      1, // 默认成功
	}
	if input.OrderID != nil {
		log.OrderID = *input.OrderID
	}
	if input.AmountDeducted != nil {
		log.AmountDeducted = *input.AmountDeducted
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 创建日志记录
		if err := tx.Create(&log).Error; err != nil {
			return err
		}

		// 2. 根据行为类型，可能需要更新 `coupon` 表的统计数据
		switch input.ActionType {
		case "RECEIVE":
			// 每领取一张，已发行数量 +1
			return tx.Model(&models.Coupon{}).Where("coupon_id = ?", input.CouponID).
				UpdateColumn("issued_quantity", gorm.Expr("issued_quantity + ?", 1)).Error
		case "USE":
			// "使用"通常意味着已经领取，所以这里不增减 issued_quantity。
			// 具体的核销逻辑可能更复杂，取决于业务需求。
			// 这里仅作示例。
		case "REFUND":
			// 如果是"退券"，已发行数量 -1
			return tx.Model(&models.Coupon{}).Where("coupon_id = ?", input.CouponID).
				UpdateColumn("issued_quantity", gorm.Expr("issued_quantity - ?", 1)).Error
		}

		// 对于其他 action_type，不执行额外操作
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &log, nil
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
