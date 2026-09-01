package handlers

import (
	"net/http"

	"infrapilot/backend/internal/reports"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReportsHandler struct {
	generator *reports.ReportGenerator
}

func NewReportsHandler() *ReportsHandler {
	baseDir := "./storage/reports"
	return &ReportsHandler{
		generator: reports.NewReportGenerator(baseDir),
	}
}

// GenerateReport godoc
// @Summary Generate a new report
// @Description Generate a report in PDF, Excel, or CSV format
// @Tags Reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body reports.GenerateReportRequest true "Report configuration"
// @Success 200 {object} models.Report
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Router /reports/generate [post]
func (h *ReportsHandler) GenerateReport(c *gin.Context) {
	userIDStr := c.GetString("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req reports.GenerateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	report, err := h.generator.Generate(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate report: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// ListReports godoc
// @Summary List all reports
// @Description Get a list of all generated reports
// @Tags Reports
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Report
// @Router /reports [get]
func (h *ReportsHandler) ListReports(c *gin.Context) {
	userIDStr := c.GetString("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	reportList, err := h.generator.Repo.FindByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch reports"})
		return
	}

	c.JSON(http.StatusOK, reportList)
}

// DownloadReport godoc
// @Summary Download a report
// @Description Download a generated report file
// @Tags Reports
// @Produce application/octet-stream
// @Security BearerAuth
// @Param id path string true "Report ID"
// @Success 200 {file} binary
// @Failure 404 {object} gin.H
// @Router /reports/{id}/download [get]
func (h *ReportsHandler) DownloadReport(c *gin.Context) {
	reportID := c.Param("id")

	report, err := h.generator.Repo.FindByID(mustParseUUID(reportID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "report not found"})
		return
	}

	if report.Status != reports.ReportStatusCompleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "report is not ready for download"})
		return
	}

	c.File(report.FilePath)
}

// DeleteReport godoc
// @Summary Delete a report
// @Description Delete a generated report
// @Tags Reports
// @Security BearerAuth
// @Param id path string true "Report ID"
// @Success 204
// @Failure 404 {object} gin.H
// @Router /reports/{id} [delete]
func (h *ReportsHandler) DeleteReport(c *gin.Context) {
	reportID := c.Param("id")

	report, err := h.generator.Repo.FindByID(mustParseUUID(reportID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "report not found"})
		return
	}

	// Cleanup file
	if report.FilePath != "" {
		h.generator.PDF.Cleanup(report.FilePath)
	}

	// Delete from database
	if err := h.generator.Repo.Delete(report.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete report"})
		return
	}

	c.Status(http.StatusNoContent)
}

// ScheduleReport godoc
// @Summary Schedule a recurring report
// @Description Configure automatic report generation
// @Tags Reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body reports.ScheduleReportRequest true "Schedule configuration"
// @Success 200 {object} models.Report
// @Failure 400 {object} gin.H
// @Router /reports/schedule [post]
func (h *ReportsHandler) ScheduleReport(c *gin.Context) {
	userIDStr := c.GetString("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req reports.ScheduleReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	report, err := h.generator.Generate(userID, reports.GenerateReportRequest{
		Name:       req.Name,
		Type:       req.Type,
		Format:     req.Format,
		Schedule:   req.Schedule,
		DateFrom:   req.DateFrom,
		DateTo:     req.DateTo,
		MachineIDs: req.MachineIDs,
		Filters:    req.Filters,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to schedule report: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

func mustParseUUID(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return id
}
