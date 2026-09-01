package services

import (
	"log"

	"infrapilot/backend/internal/models"
)

type CorrelationEngine struct {
	incidentService *IncidentService
}

func NewCorrelationEngine() *CorrelationEngine {
	return &CorrelationEngine{
		incidentService: NewIncidentService(),
	}
}

func (c *CorrelationEngine) Correlate(alert models.LinuxAlert) (*models.Incident, error) {
	// 1. Check if an active open incident already exists for this machine within correlation window
	existingIncident, err := c.incidentService.FindOpenIncident(alert.MachineID)
	if err == nil && existingIncident != nil {
		log.Printf("[CorrelationEngine] Correlating alert %s into existing Incident #%s", alert.Title, existingIncident.ID)
		err = c.incidentService.AddAlertToIncident(existingIncident, alert)
		return existingIncident, err
	}

	// 2. Otherwise create a new Incident
	log.Printf("[CorrelationEngine] Creating new Incident for alert %s on machine %s", alert.Title, alert.MachineID)
	newIncident, err := c.incidentService.CreateIncident(alert)
	return newIncident, err
}
