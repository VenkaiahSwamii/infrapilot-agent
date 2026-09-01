package handlers

import (
	"net/http"
	"strconv"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetMachineHistory returns the historical metric series for a machine
func GetMachineHistory(c *gin.Context) {
	machineID := c.Param("id")
	hours := c.DefaultQuery("hours", "24")
	h, _ := strconv.Atoi(hours)

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

	var metrics []models.Metric
	database.DB.
		Where("machine_id = ?", machineUUID).
		Where("created_at > ?", time.Now().Add(-time.Duration(h)*time.Hour)).
		Order("created_at ASC").
		Find(&metrics)

	c.JSON(http.StatusOK, metrics)
}
