package utils

import (
	"github.com/google/uuid"
	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"
	"time"
)

// LogAudit creates an audit log entry.
func LogAudit(username string, machineID uuid.UUID, action string, result string) {
	// If the database connection is not initialized (e.g., during unit tests), skip logging to avoid panic.
	if database.DB == nil {
		return
	}
	entry := models.AuditLog{ID: uuid.New(), Username: username, MachineID: machineID, Action: action, Result: result, CreatedAt: time.Now()}
	database.DB.Create(&entry)
}
