package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	PermissionFullControl = "Full Control"
	PermissionOperator    = "Operator"
	PermissionViewer      = "Viewer"
	PermissionNone        = "None"
)

// UserHostPermission model defines user access levels per machine
type UserHostPermission struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID          uuid.UUID `gorm:"type:uuid;index:idx_user_host,unique;not null" json:"user_id"`
	MachineID       uuid.UUID `gorm:"type:uuid;index:idx_user_host,unique;not null" json:"machine_id"`
	PermissionLevel string    `gorm:"size:50;not null;default:'Viewer'" json:"permission_level"` // "Full Control", "Operator", "Viewer"
	GrantedBy       uuid.UUID `gorm:"type:uuid" json:"granted_by"`
	GrantedByName   string    `gorm:"size:100" json:"granted_by_name"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// Associations
	User    *User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Machine *Server `gorm:"foreignKey:MachineID" json:"machine,omitempty"`
}

func (UserHostPermission) TableName() string {
	return "user_host_permissions"
}

func (p *UserHostPermission) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	if p.PermissionLevel == "" {
		p.PermissionLevel = PermissionViewer
	}
	return nil
}
