package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"
	"infrapilot/backend/internal/repository"

	"github.com/google/uuid"
)

var (
	ErrRuleNotFound    = errors.New("alert rule not found")
	ErrInvalidMetric   = errors.New("invalid metric")
	ErrInvalidOperator = errors.New("invalid operator")
	ErrInvalidSeverity = errors.New("invalid severity")
	ErrDuplicateName   = errors.New("alert rule with same name already exists")
)

type AlertRuleService struct {
	repo *repository.AlertRuleRepository
}

func NewAlertRuleService() *AlertRuleService {
	return &AlertRuleService{
		repo: repository.NewAlertRuleRepository(),
	}
}

var allowedMetrics = map[string]bool{
	"cpu_usage":       true,
	"memory_percent":  true,
	"disk_percent":    true,
	"latency_ms":      true,
	"cpu_temperature": true,
	"packet_loss":     true,
	"disk_read_bps":   true,
	"disk_write_bps":  true,
	"upload_mbps":     true,
	"download_mbps":   true,
}

var allowedOperators = map[string]bool{
	">":  true,
	">=": true,
	"<":  true,
	"<=": true,
	"==": true,
	"!=": true,
}

var allowedSeverities = map[string]bool{
	"critical": true,
	"warning":  true,
	"info":     true,
}

type CreateRuleRequest struct {
	OrganizationID string  `json:"organization_id" binding:"required"`
	Name           string  `json:"name" binding:"required"`
	Metric         string  `json:"metric" binding:"required"`
	Operator       string  `json:"operator" binding:"required"`
	Threshold      float64 `json:"threshold" binding:"required"`
	Severity       string  `json:"severity" binding:"required"`
	Enabled        bool    `json:"enabled"`
}

type UpdateRuleRequest struct {
	Name      string  `json:"name"`
	Metric    string  `json:"metric"`
	Operator  string  `json:"operator"`
	Threshold float64 `json:"threshold"`
	Severity  string  `json:"severity"`
	Enabled   *bool   `json:"enabled"`
}

func (s *AlertRuleService) validateRuleFields(metric, operator, severity string) error {
	if !allowedMetrics[strings.ToLower(metric)] {
		return ErrInvalidMetric
	}
	if !allowedOperators[operator] {
		return ErrInvalidOperator
	}
	if !allowedSeverities[strings.ToLower(severity)] {
		return ErrInvalidSeverity
	}
	return nil
}

func (s *AlertRuleService) CreateRule(req CreateRuleRequest) (*models.AlertRule, error) {
	if err := s.validateRuleFields(req.Metric, req.Operator, req.Severity); err != nil {
		return nil, err
	}

	existing, _ := s.repo.GetAll(req.OrganizationID)
	for _, r := range existing {
		if strings.EqualFold(r.Name, req.Name) {
			return nil, ErrDuplicateName
		}
	}

	rule := &models.AlertRule{
		ID:             uuid.New(),
		OrganizationID: req.OrganizationID,
		Name:           req.Name,
		Metric:         strings.ToLower(req.Metric),
		Operator:       req.Operator,
		Threshold:      req.Threshold,
		Severity:       strings.ToLower(req.Severity),
		IsEnabled:      req.Enabled,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.repo.Create(rule); err != nil {
		return nil, err
	}

	return rule, nil
}

func (s *AlertRuleService) UpdateRule(id uuid.UUID, req UpdateRuleRequest) (*models.AlertRule, error) {
	rule, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrRuleNotFound
	}

	if req.Name != "" {
		existing, _ := s.repo.GetAll(rule.OrganizationID)
		for _, r := range existing {
			if r.ID != id && strings.EqualFold(r.Name, req.Name) {
				return nil, ErrDuplicateName
			}
		}
		rule.Name = req.Name
	}

	if req.Metric != "" {
		rule.Metric = strings.ToLower(req.Metric)
	}
	if req.Operator != "" {
		rule.Operator = req.Operator
	}
	if req.Severity != "" {
		rule.Severity = strings.ToLower(req.Severity)
	}
	if req.Enabled != nil {
		rule.IsEnabled = *req.Enabled
	}
	if req.Threshold != 0 {
		rule.Threshold = req.Threshold
	}

	if err := s.validateRuleFields(rule.Metric, rule.Operator, rule.Severity); err != nil {
		return nil, err
	}

	rule.UpdatedAt = time.Now()
	if err := s.repo.Update(rule); err != nil {
		return nil, err
	}

	return rule, nil
}

func (s *AlertRuleService) DeleteRule(id uuid.UUID) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return ErrRuleNotFound
	}
	return s.repo.Delete(id)
}

func (s *AlertRuleService) EnableRule(id uuid.UUID) (*models.AlertRule, error) {
	rule, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrRuleNotFound
	}

	if rule.IsEnabled {
		return rule, nil
	}

	rule.IsEnabled = true
	rule.UpdatedAt = time.Now()
	if err := s.repo.Update(rule); err != nil {
		return nil, err
	}
	return rule, nil
}

func (s *AlertRuleService) DisableRule(id uuid.UUID) (*models.AlertRule, error) {
	rule, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrRuleNotFound
	}

	if !rule.IsEnabled {
		return rule, nil
	}

	rule.IsEnabled = false
	rule.UpdatedAt = time.Now()
	if err := s.repo.Update(rule); err != nil {
		return nil, err
	}
	return rule, nil
}

func (s *AlertRuleService) ListRules(org string) ([]models.AlertRule, error) {
	rules, err := s.repo.GetAll(org)
	if err != nil {
		return nil, err
	}
	return rules, nil
}

func (s *AlertRuleService) EvaluateRule(rule *models.AlertRule, value float64) bool {
	if !rule.IsEnabled {
		return false
	}

	switch rule.Operator {
	case ">":
		return value > rule.Threshold
	case ">=":
		return value >= rule.Threshold
	case "<":
		return value < rule.Threshold
	case "<=":
		return value <= rule.Threshold
	case "==":
		return value == rule.Threshold
	case "!=":
		return value != rule.Threshold
	default:
		return false
	}
}

func (s *AlertRuleService) CheckAlerts(machineID uuid.UUID, metrics map[string]float64) {
	rules, err := s.repo.GetEnabledRules(machineID.String())
	if err != nil {
		return
	}

	for _, rule := range rules {
		value, ok := metrics[rule.Metric]
		if !ok {
			continue
		}

		if s.EvaluateRule(&rule, value) {
			if database.DB != nil {
				var existing models.LinuxAlert
				if err := database.DB.Where("machine_id = ? AND (rule_id = ? OR title = ?) AND (status = 'ACTIVE' OR status = 'OPEN')", machineID, rule.ID, rule.Name).First(&existing).Error; err == nil {
					// Update existing alert timestamp and metric value without creating duplicates
					database.DB.Model(&models.LinuxAlert{}).Where("id = ?", existing.ID).Updates(map[string]interface{}{
						"metric_value": value,
						"updated_at":   time.Now(),
					})
					continue
				}
			}

			alert := models.LinuxAlert{
				ID:          uuid.New(),
				MachineID:   machineID,
				Type:        rule.Name,
				Title:       rule.Name,
				Category:    rule.Metric,
				Severity:    rule.Severity,
				Priority:    models.MapSeverityToPriority(rule.Severity),
				Message:     fmt.Sprintf("%s breached threshold: %.2f", rule.Name, value),
				Status:      "ACTIVE",
				RuleID:      rule.ID,
				MetricValue: value,
				Threshold:   rule.Threshold,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}

			if database.DB != nil {
				database.DB.Create(&alert)
			}
			SendAlert(alert)
		}
	}
}
