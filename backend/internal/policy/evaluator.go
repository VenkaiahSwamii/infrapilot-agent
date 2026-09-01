package policy

import (
	"strings"

	"infrapilot/backend/internal/models"
)

type ABACContext struct {
	UserRole     string
	UserOrgID    string
	MFAVerified  bool
	ResourceType string
	Action       string
	TargetOrgID  string
	RiskLevel    string
	ClientIP     string
}

// EvaluateABAC checks if the requested action is permitted under ABAC rules
func EvaluateABAC(ctx ABACContext, rules []models.PolicyRuleRecord) (bool, string) {
	// 1. Cross-tenant isolation check
	if ctx.UserRole != "superadmin" && ctx.UserRole != "super_admin" {
		if ctx.UserOrgID != "" && ctx.TargetOrgID != "" && ctx.UserOrgID != ctx.TargetOrgID && ctx.TargetOrgID != "default" {
			return false, "ABAC Deny: Cross-organization access forbidden"
		}
	}

	// 2. Evaluate specific rules
	for _, r := range rules {
		if !r.Enabled {
			continue
		}

		if strings.EqualFold(r.RolePattern, ctx.UserRole) && strings.EqualFold(r.ResourceType, ctx.ResourceType) {
			if r.Effect == "DENY" {
				return false, "ABAC Deny: Explicit deny rule matched (" + r.RuleName + ")"
			}
		}
	}

	return true, "ABAC Allow: Access granted"
}
