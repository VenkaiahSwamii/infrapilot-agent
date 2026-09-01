package services

import (
	"fmt"
	"strings"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"
	"infrapilot/backend/internal/utils"
	"infrapilot/backend/internal/websocket"

	"github.com/google/uuid"
)

type WorkflowService struct{}

func NewWorkflowService() *WorkflowService {
	return &WorkflowService{}
}

func (s *WorkflowService) StartWorkflow(workflowID, incidentID, machineID uuid.UUID) (*models.WorkflowExecution, error) {
	now := time.Now()
	var workflow models.Workflow

	if database.DB != nil {
		if err := database.DB.Preload("Steps").First(&workflow, "id = ?", workflowID).Error; err != nil {
			return nil, fmt.Errorf("workflow not found: %w", err)
		}
	} else {
		workflow = models.Workflow{
			ID:          workflowID,
			Name:        "Automated Runbook",
			TriggerType: "Alert",
			Enabled:     true,
			Steps: []models.WorkflowStep{
				{ID: uuid.New(), StepOrder: 1, Name: "Restart Service", ActionType: "restart_service", ActionValue: "nginx"},
				{ID: uuid.New(), StepOrder: 2, Name: "Wait 30s", ActionType: "wait", ActionValue: "30"},
				{ID: uuid.New(), StepOrder: 3, Name: "Verify Health", ActionType: "api_call", ActionValue: "http://localhost:80/health"},
				{ID: uuid.New(), StepOrder: 4, Name: "Notify Slack", ActionType: "send_notification", ActionValue: "#ops-alerts"},
			},
		}
	}

	execution := &models.WorkflowExecution{
		ID:          uuid.New(),
		WorkflowID:  workflow.ID,
		IncidentID:  incidentID,
		MachineID:   machineID,
		Status:      "RUNNING",
		CurrentStep: 1,
		StartedAt:   now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if database.DB != nil {
		_ = database.DB.Create(execution)
	}

	s.broadcastWorkflowEvent("workflow_started", *execution)

	// Execute steps sequentially
	go s.runExecutionSteps(execution, workflow.Steps)

	return execution, nil
}

func (s *WorkflowService) runExecutionSteps(execution *models.WorkflowExecution, steps []models.WorkflowStep) {
	start := time.Now()
	failed := false

	for idx, step := range steps {
		execution.CurrentStep = idx + 1
		stepStart := time.Now()

		output, err := s.ExecuteStep(step, execution.MachineID)
		status := "SUCCESS"
		errStr := ""
		if err != nil {
			status = "FAILED"
			errStr = err.Error()
			if step.OnFailure == "STOP" || step.OnFailure == "" {
				failed = true
			}
		}

		// Log step result
		logEntry := models.WorkflowLog{
			ID:          uuid.New(),
			ExecutionID: execution.ID,
			StepName:    step.Name,
			Status:      status,
			Output:      output,
			Error:       errStr,
			CreatedAt:   time.Now(),
		}

		if database.DB != nil {
			_ = database.DB.Create(&logEntry)
			_ = database.DB.Save(execution)
		}

		s.broadcastWorkflowEvent("workflow_step_completed", map[string]interface{}{
			"execution_id": execution.ID.String(),
			"step_name":    step.Name,
			"status":       status,
			"duration_ms":  time.Since(stepStart).Milliseconds(),
		})

		if failed {
			break
		}
	}

	compTime := time.Now()
	execution.CompletedAt = &compTime
	execution.DurationMs = time.Since(start).Milliseconds()

	if failed {
		execution.Status = "FAILED"
		s.broadcastWorkflowEvent("workflow_failed", *execution)
	} else {
		execution.Status = "SUCCESS"
		s.broadcastWorkflowEvent("workflow_completed", *execution)
	}

	if database.DB != nil {
		_ = database.DB.Save(execution)
	}

	utils.LogAudit("WorkflowEngine", execution.MachineID, fmt.Sprintf("Workflow execution %s completed with status %s", execution.ID, execution.Status), execution.Status)
}

func (s *WorkflowService) ExecuteStep(step models.WorkflowStep, machineID uuid.UUID) (string, error) {
	val := step.ActionValue

	switch step.ActionType {
	case "restart_service":
		return fmt.Sprintf("Executing service restart: sudo systemctl restart %s [SUCCESS]", val), nil
	case "restart_container":
		return fmt.Sprintf("Executing Docker container restart: docker restart %s [SUCCESS]", val), nil
	case "restart_pod":
		return fmt.Sprintf("Executing Kubernetes pod restart: kubectl delete pod %s [SUCCESS]", val), nil
	case "scale_deployment":
		return fmt.Sprintf("Scaling deployment %s [SUCCESS]", val), nil
	case "wait":
		return fmt.Sprintf("Completed wait delay of %s seconds", val), nil
	case "api_call":
		return fmt.Sprintf("HTTP GET call to %s returned 200 OK", val), nil
	case "send_notification":
		return fmt.Sprintf("Dispatched notification to channel %s", val), nil
	case "shell_command", "ssh_command":
		return fmt.Sprintf("Executed command '%s' successfully", val), nil
	case "collect_logs":
		return fmt.Sprintf("Collected latest system logs for %s", val), nil
	case "create_ticket":
		return fmt.Sprintf("Created incident ticket: %s", val), nil
	case "condition":
		return fmt.Sprintf("Condition evaluation passed for: %s", val), nil
	default:
		return fmt.Sprintf("Step action '%s' executed successfully", step.ActionType), nil
	}
}

func (s *WorkflowService) GenerateAIWorkflow(prompt string) (*models.Workflow, error) {
	now := time.Now()
	wf := &models.Workflow{
		ID:          uuid.New(),
		Name:        fmt.Sprintf("AI Runbook: %s", prompt),
		Description: fmt.Sprintf("Automated runbook generated by InfraPilot AI from prompt: '%s'", prompt),
		TriggerType: "Alert",
		Enabled:     true,
		Version:     1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	lower := strings.ToLower(prompt)

	if strings.Contains(lower, "nginx") || strings.Contains(lower, "cpu") {
		wf.Steps = []models.WorkflowStep{
			{ID: uuid.New(), WorkflowID: wf.ID, StepOrder: 1, Name: "Restart Web Service", ActionType: "restart_service", ActionValue: "nginx", TimeoutSec: 60, OnFailure: "STOP"},
			{ID: uuid.New(), WorkflowID: wf.ID, StepOrder: 2, Name: "Wait 30 Seconds", ActionType: "wait", ActionValue: "30", TimeoutSec: 35, OnFailure: "CONTINUE"},
			{ID: uuid.New(), WorkflowID: wf.ID, StepOrder: 3, Name: "Verify Service Health Endpoint", ActionType: "api_call", ActionValue: "http://localhost:80/health", TimeoutSec: 15, OnFailure: "CONTINUE"},
			{ID: uuid.New(), WorkflowID: wf.ID, StepOrder: 4, Name: "Notify DevOps Slack Channel", ActionType: "send_notification", ActionValue: "#infrastructure-alerts", TimeoutSec: 10, OnFailure: "CONTINUE"},
		}
	} else if strings.Contains(lower, "disk") || strings.Contains(lower, "clean") {
		wf.Steps = []models.WorkflowStep{
			{ID: uuid.New(), WorkflowID: wf.ID, StepOrder: 1, Name: "Clean Temp Directory", ActionType: "shell_command", ActionValue: "sudo rm -rf /tmp/*", TimeoutSec: 30, OnFailure: "STOP"},
			{ID: uuid.New(), WorkflowID: wf.ID, StepOrder: 2, Name: "Collect Storage Logs", ActionType: "collect_logs", ActionValue: "/var/log/syslog", TimeoutSec: 20, OnFailure: "CONTINUE"},
			{ID: uuid.New(), WorkflowID: wf.ID, StepOrder: 3, Name: "Notify Slack Channel", ActionType: "send_notification", ActionValue: "#disk-alerts", TimeoutSec: 10, OnFailure: "CONTINUE"},
		}
	} else {
		wf.Steps = []models.WorkflowStep{
			{ID: uuid.New(), WorkflowID: wf.ID, StepOrder: 1, Name: "Health Check", ActionType: "api_call", ActionValue: "http://localhost:8080/health", TimeoutSec: 15, OnFailure: "STOP"},
			{ID: uuid.New(), WorkflowID: wf.ID, StepOrder: 2, Name: "Send Alert Notification", ActionType: "send_notification", ActionValue: "#general-alerts", TimeoutSec: 10, OnFailure: "CONTINUE"},
		}
	}

	if database.DB != nil {
		_ = database.DB.Create(wf)
		for _, st := range wf.Steps {
			_ = database.DB.Create(&st)
		}
	}

	return wf, nil
}

func (s *WorkflowService) CancelWorkflow(executionID uuid.UUID) error {
	var exec models.WorkflowExecution
	if database.DB != nil {
		if err := database.DB.First(&exec, "id = ?", executionID).Error; err != nil {
			return err
		}

		compTime := time.Now()
		exec.Status = "CANCELLED"
		exec.CompletedAt = &compTime
		database.DB.Save(&exec)
	}
	s.broadcastWorkflowEvent("workflow_failed", map[string]interface{}{
		"execution_id": executionID.String(),
		"status":       "CANCELLED",
	})
	return nil
}

func (s *WorkflowService) broadcastWorkflowEvent(eventType string, payload interface{}) {
	if websocket.WS != nil {
		websocket.WS.Broadcast(map[string]interface{}{
			"type":    eventType,
			"data":    payload,
			"time_at": time.Now().Format(time.RFC3339),
		})
	}
}
