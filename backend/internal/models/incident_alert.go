package models

import (
	"time"

	"github.com/google/uuid"
)

type IncidentAlert struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	IncidentID uuid.UUID `gorm:"type:uuid;index;not null" json:"incident_id"`
	AlertID    uuid.UUID `gorm:"type:uuid;index;not null" json:"alert_id"`
	CreatedAt  time.Time `json:"created_at"`
}
