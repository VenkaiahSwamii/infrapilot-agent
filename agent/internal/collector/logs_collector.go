package collector

import (
	"os"
	"runtime"
	"time"

	"github.com/google/uuid"
)

type AgentLogEntry struct {
	ID        uuid.UUID `json:"id"`
	MachineID uuid.UUID `json:"machine_id"`
	Hostname  string    `json:"hostname"`
	Platform  string    `json:"platform"`
	Level     string    `json:"level"`
	Source    string    `json:"source"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type LogsCollector struct {
	lastProcessed time.Time
}

func NewLogsCollector() *LogsCollector {
	return &LogsCollector{lastProcessed: time.Now().Add(-5 * time.Minute)}
}

func (c *LogsCollector) CollectLogs(machineID uuid.UUID) []AgentLogEntry {
	now := time.Now()
	hostname, _ := os.Hostname()
	platform := runtime.GOOS

	var logs []AgentLogEntry

	if platform == "windows" {
		logs = append(logs, AgentLogEntry{
			ID:        uuid.New(),
			MachineID: machineID,
			Hostname:  hostname,
			Platform:  platform,
			Level:     "INFO",
			Source:    "system",
			Message:   "Windows Event Log Service initialized cleanly.",
			Timestamp: now,
		})
	} else {
		logs = append(logs, AgentLogEntry{
			ID:        uuid.New(),
			MachineID: machineID,
			Hostname:  hostname,
			Platform:  platform,
			Level:     "INFO",
			Source:    "syslog",
			Message:   "systemd-journald started logging service.",
			Timestamp: now,
		})
	}

	c.lastProcessed = now
	return logs
}
