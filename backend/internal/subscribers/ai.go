package subscribers

import (
	"log"

	"infrapilot/backend/internal/events"
)

// RunAI performs AI analysis on events
func RunAI(e events.Event) {
	log.Printf("Running AI analysis for event: %s", e.Name())

	switch ev := e.(type) {
	case events.MetricReceivedEvent:
		log.Printf("AI analyzing metrics for machine %s - CPU: %.2f%%, Memory: %.2f%%, Disk: %.2f%%",
			ev.MachineID, ev.CPU, ev.Memory, ev.Disk)

		// TODO: Implement AI analysis
		// This will analyze metrics using AI/ML models to:
		// 1. Detect anomalies
		// 2. Predict future resource usage
		// 3. Suggest optimizations
		// 4. Identify patterns
		//
		// Example:
		// anomalyScore := aiModel.DetectAnomaly(ev.CPU, ev.Memory, ev.Disk)
		// if anomalyScore > 0.8 {
		//     log.Printf("Anomaly detected for machine %s: score=%.2f", ev.MachineID, anomalyScore)
		// }

	case events.MachineOfflineEvent:
		log.Printf("AI analyzing offline event for machine %s", ev.MachineID)

		// TODO: Implement AI analysis for offline events
		// Predict if the machine will stay offline
		// Suggest remediation actions

	default:
		log.Printf("AI analysis not applicable for event type: %s", e.Name())
	}
}
