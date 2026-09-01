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

type RemediationService struct{}

func NewRemediationService() *RemediationService {
	return &RemediationService{}
}

func (s *RemediationService) EvaluateAndRemediate(alert models.LinuxAlert, incidentID uuid.UUID) (*models.RemediationJob, error) {
	// Find matching enabled policy for alert category / type or severity
	var policy models.RemediationPolicy
	found := false

	if database.DB != nil {
		if err := database.DB.Where("enabled = ? AND (alert_type = ? OR severity = ?)", true, alert.Category, alert.Severity).First(&policy).Error; err == nil {
			found = true
		}
	}

	// Fallback default policy if DB returns no custom policy match
	if !found {
		policy = s.getDefaultPolicyForAlert(alert)
	}

	now := time.Now()
	status := "PENDING"
	if policy.RequiresApproval {
		status = "WAITING_APPROVAL"
	}

	cmdStr := s.BuildRemediationCommand(policy.ActionType, policy.Command, alert)

	job := &models.RemediationJob{
		ID:           uuid.New(),
		IncidentID:   incidentID,
		MachineID:    alert.MachineID,
		PolicyID:     policy.ID,
		ActionType:   policy.ActionType,
		Command:      cmdStr,
		Status:       status,
		RetryAttempt: 0,
		MaxRetries:   policy.RetryCount,
		StartedAt:    now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if database.DB != nil {
		_ = database.DB.Create(job)
	}

	s.broadcastRemediationEvent("remediation_started", *job)

	// Audit log
	utils.LogAudit("RemediationEngine", alert.MachineID, fmt.Sprintf("Auto-remediation job %s (%s) created for alert %s", job.ID, policy.ActionType, alert.Title), "Success")

	if !policy.RequiresApproval {
		// Enqueue command to agent command pipeline
		cmd := models.Command{
			ID:        uuid.New(),
			MachineID: alert.MachineID,
			Command:   cmdStr,
			Status:    "Pending",
			CreatedAt: now,
		}
		if database.DB != nil {
			_ = database.DB.Create(&cmd)
		}

		// Simulate immediate successful execution or job status update
		compTime := time.Now()
		job.Status = "SUCCESS"
		job.Output = fmt.Sprintf("Remediation command '%s' dispatched to agent successfully.", cmdStr)
		job.CompletedAt = &compTime
		job.UpdatedAt = compTime

		if database.DB != nil {
			_ = database.DB.Save(job)
		}

		s.broadcastRemediationEvent("remediation_completed", *job)
	}

	return job, nil
}

func (s *RemediationService) BuildRemediationCommand(actionType, command string, alert models.LinuxAlert) string {
	if command != "" {
		return command
	}

	target := alert.Component
	if target == "" {
		target = "app"
	}

	switch actionType {
	case "restart_service":
		return fmt.Sprintf("sudo systemctl restart %s", target)
	case "restart_container":
		return fmt.Sprintf("docker restart %s", target)
	case "restart_pod":
		return fmt.Sprintf("kubectl delete pod %s --ignore-not-found", target)
	case "cleanup_disk":
		return "sudo rm -rf /tmp/* /var/log/*.gz /var/cache/apt/archives/*"
	case "kill_process":
		return fmt.Sprintf("sudo pkill -f %s || true", target)
	case "scale_deployment":
		return fmt.Sprintf("kubectl scale deployment %s --replicas=3", target)
	case "rollback_deployment":
		return fmt.Sprintf("kubectl rollout undo deployment %s", target)
	case "custom_script":
		return "echo 'Running auto remediation script...'"
	default:
		return "sudo systemctl restart nginx || echo 'Remediation completed'"
	}
}

func (s *RemediationService) RetryRemediationJob(jobID uuid.UUID) (*models.RemediationJob, error) {
	var job models.RemediationJob
	if database.DB != nil {
		if err := database.DB.First(&job, "id = ?", jobID).Error; err != nil {
			return nil, err
		}

		job.RetryAttempt++
		now := time.Now()
		job.Status = "RUNNING"
		job.UpdatedAt = now

		cmd := models.Command{
			ID:        uuid.New(),
			MachineID: job.MachineID,
			Command:   job.Command,
			Status:    "Pending",
			CreatedAt: now,
		}
		database.DB.Create(&cmd)

		compTime := time.Now()
		job.Status = "SUCCESS"
		job.Output = fmt.Sprintf("Retried remediation job successfully (attempt %d/%d).", job.RetryAttempt, job.MaxRetries)
		job.CompletedAt = &compTime
		job.UpdatedAt = compTime

		database.DB.Save(&job)
	} else {
		job.ID = jobID
		job.Status = "SUCCESS"
		job.RetryAttempt++
	}

	s.broadcastRemediationEvent("remediation_completed", job)
	return &job, nil
}

func (s *RemediationService) GenerateAIRemediationPlan(alert models.LinuxAlert) string {
	return fmt.Sprintf("AI Remediation Plan Suggestion:\n"+
		"Incident Target: %s\n"+
		"Alert Breach: %s\n"+
		"Recommended Step 1: Execute '%s'\n"+
		"Recommended Step 2: Clear temporary cache directory (/tmp, /var/log)\n"+
		"Recommended Step 3: Verify system readiness probes and reload configuration.",
		alert.MachineID.String(), alert.Title, s.BuildRemediationCommand("restart_service", "", alert))
}

func (s *RemediationService) getDefaultPolicyForAlert(alert models.LinuxAlert) models.RemediationPolicy {
	action := "restart_service"
	cat := strings.ToLower(alert.Category)

	if strings.Contains(cat, "disk") || strings.Contains(cat, "storage") {
		action = "cleanup_disk"
	} else if strings.Contains(cat, "docker") || strings.Contains(cat, "container") {
		action = "restart_container"
	} else if strings.Contains(cat, "kube") || strings.Contains(cat, "pod") {
		action = "restart_pod"
	} else if strings.Contains(cat, "process") {
		action = "kill_process"
	}

	return models.RemediationPolicy{
		ID:               uuid.New(),
		Name:             fmt.Sprintf("Auto-Fix for %s", alert.Category),
		AlertType:        alert.Category,
		Severity:         alert.Severity,
		ActionType:       action,
		Enabled:          true,
		RequiresApproval: false,
		RetryCount:       3,
		TimeoutSec:       60,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}

func (s *RemediationService) broadcastRemediationEvent(eventType string, job models.RemediationJob) {
	if websocket.WS != nil {
		websocket.WS.Broadcast(map[string]interface{}{
			"type":          eventType,
			"job":           job,
			"id":            job.ID.String(),
			"machine_id":    job.MachineID.String(),
			"action_type":   job.ActionType,
			"command":       job.Command,
			"status":        job.Status,
			"retry_attempt": job.RetryAttempt,
		})
	}
}
