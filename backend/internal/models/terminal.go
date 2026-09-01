package models

import (
	"time"

	"github.com/google/uuid"
)

type TerminalSession struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	MachineID uuid.UUID `gorm:"type:uuid;index;not null"`
	UserID    uuid.UUID `gorm:"type:uuid;index;not null"`

	SessionID string `gorm:"size:100;uniqueIndex"`
	Status    string `gorm:"size:20"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
