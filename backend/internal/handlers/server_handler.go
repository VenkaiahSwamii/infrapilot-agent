package handlers

import (
	"net/http"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/events"
	"infrapilot/backend/internal/models"
	"infrapilot/backend/internal/services"
	"infrapilot/backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ServerHandler struct {
	serverService *services.ServerService
}

func NewServerHandler(serverService *services.ServerService) *ServerHandler {
	return &ServerHandler{serverService: serverService}
}

type ServerRegisterRequest struct {
	EnrollmentToken string `json:"enrollment_token"`
	MachineID       string `json:"machine_id"`
	Hostname        string `json:"hostname" binding:"required"`
	IP              string `json:"ip"`
	IPAddress       string `json:"ip_address"`
	OS              string `json:"os" binding:"required"`
	Platform        string `json:"platform"`
	AgentVersion    string `json:"agent_version"`
	ResourceType    string `json:"resource_type"`
	Organization    string `json:"organization"`
	Kernel          string `json:"kernel"`
	Architecture    string `json:"architecture"`
	MACAddress      string `json:"mac_address"`
	CPUModel        string `json:"cpu_model"`
	TotalMemoryGB   uint64 `json:"total_memory_gb"`
	TotalDiskGB     uint64 `json:"total_disk_gb"`
	GPU             string `json:"gpu"`
	Virtualization  string `json:"virtualization"`
	CloudProvider   string `json:"cloud_provider"`
}

func (h *ServerHandler) EnrollServer(c *gin.Context) {
	var req ServerRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := ValidateEnrollmentToken(req.EnrollmentToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired enrollment token"})
		return
	}

	ip := req.IP
	if ip == "" {
		ip = req.IPAddress
	}

	machineID := uuid.New()
	if req.MachineID != "" {
		parsed, parseErr := uuid.Parse(req.MachineID)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid machine_id format"})
			return
		}
		machineID = parsed
	}

	input := services.RegisterServerInput{
		ID:             machineID,
		Hostname:       req.Hostname,
		OS:             req.OS,
		Platform:       req.Platform,
		IPAddress:      ip,
		AgentVersion:   req.AgentVersion,
		ResourceType:   req.ResourceType,
		Organization:   token.OrganizationID,
		Kernel:         req.Kernel,
		Architecture:   req.Architecture,
		MACAddress:     req.MACAddress,
		CPUModel:       req.CPUModel,
		TotalMemoryGB:  req.TotalMemoryGB,
		TotalDiskGB:    req.TotalDiskGB,
		GPU:            req.GPU,
		Virtualization: req.Virtualization,
		CloudProvider:  req.CloudProvider,
	}

	server, err := h.serverService.RegisterOrUpdateServer(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enroll server"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"machine_id":      server.ID.String(),
		"api_key":         server.APIKey,
		"organization_id": token.OrganizationID,
		"config": gin.H{
			"metrics_interval_seconds":   5,
			"heartbeat_interval_seconds": 15,
			"log_collection_enabled":     true,
		},
	})
}

