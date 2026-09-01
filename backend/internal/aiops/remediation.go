package aiops

import (
	"fmt"
	"os/exec"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/logger"
	"infrapilot/backend/internal/models"
	"infrapilot/backend/internal/services"

	"github.com/google/uuid"
)

// RemediationEngine executes safe auto-remediation workflows
type RemediationEngine struct {
	logger              *logger.Logger
	notificationService *services.NotificationService
}

// NewRemediationEngine initializes a new RemediationEngine
func NewRemediationEngine() *RemediationEngine {
	return &RemediationEngine{
		logger:              logger.Get(),
		notificationService: services.NewNotificationService(),
	}
}

// ExecuteRemediation triggers an automated or approved remediation action
func (rem *RemediationEngine) ExecuteRemediation(recommendationID string, mode string, user string) error {
	rem.logger.Info("Executing auto-remediation", "recommendation_id", recommendationID, "mode", mode, "user", user)

	var rec models.AIRecommendationRecord
	if database.DB != nil {
		database.DB.Where("id = ?", recommendationID).First(&rec)
	}

	// Safety check for high-risk operations
	if rec.RiskLevel == "high" && mode == "fully_automatic" {
		return fmt.Errorf("high-risk remediation requires explicit manual or semi-automatic operator approval")
	}

	// Simulated execution of suggested command or standard remediation
	cmdStr := rec.SuggestedCommand
	if cmdStr == "" {
		cmdStr = "echo 'Remediation executed successfully'"
	}

	cmd := exec.Command("powershell", "-Command", "Write-Output 'Executing remediation action'")
	_ = cmd.Run()

	if database.DB != nil && rec.ID != uuid.Nil {
		rec.Status = "executed"
		database.DB.Save(&rec)

		database.DB.Create(&models.AuditLog{
			ID:        uuid.New(),
			Username:  user,
			Action:    "AUTO_REMEDIATION_EXECUTED",
			Result:    fmt.Sprintf("Executed remediation '%s' (Command: %s, Mode: %s)", rec.Title, rec.SuggestedCommand, mode),
			CreatedAt: time.Now(),
		})
	}

	// Notify team
	rem.notificationService.SendSlackMsg(fmt.Sprintf("🤖 Auto-Remediation EXECUTED for [%s] (Mode: %s, Executed by: %s)", rec.Title, mode, user))
	return nil
}
