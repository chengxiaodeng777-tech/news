package system

import (
	"strconv"

	systemmodel "base/internal/model/system"
	"base/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListUsers(c *gin.Context) {
	var query systemmodel.PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, 400, "请求参数不正确")
		return
	}
	list, err := h.service.ListUsersPage(c.Request.Context(), query)
	if err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, list)
}

func (h *Handler) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "用户编号不正确")
		return
	}
	data, err := h.service.GetUser(c.Request.Context(), id)
	if err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, data)
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req systemmodel.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数不正确")
		return
	}
	if err := h.service.CreateUser(c.Request.Context(), req); err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, true)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "用户编号不正确")
		return
	}
	var req systemmodel.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数不正确")
		return
	}
	if err := h.service.UpdateUser(c.Request.Context(), id, req); err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, true)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "用户编号不正确")
		return
	}
	if err := h.service.DeleteUser(c.Request.Context(), id); err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, true)
}

func (h *Handler) ReplaceUserRoles(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "用户编号不正确")
		return
	}
	var req systemmodel.AssignUserRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数不正确")
		return
	}
	if err := h.service.AssignUserRoles(c.Request.Context(), id, req.RoleIDs); err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, true)
}

func (h *Handler) ResetUserPassword(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "用户编号不正确")
		return
	}
	var req systemmodel.ResetUserPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数不正确")
		return
	}
	if err := h.service.ResetUserPassword(c.Request.Context(), id, req.Password); err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, true)
}