func (h *ServerHandler) RegisterServer(c *gin.Context) {
	var req ServerRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ip := req.IP
	if ip == "" {
		ip = req.IPAddress
	}

	if req.MachineID != "" {
		macID, err := uuid.Parse(req.MachineID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid machine_id format"})
			return
		}

		input := services.RegisterServerInput{
			ID:             macID,
			Hostname:       req.Hostname,
			OS:             req.OS,
			Platform:       req.Platform,
			IPAddress:      ip,
			AgentVersion:   req.AgentVersion,
			ResourceType:   req.ResourceType,
			Organization:   req.Organization,
			Kernel:         req.Kernel,
			Architecture:   req.Architecture,
			MACAddress:     req.MACAddress,
			CPUModel:       req.CPUModel,
			TotalMemoryGB:  req.TotalMemoryGB,
			TotalDiskGB:    req.TotalDiskGB,
			GPU:            req.GPU,
			Virtualization: req.Virtualization,
			CloudProvider:  req.CloudProvider,
		}

		server, err := h.serverService.RegisterOrUpdateServer(input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"machine_id": server.ID.String(),
			"api_key":    server.APIKey,
		})
		return
	}

	machineID := uuid.New()
	input := services.RegisterServerInput{
		ID:             machineID,
		Hostname:       req.Hostname,
		OS:             req.OS,
		Platform:       req.Platform,
		IPAddress:      ip,
		AgentVersion:   req.AgentVersion,
		ResourceType:   req.ResourceType,
		Organization:   req.Organization,
		Kernel:         req.Kernel,
		Architecture:   req.Architecture,
		MACAddress:     req.MACAddress,
		CPUModel:       req.CPUModel,
		TotalMemoryGB:  req.TotalMemoryGB,
		TotalDiskGB:    req.TotalDiskGB,
		GPU:            req.GPU,
		Virtualization: req.Virtualization,
		CloudProvider:  req.CloudProvider,
	}
	server, err := h.serverService.RegisterOrUpdateServer(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register server"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"machine_id": server.ID.String(),
		"api_key":    server.APIKey,
	})
}

func (h *ServerHandler) GetServers(c *gin.Context) {
	servers, err := h.serverService.GetServers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query servers"})
		return
	}
	c.JSON(http.StatusOK, servers)
}

func (h *ServerHandler) GetServerByID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing server ID"})
		return
	}

	server, err := h.serverService.GetServerSnapshotByIDOrHostname(idStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}
	c.JSON(http.StatusOK, server)
}

type UpdateServerRequest struct {
	Hostname string `json:"hostname"`
	Status   string `json:"status"`
}

func (h *ServerHandler) UpdateServer(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing server ID"})
		return
	}

	var req UpdateServerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	server, err := h.serverService.GetServerByIDOrHostname(idStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	if req.Hostname != "" {
		server.Hostname = req.Hostname
	}
	if req.Status != "" {
		server.Status = req.Status
	}

	if err := h.serverService.Save(server); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update server"})
		return
	}

	c.JSON(http.StatusOK, server)
}

func (h *ServerHandler) DeleteServer(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing server ID"})
		return
	}

	server, err := h.serverService.GetServerByIDOrHostname(idStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	err = database.DB.Delete(&models.Server{}, "id = ?", server.ID).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete server"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Server deleted successfully"})
}

func (h *ServerHandler) ServerHeartbeat(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing server ID"})
		return
	}

	server, err := h.serverService.GetServerByIDOrHostname(idStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	if err := h.serverService.Heartbeat(server.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *ServerHandler) RotateServerKey(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		idStr = c.Query("id")
	}

	var server *models.Server

	if idStr == "" {
		if machineVal, exists := c.Get("machine"); exists {
			if m, ok := machineVal.(*models.Machine); ok {
				server = m
			}
		}
		if server == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing server ID"})
			return
		}
	} else {
		var getErr error
		server, getErr = h.serverService.GetServerByIDOrHostname(idStr)
		if getErr != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
			return
		}
	}

	newKey := utils.GenerateAPIKey()
	server.APIKey = newKey
	server.KeyVersion++
	server.LastKeyRotate = time.Now()

	if err := h.serverService.Save(server); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to rotate key"})
		return
	}

	if h.serverService != nil {
		h.serverService.PublishEvent(events.APIKeyRotatedEvent{
			MachineID: server.ID.String(),
			UserID:    c.GetString("userId"),
			Time:      time.Now(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"machine_id": server.ID.String(),
		"api_key":    server.APIKey,
		"version":    server.KeyVersion,
		"rotate":     true,
	})
}

func (h *ServerHandler) GetServerKeyRotationStatus(c *gin.Context) {
	machineVal, exists := c.Get("machine")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	server, ok := machineVal.(*models.Server)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid server context"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rotate":     false,
		"api_key":    "",
		"version":    server.KeyVersion,
		"machine_id": server.ID.String(),
	})
}
