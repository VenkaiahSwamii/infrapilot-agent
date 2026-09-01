package models

import (
	"time"

	"github.com/google/uuid"
)

type RemediationJob struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	IncidentID   uuid.UUID  `gorm:"type:uuid;index" json:"incident_id"`
	MachineID    uuid.UUID  `gorm:"type:uuid;index;not null" json:"machine_id"`
	PolicyID     uuid.UUID  `gorm:"type:uuid;index" json:"policy_id"`
	ActionType   string     `gorm:"size:50" json:"action_type"`
	Command      string     `gorm:"type:text" json:"command"`
	Status       string     `gorm:"size:30;index;default:'PENDING'" json:"status"` // PENDING, RUNNING, SUCCESS, FAILED, WAITING_APPROVAL
	RetryAttempt int        `gorm:"default:0" json:"retry_attempt"`
	MaxRetries   int        `gorm:"default:3" json:"max_retries"`
	Output       string     `gorm:"type:text" json:"output"`
	Error        string     `gorm:"type:text" json:"error"`
	StartedAt    time.Time  `json:"started_at"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
