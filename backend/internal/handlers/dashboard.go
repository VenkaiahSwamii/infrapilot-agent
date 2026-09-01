package handlers

import (
	"net/http"
	"strconv"
	"time"

	"infrapilot/backend/internal/models"
	"infrapilot/backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DashboardHandler struct {
	service *services.DashboardService
}

func NewDashboardHandler(service *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		service: service,
	}
}

// Standard API response helpers
type apiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type apiError struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// GetDashboard returns real-time aggregated infrastructure statistics (legacy helper)
// GET /api/v1/dashboard
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	stats, err := services.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// Dashboard returns real-time aggregated infrastructure statistics (legacy helper)
func Dashboard(c *gin.Context) {
	stats, err := services.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetOverview alias for Dashboard
func GetOverview(c *gin.Context) {
	Dashboard(c)
}

// GetDashboards retrieves all dashboards for the current user
// GET /api/v1/dashboards
func (h *DashboardHandler) GetDashboards(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, apiError{Success: false, Error: "unauthorized"})
		return
	}

	ownerID := userID.(uuid.UUID)
	organizationIDStr := c.Query("organization_id")
	includeSharedStr := c.DefaultQuery("include_shared", "true")
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	var organizationID *uuid.UUID
	if organizationIDStr != "" {
		orgID, err := uuid.Parse(organizationIDStr)
		if err == nil {
			organizationID = &orgID
		}
	}

	includeShared := includeSharedStr == "true"

	dashboards, total, err := h.service.ListDashboards(ownerID, organizationID, includeShared)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiError{Success: false, Error: err.Error()})
		return
	}

	// Apply search filter if provided
	if search != "" {
		var filtered []models.Dashboard
		for _, d := range dashboards {
			if containsString(d.Name, search) || containsString(d.Description, search) {
				filtered = append(filtered, d)
			}
		}
		dashboards = filtered
	}

	// Calculate pagination
	totalPages := (int(total) + limit - 1) / limit
	offset := (page - 1) * limit
	end := offset + limit
	if end > len(dashboards) {
		end = len(dashboards)
	}
	if offset > len(dashboards) {
		offset = len(dashboards)
	}

	paginatedDashboards := dashboards[offset:end]

	c.JSON(http.StatusOK, apiResponse{
		Success: true,
		Data: gin.H{
			"dashboards":  paginatedDashboards,
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": totalPages,
		},
	})
}

// GetDashboardByID retrieves a specific dashboard by ID
// GET /api/v1/dashboards/:id
func (h *DashboardHandler) GetDashboardByID(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, apiError{Success: false, Error: "unauthorized"})
		return
	}

	dashboardIDStr := c.Param("id")
	dashboardID, err := uuid.Parse(dashboardIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, apiError{Success: false, Error: "invalid dashboard ID"})
		return
	}

	dashboard, err := h.service.GetDashboard(dashboardID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, apiError{Success: false, Error: "dashboard not found"})
		return
	}

	c.JSON(http.StatusOK, apiResponse{
		Success: true,
		Data:    dashboard,
	})
}

// CreateDashboard creates a new dashboard
// POST /api/v1/dashboards
func (h *DashboardHandler) CreateDashboard(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, apiError{Success: false, Error: "unauthorized"})
		return
	}

	var req services.CreateDashboardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiError{Success: false, Error: err.Error()})
		return
	}

	dashboard, err := h.service.CreateDashboard(userID.(uuid.UUID), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiError{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, apiResponse{
		Success: true,
		Message: "Dashboard created successfully",
		Data:    dashboard,
	})
}

