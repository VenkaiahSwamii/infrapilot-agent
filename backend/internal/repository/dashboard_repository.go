package repository

import (
	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DashboardRepository handles database operations for dashboards and widgets
type DashboardRepository struct {
	db *gorm.DB
}

// NewDashboardRepository creates a new dashboard repository
func NewDashboardRepository() *DashboardRepository {
	return &DashboardRepository{
		db: database.DB,
	}
}

// DashboardFilters represents filters for listing dashboards
type DashboardFilters struct {
	OwnerID        *uuid.UUID
	OrganizationID *uuid.UUID
	IsDefault      *bool
	IsShared       *bool
	Search         string
	Limit          int
	Offset         int
}

// CreateDashboard creates a new dashboard
func (r *DashboardRepository) CreateDashboard(dashboard *models.Dashboard) error {
	return r.db.Create(dashboard).Error
}

// GetDashboardByID retrieves a dashboard by ID with its widgets
func (r *DashboardRepository) GetDashboardByID(id uuid.UUID) (*models.Dashboard, error) {
	var dashboard models.Dashboard
	err := r.db.Where("id = ?", id).First(&dashboard).Error
	if err != nil {
		return nil, err
	}
	return &dashboard, nil
}

// GetDashboardWithWidgets retrieves a dashboard with all its widgets
func (r *DashboardRepository) GetDashboardWithWidgets(id uuid.UUID) (*models.Dashboard, error) {
	var dashboard models.Dashboard
	err := r.db.
		Preload("Widgets").
		Where("id = ?", id).
		First(&dashboard).Error
	if err != nil {
		return nil, err
	}
	return &dashboard, nil
}

// ListDashboards retrieves a list of dashboards based on filters
func (r *DashboardRepository) ListDashboards(filters DashboardFilters) ([]models.Dashboard, int64, error) {
	var dashboards []models.Dashboard
	var total int64

	query := r.db.Model(&models.Dashboard{})

	// Apply filters
	if filters.OwnerID != nil {
		query = query.Where("owner_id = ?", *filters.OwnerID)
	}
	if filters.OrganizationID != nil {
		query = query.Where("organization_id = ?", *filters.OrganizationID)
	}
	if filters.IsDefault != nil {
		query = query.Where("is_default = ?", *filters.IsDefault)
	}
	if filters.IsShared != nil {
		query = query.Where("is_shared = ?", *filters.IsShared)
	}
	if filters.Search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+filters.Search+"%", "%"+filters.Search+"%")
	}

	// Count total
	query.Count(&total)

	// Apply pagination and ordering
	query = query.Order("updated_at DESC")
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	err := query.Find(&dashboards).Error
	return dashboards, total, err
}

