package services

import (
	"errors"
	"strings"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

type HostAccessService struct{}

func NewHostAccessService() *HostAccessService {
	return &HostAccessService{}
}

type UserPermissionDTO struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	Username        string    `json:"username"`
	Email           string    `json:"email"`
	UserRole        string    `json:"user_role"`
	MachineID       uuid.UUID `json:"machine_id"`
	PermissionLevel string    `json:"permission_level"`
	GrantedByName   string    `json:"granted_by_name"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// GetHostPermissionsForMachine lists all permissions for a specific machine
func (s *HostAccessService) GetHostPermissionsForMachine(machineID uuid.UUID) ([]UserPermissionDTO, error) {
	if database.DB == nil {
		return nil, errors.New("database unavailable")
	}

	var users []models.User
	if err := database.DB.Where("is_active = ?", true).Order("username ASC").Find(&users).Error; err != nil {
		return nil, err
	}

	var existing []models.UserHostPermission
	_ = database.DB.Where("machine_id = ?", machineID).Find(&existing).Error

	permMap := make(map[uuid.UUID]models.UserHostPermission)
	for _, p := range existing {
		permMap[p.UserID] = p
	}

	res := make([]UserPermissionDTO, 0, len(users))
	for _, u := range users {
		level := models.PermissionNone
		grantedBy := ""
		updatedAt := u.CreatedAt
		permID := uuid.Nil

		// Admins always have Full Control by role default
		if isSystemAdminRole(u.Role) {
			level = models.PermissionFullControl
			grantedBy = "System Role (" + u.Role + ")"
		} else if p, ok := permMap[u.ID]; ok {
			level = p.PermissionLevel
			grantedBy = p.GrantedByName
			updatedAt = p.UpdatedAt
			permID = p.ID
		}

		res = append(res, UserPermissionDTO{
			ID:              permID,
			UserID:          u.ID,
			Username:        u.Username,
			Email:           u.Email,
			UserRole:        u.Role,
			MachineID:       machineID,
			PermissionLevel: level,
			GrantedByName:   grantedBy,
			UpdatedAt:       updatedAt,
		})
	}

	return res, nil
}

// GrantOrUpdatePermission sets a user's machine permission
func (s *HostAccessService) GrantOrUpdatePermission(userID, machineID, grantedBy uuid.UUID, permissionLevel, grantedByName string) (*models.UserHostPermission, error) {
	if database.DB == nil {
		return nil, errors.New("database unavailable")
	}

	permissionLevel = strings.TrimSpace(permissionLevel)
	if permissionLevel == "" {
		permissionLevel = models.PermissionViewer
	}

	// If level is None, delete permission record if it exists
	if permissionLevel == models.PermissionNone {
		_ = database.DB.Where("user_id = ? AND machine_id = ?", userID, machineID).Delete(&models.UserHostPermission{}).Error
		return nil, nil
	}

	var existing models.UserHostPermission
	err := database.DB.Where("user_id = ? AND machine_id = ?", userID, machineID).First(&existing).Error
	if err == nil {
		// Update existing
		existing.PermissionLevel = permissionLevel
		existing.GrantedBy = grantedBy
		existing.GrantedByName = grantedByName
		existing.UpdatedAt = time.Now()
		if err := database.DB.Save(&existing).Error; err != nil {
			return nil, err
		}
		return &existing, nil
	}

	// Create new
	perm := &models.UserHostPermission{
		ID:              uuid.New(),
		UserID:          userID,
		MachineID:       machineID,
		PermissionLevel: permissionLevel,
		GrantedBy:       grantedBy,
		GrantedByName:   grantedByName,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := database.DB.Create(perm).Error; err != nil {
		return nil, err
	}

	return perm, nil
}

// RevokePermission revokes machine permission for a user
func (s *HostAccessService) RevokePermission(userID, machineID uuid.UUID) error {
	if database.DB == nil {
		return nil
	}
	return database.DB.Where("user_id = ? AND machine_id = ?", userID, machineID).Delete(&models.UserHostPermission{}).Error
}

// GetUserAllowedMachineIDs returns map of allowed machine UUIDs and permission levels
func (s *HostAccessService) GetUserAllowedMachineIDs(userID uuid.UUID, role string) (map[uuid.UUID]string, bool) {
	if isSystemAdminRole(role) {
		return nil, true // Admin has full access to all machines
	}

	allowed := make(map[uuid.UUID]string)
	if database.DB == nil || userID == uuid.Nil {
		return allowed, false
	}

	var perms []models.UserHostPermission
	if err := database.DB.Where("user_id = ?", userID).Find(&perms).Error; err == nil {
		for _, p := range perms {
			if p.PermissionLevel != models.PermissionNone {
				allowed[p.MachineID] = p.PermissionLevel
			}
		}
	}

	return allowed, false
}

// CheckHostAccess checks if a user has sufficient permission level for a machine
func (s *HostAccessService) CheckHostAccess(userID uuid.UUID, role string, machineID uuid.UUID, minRequired string) (bool, string) {
	if isSystemAdminRole(role) {
		return true, models.PermissionFullControl
	}

	if database.DB == nil || userID == uuid.Nil || machineID == uuid.Nil {
		return false, models.PermissionNone
	}

	var perm models.UserHostPermission
	err := database.DB.Where("user_id = ? AND machine_id = ?", userID, machineID).First(&perm).Error
	if err != nil {
		return false, models.PermissionNone
	}

	if !isPermissionSufficient(perm.PermissionLevel, minRequired) {
		return false, perm.PermissionLevel
	}

	return true, perm.PermissionLevel
}

func isSystemAdminRole(role string) bool {
	r := strings.ToLower(strings.TrimSpace(role))
	return r == "superadmin" || r == "admin" || r == "orgadmin"
}

func isPermissionSufficient(actualLevel, requiredLevel string) bool {
	if actualLevel == models.PermissionNone || actualLevel == "" {
		return false
	}

	// Permission Hierarchy: Full Control > Operator > Viewer
	switch requiredLevel {
	case models.PermissionViewer, "viewer", "READ":
		return actualLevel == models.PermissionFullControl || actualLevel == models.PermissionOperator || actualLevel == models.PermissionViewer
	case models.PermissionOperator, "operator", "EXECUTE":
		return actualLevel == models.PermissionFullControl || actualLevel == models.PermissionOperator
	case models.PermissionFullControl, "full control", "ADMIN":
		return actualLevel == models.PermissionFullControl
	default:
		return actualLevel == models.PermissionFullControl || actualLevel == models.PermissionOperator || actualLevel == models.PermissionViewer
	}
}
