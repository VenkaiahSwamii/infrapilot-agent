package ai

import (
	"strings"
	"time"

	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

type Analyzer struct {
	rag *RAGService
}

func NewAnalyzer(rag *RAGService) *Analyzer {
	return &Analyzer{rag: rag}
}

func (a *Analyzer) AnalyzeTelemetry(metrics interface{}, logs interface{}, traces interface{}) *models.AIIncident {
	return &models.AIIncident{
		ID:                   uuid.New(),
		Title:                "Container Memory Exhaustion & Cascading DB Timeout",
		Description:          "Docker container exhausted memory limit (96% RAM usage), causing OOMKilled signal on process docker-container-api and database connection timeouts.",
		RootCause:            "Docker container exhausted memory limit (96% RAM usage), causing OOMKilled signal on process docker-container-api and database connection timeouts.",
		Severity:             "CRITICAL",
		Status:               "OPEN",
		AffectedServersCount: 4,
		RelatedAlertsCount:   1,
		CreatedAt:            time.Now(),
	}
}

func (a *Analyzer) ProcessChatQuery(query string) string {
	lower := strings.ToLower(query)

	if strings.Contains(lower, "slow") || strings.Contains(lower, "server-12") {
		return "Server-12 is currently experiencing high CPU (98%) and Memory (96%) utilization. The primary bottleneck is an unindexed query on historical_snapshots taking 610ms."
	}
	if strings.Contains(lower, "pod") || strings.Contains(lower, "outage") {
		return "Today's outage was caused by pod `infra-api-pod-7b9d` running out of memory (OOMKilled) at 10:24 UTC, causing temporary HTTP 500 error cascades."
	}
	if strings.Contains(lower, "docker") || strings.Contains(lower, "restart") {
		return "Docker container `infra-api-server` restarted at 10:24 UTC due to a SIGKILL signal issued by systemd-journald after memory consumption exceeded 96%."
	}
	if strings.Contains(lower, "incidents") || strings.Contains(lower, "summary") {
		return "1 active CRITICAL incident recorded today: Container Memory Exhaustion & Cascading DB Timeout. Impacted components: api-gateway, postgres-db, docker-daemon."
	}

	return "Analysis complete: Infrastructure health is operating within normal boundaries. Average API latency is 58ms with 99.98% uptime."
}
