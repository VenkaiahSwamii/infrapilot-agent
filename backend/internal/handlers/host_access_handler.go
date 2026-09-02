package handlers

import (
	"net/http"
	"strings"

	"infrapilot/backend/internal/models"
	"infrapilot/backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HostAccessHandler struct {
	accessService *services.HostAccessService
}

func NewHostAccessHandler(accessService *services.HostAccessService) *HostAccessHandler {
	return &HostAccessHandler{accessService: accessService}
}

type UpdateMachineAccessRequest struct {
	UserID          string `json:"user_id" binding:"required"`
	PermissionLevel string `json:"permission_level" binding:"required"` // "Full Control", "Operator", "Viewer", "None"
}

// GetMachineAccess lists all user permissions for a target machine
func (h *HostAccessHandler) GetMachineAccess(c *gin.Context) {
	idStr := c.Param("id")
	machineID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid machine_id UUID format"})
		return
	}

	perms, err := h.accessService.GetHostPermissionsForMachine(machineID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch machine permissions: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"machine_id":  machineID.String(),
		"permissions": perms,
	})
}

// UpdateMachineAccess grants or updates user permission for a machine
func (h *HostAccessHandler) UpdateMachineAccess(c *gin.Context) {
	idStr := c.Param("id")
	machineID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid machine_id UUID format"})
		return
	}

	var req UpdateMachineAccessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	targetUserID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target user_id UUID format"})
		return
	}

	// Get actor info from auth context
	actorIDVal, _ := c.Get("userId")
	actorNameVal, _ := c.Get("username")
	actorRoleVal, _ := c.Get("role")

	actorIDStr, _ := actorIDVal.(string)
	actorID, _ := uuid.Parse(actorIDStr)
	actorName, _ := actorNameVal.(string)
	actorRole, _ := actorRoleVal.(string)

	if actorName == "" {
		actorName = "Admin User"
	}

	// Enforce Admin role for managing access
	r := strings.ToLower(strings.TrimSpace(actorRole))
	if r != "superadmin" && r != "admin" && r != "orgadmin" && actorRole != "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only Admins can delegate machine access"})
		return
	}

	perm, err := h.accessService.GrantOrUpdatePermission(targetUserID, machineID, actorID, req.PermissionLevel, actorName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update machine permission: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "Machine access updated successfully",
		"permission":       perm,
		"machine_id":       machineID.String(),
		"user_id":          targetUserID.String(),
		"permission_level": req.PermissionLevel,
	})
}

// RevokeMachineAccess deletes a user's machine permission
func (h *HostAccessHandler) RevokeMachineAccess(c *gin.Context) {
	idStr := c.Param("id")
	machineID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid machine_id format"})
		return
	}

	userIDStr := c.Param("userId")
	targetUserID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id format"})
		return
	}

	actorRoleVal, _ := c.Get("role")
	actorRole, _ := actorRoleVal.(string)
	r := strings.ToLower(strings.TrimSpace(actorRole))
	if r != "superadmin" && r != "admin" && r != "orgadmin" && actorRole != "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only Admins can revoke machine access"})
		return
	}

	if err := h.accessService.RevokePermission(targetUserID, machineID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke machine access"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "Machine access revoked successfully",
		"machine_id":       machineID.String(),
		"user_id":          targetUserID.String(),
		"permission_level": models.PermissionNone,
	})
}

// GetMyMachineAccess lists allowed machine IDs and levels for the logged in user
func (h *HostAccessHandler) GetMyMachineAccess(c *gin.Context) {
	actorIDVal, _ := c.Get("userId")
	actorRoleVal, _ := c.Get("role")

	actorIDStr, _ := actorIDVal.(string)
	actorID, _ := uuid.Parse(actorIDStr)
	actorRole, _ := actorRoleVal.(string)

	allowed, isSuperAdmin := h.accessService.GetUserAllowedMachineIDs(actorID, actorRole)

	c.JSON(http.StatusOK, gin.H{
		"is_admin":          isSuperAdmin,
		"allowed_machines": allowed,
	})
}
