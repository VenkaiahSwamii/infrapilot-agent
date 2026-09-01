package models

import (
	"time"

	"github.com/google/uuid"
)

type AIIncident struct {
	ID                   uuid.UUID          `gorm:"type:uuid;primaryKey" json:"id"`
	Title                string             `gorm:"size:255;not null" json:"title"`
	Description          string             `gorm:"type:text" json:"description"`
	RootCause            string             `gorm:"size:255" json:"root_cause"`
	Severity             string             `gorm:"size:50;not null" json:"severity"`     // Critical, Warning, Info
	Status               string             `gorm:"size:50;default:'OPEN'" json:"status"` // OPEN, RESOLVED
	AffectedServersCount int                `json:"affected_servers_count"`
	RelatedAlertsCount   int                `json:"related_alerts_count"`
	CreatedAt            time.Time          `json:"created_at"`
	ResolvedAt           *time.Time         `json:"resolved_at,omitempty"`
	Timelines            []IncidentTimeline `gorm:"foreignKey:IncidentID;constraint:OnDelete:CASCADE" json:"timelines,omitempty"`
}

func (AIIncident) TableName() string {
	return "ai_incidents"
}

type AIRecommendation struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ServerID  uuid.UUID `gorm:"type:uuid;index;not null" json:"server_id"`
	Hostname  string    `gorm:"size:255" json:"hostname"`
	Title     string    `gorm:"size:255;not null" json:"title"`
	Reason    string    `gorm:"type:text;not null" json:"reason"`
	Priority  string    `gorm:"size:50;default:'Medium'" json:"priority"` // High, Medium, Low
	CreatedAt time.Time `json:"created_at"`
}

func (AIRecommendation) TableName() string {
	return "ai_recommendations"
}

type AIPrediction struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ServerID      uuid.UUID `gorm:"type:uuid;index;not null" json:"server_id"`
	Hostname      string    `gorm:"size:255" json:"hostname"`
	MetricName    string    `gorm:"size:100;index" json:"metric_name"` // CPU, Memory, Disk storage
	CurrentValue  float64   `json:"current_value"`
	ExpectedValue float64   `json:"expected_value"`
	TimeframeDays int       `json:"timeframe_days"`
	PredictedAt   time.Time `json:"predicted_at"`
	CreatedAt     time.Time `json:"created_at"`
}

func (AIPrediction) TableName() string {
	return "ai_predictions"
}

type AIHealthScore struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ServerID  uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"server_id"`
	Hostname  string    `gorm:"size:255" json:"hostname"`
	Score     int       `json:"score"`                    // 0-100
	Risk      string    `gorm:"size:50" json:"risk"`      // Low, Medium, High
	Factors   string    `gorm:"type:text" json:"factors"` // comma-separated reasons
	UpdatedAt time.Time `json:"updated_at"`
}

func (AIHealthScore) TableName() string {
	return "ai_health_scores"
}

type AIChatMessage struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	SessionID  string    `gorm:"size:255;index" json:"session_id"`
	UserQuery  string    `gorm:"type:text" json:"user_query"`
	AIResponse string    `gorm:"type:text" json:"ai_response"`
	CreatedAt  time.Time `json:"created_at"`
}

func (AIChatMessage) TableName() string {
	return "chat_history"
}
