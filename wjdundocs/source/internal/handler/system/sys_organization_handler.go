package system

import (
	"strconv"

	systemmodel "base/internal/model/system"
	"base/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListOrganizations(c *gin.Context) {
	list, err := h.service.ListOrganizations(c.Request.Context())
	if err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, list)
}

func (h *Handler) ListOrganizationTree(c *gin.Context) {
	list, err := h.service.ListOrganizationTree(c.Request.Context())
	if err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, list)
}

func (h *Handler) GetOrganization(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "组织编号不正确")
		return
	}
	data, err := h.service.GetOrganization(c.Request.Context(), id)
	if err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, data)
}

func (h *Handler) CreateOrganization(c *gin.Context) {
	var req systemmodel.CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数不正确")
		return
	}
	if err := h.service.CreateOrganization(c.Request.Context(), req); err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, true)
}

func (h *Handler) UpdateOrganization(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "组织编号不正确")
		return
	}
	var req systemmodel.UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数不正确")
		return
	}
	if err := h.service.UpdateOrganization(c.Request.Context(), id, req); err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, true)
}

func (h *Handler) DeleteOrganization(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "组织编号不正确")
		return
	}
	if err := h.service.DeleteOrganization(c.Request.Context(), id); err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, true)
}
