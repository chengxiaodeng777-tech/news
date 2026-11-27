package system

import systemservice "base/internal/service/system"

type Handler struct {
	service        *systemservice.Service
	monitorService *systemservice.MonitorService
}

func New(service *systemservice.Service, monitorService *systemservice.MonitorService) *Handler {
	return &Handler{
		service:        service,
		monitorService: monitorService,
	}
}
