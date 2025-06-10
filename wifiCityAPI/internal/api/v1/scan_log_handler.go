package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/wifiCityAPI/internal/service"
	"github.com/gin-gonic/gin/wifiCityAPI/pkg/security"
	"gorm.io/gorm"
)

// ScanLogHandler 负责处理扫码日志相关的API请求
type ScanLogHandler struct {
	service *service.ScanLogService
}

// NewScanLogHandler 创建一个新的 ScanLogHandler
func NewScanLogHandler() *ScanLogHandler {
	return &ScanLogHandler{
		service: &service.ScanLogService{},
	}
}

// CreateScanLog
// @Summary 记录用户扫码连接日志
// @Accept json
// @Produce json
// @Param log body service.CreateScanLogInput true "扫码日志信息"
// @Success 201 {object} models.ScanLog
// @Router /scan-logs [post]
func (h *ScanLogHandler) CreateScanLog(c *gin.Context) {
	var input service.CreateScanLogInput
	if err := c.ShouldBindJSON(&input); err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: err.Error()})
		return
	}
	// 补充IP地址
	input.IPAddress = c.ClientIP()

	logEntry, err := h.service.CreateScanLog(&input)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		return
	}

	security.SendEncryptedResponse(c, http.StatusCreated, logEntry)
}

// GetScanLogs
// @Summary 查询扫码日志列表
// @Description 分页获取扫码日志
// @Tags scan-logs
// @Accept  json
// @Produce  json
// @Param store_id query int false "门店ID"
// @Param user_union_id query string false "用户UnionID"
// @Param success_flag query boolean false "是否成功连接"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} object{logs=[]models.ScanLog, total=int64} "成功响应"
// @Failure 400 {object} security.ErrorResponse "请求参数错误"
// @Failure 500 {object} security.ErrorResponse "服务器内部错误"
// @Router /api/v1/scan-logs [get]
func (h *ScanLogHandler) GetScanLogs(c *gin.Context) {
	var input service.GetScanLogsInput
	if err := c.ShouldBindQuery(&input); err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: err.Error()})
		return
	}

	logs, total, err := h.service.GetScanLogs(&input)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		return
	}

	security.SendEncryptedResponse(c, http.StatusOK, gin.H{
		"logs":  logs,
		"total": total,
	})
}

// UpdateScanLogResult
// @Summary 更新扫码日志连接结果
// @Description 更新指定扫码日志的WIFI连接结果
// @Tags scan-logs
// @Accept  json
// @Produce  json
// @Param id path int true "扫码日志ID"
// @Param result body service.UpdateScanLogResultInput true "连接结果"
// @Success 200 {object} models.ScanLog "成功响应"
// @Failure 400 {object} security.ErrorResponse "请求参数错误"
// @Failure 404 {object} security.ErrorResponse "日志未找到"
// @Failure 500 {object} security.ErrorResponse "服务器内部错误"
// @Router /api/v1/scan-logs/{id}/result [patch]
func (h *ScanLogHandler) UpdateScanLogResult(c *gin.Context) {
	logId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: "无效的日志ID"})
		return
	}

	var input service.UpdateScanLogResultInput
	if err := c.ShouldBindJSON(&input); err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: err.Error()})
		return
	}

	err = h.service.UpdateScanLogResult(logId, &input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			security.SendEncryptedResponse(c, http.StatusNotFound, security.ErrorResponse{Error: "扫码日志未找到"})
		} else {
			security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		}
		return
	}

	security.SendEncryptedResponse(c, http.StatusOK, gin.H{"message": "更新成功"})
}

// GetDailyScanCountByStore
// @Summary 查询门店的每日扫码量
// @Description 获取指定门店过去N天的每日扫码统计
// @Tags stores
// @Accept  json
// @Produce  json
// @Param storeId path int true "门店ID"
// @Param days query int false "查询天数" default(7)
// @Success 200 {array} service.DailyScanCountResult "成功响应"
// @Failure 400 {object} security.ErrorResponse "请求参数错误"
// @Failure 500 {object} security.ErrorResponse "服务器内部错误"
// @Router /api/v1/stores/{storeId}/scans/daily-count [get]
func (h *ScanLogHandler) GetDailyScanCountByStore(c *gin.Context) {
	storeId, err := strconv.ParseUint(c.Param("storeId"), 10, 32)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: "无效的门店ID"})
		return
	}

	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))

	stats, err := h.service.GetDailyScanCountByStore(uint(storeId), days)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		return
	}

	security.SendEncryptedResponse(c, http.StatusOK, stats)
}
