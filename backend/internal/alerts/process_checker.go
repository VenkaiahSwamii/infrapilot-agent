package alerts

import (
	"fmt"
	"strings"

	"infrapilot/backend/internal/models"
	"infrapilot/backend/internal/services"
)

// CheckProcessHealth evaluates high CPU processes, zombie processes, and process table volume.
func CheckProcessHealth(machine models.Machine, metric models.Metric, input services.SaveMetricInput) []models.LinuxAlert {
	var alerts []models.LinuxAlert

	zombieCount := 0
	highCPUProcesses := []string{}

	for _, proc := range input.Processes {
		// 1. High CPU Process Check (>80%)
		if proc.CPUPercent > 80.0 {
			highCPUProcesses = append(highCPUProcesses, fmt.Sprintf("%s (PID %d: %.1f%%)", proc.Name, proc.PID, proc.CPUPercent))
		}

		// 2. Zombie Process Check
		st := strings.ToLower(proc.Status)
		if st == "zombie" || st == "z" || strings.Contains(st, "defunct") {
			zombieCount++
		}
	}

	if len(highCPUProcesses) > 0 {
		procDetails := strings.Join(highCPUProcesses, ", ")
		alerts = append(alerts, models.LinuxAlert{
			Title:              "High CPU Process",
			Description:        fmt.Sprintf("Processes consuming >80%% CPU detected: %s", procDetails),
			Category:           "Process",
			Component:          "process_table",
			Source:             "ProcessChecker",
			Type:               "process_cpu",
			Severity:           "Warning",
			Priority:           models.MapSeverityToPriority("Warning"),
			Message:            fmt.Sprintf("High CPU Process: %s", procDetails),
			Status:             "OPEN",
			RecoverySuggestion: "Inspect process execution behavior (`top -p <pid>`). Terminate or re-nice rogue process if unneeded.",
		})
	}

	if zombieCount > 0 {
		alerts = append(alerts, models.LinuxAlert{
			Title:              "Zombie Process Detected",
			Description:        fmt.Sprintf("Detected %d defunct/zombie process(es) awaiting parent harvest", zombieCount),
			Category:           "Process",
			Component:          "process_table",
			Source:             "ProcessChecker",
			Type:               "zombie_process",
			Severity:           "Warning",
			Priority:           models.MapSeverityToPriority("Warning"),
			Message:            fmt.Sprintf("%d Zombie Process(es) Detected", zombieCount),
			MetricValue:        float64(zombieCount),
			Status:             "OPEN",
			RecoverySuggestion: "Identify parent PID (`ps -ef | grep defunct`) and signal parent process (`kill -CHLD <ppid>`) to clean up zombies.",
		})
	}

	// 3. High Process Count Check (>500 processes)
	if len(input.Processes) > 500 {
		alerts = append(alerts, models.LinuxAlert{
			Title:              "High Process Count",
			Description:        fmt.Sprintf("Total active system processes count high: %d processes", len(input.Processes)),
			Category:           "Process",
			Component:          "process_table",
			Source:             "ProcessChecker",
			Type:               "process_count",
			Severity:           "Warning",
			Priority:           models.MapSeverityToPriority("Warning"),
			Message:            fmt.Sprintf("High Process Count: %d > 500", len(input.Processes)),
			MetricValue:        float64(len(input.Processes)),
			Threshold:          500.0,
			Status:             "OPEN",
			RecoverySuggestion: "Check for process fork-bombs or runaway daemon worker spawns.",
		})
	}

	return alerts
}
