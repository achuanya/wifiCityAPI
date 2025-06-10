package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/wifiCityAPI/internal/service"
	"github.com/gin-gonic/gin/wifiCityAPI/pkg/security"
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
