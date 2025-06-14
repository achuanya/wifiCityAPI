package v1

import (
	"app/internal/service"
	"app/pkg/security"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StoreHandler 负责处理门店相关的API请求
type StoreHandler struct {
	service *service.StoreService
}

// NewStoreHandler 创建一个新的 StoreHandler
func NewStoreHandler() *StoreHandler {
	return &StoreHandler{
		service: &service.StoreService{},
	}
}

// CreateStore
// @Summary 新增门店
// @Accept json
// @Produce json
// @Param store body service.CreateStoreInput true "门店信息"
// @Success 201 {object} models.Store
// @Router /stores [post]
func (h *StoreHandler) CreateStore(c *gin.Context) {
	var input service.CreateStoreInput
	if err := c.ShouldBindJSON(&input); err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: err.Error()})
		return
	}

	store, err := h.service.CreateStore(&input)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		return
	}

	security.SendEncryptedResponse(c, http.StatusCreated, store)
}

// CreateStoreWithWifi 同时创建门店和WIFI配置
// @Summary 新增门店及WIFI
// @Description 在一个事务中同时创建门店及其初始WIFI配置
// @Tags Stores
// @Accept  json
// @Produce  json
// @Param   store_with_wifi body service.CreateStoreWithWifiInput true "门店及WIFI配置信息"
// @Success 201 {object} models.Store
// @Failure 400 {object} security.ErrorResponse
// @Failure 500 {object} security.ErrorResponse
// @Router /stores/with-wifi [post]
func (h *StoreHandler) CreateStoreWithWifi(c *gin.Context) {
	var input service.CreateStoreWithWifiInput
	if err := c.ShouldBindJSON(&input); err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, gin.H{"error": "无效的请求数据: " + err.Error()})
		return
	}

	store, err := h.service.CreateStoreWithWifi(&input)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusInternalServerError, gin.H{"error": "创建失败: " + err.Error()})
		return
	}

	security.SendEncryptedResponse(c, http.StatusCreated, store)
}

// GetStore
// @Summary 查询门店详情
// @Produce json
// @Param id path int true "门店ID"
// @Success 200 {object} models.Store
// @Router /stores/{id} [get]
func (h *StoreHandler) GetStore(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("storeId"), 10, 32)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: "无效的门店ID格式"})
		return
	}

	store, err := h.service.GetStoreByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			security.SendEncryptedResponse(c, http.StatusNotFound, security.ErrorResponse{Error: "门店未找到"})
		} else {
			security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		}
		return
	}

	security.SendEncryptedResponse(c, http.StatusOK, store)
}

// GetStores
// @Summary 查询门店列表
// @Description 支持分页、按区域筛选（province, city, district）和查询附近门店（lat, lng, radius）
// @Tags Stores
// @Accept  json
// @Produce  json
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param province query string false "省份"
// @Param city query string false "城市"
// @Param district query string false "区/县"
// @Param lat query number false "纬度"
// @Param lng query number false "经度"
// @Param radius query number false "半径（公里）"
// @Success 200 {object} object{stores=[]models.Store, total=int64}
// @Failure 400 {object} security.ErrorResponse
// @Failure 500 {object} security.ErrorResponse
// @Router /stores [get]
func (h *StoreHandler) GetStores(c *gin.Context) {
	var input service.GetStoresInput
	if err := c.ShouldBindQuery(&input); err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, gin.H{"error": "无效的查询参数: " + err.Error()})
		return
	}

	stores, total, err := h.service.GetStores(&input)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusInternalServerError, gin.H{"error": "获取门店列表失败: " + err.Error()})
		return
	}

	security.SendEncryptedResponse(c, http.StatusOK, gin.H{
		"stores": stores,
		"total":  total,
	})
}

// UpdateStore
// @Summary 更新门店信息
// @Accept json
// @Produce json
// @Param id path int true "门店ID"
// @Param store body service.UpdateStoreInput true "要更新的门店信息"
// @Success 200 {object} models.Store
// @Router /stores/{id} [put]
func (h *StoreHandler) UpdateStore(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("storeId"), 10, 32)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: "无效的门店ID格式"})
		return
	}

	var input service.UpdateStoreInput
	if err := c.ShouldBindJSON(&input); err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: err.Error()})
		return
	}

	store, err := h.service.UpdateStore(uint(id), &input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			security.SendEncryptedResponse(c, http.StatusNotFound, security.ErrorResponse{Error: "门店未找到"})
		} else {
			security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		}
		return
	}

	security.SendEncryptedResponse(c, http.StatusOK, store)
}

