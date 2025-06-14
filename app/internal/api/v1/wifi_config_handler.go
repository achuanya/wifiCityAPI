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

// WifiConfigHandler 负责处理WIFI配置相关的API请求
type WifiConfigHandler struct {
	service *service.WifiConfigService
}

// NewWifiConfigHandler 创建一个新的 WifiConfigHandler
func NewWifiConfigHandler() *WifiConfigHandler {
	return &WifiConfigHandler{
		service: &service.WifiConfigService{},
	}
}

// CreateWifiConfig
// @Summary 新增WIFI配置
// @Accept json
// @Produce json
// @Param wifi_config body service.CreateWifiConfigInput true "WIFI配置信息"
// @Success 201 {object} models.WifiConfig
// @Router /wifis [post]
func (h *WifiConfigHandler) CreateWifiConfig(c *gin.Context) {
	var input service.CreateWifiConfigInput
	if err := c.ShouldBindJSON(&input); err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: err.Error()})
		return
	}

	wifiConfig, err := h.service.CreateWifiConfig(&input)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		return
	}

	security.SendEncryptedResponse(c, http.StatusCreated, wifiConfig)
}

// GetWifiConfig
// @Summary 查询单个WIFI配置详情
// @Produce json
// @Param id path int true "WIFI配置ID"
// @Success 200 {object} models.WifiConfig
// @Router /wifis/{id} [get]
func (h *WifiConfigHandler) GetWifiConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: "无效的WIFI配置ID"})
		return
	}

	wifiConfig, err := h.service.GetWifiConfigByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			security.SendEncryptedResponse(c, http.StatusNotFound, security.ErrorResponse{Error: "WIFI配置未找到"})
		} else {
			security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		}
		return
	}

	security.SendEncryptedResponse(c, http.StatusOK, wifiConfig)
}

// GetWifiConfigsByStore
// @Summary 查询门店所有WIFI配置列表
// @Produce json
// @Param storeId path int true "门店ID"
// @Success 200 {array} models.WifiConfig
// @Router /stores/{storeId}/wifis [get]
func (h *WifiConfigHandler) GetWifiConfigsByStore(c *gin.Context) {
	storeId, err := strconv.ParseUint(c.Param("storeId"), 10, 32)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: "无效的门店ID"})
		return
	}

	wifiConfigs, err := h.service.GetWifiConfigsByStoreID(uint(storeId))
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		return
	}

	security.SendEncryptedResponse(c, http.StatusOK, wifiConfigs)
}

// UpdateWifiConfig
// @Summary 更新WIFI配置
// @Accept json
// @Produce json
// @Param id path int true "WIFI配置ID"
// @Param wifi_config body service.UpdateWifiConfigInput true "要更新的WIFI配置信息"
// @Success 200 {object} models.WifiConfig
// @Router /wifis/{id} [put]
func (h *WifiConfigHandler) UpdateWifiConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: "无效的WIFI配置ID"})
		return
	}

	var input service.UpdateWifiConfigInput
	if err := c.ShouldBindJSON(&input); err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: err.Error()})
		return
	}

	wifiConfig, err := h.service.UpdateWifiConfig(uint(id), &input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			security.SendEncryptedResponse(c, http.StatusNotFound, security.ErrorResponse{Error: "WIFI配置未找到"})
		} else {
			security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		}
		return
	}

	security.SendEncryptedResponse(c, http.StatusOK, wifiConfig)
}

// DeleteWifiConfig
// @Summary 删除WIFI配置
// @Produce json
// @Param id path int true "WIFI配置ID"
// @Success 204
// @Router /wifis/{id} [delete]
func (h *WifiConfigHandler) DeleteWifiConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: "无效的WIFI配置ID"})
		return
	}

	err = h.service.DeleteWifiConfig(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			security.SendEncryptedResponse(c, http.StatusNotFound, security.ErrorResponse{Error: "WIFI配置未找到"})
		} else {
			security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// CreateBatchWifiConfigs 批量创建WIFI配置
// @Summary 批量新增WIFI配置
// @Description 一次性为门店添加多个WIFI配置
// @Tags WifiConfigs
// @Accept  json
// @Produce  json
// @Param   configs body []service.CreateWifiConfigInput true "WIFI配置数组"
// @Success 201 {object} object{configs=[]models.WifiConfig}
// @Failure 400 {object} security.ErrorResponse
// @Failure 500 {object} security.ErrorResponse
// @Router /wifi-configs/batch [post]
func (h *WifiConfigHandler) CreateBatchWifiConfigs(c *gin.Context) {
	var inputs []*service.CreateWifiConfigInput
	if err := c.ShouldBindJSON(&inputs); err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, gin.H{"error": "无效的请求数据: " + err.Error()})
		return
	}

	if len(inputs) == 0 {
		security.SendEncryptedResponse(c, http.StatusBadRequest, gin.H{"error": "请求体不能为空数组"})
		return
	}

	createdConfigs, err := h.service.CreateBatchWifiConfigs(inputs)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusInternalServerError, gin.H{"error": "批量创建失败: " + err.Error()})
		return
	}

	security.SendEncryptedResponse(c, http.StatusCreated, gin.H{"configs": createdConfigs})
}

// DeleteBatchWifiConfigs godoc
// @Summary      批量删除WIFI配置
// @Description  批量删除多个WIFI配置
// @Tags         WiFiConfigs
// @Accept       json
// @Produce      json
// @Param        ids   body      []uint  true  "WIFI配置ID数组"
// @Success      204  {object}  nil
// @Failure      400  {object}  security.ErrorResponse
// @Failure      500  {object}  security.ErrorResponse
// @Router       /wifi-configs/batch [delete]
func (h *WifiConfigHandler) DeleteBatchWifiConfigs(c *gin.Context) {
	var ids []uint
	if err := c.ShouldBindJSON(&ids); err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: "无效的请求数据: " + err.Error()})
		return
	}

	if len(ids) == 0 {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: "请求体不能为空数组"})
		return
	}

	if err := h.service.DeleteBatchWifiConfigs(ids); err != nil {
		security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: "批量删除失败: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetWifiConfigsByStoreAndType godoc
// @Summary 查询门店特定类型的WIFI配置
// @Description 获取指定门店特定类型的WIFI配置列表
// @Tags WifiConfigs
// @Accept  json
// @Produce  json
// @Param store_id query int true "门店ID"
// @Param wifi_type query string true "WIFI类型"
// @Success 200 {array} models.WifiConfig
// @Failure 400 {object} security.ErrorResponse
// @Failure 500 {object} security.ErrorResponse
// @Router /wifi-configs/type [get]
func (h *WifiConfigHandler) GetWifiConfigsByStoreAndType(c *gin.Context) {
	var input service.GetWifiConfigsByStoreAndTypeInput
	if err := c.ShouldBindQuery(&input); err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: "无效的查询参数: " + err.Error()})
		return
	}

	if input.StoreID == 0 {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: "门店ID不能为空"})
		return
	}

	if input.WifiType == "" {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: "WIFI类型不能为空"})
		return
	}

	wifiConfigs, err := h.service.GetWifiConfigsByStoreAndType(&input)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		return
	}

	security.SendEncryptedResponse(c, http.StatusOK, wifiConfigs)
}
