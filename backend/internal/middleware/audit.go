package middleware

import (
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuditMiddleware automatically logs all state-modifying HTTP requests.
func AuditMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method == "GET" || method == "OPTIONS" || method == "HEAD" {
			c.Next()
			return
		}

		c.Next()

		usernameVal, exists := c.Get("username")
		username := "anonymous"
		if exists {
			username = usernameVal.(string)
		}

		status := c.Writer.Status()
		result := "Success"
		if status >= 400 {
			result = "Failure: HTTP Status Code"
		}

		if database.DB != nil {
			entry := models.AuditLog{
				ID:        uuid.New(),
				Username:  username,
				IP:        c.ClientIP(),
				Resource:  c.Request.URL.Path,
				MachineID: uuid.Nil,
				Action:    method + " request",
				Result:    result,
				CreatedAt: time.Now(),
			}
			database.DB.Create(&entry)
		}
	}
}
