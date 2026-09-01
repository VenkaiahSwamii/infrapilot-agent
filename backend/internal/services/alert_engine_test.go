package services

import (
	"testing"

	"github.com/google/uuid"
)

func TestAlertEngine_Evaluate(t *testing.T) {
	engine := NewAlertEngine()

	testMachineID := uuid.New().String()

	// 1. Normal metrics (below all thresholds)
	engine.Evaluate(testMachineID, 45.0, 50.0, 60.0, 30.0, 0.0)

	// 2. High metrics (CPU=95, Memory=92)
	engine.Evaluate(testMachineID, 95.0, 92.0, 60.0, 350.0, 6.0)

	// If no panic or error occurs, the test passes safely
}
