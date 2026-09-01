package analytics

import (
	"infrapilot/backend/internal/models"
)

type SLAEngine struct {
	repo Repository
}

func NewSLAEngine(repo Repository) *SLAEngine {
	return &SLAEngine{repo: repo}
}

func (s *SLAEngine) GetSLASummary(orgID string) (models.SLAReport, error) {
	reports, err := s.repo.GetSLAReports(orgID)
	if err != nil || len(reports) == 0 {
		return generateMockSLAReport(orgID), nil
	}
	return reports[0], nil
}
