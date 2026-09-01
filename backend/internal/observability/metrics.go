package observability

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all Prometheus metrics for InfraPilot
var Metrics = struct {
	// HTTP metrics
	HTTPRequestsTotal   *prometheus.CounterVec
	HTTPRequestDuration *prometheus.HistogramVec
	HTTPResponseSize    *prometheus.HistogramVec

	// WebSocket metrics
	WebSocketActiveConnections prometheus.Gauge
	WebSocketConnectionsTotal  *prometheus.CounterVec
	WebSocketMessagesTotal     *prometheus.CounterVec

	// Database metrics
	DBQueryDuration     *prometheus.HistogramVec
	DBConnectionsActive prometheus.Gauge
	DBConnectionsIdle   prometheus.Gauge
	DBQueriesTotal      *prometheus.CounterVec

	// Agent metrics
	AgentRegistrationsTotal     prometheus.Counter
	AgentHeartbeatsTotal        *prometheus.CounterVec
	AgentLastHeartbeatTimestamp prometheus.Gauge
	AgentStatus                 *prometheus.GaugeVec
	AgentMetricsReceivedTotal   prometheus.Counter

	// System metrics
	HealthScore   prometheus.Gauge
	UptimeSeconds prometheus.Counter

	// Error metrics
	ErrorsTotal *prometheus.CounterVec

	// Background job metrics
	JobsTotal   *prometheus.CounterVec
	JobDuration *prometheus.HistogramVec
}{
	// HTTP metrics
	HTTPRequestsTotal: promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "infrapilot_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	),
	HTTPRequestDuration: promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "infrapilot_http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: []float64{0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path"},
	),
	HTTPResponseSize: promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "infrapilot_http_response_size_bytes",
			Help:    "HTTP response size in bytes",
			Buckets: []float64{100, 1000, 10000, 100000, 1000000},
		},
		[]string{"method", "path"},
	),

	// WebSocket metrics
	WebSocketActiveConnections: promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "infrapilot_websocket_active_connections",
			Help: "Current number of active WebSocket connections",
		},
	),
	WebSocketConnectionsTotal: promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "infrapilot_websocket_connections_total",
			Help: "Total number of WebSocket connections",
		},
		[]string{"status"},
	),
	WebSocketMessagesTotal: promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "infrapilot_websocket_messages_total",
			Help: "Total number of WebSocket messages",
		},
		[]string{"direction", "type"},
	),

	// Database metrics
	DBQueryDuration: promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "infrapilot_db_query_duration_seconds",
			Help:    "Database query latency in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 5},
		},
		[]string{"operation", "table"},
	),
	DBConnectionsActive: promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "infrapilot_db_connections_active",
			Help: "Number of active database connections",
		},
	),
	DBConnectionsIdle: promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "infrapilot_db_connections_idle",
			Help: "Number of idle database connections",
		},
	),
	DBQueriesTotal: promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "infrapilot_db_queries_total",
			Help: "Total number of database queries",
		},
		[]string{"operation", "table", "status"},
	),

	// Agent metrics
	AgentRegistrationsTotal: promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "infrapilot_agent_registrations_total",
			Help: "Total number of agent registrations",
		},
	),
	AgentHeartbeatsTotal: promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "infrapilot_agent_heartbeats_total",
			Help: "Total number of agent heartbeats",
		},
		[]string{"status"},
	),
	AgentLastHeartbeatTimestamp: promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "infrapilot_agent_last_heartbeat_timestamp",
			Help: "Timestamp of the last agent heartbeat",
		},
	),
	AgentStatus: promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "infrapilot_agent_status",
			Help: "Agent status (1=online, 0=offline)",
		},
		[]string{"agent_id", "host", "status"},
	),
	AgentMetricsReceivedTotal: promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "infrapilot_agent_metrics_received_total",
			Help: "Total number of metrics received from agents",
		},
	),

	// System metrics
	HealthScore: promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "infrapilot_health_score",
			Help: "Overall system health score (0-100)",
		},
	),
	UptimeSeconds: promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "infrapilot_uptime_seconds_total",
			Help: "Total uptime in seconds",
		},
	),

	// Error metrics
	ErrorsTotal: promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "infrapilot_errors_total",
			Help: "Total number of errors",
		},
		[]string{"type", "severity"},
	),

	// Background job metrics
	JobsTotal: promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "infrapilot_jobs_total",
			Help: "Total number of background jobs",
		},
		[]string{"job_type", "status"},
	),
	JobDuration: promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "infrapilot_job_duration_seconds",
			Help:    "Background job duration in seconds",
			Buckets: []float64{1, 5, 10, 30, 60, 300, 600},
		},
		[]string{"job_type"},
	),
}

