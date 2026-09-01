package automation

import (
	"fmt"
	"time"

	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

type Engine struct {
	repo     Repository
	runner   CommandRunner
	approval *ApprovalManager
	notifier *NotificationEngine
}

func NewEngine(repo Repository, approval *ApprovalManager, notifier *NotificationEngine) *Engine {
	return &Engine{
		repo:     repo,
		runner:   NewCommandRunner(),
		approval: approval,
		notifier: notifier,
	}
}

func (e *Engine) EvaluateAlertAndRemediate(orgID string, eventType string, host string) (*models.AutomationHistory, error) {
	rules, err := e.repo.GetRules(orgID)
	if err != nil || len(rules) == 0 {
		return nil, fmt.Errorf("no matching automation rule found")
	}

	var targetRule *models.AutomationRule
	for _, r := range rules {
		if r.TriggerEventType == eventType && r.Enabled {
			targetRule = &r
			break
		}
	}

	if targetRule == nil {
		targetRule = &rules[0]
	}

	// Check if approval is required
	if targetRule.RequiresApproval {
		_, err := e.approval.CreateRequest(orgID, uuid.New(), targetRule.Name, targetRule.ActionType, targetRule.TargetResource, "AIOps Engine")
		if err != nil {
			return nil, err
		}
		_ = e.notifier.NotifyAll("Approval Required", fmt.Sprintf("Rule '%s' requires admin approval before executing on %s", targetRule.Name, host))
		return nil, fmt.Errorf("remediation action requires approval: request submitted")
	}

	// Execute action
	start := time.Now()
	res := e.runner.ExecuteSSH(host, targetRule.TargetResource)

	history := &models.AutomationHistory{
		ID:             uuid.New(),
		OrganizationID: orgID,
		RuleName:       targetRule.Name,
		Trigger:        eventType,
		User:           "AIOps Engine",
		Action:         targetRule.ActionType,
		Target:         host,
		StartTime:      start,
		EndTime:        time.Now(),
		Result:         res.Status,
		OutputSnippet:  res.Output,
	}

	_ = e.repo.CreateHistory(history)
	_ = e.notifier.NotifyAll("Automation Completed", fmt.Sprintf("Rule '%s' executed on %s with status %s", targetRule.Name, host, res.Status))

	return history, nil
}
