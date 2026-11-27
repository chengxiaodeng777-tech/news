package system

import (
	"base/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetMonitor(c *gin.Context) {
	monitorService := h.monitorService
	if monitorService == nil {
		response.Error(c, 500, "监控服务未初始化")
		return
	}

	data, err := monitorService.GetMonitor(c.Request.Context())
	if err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, data)
}
