package main

import (
	"log"
	"time"

	"infrapilot/backend/internal/config"
	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func main() {
	log.Println("Starting InfraPilot Scheduler Process...")

	config.LoadEnv()
	database.Connect()

	log.Println("Scheduler initialized. Starting background tasks...")

	// Run health check (offline detector) loop every 15 seconds
	go runOfflineDetector(database.DB)

	// Run database maintenance every 60 seconds.
	for {
		log.Println("[Scheduler] Running periodic database maintenance and rollups...")

		var count int64
		database.DB.Model(&models.HistoricalMetric{}).Count(&count)
		log.Printf("[Scheduler] Maintenance: Found %d historical metric records.", count)

		time.Sleep(60 * time.Second)
	}
}

func runOfflineDetector(db *gorm.DB) {
	for {
		var machines []models.Machine
		// Find machines that haven't sent a heartbeat for 30 seconds and are still marked ONLINE
		threshold := time.Now().Add(-30 * time.Second)
		err := db.Where("last_seen < ? AND (status = 'ONLINE' OR online = true)", threshold).Find(&machines).Error
		if err != nil {
			log.Printf("[Scheduler] Error querying active machines: %v", err)
			time.Sleep(15 * time.Second)
			continue
		}

		for _, machine := range machines {
			machine.Status = "OFFLINE"
			machine.Online = false
			if err := db.Save(&machine).Error; err != nil {
				log.Printf("[Scheduler] Error updating machine status: %v", err)
				continue
			}

			log.Printf("[Scheduler] Machine %s has gone OFFLINE (no heartbeat for 30s)", machine.Hostname)

			// Create critical alert for the offline machine
			alert := models.LinuxAlert{
				ID:        uuid.New(),
				MachineID: machine.ID,
				Type:      "machine_offline",
				Severity:  "critical",
				Message:   "Machine has gone OFFLINE (no heartbeat for 30s)",
				Status:    "OPEN",
				CreatedAt: time.Now(),
			}
			if err := db.Create(&alert).Error; err != nil {
				log.Printf("[Scheduler] Error creating offline alert: %v", err)
			}
		}

		time.Sleep(15 * time.Second)
	}
}