// UpdateDashboard updates an existing dashboard
// PUT /api/v1/dashboards/:id
func (h *DashboardHandler) UpdateDashboard(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, apiError{Success: false, Error: "unauthorized"})
		return
	}

	dashboardIDStr := c.Param("id")
	dashboardID, err := uuid.Parse(dashboardIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, apiError{Success: false, Error: "invalid dashboard ID"})
		return
	}

	var req services.UpdateDashboardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiError{Success: false, Error: err.Error()})
		return
	}

	dashboard, err := h.service.UpdateDashboard(dashboardID, req, userID.(uuid.UUID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, apiError{Success: false, Error: "dashboard not found"})
		} else {
			c.JSON(http.StatusInternalServerError, apiError{Success: false, Error: err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, apiResponse{
		Success: true,
		Message: "Dashboard updated successfully",
		Data:    dashboard,
	})
}

// DeleteDashboard deletes a dashboard
// DELETE /api/v1/dashboards/:id
func (h *DashboardHandler) DeleteDashboard(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, apiError{Success: false, Error: "unauthorized"})
		return
	}

	dashboardIDStr := c.Param("id")
	dashboardID, err := uuid.Parse(dashboardIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, apiError{Success: false, Error: "invalid dashboard ID"})
		return
	}

	err = h.service.DeleteDashboard(dashboardID, userID.(uuid.UUID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, apiError{Success: false, Error: "dashboard not found"})
		} else {
			c.JSON(http.StatusInternalServerError, apiError{Success: false, Error: err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, apiResponse{
		Success: true,
		Message: "Dashboard deleted successfully",
	})
}

// SetDefaultDashboard sets a dashboard as the default
// POST /api/v1/dashboards/:id/set-default
func (h *DashboardHandler) SetDefaultDashboard(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, apiError{Success: false, Error: "unauthorized"})
		return
	}

	dashboardIDStr := c.Param("id")
	dashboardID, err := uuid.Parse(dashboardIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, apiError{Success: false, Error: "invalid dashboard ID"})
		return
	}

	err = h.service.SetDefaultDashboard(dashboardID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiError{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, apiResponse{
		Success: true,
		Message: "Default dashboard set successfully",
	})
}

// DuplicateDashboard duplicates an existing dashboard
// POST /api/v1/dashboards/:id/duplicate
func (h *DashboardHandler) DuplicateDashboard(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, apiError{Success: false, Error: "unauthorized"})
		return
	}

	dashboardIDStr := c.Param("id")
	dashboardID, err := uuid.Parse(dashboardIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, apiError{Success: false, Error: "invalid dashboard ID"})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiError{Success: false, Error: err.Error()})
		return
	}

	dashboard, err := h.service.DuplicateDashboard(dashboardID, userID.(uuid.UUID), req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiError{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, apiResponse{
		Success: true,
		Message: "Dashboard duplicated successfully",
		Data:    dashboard,
	})
}

// GetDefaultDashboard retrieves the default dashboard for the current user
// GET /api/v1/dashboards/default
func (h *DashboardHandler) GetDefaultDashboard(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, apiError{Success: false, Error: "unauthorized"})
		return
	}

	dashboard, err := h.service.GetDefaultDashboardForUser(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, apiError{Success: false, Error: "no default dashboard found"})
		return
	}

	c.JSON(http.StatusOK, apiResponse{
		Success: true,
		Data:    dashboard,
	})
}

