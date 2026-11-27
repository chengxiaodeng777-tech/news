package system

import (
	"strconv"

	systemmodel "base/internal/model/system"
	"base/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListMenus(c *gin.Context) {
	list, err := h.service.ListMenus(c.Request.Context())
	if err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, list)
}

func (h *Handler) ListMenuTree(c *gin.Context) {
	list, err := h.service.ListMenuTree(c.Request.Context())
	if err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, list)
}

func (h *Handler) GetMenu(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "菜单编号不正确")
		return
	}
	data, err := h.service.GetMenu(c.Request.Context(), id)
	if err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, data)
}

func (h *Handler) CreateMenu(c *gin.Context) {
	var req systemmodel.CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数不正确")
		return
	}
	if err := h.service.CreateMenu(c.Request.Context(), req); err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, true)
}

func (h *Handler) UpdateMenu(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "菜单编号不正确")
		return
	}
	var req systemmodel.UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数不正确")
		return
	}
	if err := h.service.UpdateMenu(c.Request.Context(), id, req); err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, true)
}

func (h *Handler) DeleteMenu(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "菜单编号不正确")
		return
	}
	if err := h.service.DeleteMenu(c.Request.Context(), id); err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, true)
}
