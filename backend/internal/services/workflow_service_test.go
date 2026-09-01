package services

import (
	"testing"

	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

func TestWorkflowService_ExecutionAndAIGenerator(t *testing.T) {
	wfSvc := NewWorkflowService()

	// 1. Test ExecuteStep for various actions
	step1 := models.WorkflowStep{
		Name:        "Restart Service Step",
		ActionType:  "restart_service",
		ActionValue: "nginx",
	}
	out1, err1 := wfSvc.ExecuteStep(step1, uuid.New())
	if err1 != nil {
		t.Errorf("Expected nil error for ExecuteStep, got %v", err1)
	}
	if out1 == "" {
		t.Errorf("Expected non-empty output for ExecuteStep")
	}

	step2 := models.WorkflowStep{
		Name:        "Wait Step",
		ActionType:  "wait",
		ActionValue: "30",
	}
	out2, err2 := wfSvc.ExecuteStep(step2, uuid.New())
	if err2 != nil {
		t.Errorf("Expected nil error for wait step, got %v", err2)
	}
	if out2 == "" {
		t.Errorf("Expected non-empty output for wait step")
	}

	// 2. Test GenerateAIWorkflow
	aiWf, errAI := wfSvc.GenerateAIWorkflow("When CPU exceeds 90%, restart nginx and notify Slack")
	if errAI != nil {
		t.Fatalf("Expected nil error for GenerateAIWorkflow, got %v", errAI)
	}
	if len(aiWf.Steps) == 0 {
		t.Errorf("Expected AI workflow to generate steps, got 0")
	}

	// 3. Test StartWorkflow
	exec, errExec := wfSvc.StartWorkflow(aiWf.ID, uuid.Nil, uuid.New())
	if errExec != nil {
		t.Errorf("Expected nil error for StartWorkflow, got %v", errExec)
	}
	if exec.Status != "RUNNING" && exec.Status != "SUCCESS" {
		t.Errorf("Expected execution status RUNNING/SUCCESS, got '%s'", exec.Status)
	}
}
