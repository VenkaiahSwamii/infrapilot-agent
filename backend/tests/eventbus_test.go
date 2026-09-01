package tests

import (
	"sync"
	"testing"
	"time"

	"infrapilot/backend/internal/events"
)

// testEvent implements the events.Event interface for testing
type testEvent struct {
	name      string
	timestamp time.Time
	payload   string
}

func (e testEvent) Name() string         { return e.name }
func (e testEvent) Timestamp() time.Time { return e.timestamp }

func TestEventBus(t *testing.T) {
	t.Run("publish and subscribe", func(t *testing.T) {
		bus := events.NewEventBus()
		received := make(chan string, 1)

		bus.Subscribe("test.event", func(e events.Event) {
			if te, ok := e.(testEvent); ok {
				received <- te.payload
			}
		})

		bus.Publish(testEvent{
			name:      "test.event",
			timestamp: time.Now(),
			payload:   "hello world",
		})

		select {
		case msg := <-received:
			if msg != "hello world" {
				t.Errorf("Expected 'hello world', got '%s'", msg)
			}
		case <-time.After(time.Second):
			t.Error("Timed out waiting for event")
		}
	})

	t.Run("multiple subscribers", func(t *testing.T) {
		bus := events.NewEventBus()
		var mu sync.Mutex
		count := 0

		for i := 0; i < 3; i++ {
			bus.Subscribe("test.multi", func(e events.Event) {
				mu.Lock()
				count++
				mu.Unlock()
			})
		}

		bus.Publish(testEvent{
			name:      "test.multi",
			timestamp: time.Now(),
		})

		time.Sleep(200 * time.Millisecond)

		mu.Lock()
		if count != 3 {
			t.Errorf("Expected 3 subscribers to fire, got %d", count)
		}
		mu.Unlock()
	})

	t.Run("different event types don't cross", func(t *testing.T) {
		bus := events.NewEventBus()
		received := make(chan string, 1)

		bus.Subscribe("event.a", func(e events.Event) {
			received <- "a"
		})

		bus.Publish(testEvent{
			name:      "event.b",
			timestamp: time.Now(),
		})

		select {
		case <-received:
			t.Error("Should not have received event for different type")
		case <-time.After(200 * time.Millisecond):
			// Expected - no event received
		}
	})

	t.Run("no handlers registered", func(t *testing.T) {
		bus := events.NewEventBus()

		// Should not panic
		bus.Publish(testEvent{
			name:      "nonexistent.event",
			timestamp: time.Now(),
		})

		// If we got here without panic, test passes
	})
}
