package apm

import "infrapilot/backend/internal/models"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RecordAPIRequest(req *models.APIRequest) error {
	return s.repo.RecordAPIRequest(req)
}

func (s *Service) GetAPMData() (*APMDataResponse, error) {
	overview, err := s.repo.GetOverviewStats()
	if err != nil {
		return nil, err
	}

	endpoints, _ := s.repo.GetEndpoints(10)
	slowEndpoints, _ := s.repo.GetSlowEndpoints(5)
	topAPIs, _ := s.repo.GetTopAPIs(5)
	services, _ := s.repo.GetServiceMetrics()
	aiRecs := GenerateAIRecommendations(slowEndpoints)

	return &APMDataResponse{
		Overview:        *overview,
		Endpoints:       endpoints,
		SlowEndpoints:   slowEndpoints,
		TopAPIs:         topAPIs,
		Recommendations: aiRecs,
		Services:        services,
	}, nil
}
