package v1

import (
	"app/internal/service"
	"app/pkg/security"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// StatsHandler 负责处理统计相关的API请求
type StatsHandler struct {
	service *service.StatsService
}

// NewStatsHandler 创建一个新的 StatsHandler
func NewStatsHandler() *StatsHandler {
	return &StatsHandler{
		service: &service.StatsService{},
	}
}

// GetStoreStats godoc
// @Summary      获取门店统计数据
// @Description  获取门店总数、按省份和城市分布的统计信息
// @Tags         Statistics
// @Produce      json
// @Success      200  {object}  security.EncryptedData
// @Failure      500  {object}  security.EncryptedData
// @Router       /stats/stores [get]
func (h *StatsHandler) GetStoreStats(c *gin.Context) {
	stats, err := h.service.GetStoreStats()
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		return
	}
	security.SendEncryptedResponse(c, http.StatusOK, stats)
}

// GetWifiUsageStats godoc
// @Summary 获取WIFI使用统计
// @Description 获取WIFI使用情况的统计数据，可按门店ID筛选
// @Tags stats
// @Accept  json
// @Produce  json
// @Param store_id query int false "门店ID"
// @Success 200 {object} object "成功响应，返回多种统计数据"
// @Failure 500 {object} security.ErrorResponse "服务器内部错误"
// @Router /api/v1/stats/wifi-usage [get]
func (h *StatsHandler) GetWifiUsageStats(c *gin.Context) {
	var storeID *uint
	if idStr := c.Query("store_id"); idStr != "" {
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err == nil {
			u_id := uint(id)
			storeID = &u_id
		}
	}

	stats, err := h.service.GetWifiUsageStats(storeID)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		return
	}
	security.SendEncryptedResponse(c, http.StatusOK, stats)
}

// GetUserBehaviorStats godoc
// @Summary 获取用户行为统计
// @Description 获取用户行为的统计数据，可按日期范围筛选
// @Tags stats
// @Accept  json
// @Produce  json
// @Param start_date query string false "开始日期 (格式: YYYY-MM-DD)"
// @Param end_date query string false "结束日期 (格式: YYYY-MM-DD)"
// @Success 200 {object} object "成功响应，返回多种统计数据"
// @Failure 400 {object} security.ErrorResponse "请求参数错误"
// @Failure 500 {object} security.ErrorResponse "服务器内部错误"
// @Router /api/v1/stats/user-behavior [get]
func (h *StatsHandler) GetUserBehaviorStats(c *gin.Context) {
	var input service.GetUserBehaviorStatsInput
	if err := c.ShouldBindQuery(&input); err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: err.Error()})
		return
	}

	stats, err := h.service.GetUserBehaviorStats(&input)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		return
	}
	security.SendEncryptedResponse(c, http.StatusOK, stats)
}

// GetCouponStats godoc
// @Summary 获取优惠券统计
// @Description 获取优惠券相关的统计数据，可按日期范围和门店ID筛选
// @Tags stats
// @Accept  json
// @Produce  json
// @Param start_date query string false "开始日期 (格式: YYYY-MM-DD)"
// @Param end_date query string false "结束日期 (格式: YYYY-MM-DD)"
// @Param store_id query int false "门店ID"
// @Success 200 {object} object "成功响应，返回多种统计数据"
// @Failure 400 {object} security.ErrorResponse "请求参数错误"
// @Failure 500 {object} security.ErrorResponse "服务器内部错误"
// @Router /api/v1/stats/coupons [get]
func (h *StatsHandler) GetCouponStats(c *gin.Context) {
	var input service.GetCouponStatsInput
	if err := c.ShouldBindQuery(&input); err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: err.Error()})
		return
	}
	stats, err := h.service.GetCouponStats(&input)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		return
	}
	security.SendEncryptedResponse(c, http.StatusOK, stats)
}

// GetPopularWifi godoc
// @Summary 最受欢迎WIFI统计
// @Description 获取最受欢迎的WIFI配置统计，可按门店、时间段筛选
// @Tags Stats
// @Accept  json
// @Produce  json
// @Param store_id query int false "门店ID"
// @Param start_date query string false "开始日期（格式：YYYY-MM-DD）"
// @Param end_date query string false "结束日期（格式：YYYY-MM-DD）"
// @Param limit query int false "返回记录数量（默认10）"
// @Success 200 {object} []service.WifiPopularityItem
// @Failure 400 {object} security.ErrorResponse
// @Failure 500 {object} security.ErrorResponse
// @Router /stats/popular-wifi [get]
func (h *StatsHandler) GetPopularWifi(c *gin.Context) {
	var input service.GetPopularWifiInput
	if err := c.ShouldBindQuery(&input); err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: "无效的查询参数: " + err.Error()})
		return
	}

	stats, err := h.service.GetPopularWifi(&input)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		return
	}

	security.SendEncryptedResponse(c, http.StatusOK, stats)
}

// GetScanTimeDistribution godoc
// @Summary 扫码时段分布统计
// @Description 获取扫码时段的分布统计，可按门店、时间段筛选
// @Tags Stats
// @Accept  json
// @Produce  json
// @Param store_id query int false "门店ID"
// @Param start_date query string false "开始日期（格式：YYYY-MM-DD）"
// @Param end_date query string false "结束日期（格式：YYYY-MM-DD）"
// @Success 200 {object} []service.HourlyDistribution
// @Failure 400 {object} security.ErrorResponse
// @Failure 500 {object} security.ErrorResponse
// @Router /stats/scan-time-distribution [get]
func (h *StatsHandler) GetScanTimeDistribution(c *gin.Context) {
	var input service.GetScanTimeDistributionInput
	if err := c.ShouldBindQuery(&input); err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: "无效的查询参数: " + err.Error()})
		return
	}

	stats, err := h.service.GetScanTimeDistribution(&input)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		return
	}

	security.SendEncryptedResponse(c, http.StatusOK, stats)
}
