package subscribers

import (
	"log"

	"infrapilot/backend/internal/events"
)

// CheckAlerts evaluates metrics against alert rules and triggers alerts
func CheckAlerts(e events.Event) {
	log.Printf("Checking alerts for event: %s", e.Name())

	switch ev := e.(type) {
	case events.MetricReceivedEvent:
		log.Printf("Checking alert rules for machine %s - CPU: %.2f%%, Memory: %.2f%%, Disk: %.2f%%",
			ev.MachineID, ev.CPU, ev.Memory, ev.Disk)

		// TODO: Implement alert rule evaluation
		// This will check the metric values against configured alert rules
		// and publish AlertCreatedEvent if thresholds are breached
		//
		// Example:
		// if ev.CPU > 90.0 {
		//     alertEvent := events.AlertCreatedEvent{
		//         AlertID:     generateAlertID(),
		//         MachineID:   ev.MachineID,
		//         AlertType:   "cpu_high",
		//         Severity:    "critical",
		//         Message:     "CPU usage exceeded 90%",
		//         MetricValue: ev.CPU,
		//         Threshold:   90.0,
		//         Time:        time.Now(),
		//     }
		//     eventBus.Publish(alertEvent)
		// }

	default:
		log.Printf("Alert check not applicable for event type: %s", e.Name())
	}
}
