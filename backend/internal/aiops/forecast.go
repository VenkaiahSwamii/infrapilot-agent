package aiops

import (
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/logger"
	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

// ForecastEngine manages capacity projections across 7-day, 30-day, and 90-day horizons
type ForecastEngine struct {
	logger *logger.Logger
}

// NewForecastEngine creates a new ForecastEngine
func NewForecastEngine() *ForecastEngine {
	return &ForecastEngine{
		logger: logger.Get(),
	}
}

// GetCapacityForecasts returns resource capacity trends and expansion recommendations
func (f *ForecastEngine) GetCapacityForecasts(orgID string) ([]models.CapacityForecastRecord, error) {
	var forecasts []models.CapacityForecastRecord

	if database.DB != nil {
		database.DB.Where("organization_id = ? OR organization_id = ''", orgID).Find(&forecasts)
	}

	if len(forecasts) == 0 {
		forecasts = f.generateSampleForecasts(orgID)
	}

	return forecasts, nil
}

func (f *ForecastEngine) generateSampleForecasts(orgID string) []models.CapacityForecastRecord {
	now := time.Now()
	return []models.CapacityForecastRecord{
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			ResourceType:   "CPU Fleet Utilization",
			CurrentUsage:   48.2,
			Forecast7D:     54.0,
			Forecast30D:    68.5,
			Forecast90D:    89.0,
			Recommendation: "Plan to add 2 worker nodes before 90-day window to maintain < 70% CPU headroom",
			UpdatedAt:      now,
		},
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			ResourceType:   "RAM Memory Utilization",
			CurrentUsage:   62.0,
			Forecast7D:     65.5,
			Forecast30D:    78.0,
			Forecast90D:    94.5,
			Recommendation: "Upgrade server RAM nodes from 64GB to 128GB within 30 days",
			UpdatedAt:      now,
		},
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			ResourceType:   "Persistent Volume Storage",
			CurrentUsage:   71.5,
			Forecast7D:     76.0,
			Forecast30D:    86.2,
			Forecast90D:    98.0,
			Recommendation: "Expand Postgres & Qdrant PVC disk sizes by 500GB prior to day 30",
			UpdatedAt:      now,
		},
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			ResourceType:   "Kubernetes Pod Count",
			CurrentUsage:   142,
			Forecast7D:     155,
			Forecast30D:    190,
			Forecast90D:    260,
			Recommendation: "Scale EKS/GKE cluster node pool by 3 additional instances",
			UpdatedAt:      now,
		},
	}
}
