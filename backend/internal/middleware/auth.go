package middleware

import (
	"net/http"
	"strings"

	"infrapilot/backend/internal/auth"
	"infrapilot/backend/internal/cache"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header must be Bearer {token}",
			})
			c.Abort()
			return
		}

		tokenStr := parts[1]

		if cache.IsTokenRevoked(tokenStr) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token has been revoked",
			})
			c.Abort()
			return
		}

		claims, err := auth.ValidateToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired JWT token",
			})
			c.Abort()
			return
		}

		// Store authenticated user information
		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("organizationId", claims.OrganizationID)

		c.Next()
	}
}