// DeleteStore
// @Summary 删除门店
// @Produce json
// @Param id path int true "门店ID"
// @Success 204
// @Router /stores/{id} [delete]
func (h *StoreHandler) DeleteStore(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("storeId"), 10, 32)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: "无效的门店ID格式"})
		return
	}

	err = h.service.DeleteStore(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			security.SendEncryptedResponse(c, http.StatusNotFound, security.ErrorResponse{Error: "门店未找到"})
		} else {
			security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateStoreStatus godoc
// @Summary 更新门店状态
// @Description 仅更新门店的状态（启用/停用）
// @Tags Stores
// @Accept  json
// @Produce  json
// @Param storeId path int true "门店ID"
// @Param status body object{status=int8} true "门店状态（1:启用, 0:停用）"
// @Success 200 {object} models.Store
// @Failure 400 {object} security.ErrorResponse
// @Failure 404 {object} security.ErrorResponse
// @Failure 500 {object} security.ErrorResponse
// @Router /stores/{storeId}/status [patch]
func (h *StoreHandler) UpdateStoreStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("storeId"), 10, 32)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: "无效的门店ID格式"})
		return
	}

	var input struct {
		Status int8 `json:"status" binding:"required,oneof=0 1"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: "无效的状态值，应为0或1"})
		return
	}

	store, err := h.service.UpdateStoreStatus(uint(id), input.Status)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			security.SendEncryptedResponse(c, http.StatusNotFound, security.ErrorResponse{Error: "门店未找到"})
		} else {
			security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		}
		return
	}

	security.SendEncryptedResponse(c, http.StatusOK, store)
}

// UpdateStorePhone godoc
// @Summary 更新门店联系电话
// @Description 仅更新门店的联系电话
// @Tags Stores
// @Accept  json
// @Produce  json
// @Param storeId path int true "门店ID"
// @Param phone body object{phone=string} true "门店联系电话"
// @Success 200 {object} models.Store
// @Failure 400 {object} security.ErrorResponse
// @Failure 404 {object} security.ErrorResponse
// @Failure 500 {object} security.ErrorResponse
// @Router /stores/{storeId}/phone [patch]
func (h *StoreHandler) UpdateStorePhone(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("storeId"), 10, 32)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: "无效的门店ID格式"})
		return
	}

	var input struct {
		Phone string `json:"phone" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: "电话号码不能为空"})
		return
	}

	store, err := h.service.UpdateStorePhone(uint(id), input.Phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			security.SendEncryptedResponse(c, http.StatusNotFound, security.ErrorResponse{Error: "门店未找到"})
		} else {
			security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		}
		return
	}

	security.SendEncryptedResponse(c, http.StatusOK, store)
}

// UpdateStoreLocation godoc
// @Summary 更新门店地理位置
// @Description 更新门店的地理坐标和地址信息
// @Tags Stores
// @Accept  json
// @Produce  json
// @Param storeId path int true "门店ID"
// @Param location body service.UpdateStoreLocationInput true "门店地理位置信息"
// @Success 200 {object} models.Store
// @Failure 400 {object} security.ErrorResponse
// @Failure 404 {object} security.ErrorResponse
// @Failure 500 {object} security.ErrorResponse
// @Router /stores/{storeId}/location [patch]
func (h *StoreHandler) UpdateStoreLocation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("storeId"), 10, 32)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: "无效的门店ID格式"})
		return
	}

	var input service.UpdateStoreLocationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: "无效的位置信息: " + err.Error()})
		return
	}

	store, err := h.service.UpdateStoreLocation(uint(id), &input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			security.SendEncryptedResponse(c, http.StatusNotFound, security.ErrorResponse{Error: "门店未找到"})
		} else {
			security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		}
		return
	}

	security.SendEncryptedResponse(c, http.StatusOK, store)
}
