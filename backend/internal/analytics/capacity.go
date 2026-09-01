package analytics

import (
	"infrapilot/backend/internal/models"
)

type CapacityEngine struct {
	repo Repository
}

func NewCapacityEngine(repo Repository) *CapacityEngine {
	return &CapacityEngine{repo: repo}
}

func (c *CapacityEngine) GetForecasts(orgID string) ([]models.CapacityPrediction, error) {
	return c.repo.GetCapacityPredictions(orgID)
}
