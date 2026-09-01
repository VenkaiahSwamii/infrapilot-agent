package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var RolePermissions = map[string][]string{
	"SuperAdmin": {
		"View Dashboard", "Manage Servers", "Deploy Agents", "View Logs", "Manage Alerts",
		"Manage Users", "Configure AI", "Configure Kubernetes",
	},
	"OrganizationAdmin": {
		"View Dashboard", "Manage Servers", "Deploy Agents", "View Logs", "Manage Alerts",
		"Manage Users", "Configure AI",
	},
	"DevOps": {
		"View Dashboard", "Manage Servers", "Deploy Agents", "View Logs", "Manage Alerts",
		"Configure AI", "Configure Kubernetes",
	},
	"Operator": {
		"View Dashboard", "Manage Servers", "View Logs", "Manage Alerts",
	},
	"ReadOnly": {
		"View Dashboard",
	},
}

// RequirePermission checks if the authenticated user has the specified permission.
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Role not found in context"})
			c.Abort()
			return
		}

		role := roleVal.(string)
		perms := RolePermissions[role]

		// Fallbacks for default legacy roles
		if len(perms) == 0 {
			if role == "Admin" {
				perms = RolePermissions["OrganizationAdmin"]
			} else if role == "Viewer" {
				perms = RolePermissions["ReadOnly"]
			}
		}

		for _, p := range perms {
			if p == permission {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied: requires " + permission})
		c.Abort()
	}
}
