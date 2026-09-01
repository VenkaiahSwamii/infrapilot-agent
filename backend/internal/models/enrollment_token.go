package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EnrollmentToken struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	OrganizationID string     `gorm:"index" json:"organization_id"`
	Name           string     `json:"name"`
	TokenPrefix    string     `gorm:"index" json:"token_prefix"`
	TokenHash      string     `gorm:"uniqueIndex" json:"-"`
	Token          string     `gorm:"index" json:"token"` // For Sprint 2 compatibility
	Used           bool       `json:"used"`               // For Sprint 2 compatibility
	UsedCount      int        `json:"used_count"`
	MaxUses        int        `json:"max_uses"`
	ExpiresAt      *time.Time `json:"expires_at"`
	RevokedAt      *time.Time `json:"revoked_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (e *EnrollmentToken) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}
