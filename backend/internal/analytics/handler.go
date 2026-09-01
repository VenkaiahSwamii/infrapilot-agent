package analytics

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
	orgID := c.Query("organization_id")
	overview, err := h.service.GetOverview(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, overview)
}

func (h *Handler) GetTrends(c *gin.Context) {
	orgID := c.Query("organization_id")
	rangeStr := c.DefaultQuery("range", "1d")

	trends, err := h.service.GetTrends(orgID, rangeStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"trends": trends})
}

func (h *Handler) GetCapacity(c *gin.Context) {
	orgID := c.Query("organization_id")
	capacity, err := h.service.GetCapacity(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"capacity_predictions": capacity})
}

func (h *Handler) GetSLA(c *gin.Context) {
	orgID := c.Query("organization_id")
	sla, err := h.service.GetSLA(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sla)
}

func (h *Handler) GetIncidents(c *gin.Context) {
	orgID := c.Query("organization_id")
	incidents, err := h.service.GetIncidents(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"incidents": incidents})
}

func (h *Handler) GetReports(c *gin.Context) {
	orgID := c.Query("organization_id")
	reports, err := h.service.GetReports(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"reports": reports})
}

func (h *Handler) GenerateReport(c *gin.Context) {
	var req struct {
		OrganizationID string `json:"organization_id"`
		Name           string `json:"name" binding:"required"`
		Type           string `json:"type" binding:"required"`
		Format         string `json:"format" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := c.GetString("username")
	if username == "" {
		username = "admin"
	}

	report, err := h.service.GenerateReport(req.OrganizationID, req.Name, req.Type, req.Format, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, report)
}

func (h *Handler) ScheduleReport(c *gin.Context) {
	var req struct {
		OrganizationID string `json:"organization_id"`
		Name           string `json:"name" binding:"required"`
		Type           string `json:"type" binding:"required"`
		Format         string `json:"format" binding:"required"`
		Frequency      string `json:"frequency" binding:"required"`
		Recipients     string `json:"recipients"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sReport, err := h.service.ScheduleReport(req.OrganizationID, req.Name, req.Type, req.Format, req.Frequency, req.Recipients)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sReport)
}

func (h *Handler) GetReportByID(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"id":           id,
		"name":         "Infrastructure Summary Report",
		"status":       "completed",
		"download_url": "/storage/reports/report_" + id + ".pdf",
	})
}
