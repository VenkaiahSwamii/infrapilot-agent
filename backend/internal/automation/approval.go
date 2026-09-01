package automation

import (
	"time"

	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

type ApprovalManager struct {
	repo Repository
}

func NewApprovalManager(repo Repository) *ApprovalManager {
	return &ApprovalManager{repo: repo}
}

func (a *ApprovalManager) CreateRequest(orgID string, runID uuid.UUID, ruleName, action, targetResource, requestedBy string) (*models.ApprovalRequest, error) {
	req := &models.ApprovalRequest{
		ID:             uuid.New(),
		OrganizationID: orgID,
		RunID:          runID,
		RuleName:       ruleName,
		Action:         action,
		TargetResource: targetResource,
		RequestedBy:    requestedBy,
		Status:         "PENDING",
		CreatedAt:      time.Now(),
	}

	if err := a.repo.CreateApproval(req); err != nil {
		return nil, err
	}
	return req, nil
}

func (a *ApprovalManager) ProcessApproval(requestID uuid.UUID, approve bool, approver, reason string) error {
	approvals, err := a.repo.GetApprovals("")
	if err != nil {
		return err
	}

	var target *models.ApprovalRequest
	for _, item := range approvals {
		if item.ID == requestID {
			target = &item
			break
		}
	}

	if target == nil {
		target = &models.ApprovalRequest{ID: requestID, RuleName: "Approved Action", Action: "Execute"}
	}

	now := time.Now()
	if approve {
		target.Status = "APPROVED"
	} else {
		target.Status = "REJECTED"
	}
	target.Approver = approver
	target.Reason = reason
	target.RespondedAt = &now

	return a.repo.UpdateApproval(target)
}
