package organization

import (
	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

type SettingsService struct {
	repo Repository
}

func NewSettingsService(repo Repository) *SettingsService {
	return &SettingsService{repo: repo}
}

func (s *SettingsService) Get(orgID uuid.UUID) (*models.OrganizationSettings, error) {
	return s.repo.GetSettings(orgID)
}

func (s *SettingsService) Update(orgID uuid.UUID, settings *models.OrganizationSettings) (*models.OrganizationSettings, error) {
	settings.OrganizationID = orgID
	if err := s.repo.UpdateSettings(settings); err != nil {
		return nil, err
	}
	return s.repo.GetSettings(orgID)
}
