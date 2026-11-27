package system

import (
	"strconv"

	systemmodel "base/internal/model/system"
	"base/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListDictTypes(c *gin.Context) {
	var query systemmodel.PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, 400, "请求参数不正确")
		return
	}
	list, err := h.service.ListDictTypesPage(c.Request.Context(), query)
	if err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, list)
}

func (h *Handler) GetDictType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "字典类型编号不正确")
		return
	}
	data, err := h.service.GetDictType(c.Request.Context(), id)
	if err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, data)
}

func (h *Handler) CreateDictType(c *gin.Context) {
	var req systemmodel.CreateDictTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数不正确")
		return
	}
	if err := h.service.CreateDictType(c.Request.Context(), req); err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, true)
}

func (h *Handler) UpdateDictType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "字典类型编号不正确")
		return
	}
	var req systemmodel.UpdateDictTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数不正确")
		return
	}
	if err := h.service.UpdateDictType(c.Request.Context(), id, req); err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, true)
}

func (h *Handler) DeleteDictType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "字典类型编号不正确")
		return
	}
	if err := h.service.DeleteDictType(c.Request.Context(), id); err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, true)
}

func (h *Handler) ListDictData(c *gin.Context) {
	var query systemmodel.PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, 400, "请求参数不正确")
		return
	}
	dictTypeID, _ := strconv.Atoi(c.Query("dictTypeId"))
	list, err := h.service.ListDictDataPage(c.Request.Context(), query, dictTypeID)
	if err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, list)
}

func (h *Handler) GetDictData(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "字典数据编号不正确")
		return
	}
	data, err := h.service.GetDictData(c.Request.Context(), id)
	if err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, data)
}

func (h *Handler) CreateDictData(c *gin.Context) {
	var req systemmodel.CreateDictDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数不正确")
		return
	}
	if err := h.service.CreateDictData(c.Request.Context(), req); err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, true)
}

func (h *Handler) UpdateDictData(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "字典数据编号不正确")
		return
	}
	var req systemmodel.UpdateDictDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数不正确")
		return
	}
	if err := h.service.UpdateDictData(c.Request.Context(), id, req); err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, true)
}

func (h *Handler) DeleteDictData(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "字典数据编号不正确")
		return
	}
	if err := h.service.DeleteDictData(c.Request.Context(), id); err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, true)
}
