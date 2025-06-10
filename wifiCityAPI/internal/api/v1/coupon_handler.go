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
