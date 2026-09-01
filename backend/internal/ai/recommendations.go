package ai

import (
	"net/http"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetRecommendationsHandler returns active recommendations.
func GetRecommendationsHandler(c *gin.Context) {
	db := database.DB
	if db == nil {
		c.JSON(http.StatusOK, generateSampleRecommendations())
		return
	}

	var recommendations []models.AIRecommendation
	err := db.Order("created_at desc").Find(&recommendations).Error
	if err != nil || len(recommendations) == 0 {
		c.JSON(http.StatusOK, generateSampleRecommendations())
		return
	}

	c.JSON(http.StatusOK, recommendations)
}

// PredictionRequest represents the request for predictions
type PredictionRequest struct {
	MachineID string `json:"machine_id"`
	Metric    string `json:"metric"` // cpu, memory, disk
	Days      int    `json:"days"`
}

// PredictionResponse represents prediction results
type PredictionResponse struct {
	MachineID      string  `json:"machine_id"`
	Metric         string  `json:"metric"`
	CurrentValue   float64 `json:"current_value"`
	PredictedValue float64 `json:"predicted_value"`
	GrowthRate     float64 `json:"growth_rate"`
	DaysRemaining  int     `json:"days_remaining"`
	Threshold      float64 `json:"threshold"`
	Confidence     float64 `json:"confidence"`
	Message        string  `json:"message"`
}

// GetPredictionsHandler returns AI-generated predictions for resource exhaustion
func GetPredictionsHandler(c *gin.Context) {
	// Placeholder implementation returning sample predictions
	predictions := []PredictionResponse{
		{
			MachineID:      "srv-001",
			Metric:         "disk",
			CurrentValue:   89.0,
			PredictedValue: 95.0,
			GrowthRate:     0.03,
			DaysRemaining:  4,
			Threshold:      95.0,
			Confidence:     0.92,
			Message:        "Disk usage predicted to reach 95% in 4 days",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"predictions": predictions,
	})
}

func generateSampleRecommendations() []models.AIRecommendation {
	now := time.Now()
	srvID := uuid.NewSHA1(uuid.NameSpaceDNS, []byte("server-01"))

	return []models.AIRecommendation{
		{
			ID:        uuid.New(),
			ServerID:  srvID,
			Hostname:  "web-prod-01",
			Title:     "Scale Ingress Controller Deployment",
			Reason:    "CPU usage spikes consistently every day between 9 AM and 11 AM due to peak morning traffic, reaching 94.2%. Increasing replicas will balance load.",
			Priority:  "High",
			CreatedAt: now.Add(-1 * time.Hour),
		},
		{
			ID:        uuid.New(),
			ServerID:  srvID,
			Hostname:  "db-primary",
			Title:     "Increase RAM Allocation",
			Reason:    "Memory has exceeded 90% utilization continuously for the last 7 days. SQL buffer cache hit ratio dropped to 72%.",
			Priority:  "High",
			CreatedAt: now.Add(-12 * time.Hour),
		},
		{
			ID:        uuid.New(),
			ServerID:  srvID,
			Hostname:  "api-gateway",
			Title:     "Rotate Application Container logs",
			Reason:    "Disk space consumption shows a linear growth of 4.2GB daily. Log rotation config is disabled.",
			Priority:  "Medium",
			CreatedAt: now.Add(-24 * time.Hour),
		},
	}
}
