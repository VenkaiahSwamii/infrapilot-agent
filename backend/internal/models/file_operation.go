package models

import (
	"time"

	"github.com/google/uuid"
)

type FileOperation struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	MachineID uuid.UUID `gorm:"type:uuid;index"`
	Operation string    `gorm:"size:30"`
	Path      string    `gorm:"size:500"`
	Result    string    `gorm:"type:text"`
	Status    string    `gorm:"size:20"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
