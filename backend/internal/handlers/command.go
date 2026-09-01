package handlers

import (
	"net/http"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateCommand queues a new remote command execution
func CreateCommand(c *gin.Context) {
	var req struct {
		MachineID string `json:"machine_id" binding:"required"`
		Command   string `json:"command" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var machineUUID uuid.UUID
	var err error
	if machineUUID, err = uuid.Parse(req.MachineID); err != nil {
		var machine models.Machine
		if database.DB.Where("id::text LIKE ? OR hostname = ?", req.MachineID+"%", req.MachineID).First(&machine).Error == nil {
			machineUUID = machine.ID
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid machine ID"})
			return
		}
	}

	cmd := models.Command{
		ID:        uuid.New(),
		MachineID: machineUUID,
		Command:   req.Command,
		Status:    "Pending",
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&cmd).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create command"})
		return
	}

	c.JSON(http.StatusCreated, cmd)
}

// GetPendingCommands returns command queue items targeting a specific machine
func GetPendingCommands(c *gin.Context) {
	machineID := c.Query("machine_id")
	if machineID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "machine_id query parameter is required"})
		return
	}

	var machineUUID uuid.UUID
	var err error
	if machineUUID, err = uuid.Parse(machineID); err != nil {
		var machine models.Machine
		if database.DB.Where("id::text LIKE ? OR hostname = ?", machineID+"%", machineID).First(&machine).Error == nil {
			machineUUID = machine.ID
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid machine ID"})
			return
		}
	}

	var commands []models.Command
	if err := database.DB.Where("machine_id = ? AND status = 'Pending'", machineUUID).Order("created_at ASC").Find(&commands).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch commands"})
		return
	}

	// Update status to Running to prevent duplicated delivery
	for i := range commands {
		commands[i].Status = "Running"
		database.DB.Save(&commands[i])
	}

	c.JSON(http.StatusOK, commands)
}

// PostCommandResult saves execution output and transitions status to Completed or Failed
func PostCommandResult(c *gin.Context) {
	var req struct {
		CommandID string `json:"command_id" binding:"required"`
		Output    string `json:"output"`
		Status    string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmdUUID, err := uuid.Parse(req.CommandID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid command ID"})
		return
	}

	var cmd models.Command
	if err := database.DB.First(&cmd, "id = ?", cmdUUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "command not found"})
		return
	}

	cmd.Output = req.Output
	cmd.Status = req.Status
	if err := database.DB.Save(&cmd).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save result"})
		return
	}

	c.JSON(http.StatusOK, cmd)
}

// GetCommandLogs returns execution logs history
func GetCommandLogs(c *gin.Context) {
	var logs []models.Command
	database.DB.Order("created_at desc").Limit(50).Find(&logs)
	c.JSON(http.StatusOK, logs)
}
