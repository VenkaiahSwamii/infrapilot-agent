package automation

import (
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

type Repository interface {
	GetRules(orgID string) ([]models.AutomationRule, error)
	CreateRule(rule *models.AutomationRule) error
	GetRuleByID(id uuid.UUID) (*models.AutomationRule, error)
	UpdateRule(rule *models.AutomationRule) error
	DeleteRule(id uuid.UUID) error

	GetRunbooks(orgID string) ([]models.Runbook, error)
	CreateRunbook(runbook *models.Runbook) error

	GetHistory(orgID string) ([]models.AutomationHistory, error)
	CreateHistory(history *models.AutomationHistory) error

	GetApprovals(orgID string) ([]models.ApprovalRequest, error)
	CreateApproval(app *models.ApprovalRequest) error
	UpdateApproval(app *models.ApprovalRequest) error
}

type postgresRepository struct{}

func NewRepository() Repository {
	return &postgresRepository{}
}

func (r *postgresRepository) GetRules(orgID string) ([]models.AutomationRule, error) {
	if database.DB == nil {
		return generateMockRules(orgID), nil
	}
	var rules []models.AutomationRule
	db := database.DB.Model(&models.AutomationRule{})
	if orgID != "" && orgID != "default" {
		db = db.Where("organization_id = ?", orgID)
	}
	err := db.Order("created_at desc").Find(&rules).Error
	if err != nil || len(rules) == 0 {
		return generateMockRules(orgID), nil
	}
	return rules, nil
}

func (r *postgresRepository) CreateRule(rule *models.AutomationRule) error {
	if database.DB == nil {
		return nil
	}
	if rule.ID == uuid.Nil {
		rule.ID = uuid.New()
	}
	return database.DB.Create(rule).Error
}

func (r *postgresRepository) GetRuleByID(id uuid.UUID) (*models.AutomationRule, error) {
	if database.DB == nil {
		return &models.AutomationRule{ID: id, Name: "Mock Rule", TriggerEventType: "high_cpu", ActionType: "restart_service"}, nil
	}
	var rule models.AutomationRule
	if err := database.DB.First(&rule, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *postgresRepository) UpdateRule(rule *models.AutomationRule) error {
	if database.DB == nil {
		return nil
	}
	return database.DB.Save(rule).Error
}

func (r *postgresRepository) DeleteRule(id uuid.UUID) error {
	if database.DB == nil {
		return nil
	}
	return database.DB.Delete(&models.AutomationRule{}, "id = ?", id).Error
}

func (r *postgresRepository) GetRunbooks(orgID string) ([]models.Runbook, error) {
	if database.DB == nil {
		return generateMockRunbooks(orgID), nil
	}
	var runbooks []models.Runbook
	db := database.DB.Model(&models.Runbook{})
	if orgID != "" && orgID != "default" {
		db = db.Where("organization_id = ?", orgID)
	}
	err := db.Order("created_at desc").Find(&runbooks).Error
	if err != nil || len(runbooks) == 0 {
		return generateMockRunbooks(orgID), nil
	}
	return runbooks, nil
}

func (r *postgresRepository) CreateRunbook(runbook *models.Runbook) error {
	if database.DB == nil {
		return nil
	}
	if runbook.ID == uuid.Nil {
		runbook.ID = uuid.New()
	}
	return database.DB.Create(runbook).Error
}

func (r *postgresRepository) GetHistory(orgID string) ([]models.AutomationHistory, error) {
	if database.DB == nil {
		return generateMockHistory(orgID), nil
	}
	var history []models.AutomationHistory
	db := database.DB.Model(&models.AutomationHistory{})
	if orgID != "" && orgID != "default" {
		db = db.Where("organization_id = ?", orgID)
	}
	err := db.Order("start_time desc").Find(&history).Error
	if err != nil || len(history) == 0 {
		return generateMockHistory(orgID), nil
	}
	return history, nil
}

func (r *postgresRepository) CreateHistory(history *models.AutomationHistory) error {
	if database.DB == nil {
		return nil
	}
	if history.ID == uuid.Nil {
		history.ID = uuid.New()
	}
	return database.DB.Create(history).Error
}

func (r *postgresRepository) GetApprovals(orgID string) ([]models.ApprovalRequest, error) {
	if database.DB == nil {
		return generateMockApprovals(orgID), nil
	}
	var approvals []models.ApprovalRequest
	db := database.DB.Model(&models.ApprovalRequest{})
	if orgID != "" && orgID != "default" {
		db = db.Where("organization_id = ?", orgID)
	}
	err := db.Order("created_at desc").Find(&approvals).Error
	if err != nil || len(approvals) == 0 {
		return generateMockApprovals(orgID), nil
	}
	return approvals, nil
}

func (r *postgresRepository) CreateApproval(app *models.ApprovalRequest) error {
	if database.DB == nil {
		return nil
	}
	if app.ID == uuid.Nil {
		app.ID = uuid.New()
	}
	return database.DB.Create(app).Error
}

func (r *postgresRepository) UpdateApproval(app *models.ApprovalRequest) error {
	if database.DB == nil {
		return nil
	}
	return database.DB.Save(app).Error
}

// Fallback Mock Data Generators
func generateMockRules(orgID string) []models.AutomationRule {
	now := time.Now()
	return []models.AutomationRule{
		{
			ID:               uuid.New(),
			OrganizationID:   orgID,
			Name:             "Auto-Restart Nginx on High CPU",
			Description:      "Restart web server daemon when CPU exceeds 95% for 10 minutes",
			TriggerEventType: "high_cpu",
			ConditionJSON:    `{"metric":"cpu","operator":">","threshold":95,"duration_min":10}`,
			ActionType:       "restart_service",
			TargetResource:   "nginx.service",
			RequiresApproval: false,
			Enabled:          true,
			CreatedAt:        now,
		},
		{
			ID:               uuid.New(),
			OrganizationID:   orgID,
			Name:             "Auto-Prune Temp Files on Low Disk",
			Description:      "Purge /tmp and journal logs when disk usage exceeds 90%",
			TriggerEventType: "low_disk",
			ConditionJSON:    `{"metric":"disk","operator":">","threshold":90}`,
			ActionType:       "clean_disk",
			TargetResource:   "/tmp, /var/log/journal",
			RequiresApproval: false,
			Enabled:          true,
			CreatedAt:        now,
		},
		{
			ID:               uuid.New(),
			OrganizationID:   orgID,
			Name:             "K8s Pod CrashLoopBackOff Auto-Rollout",
			Description:      "Trigger deployment rollout restart on CrashLoopBackOff status",
			TriggerEventType: "pod_crash",
			ConditionJSON:    `{"status":"CrashLoopBackOff"}`,
			ActionType:       "restart_pod",
			TargetResource:   "api-deployment",
			RequiresApproval: true,
			Enabled:          true,
			CreatedAt:        now,
		},
		{
			ID:               uuid.New(),
			OrganizationID:   orgID,
			Name:             "Auto-Restart Stopped Docker Containers",
			Description:      "Automatically restart critical stopped docker container instances",
			TriggerEventType: "container_down",
			ConditionJSON:    `{"state":"exited"}`,
			ActionType:       "restart_container",
			TargetResource:   "redis-prod-01",
			RequiresApproval: false,
			Enabled:          true,
			CreatedAt:        now,
		},
	}
}

func generateMockRunbooks(orgID string) []models.Runbook {
	now := time.Now()
	return []models.Runbook{
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			Name:           "Restart Nginx Web Server",
			Description:    "Executes systemctl restart nginx via SSH",
			Category:       "service",
			Command:        "sudo systemctl restart nginx",
			TargetOs:       "linux",
			TimeoutSec:     30,
			CreatedAt:      now,
		},
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			Name:           "Restart Apache HTTPD",
			Description:    "Executes systemctl restart apache2 via SSH",
			Category:       "service",
			Command:        "sudo systemctl restart apache2",
			TargetOs:       "linux",
			TimeoutSec:     30,
			CreatedAt:      now,
		},
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			Name:           "Restart Docker Daemon & Containers",
			Description:    "Restarts docker service and recovers dead containers",
			Category:       "docker",
			Command:        "sudo systemctl restart docker && docker restart $(docker ps -a -q -f status=exited)",
			TargetOs:       "linux",
			TimeoutSec:     60,
			CreatedAt:      now,
		},
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			Name:           "Restart Kubernetes Deployment",
			Description:    "Triggers kubectl rollout restart for target deployment",
			Category:       "kubernetes",
			Command:        "kubectl rollout restart deployment/{{target}}",
			TargetOs:       "linux",
			TimeoutSec:     120,
			CreatedAt:      now,
		},
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			Name:           "Clear System Cache & Buffers",
			Description:    "Flushes Linux page cache, dentries, and inodes",
			Category:       "maintenance",
			Command:        "sudo sync && echo 3 | sudo tee /proc/sys/vm/drop_caches",
			TargetOs:       "linux",
			TimeoutSec:     15,
			CreatedAt:      now,
		},
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			Name:           "Rotate System Logs",
			Description:    "Triggers logrotate for system journals",
			Category:       "maintenance",
			Command:        "sudo logrotate -f /etc/logrotate.conf",
			TargetOs:       "linux",
			TimeoutSec:     30,
			CreatedAt:      now,
		},
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			Name:           "Cleanup Disk & Temp Files",
			Description:    "Purges apt cache, journal logs older than 3 days, and /tmp",
			Category:       "maintenance",
			Command:        "sudo apt-get clean && sudo journalctl --vacuum-time=3d",
			TargetOs:       "linux",
			TimeoutSec:     45,
			CreatedAt:      now,
		},
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			Name:           "Restart PostgreSQL Database",
			Description:    "Gracefully restarts PostgreSQL database cluster",
			Category:       "database",
			Command:        "sudo systemctl restart postgresql",
			TargetOs:       "linux",
			TimeoutSec:     60,
			CreatedAt:      now,
		},
	}
}

