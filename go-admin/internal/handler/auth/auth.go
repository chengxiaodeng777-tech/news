package auth

import (
	"errors"

	authmodel "base/internal/model/auth"
	"base/internal/pkg/response"
	authservice "base/internal/service/auth"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *authservice.Service
}

func New(service *authservice.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Login(c *gin.Context) {
	var req authmodel.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数不正确")
		return
	}

	token, err := h.service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, authmodel.LoginResponse{AccessToken: token})
}

func (h *Handler) UserInfo(c *gin.Context) {
	userID, exists := c.Get("current_user_id")
	if !exists {
		response.Error(c, 401, "未获取到当前用户")
		return
	}

	info, err := h.service.GetUserInfo(c.Request.Context(), userID.(int))
	if err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, info)
}

func (h *Handler) AccessCodes(c *gin.Context) {
	userID, exists := c.Get("current_user_id")
	if !exists {
		response.Error(c, 401, "未获取到当前用户")
		return
	}

	codes, err := h.service.GetPermissionCodes(c.Request.Context(), userID.(int))
	if err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, codes)
}

func (h *Handler) Menus(c *gin.Context) {
	userID, exists := c.Get("current_user_id")
	if !exists {
		response.Error(c, 401, "未获取到当前用户")
		return
	}

	menus, err := h.service.GetMenuTree(c.Request.Context(), userID.(int))
	if err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, menus)
}

func ParseCurrentUserID(c *gin.Context) (int, error) {
	userID, exists := c.Get("current_user_id")
	if !exists {
		return 0, errors.New("未获取到当前用户")
	}
	id, ok := userID.(int)
	if !ok {
		return 0, errors.New("当前用户类型不正确")
	}
	return id, nil
}
