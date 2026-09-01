package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SIEMAlert struct {
	ID        string    `json:"id"`
	Machine   string    `json:"machine"`
	Threat    string    `json:"threat"`
	Severity  string    `json:"severity"`
	SourceIP  string    `json:"source_ip"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

type PluginExtension struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Version     string  `json:"version"`
	Description string  `json:"description"`
	Rating      float64 `json:"rating"`
	Installed   bool    `json:"installed"`
	Category    string  `json:"category"`
}

type WebhookRequest struct {
	URL    string `json:"url" binding:"required"`
	Event  string `json:"event" binding:"required"`
	Active bool   `json:"active"`
}

// GetSIEMAlerts returns centralized security threats and intrusion indicators
func GetSIEMAlerts(c *gin.Context) {
	alerts := []SIEMAlert{
		{
			ID:        uuid.New().String(),
			Machine:   "Win-Workstation",
			Threat:    "Brute-Force Attack: 25 failed login attempts in 10 seconds",
			Severity:  "Critical",
			SourceIP:  "192.168.1.145",
			Status:    "Active",
			Timestamp: time.Now().Add(-4 * time.Minute),
		},
		{
			ID:        uuid.New().String(),
			Machine:   "Ubuntu-Web-01",
			Threat:    "File Integrity Violation: Modification detected in /etc/passwd",
			Severity:  "Critical",
			SourceIP:  "Localhost",
			Status:    "Investigating",
			Timestamp: time.Now().Add(-18 * time.Minute),
		},
		{
			ID:        uuid.New().String(),
			Machine:   "Postgres-Replica-01",
			Threat:    "Unauthorized Software Detected: Crypto miner running in background",
			Severity:  "Warning",
			SourceIP:  "Localhost",
			Status:    "Active",
			Timestamp: time.Now().Add(-45 * time.Minute),
		},
	}

	c.JSON(http.StatusOK, alerts)
}

// ListPlugins returns installed plugins and marketplace extensions
func ListPlugins(c *gin.Context) {
	plugins := []PluginExtension{
		{
			ID:          "docker-monitor",
			Name:        "Docker Container Monitor",
			Version:     "1.4.2",
			Description: "Probes, lists and controls Docker containers and local images from the dashboard",
			Rating:      4.8,
			Installed:   true,
			Category:    "Containers",
		},
		{
			ID:          "kubernetes-discovery",
			Name:        "Kubernetes Cluster Discovery",
			Version:     "2.1.0",
			Description: "Exposes namespaces, nodes, pods, and schedules container deployments",
			Rating:      4.9,
			Installed:   false,
			Category:    "Containers",
		},
		{
			ID:          "slack-notifications",
			Name:        "Slack Notifications Bridge",
			Version:     "1.0.5",
			Description: "Dispatches system metric threshold breaches and incidents to Slack channels",
			Rating:      4.6,
			Installed:   true,
			Category:    "Alerts",
		},
		{
			ID:          "aws-cost-monitor",
			Name:        "AWS Resource Cost Optimizer",
			Version:     "0.9.8",
			Description: "Monitors EC2, S3, RDS utilization and calculates cost reduction suggestions",
			Rating:      4.3,
			Installed:   false,
			Category:    "Cloud",
		},
	}

	c.JSON(http.StatusOK, plugins)
}

// InstallPlugin handles plugin installation triggers
func InstallPlugin(c *gin.Context) {
	pluginID := c.Param("pluginId")
	c.JSON(http.StatusOK, gin.H{
		"message":   fmt.Sprintf("Plugin %s installed successfully", pluginID),
		"plugin_id": pluginID,
		"status":    "Installed",
		"time":      time.Now().Format(time.RFC3339),
	})
}

// RegisterWebhook stores outbound event endpoints
func RegisterWebhook(c *gin.Context) {
	var req WebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Outbound webhook registered successfully",
		"webhook_id": uuid.New().String(),
		"url":        req.URL,
		"event":      req.Event,
		"active":     true,
	})
}
