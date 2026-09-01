package services

import (
	"fmt"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/events"
	"infrapilot/backend/internal/models"
	"infrapilot/backend/internal/websocket"

	"github.com/google/uuid"
)

type AlertEngine struct {
	eventBus *events.EventBus
}

func NewAlertEngine(bus ...*events.EventBus) *AlertEngine {
	var b *events.EventBus
	if len(bus) > 0 {
		b = bus[0]
	}
	return &AlertEngine{
		eventBus: b,
	}
}

func (a *AlertEngine) Evaluate(machineID string, cpu, memory, disk, latency, packetLoss float64) {
	mUUID, _ := uuid.Parse(machineID)

	check := func(metric string, value, threshold float64) {
		if value < threshold {
			return
		}

		title := fmt.Sprintf("%s Threshold Exceeded", metric)

		// Prevent duplicate open alert spam
		if database.DB != nil {
			var count int64
			if err := database.DB.Model(&models.LinuxAlert{}).Where("machine_id = ? AND (title = ? OR category = ?) AND LOWER(status) IN ('open', 'active')", mUUID, title, metric).Count(&count).Error; err == nil && count > 0 {
				database.DB.Model(&models.LinuxAlert{}).Where("machine_id = ? AND (title = ? OR category = ?) AND LOWER(status) IN ('open', 'active')", mUUID, title, metric).Updates(map[string]interface{}{
					"metric_value": value,
					"updated_at":   time.Now(),
				})
				return
			}
		}

		alert := models.LinuxAlert{
			ID:          uuid.New(),
			MachineID:   mUUID,
			Category:    metric,
			Title:       title,
			Severity:    "Critical",
			Priority:    models.MapSeverityToPriority("Critical"),
			Message:     fmt.Sprintf("%s is %.2f%% (Threshold %.2f%%)", metric, value, threshold),
			Status:      "OPEN",
			MetricValue: value,
			Threshold:   threshold,
			Source:      "AlertEngine",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if database.DB != nil {
			database.DB.Create(&alert)
		}

		// Sprint 7.7: Correlate alert into incident
		correlation := NewCorrelationEngine()
		inc, _ := correlation.Correlate(alert)

		// Sprint 7.8: Trigger Intelligent Auto Remediation Engine
		remediation := NewRemediationService()
		var incID uuid.UUID
		if inc != nil {
			incID = inc.ID
		}
		_, _ = remediation.EvaluateAndRemediate(alert, incID)

		// Sprint 7.9: Trigger Enterprise Workflow & Runbook Automation Engine
		workflowSvc := NewWorkflowService()
		var targetWfID uuid.UUID
		if database.DB != nil {
			var wf models.Workflow
			if database.DB.Where("enabled = ? AND (trigger_type = 'Alert' OR trigger_type = 'Incident')", true).First(&wf).Error == nil {
				targetWfID = wf.ID
			}
		}
		if targetWfID != uuid.Nil {
			_, _ = workflowSvc.StartWorkflow(targetWfID, incID, alert.MachineID)
		}

		// Sprint 7.6: Trigger enterprise notification engine
		notification := NewNotificationService()
		_ = notification.Send(alert)

		// Publish AlertCreatedEvent to EventBus
		if a.eventBus != nil {
			a.eventBus.Publish(events.AlertCreatedEvent{
				ID:        alert.ID.String(),
				MachineID: machineID,
				Severity:  alert.Severity,
				Title:     alert.Title,
				Message:   alert.Message,
				Status:    alert.Status,
				CreatedAt: alert.CreatedAt,
			})
		}

		// Direct Broadcast alert via WebSocket Hub if legacy
		if websocket.WS != nil {
			websocket.WS.Broadcast(map[string]interface{}{
				"type":       "alert",
				"alert":      alert,
				"id":         alert.ID.String(),
				"machine_id": machineID,
				"title":      alert.Title,
				"severity":   alert.Severity,
				"status":     alert.Status,
			})
		}
	}

	check("CPU", cpu, 90)
	check("Memory", memory, 90)
	check("Disk", disk, 90)
	check("Latency", latency, 300)
	check("Packet Loss", packetLoss, 5)
}
