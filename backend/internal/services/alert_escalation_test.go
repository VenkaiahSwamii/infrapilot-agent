package services

import (
	"testing"
)

func TestAlertEscalation_Escalate(t *testing.T) {
	escalation := NewAlertEscalation()

	if escalation == nil {
		t.Fatalf("Expected non-nil AlertEscalation instance")
	}

	// Test escalation dispatch across Email, Slack, Teams
	escalation.Escalate("Critical CPU Usage: 97.5%")
	escalation.Escalate("Critical Memory Usage: 96.0%")
	escalation.Escalate("Critical Disk Usage: 98.2%")
}
