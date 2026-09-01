package cloud

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

func (h *Handler) ConnectAWS(c *gin.Context) {
	var req struct {
		OrganizationID string `json:"organization_id"`
		AccountID      string `json:"account_id" binding:"required"`
		AccountName    string `json:"account_name" binding:"required"`
		AccessKey      string `json:"access_key"`
		SecretKey      string `json:"secret_key"`
		RoleARN        string `json:"role_arn"`
		Regions        string `json:"regions"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	acc, err := h.service.ConnectAWS(req.OrganizationID, req.AccountID, req.AccountName, req.AccessKey, req.SecretKey, req.RoleARN, req.Regions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, acc)
}

func (h *Handler) ConnectAzure(c *gin.Context) {
	var req struct {
		OrganizationID string `json:"organization_id"`
		SubscriptionID string `json:"subscription_id" binding:"required"`
		AccountName    string `json:"account_name" binding:"required"`
		TenantID       string `json:"tenant_id"`
		ClientID       string `json:"client_id"`
		ClientSecret   string `json:"client_secret"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	acc, err := h.service.ConnectAzure(req.OrganizationID, req.SubscriptionID, req.AccountName, req.TenantID, req.ClientID, req.ClientSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, acc)
}

func (h *Handler) ConnectGCP(c *gin.Context) {
	var req struct {
		OrganizationID     string `json:"organization_id"`
		ProjectID          string `json:"project_id" binding:"required"`
		AccountName        string `json:"account_name" binding:"required"`
		ServiceAccountJSON string `json:"service_account_json"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	acc, err := h.service.ConnectGCP(req.OrganizationID, req.ProjectID, req.AccountName, req.ServiceAccountJSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, acc)
}

func (h *Handler) GetAccounts(c *gin.Context) {
	orgID := c.Query("organization_id")
	accounts, err := h.service.GetAccounts(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"accounts": accounts})
}

func (h *Handler) GetResources(c *gin.Context) {
	orgID := c.Query("organization_id")
	provider := c.Query("provider")

	if provider == "aws" {
		awsRes, _ := h.service.GetAWSResources(orgID)
		c.JSON(http.StatusOK, gin.H{"resources": awsRes})
		return
	}
	if provider == "azure" {
		azRes, _ := h.service.GetAzureResources(orgID)
		c.JSON(http.StatusOK, gin.H{"resources": azRes})
		return
	}
	if provider == "gcp" {
		gcpRes, _ := h.service.GetGCPResources(orgID)
		c.JSON(http.StatusOK, gin.H{"resources": gcpRes})
		return
	}

	hybridSummary, _ := h.service.GetHybridInventory(orgID)
	c.JSON(http.StatusOK, gin.H{"hybrid_summary": hybridSummary})
}

func (h *Handler) GetCosts(c *gin.Context) {
	orgID := c.Query("organization_id")
	costs, err := h.service.GetCosts(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, costs)
}

func (h *Handler) GetSecurity(c *gin.Context) {
	orgID := c.Query("organization_id")
	findings, err := h.service.GetSecurity(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"security_findings": findings})
}

func (h *Handler) GetPerformance(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ec2_cpu_percent":    42.8,
		"rds_cpu_percent":    38.4,
		"lambda_invocations": 128450,
		"cluster_health":     "HEALTHY",
	})
}

func (h *Handler) GetAlerts(c *gin.Context) {
	orgID := c.Query("organization_id")
	alerts, err := h.service.GetAlerts(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"alerts": alerts})
}
