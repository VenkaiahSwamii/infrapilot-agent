package events

import "time"

type AlertCreatedEvent struct {
	ID        string    `json:"id"`
	MachineID string    `json:"machine_id"`
	Severity  string    `json:"severity"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (e AlertCreatedEvent) Name() string {
	return "alert.created"
}

func (e AlertCreatedEvent) Timestamp() time.Time {
	return e.CreatedAt
}
