package v1

import (
	"app/internal/service"
	"app/pkg/security"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CouponLogHandler 负责处理优惠券日志相关的API请求
type CouponLogHandler struct {
	service *service.CouponLogService
}

// NewCouponLogHandler 创建一个新的 CouponLogHandler
func NewCouponLogHandler() *CouponLogHandler {
	return &CouponLogHandler{
		service: &service.CouponLogService{},
	}
}

// CreateCouponLog godoc
// @Summary      记录优惠券行为日志
// @Description  用于记录用户领取、使用、核销优惠券等行为
// @Tags         CouponLogs
// @Accept       json
// @Produce      json
// @Param        log   body      service.LogActionInput  true  "日志信息"
// @Success      201  {object}  security.EncryptedData
// @Failure      400  {object}  security.EncryptedData
// @Failure      500  {object}  security.EncryptedData
// @Router       /coupon-logs [post]
func (h *CouponLogHandler) CreateCouponLog(c *gin.Context) {
	var input service.LogActionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: err.Error()})
		return
	}

	logEntry, err := h.service.CreateCouponLog(&input)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		return
	}

	security.SendEncryptedResponse(c, http.StatusCreated, logEntry)
}

// GetCouponLogs godoc
// @Summary 查询优惠券日志
// @Description 分页获取优惠券的发放、领取、使用等日志
// @Tags coupon-logs
// @Accept  json
// @Produce  json
// @Param user_union_id query string false "用户UnionID"
// @Param coupon_id query int false "优惠券ID"
// @Param store_id query int false "门店ID"
// @Param action_type query string false "行为类型 (ISSUE, RECEIVE, USE, EXPIRE, REFUND)"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} object{logs=[]models.CouponLog, total=int64} "成功响应"
// @Failure 400 {object} security.ErrorResponse "请求参数错误"
// @Failure 500 {object} security.ErrorResponse "服务器内部错误"
// @Router /api/v1/coupon-logs [get]
func (h *CouponLogHandler) GetCouponLogs(c *gin.Context) {
	var input service.GetCouponLogsInput
	if err := c.ShouldBindQuery(&input); err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: err.Error()})
		return
	}

	logs, total, err := h.service.GetCouponLogs(&input)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		return
	}

	security.SendEncryptedResponse(c, http.StatusOK, gin.H{
		"logs":  logs,
		"total": total,
	})
}
