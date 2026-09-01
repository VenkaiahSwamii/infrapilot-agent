package models

import (
	"time"

	"github.com/google/uuid"
)

type QueuedMetric struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Payload   string    `gorm:"type:text;not null" json:"payload"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}
