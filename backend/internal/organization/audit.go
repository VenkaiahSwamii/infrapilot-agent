package organization

import (
	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

type AuditService struct {
	repo Repository
}

func NewAuditService(repo Repository) *AuditService {
	return &AuditService{repo: repo}
}

func (a *AuditService) Log(orgID uuid.UUID, username, ip, action, resource, result, details string) error {
	entry := &models.OrganizationAuditLog{
		OrganizationID: orgID,
		Username:       username,
		IP:             ip,
		Action:         action,
		Resource:       resource,
		Result:         result,
		Details:        details,
	}
	return a.repo.LogAudit(entry)
}

func (a *AuditService) GetLogs(orgID uuid.UUID, limit, offset int) ([]models.OrganizationAuditLog, int64, error) {
	if limit <= 0 {
		limit = 50
	}
	return a.repo.GetAuditLogs(orgID, limit, offset)
}
