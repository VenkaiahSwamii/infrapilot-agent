package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	AlertID   uuid.UUID  `gorm:"type:uuid;index" json:"alert_id"`
	Channel   string     `gorm:"size:50;index" json:"channel"` // Email, Slack, Telegram, Discord, Teams, Webhook
	Recipient string     `json:"recipient"`
	Status    string     `gorm:"size:20;index" json:"status"` // SENT, FAILED, PENDING
	Error     string     `gorm:"type:text" json:"error"`
	SentAt    *time.Time `json:"sent_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}
