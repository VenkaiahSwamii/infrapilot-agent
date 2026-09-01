package subscribers

import (
	"log"

	"infrapilot/backend/internal/events"
)

// WriteAudit logs events for audit purposes
func WriteAudit(e events.Event) {
	log.Printf("Writing audit log for event: %s at %v", e.Name(), e.Timestamp())

	// The actual database persistence is handled by SaveAuditLog in database.go
	// This subscriber can be used for additional audit operations like:
	// 1. Sending audit logs to external SIEM systems
	// 2. Writing to audit log files
	// 3. Triggering compliance checks
	// 4. Notifying security teams

	switch ev := e.(type) {
	case events.UserLoginEvent:
		log.Printf("AUDIT: User %s (ID: %s) logged in from %s", ev.Username, ev.UserID, ev.IPAddress)
	case events.UserLogoutEvent:
		log.Printf("AUDIT: User %s (ID: %s) logged out", ev.Username, ev.UserID)
	case events.CommandExecutedEvent:
		log.Printf("AUDIT: Command executed on machine %s by %s: %s (exit code: %d)",
			ev.MachineID, ev.ExecutedBy, ev.Command, ev.ExitCode)
	case events.FileUploadedEvent:
		log.Printf("AUDIT: File %s uploaded to machine %s by %s (size: %d bytes)",
			ev.FileName, ev.MachineID, ev.UploadedBy, ev.FileSize)
	case events.FileDeletedEvent:
		log.Printf("AUDIT: File %s deleted from machine %s by %s",
			ev.FileName, ev.MachineID, ev.DeletedBy)
	case events.TerminalOpenedEvent:
		log.Printf("AUDIT: Terminal session opened on machine %s by user %s (session: %s)",
			ev.MachineID, ev.Username, ev.SessionID)
	case events.TerminalClosedEvent:
		log.Printf("AUDIT: Terminal session closed on machine %s by user %s (session: %s, duration: %d seconds)",
			ev.MachineID, ev.Username, ev.SessionID, ev.Duration)
	case events.DockerContainerStoppedEvent:
		log.Printf("AUDIT: Docker container %s stopped on machine %s (image: %s, exit code: %d)",
			ev.ContainerName, ev.MachineID, ev.Image, ev.ExitCode)
	case events.KubernetesPodFailedEvent:
		log.Printf("AUDIT: Kubernetes pod %s failed in namespace %s on machine %s (reason: %s)",
			ev.PodName, ev.Namespace, ev.MachineID, ev.Reason)
	default:
		log.Printf("AUDIT: Event %s occurred at %v", e.Name(), e.Timestamp())
	}
}
