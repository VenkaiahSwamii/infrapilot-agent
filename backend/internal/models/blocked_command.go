package models

import (
	"time"

	"github.com/google/uuid"
)

type BlockedCommand struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey"`
	Pattern      string     `gorm:"type:text;uniqueIndex:idx_blocked_commands_org_pattern"`
	Reason       string     `gorm:"type:text"`
	Severity     string     `gorm:"size:30;not null;default:'BLOCKED'"`
	Organization string     `gorm:"size:255;not null;default:'Default Organization';uniqueIndex:idx_blocked_commands_org_pattern"`
	CreatedBy    *uuid.UUID `gorm:"type:uuid"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
