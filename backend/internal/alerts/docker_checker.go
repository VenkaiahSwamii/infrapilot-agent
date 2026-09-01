package alerts

import (
	"fmt"
	"strings"

	"infrapilot/backend/internal/models"
	"infrapilot/backend/internal/services"
)

// CheckDockerHealth evaluates Docker container states, restart loops, and container resource utilization.
func CheckDockerHealth(machine models.Machine, metric models.Metric, input services.SaveMetricInput) []models.LinuxAlert {
	var alerts []models.LinuxAlert

	for _, container := range input.DockerContainers {
		name := container.Name
		if name == "" {
			name = container.ID
			if len(name) > 12 {
				name = name[:12]
			}
		}

		st := strings.ToLower(container.State)
		status := strings.ToLower(container.Status)

		// 1. Container Exited / Dead Check
		if st == "exited" || st == "dead" || strings.Contains(status, "exited") || strings.Contains(status, "unhealthy") {
			alerts = append(alerts, models.LinuxAlert{
				Title:              "Docker Container Exited",
				Description:        fmt.Sprintf("Docker container '%s' (image: %s) is exited or dead", name, container.Image),
				Category:           "Docker",
				Component:          fmt.Sprintf("docker_%s", name),
				Source:             "DockerChecker",
				Type:               "docker_status",
				Severity:           "Critical",
				Priority:           models.MapSeverityToPriority("Critical"),
				Message:            fmt.Sprintf("Container %s Exited (%s)", name, container.Status),
				Status:             "OPEN",
				RecoverySuggestion: fmt.Sprintf("Inspect container logs (`docker logs %s`) and restart container (`docker start %s`).", name, name),
			})
		}

		// 2. Restart Loop Check (>3 restarts)
		if container.RestartCount > 3 {
			alerts = append(alerts, models.LinuxAlert{
				Title:              "Docker Container Restart Loop",
				Description:        fmt.Sprintf("Docker container '%s' has restarted %d times", name, container.RestartCount),
				Category:           "Docker",
				Component:          fmt.Sprintf("docker_%s", name),
				Source:             "DockerChecker",
				Type:               "docker_restart_loop",
				Severity:           "Warning",
				Priority:           models.MapSeverityToPriority("Warning"),
				Message:            fmt.Sprintf("Container %s Restarting (%d times)", name, container.RestartCount),
				MetricValue:        float64(container.RestartCount),
				Threshold:          3.0,
				Status:             "OPEN",
				RecoverySuggestion: fmt.Sprintf("Check container healthcheck, memory limits, and application entrypoint crash logs (`docker logs --tail 100 %s`).", name),
			})
		}

		// 3. Container High Resource Usage Check (>90% CPU)
		if container.CPUPercent > 90.0 {
			alerts = append(alerts, models.LinuxAlert{
				Title:              "Docker Container High CPU",
				Description:        fmt.Sprintf("Docker container '%s' CPU usage high: %.1f%%", name, container.CPUPercent),
				Category:           "Docker",
				Component:          fmt.Sprintf("docker_%s", name),
				Source:             "DockerChecker",
				Type:               "docker_cpu",
				Severity:           "Warning",
				Priority:           models.MapSeverityToPriority("Warning"),
				Message:            fmt.Sprintf("Container %s CPU High: %.1f%% > 90.0%%", name, container.CPUPercent),
				MetricValue:        container.CPUPercent,
				Threshold:          90.0,
				Status:             "OPEN",
				RecoverySuggestion: fmt.Sprintf("Check container processes using `docker top %s`. Apply CPU quota limits if necessary.", name),
			})
		}
	}

	return alerts
}
