package services

import (
	"testing"
	"time"

	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

func TestNotificationService_MultiChannelSend(t *testing.T) {
	svc := NewNotificationService()

	criticalAlert := models.LinuxAlert{
		ID:        uuid.New(),
		Title:     "CPU Threshold Exceeded",
		Severity:  "Critical",
		Message:   "CPU usage is 95.0%",
		Status:    "OPEN",
		CreatedAt: time.Now(),
	}

	if err := svc.Send(criticalAlert); err != nil {
		t.Errorf("Expected nil error for Critical notification send, got %v", err)
	}

	warningAlert := models.LinuxAlert{
		ID:        uuid.New(),
		Title:     "Swap Usage High",
		Severity:  "Warning",
		Message:   "Swap usage is 85.0%",
		Status:    "OPEN",
		CreatedAt: time.Now(),
	}

	if err := svc.Send(warningAlert); err != nil {
		t.Errorf("Expected nil error for Warning notification send, got %v", err)
	}
}
