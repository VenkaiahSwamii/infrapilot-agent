package ai

type TimelineStep struct {
	Time    string `json:"time"`
	Event   string `json:"event"`
	Details string `json:"details"`
}

type IncidentSummary struct {
	Summary   string         `json:"summary"`
	Timeline  []TimelineStep `json:"timeline"`
	RootCause string         `json:"root_cause"`
}

func DetectIncidentSummary() IncidentSummary {
	return IncidentSummary{
		Summary:   "A database connection pool was exhausted, triggering API response spikes and a Docker container restart.",
		RootCause: "Docker container exhausted RAM memory limit during vector embedding calculation.",
		Timeline: []TimelineStep{
			{Time: "10:21 UTC", Event: "CPU Usage Spike", Details: "Server-12 CPU reached 98% utilization"},
			{Time: "10:22 UTC", Event: "DB Latency Increase", Details: "PostgreSQL query latency increased from 15ms to 240ms"},
			{Time: "10:23 UTC", Event: "API Errors Appeared", Details: "HTTP 500 status rate jumped to 4.2% on /api/v1/aiops/analyze"},
			{Time: "10:24 UTC", Event: "Service Restarted", Details: "Docker daemon restarted container infra-api-server"},
		},
	}
}
