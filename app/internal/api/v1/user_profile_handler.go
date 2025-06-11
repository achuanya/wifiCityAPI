package v1

import (
	"app/internal/models"
	"app/internal/service"
	"app/pkg/security"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserProfileHandler 负责处理用户相关的API请求
type UserProfileHandler struct {
	service *service.UserProfileService
}

// NewUserProfileHandler 创建一个新的 UserProfileHandler
func NewUserProfileHandler() *UserProfileHandler {
	return &UserProfileHandler{
		service: &service.UserProfileService{},
	}
}

// CreateOrUpdateUser
// @Summary 创建或更新用户档案
// @Description 根据UnionID创建或更新用户。如果用户已存在，则更新信息；否则创建新用户。
// @Accept json
// @Produce json
// @Param user body service.CreateOrUpdateUserInput true "用户信息"
// @Success 200 {object} models.UserProfile
// @Success 201 {object} models.UserProfile
// @Router /users [post]
func (h *UserProfileHandler) CreateOrUpdateUser(c *gin.Context) {
	var input service.CreateOrUpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.service.CreateOrUpdateUserProfile(&input)
	if err != nil {
		security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		return
	}

	security.SendEncryptedResponse(c, http.StatusOK, user)
}

// GetUser
// @Summary 获取用户详情
// @Description 可以通过 union_id, open_id 或 phone 查询用户。参数三选一。
// @Produce json
// @Param union_id query string false "用户 UnionID"
// @Param open_id query string false "用户 OpenID"
// @Param phone query string false "用户手机号"
// @Success 200 {object} models.UserProfile
// @Router /users [get]
func (h *UserProfileHandler) GetUser(c *gin.Context) {
	unionID := c.Query("union_id")
	openID := c.Query("open_id")
	phone := c.Query("phone")

	var user *models.UserProfile
	var err error

	switch {
	case unionID != "":
		user, err = h.service.GetUserByUnionID(unionID)
	case openID != "":
		user, err = h.service.GetUserByOpenID(openID)
	case phone != "":
		user, err = h.service.GetUserByPhone(phone)
	default:
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: "必须提供 union_id, open_id 或 phone 中的至少一个查询参数"})
		return
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			security.SendEncryptedResponse(c, http.StatusNotFound, security.ErrorResponse{Error: "用户未找到"})
		} else {
			security.SendEncryptedResponse(c, http.StatusInternalServerError, security.ErrorResponse{Error: err.Error()})
		}
		return
	}

	security.SendEncryptedResponse(c, http.StatusOK, user)
}

// BindPhoneNumber godoc
// @Summary      用户绑定手机号
// @Description  为指定用户绑定手机号
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        input body struct{user_union_id string "用户UnionID" phone_number string "手机号" phone_country_code string "国家代码"} true "绑定手机号请求体"
// @Success      200  {object}  security.EncryptedData
// @Failure      400  {object}  security.EncryptedData
// @Failure      404  {object}  security.EncryptedData
// @Failure      500  {object}  security.EncryptedData
// @Router       /users/bind-phone [post]
func (h *UserProfileHandler) BindPhoneNumber(c *gin.Context) {
	var input struct {
		UserUnionID      string `json:"user_union_id" binding:"required"`
		PhoneNumber      string `json:"phone_number" binding:"required"`
		PhoneCountryCode string `json:"phone_country_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: "无效的请求数据: " + err.Error()})
		return
	}

	user, err := h.service.BindPhoneNumber(input.UserUnionID, input.PhoneNumber, input.PhoneCountryCode)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "用户不存在") {
			status = http.StatusNotFound
		} else if strings.Contains(err.Error(), "已被其他用户绑定") {
			status = http.StatusBadRequest
		}

		security.SendEncryptedResponse(c, status, security.ErrorResponse{Error: err.Error()})
		return
	}

	security.SendEncryptedResponse(c, http.StatusOK, user)
}

// UnbindPhoneNumber godoc
// @Summary      用户解绑手机号
// @Description  解除指定用户的手机号绑定
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        input body struct{user_union_id string "用户UnionID"} true "解绑手机号请求体"
// @Success      200  {object}  security.EncryptedData
// @Failure      400  {object}  security.EncryptedData
// @Failure      404  {object}  security.EncryptedData
// @Failure      500  {object}  security.EncryptedData
// @Router       /users/unbind-phone [post]
func (h *UserProfileHandler) UnbindPhoneNumber(c *gin.Context) {
	var input struct {
		UserUnionID string `json:"user_union_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: "无效的请求数据: " + err.Error()})
		return
	}

	user, err := h.service.UnbindPhoneNumber(input.UserUnionID)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "用户不存在") {
			status = http.StatusNotFound
		}

		security.SendEncryptedResponse(c, status, security.ErrorResponse{Error: err.Error()})
		return
	}

	security.SendEncryptedResponse(c, http.StatusOK, user)
}

// GetUserScanHistory godoc
// @Summary 查询用户扫码门店历史
// @Description 获取指定用户的扫码历史记录
// @Tags Users
// @Accept  json
// @Produce  json
// @Param user_union_id query string true "用户UnionID"
// @Param start_date query string false "开始日期（格式：YYYY-MM-DD）"
// @Param end_date query string false "结束日期（格式：YYYY-MM-DD）"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} object{scan_history=[]service.UserScanHistoryItem, total=int64}
// @Failure 400 {object} security.ErrorResponse
// @Failure 404 {object} security.ErrorResponse
// @Failure 500 {object} security.ErrorResponse
// @Router /users/scan-history [get]
func (h *UserProfileHandler) GetUserScanHistory(c *gin.Context) {
	var input service.GetUserScanHistoryInput
	if err := c.ShouldBindQuery(&input); err != nil {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: "无效的查询参数: " + err.Error()})
		return
	}

	if input.UserUnionID == "" {
		security.SendEncryptedResponse(c, http.StatusBadRequest, security.ErrorResponse{Error: "用户ID不能为空"})
		return
	}

	history, total, err := h.service.GetUserScanHistory(&input)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "用户不存在") {
			status = http.StatusNotFound
		}
		security.SendEncryptedResponse(c, status, security.ErrorResponse{Error: err.Error()})
		return
	}

	security.SendEncryptedResponse(c, http.StatusOK, gin.H{
		"scan_history": history,
		"total":        total,
	})
}
