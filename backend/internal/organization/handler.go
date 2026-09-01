package organization

import (
	"net/http"
	"strconv"
	"time"

	"infrapilot/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateOrganization(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
		Slug string `json:"slug"`
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

func (h *Handler) ListOrganizations(c *gin.Context) {
	userID := c.GetString("userId")
	role := c.GetString("role")

	orgs, err := h.service.ListOrganizations(userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"organizations": orgs})
}

func (h *Handler) GetOrganization(c *gin.Context) {
	idOrSlug := c.Param("id")
	org, err := h.service.GetOrganization(idOrSlug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}
	c.JSON(http.StatusOK, org)
}

func (h *Handler) UpdateOrganization(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

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

func (h *Handler) DeleteOrganization(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	if err := h.service.DeleteOrganization(orgID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Organization suspended/archived successfully"})
}

func (h *Handler) GetDashboard(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	metrics, err := h.service.GetDashboardMetrics(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

func (h *Handler) GetSettings(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, _ := uuid.Parse(orgIDStr)

	settings, err := h.service.repo.GetSettings(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, settings)
}

func (h *Handler) UpdateSettings(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, _ := uuid.Parse(orgIDStr)

	var req models.OrganizationSettings
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.OrganizationID = orgID
	if err := h.service.repo.UpdateSettings(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

func (h *Handler) GetQuota(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, _ := uuid.Parse(orgIDStr)

	report, err := h.service.GetQuotaEngine().GetQuotaReport(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"quota_report": report})
}

func (h *Handler) UpdateQuota(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, _ := uuid.Parse(orgIDStr)

	var req models.OrganizationQuota
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.OrganizationID = orgID
	if err := h.service.repo.UpdateQuotas(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

func (h *Handler) GetAuditLogs(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, _ := uuid.Parse(orgIDStr)

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	logs, total, err := h.service.audit.GetLogs(orgID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"audit_logs": logs,
		"total":      total,
	})
}

func (h *Handler) InviteUser(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, _ := uuid.Parse(orgIDStr)

	var req struct {
		Email string `json:"email" binding:"required,email"`
		Role  string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	invitedBy := c.GetString("username")
	inv, err := h.service.CreateInvitation(orgID, req.Email, req.Role, invitedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, inv)
}

func (h *Handler) GetEnrollmentTokens(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, _ := uuid.Parse(orgIDStr)

	tokens, err := h.service.repo.GetEnrollmentTokens(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"enrollment_tokens": tokens})
}

func (h *Handler) CreateEnrollmentToken(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, _ := uuid.Parse(orgIDStr)

	var req struct {
		Name    string `json:"name"`
		MaxUses int    `json:"max_uses"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org, _ := h.service.repo.GetByID(orgID)
	slug := "tenant"
	if org != nil {
		slug = org.Slug
	}

	tokenStr := h.service.GenerateEnrollmentTokenString(slug)
	token := &models.EnrollmentToken{
		ID:             uuid.New(),
		OrganizationID: orgID.String(),
		Name:           req.Name,
		TokenPrefix:    "INFRA-" + slug,
		Token:          tokenStr,
		MaxUses:        req.MaxUses,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := h.service.repo.CreateEnrollmentToken(token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, token)
}
