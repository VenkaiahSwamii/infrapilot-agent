package automation

import (
	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

type Service struct {
	repo            Repository
	engine          *Engine
	runner          CommandRunner
	runbookEngine   *RunbookEngine
	approvalManager *ApprovalManager
	scheduler       *AutomationScheduler
	notifier        *NotificationEngine
	historyService  *HistoryService
}

func NewService(repo Repository) *Service {
	runner := NewCommandRunner()
	notifier := NewNotificationEngine()
	approval := NewApprovalManager(repo)
	engine := NewEngine(repo, approval, notifier)
	runbookEngine := NewRunbookEngine(repo, runner)
	scheduler := NewAutomationScheduler()
	historyService := NewHistoryService(repo)

	return &Service{
		repo:            repo,
		engine:          engine,
		runner:          runner,
		runbookEngine:   runbookEngine,
		approvalManager: approval,
		scheduler:       scheduler,
		notifier:        notifier,
		historyService:  historyService,
	}
}

func (s *Service) GetRules(orgID string) ([]models.AutomationRule, error) {
	return s.repo.GetRules(orgID)
}

func (s *Service) CreateRule(rule *models.AutomationRule) error {
	return s.repo.CreateRule(rule)
}

func (s *Service) UpdateRule(rule *models.AutomationRule) error {
	return s.repo.UpdateRule(rule)
}

func (s *Service) DeleteRule(id uuid.UUID) error {
	return s.repo.DeleteRule(id)
}

func (s *Service) GetRunbooks(orgID string) ([]models.Runbook, error) {
	return s.repo.GetRunbooks(orgID)
}

func (s *Service) RunAutomation(orgID, target, actionType string) (ExecutionResult, error) {
	return s.runner.ExecuteSSH(target, actionType), nil
}

func (s *Service) GetHistory(orgID string) ([]models.AutomationHistory, error) {
	return s.historyService.GetHistory(orgID)
}

func (s *Service) GetApprovals(orgID string) ([]models.ApprovalRequest, error) {
	return s.repo.GetApprovals(orgID)
}

func (s *Service) ApproveAction(requestID uuid.UUID, approve bool, approver, reason string) error {
	return s.approvalManager.ProcessApproval(requestID, approve, approver, reason)
}
