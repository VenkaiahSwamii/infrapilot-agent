package services

import (
	"fmt"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"
	"infrapilot/backend/internal/repository"
	"infrapilot/backend/internal/websocket"

	"github.com/google/uuid"
)

type AlertService struct {
	alertRepo *repository.AlertRepository
}

func NewAlertService(alertRepo *repository.AlertRepository) *AlertService {
	return &AlertService{
		alertRepo: alertRepo,
	}
}

func (s *AlertService) CreateAlert(alert *models.LinuxAlert) error {
	if alert.ID == uuid.Nil {
		alert.ID = uuid.New()
	}
	if alert.CreatedAt.IsZero() {
		alert.CreatedAt = time.Now()
	}
	if alert.UpdatedAt.IsZero() {
		alert.UpdatedAt = time.Now()
	}
	if alert.Status == "" {
		alert.Status = "OPEN"
	}
	if s.alertRepo != nil {
		return s.alertRepo.Create(alert)
	}
	return database.DB.Create(alert).Error
}

func (s *AlertService) ResolveAlert(alertID uuid.UUID, resolutionNote string) error {
	if s.alertRepo != nil {
		return s.alertRepo.ResolveAlert(alertID, resolutionNote, "System (Auto-Resolved)")
	}
	now := time.Now()
	return database.DB.Model(&models.LinuxAlert{}).
		Where("id = ?", alertID).
		Updates(map[string]interface{}{
			"status":          "RESOLVED",
			"resolution_note": resolutionNote,
			"resolved_by":     "System (Auto-Resolved)",
			"resolved_at":     now,
			"updated_at":      now,
		}).Error
}

func (s *AlertService) FindOpenAlert(machineID uuid.UUID, ruleID uuid.UUID) (*models.LinuxAlert, error) {
	if s.alertRepo != nil {
		return s.alertRepo.FindOpenByMachineAndRule(machineID, ruleID)
	}
	var alert models.LinuxAlert
	err := database.DB.Where("machine_id = ? AND rule_id = ? AND status = ?", machineID, ruleID, "OPEN").First(&alert).Error
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

func (s *AlertService) ListOpenAlerts(machineID uuid.UUID) ([]models.LinuxAlert, error) {
	if s.alertRepo != nil {
		return s.alertRepo.ListOpenByMachine(machineID)
	}
	var alerts []models.LinuxAlert
	err := database.DB.Where("machine_id = ? AND status = ?", machineID, "OPEN").
		Order("created_at DESC").
		Find(&alerts).Error
	return alerts, err
}

// CheckAlerts evaluates metrics against critical thresholds with deduplication & auto-resolution.
func CheckAlerts(machineID uuid.UUID, cpu, memory, disk, latency, packetLoss float64) {
	if database.DB == nil {
		return
	}

	evaluateMetric := func(title, category, severity string, val, threshold, recoveryThreshold float64, isHigherBreach bool) {
		isBreached := false
		if isHigherBreach {
			isBreached = val >= threshold
		} else {
			isBreached = val <= threshold
		}

		// Look for existing active/open alert for this machine and title
		var existing models.LinuxAlert
		err := database.DB.Where("machine_id = ? AND (title = ? OR type = ? OR category = ?) AND LOWER(status) IN ('open', 'active')", machineID, title, title, category).
			Order("created_at desc").
			First(&existing).Error

		if isBreached {
			if err == nil {
				// Alert already exists and is active! Update value and timestamp without duplicating
				database.DB.Model(&models.LinuxAlert{}).Where("id = ?", existing.ID).Updates(map[string]interface{}{
					"metric_value": val,
					"message":      fmt.Sprintf("%s threshold breached: %.1f%% (threshold %.1f%%)", title, val, threshold),
					"updated_at":   time.Now(),
				})
				return
			}

			// New alert breach! Create alert row
			alert := models.LinuxAlert{
				ID:                 uuid.New(),
				MachineID:          machineID,
				Title:              title,
				Type:               title,
				Category:           category,
				Message:            fmt.Sprintf("%s threshold breached: %.1f%% (threshold %.1f%%)", title, val, threshold),
				Severity:           severity,
				Priority:           models.MapSeverityToPriority(severity),
				Status:             "ACTIVE",
				MetricValue:        val,
				Threshold:          threshold,
				RecoverySuggestion: fmt.Sprintf("Inspect top consuming processes or scale machine resources."),
				Source:             "AlertEngine",
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
			}

			database.DB.Create(&alert)
			SendAlert(alert)

			// Broadcast alert immediately through WebSockets
			if websocket.WS != nil {
				websocket.WS.Broadcast(map[string]interface{}{
					"type":        "alert",
					"id":          alert.ID.String(),
					"machine_id":  alert.MachineID.String(),
					"title":       alert.Title,
					"description": alert.Message,
					"severity":    alert.Severity,
					"status":      alert.Status,
					"value":       alert.MetricValue,
					"created_at":  alert.CreatedAt,
				})
			}
		} else {
			// If not breached and value is safely below recovery threshold, auto-resolve existing open alert
			isRecovered := false
			if isHigherBreach {
				isRecovered = val < recoveryThreshold
			} else {
				isRecovered = val > recoveryThreshold
			}

			if isRecovered && err == nil {
				now := time.Now()
				database.DB.Model(&models.LinuxAlert{}).Where("id = ?", existing.ID).Updates(map[string]interface{}{
					"status":          "RESOLVED",
					"resolution_note": fmt.Sprintf("Auto-resolved: Metric normalized to %.1f%% (recovery threshold %.1f%%)", val, recoveryThreshold),
					"resolved_by":     "System (Auto-Resolved)",
					"resolved_at":     &now,
					"updated_at":      now,
				})

				if websocket.WS != nil {
					websocket.WS.Broadcast(map[string]interface{}{
						"type":        "alert.resolved",
						"id":          existing.ID.String(),
						"machine_id":  existing.MachineID.String(),
						"title":       existing.Title,
						"status":      "RESOLVED",
						"resolved_at": now,
					})
				}
			}
		}
	}

	evaluateMetric("High CPU Usage", "CPU", "Critical", cpu, 90, 80, true)
	evaluateMetric("High Memory Usage", "Memory", "Critical", memory, 90, 80, true)
	evaluateMetric("Disk Almost Full", "Storage", "Critical", disk, 95, 88, true)
	evaluateMetric("High Latency", "Network", "Warning", latency, 100, 70, true)
	evaluateMetric("High Packet Loss", "Network", "Critical", packetLoss, 5, 2, true)
}
