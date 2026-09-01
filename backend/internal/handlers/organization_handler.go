package handlers

import (
	"net/http"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"
	"infrapilot/backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrganizationHandler struct {
	service *services.OrganizationService
}

func NewOrganizationHandler(service *services.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{
		service: service,
	}
}

// ListOrganizations lists tenant organizations accessible to the user
func (h *OrganizationHandler) ListOrganizations(c *gin.Context) {
	userID := c.GetString("userId")
	role := c.GetString("role")

	orgs, err := h.service.ListOrganizations(userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"organizations": orgs})
}

// CreateOrganization creates a new tenant organization
func (h *OrganizationHandler) CreateOrganization(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
		Slug string `json:"slug" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ownerIDStr := c.GetString("userId")
	ownerID, _ := uuid.Parse(ownerIDStr)

	org, err := h.service.CreateOrganization(req.Name, req.Slug, ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, org)
}

// GetOrganization retrieves tenant organization details
func (h *OrganizationHandler) GetOrganization(c *gin.Context) {
	orgID := c.Param("orgId")
	org, err := h.service.GetOrganization(orgID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}
	c.JSON(http.StatusOK, org)
}

// UpdateOrganization updates tenant organization parameters
func (h *OrganizationHandler) UpdateOrganization(c *gin.Context) {
	orgID := c.Param("orgId")
	var req struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org, err := h.service.UpdateOrganization(orgID, req.Name, req.Slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, org)
}

// DeleteOrganization purges a tenant organization
func (h *OrganizationHandler) DeleteOrganization(c *gin.Context) {
	orgID := c.Param("orgId")
	if err := h.service.DeleteOrganization(orgID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Organization deleted successfully"})
}

// ListOrganizationUsers retrieves users belonging to a tenant organization
func (h *OrganizationHandler) ListOrganizationUsers(c *gin.Context) {
	orgIDStr := c.Param("orgId")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	type MemberResult struct {
		ID        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		Role      string    `json:"role"`
		CreatedAt time.Time `json:"created_at"`
	}

	var members []MemberResult
	if database.DB != nil {
		database.DB.Table("organization_users").
			Select("users.id, users.name, users.email, organization_users.role, organization_users.created_at").
			Joins("join users on users.id = organization_users.user_id").
			Where("organization_users.organization_id = ?", orgID).
			Scan(&members)
	}

	c.JSON(http.StatusOK, gin.H{"members": members})
}

// GetOrganizationSettings retrieves tenant white-label settings
func (h *OrganizationHandler) GetOrganizationSettings(c *gin.Context) {
	orgIDStr := c.Param("orgId")
	orgID, _ := uuid.Parse(orgIDStr)

	settings, err := h.service.GetSettings(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, settings)
}

// UpdateOrganizationSettings updates white-label branding and policies
func (h *OrganizationHandler) UpdateOrganizationSettings(c *gin.Context) {
	orgIDStr := c.Param("orgId")
	orgID, _ := uuid.Parse(orgIDStr)

	var req models.OrganizationSettings
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	settings, err := h.service.UpdateSettings(orgID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, settings)
}

// GetOrganizationQuotas retrieves tenant resource limits
func (h *OrganizationHandler) GetOrganizationQuotas(c *gin.Context) {
	orgIDStr := c.Param("orgId")
	orgID, _ := uuid.Parse(orgIDStr)

	quota, err := h.service.GetQuotas(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, quota)
}

// UpdateOrganizationQuotas updates tenant quota limits
func (h *OrganizationHandler) UpdateOrganizationQuotas(c *gin.Context) {
	orgIDStr := c.Param("orgId")
	orgID, _ := uuid.Parse(orgIDStr)

	var req models.OrganizationQuota
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.OrganizationID = orgID
	if database.DB != nil {
		database.DB.Where("organization_id = ?", orgID).Assign(req).FirstOrCreate(&req)
	}

	c.JSON(http.StatusOK, req)
}

// GetOrganizationBilling retrieves subscription and billing data
func (h *OrganizationHandler) GetOrganizationBilling(c *gin.Context) {
	orgIDStr := c.Param("orgId")
	orgID, _ := uuid.Parse(orgIDStr)

	var billing models.OrganizationBilling
	if database.DB != nil {
		database.DB.Where("organization_id = ?", orgID).First(&billing)
	}

	if billing.ID == uuid.Nil {
		billing = models.OrganizationBilling{
			OrganizationID:  orgID,
			Plan:            "enterprise",
			Status:          "active",
			PaymentProvider: "stripe",
			RenewalDate:     time.Now().AddDate(1, 0, 0),
		}
	}

	c.JSON(http.StatusOK, billing)
}

// CreateInvitation generates team email invitations
func (h *OrganizationHandler) CreateInvitation(c *gin.Context) {
	orgIDStr := c.Param("orgId")
	orgID, _ := uuid.Parse(orgIDStr)

	var req struct {
		Email string `json:"email" binding:"required,email"`
		Role  string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Role == "" {
		req.Role = "operator"
	}

	invitedBy := c.GetString("username")
	inv, err := h.service.CreateInvitation(orgID, req.Email, req.Role, invitedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, inv)
}

// GetSuperAdminMetrics retrieves platform-wide statistics for Super Admin
func (h *OrganizationHandler) GetSuperAdminMetrics(c *gin.Context) {
	metrics := h.service.GetSuperAdminMetrics()
	c.JSON(http.StatusOK, metrics)
}
