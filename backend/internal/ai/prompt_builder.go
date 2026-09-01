package ai

import "fmt"

func BuildOperationalPrompt(userQuery string, telemetryContext string) string {
	return fmt.Sprintf(`You are InfraPilot AI, an elite SRE & Systems Architect assistant.
Analyze the following live telemetry context and answer the user query concisely with root cause, affected components, and recommended fix.

TELEMETRY CONTEXT:
%s

USER QUERY:
%s
`, telemetryContext, userQuery)
}
