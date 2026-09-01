package handlers

import (
	"net/http"

	"infrapilot/backend/internal/hardening"

	"github.com/gin-gonic/gin"
)

type HardeningHandler struct {
	engine *hardening.HardeningEngine
}

func NewHardeningHandler() *HardeningHandler {
	return &HardeningHandler{
		engine: hardening.NewHardeningEngine(),
	}
}

// GetReadiness returns production engineering benchmarks and readiness status
func (h *HardeningHandler) GetReadiness(c *gin.Context) {
	report := h.engine.GenerateProductionReport()
	c.JSON(http.StatusOK, report)
}

// GetBenchmarks returns target vs observed latencies
func (h *HardeningHandler) GetBenchmarks(c *gin.Context) {
	report := h.engine.GenerateProductionReport()
	c.JSON(http.StatusOK, gin.H{
		"benchmarks": report.Benchmarks,
		"status":     report.Status,
	})
}