// ShareDashboard shares a dashboard with users or teams
// POST /api/v1/dashboards/:id/share
func (h *DashboardHandler) ShareDashboard(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, apiError{Success: false, Error: "unauthorized"})
		return
	}

	dashboardIDStr := c.Param("id")
	dashboardID, err := uuid.Parse(dashboardIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, apiError{Success: false, Error: "invalid dashboard ID"})
		return
	}

	var req struct {
		SharedWith uuid.UUID `json:"shared_with" binding:"required"`
		ShareType  string    `json:"share_type" binding:"required,oneof=user team organization"`
		Permission string    `json:"permission" binding:"required,oneof=view edit admin"`
		ExpiresAt  *string   `json:"expires_at"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiError{Success: false, Error: err.Error()})
		return
	}

	var expiresAt *time.Time
	if req.ExpiresAt != nil && *req.ExpiresAt != "" {
		t, err := time.Parse(time.RFC3339, *req.ExpiresAt)
		if err == nil {
			expiresAt = &t
		}
	}

	err = h.service.ShareDashboard(dashboardID, userID.(uuid.UUID), req.SharedWith, req.ShareType, req.Permission, expiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiError{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, apiResponse{
		Success: true,
		Message: "Dashboard shared successfully",
	})
}

// GetWidgets retrieves all widgets for a dashboard
// GET /api/v1/dashboards/:id/widgets
func (h *DashboardHandler) GetWidgets(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, apiError{Success: false, Error: "unauthorized"})
		return
	}

	dashboardIDStr := c.Param("id")
	dashboardID, err := uuid.Parse(dashboardIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, apiError{Success: false, Error: "invalid dashboard ID"})
		return
	}

	// Verify permission
	_, err = h.service.GetDashboard(dashboardID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, apiError{Success: false, Error: "dashboard not found"})
		return
	}

	// Get widgets (this would need to be added to the service)
	// For now, return placeholder
	c.JSON(http.StatusOK, apiResponse{
		Success: true,
		Data:    []interface{}{},
	})
}

// CreateWidget creates a new widget in a dashboard
// POST /api/v1/dashboards/:id/widgets
func (h *DashboardHandler) CreateWidget(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, apiError{Success: false, Error: "unauthorized"})
		return
	}

	dashboardIDStr := c.Param("id")
	dashboardID, err := uuid.Parse(dashboardIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, apiError{Success: false, Error: "invalid dashboard ID"})
		return
	}

	var req services.CreateWidgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiError{Success: false, Error: err.Error()})
		return
	}

	widget, err := h.service.CreateWidget(dashboardID, req, userID.(uuid.UUID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, apiError{Success: false, Error: "dashboard not found"})
		} else {
			c.JSON(http.StatusInternalServerError, apiError{Success: false, Error: err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, apiResponse{
		Success: true,
		Message: "Widget created successfully",
		Data:    widget,
	})
}

// UpdateWidget updates a widget
// PATCH /api/v1/widgets/:id
func (h *DashboardHandler) UpdateWidget(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, apiError{Success: false, Error: "unauthorized"})
		return
	}

	widgetIDStr := c.Param("id")
	widgetID, err := uuid.Parse(widgetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, apiError{Success: false, Error: "invalid widget ID"})
		return
	}

	var req services.UpdateWidgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiError{Success: false, Error: err.Error()})
		return
	}

	// Get widget to find dashboard ID
	// This would need to be fetched from repository
	// For now, using a placeholder
	widget, err := h.service.UpdateWidget(uuid.Nil, widgetID, req, userID.(uuid.UUID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, apiError{Success: false, Error: "widget not found"})
		} else {
			c.JSON(http.StatusInternalServerError, apiError{Success: false, Error: err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, apiResponse{
		Success: true,
		Message: "Widget updated successfully",
		Data:    widget,
	})
}

// DeleteWidget deletes a widget
// DELETE /api/v1/widgets/:id
func (h *DashboardHandler) DeleteWidget(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, apiError{Success: false, Error: "unauthorized"})
		return
	}

	widgetIDStr := c.Param("id")
	widgetID, err := uuid.Parse(widgetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, apiError{Success: false, Error: "invalid widget ID"})
		return
	}

	// Get widget to find dashboard ID - this would need to be fetched from repository
	dashboardID := uuid.Nil // Placeholder

	err = h.service.DeleteWidget(dashboardID, widgetID, userID.(uuid.UUID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, apiError{Success: false, Error: "widget not found"})
		} else {
			c.JSON(http.StatusInternalServerError, apiError{Success: false, Error: err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, apiResponse{
		Success: true,
		Message: "Widget deleted successfully",
	})
}

// GetTemplates retrieves all dashboard templates
// GET /api/v1/dashboard-templates
func (h *DashboardHandler) GetTemplates(c *gin.Context) {
	category := c.Query("category")

	templates, err := h.service.GetTemplates(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiError{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, apiResponse{
		Success: true,
		Data:    templates,
	})
}

// GetTemplateWidgets retrieves widgets for a specific template
// GET /api/v1/dashboard-templates/:id/widgets
func (h *DashboardHandler) GetTemplateWidgets(c *gin.Context) {
	templateIDStr := c.Param("id")
	templateID, err := uuid.Parse(templateIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, apiError{Success: false, Error: "invalid template ID"})
		return
	}

	template, err := h.service.GetTemplateByID(templateID)
	if err != nil {
		c.JSON(http.StatusNotFound, apiError{Success: false, Error: "template not found"})
		return
	}

	c.JSON(http.StatusOK, apiResponse{
		Success: true,
		Data:    template.Widgets,
	})
}

// Helper function to check if a string contains another
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
