package handlers

import (
	"net/http"

	"infrapilot/backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HostSecurityHandler struct {
	securityService *services.HostSecurityService
}

func NewHostSecurityHandler(securityService *services.HostSecurityService) *HostSecurityHandler {
	return &HostSecurityHandler{securityService: securityService}
}

type ExecuteRemediationRequest struct {
	PlaybookID string `json:"playbook_id" binding:"required"`
}

// GetSecurityAudit runs or returns CIS security audit report for a host
func (h *HostSecurityHandler) GetSecurityAudit(c *gin.Context) {
	idStr := c.Param("id")
	machineID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid machine_id UUID format"})
		return
	}

	report, err := h.securityService.RunHostSecurityScan(machineID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to run host security audit: " + err.Error()})
		return
	}

	playbooks := h.securityService.GetAvailablePlaybooks()

	c.JSON(http.StatusOK, gin.H{
		"report":    report,
		"playbooks": playbooks,
	})
}

// RunSecurityScan triggers a fresh security scan
func (h *HostSecurityHandler) RunSecurityScan(c *gin.Context) {
	h.GetSecurityAudit(c)
}

// ListPlaybooks returns available remediation playbooks
func (h *HostSecurityHandler) ListPlaybooks(c *gin.Context) {
	playbooks := h.securityService.GetAvailablePlaybooks()
	c.JSON(http.StatusOK, gin.H{"playbooks": playbooks})
}

// ExecuteRemediation executes a 1-click remediation playbook on a target machine
func (h *HostSecurityHandler) ExecuteRemediation(c *gin.Context) {
	idStr := c.Param("id")
	machineID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid machine_id UUID format"})
		return
	}

	var req ExecuteRemediationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	usernameVal, _ := c.Get("username")
	username, _ := usernameVal.(string)
	if username == "" {
		username = "Admin User"
	}

	res, err := h.securityService.ExecuteRemediationPlaybook(machineID, req.PlaybookID, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute playbook: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Playbook executed successfully",
		"result":  res,
	})
}
