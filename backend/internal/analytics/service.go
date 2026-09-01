package analytics

import (
	"infrapilot/backend/internal/models"
)

type Service struct {
	repo           Repository
	trendEngine    *TrendEngine
	capacityEngine *CapacityEngine
	slaEngine      *SLAEngine
	reportEngine   *ReportEngine
}

func NewService(repo Repository) *Service {
	return &Service{
		repo:           repo,
		trendEngine:    NewTrendEngine(repo),
		capacityEngine: NewCapacityEngine(repo),
		slaEngine:      NewSLAEngine(repo),
		reportEngine:   NewReportEngine(repo),
	}
}

func (s *Service) GetOverview(orgID string) (map[string]interface{}, error) {
	sla, _ := s.slaEngine.GetSLASummary(orgID)

	return map[string]interface{}{
		"avg_cpu_percent":    48.0,
		"avg_memory_percent": 61.0,
		"availability_pct":   sla.AvailabilityPct,
		"alerts_today":       14,
		"critical_incidents": 2,
		"total_servers":      480,
		"online_servers":     472,
		"offline_servers":    8,
		"docker_containers":  3200,
		"k8s_clusters":       28,
		"scorecard": map[string]interface{}{
			"availability": sla.ScoreAvailability,
			"performance":  sla.ScorePerformance,
			"security":     sla.ScoreSecurity,
			"reliability":  sla.ScoreReliability,
			"overall":      sla.OverallScore,
		},
	}, nil
}

func (s *Service) GetTrends(orgID, rangeStr string) ([]TrendPoint, error) {
	return s.trendEngine.GetTrends(orgID, rangeStr)
}

func (s *Service) GetCapacity(orgID string) ([]models.CapacityPrediction, error) {
	return s.capacityEngine.GetForecasts(orgID)
}

func (s *Service) GetSLA(orgID string) (models.SLAReport, error) {
	return s.slaEngine.GetSLASummary(orgID)
}

func (s *Service) GetIncidents(orgID string) ([]models.IncidentStatistic, error) {
	return s.repo.GetIncidentStatistics(orgID)
}

func (s *Service) GetReports(orgID string) ([]models.GeneratedReport, error) {
	return s.reportEngine.ListGenerated(orgID)
}

func (s *Service) GenerateReport(orgID, name, reportType, format, username string) (*models.GeneratedReport, error) {
	return s.reportEngine.GenerateReport(orgID, name, reportType, format, username)
}

func (s *Service) ScheduleReport(orgID, name, reportType, format, frequency, recipients string) (*models.ScheduledReport, error) {
	return s.reportEngine.ScheduleReport(orgID, name, reportType, format, frequency, recipients)
}
