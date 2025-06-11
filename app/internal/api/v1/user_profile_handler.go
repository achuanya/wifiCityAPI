package v1

import (
	"app/internal/models"
	"app/internal/service"
	"app/pkg/security"
	"errors"
	"net/http"

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