// UpdateDashboard updates a dashboard
func (r *DashboardRepository) UpdateDashboard(id uuid.UUID, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	return r.db.Model(&models.Dashboard{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteDashboard deletes a dashboard (soft delete)
func (r *DashboardRepository) DeleteDashboard(id uuid.UUID) error {
	return r.db.Delete(&models.Dashboard{}, id).Error
}

// SetDefaultDashboard sets a dashboard as default and unsets others
func (r *DashboardRepository) SetDefaultDashboard(id uuid.UUID, ownerID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Unset all dashboards for this owner
		if err := tx.Model(&models.Dashboard{}).
			Where("owner_id = ?", ownerID).
			Update("is_default", false).Error; err != nil {
			return err
		}

		// Set this dashboard as default
		if err := tx.Model(&models.Dashboard{}).
			Where("id = ?", id).
			Updates(map[string]interface{}{
				"is_default": true,
				"updated_at": time.Now(),
			}).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetDefaultDashboard retrieves the default dashboard for a user
func (r *DashboardRepository) GetDefaultDashboard(ownerID uuid.UUID) (*models.Dashboard, error) {
	var dashboard models.Dashboard
	err := r.db.
		Where("owner_id = ? AND is_default = ?", ownerID, true).
		First(&dashboard).Error
	if err != nil {
		return nil, err
	}
	return &dashboard, nil
}

// Widget operations

// CreateWidget creates a new widget
func (r *DashboardRepository) CreateWidget(widget *models.DashboardWidget) error {
	return r.db.Create(widget).Error
}

// GetWidgetByID retrieves a widget by ID
func (r *DashboardRepository) GetWidgetByID(id uuid.UUID) (*models.DashboardWidget, error) {
	var widget models.DashboardWidget
	err := r.db.Where("id = ?", id).First(&widget).Error
	if err != nil {
		return nil, err
	}
	return &widget, nil
}

// UpdateWidget updates a widget
func (r *DashboardRepository) UpdateWidget(id uuid.UUID, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	return r.db.Model(&models.DashboardWidget{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteWidget deletes a widget
func (r *DashboardRepository) DeleteWidget(id uuid.UUID) error {
	return r.db.Delete(&models.DashboardWidget{}, id).Error
}

// DeleteWidgetsByDashboardID deletes all widgets for a dashboard
func (r *DashboardRepository) DeleteWidgetsByDashboardID(dashboardID uuid.UUID) error {
	return r.db.Where("dashboard_id = ?", dashboardID).Delete(&models.DashboardWidget{}).Error
}

// GetWidgetsByDashboardID retrieves all widgets for a dashboard
func (r *DashboardRepository) GetWidgetsByDashboardID(dashboardID uuid.UUID) ([]models.DashboardWidget, error) {
	var widgets []models.DashboardWidget
	err := r.db.
		Where("dashboard_id = ?", dashboardID).
		Order("y ASC, x ASC").
		Find(&widgets).Error
	return widgets, err
}

// ReorderWidgets reorders widgets in a dashboard
func (r *DashboardRepository) ReorderWidgets(dashboardID uuid.UUID, widgetIDs []uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for i, widgetID := range widgetIDs {
			if err := tx.Model(&models.DashboardWidget{}).
				Where("id = ? AND dashboard_id = ?", widgetID, dashboardID).
				Update("order", i).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// Share operations

// CreateShare creates a new dashboard share
func (r *DashboardRepository) CreateShare(share *models.DashboardShare) error {
	return r.db.Create(share).Error
}

// GetSharesByDashboardID retrieves all shares for a dashboard
func (r *DashboardRepository) GetSharesByDashboardID(dashboardID uuid.UUID) ([]models.DashboardShare, error) {
	var shares []models.DashboardShare
	err := r.db.
		Where("dashboard_id = ?", dashboardID).
		Find(&shares).Error
	return shares, err
}

// DeleteShare deletes a share
func (r *DashboardRepository) DeleteShare(id uuid.UUID) error {
	return r.db.Delete(&models.DashboardShare{}, id).Error
}

// CheckSharePermission checks if a user has permission to access a dashboard
func (r *DashboardRepository) CheckSharePermission(dashboardID, userID uuid.UUID) (string, error) {
	var share models.DashboardShare
	err := r.db.
		Where("dashboard_id = ? AND shared_with = ?", dashboardID, userID).
		First(&share).Error
	if err != nil {
		// Check if user is the owner
		var dashboard models.Dashboard
		err := r.db.Where("id = ? AND owner_id = ?", dashboardID, userID).First(&dashboard).Error
		if err == nil {
			return "owner", nil
		}
		return "", err
	}
	return share.Permission, nil
}

// DuplicateDashboard duplicates a dashboard with all its widgets
func (r *DashboardRepository) DuplicateDashboard(dashboardID uuid.UUID, newOwnerID uuid.UUID, newName string) (*models.Dashboard, error) {
	var original models.Dashboard
	var newDashboard *models.Dashboard

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Get original dashboard with widgets
		if err := tx.Preload("Widgets").Where("id = ?", dashboardID).First(&original).Error; err != nil {
			return err
		}

		// Create new dashboard
		newDashboard = &models.Dashboard{
			OrganizationID: original.OrganizationID,
			OwnerID:        newOwnerID,
			Name:           newName,
			Description:    original.Description,
			Config:         original.Config,
			TemplateID:     original.TemplateID,
			CreatedBy:      newOwnerID,
		}
		if err := tx.Create(newDashboard).Error; err != nil {
			return err
		}

		// Duplicate widgets
		for _, widget := range original.Widgets {
			newWidget := models.DashboardWidget{
				DashboardID: newDashboard.ID,
				WidgetType:  widget.WidgetType,
				Title:       widget.Title,
				X:           widget.X,
				Y:           widget.Y,
				Width:       widget.Width,
				Height:      widget.Height,
				Config:      widget.Config,
				Order:       widget.Order,
			}
			if err := tx.Create(&newWidget).Error; err != nil {
				return err
			}
		}

		return nil
	})

	return newDashboard, err
}

// Template operations

// GetTemplates retrieves all dashboard templates
func (r *DashboardRepository) GetTemplates(category string) ([]models.DashboardTemplate, error) {
	var templates []models.DashboardTemplate

	query := r.db.Where("is_default = ?", true)
	if category != "" {
		query = query.Where("category = ?", category)
	}

	err := query.Find(&templates).Error
	return templates, err
}

// GetTemplateByID retrieves a template by ID
func (r *DashboardRepository) GetTemplateByID(id uuid.UUID) (*models.DashboardTemplate, error) {
	var template models.DashboardTemplate
	err := r.db.Where("id = ?", id).First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}
