package cloud

import (
	"infrapilot/backend/internal/models"
)

type AlertEngine struct {
	repo Repository
}

func NewAlertEngine(repo Repository) *AlertEngine {
	return &AlertEngine{repo: repo}
}

func (a *AlertEngine) GetAlerts(orgID string) ([]models.CloudAlert, error) {
	return a.repo.GetAlerts(orgID)
}
