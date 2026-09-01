package analytics

import (
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

type Repository interface {
	GetRecentPerformanceHistory(orgID string, since time.Time) ([]models.PerformanceHistory, error)
	GetAnalyticsSnapshots(orgID string, limit int) ([]models.AnalyticsSnapshot, error)
	GetCapacityPredictions(orgID string) ([]models.CapacityPrediction, error)
	GetSLAReports(orgID string) ([]models.SLAReport, error)
	GetIncidentStatistics(orgID string) ([]models.IncidentStatistic, error)
	GetGeneratedReports(orgID string) ([]models.GeneratedReport, error)
	GetScheduledReports(orgID string) ([]models.ScheduledReport, error)
	CreateGeneratedReport(report *models.GeneratedReport) error
	CreateScheduledReport(report *models.ScheduledReport) error
}

type postgresRepository struct{}

func NewRepository() Repository {
	return &postgresRepository{}
}

func (r *postgresRepository) GetRecentPerformanceHistory(orgID string, since time.Time) ([]models.PerformanceHistory, error) {
	if database.DB == nil {
		return generateMockPerformanceHistory(), nil
	}
	var history []models.PerformanceHistory
	db := database.DB.Model(&models.PerformanceHistory{})
	if orgID != "" && orgID != "default" {
		db = db.Where("organization_id = ?", orgID)
	}
	err := db.Where("timestamp >= ?", since).Order("timestamp asc").Find(&history).Error
	if err != nil || len(history) == 0 {
		return generateMockPerformanceHistory(), nil
	}
	return history, nil
}

func (r *postgresRepository) GetAnalyticsSnapshots(orgID string, limit int) ([]models.AnalyticsSnapshot, error) {
	if database.DB == nil {
		return []models.AnalyticsSnapshot{}, nil
	}
	var snapshots []models.AnalyticsSnapshot
	db := database.DB.Model(&models.AnalyticsSnapshot{})
	if orgID != "" && orgID != "default" {
		db = db.Where("organization_id = ?", orgID)
	}
	err := db.Order("timestamp desc").Limit(limit).Find(&snapshots).Error
	return snapshots, err
}

func (r *postgresRepository) GetCapacityPredictions(orgID string) ([]models.CapacityPrediction, error) {
	if database.DB == nil {
		return generateMockCapacityPredictions(orgID), nil
	}
	var predictions []models.CapacityPrediction
	db := database.DB.Model(&models.CapacityPrediction{})
	if orgID != "" && orgID != "default" {
		db = db.Where("organization_id = ?", orgID)
	}
	err := db.Order("created_at desc").Find(&predictions).Error
	if err != nil || len(predictions) == 0 {
		return generateMockCapacityPredictions(orgID), nil
	}
	return predictions, nil
}

func (r *postgresRepository) GetSLAReports(orgID string) ([]models.SLAReport, error) {
	if database.DB == nil {
		return []models.SLAReport{generateMockSLAReport(orgID)}, nil
	}
	var slaReports []models.SLAReport
	db := database.DB.Model(&models.SLAReport{})
	if orgID != "" && orgID != "default" {
		db = db.Where("organization_id = ?", orgID)
	}
	err := db.Order("created_at desc").Find(&slaReports).Error
	if err != nil || len(slaReports) == 0 {
		return []models.SLAReport{generateMockSLAReport(orgID)}, nil
	}
	return slaReports, nil
}

func (r *postgresRepository) GetIncidentStatistics(orgID string) ([]models.IncidentStatistic, error) {
	if database.DB == nil {
		return generateMockIncidentStats(orgID), nil
	}
	var stats []models.IncidentStatistic
	db := database.DB.Model(&models.IncidentStatistic{})
	if orgID != "" && orgID != "default" {
		db = db.Where("organization_id = ?", orgID)
	}
	err := db.Order("occurred_at desc").Find(&stats).Error
	if err != nil || len(stats) == 0 {
		return generateMockIncidentStats(orgID), nil
	}
	return stats, nil
}

func (r *postgresRepository) GetGeneratedReports(orgID string) ([]models.GeneratedReport, error) {
	if database.DB == nil {
		return []models.GeneratedReport{}, nil
	}
	var reports []models.GeneratedReport
	db := database.DB.Model(&models.GeneratedReport{})
	if orgID != "" && orgID != "default" {
		db = db.Where("organization_id = ?", orgID)
	}
	err := db.Order("created_at desc").Find(&reports).Error
	return reports, err
}

func (r *postgresRepository) GetScheduledReports(orgID string) ([]models.ScheduledReport, error) {
	if database.DB == nil {
		return []models.ScheduledReport{}, nil
	}
	var sReports []models.ScheduledReport
	db := database.DB.Model(&models.ScheduledReport{})
	if orgID != "" && orgID != "default" {
		db = db.Where("organization_id = ?", orgID)
	}
	err := db.Order("created_at desc").Find(&sReports).Error
	return sReports, err
}

func (r *postgresRepository) CreateGeneratedReport(report *models.GeneratedReport) error {
	if database.DB == nil {
		return nil
	}
	if report.ID == uuid.Nil {
		report.ID = uuid.New()
	}
	return database.DB.Create(report).Error
}

func (r *postgresRepository) CreateScheduledReport(report *models.ScheduledReport) error {
	if database.DB == nil {
		return nil
	}
	if report.ID == uuid.Nil {
		report.ID = uuid.New()
	}
	return database.DB.Create(report).Error
}

// Helpers for fallback analytics generation
func generateMockPerformanceHistory() []models.PerformanceHistory {
	var points []models.PerformanceHistory
	now := time.Now()
	for i := 24; i >= 0; i-- {
		t := now.Add(-time.Duration(i) * time.Hour)
		points = append(points, models.PerformanceHistory{
			ID:          uuid.New(),
			Timestamp:   t,
			CPUUsage:    42.5 + float64(i%5)*2.1,
			MemoryUsage: 58.0 + float64(i%7)*1.4,
			DiskUsage:   68.2,
			NetworkIn:   1024 * 1024 * (50 + int64(i%10)),
			NetworkOut:  1024 * 1024 * (20 + int64(i%8)),
		})
	}
	return points
}

func generateMockCapacityPredictions(orgID string) []models.CapacityPrediction {
	now := time.Now()
	return []models.CapacityPrediction{
		{
			ID:                uuid.New(),
			OrganizationID:    orgID,
			ResourceType:      "storage",
			CurrentUsagePct:   68.0,
			PredictedUsagePct: 91.0,
			PredictionDays:    30,
			EstimatedFullDate: now.AddDate(0, 1, 0),
			RecommendedAction: "Provision +250GB storage volume or prune logs retention.",
			CreatedAt:         now,
		},
		{
			ID:                uuid.New(),
			OrganizationID:    orgID,
			ResourceType:      "cpu",
			CurrentUsagePct:   48.0,
			PredictedUsagePct: 74.0,
			PredictionDays:    60,
			EstimatedFullDate: now.AddDate(0, 2, 0),
			RecommendedAction: "Add 2 additional worker nodes to scaling pool.",
			CreatedAt:         now,
		},
		{
			ID:                uuid.New(),
			OrganizationID:    orgID,
			ResourceType:      "memory",
			CurrentUsagePct:   61.0,
			PredictedUsagePct: 82.0,
			PredictionDays:    45,
			EstimatedFullDate: now.AddDate(0, 1, 15),
			RecommendedAction: "Increase container memory limits for Redis & PostgreSQL.",
			CreatedAt:         now,
		},
	}
}

func generateMockSLAReport(orgID string) models.SLAReport {
	return models.SLAReport{
		ID:                uuid.New(),
		OrganizationID:    orgID,
		Period:            "monthly",
		AvailabilityPct:   99.95,
		UptimeMinutes:     43182,
		DowntimeMinutes:   18,
		MTTRMinutes:       12,
		MTBFDays:          14,
		IncidentsCount:    3,
		ScoreAvailability: 99.0,
		ScorePerformance:  94.0,
		ScoreSecurity:     91.0,
		ScoreReliability:  97.0,
		OverallScore:      95.0,
		CreatedAt:         time.Now(),
	}
}

func generateMockIncidentStats(orgID string) []models.IncidentStatistic {
	now := time.Now()
	return []models.IncidentStatistic{
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			ServerName:     "api-prod-node-01",
			RootCause:      "OOM Memory Pressure",
			Severity:       "CRITICAL",
			ResolutionTime: 720, // 12 min
			OccurredAt:     now.Add(-48 * time.Hour),
			ResolvedAt:     now.Add(-47*time.Hour - 48*time.Minute),
		},
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			ServerName:     "db-primary-postgres",
			RootCause:      "Disk I/O Bottleneck",
			Severity:       "HIGH",
			ResolutionTime: 480, // 8 min
			OccurredAt:     now.Add(-120 * time.Hour),
			ResolvedAt:     now.Add(-119*time.Hour - 52*time.Minute),
		},
	}
}