// StartUptimeCounter starts the uptime counter
func StartUptimeCounter() {
	go func() {
		startTime := time.Now()
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			Metrics.UptimeSeconds.Add(60)
			_ = startTime
		}
	}()
}

// RecordHTTPRequest records an HTTP request metric
func RecordHTTPRequest(method, path, status string, duration time.Duration, size int64) {
	Metrics.HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
	Metrics.HTTPRequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
	Metrics.HTTPResponseSize.WithLabelValues(method, path).Observe(float64(size))
}

// RecordWebSocketConnection records a WebSocket connection event
func RecordWebSocketConnection(status string) {
	Metrics.WebSocketConnectionsTotal.WithLabelValues(status).Inc()
	if status == "connected" {
		Metrics.WebSocketActiveConnections.Inc()
	} else {
		Metrics.WebSocketActiveConnections.Dec()
	}
}

// RecordWebSocketMessage records a WebSocket message
func RecordWebSocketMessage(direction, messageType string) {
	Metrics.WebSocketMessagesTotal.WithLabelValues(direction, messageType).Inc()
}

// RecordDBQuery records a database query
func RecordDBQuery(operation, table, status string, duration time.Duration) {
	Metrics.DBQueriesTotal.WithLabelValues(operation, table, status).Inc()
	Metrics.DBQueryDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
}

// RecordAgentRegistration records an agent registration
func RecordAgentRegistration() {
	Metrics.AgentRegistrationsTotal.Inc()
}

// RecordAgentHeartbeat records an agent heartbeat
func RecordAgentHeartbeat(status string) {
	Metrics.AgentHeartbeatsTotal.WithLabelValues(status).Inc()
	Metrics.AgentLastHeartbeatTimestamp.Set(float64(time.Now().Unix()))
}

// SetAgentStatus sets the status of an agent
func SetAgentStatus(agentID, host, status string) {
	value := 0.0
	if status == "online" || status == "active" {
		value = 1.0
	}
	Metrics.AgentStatus.WithLabelValues(agentID, host, status).Set(value)
}

// RecordAgentMetricsReceived records metrics received from an agent
func RecordAgentMetricsReceived(count int) {
	Metrics.AgentMetricsReceivedTotal.Add(float64(count))
}

// RecordErrorMetric records an error metric
func RecordErrorMetric(errorType, severity string) {
	Metrics.ErrorsTotal.WithLabelValues(errorType, severity).Inc()
}

// RecordJob records a background job execution
func RecordJob(jobType, status string, duration time.Duration) {
	Metrics.JobsTotal.WithLabelValues(jobType, status).Inc()
	Metrics.JobDuration.WithLabelValues(jobType).Observe(duration.Seconds())
}

// SetHealthScore sets the overall health score
func SetHealthScore(score float64) {
	Metrics.HealthScore.Set(score)
}

// CalculateHealthScore calculates the overall health score based on weighted metrics
// Weights: CPU 20%, Memory 20%, Disk 20%, Network 15%, Services 15%, Alerts 10%
func CalculateHealthScore(
	cpuUsage, memoryUsage, diskUsage, networkReliability, servicesUptime, alertsScore float64,
) float64 {
	// Invert resource usage to get scores (100% usage = 0 score)
	cpuScore := 100 - cpuUsage
	memoryScore := 100 - memoryUsage
	diskScore := 100 - diskUsage

	// Weighted calculation
	healthScore := (cpuScore * 0.20) +
		(memoryScore * 0.20) +
		(diskScore * 0.20) +
		(networkReliability * 0.15) +
		(servicesUptime * 0.15) +
		(alertsScore * 0.10)

	return healthScore
}

type BackendSystemMetrics struct {
	KubernetesNodes int     `json:"kubernetes_nodes"`
	UptimeSeconds   int     `json:"uptime_seconds"`
	CPUUsage        float64 `json:"cpu_usage"`
	MemoryUsage     float64 `json:"memory_usage"`
	APILatency      float64 `json:"api_latency"`
	DBConnections   int     `json:"db_connections"`
	RedisUsage      float64 `json:"redis_usage"`
	RequestsPerSec  int64   `json:"requests_per_sec"`
	ConnectedAgents int     `json:"connected_agents"`
}

func GetBackendSystemMetrics() BackendSystemMetrics {
	return BackendSystemMetrics{
		KubernetesNodes: 3,
		UptimeSeconds:   86400,
		CPUUsage:        12.5,
		MemoryUsage:     34.2,
		APILatency:      4.2,
		DBConnections:   18,
		RedisUsage:      8.4,
		RequestsPerSec:  125,
		ConnectedAgents: 12,
	}
}
