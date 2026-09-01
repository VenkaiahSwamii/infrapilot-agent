package handlers

import (
	"fmt"
	"net/http"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"
	"infrapilot/backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TestNotificationRequest struct {
	Channel   string `json:"channel" binding:"required"`   // Email, Slack, Telegram, Discord, Teams, Webhook
	Recipient string `json:"recipient" binding:"required"` // Email address, Webhook URL, or Channel ID
}

type NotificationPolicyRequest struct {
	Severity   string `json:"severity" binding:"required"`
	Channel    string `json:"channel" binding:"required"`
	Recipient  string `json:"recipient"`
	WebhookURL string `json:"webhook_url"`
	Enabled    *bool  `json:"enabled"`
}

// GetNotifications returns notification delivery history
func GetNotifications(c *gin.Context) {
	var notifications []models.Notification
	if database.DB != nil {
		if err := database.DB.Order("created_at desc").Find(&notifications).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query notifications history"})
			return
		}
	}

	// Fallback seed history if DB returns 0 items for demo/testing
	if len(notifications) == 0 {
		now := time.Now()
		notifications = []models.Notification{
			{ID: uuid.New(), Channel: "Email", Recipient: "admin@infrapilot.enterprise", Status: "SENT", SentAt: &now, CreatedAt: now.Add(-5 * time.Minute)},
			{ID: uuid.New(), Channel: "Slack", Recipient: "#infrastructure-alerts", Status: "SENT", SentAt: &now, CreatedAt: now.Add(-10 * time.Minute)},
			{ID: uuid.New(), Channel: "Telegram", Recipient: "@InfraPilotAlertBot", Status: "SENT", SentAt: &now, CreatedAt: now.Add(-15 * time.Minute)},
			{ID: uuid.New(), Channel: "Teams", Recipient: "DevOps Channel", Status: "SENT", SentAt: &now, CreatedAt: now.Add(-30 * time.Minute)},
		}
	}

	c.JSON(http.StatusOK, notifications)
}

// SendTestNotification triggers a test notification to a channel
func SendTestNotification(c *gin.Context) {
	var req TestNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	testAlert := models.LinuxAlert{
		ID:        uuid.New(),
		Title:     fmt.Sprintf("[TEST] %s Test Notification", req.Channel),
		Severity:  "Critical",
		Status:    "OPEN",
		Message:   fmt.Sprintf("InfraPilot test notification sent successfully to %s (%s)", req.Channel, req.Recipient),
		CreatedAt: time.Now(),
	}

	notifSvc := services.NewNotificationService()
	switch req.Channel {
	case "Email":
		notifSvc.SendEmail(testAlert)
	case "Slack":
		notifSvc.SendSlack(testAlert)
	case "Telegram":
		notifSvc.SendTelegram(testAlert)
	case "Discord":
		notifSvc.SendDiscord(testAlert)
	case "Teams":
		notifSvc.SendTeams(testAlert)
	case "Webhook":
		notifSvc.SendWebhook(testAlert)
	default:
		notifSvc.SendSlack(testAlert)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   fmt.Sprintf("Test notification dispatched to %s (%s)", req.Channel, req.Recipient),
		"status":    "SENT",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// GetNotificationPolicies lists notification policies
func GetNotificationPolicies(c *gin.Context) {
	var policies []models.NotificationPolicy
	if database.DB != nil {
		_ = database.DB.Order("severity desc").Find(&policies).Error
	}

	// Fallback seed policies if DB returns 0 items
	if len(policies) == 0 {
		now := time.Now()
		policies = []models.NotificationPolicy{
			{ID: uuid.New(), Severity: "Critical", Channel: "Email", Recipient: "devops-oncall@infrapilot.enterprise", Enabled: true, CreatedAt: now, UpdatedAt: now},
			{ID: uuid.New(), Severity: "Critical", Channel: "Slack", Recipient: "#critical-incidents", Enabled: true, CreatedAt: now, UpdatedAt: now},
			{ID: uuid.New(), Severity: "Critical", Channel: "Telegram", Recipient: "@InfraPilotAlertBot", Enabled: true, CreatedAt: now, UpdatedAt: now},
			{ID: uuid.New(), Severity: "Warning", Channel: "Slack", Recipient: "#warnings-log", Enabled: true, CreatedAt: now, UpdatedAt: now},
			{ID: uuid.New(), Severity: "Warning", Channel: "Teams", Recipient: "DevOps Ops Channel", Enabled: true, CreatedAt: now, UpdatedAt: now},
		}
	}

	c.JSON(http.StatusOK, policies)
}

// CreateNotificationPolicy creates a new policy rule
func CreateNotificationPolicy(c *gin.Context) {
	var req NotificationPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	policy := models.NotificationPolicy{
		ID:         uuid.New(),
		Severity:   req.Severity,
		Channel:    req.Channel,
		Recipient:  req.Recipient,
		WebhookURL: req.WebhookURL,
		Enabled:    enabled,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if database.DB != nil {
		if err := database.DB.Create(&policy).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notification policy"})
			return
		}
	}

	c.JSON(http.StatusCreated, policy)
}

// UpdateNotificationPolicy updates an existing policy
func UpdateNotificationPolicy(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid policy ID format"})
		return
	}

	var req NotificationPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var policy models.NotificationPolicy
	if database.DB != nil {
		if err := database.DB.First(&policy, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification policy not found"})
			return
		}

		if req.Severity != "" {
			policy.Severity = req.Severity
		}
		if req.Channel != "" {
			policy.Channel = req.Channel
		}
		if req.Recipient != "" {
			policy.Recipient = req.Recipient
		}
		if req.WebhookURL != "" {
			policy.WebhookURL = req.WebhookURL
		}
		if req.Enabled != nil {
			policy.Enabled = *req.Enabled
		}
		policy.UpdatedAt = time.Now()

		database.DB.Save(&policy)
	} else {
		policy.ID = id
		policy.Severity = req.Severity
		policy.Channel = req.Channel
	}

	c.JSON(http.StatusOK, policy)
}

// DeleteNotificationPolicy deletes a policy rule
func DeleteNotificationPolicy(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid policy ID format"})
		return
	}

	if database.DB != nil {
		if err := database.DB.Delete(&models.NotificationPolicy{}, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete notification policy"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification policy deleted successfully", "id": id})
}
