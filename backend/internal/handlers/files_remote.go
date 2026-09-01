package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func resolveRemoteFileMachineID(c *gin.Context) (uuid.UUID, error) {
	idStr := strings.TrimSpace(c.Param("id"))
	if idStr == "" {
		return uuid.Nil, fmt.Errorf("missing machine ID")
	}
	if parsed, err := uuid.Parse(idStr); err == nil {
		return parsed, nil
	}
	if database.DB != nil {
		var machine models.Machine
		if err := database.DB.Where("LOWER(hostname) = LOWER(?) OR LOWER(name) = LOWER(?) OR id::text LIKE ? OR ip_address = ?", idStr, idStr, idStr+"%", idStr).First(&machine).Error; err == nil {
			return machine.ID, nil
		}
	}
	return uuid.Nil, fmt.Errorf("machine not found")
}

// GetMachineFiles handles GET /api/v1/machines/:id/files
func GetMachineFiles(c *gin.Context) {
	machineID, err := resolveRemoteFileMachineID(c)
	if err != nil {
		machineID = uuid.Nil
	}

	path := c.Query("path")
	if path == "" {
		path = "/"
	}

	// Create audit record in file_operations table
	if machineID != uuid.Nil && database.DB != nil {
		fo := models.FileOperation{
			ID:        uuid.New(),
			MachineID: machineID,
			Operation: "LIST",
			Path:      path,
			Status:    "COMPLETED",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		_ = database.DB.Create(&fo)
	}

	// Delegates to local ListFiles
	ListFiles(c)
}

// CreateFileOperation handles POST /api/v1/machines/:id/files/operation
func CreateFileOperation(c *gin.Context) {
	machineID, _ := resolveRemoteFileMachineID(c)

	var req struct {
		Operation string `json:"operation" binding:"required"`
		Path      string `json:"path" binding:"required"`
		Result    string `json:"result"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fo := models.FileOperation{
		ID:        uuid.New(),
		MachineID: machineID,
		Operation: req.Operation,
		Path:      req.Path,
		Result:    req.Result,
		Status:    "COMPLETED",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if database.DB != nil {
		if err := database.DB.Create(&fo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to record file operation: %v", err)})
			return
		}
	}

	c.JSON(http.StatusOK, fo)
}
