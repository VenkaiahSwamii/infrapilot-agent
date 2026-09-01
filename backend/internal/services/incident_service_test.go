package services

import (
	"strings"
	"testing"
	"time"

	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

func TestCorrelationEngine_And_IncidentService(t *testing.T) {
	corr := NewCorrelationEngine()
	incSvc := NewIncidentService()

	testMachineID := uuid.New()

	alert1 := models.LinuxAlert{
		ID:        uuid.New(),
		MachineID: testMachineID,
		Title:     "CPU Threshold Exceeded",
		Category:  "CPU",
		Severity:  "Critical",
		Status:    "OPEN",
		CreatedAt: time.Now(),
	}

	// 1. Correlate first alert -> creates new incident
	inc1, err := corr.Correlate(alert1)
	if err != nil {
		t.Fatalf("Expected no error on first alert correlation, got %v", err)
	}

	if inc1.MachineID != testMachineID {
		t.Errorf("Expected machine_id %s, got %s", testMachineID, inc1.MachineID)
	}

	// 2. Correlate second alert on same machine -> groups into existing incident
	alert2 := models.LinuxAlert{
		ID:        uuid.New(),
		MachineID: testMachineID,
		Title:     "Memory Threshold Exceeded",
		Category:  "Memory",
		Severity:  "Critical",
		Status:    "OPEN",
		CreatedAt: time.Now(),
	}

	inc2, err := corr.Correlate(alert2)
	if err != nil {
		t.Fatalf("Expected no error on second alert correlation, got %v", err)
	}

	if inc2.ID != inc1.ID {
		t.Errorf("Expected alert2 to correlate into existing incident %s, got %s", inc1.ID, inc2.ID)
	}

	if inc2.AlertCount != 2 {
		t.Errorf("Expected incident alert_count to be 2, got %d", inc2.AlertCount)
	}

	// 3. Root Cause Classification
	rc := incSvc.ClassifyRootCause([]string{"CPU", "Memory"})
	if !strings.Contains(rc, "Resource Exhaustion") {
		t.Errorf("Expected root cause to contain 'Resource Exhaustion', got '%s'", rc)
	}

	// 4. AI Summary Generation
	summary := incSvc.GenerateAIIncidentSummary(*inc2)
	if summary == "" {
		t.Errorf("Expected non-empty AI incident summary")
	}
}
