package handlers

import (
	"net/http"

	"github.com/FebryanHernanda/BE-EventOrganizer/internal/services"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	service services.HealthService
}

func NewHealthHandler(service services.HealthService) *HealthHandler {
	return &HealthHandler{service}
}

func (h *HealthHandler) HealthCheck(ctx *gin.Context) {
	result := h.service.GetHealth()
	ctx.JSON(http.StatusOK, result)
}
