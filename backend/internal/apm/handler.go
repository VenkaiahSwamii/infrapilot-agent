package apm

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetOverview(c *gin.Context) {
	data, err := h.service.GetAPMData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data.Overview)
}

func (h *Handler) GetEndpoints(c *gin.Context) {
	data, err := h.service.GetAPMData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"endpoints": data.Endpoints, "slow_endpoints": data.SlowEndpoints})
}

func (h *Handler) GetServices(c *gin.Context) {
	data, err := h.service.GetAPMData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data.Services)
}

func (h *Handler) GetErrors(c *gin.Context) {
	data, err := h.service.GetAPMData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error_rate": data.Overview.ErrorRate, "recommendations": data.Recommendations})
}

func (h *Handler) GetLatency(c *gin.Context) {
	data, err := h.service.GetAPMData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"avg_latency_ms": data.Overview.AverageLatency, "slow_endpoints": data.SlowEndpoints})
}

func (h *Handler) GetTop(c *gin.Context) {
	data, err := h.service.GetAPMData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data.TopAPIs)
}

func (h *Handler) GetAPMData(c *gin.Context) {
	data, err := h.service.GetAPMData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
