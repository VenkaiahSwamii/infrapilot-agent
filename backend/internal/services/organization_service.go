package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/logger"
	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

type OrganizationService struct {
	logger *logger.Logger
}

func NewOrganizationService() *OrganizationService {
	return &OrganizationService{
		logger: logger.Get(),
	}
}

// CreateOrganization initializes a new tenant organization with default settings, quotas, and billing records
func (s *OrganizationService) CreateOrganization(name, slug string, ownerID uuid.UUID) (*models.Organization, error) {
	if database.DB == nil {
		return nil, fmt.Errorf("database connection unavailable")
	}

	orgID := uuid.New()
	org := &models.Organization{
		ID:        orgID,
		Name:      name,
		Slug:      slug,
		Status:    "active",
		OwnerID:   ownerID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := database.DB.Create(org).Error; err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	// 1. Add owner to organization_users
	if ownerID != uuid.Nil {
		member := &models.OrganizationUser{
			ID:             uuid.New(),
			OrganizationID: orgID,
			UserID:         ownerID,
			Role:           "owner",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		database.DB.Create(member)
	}

	// 2. Create default white-label settings
	settings := &models.OrganizationSettings{
		ID:             uuid.New(),
		OrganizationID: orgID,
		CompanyName:    name,
		Timezone:       "UTC",
		Language:       "en",
		Theme:          "dark",
		PrimaryColor:   "#58a6ff",
		SidebarColor:   "#0d1117",
		LicenseTier:    "enterprise",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	database.DB.Create(settings)

	// 3. Create default quotas
	quota := &models.OrganizationQuota{
		ID:             uuid.New(),
		OrganizationID: orgID,
		MaxMachines:    50,
		MaxUsers:       20,
		MaxStorageGB:   500,
		MaxReports:     100,
		MaxAlerts:      1000,
		MaxAPICalls:    100000,
		RetentionDays:  90,
		MaxProjects:    10,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	database.DB.Create(quota)

	// 4. Create default billing record
	billing := &models.OrganizationBilling{
		ID:              uuid.New(),
		OrganizationID:  orgID,
		Plan:            "enterprise",
		Status:          "active",
		PaymentProvider: "stripe",
		RenewalDate:     time.Now().AddDate(1, 0, 0),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	database.DB.Create(billing)

	s.logger.Info("Organization created successfully", "org_id", orgID.String(), "name", name)
	return org, nil
}

// GetOrganization retrieves an organization by ID or Slug
func (s *OrganizationService) GetOrganization(identifier string) (*models.Organization, error) {
	var org models.Organization
	if database.DB == nil {
		return nil, fmt.Errorf("database unavailable")
	}

	if err := database.DB.Where("id = ? OR slug = ?", identifier, identifier).First(&org).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

// ListOrganizations lists all organizations (or filtered by user)
func (s *OrganizationService) ListOrganizations(userID string, userRole string) ([]models.Organization, error) {
	var orgs []models.Organization
	if database.DB == nil {
		return orgs, nil
	}

	if userRole == string(models.RoleSuperAdmin) || userRole == "super_admin" || userID == "" {
		database.DB.Order("created_at desc").Find(&orgs)
		return orgs, nil
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return orgs, nil
	}

	database.DB.Table("organizations").
		Joins("JOIN organization_users ON organization_users.organization_id = organizations.id").
		Where("organization_users.user_id = ? AND organizations.status = 'active'", parsedUserID).
		Find(&orgs)

	return orgs, nil
}

// UpdateOrganization updates organization details
func (s *OrganizationService) UpdateOrganization(orgID string, name, slug string) (*models.Organization, error) {
	org, err := s.GetOrganization(orgID)
	if err != nil {
		return nil, err
	}

	org.Name = name
	if slug != "" {
		org.Slug = slug
	}
	org.UpdatedAt = time.Now()

	if err := database.DB.Save(org).Error; err != nil {
		return nil, err
	}
	return org, nil
}

// DeleteOrganization removes an organization
func (s *OrganizationService) DeleteOrganization(orgID string) error {
	org, err := s.GetOrganization(orgID)
	if err != nil {
		return err
	}

	database.DB.Where("organization_id = ?", org.ID).Delete(&models.OrganizationUser{})
	database.DB.Where("organization_id = ?", org.ID).Delete(&models.OrganizationSettings{})
	database.DB.Where("organization_id = ?", org.ID).Delete(&models.OrganizationQuota{})
	database.DB.Where("organization_id = ?", org.ID).Delete(&models.OrganizationBilling{})
	database.DB.Delete(org)
	return nil
}

// GetSettings retrieves white-label branding and settings
func (s *OrganizationService) GetSettings(orgID uuid.UUID) (*models.OrganizationSettings, error) {
	var settings models.OrganizationSettings
	if database.DB == nil {
		return &settings, nil
	}

	if err := database.DB.Where("organization_id = ?", orgID).First(&settings).Error; err != nil {
		// Return default if not found
		return &models.OrganizationSettings{
			OrganizationID: orgID,
			CompanyName:    "InfraPilot Enterprise",
			PrimaryColor:   "#58a6ff",
			SidebarColor:   "#0d1117",
			LicenseTier:    "enterprise",
		}, nil
	}
	return &settings, nil
}

// UpdateSettings updates white-label branding and policies
func (s *OrganizationService) UpdateSettings(orgID uuid.UUID, updated *models.OrganizationSettings) (*models.OrganizationSettings, error) {
	var settings models.OrganizationSettings
	if database.DB == nil {
		return nil, fmt.Errorf("database unavailable")
	}

	if database.DB.Where("organization_id = ?", orgID).First(&settings).Error != nil {
		settings.ID = uuid.New()
		settings.OrganizationID = orgID
	}

	if updated.CompanyName != "" {
		settings.CompanyName = updated.CompanyName
	}
	if updated.LogoURL != "" {
		settings.LogoURL = updated.LogoURL
	}
	if updated.PrimaryColor != "" {
		settings.PrimaryColor = updated.PrimaryColor
	}
	if updated.SidebarColor != "" {
		settings.SidebarColor = updated.SidebarColor
	}
	if updated.CustomDomain != "" {
		settings.CustomDomain = updated.CustomDomain
	}
	if updated.LicenseTier != "" {
		settings.LicenseTier = updated.LicenseTier
	}
	settings.UpdatedAt = time.Now()

	database.DB.Save(&settings)
	return &settings, nil
}

// GetQuotas retrieves tenant resource limits
func (s *OrganizationService) GetQuotas(orgID uuid.UUID) (*models.OrganizationQuota, error) {
	var quota models.OrganizationQuota
	if database.DB == nil {
		return &quota, nil
	}
	database.DB.Where("organization_id = ?", orgID).First(&quota)
	return &quota, nil
}

// CheckQuotaLimit verifies if tenant has capacity to create resourceType
func (s *OrganizationService) CheckQuotaLimit(orgID uuid.UUID, resourceType string) (bool, error) {
	quota, err := s.GetQuotas(orgID)
	if err != nil || quota == nil || quota.ID == uuid.Nil {
		return true, nil // Unlimited by default
	}

	if database.DB == nil {
		return true, nil
	}

	switch resourceType {
	case "machines":
		if quota.MaxMachines == 0 {
			return true, nil
		}
		var count int64
		database.DB.Model(&models.Machine{}).Where("organization_id = ?", orgID.String()).Count(&count)
		return count < int64(quota.MaxMachines), nil
	case "users":
		if quota.MaxUsers == 0 {
			return true, nil
		}
		var count int64
		database.DB.Model(&models.OrganizationUser{}).Where("organization_id = ?", orgID).Count(&count)
		return count < int64(quota.MaxUsers), nil
	}

	return true, nil
}

// CreateInvitation generates a user invitation token
func (s *OrganizationService) CreateInvitation(orgID uuid.UUID, email, role, invitedBy string) (*models.OrganizationInvitation, error) {
	tokenBytes := make([]byte, 16)
	rand.Read(tokenBytes)
	inviteToken := hex.EncodeToString(tokenBytes)

	inv := &models.OrganizationInvitation{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Email:          email,
		Role:           role,
		InviteToken:    inviteToken,
		InvitedBy:      invitedBy,
		ExpiresAt:      time.Now().Add(7 * 24 * time.Hour),
		Accepted:       false,
		CreatedAt:      time.Now(),
	}

	if database.DB != nil {
		database.DB.Create(inv)
	}

	s.logger.Info("Organization invitation created", "email", email, "org_id", orgID.String())
	return inv, nil
}

// GetSuperAdminMetrics aggregates platform-wide statistics for Super Admin dashboard
func (s *OrganizationService) GetSuperAdminMetrics() map[string]interface{} {
	metrics := map[string]interface{}{
		"total_organizations": int64(0),
		"total_machines":      int64(0),
		"total_users":         int64(0),
		"active_agents":       int64(0),
		"storage_used_bytes":  int64(0),
		"monthly_revenue":     125000,
		"subscriptions":       map[string]int{"starter": 12, "professional": 8, "enterprise": 5},
	}

	if database.DB != nil {
		var orgCount, machineCount, userCount, agentCount int64
		database.DB.Model(&models.Organization{}).Count(&orgCount)
		database.DB.Model(&models.Machine{}).Count(&machineCount)
		database.DB.Model(&models.User{}).Count(&userCount)
		database.DB.Model(&models.Machine{}).Where("status = ?", "online").Count(&agentCount)

		metrics["total_organizations"] = orgCount
		metrics["total_machines"] = machineCount
		metrics["total_users"] = userCount
		metrics["active_agents"] = agentCount
	}

	return metrics
}
