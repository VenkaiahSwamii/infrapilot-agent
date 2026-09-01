package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NetworkScanRequest struct {
	Subnet string `json:"subnet" binding:"required"`
}

type DiscoveredMachine struct {
	IPAddress string `json:"ip_address"`
	Hostname  string `json:"hostname"`
	OS        string `json:"os"`
	Status    string `json:"status"`
	OpenPorts []int  `json:"open_ports"`
}

type ApproveMachineRequest struct {
	IPAddress string `json:"ip_address" binding:"required"`
	Hostname  string `json:"hostname" binding:"required"`
	OS        string `json:"os" binding:"required"`
}

// ScanNetwork scans a subnet range and returns discovered machines
func ScanNetwork(c *gin.Context) {
	var req NetworkScanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNotImplemented, gin.H{
		"error":  "network discovery is unavailable until a real scanner is configured",
		"subnet": req.Subnet,
	})
}

// ApproveMachine enrolls a discovered machine into the database
func ApproveMachine(c *gin.Context) {
	var req ApproveMachineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	machine := models.Machine{
		ID:           uuid.New(),
		Hostname:     req.Hostname,
		IPAddress:    req.IPAddress,
		OS:           strings.ToLower(req.OS),
		Platform:     req.OS,
		Status:       "OFFLINE",
		Online:       false,
		ResourceType: "virtual_machine",
		Organization: "Default Organization",
		APIKey:       "ip_live_" + strings.ReplaceAll(uuid.New().String(), "-", ""),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := database.DB.Create(&machine).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to approve and register machine"})
		return
	}

	usernameVal, exists := c.Get("userId")
	username := "admin"
	if exists {
		var user models.User
		if database.DB.First(&user, "id = ?", usernameVal).Error == nil {
			username = user.Username
		}
	}

	audit := models.AuditLog{
		ID:        uuid.New(),
		Username:  username,
		MachineID: machine.ID,
		Action:    fmt.Sprintf("Approve discovered machine %s (%s)", req.Hostname, req.IPAddress),
		CreatedAt: time.Now(),
	}
	database.DB.Create(&audit)

	scheme := "http"
	if c.Request.TLS != nil || strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https") {
		scheme = "https"
	}
	host := c.Request.Host
	if host == "" {
		host = "localhost:8080"
	}

	installCmd := fmt.Sprintf("curl -sSL %s://%s/downloads/install.sh | sudo bash -s %s://%s", scheme, host, scheme, host)
	if strings.Contains(strings.ToLower(req.OS), "win") {
		installCmd = fmt.Sprintf("irm %s://%s/downloads/install.ps1 | iex", scheme, host)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         fmt.Sprintf("Machine %s approved for management.", req.Hostname),
		"machine_id":      machine.ID.String(),
		"api_key":         machine.APIKey,
		"install_command": installCmd,
	})
}
