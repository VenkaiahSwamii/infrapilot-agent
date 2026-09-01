package models

import (
	"time"

	"github.com/google/uuid"
)

type Workflow struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name             string         `gorm:"size:255;not null" json:"name"`
	Description      string         `gorm:"type:text" json:"description"`
	TriggerType      string         `gorm:"size:100;index" json:"trigger_type"` // Alert, Incident, Scheduled, Manual, MachineOffline
	Enabled          bool           `gorm:"default:true" json:"enabled"`
	Version          int            `gorm:"default:1" json:"version"`
	RequiresApproval bool           `gorm:"default:false" json:"requires_approval"`
	Steps            []WorkflowStep `gorm:"foreignKey:WorkflowID;constraint:OnDelete:CASCADE" json:"steps,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

type WorkflowStep struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	WorkflowID  uuid.UUID `gorm:"type:uuid;index;not null" json:"workflow_id"`
	StepOrder   int       `gorm:"not null" json:"step_order"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	ActionType  string    `gorm:"size:50;not null" json:"action_type"` // shell_command, ssh_command, restart_service, restart_container, restart_pod, scale_deployment, api_call, send_notification, wait, condition, collect_logs, create_ticket
	ActionValue string    `gorm:"type:text" json:"action_value"`
	TimeoutSec  int       `gorm:"default:60" json:"timeout_sec"`
	OnFailure   string    `gorm:"size:30;default:'STOP'" json:"on_failure"` // STOP, CONTINUE, RETRY
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type WorkflowExecution struct {
	ID          uuid.UUID     `gorm:"type:uuid;primaryKey" json:"id"`
	WorkflowID  uuid.UUID     `gorm:"type:uuid;index;not null" json:"workflow_id"`
	IncidentID  uuid.UUID     `gorm:"type:uuid;index" json:"incident_id"`
	MachineID   uuid.UUID     `gorm:"type:uuid;index" json:"machine_id"`
	Status      string        `gorm:"size:30;index;default:'RUNNING'" json:"status"` // RUNNING, SUCCESS, FAILED, CANCELLED
	CurrentStep int           `gorm:"default:1" json:"current_step"`
	StartedAt   time.Time     `json:"started_at"`
	CompletedAt *time.Time    `json:"completed_at,omitempty"`
	DurationMs  int64         `json:"duration_ms"`
	Logs        []WorkflowLog `gorm:"foreignKey:ExecutionID;constraint:OnDelete:CASCADE" json:"logs,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type WorkflowLog struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ExecutionID uuid.UUID `gorm:"type:uuid;index;not null" json:"execution_id"`
	StepName    string    `gorm:"size:255" json:"step_name"`
	Status      string    `gorm:"size:30" json:"status"` // SUCCESS, FAILED, SKIPPED
	Output      string    `gorm:"type:text" json:"output"`
	Error       string    `gorm:"type:text" json:"error"`
	CreatedAt   time.Time `json:"created_at"`
}

type WorkflowStepRequest struct {
	StepOrder   int    `json:"step_order"`
	Name        string `json:"name" binding:"required"`
	ActionType  string `json:"action_type" binding:"required"`
	ActionValue string `json:"action_value"`
	TimeoutSec  int    `json:"timeout_sec"`
	OnFailure   string `json:"on_failure"`
}
