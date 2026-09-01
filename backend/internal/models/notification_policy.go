package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationPolicy struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Severity   string    `gorm:"size:30;index" json:"severity"` // Critical, Major, Warning, Info
	Channel    string    `gorm:"size:50" json:"channel"`        // Email, Slack, Telegram, Discord, Teams, Webhook
	Recipient  string    `json:"recipient"`
	WebhookURL string    `json:"webhook_url"`
	Enabled    bool      `gorm:"default:true" json:"enabled"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
