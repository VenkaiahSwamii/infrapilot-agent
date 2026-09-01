package services

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"
	"infrapilot/backend/internal/websocket"

	"github.com/google/uuid"
)

var (
	memIncidentsMu sync.RWMutex
	memIncidents   = make(map[uuid.UUID]*models.Incident)
)

type IncidentService struct{}

func NewIncidentService() *IncidentService {
	return &IncidentService{}
}

func (s *IncidentService) FindOpenIncident(machineID uuid.UUID) (*models.Incident, error) {
	if database.DB != nil {
		var incident models.Incident
		fiveMinsAgo := time.Now().Add(-5 * time.Minute)
		err := database.DB.Where("machine_id = ? AND status IN ('OPEN', 'ACKNOWLEDGED') AND updated_at >= ?", machineID, fiveMinsAgo).Order("updated_at desc").First(&incident).Error
		if err != nil {
			return nil, err
		}
		return &incident, nil
	}

	// Fallback in-memory lookup for unit tests without PostgreSQL
	memIncidentsMu.RLock()
	defer memIncidentsMu.RUnlock()

	for _, inc := range memIncidents {
		if inc.MachineID == machineID && (inc.Status == "OPEN" || inc.Status == "ACKNOWLEDGED") {
			return inc, nil
		}
	}

	return nil, fmt.Errorf("no open incident found")
}

