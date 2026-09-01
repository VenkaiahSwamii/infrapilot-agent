package services

import (
	"testing"
	"time"

	"infrapilot/backend/internal/models"
)

func TestEvaluateMachineStatus(t *testing.T) {
	now := time.Now()

	// 1. ONLINE Threshold (<= 60s)
	lastSeenOnline := now.Add(-30 * time.Second)
	statusOnline := EvaluateMachineStatus(lastSeenOnline, now)
	if statusOnline != "ONLINE" {
		t.Errorf("Expected ONLINE status for 30s difference, got '%s'", statusOnline)
	}

	// 2. OFFLINE Threshold (> 60s)
	lastSeenOffline := now.Add(-90 * time.Second)
	statusOffline := EvaluateMachineStatus(lastSeenOffline, now)
	if statusOffline != "OFFLINE" {
		t.Errorf("Expected OFFLINE status for 90s difference, got '%s'", statusOffline)
	}
}

func TestMachineStatusReconcilesFromLastSeen(t *testing.T) {
	now := time.Now().UTC()
	machine := &models.Machine{
		Status:    "OFFLINE",
		Online:    false,
		LastSeen:  now.Add(-30 * time.Second),
		Hostname:  "machine-b",
	}

	status := ReconcileMachineStatus(machine, now, time.Minute)
	if status != "ONLINE" {
		t.Fatalf("expected ONLINE based on recent last_seen, got %q", status)
	}
	if machine.Status != "ONLINE" {
		t.Fatalf("expected machine status to be corrected to ONLINE, got %q", machine.Status)
	}
	if !machine.Online {
		t.Fatal("expected machine to be marked online when last_seen is recent")
	}
}
