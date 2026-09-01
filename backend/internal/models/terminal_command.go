package models

import (
	"time"

	"github.com/google/uuid"
)

type TerminalCommand struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	SessionID uuid.UUID `gorm:"type:uuid;index"`
	MachineID uuid.UUID `gorm:"type:uuid;index"`

	Command string `gorm:"type:text"`
	Stdout  string `gorm:"type:text"`
	Stderr  string `gorm:"type:text"`
	Output  string `gorm:"type:text"`

	ExitCode *int `gorm:"type:integer"`

	Status string `gorm:"size:30;index"`

	CreatedAt  time.Time
	UpdatedAt  time.Time
	StartedAt  *time.Time
	FinishedAt *time.Time
}
