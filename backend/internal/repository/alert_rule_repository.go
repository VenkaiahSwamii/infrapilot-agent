package repository

import (
	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

type AlertRuleRepository struct{}

func NewAlertRuleRepository() *AlertRuleRepository {
	return &AlertRuleRepository{}
}

func (r *AlertRuleRepository) Create(rule *models.AlertRule) error {
	return database.DB.Create(rule).Error
}

func (r *AlertRuleRepository) Update(rule *models.AlertRule) error {
	return database.DB.Save(rule).Error
}

func (r *AlertRuleRepository) Delete(id uuid.UUID) error {
	return database.DB.Delete(&models.AlertRule{}, id).Error
}

func (r *AlertRuleRepository) GetByID(id uuid.UUID) (*models.AlertRule, error) {
	var rule models.AlertRule
	err := database.DB.First(&rule, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *AlertRuleRepository) GetAll(org string) ([]models.AlertRule, error) {
	var rules []models.AlertRule
	err := database.DB.Where("organization_id = ?", org).
		Order("created_at DESC").
		Find(&rules).Error
	return rules, err
}

func (r *AlertRuleRepository) GetEnabledRules(org string) ([]models.AlertRule, error) {
	var rules []models.AlertRule
	err := database.DB.Where("organization_id = ? AND is_enabled = ?", org, true).
		Order("created_at DESC").
		Find(&rules).Error
	return rules, err
}

func (r *AlertRuleRepository) GetAllEnabledRules() ([]models.AlertRule, error) {
	var rules []models.AlertRule
	err := database.DB.Where("is_enabled = ?", true).
		Order("created_at DESC").
		Find(&rules).Error
	return rules, err
}
