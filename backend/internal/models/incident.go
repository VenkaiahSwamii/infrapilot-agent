package models

import (
	"time"

	"github.com/google/uuid"
)

// IncidentStatus represents the lifecycle stages of an incident
type IncidentStatus string

const (
	IncidentStatusOpen         IncidentStatus = "OPEN"
	IncidentStatusACKNOWLEDGED IncidentStatus = "ACKNOWLEDGED"
	IncidentStatusInvestigate  IncidentStatus = "INVESTIGATING"
	IncidentStatusMitigated    IncidentStatus = "MITIGATED"
	IncidentStatusResolved     IncidentStatus = "RESOLVED"
	IncidentStatusClosed       IncidentStatus = "CLOSED"
)

// IncidentSeverity represents incident severity levels
type IncidentSeverity string

const (
	IncidentSeverityP1 IncidentSeverity = "P1" // Critical
	IncidentSeverityP2 IncidentSeverity = "P2" // High
	IncidentSeverityP3 IncidentSeverity = "P3" // Medium
	IncidentSeverityP4 IncidentSeverity = "P4" // Low
	IncidentSeverityP5 IncidentSeverity = "P5" // Informational
)

// IncidentSource represents how an incident was created
type IncidentSource string

const (
	IncidentSourceAlert       IncidentSource = "alert"
	IncidentSourceManual      IncidentSource = "manual"
	IncidentSourceAI          IncidentSource = "ai"
	IncidentSourceIntegration IncidentSource = "integration"
)

// Incident represents a tracking record for an operational incident
type Incident struct {
	ID               uuid.UUID        `gorm:"type:uuid;primaryKey" json:"id"`
	OrganizationID   uuid.UUID        `gorm:"type:uuid;index;not null" json:"organization_id"`
	Title            string           `gorm:"type:text;not null" json:"title"`
	Description      string           `gorm:"type:text" json:"description"`
	Severity         IncidentSeverity `gorm:"size:10;index;not null;default:'P3'" json:"severity"`
	Status           IncidentStatus   `gorm:"size:30;index;not null;default:'OPEN'" json:"status"`
	Source           IncidentSource   `gorm:"size:50;not null;default:'manual'" json:"source"`
	MachineID        uuid.UUID        `gorm:"type:uuid;index" json:"machine_id"`
	AssignedTo       uuid.UUID        `gorm:"type:uuid;index" json:"assigned_to,omitempty"`
	CreatedBy        uuid.UUID        `gorm:"type:uuid;index" json:"created_by"`
	RootCause        string           `gorm:"type:text" json:"root_cause"`
	AIConfidence     float64          `json:"ai_confidence"`
	AIRecommendation string           `gorm:"type:text" json:"ai_recommendation"`
	AIImpact         string           `gorm:"type:text" json:"ai_impact"`
	AlertCount       int              `gorm:"default:1" json:"alert_count"`
	CommentCount     int              `gorm:"default:0" json:"comment_count"`
	MTTDSeconds      int64            `json:"mttd_seconds"`
	MTTRSeconds      int64            `json:"mttr_seconds"`
	MTTASeconds      int64            `json:"mtta_seconds"` // Mean Time To Acknowledge
	StartedAt        time.Time        `gorm:"index" json:"started_at"`
	AcknowledgedAt   *time.Time       `json:"acknowledged_at,omitempty"`
	ResolvedAt       *time.Time       `json:"resolved_at,omitempty"`
	ClosedAt         *time.Time       `json:"closed_at,omitempty"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`

	// Relations
	Timeline []IncidentTimeline `gorm:"foreignKey:IncidentID;constraint:OnDelete:CASCADE" json:"timeline,omitempty"`
	Comments []IncidentComment  `gorm:"foreignKey:IncidentID;constraint:OnDelete:CASCADE" json:"comments,omitempty"`
}

// IncidentTimeline represents an event in an incident's history
type IncidentTimeline struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	IncidentID uuid.UUID `gorm:"type:uuid;index;not null" json:"incident_id"`
	EventType  string    `gorm:"size:50;not null" json:"event_type"` // CREATED, ASSIGNED, ACKNOWLEDGED, COMMENTED, STATUS_CHANGED, RESOLVED, CLOSED
	Message    string    `gorm:"type:text;not null" json:"message"`
	CreatedBy  uuid.UUID `gorm:"type:uuid" json:"created_by"`
	CreatedAt  time.Time `gorm:"index" json:"created_at"`
}

// IncidentComment represents a comment on an incident
type IncidentComment struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	IncidentID uuid.UUID `gorm:"type:uuid;index;not null" json:"incident_id"`
	UserID     uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	Comment    string    `gorm:"type:text;not null" json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// IncidentAttachment represents a file attached to an incident
type IncidentAttachment struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	IncidentID uuid.UUID `gorm:"type:uuid;index;not null" json:"incident_id"`
	UserID     uuid.UUID `gorm:"type:uuid" json:"user_id"`
	FileName   string    `gorm:"not null" json:"file_name"`
	FilePath   string    `gorm:"not null" json:"file_path"`
	FileSize   int64     `json:"file_size"`
	FileType   string    `gorm:"not null" json:"file_type"` // screenshot, log, trace, report
	CreatedAt  time.Time `json:"created_at"`
}
