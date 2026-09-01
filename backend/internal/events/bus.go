package events

import "log"

// Handler is a function that processes events
type Handler func(Event)

// EventBus manages event subscriptions and publishing
type EventBus struct {
	handlers map[string][]Handler
}

// NewEventBus creates a new event bus instance
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]Handler),
	}
}

// Subscribe registers a handler for a specific event type
func (b *EventBus) Subscribe(eventName string, h Handler) {
	b.handlers[eventName] = append(b.handlers[eventName], h)
	log.Printf("Subscribed handler to event: %s", eventName)
}

// Publish sends an event to all registered handlers
func (b *EventBus) Publish(e Event) {
	eventName := e.Name()
	if handlers, ok := b.handlers[eventName]; ok {
		log.Printf("Publishing event: %s at %v", eventName, e.Timestamp())
		for _, h := range handlers {
			go h(e) // Process handlers asynchronously
		}
	} else {
		log.Printf("No handlers registered for event: %s", eventName)
	}
}
