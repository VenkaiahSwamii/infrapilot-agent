package alerts

import (
	"fmt"
	"strings"

	"infrapilot/backend/internal/models"
	"infrapilot/backend/internal/services"
)

// GenerateLinuxAlerts runs all 9 health checkers and applies Phase 13 Incident Correlation.
func GenerateLinuxAlerts(machine models.Machine, metric models.Metric, input services.SaveMetricInput) []models.LinuxAlert {
	var rawAlerts []models.LinuxAlert

	// 1. Run individual domain health checkers
	rawAlerts = append(rawAlerts, CheckCPUHealth(machine, metric, input)...)
	rawAlerts = append(rawAlerts, CheckMemoryHealth(machine, metric, input)...)
	rawAlerts = append(rawAlerts, CheckStorageHealth(machine, metric, input)...)
	rawAlerts = append(rawAlerts, CheckNetworkHealth(machine, metric, input)...)
	rawAlerts = append(rawAlerts, CheckProcessHealth(machine, metric, input)...)
	rawAlerts = append(rawAlerts, CheckServiceHealth(machine, metric, input)...)
	rawAlerts = append(rawAlerts, CheckDockerHealth(machine, metric, input)...)
	rawAlerts = append(rawAlerts, CheckKubernetesHealth(machine, metric, input)...)
	rawAlerts = append(rawAlerts, CheckSecurityHealth(machine, metric, input)...)

	// 2. Phase 13: Incident Correlation Engine
	return CorrelateAlerts(machine, rawAlerts)
}

// CorrelateAlerts aggregates multiple simultaneous critical alerts into a single incident to reduce alert noise.
func CorrelateAlerts(machine models.Machine, alerts []models.LinuxAlert) []models.LinuxAlert {
	if len(alerts) <= 1 {
		return alerts
	}

	hasHighCPU := false
	hasHighMemory := false
	hasHighLoad := false
	var correlatedIndices []int

	for i, alert := range alerts {
		if alert.Category == "CPU" && alert.Type == "cpu_usage" && alert.Severity == "Critical" {
			hasHighCPU = true
			correlatedIndices = append(correlatedIndices, i)
		}
		if alert.Category == "Memory" && alert.Type == "memory_percent" && alert.Severity == "Critical" {
			hasHighMemory = true
			correlatedIndices = append(correlatedIndices, i)
		}
		if alert.Category == "CPU" && alert.Type == "load_average" {
			hasHighLoad = true
			correlatedIndices = append(correlatedIndices, i)
		}
	}

	// If CPU, Memory, and Load are all critical simultaneously -> Correlate into "Infrastructure Overloaded"
	if hasHighCPU && hasHighMemory && hasHighLoad {
		var uncorrelatedAlerts []models.LinuxAlert
		var correlatedIDs []string

		for i, alert := range alerts {
			isIncludedInCorrelation := false
			for _, idx := range correlatedIndices {
				if i == idx {
					isIncludedInCorrelation = true
					correlatedIDs = append(correlatedIDs, alert.ID.String())
					break
				}
			}
			if !isIncludedInCorrelation {
				uncorrelatedAlerts = append(uncorrelatedAlerts, alert)
			}
		}

		hostname := machine.Hostname
		if hostname == "" {
			hostname = machine.ID.String()
		}

		incidentAlert := models.LinuxAlert{
			Title:              "Infrastructure Overloaded",
			Description:        fmt.Sprintf("Correlated Incident on %s: Simultaneous Critical CPU Usage, Memory Exhaustion, and System Load Spike.", hostname),
			Category:           "Infrastructure",
			Component:          "system_core",
			Source:             "CorrelationEngine",
			Type:               "correlated_incident",
			Severity:           "Critical",
			Priority:           "P1",
			Message:            fmt.Sprintf("Infrastructure Overloaded on %s (Correlated CPU + Memory + Load Incident)", hostname),
			Status:             "OPEN",
			IsCorrelated:       true,
			CorrelatedAlertIDs: strings.Join(correlatedIDs, ","),
			RecoverySuggestion: "SYSTEM OVERLOAD: Scale compute resources immediately. Terminate non-critical background jobs (`killall -9 <name>`) and clear page cache (`sync; echo 3 > /proc/sys/vm/drop_caches`).",
		}

		return append([]models.LinuxAlert{incidentAlert}, uncorrelatedAlerts...)
	}

	return alerts
}
