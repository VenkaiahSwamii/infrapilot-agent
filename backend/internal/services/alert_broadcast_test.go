package services

import (
	"sync"
	"testing"

	"infrapilot/backend/internal/events"

	"github.com/google/uuid"
)

func TestAlertEngine_PublishAlertCreatedEvent(t *testing.T) {
	bus := events.NewEventBus()

	var wg sync.WaitGroup
	var receivedEvent events.AlertCreatedEvent
	var receivedMutex sync.Mutex

	wg.Add(1)
	bus.Subscribe("alert.created", func(e events.Event) {
		if alertEvt, ok := e.(events.AlertCreatedEvent); ok {
			receivedMutex.Lock()
			receivedEvent = alertEvt
			receivedMutex.Unlock()
			wg.Done()
		}
	})

	engine := NewAlertEngine(bus)
	testMachineID := uuid.New().String()

	// Evaluate with breaching CPU=95%
	engine.Evaluate(testMachineID, 95.0, 50.0, 60.0, 30.0, 0.0)

	wg.Wait()

	receivedMutex.Lock()
	defer receivedMutex.Unlock()

	if receivedEvent.MachineID != testMachineID {
		t.Errorf("Expected machine_id %s, got %s", testMachineID, receivedEvent.MachineID)
	}
	if receivedEvent.Severity != "Critical" {
		t.Errorf("Expected severity Critical, got %s", receivedEvent.Severity)
	}
	if receivedEvent.Title != "CPU Threshold Exceeded" {
		t.Errorf("Expected title 'CPU Threshold Exceeded', got '%s'", receivedEvent.Title)
	}
	if receivedEvent.Status != "OPEN" {
		t.Errorf("Expected status OPEN, got %s", receivedEvent.Status)
	}
}