func (s *IncidentService) CreateIncident(alert models.LinuxAlert) (*models.Incident, error) {
	now := time.Now()
	rootCause := s.ClassifyRootCause([]string{alert.Category, alert.Title})

	incident := &models.Incident{
		ID:          uuid.New(),
		Title:       fmt.Sprintf("Incident: %s", alert.Title),
		MachineID:   alert.MachineID,
		Severity:    models.IncidentSeverity(alert.Severity),
		Status:      models.IncidentStatusOpen,
		RootCause:   rootCause,
		AlertCount:  1,
		MTTDSeconds: 30,
		StartedAt:   now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if database.DB != nil {
		if err := database.DB.Create(incident).Error; err != nil {
			return nil, err
		}

		incAlert := models.IncidentAlert{
			ID:         uuid.New(),
			IncidentID: incident.ID,
			AlertID:    alert.ID,
			CreatedAt:  now,
		}
		database.DB.Create(&incAlert)
		s.AddTimelineEvent(incident.ID, fmt.Sprintf("Incident triggered by alert: %s (%s)", alert.Title, alert.Severity), "ALERT_ADDED")
	}

	memIncidentsMu.Lock()
	memIncidents[incident.MachineID] = incident
	memIncidentsMu.Unlock()

	s.broadcastIncidentEvent("incident_created", *incident)
	return incident, nil
}

func (s *IncidentService) AddAlertToIncident(incident *models.Incident, alert models.LinuxAlert) error {
	now := time.Now()
	incident.AlertCount++
	incident.UpdatedAt = now

	if (alert.Severity == "Critical" || alert.Severity == "CRITICAL") && incident.Severity != "Critical" {
		incident.Severity = "Critical"
	}

	var categories []string
	categories = append(categories, alert.Category, alert.Title)

	if database.DB != nil {
		var mappedAlerts []models.IncidentAlert
		database.DB.Where("incident_id = ?", incident.ID).Find(&mappedAlerts)
		for _, ma := range mappedAlerts {
			var a models.LinuxAlert
			if database.DB.First(&a, "id = ?", ma.AlertID).Error == nil {
				categories = append(categories, a.Category, a.Title)
			}
		}

		incident.RootCause = s.ClassifyRootCause(categories)
		database.DB.Save(incident)

		incAlert := models.IncidentAlert{
			ID:         uuid.New(),
			IncidentID: incident.ID,
			AlertID:    alert.ID,
			CreatedAt:  now,
		}
		database.DB.Create(&incAlert)

		s.AddTimelineEvent(incident.ID, fmt.Sprintf("Alert correlated into incident: %s (%s)", alert.Title, alert.Severity), "ALERT_ADDED")
	} else {
		incident.RootCause = s.ClassifyRootCause(categories)
	}

	memIncidentsMu.Lock()
	memIncidents[incident.MachineID] = incident
	memIncidentsMu.Unlock()

	s.broadcastIncidentEvent("incident_updated", *incident)
	return nil
}

func (s *IncidentService) ResolveIncident(incidentID uuid.UUID, resolutionNote string) error {
	now := time.Now()
	var incident models.Incident

	if database.DB != nil {
		if err := database.DB.First(&incident, "id = ?", incidentID).Error; err != nil {
			return err
		}

		incident.Status = "RESOLVED"
		incident.ResolvedAt = &now
		mttr := int64(now.Sub(incident.StartedAt).Seconds())
		incident.MTTRSeconds = mttr
		incident.UpdatedAt = now

		database.DB.Save(&incident)
		s.AddTimelineEvent(incident.ID, fmt.Sprintf("Incident resolved: %s", resolutionNote), "RESOLVED")
	} else {
		incident.ID = incidentID
		incident.Status = "RESOLVED"
		incident.ResolvedAt = &now
	}

	memIncidentsMu.Lock()
	for k, inc := range memIncidents {
		if inc.ID == incidentID {
			inc.Status = "RESOLVED"
			inc.ResolvedAt = &now
			delete(memIncidents, k)
		}
	}
	memIncidentsMu.Unlock()

	s.broadcastIncidentEvent("incident_updated", incident)
	return nil
}

func (s *IncidentService) CloseIncident(incidentID uuid.UUID) error {
	now := time.Now()
	var incident models.Incident

	if database.DB != nil {
		if err := database.DB.First(&incident, "id = ?", incidentID).Error; err != nil {
			return err
		}

		incident.Status = "CLOSED"
		incident.ClosedAt = &now
		incident.UpdatedAt = now

		database.DB.Save(&incident)
		s.AddTimelineEvent(incident.ID, "Incident closed by administrator", "CLOSED")
	}

	memIncidentsMu.Lock()
	for k, inc := range memIncidents {
		if inc.ID == incidentID {
			inc.Status = "CLOSED"
			delete(memIncidents, k)
		}
	}
	memIncidentsMu.Unlock()

	s.broadcastIncidentEvent("incident_closed", incident)
	return nil
}

func (s *IncidentService) AddTimelineEvent(incidentID uuid.UUID, message, eventType string) {
	timeline := models.IncidentTimeline{
		ID:         uuid.New(),
		IncidentID: incidentID,
		Message:    message,
		EventType:  eventType,
		CreatedAt:  time.Now(),
	}
	if database.DB != nil {
		database.DB.Create(&timeline)
	}
}

func (s *IncidentService) ClassifyRootCause(terms []string) string {
	combined := strings.ToLower(strings.Join(terms, " "))

	if (strings.Contains(combined, "cpu") && strings.Contains(combined, "memory")) || (strings.Contains(combined, "cpu") && strings.Contains(combined, "ram")) {
		return "Resource Exhaustion (CPU & RAM Saturation)"
	}
	if strings.Contains(combined, "cpu") && strings.Contains(combined, "disk") {
		return "Storage IO Bottleneck & CPU Contention"
	}
	if strings.Contains(combined, "network") || strings.Contains(combined, "latency") || strings.Contains(combined, "packet loss") {
		return "Network Infrastructure Degradation"
	}
	if strings.Contains(combined, "kube") || strings.Contains(combined, "pod") {
		return "Kubernetes Cluster Workload Failure"
	}
	if strings.Contains(combined, "docker") || strings.Contains(combined, "container") {
		return "Container System Overload"
	}
	if strings.Contains(combined, "cpu") {
		return "CPU Compute Saturation"
	}
	if strings.Contains(combined, "memory") || strings.Contains(combined, "ram") {
		return "System Memory Exhaustion"
	}
	if strings.Contains(combined, "disk") || strings.Contains(combined, "storage") {
		return "Storage Filesystem Capacity Limit"
	}

	return "System Hardware / Service Degradation"
}

func (s *IncidentService) GenerateAIIncidentSummary(incident models.Incident) string {
	return fmt.Sprintf("AI Incident Root Cause Analysis Report:\n"+
		"Target Machine: %s\n"+
		"Identified Cause: %s\n"+
		"Impact: Multi-vector alert degradation (%d alerts correlated).\n"+
		"Recommendation: Inspect high CPU/RAM process tree, free disk/swap space, and verify systemd service stability.",
		incident.MachineID.String(), incident.RootCause, incident.AlertCount)
}

func (s *IncidentService) broadcastIncidentEvent(eventType string, incident models.Incident) {
	if websocket.WS != nil {
		websocket.WS.Broadcast(map[string]interface{}{
			"type":        eventType,
			"incident":    incident,
			"id":          incident.ID.String(),
			"machine_id":  incident.MachineID.String(),
			"title":       incident.Title,
			"severity":    incident.Severity,
			"status":      incident.Status,
			"root_cause":  incident.RootCause,
			"alert_count": incident.AlertCount,
		})
	}
}
