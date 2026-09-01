package policy

import (
	"infrapilot/backend/internal/logger"
	"infrapilot/backend/internal/models"
)

type PolicyEngine struct {
	logger *logger.Logger
}

func NewPolicyEngine() *PolicyEngine {
	return &PolicyEngine{
		logger: logger.Get(),
	}
}

// VerifyRequest performs continuous Zero Trust evaluation
func (p *PolicyEngine) VerifyRequest(abacCtx ABACContext) (bool, string) {
	rules := GetDefaultRules(abacCtx.UserOrgID)
	allowed, reason := EvaluateABAC(abacCtx, rules)

	p.logger.Info("Zero Trust Policy Evaluation", "user_role", abacCtx.UserRole, "action", abacCtx.Action, "allowed", allowed, "reason", reason)
	return allowed, reason
}

// GetRules lists all active policy rules
func (p *PolicyEngine) GetRules(orgID string) []models.PolicyRuleRecord {
	return GetDefaultRules(orgID)
}
