package tracing

import (
	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RecordTrace(trace *models.Trace) error {
	return s.repo.CreateTrace(trace)
}

func (s *Service) GetTraces(serviceName, query string, onlySlow bool, limit, offset int) ([]models.Trace, int64, error) {
	return s.repo.GetTraces(serviceName, query, onlySlow, limit, offset)
}

func (s *Service) GetTraceByID(id uuid.UUID) (*models.Trace, error) {
	return s.repo.GetTraceByID(id)
}

func (s *Service) GetServiceTopology() (map[string]interface{}, error) {
	return s.repo.GetServiceTopology()
}
