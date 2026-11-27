package system

import (
	"strconv"

	systemmodel "base/internal/model/system"
	"base/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListRoles(c *gin.Context) {
	var query systemmodel.PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, 400, "请求参数不正确")
		return
	}
	list, err := h.service.ListRolesPage(c.Request.Context(), query)
	if err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, list)
}

func (h *Handler) GetRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "角色编号不正确")
		return
	}
	data, err := h.service.GetRole(c.Request.Context(), id)
	if err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, data)
}

func (h *Handler) CreateRole(c *gin.Context) {
	var req systemmodel.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数不正确")
		return
	}
	if err := h.service.CreateRole(c.Request.Context(), req); err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, true)
}

func (h *Handler) UpdateRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "角色编号不正确")
		return
	}
	var req systemmodel.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数不正确")
		return
	}
	if err := h.service.UpdateRole(c.Request.Context(), id, req); err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, true)
}

func (h *Handler) DeleteRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "角色编号不正确")
		return
	}
	if err := h.service.DeleteRole(c.Request.Context(), id); err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, true)
}

func (h *Handler) ReplaceRoleMenus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "角色编号不正确")
		return
	}
	var req systemmodel.AssignRoleMenusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数不正确")
		return
	}
	if err := h.service.AssignRoleMenus(c.Request.Context(), id, req.MenuIDs); err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, true)
}
