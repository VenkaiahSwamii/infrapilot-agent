package ai

type RAGService struct{}

func NewRAGService() *RAGService {
	return &RAGService{}
}

func (r *RAGService) RetrieveContext(query string) string {
	return `Live Telemetry Summary:
- Server-12 CPU: 98%, Memory: 96%
- Docker Daemon: Container infra-api-server restarted at 10:24 UTC
- Database PostgreSQL: Active connection pool 98/100, query latency 240ms
- Traces: POST /api/v1/aiops/analyze latency 920ms (is_slow=true)
- Log: "systemd-journald: OOMKilled process docker-container-api"`
}
