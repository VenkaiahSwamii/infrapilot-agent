package middleware

import (
	"net/http"
	"strings"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// OrganizationMiddleware extracts, validates, and injects tenant organization context
func OrganizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		userIDStr := c.GetString("userId")

		// 1. Determine requested organization_id from Header, Query, or Path parameter
		orgID := c.GetHeader("X-Organization-ID")
		if orgID == "" {
			orgID = c.Param("orgId")
		}
		if orgID == "" {
			orgID = c.GetString("organizationId")
		}
		if orgID == "" {
			orgID = "default"
		}

		// Clean orgID
		orgID = strings.TrimSpace(orgID)

		// SuperAdmin bypasses tenant restriction check for cross-tenant management
		if role == string(models.RoleSuperAdmin) || role == "super_admin" {
			c.Set("organizationId", orgID)
			c.Next()
			return
		}

		// Default organization is universally accessible if no specific org enforced
		if orgID == "default" || orgID == "00000000-0000-0000-0000-000000000000" {
			c.Set("organizationId", "default")
			c.Next()
			return
		}

		// 2. Validate membership if valid UUID
		parsedOrgID, errOrg := uuid.Parse(orgID)
		parsedUserID, errUser := uuid.Parse(userIDStr)

		if errOrg == nil && errUser == nil && database.DB != nil {
			var membership models.OrganizationUser
			err := database.DB.Where("organization_id = ? AND user_id = ?", parsedOrgID, parsedUserID).First(&membership).Error
			if err != nil {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "Access denied: User is not a member of the requested organization",
				})
				c.Abort()
				return
			}
			c.Set("tenantRole", membership.Role)
		}

		c.Set("organizationId", orgID)
		c.Next()
	}
}
