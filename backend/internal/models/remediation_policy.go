package models

import (
	"time"

	"github.com/google/uuid"
)

type RemediationPolicy struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name             string    `gorm:"size:255;not null" json:"name"`
	AlertType        string    `gorm:"size:100;index" json:"alert_type"`    // CPU, Memory, Disk, Service, Docker, Kubernetes, etc.
	Severity         string    `gorm:"size:30;index" json:"severity"`       // Critical, Major, Warning, Info
	ActionType       string    `gorm:"size:50;not null" json:"action_type"` // restart_service, restart_container, restart_pod, cleanup_disk, kill_process, scale_deployment, rollback_deployment, custom_script
	Command          string    `gorm:"type:text" json:"command"`
	Enabled          bool      `gorm:"default:true" json:"enabled"`
	RequiresApproval bool      `gorm:"default:false" json:"requires_approval"`
	RetryCount       int       `gorm:"default:3" json:"retry_count"`
	TimeoutSec       int       `gorm:"default:60" json:"timeout_sec"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
