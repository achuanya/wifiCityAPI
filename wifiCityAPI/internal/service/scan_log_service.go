package service

import (
	"context"
	"time"

	"github.com/gin-gonic/gin/wifiCityAPI/internal/models"
	"github.com/gin-gonic/gin/wifiCityAPI/pkg/database"
	"gorm.io/gorm"
)

// ScanLogService 提供了扫码日志相关的业务逻辑
type ScanLogService struct{}

// CreateScanLogInput 定义了记录扫码日志的输入
type CreateScanLogInput struct {
	StoreID            uint    `json:"store_id" binding:"required"`
	UserUnionID        string  `json:"user_union_id" binding:"required"`
	DeviceInfo         string  `json:"device_info"`
	IPAddress          string  `json:"ip_address"`
	NetworkType        string  `json:"network_type"`
	LocationLat        float64 `json:"location_lat"`
	LocationLng        float64 `json:"location_lng"`
	MiniProgramVersion string  `json:"mini_program_version"`
	QrCodeType         string  `json:"qr_code_type"`
	QrCodeID           string  `json:"qr_code_id"`
	SystemInfo         string  `json:"system_info"`
	Brand              string  `json:"brand"`
	Model              string  `json:"model"`
	PagePath           string  `json:"page_path"`
	Referer            string  `json:"referer"`
}

// CreateScanLog 创建一条新的扫码日志
func (s *ScanLogService) CreateScanLog(input *CreateScanLogInput) (*models.ScanLog, error) {
	log := models.ScanLog{
		StoreID:            input.StoreID,
		UserUnionID:        input.UserUnionID,
		DeviceInfo:         input.DeviceInfo,
		IPAddress:          input.IPAddress,
		NetworkType:        input.NetworkType,
		LocationLat:        input.LocationLat,
		LocationLng:        input.LocationLng,
		MiniProgramVersion: input.MiniProgramVersion,
		QrCodeType:         input.QrCodeType,
		QrCodeID:           input.QrCodeID,
		SystemInfo:         input.SystemInfo,
		Brand:              input.Brand,
		Model:              input.Model,
		PagePath:           input.PagePath,
		Referer:            input.Referer,
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&log).Error
	})
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// GetScanLogsInput 定义了查询扫码日志的输入
type GetScanLogsInput struct {
	StoreID     uint   `form:"store_id"`
	UserUnionID string `form:"user_union_id"`
	SuccessFlag *bool  `form:"success_flag"`
	Page        int    `form:"page"`
	PageSize    int    `form:"pageSize"`
}

// GetScanLogs 查询扫码日志列表（分页和过滤）
func (s *ScanLogService) GetScanLogs(input *GetScanLogsInput) ([]models.ScanLog, int64, error) {
	var logs []models.ScanLog
	var total int64

	db := database.DB.WithContext(context.Background()).Model(&models.ScanLog{})

	if input.StoreID != 0 {
		db = db.Where("store_id = ?", input.StoreID)
	}
	if input.UserUnionID != "" {
		db = db.Where("user_union_id = ?", input.UserUnionID)
	}
	if input.SuccessFlag != nil {
		db = db.Where("success_flag = ?", *input.SuccessFlag)
	}

	// 计算总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 处理分页
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 10
	}
	offset := (input.Page - 1) * input.PageSize

	// 查询数据
	err := db.Order("scan_time DESC").Offset(offset).Limit(input.PageSize).Find(&logs).Error
	return logs, total, err
}

// UpdateScanLogResultInput 定义了更新扫码日志结果的输入
type UpdateScanLogResultInput struct {
	SuccessFlag       bool   `json:"success_flag"`
	FailReasonCode    string `json:"fail_reason_code,omitempty"`
	FailReasonMessage string `json:"fail_reason_message,omitempty"`
	WifiSSID          string `json:"wifi_ssid,omitempty"`
	WifiMac           string `json:"wifi_mac,omitempty"`
	WifiSignal        int8   `json:"wifi_signal,omitempty"`
}

// UpdateScanLogResult 更新扫码日志的连接结果
func (s *ScanLogService) UpdateScanLogResult(logID uint64, input *UpdateScanLogResultInput) error {
	updateData := map[string]interface{}{
		"success_flag":        input.SuccessFlag,
		"fail_reason_code":    input.FailReasonCode,
		"fail_reason_message": input.FailReasonMessage,
		"wifi_ssid":           input.WifiSSID,
		"wifi_mac":            input.WifiMac,
		"wifi_signal":         input.WifiSignal,
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&models.ScanLog{}).Where("log_id = ?", logID).UpdateColumns(updateData)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

// DailyScanCountResult 定义了每日扫码量的返回结构
type DailyScanCountResult struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// GetDailyScanCountByStore 查询指定门店的每日扫码量
func (s *ScanLogService) GetDailyScanCountByStore(storeID uint, days int) ([]DailyScanCountResult, error) {
	var results []DailyScanCountResult
	if days <= 0 {
		days = 7 // 默认查询最近7天
	}
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	err := database.DB.WithContext(context.Background()).
		Model(&models.ScanLog{}).
		Select("DATE(scan_time) as date, COUNT(*) as count").
		Where("store_id = ? AND scan_time BETWEEN ? AND ?", storeID, startDate, endDate).
		Group("DATE(scan_time)").
		Order("date DESC").
		Scan(&results).Error

	return results, err
}
