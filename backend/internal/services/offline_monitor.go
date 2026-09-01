package services

import (
	"encoding/json"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"
	"infrapilot/backend/internal/websocket"
)

// OfflineMonitor periodically checks for offline machines
type OfflineMonitor struct {
	hub          *websocket.Hub
	threshold    time.Duration
	pollInterval time.Duration
}

// NewOfflineMonitor creates a new offline monitor
func NewOfflineMonitor(hub *websocket.Hub, pollInterval, threshold time.Duration) *OfflineMonitor {
	return &OfflineMonitor{
		hub:          hub,
		threshold:    threshold,
		pollInterval: pollInterval,
	}
}

// Start begins the offline detection loop
func (m *OfflineMonitor) Start() {
	// Status checking consolidated into status_monitor.go to avoid database write race conditions
}

func (m *OfflineMonitor) checkOfflineMachines() {
	var machines []models.Machine
	if err := database.DB.Find(&machines).Error; err != nil {
		return
	}

	now := time.Now()
	for _, machine := range machines {
		if machine.Status == "OFFLINE" {
			continue
		}

		if now.Sub(machine.LastSeen) <= m.threshold {
			continue
		}

		machine.Status = "OFFLINE"
		if err := database.DB.Save(&machine).Error; err != nil {
			continue
		}

		payload := map[string]interface{}{
			"type":       "machine_offline",
			"machine_id": machine.ID.String(),
			"hostname":   machine.Hostname,
			"ip_address": machine.IPAddress,
			"os":         machine.OS,
			"platform":   machine.Platform,
			"status":     "OFFLINE",
			"last_seen":  machine.LastSeen,
		}
		jsonData := jsonMarshal(payload)
		m.hub.Broadcast(jsonData)
	}
}

func jsonMarshal(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}
