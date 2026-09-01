package handlers

import (
	"log"
	"net/http"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HeartbeatRequest struct {
	MachineID string `json:"machine_id"`
}

func Heartbeat(c *gin.Context) {

	var req HeartbeatRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	machineID, err := uuid.Parse(req.MachineID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid machine_id format"})
		return
	}

	var machine models.Machine

	if err := database.DB.First(&machine, "id = ?", machineID).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Machine not found",
		})

		return
	}

	machine.LastSeen = time.Now().UTC()
	machine.Status = "ONLINE"
	machine.Online = true
	machine.RetryCount = 0

	if err := database.DB.Save(&machine).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to persist heartbeat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Heartbeat received",
	})

	log.Printf("HEARTBEAT RECEIVED machine_id=%s hostname=%s ip=%s last_seen=%s status=ONLINE",
		machine.ID.String(), machine.Hostname, machine.IPAddress, machine.LastSeen.Format(time.RFC3339))
}
