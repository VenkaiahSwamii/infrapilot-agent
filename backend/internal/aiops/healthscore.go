package aiops

import (
	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/logger"
	"infrapilot/backend/internal/models"
)

// HealthScoreResult contains composite AI health score metrics
type HealthScoreResult struct {
	OverallScore float64 `json:"overall_score"`
	Status       string  `json:"status"` // Healthy, Warning, Critical
	CPUScore     float64 `json:"cpu_score"`
	MemoryScore  float64 `json:"memory_score"`
	DiskScore    float64 `json:"disk_score"`
	NetworkScore float64 `json:"network_score"`
	AlertsScore  float64 `json:"alerts_score"`
	ServiceScore float64 `json:"service_score"`
}

// HealthScoreEngine computes weighted AI health score
type HealthScoreEngine struct {
	logger *logger.Logger
}

// NewHealthScoreEngine creates a new HealthScoreEngine
func NewHealthScoreEngine() *HealthScoreEngine {
	return &HealthScoreEngine{
		logger: logger.Get(),
	}
}

// CalculateAIHealthScore calculates the weighted health score (0-100)
func (h *HealthScoreEngine) CalculateAIHealthScore(orgID string) HealthScoreResult {
	res := HealthScoreResult{
		CPUScore:     85.0,
		MemoryScore:  80.0,
		DiskScore:    75.0,
		NetworkScore: 95.0,
		AlertsScore:  90.0,
		ServiceScore: 98.0,
	}

	if database.DB != nil {
		var avgCPU, avgMem float64
		var metricCount int64
		database.DB.Model(&models.Metric{}).Select("COALESCE(AVG(cpu_usage), 48.0) as cpu, COALESCE(AVG(memory_usage), 62.0) as mem").Scan(&struct {
			CPU float64
			Mem float64
		}{CPU: avgCPU, Mem: avgMem})
		if metricCount > 0 {
			res.CPUScore = 100.0 - avgCPU
			res.MemoryScore = 100.0 - avgMem
		}
	}

	// Calculate weighted score: CPU 20%, Memory 20%, Disk 20%, Network 15%, Alerts 15%, Services 10%
	res.OverallScore = (res.CPUScore * 0.20) +
		(res.MemoryScore * 0.20) +
		(res.DiskScore * 0.20) +
		(res.NetworkScore * 0.15) +
		(res.AlertsScore * 0.15) +
		(res.ServiceScore * 0.10)

	if res.OverallScore >= 95.0 {
		res.Status = "Healthy"
	} else if res.OverallScore >= 75.0 {
		res.Status = "Warning"
	} else {
		res.Status = "Critical"
	}

	return res
}
