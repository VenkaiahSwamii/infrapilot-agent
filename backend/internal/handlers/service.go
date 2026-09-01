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

type ServiceActionRequest struct {
	MachineID string `json:"machine_id" binding:"required"`
	Service   string `json:"service" binding:"required"`
}

// GetMachineServices returns the service status metrics collected for a machine
func GetMachineServices(c *gin.Context) {
	machineID := c.Param("id")
	var machineUUID uuid.UUID
	var err error
	if machineUUID, err = uuid.Parse(machineID); err != nil {
		var machine models.Machine
		if database.DB.Where("id::text LIKE ? OR hostname = ?", machineID+"%", machineID).First(&machine).Error == nil {
			machineUUID = machine.ID
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid machine ID"})
			return
		}
	}

	var services []models.LinuxService
	if err := database.DB.Where("machine_id = ?", machineUUID).Order("name ASC").Find(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query services"})
		return
	}

	type ServiceItem struct {
		Name   string `json:"name"`
		Status string `json:"status"`
	}

	response := make([]ServiceItem, 0, len(services))
	for _, s := range services {
		response = append(response, ServiceItem{
			Name:   s.Name,
			Status: s.Status,
		})
	}
	c.JSON(http.StatusOK, response)
}

func ServiceStart(c *gin.Context) {
	handleServiceAction(c, "start")
}

func ServiceStop(c *gin.Context) {
	handleServiceAction(c, "stop")
}

func ServiceRestart(c *gin.Context) {
	handleServiceAction(c, "restart")
}

func handleServiceAction(c *gin.Context, action string) {
	var req ServiceActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var machineUUID uuid.UUID
	var err error
	if machineUUID, err = uuid.Parse(req.MachineID); err != nil {
		var machine models.Machine
		if database.DB.Where("id::text LIKE ? OR hostname = ?", req.MachineID+"%", req.MachineID).First(&machine).Error == nil {
			machineUUID = machine.ID
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid machine ID"})
			return
		}
	}

	var machine models.Machine
	if err := database.DB.First(&machine, "id = ?", machineUUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "machine not found"})
		return
	}

	isWindows := strings.Contains(strings.ToLower(machine.OS), "win") || strings.Contains(strings.ToLower(machine.Platform), "win")

	var shellCmd string
	switch action {
	case "start":
		if isWindows {
			shellCmd = fmt.Sprintf("net start \"%s\"", req.Service)
		} else {
			shellCmd = fmt.Sprintf("sudo systemctl start %s", req.Service)
		}
	case "stop":
		if isWindows {
			shellCmd = fmt.Sprintf("net stop \"%s\"", req.Service)
		} else {
			shellCmd = fmt.Sprintf("sudo systemctl stop %s", req.Service)
		}
	case "restart":
		if isWindows {
			shellCmd = fmt.Sprintf("powershell -Command \"Restart-Service -Name '%s'\"", req.Service)
		} else {
			shellCmd = fmt.Sprintf("sudo systemctl restart %s", req.Service)
		}
	}

	cmd := models.Command{
		ID:        uuid.New(),
		MachineID: machineUUID,
		Command:   shellCmd,
		Status:    "Pending",
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&cmd).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enqueue service execution"})
		return
	}

	// Retrieve username from auth context claims
	usernameVal, exists := c.Get("userId")
	username := "admin"
	if exists {
		var user models.User
		if database.DB.First(&user, "id = ?", usernameVal).Error == nil {
			username = user.Username
		}
	}

	// Save log in audit_logs table
	audit := models.AuditLog{
		ID:        uuid.New(),
		Username:  username,
		MachineID: machineUUID,
		Action:    fmt.Sprintf("%s service %s", strings.Title(action), req.Service),
		CreatedAt: time.Now(),
	}
	database.DB.Create(&audit)

	c.JSON(http.StatusOK, gin.H{
		"message":    fmt.Sprintf("Service %s request queued", action),
		"command_id": cmd.ID,
		"status":     "queued",
	})
}

// Duplicate GetAuditLogs removed; use handlers/audit.go implementation
