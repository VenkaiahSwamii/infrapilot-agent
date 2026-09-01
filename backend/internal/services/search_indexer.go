package services

import (
	"encoding/json"

	"infrapilot/backend/internal/models"
	"infrapilot/backend/internal/repository"

	"github.com/google/uuid"
)

// SearchIndexer handles indexing of resources for search
type SearchIndexer struct {
	searchRepo *repository.SearchRepository
}

// NewSearchIndexer creates a new search indexer
func NewSearchIndexer(searchRepo *repository.SearchRepository) *SearchIndexer {
	return &SearchIndexer{
		searchRepo: searchRepo,
	}
}

// IndexMachine indexes a machine for search
func (i *SearchIndexer) IndexMachine(orgID uuid.UUID, machine models.Machine) error {
	title := machine.Hostname
	if title == "" {
		title = "Machine " + machine.ID.String()
	}

	description := ""
	if machine.OS != "" {
		description += "OS: " + machine.OS + ". "
	}
	if machine.Platform != "" {
		description += "Platform: " + machine.Platform + ". "
	}
	if machine.Status != "" {
		description += "Status: " + machine.Status + ". "
	}

	keywords := []string{
		machine.Hostname,
		machine.OS,
		machine.Platform,
		machine.Virtualization,
		machine.CloudProvider,
		machine.Status,
	}

	metadata, _ := json.Marshal(map[string]interface{}{
		"hostname":  machine.Hostname,
		"os":        machine.OS,
		"platform":  machine.Platform,
		"status":    machine.Status,
		"last_seen": machine.LastSeen,
	})

	return i.searchRepo.Index(
		"machine",
		machine.ID,
		title,
		description,
		keywords,
		orgID,
		"infrastructure",
		"",
		machine.Status,
		"/machines/"+machine.ID.String(),
		metadata,
	)
}

// IndexIncident indexes an incident for search
func (i *SearchIndexer) IndexIncident(orgID uuid.UUID, incident models.Incident) error {
	title := incident.Title
	description := incident.Description
	if description == "" {
		description = "Incident " + string(incident.Severity) + " - " + string(incident.Status)
	}

	keywords := []string{
		incident.Title,
		string(incident.Severity),
		string(incident.Status),
		string(incident.Source),
	}

	metadata, _ := json.Marshal(map[string]interface{}{
		"severity":    string(incident.Severity),
		"status":      string(incident.Status),
		"source":      string(incident.Source),
		"created_at":  incident.CreatedAt,
		"assigned_to": incident.AssignedTo,
	})

	return i.searchRepo.Index(
		"incident",
		incident.ID,
		title,
		description,
		keywords,
		orgID,
		"incidents",
		string(incident.Severity),
		string(incident.Status),
		"/incidents/"+incident.ID.String(),
		metadata,
	)
}

// IndexAlert indexes an alert for search
func (i *SearchIndexer) IndexAlert(orgID uuid.UUID, alert models.LinuxAlert) error {
	title := alert.ID.String()
	if title == "" {
		title = "Alert"
	}

	description := alert.Message
	if description == "" {
		description = "Alert severity: " + alert.Severity
	}

	keywords := []string{
		alert.Severity,
		alert.Status,
		alert.Message,
	}

	metadata, _ := json.Marshal(map[string]interface{}{
		"severity":   alert.Severity,
		"status":     alert.Status,
		"message":    alert.Message,
		"created_at": alert.CreatedAt,
	})

	return i.searchRepo.Index(
		"alert",
		alert.ID,
		title,
		description,
		keywords,
		orgID,
		"alerts",
		alert.Severity,
		alert.Status,
		"/alerts/"+alert.ID.String(),
		metadata,
	)
}

// IndexLog indexes a log entry for search
func (i *SearchIndexer) IndexLog(orgID uuid.UUID, machineID uuid.UUID, log models.LinuxLog) error {
	title := log.Source + " - " + log.Level
	description := log.Message

	keywords := []string{
		log.Source,
		log.Level,
		log.Message,
	}

	metadata, _ := json.Marshal(map[string]interface{}{
		"source":     log.Source,
		"level":      log.Level,
		"message":    log.Message,
		"timestamp":  log.Timestamp,
		"machine_id": machineID,
	})

	return i.searchRepo.Index(
		"log",
		log.ID,
		title,
		description,
		keywords,
		orgID,
		"logs",
		log.Level,
		"",
		"/machines/"+machineID.String()+"/logs",
		metadata,
	)
}

// RemoveFromIndex removes a resource from the search index
func (i *SearchIndexer) RemoveFromIndex(resourceType string, resourceID uuid.UUID, orgID uuid.UUID) error {
	return i.searchRepo.RemoveFromIndex(resourceType, resourceID, orgID)
}

// ReindexAll reindexes all resources (admin operation)
func (i *SearchIndexer) ReindexAll(orgID uuid.UUID) error {
	// This would be a long-running operation
	// In production, this would use a background worker
	// For now, this is a placeholder
	return nil
}
