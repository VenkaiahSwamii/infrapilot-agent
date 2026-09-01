package middleware

import (
	"net/http"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"

	"github.com/gin-gonic/gin"
)

// APIKeyAuthMiddleware verifies the client's Bearer token against the machines table
func APIKeyAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: missing or invalid Bearer token"})
			return
		}
		apiKey := authHeader[7:]

		var machine models.Machine
		if err := database.DB.Where("api_key = ?", apiKey).First(&machine).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: invalid API Key"})
			return
		}

		machine.LastSeen = time.Now()
		database.DB.Save(machine)

		// Store machine metadata context in Gin context values
		c.Set("machine", &machine)
		c.Set("api_key", apiKey)
		c.Next()
	}
}
