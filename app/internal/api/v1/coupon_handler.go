package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"app/internal/service"
	"app/pkg/security"
)

// CouponHandler 负责处理优惠券相关的API请求
type CouponHandler struct {
	service *service.CouponService
}

// NewCouponHandler 创建一个新的 CouponHandler
func NewCouponHandler() *CouponHandler {
	return &CouponHandler{
		service: &service.CouponService{},
	}
}

// CreateCoupon
// @Summary 创建优惠券
// @Accept json
// @Produce json
// @Param coupon body service.CreateCouponInput true "优惠券信息"
// @Success 201 {object} models.Coupon
// @Router /coupons [post]
func (h *CouponHandler) CreateCoupon(c *gin.Context) {
	var input service.CreateCouponInput
	if err := c.ShouldBindJSON(&input); err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: err.Error()})
		return
	}

	coupon, err := h.service.CreateCoupon(&input)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		return
	}

	security.SendEncryptedResponse(c, http.StatusCreated, coupon)
}

// CreateBatchCoupons 批量创建优惠券
// @Summary 批量创建优惠券
// @Description 一次性创建多个优惠券
// @Tags Coupons
// @Accept  json
// @Produce  json
// @Param   coupons body []service.CreateCouponInput true "优惠券数组"
// @Success 201 {object} object{coupons=[]models.Coupon}
// @Failure 400 {object} security.ErrorResponse
// @Failure 500 {object} security.ErrorResponse
// @Router /coupons/batch [post]
func (h *CouponHandler) CreateBatchCoupons(c *gin.Context) {
	var inputs []*service.CreateCouponInput
	if err := c.ShouldBindJSON(&inputs); err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, gin.H{"error": "无效的请求数据: " + err.Error()})
		return
	}

	if len(inputs) == 0 {
		security.SendEncryptedResponse(c, http.StatusBadRequest, gin.H{"error": "请求体不能为空数组"})
		return
	}

	createdCoupons, err := h.service.CreateBatchCoupons(inputs)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusInternalServerError, gin.H{"error": "批量创建失败: " + err.Error()})
		return
	}

	security.SendEncryptedResponse(c, http.StatusCreated, gin.H{"coupons": createdCoupons})
}

// GetCoupon
// @Summary 查询优惠券详情
// @Produce json
// @Param id path int true "优惠券ID"
// @Success 200 {object} models.Coupon
// @Router /coupons/{id} [get]
func (h *CouponHandler) GetCoupon(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: "无效的优惠券ID"})
		return
	}

	coupon, err := h.service.GetCouponByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			security.SendEncryptedResponse(c, http.StatusNotFound, security.ErrorResponse{Error: "优惠券未找到"})
		} else {
			security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		}
		return
	}

	security.SendEncryptedResponse(c, http.StatusOK, coupon)
}

// GetCoupons
// @Summary 查询优惠券列表
// @Produce json
// @Param store_id query int false "适用门店ID"
// @Param status query int false "状态 (1:启用, 0:禁用)"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} gin.H{"coupons": []models.Coupon, "total": int64}
// @Router /coupons [get]
func (h *CouponHandler) GetCoupons(c *gin.Context) {
	var input service.GetCouponsInput
	if err := c.ShouldBindQuery(&input); err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: err.Error()})
		return
	}

	coupons, total, err := h.service.GetCoupons(&input)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		return
	}

	security.SendEncryptedResponse(c, http.StatusOK, gin.H{
		"coupons": coupons,
		"total":   total,
	})
}

// UpdateCoupon
// @Summary 更新优惠券
// @Accept json
// @Produce json
// @Param id path int true "优惠券ID"
// @Param coupon body service.UpdateCouponInput true "要更新的优惠券信息"
// @Success 200 {object} models.Coupon
// @Router /coupons/{id} [put]
func (h *CouponHandler) UpdateCoupon(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: "无效的优惠券ID"})
		return
	}

	var input service.UpdateCouponInput
	if err := c.ShouldBindJSON(&input); err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: err.Error()})
		return
	}

	coupon, err := h.service.UpdateCoupon(uint(id), &input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			security.SendEncryptedResponse(c, http.StatusNotFound, security.ErrorResponse{Error: "优惠券未找到"})
		} else {
			security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		}
		return
	}

	security.SendEncryptedResponse(c, http.StatusOK, coupon)
}

// DeleteCoupon
// @Summary 删除优惠券（软删除）
// @Produce json
// @Param id path int true "优惠券ID"
// @Success 204
// @Router /coupons/{id} [delete]
func (h *CouponHandler) DeleteCoupon(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: "无效的优惠券ID"})
		return
	}

	err = h.service.DeleteCoupon(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			security.SendEncryptedResponse(c, http.StatusNotFound, security.ErrorResponse{Error: "优惠券未找到"})
		} else {
			security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAvailableCouponsForUser 获取用户可领取的优惠券列表
// @Summary 获取用户可领取的优惠券
// @Description 查询指定用户可领取的优惠券列表，支持按门店筛选
// @Tags Coupons
// @Accept  json
// @Produce  json
// @Param user_id query string true "用户ID"
// @Param store_id query int false "门店ID"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} object{coupons=[]models.Coupon, total=int64}
// @Failure 400 {object} security.ErrorResponse
// @Failure 500 {object} security.ErrorResponse
// @Router /coupons/available-for-user [get]
func (h *CouponHandler) GetAvailableCouponsForUser(c *gin.Context) {
	var input service.GetAvailableCouponsForUserInput
	if err := c.ShouldBindQuery(&input); err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, gin.H{"error": "无效的查询参数: " + err.Error()})
		return
	}

	coupons, total, err := h.service.GetAvailableCouponsForUser(&input)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusInternalServerError, gin.H{"error": "查询可领取优惠券失败: " + err.Error()})
		return
	}

	security.SendEncryptedResponse(c, http.StatusOK, gin.H{
		"coupons": coupons,
		"total":   total,
	})
}
