package repository

import (
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

type AlertRepository struct{}

func NewAlertRepository() *AlertRepository {
	return &AlertRepository{}
}

func (r *AlertRepository) Create(alert *models.LinuxAlert) error {
	return database.DB.Create(alert).Error
}

func (r *AlertRepository) FindOpenByMachineAndRule(machineID uuid.UUID, ruleID uuid.UUID) (*models.LinuxAlert, error) {
	var alert models.LinuxAlert
	err := database.DB.Where("machine_id = ? AND rule_id = ? AND status = ?", machineID, ruleID, "OPEN").First(&alert).Error
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

func (r *AlertRepository) ListOpenByMachine(machineID uuid.UUID) ([]models.LinuxAlert, error) {
	var alerts []models.LinuxAlert
	err := database.DB.Where("machine_id = ? AND status = ?", machineID, "OPEN").
		Order("created_at DESC").
		Find(&alerts).Error
	return alerts, err
}

func (r *AlertRepository) ResolveAlert(alertID uuid.UUID, resolutionNote string, resolvedBy string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":          "RESOLVED",
		"resolution_note": resolutionNote,
		"resolved_by":     resolvedBy,
		"resolved_at":     now,
		"updated_at":      now,
	}
	return database.DB.Model(&models.LinuxAlert{}).Where("id = ?", alertID).Updates(updates).Error
}

func (r *AlertRepository) GetByID(alertID uuid.UUID) (*models.LinuxAlert, error) {
	var alert models.LinuxAlert
	err := database.DB.Where("id = ?", alertID).First(&alert).Error
	if err != nil {
		return nil, err
	}
	return &alert, nil
}
