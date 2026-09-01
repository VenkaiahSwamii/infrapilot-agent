package services

import (
	"encoding/json"
	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"
	"infrapilot/backend/internal/repository"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DashboardService handles dashboard business logic
type DashboardService struct {
	dashboardRepo *repository.DashboardRepository
}

// NewDashboardService creates a new dashboard service
func NewDashboardService() *DashboardService {
	return &DashboardService{
		dashboardRepo: repository.NewDashboardRepository(),
	}
}

// DashboardResponse represents the response for dashboard operations
type DashboardResponse struct {
	Success   bool                       `json:"success"`
	Message   string                     `json:"message,omitempty"`
	Data      *models.Dashboard          `json:"data,omitempty"`
	Widgets   []models.DashboardWidget   `json:"widgets,omitempty"`
	Templates []models.DashboardTemplate `json:"templates,omitempty"`
	Total     int64                      `json:"total,omitempty"`
}

// CreateDashboardRequest represents a request to create a dashboard
type CreateDashboardRequest struct {
	Name           string                 `json:"name" binding:"required"`
	Description    string                 `json:"description"`
	OrganizationID *uuid.UUID             `json:"organization_id"`
	TemplateID     *uuid.UUID             `json:"template_id,omitempty"`
	Config         map[string]interface{} `json:"config,omitempty"`
	IsDefault      bool                   `json:"is_default"`
}

// UpdateDashboardRequest represents a request to update a dashboard
type UpdateDashboardRequest struct {
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	IsShared    *bool                  `json:"is_shared,omitempty"`
	IsDefault   *bool                  `json:"is_default,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
	SharedWith  []uuid.UUID            `json:"shared_with,omitempty"`
}

// CreateWidgetRequest represents a request to create a widget
type CreateWidgetRequest struct {
	WidgetType models.WidgetType      `json:"widget_type" binding:"required"`
	Title      string                 `json:"title" binding:"required"`
	X          int                    `json:"x"`
	Y          int                    `json:"y"`
	Width      int                    `json:"width"`
	Height     int                    `json:"height"`
	Config     map[string]interface{} `json:"config,omitempty"`
}

// UpdateWidgetRequest represents a request to update a widget
type UpdateWidgetRequest struct {
	WidgetType *models.WidgetType     `json:"widget_type,omitempty"`
	Title      *string                `json:"title,omitempty"`
	X          *int                   `json:"x,omitempty"`
	Y          *int                   `json:"y,omitempty"`
	Width      *int                   `json:"width,omitempty"`
	Height     *int                   `json:"height,omitempty"`
	Config     map[string]interface{} `json:"config,omitempty"`
	Order      *int                   `json:"order,omitempty"`
}

// CreateDashboard creates a new dashboard
func (s *DashboardService) CreateDashboard(ownerID uuid.UUID, req CreateDashboardRequest) (*models.Dashboard, error) {
	dashboard := &models.Dashboard{
		OwnerID:     ownerID,
		Name:        req.Name,
		Description: req.Description,
		IsDefault:   req.IsDefault,
		CreatedBy:   ownerID,
		TemplateID:  req.TemplateID,
	}

	// Set organization ID
	if req.OrganizationID != nil {
		dashboard.OrganizationID = *req.OrganizationID
	} else {
		dashboard.OrganizationID = uuid.Nil
	}

	// Marshal config
	if req.Config != nil {
		configBytes, err := json.Marshal(req.Config)
		if err != nil {
			return nil, err
		}
		dashboard.Config = configBytes
	}

	// If this is set as default, handle the transaction
	if req.IsDefault {
		err := database.DB.Transaction(func(tx *gorm.DB) error {
			// Unset all other defaults for this owner
			if err := tx.Model(&models.Dashboard{}).
				Where("owner_id = ?", ownerID).
				Update("is_default", false).Error; err != nil {
				return err
			}

			// Create the new default dashboard
			if err := tx.Create(dashboard).Error; err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		// Just create the dashboard
		if err := database.DB.Create(dashboard).Error; err != nil {
			return nil, err
		}
	}

	// Load the dashboard with widgets
	err := database.DB.Preload("Widgets").First(dashboard, dashboard.ID).Error
	return dashboard, err
}

// GetDashboard retrieves a dashboard by ID
func (s *DashboardService) GetDashboard(id uuid.UUID, userID uuid.UUID) (*models.Dashboard, error) {
	dashboard, err := s.dashboardRepo.GetDashboardWithWidgets(id)
	if err != nil {
		return nil, err
	}

	// Check permissions
	permission, err := s.dashboardRepo.CheckSharePermission(id, userID)
	if err != nil && permission != "owner" {
		return nil, err
	}

	return dashboard, nil
}

// ListDashboards retrieves dashboards for a user
func (s *DashboardService) ListDashboards(ownerID uuid.UUID, organizationID *uuid.UUID, includeShared bool) ([]models.Dashboard, int64, error) {
	filters := repository.DashboardFilters{
		OwnerID:        &ownerID,
		OrganizationID: organizationID,
	}

	// Get owned dashboards
	dashboards, total, err := s.dashboardRepo.ListDashboards(filters)
	if err != nil {
		return nil, 0, err
	}

	// If including shared dashboards, add them
	if includeShared {
		sharedDashboards, _, err := s.dashboardRepo.ListDashboards(repository.DashboardFilters{
			IsShared: func() *bool { b := true; return &b }(),
		})
		if err == nil {
			// Filter out duplicates and add shared dashboards
			existingIDs := make(map[uuid.UUID]bool)
			for _, d := range dashboards {
				existingIDs[d.ID] = true
			}
			for _, d := range sharedDashboards {
				if !existingIDs[d.ID] {
					dashboards = append(dashboards, d)
				}
			}
		}
	}

	return dashboards, total, nil
}

// UpdateDashboard updates a dashboard
func (s *DashboardService) UpdateDashboard(id uuid.UUID, req UpdateDashboardRequest, userID uuid.UUID) (*models.Dashboard, error) {
	// Verify ownership or permission
	permission, err := s.dashboardRepo.CheckSharePermission(id, userID)
	if err != nil {
		return nil, err
	}
	if permission != "owner" && permission != "edit" && permission != "admin" {
		return nil, gorm.ErrRecordNotFound
	}

	// Build updates map
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.IsShared != nil {
		updates["is_shared"] = *req.IsShared
	}

	// Handle default flag
	if req.IsDefault != nil && *req.IsDefault {
		// First unset all other defaults
		if err := database.DB.Model(&models.Dashboard{}).
			Where("owner_id = ?", userID).
			Update("is_default", false).Error; err != nil {
			return nil, err
		}
		updates["is_default"] = true
	}

	// Update config
	if req.Config != nil {
		configBytes, err := json.Marshal(req.Config)
		if err != nil {
			return nil, err
		}
		updates["config"] = configBytes
	}

	// Update shared_with field
	if req.SharedWith != nil {
		sharedWithBytes, err := json.Marshal(req.SharedWith)
		if err != nil {
			return nil, err
		}
		updates["shared_with"] = sharedWithBytes
	}

	// Apply updates
	err = s.dashboardRepo.UpdateDashboard(id, updates)
	if err != nil {
		return nil, err
	}

	// Return updated dashboard
	return s.dashboardRepo.GetDashboardWithWidgets(id)
}

// DeleteDashboard deletes a dashboard
func (s *DashboardService) DeleteDashboard(id uuid.UUID, userID uuid.UUID) error {
	// Verify ownership or permission
	permission, err := s.dashboardRepo.CheckSharePermission(id, userID)
	if err != nil {
		return err
	}
	if permission != "owner" && permission != "admin" {
		return gorm.ErrRecordNotFound
	}

	return s.dashboardRepo.DeleteDashboard(id)
}

// SetDefaultDashboard sets a dashboard as default
func (s *DashboardService) SetDefaultDashboard(id uuid.UUID, userID uuid.UUID) error {
	return s.dashboardRepo.SetDefaultDashboard(id, userID)
}

// CreateWidget creates a new widget in a dashboard
func (s *DashboardService) CreateWidget(dashboardID uuid.UUID, req CreateWidgetRequest, userID uuid.UUID) (*models.DashboardWidget, error) {
	// Verify dashboard ownership
	permission, err := s.dashboardRepo.CheckSharePermission(dashboardID, userID)
	if err != nil {
		return nil, err
	}
	if permission != "owner" && permission != "edit" && permission != "admin" {
		return nil, gorm.ErrRecordNotFound
	}

	widget := &models.DashboardWidget{
		DashboardID: dashboardID,
		WidgetType:  req.WidgetType,
		Title:       req.Title,
		X:           req.X,
		Y:           req.Y,
		Width:       req.Width,
		Height:      req.Height,
	}

	// Marshal config
	if req.Config != nil {
		configBytes, err := json.Marshal(req.Config)
		if err != nil {
			return nil, err
		}
		widget.Config = configBytes
	}

	// Get next order value
	widgets, _ := s.dashboardRepo.GetWidgetsByDashboardID(dashboardID)
	widget.Order = len(widgets)

	err = s.dashboardRepo.CreateWidget(widget)
	if err != nil {
		return nil, err
	}

	return widget, nil
}

// UpdateWidget updates a widget
func (s *DashboardService) UpdateWidget(dashboardID, widgetID uuid.UUID, req UpdateWidgetRequest, userID uuid.UUID) (*models.DashboardWidget, error) {
	// Verify dashboard ownership
	permission, err := s.dashboardRepo.CheckSharePermission(dashboardID, userID)
	if err != nil {
		return nil, err
	}
	if permission != "owner" && permission != "edit" && permission != "admin" {
		return nil, gorm.ErrRecordNotFound
	}

	// Verify widget belongs to dashboard
	widget, err := s.dashboardRepo.GetWidgetByID(widgetID)
	if err != nil {
		return nil, err
	}
	if widget.DashboardID != dashboardID {
		return nil, gorm.ErrRecordNotFound
	}

	// Build updates map
	updates := make(map[string]interface{})

	if req.WidgetType != nil {
		updates["widget_type"] = *req.WidgetType
	}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.X != nil {
		updates["x"] = *req.X
	}
	if req.Y != nil {
		updates["y"] = *req.Y
	}
	if req.Width != nil {
		updates["width"] = *req.Width
	}
	if req.Height != nil {
		updates["height"] = *req.Height
	}
	if req.Order != nil {
		updates["order"] = *req.Order
	}

	if req.Config != nil {
		configBytes, err := json.Marshal(req.Config)
		if err != nil {
			return nil, err
		}
		updates["config"] = configBytes
	}

	err = s.dashboardRepo.UpdateWidget(widgetID, updates)
	if err != nil {
		return nil, err
	}

	return s.dashboardRepo.GetWidgetByID(widgetID)
}

// DeleteWidget deletes a widget
func (s *DashboardService) DeleteWidget(dashboardID, widgetID uuid.UUID, userID uuid.UUID) error {
	// Verify dashboard ownership
	permission, err := s.dashboardRepo.CheckSharePermission(dashboardID, userID)
	if err != nil {
		return err
	}
	if permission != "owner" && permission != "edit" && permission != "admin" {
		return gorm.ErrRecordNotFound
	}

	// Verify widget belongs to dashboard
	widget, err := s.dashboardRepo.GetWidgetByID(widgetID)
	if err != nil {
		return err
	}
	if widget.DashboardID != dashboardID {
		return gorm.ErrRecordNotFound
	}

	return s.dashboardRepo.DeleteWidget(widgetID)
}

// ShareDashboard shares a dashboard with users or teams
func (s *DashboardService) ShareDashboard(dashboardID uuid.UUID, sharedBy uuid.UUID, sharedWith uuid.UUID, shareType, permission string, expiresAt *time.Time) error {
	// Verify dashboard ownership
	_, err := s.dashboardRepo.CheckSharePermission(dashboardID, sharedBy)
	if err != nil {
		return err
	}

	share := &models.DashboardShare{
		DashboardID: dashboardID,
		SharedBy:    sharedBy,
		SharedWith:  sharedWith,
		ShareType:   shareType,
		Permission:  permission,
		ExpiresAt:   expiresAt,
		CreatedAt:   time.Now(),
	}

	return s.dashboardRepo.CreateShare(share)
}

// GetDefaultDashboardForUser retrieves the default dashboard for a user
func (s *DashboardService) GetDefaultDashboardForUser(ownerID uuid.UUID) (*models.Dashboard, error) {
	return s.dashboardRepo.GetDefaultDashboard(ownerID)
}

// DuplicateDashboard duplicates a dashboard
func (s *DashboardService) DuplicateDashboard(dashboardID uuid.UUID, newOwnerID uuid.UUID, newName string) (*models.Dashboard, error) {
	return s.dashboardRepo.DuplicateDashboard(dashboardID, newOwnerID, newName)
}

// GetTemplates retrieves dashboard templates
func (s *DashboardService) GetTemplates(category string) ([]models.DashboardTemplate, error) {
	return s.dashboardRepo.GetTemplates(category)
}

// GetTemplateByID retrieves a template by ID
func (s *DashboardService) GetTemplateByID(id uuid.UUID) (*models.DashboardTemplate, error) {
	return s.dashboardRepo.GetTemplateByID(id)
}

// AITemplateSuggestionRequest represents a request for AI template suggestions
type AITemplateSuggestionRequest struct {
	Context string `json:"context"` // Context for AI to analyze
}

// AISuggestedWidget represents a widget suggested by AI
type AISuggestedWidget struct {
	WidgetType models.WidgetType      `json:"widget_type"`
	Title      string                 `json:"title"`
	Reason     string                 `json:"reason"` // Why AI suggests this widget
	Config     map[string]interface{} `json:"config,omitempty"`
}

// GetAITemplateSuggestions gets AI-suggested widgets based on context
func (s *DashboardService) GetAITemplateSuggestions(context string) ([]AISuggestedWidget, error) {
	// This is a placeholder - in production, this would call the AI service
	// For now, return some sensible defaults based on keywords in context
	var suggestions []AISuggestedWidget

	// Simple keyword-based suggestions (replace with actual AI integration)
	if len(context) > 0 {
		contextLower := context

		// Kubernetes-related
		if len(contextLower) > 0 {
			suggestions = append(suggestions, AISuggestedWidget{
				WidgetType: models.WidgetTypeClusterHealth,
				Title:      "Cluster Health",
				Reason:     "High Kubernetes activity detected",
				Config:     map[string]interface{}{"showDetails": true},
			})
			suggestions = append(suggestions, AISuggestedWidget{
				WidgetType: models.WidgetTypeK8sPods,
				Title:      "Pod Status",
				Reason:     "Monitor pod counts and status",
			})
		}

		// High error detection
		suggestions = append(suggestions, AISuggestedWidget{
			WidgetType: models.WidgetTypeErrorLogs,
			Title:      "Error Logs",
			Reason:     "Track error patterns",
			Config:     map[string]interface{}{"tailLines": 100, "logLevel": "error"},
		})
	}

	return suggestions, nil
}

// GetWidgetData gets data for a specific widget
func (s *DashboardService) GetWidgetData(widget *models.DashboardWidget, dashboard *models.Dashboard) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"id":     widget.ID,
		"type":   widget.WidgetType,
		"title":  widget.Title,
		"config": widget.Config,
	}

	// Unmarshal config
	var widgetConfig models.WidgetConfig
	if len(widget.Config) > 0 {
		json.Unmarshal(widget.Config, &widgetConfig)
	}

	// Fetch data based on widget type
	switch widget.WidgetType {
	case models.WidgetTypeCPU, models.WidgetTypeMemory, models.WidgetTypeDisk:
		result["data"] = s.getResourceData(widgetConfig)
	case models.WidgetTypeAlerts, models.WidgetTypeActiveIncidents:
		result["data"] = s.getAlertsData(widgetConfig)
	case models.WidgetTypeK8sPods, models.WidgetTypeK8sNodes, models.WidgetTypeK8sDeployments, models.WidgetTypeK8sServices:
		result["data"] = s.getKubernetesData(widget.WidgetType, widgetConfig)
	case models.WidgetTypeContainers, models.WidgetTypeImages, models.WidgetTypeVolumes:
		result["data"] = s.getDockerData(widget.WidgetType, widgetConfig)
	case models.WidgetTypeLiveLogs, models.WidgetTypeErrorLogs:
		result["data"] = s.getLogsData(widget.WidgetType, widgetConfig)
	case models.WidgetTypeHealthScore:
		result["data"] = s.getHealthScoreData(widgetConfig)
	case models.WidgetTypeLatency, models.WidgetTypeRPS, models.WidgetTypeErrorRate, models.WidgetTypeTopAPIs:
		result["data"] = s.getAPMData(widget.WidgetType, widgetConfig)
	case models.WidgetTypeAIInsights, models.WidgetTypeRootCause, models.WidgetTypePredictions:
		result["data"] = s.getAIData(widget.WidgetType, widgetConfig)
	default:
		result["data"] = map[string]interface{}{
			"message": "Widget type not implemented",
		}
	}

	return result, nil
}

// Placeholder implementations for widget data fetching
func (s *DashboardService) getResourceData(config models.WidgetConfig) map[string]interface{} {
	// In production, this would fetch from the metric service
	return map[string]interface{}{
		"current": 45.5,
		"average": 42.3,
		"peak":    78.9,
		"unit":    "%",
		"history": []float64{40.0, 42.0, 45.0, 43.0, 45.5},
	}
}

func (s *DashboardService) getAlertsData(config models.WidgetConfig) map[string]interface{} {
	// In production, this would fetch from the alerts service
	return map[string]interface{}{
		"total":    12,
		"critical": 2,
		"warning":  5,
		"info":     5,
		"items":    []interface{}{},
	}
}

func (s *DashboardService) getKubernetesData(widgetType models.WidgetType, config models.WidgetConfig) map[string]interface{} {
	// In production, this would fetch from the Kubernetes service
	return map[string]interface{}{
		"widget_type":   widgetType,
		"cluster_count": 1,
		"items":         []interface{}{},
	}
}

func (s *DashboardService) getDockerData(widgetType models.WidgetType, config models.WidgetConfig) map[string]interface{} {
	// In production, this would fetch from the Docker service
	return map[string]interface{}{
		"widget_type":     widgetType,
		"container_count": 5,
		"items":           []interface{}{},
	}
}

func (s *DashboardService) getLogsData(widgetType models.WidgetType, config models.WidgetConfig) map[string]interface{} {
	// In production, this would fetch from the logs service
	return map[string]interface{}{
		"widget_type": widgetType,
		"lines":       []interface{}{},
		"count":       0,
	}
}

func (s *DashboardService) getHealthScoreData(config models.WidgetConfig) map[string]interface{} {
	// In production, this would fetch from the health score service
	return map[string]interface{}{
		"score": 85,
		"grade": "B+",
		"components": map[string]interface{}{
			"infrastructure": 90,
			"applications":   82,
			"security":       85,
		},
	}
}

func (s *DashboardService) getAPMData(widgetType models.WidgetType, config models.WidgetConfig) map[string]interface{} {
	// In production, this would fetch from the APM service
	return map[string]interface{}{
		"widget_type": widgetType,
		"current":     42.5,
		"unit":        "ms",
		"history":     []float64{40.0, 41.0, 43.0, 42.0, 42.5},
	}
}

func (s *DashboardService) getAIData(widgetType models.WidgetType, config models.WidgetConfig) map[string]interface{} {
	// In production, this would fetch from the AI service
	return map[string]interface{}{
		"widget_type": widgetType,
		"insights":    []interface{}{},
		"score":       75,
	}
}

// GetDashboardWithUserContext gets dashboard with user context
func (s *DashboardService) GetDashboardWithUserContext(dashboard *models.Dashboard, userID uuid.UUID) map[string]interface{} {
	result := map[string]interface{}{
		"id":              dashboard.ID,
		"name":            dashboard.Name,
		"description":     dashboard.Description,
		"is_default":      dashboard.IsDefault,
		"is_shared":       dashboard.IsShared,
		"config":          dashboard.Config,
		"created_at":      dashboard.CreatedAt,
		"updated_at":      dashboard.UpdatedAt,
		"widgets":         dashboard.Widgets,
		"user_permission": "owner", // Would be fetched from CheckSharePermission
	}

	return result
}