func generateMockHistory(orgID string) []models.AutomationHistory {
	now := time.Now()
	return []models.AutomationHistory{
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			RuleName:       "Auto-Restart Nginx on High CPU",
			Trigger:        "High CPU Alert (98%)",
			User:           "AIOps Engine",
			Action:         "Restart Nginx Service",
			Target:         "api-prod-node-01",
			StartTime:      now.Add(-2 * time.Minute),
			EndTime:        now.Add(-2*time.Minute + 320*time.Millisecond),
			Result:         "Success",
			OutputSnippet:  "systemctl restart nginx executed cleanly. CPU dropped to 34%.",
		},
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			RuleName:       "Auto-Prune Temp Files on Low Disk",
			Trigger:        "Low Disk Space Alert (92%)",
			User:           "AIOps Engine",
			Action:         "Cleanup Disk",
			Target:         "db-primary-postgres",
			StartTime:      now.Add(-45 * time.Minute),
			EndTime:        now.Add(-45*time.Minute + 1200*time.Millisecond),
			Result:         "Success",
			OutputSnippet:  "Freed 14.8GB from journalctl vacuuming and /tmp cleanup.",
		},
	}
}

func generateMockApprovals(orgID string) []models.ApprovalRequest {
	now := time.Now()
	return []models.ApprovalRequest{
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			RunID:          uuid.New(),
			RuleName:       "K8s Pod CrashLoopBackOff Auto-Rollout",
			Action:         "Rollout Restart Deployment 'api-deployment'",
			TargetResource: "k8s-cluster-prod / api-deployment",
			RequestedBy:    "AIOps Engine",
			Status:         "PENDING",
			CreatedAt:      now.Add(-10 * time.Minute),
		},
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			RunID:          uuid.New(),
			RuleName:       "Restart Database Node",
			Action:         "Restart PostgreSQL Cluster Service",
			TargetResource: "db-primary-postgres",
			RequestedBy:    "AIOps Engine",
			Status:         "PENDING",
			CreatedAt:      now.Add(-25 * time.Minute),
		},
	}
}
