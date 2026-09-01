package models

import (
	"time"

	"github.com/google/uuid"
)

type InstalledSoftware struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	MachineID       uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_software_machine_name;not null" json:"machine_id"`
	Name            string    `gorm:"uniqueIndex:idx_software_machine_name" json:"name"`
	Version         string    `json:"version"`
	Category        string    `json:"category"`
	Publisher       string    `json:"publisher"`
	UpdateAvailable bool      `json:"update_available"`
	LatestVersion   string    `json:"latest_version"`
	IsSecurityPatch bool      `json:"is_security_patch"`
	InstalledAt     time.Time `json:"installed_at"`
	LastSeenAt      time.Time `json:"last_seen_at"`
}
