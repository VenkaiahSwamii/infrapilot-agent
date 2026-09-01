package policy

import (
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

// DefaultRules returns base ABAC security policy rules
func GetDefaultRules(orgID string) []models.PolicyRuleRecord {
	var rules []models.PolicyRuleRecord

	if database.DB != nil {
		database.DB.Where("organization_id = ? OR organization_id = ''", orgID).Find(&rules)
	}

	if len(rules) == 0 {
		now := time.Now()
		rules = []models.PolicyRuleRecord{
			{
				ID:             uuid.New(),
				OrganizationID: orgID,
				RuleName:       "Operator Service Control Rule",
				Effect:         "ALLOW",
				RolePattern:    "operator",
				ResourceType:   "linux_server",
				ActionPattern:  "restart_service",
				ConditionsJSON: `{"same_organization": true, "environment": ["production", "staging"]}`,
				Enabled:        true,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			{
				ID:             uuid.New(),
				OrganizationID: orgID,
				RuleName:       "High-Risk Remote Command Approval",
				Effect:         "DENY",
				RolePattern:    "viewer",
				ResourceType:   "remote_command",
				ActionPattern:  "execute_shell",
				ConditionsJSON: `{"risk_level": "high"}`,
				Enabled:        true,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			{
				ID:             uuid.New(),
				OrganizationID: orgID,
				RuleName:       "Enforce MFA for Admin Actions",
				Effect:         "ALLOW",
				RolePattern:    "admin",
				ResourceType:   "backup_restore",
				ActionPattern:  "trigger_restore",
				ConditionsJSON: `{"mfa_verified": true}`,
				Enabled:        true,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		}
	}

	return rules
}
