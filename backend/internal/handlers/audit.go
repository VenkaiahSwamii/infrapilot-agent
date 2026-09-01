package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"
)

// GetAuditLogs returns audit log entries with optional pagination.
// Query parameters:
//   - limit: maximum number of records to return (default 100)
//   - offset: number of records to skip (default 0)
func GetAuditLogs(c *gin.Context) {
	// Parse pagination parameters
	limitStr := c.DefaultQuery("limit", "100")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 100
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	var logs []models.AuditLog
	result := database.DB.Order("created_at desc").Limit(limit).Offset(offset).Find(&logs)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}
