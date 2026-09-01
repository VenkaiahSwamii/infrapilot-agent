package events

import "time"

// Event is the core interface for all events in the system
type Event interface {
	Name() string
	Timestamp() time.Time
}
