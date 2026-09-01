package cloud

import (
	"infrapilot/backend/internal/models"
)

type SecurityEngine struct {
	repo Repository
}

func NewSecurityEngine(repo Repository) *SecurityEngine {
	return &SecurityEngine{repo: repo}
}

func (s *SecurityEngine) GetFindings(orgID string) ([]models.CloudSecurityFinding, error) {
	return s.repo.GetSecurityFindings(orgID)
}
