package events

import "time"

// MachineRegisteredEvent is published when a new machine is registered
type MachineRegisteredEvent struct {
	MachineID string
	Hostname  string
	IPAddress string
	OS        string
	Platform  string
	APIKey    string
	Time      time.Time
}

func (e MachineRegisteredEvent) Name() string {
	return "machine.registered"
}

func (e MachineRegisteredEvent) Timestamp() time.Time {
	return e.Time
}

// MachineOnlineEvent is published when a machine sends a heartbeat
type MachineOnlineEvent struct {
	MachineID string
	Hostname  string
	IPAddress string
	OS        string
	Platform  string
	LastSeen  time.Time
	Time      time.Time
}

func (e MachineOnlineEvent) Name() string {
	return "machine.online"
}

func (e MachineOnlineEvent) Timestamp() time.Time {
	return e.Time
}

// MachineOfflineEvent is published when a machine goes offline
type MachineOfflineEvent struct {
	MachineID string
	Hostname  string
	IPAddress string
	OS        string
	Platform  string
	LastSeen  time.Time
	Time      time.Time
}

func (e MachineOfflineEvent) Name() string {
	return "machine.offline"
}

func (e MachineOfflineEvent) Timestamp() time.Time {
	return e.Time
}
