package automation

import (
	"infrapilot/backend/internal/models"
)

type HistoryService struct {
	repo Repository
}

func NewHistoryService(repo Repository) *HistoryService {
	return &HistoryService{repo: repo}
}

func (h *HistoryService) GetHistory(orgID string) ([]models.AutomationHistory, error) {
	return h.repo.GetHistory(orgID)
}
