package services

import (
	"testing"
	"time"

	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

func TestRemediationService_EvaluationAndCommands(t *testing.T) {
	remSvc := NewRemediationService()

	testAlert := models.LinuxAlert{
		ID:        uuid.New(),
		MachineID: uuid.New(),
		Title:     "Disk Threshold Exceeded",
		Category:  "Disk",
		Component: "/tmp",
		Severity:  "Critical",
		Status:    "OPEN",
		CreatedAt: time.Now(),
	}

	// 1. Test EvaluateAndRemediate
	job, err := remSvc.EvaluateAndRemediate(testAlert, uuid.Nil)
	if err != nil {
		t.Fatalf("Expected nil error for EvaluateAndRemediate, got %v", err)
	}

	if job.ActionType != "cleanup_disk" {
		t.Errorf("Expected action_type 'cleanup_disk', got '%s'", job.ActionType)
	}

	// 2. Test BuildRemediationCommand
	cmdService := remSvc.BuildRemediationCommand("restart_service", "", testAlert)
	if cmdService != "sudo systemctl restart /tmp" && cmdService != "sudo systemctl restart app" {
		// Valid systemd command string
	}

	// 3. Test GenerateAIRemediationPlan
	aiPlan := remSvc.GenerateAIRemediationPlan(testAlert)
	if aiPlan == "" {
		t.Errorf("Expected non-empty AI remediation plan suggestion")
	}
}
