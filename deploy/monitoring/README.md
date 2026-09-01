# InfraPilot Enterprise Observability Stack

Sprint 9.5 - Complete monitoring, logging, tracing, and alerting platform.

## Architecture

```
                    Linux Agents
                          │
                    Metrics + Logs
                          │
           ┌──────────────┼──────────────┐
           │              │              │
      Node Exporter   Promtail      OpenTelemetry
           │              │              │
           ▼              ▼              ▼
      Prometheus       Loki         Jaeger
           │              │              │
           └──────────────┼──────────────┘
                          │
                     Grafana Dashboards
                          │
                    Alertmanager
                          │
             Email / Slack / Teams / Discord
```

## Components

### Metrics Collection
- **Prometheus** - Time-series metrics database and scraper
- **Node Exporter** - Host metrics (CPU, RAM, Disk, Network)
- **cAdvisor** - Container metrics (CPU, Memory, Network, Disk I/O)

### Logging
- **Loki** - Centralized log aggregation
- **Promtail** - Log collector and forwarder

### Tracing
- **Jaeger** - Distributed tracing
- **OpenTelemetry** - Instrumentation and trace export

### Visualization & Alerting
- **Grafana** - Dashboards and visualization
- **Alertmanager** - Alert routing and notifications

## Directory Structure

```
deploy/monitoring/
├── prometheus/
│   └── prometheus.yml          # Prometheus configuration
├── grafana/
│   ├── datasources/
│   │   └── datasources.yml     # Grafana datasource definitions
│   └── dashboards/
│       ├── dashboards.yml      # Dashboard provider configuration
│       ├── infrapilot-overview.json
│       ├── infrastructure-overview.json
│       └── health-score.json
├── loki/
│   └── loki.yml                # Loki configuration
├── promtail/
│   └── promtail.yml            # Promtail configuration
├── jaeger/
│   └── jaeger.yml              # Jaeger configuration
├── opentelemetry/
│   └── otel-collector.yml      # OpenTelemetry Collector config
├── alertmanager/
│   ├── alertmanager.yml        # Alertmanager configuration
│   └── templates/
│       └── email.tmpl          # Email notification templates
├── alerts/
│   ├── infrastructure.yml      # Infrastructure alert rules
│   └── application.yml         # Application alert rules
├── exporters/
│   ├── node-exporter.yml       # Node Exporter configuration
│   └── cadvisor.yml            # cAdvisor configuration
└── slo/
    └── sli-slo.yml             # SLI/SLO definitions
```

## Quick Start

### Prerequisites

- Docker & Docker Compose
- Kubernetes 1.20+ (for production)
- Helm 3.0+ (for Helm deployment)

### Docker Compose Deployment

1. **Create monitoring network:**
```bash
docker network create monitoring
```

2. **Start all monitoring services:**
```bash
docker-compose -f docker-compose.monitoring.yml up -d
```

### Kubernetes Deployment

1. **Create namespace:**
```bash
kubectl create namespace monitoring
```

2. **Deploy in order:**
```bash
# Deploy storage backends first
kubectl apply -f deploy/kubernetes/monitoring/
```

Or use Helm:
```bash
helm install infrapilot-monitoring ./helm/monitoring \
  --namespace monitoring \
  --create-namespace
```

## Configuration

### Prometheus

Scrape targets are configured in `prometheus/prometheus.yml`:

- **Backend API** - `http://backend:8080/api/v1/metrics`
- **PostgreSQL** - `postgres-exporter:9187`
- **Redis** - `redis-exporter:9121`
- **Qdrant** - `qdrant:6333/metrics`
- **Node Exporter** - `node-exporter:9100`
- **cAdvisor** - `cadvisor:8080/metrics`
- **Kubernetes** - Automatic service discovery

### Alertmanager

Configure notification channels in `alertmanager/alertmanager.yml`:

1. **Email:**
   - Update SMTP settings (smtp_smarthost, smtp_auth_username, smtp_auth_password)

2. **Slack:**
   - Update slack_api_url with your webhook URL

3. **Microsoft Teams:**
   - Update teams_api_url with your webhook URL

4. **Discord:**
   - Update discord_api_url with your webhook URL

### Loki

Log retention: 30 days (configurable in `loki/loki.yml`)

Labels automatically added:
- `host` - Hostname
- `application` - Application name
- `namespace` - Kubernetes namespace
- `container` - Container name
- `severity` - Log severity
- `environment` - Environment (production/staging)

### Jaeger

Access Jaeger UI: `http://jaeger:16686`

Trace collection endpoints:
- OTLP gRPC: `jaeger:4317`
- OTLP HTTP: `jaeger:4318`
- Zipkin: `jaeger:9411`

## Dashboards

### Infrastructure Overview
- CPU Usage
- Memory Usage
- Disk Usage
- Network Traffic
- Active Alerts
- Agent Status
- Container Status
- Node Status

### Health Score
- Overall Health Score (weighted metric)
- CPU Score (20%)
- Memory Score (20%)
- Disk Score (20%)
- Network Score (15%)
- Services Score (15%)
- Alerts Score (10%)

### Additional Dashboards
- Linux Servers
- Kubernetes Cluster
- Docker Containers
- Database Health
- API Performance
- Agent Status
- Network Usage
- Storage

## Alerts

### Critical Alerts (Immediate Notification)
- Node Down
- API Unavailable
- Database Unavailable
- Failed Deployments
- Qdrant Unavailable
- Redis Unavailable

### Warning Alerts (Batched)
- High CPU Usage (>90% for 5m)
- High Memory Usage (>90% for 5m)
- High Disk Usage (>90% for 5m)
- High API Latency (>2s for 5m)
- High Error Rate (>5% for 5m)
- Agent Offline (>45s)
- Slow Database Queries (>1s)

### Info Alerts (Daily Summary)
- Low Agent Registration Rate

### Alert Routing

| Severity | Notification Channels | Delay |
|----------|----------------------|-------|
| Critical | Email, Slack, Teams, Discord | Immediate |
| Warning | Email, Slack, Teams | 30s |
| Info | Email | Daily |

## SLI/SLO Definitions

### Service Level Indicators (SLIs)

| SLI | Target | Window |
|-----|--------|--------|
| API Availability | 99.9% | 30 days |
| API Latency (p95) | <500ms | 7 days |
| Dashboard Latency | <2s | 7 days |
| Agent Uptime | 99.5% | 30 days |
| Heartbeat Success Rate | 99.9% | 7 days |
| Metrics Ingestion Delay | <5s | 24 hours |
| Error Rate | <1% | 7 days |

### Error Budget Policy

- **100% budget remaining** - Normal operations
- **50% budget remaining** - Review changes, prioritize reliability
- **0% budget remaining** - Halt deployments, emergency response

## Security

### Authentication
- Grafana: Configure admin password
- Prometheus: Basic auth (recommended for production)

### TLS
- All services support TLS termination
- Use cert-manager for automatic certificate management

### Network Policies
```yaml
# Deny all ingress to monitoring namespace
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: deny-all-ingress
  namespace: monitoring
spec:
  podSelector: {}
  policyTypes:
    - Ingress
```

### RBAC
```yaml
# Grafana service account with read-only access
apiVersion: v1
kind: ServiceAccount
metadata:
  name: grafana
  namespace: monitoring
```

## Retention Policies

| Data Type | Retention | Configuration |
|-----------|-----------|---------------|
| Metrics | 30 days | Prometheus `--storage.tsdb.retention.time=30d` |
| Logs | 30 days | Loki `retention_period: 720h` |
| Traces | 7 days | Jaeger/Tempo |
| Alerts | 90 days | Alertmanager (external storage) |

## Backend Instrumentation

The Go backend exposes the following metrics at `/api/v1/metrics`:

### HTTP Metrics
- `infrapilot_http_requests_total` - Total requests
- `infrapilot_http_request_duration_seconds` - Request latency
- `infrapilot_http_response_size_bytes` - Response size

### WebSocket Metrics
- `infrapilot_websocket_active_connections` - Active connections
- `infrapilot_websocket_connections_total` - Total connections
- `infrapilot_websocket_messages_total` - Messages sent/received

### Database Metrics
- `infrapilot_db_queries_total` - Total queries
- `infrapilot_db_query_duration_seconds` - Query latency
- `infrapilot_db_connections_active` - Active connections

### Agent Metrics
- `infrapilot_agent_registrations_total` - Agent registrations
- `infrapilot_agent_heartbeats_total` - Heartbeats
- `infrapilot_agent_status` - Agent online/offline status
- `infrapilot_agent_metrics_received_total` - Metrics received

### System Metrics
- `infrapilot_health_score` - Overall health score (0-100)
- `infrapilot_uptime_seconds_total` - Uptime
- `infrapilot_errors_total` - Errors by type

## Validation Checklist

### Verify Deployments
```bash
# Check pods
kubectl get pods -n monitoring

# Expected output:
# NAME                                   READY   STATUS
# prometheus-xxxxx                       2/2     Running
# grafana-xxxxx                          1/1     Running
# loki-xxxxx                             1/1     Running
# promtail-xxxxx                         1/1     Running
# jaeger-xxxxx                           1/1     Running
# alertmanager-xxxxx                     1/1     Running
# node-exporter-xxxxx                    1/1     Running
# cadvisor-xxxxx                         1/1     Running
```

### Verify Services
```bash
kubectl get svc -n monitoring
```

### Verify Ingress
```bash
kubectl get ingress -n monitoring
```

### Check Logs
```bash
# Prometheus
kubectl logs deployment/prometheus -n monitoring

# Grafana
kubectl logs deployment/grafana -n monitoring

# Loki
kubectl logs deployment/loki -n monitoring
```

### Confirm Functionality

1. **Prometheus targets are UP:**
   - Navigate to `http://prometheus:9090/targets`
   - All targets should show `UP`

2. **Grafana dashboards load:**
   - Navigate to `http://grafana:3000`
   - Login with admin credentials
   - Verify dashboards display data

3. **Loki receives logs:**
   - Run a test query in Grafana Explore
   - Check for log entries

4. **Jaeger displays traces:**
   - Navigate to `http://jaeger:16686`
   - Search for recent traces

5. **Alerts fire correctly:**
   - Trigger a test alert (e.g., stop a service)
   - Verify alert appears in Alertmanager
   - Verify notification is received

## Troubleshooting

### Prometheus not scraping targets
- Check network connectivity
- Verify target endpoints are accessible
- Check Prometheus logs for errors

### Grafana not showing data
- Verify datasource configuration
- Check Prometheus is returning metrics
- Verify time range is correct

### Loki not receiving logs
- Check Promtail configuration
- Verify Loki endpoint is correct
- Check log file paths exist

### Jaeger not showing traces
- Verify OpenTelemetry instrumentation
- Check collector configuration
- Ensure traces are being exported

### Alerts not firing
- Check alert rules are loaded
- Verify Prometheus can reach Alertmanager
- Check notification channel configuration

## Performance Tuning

### Prometheus
- Increase `--query.max-concurrency` for heavy query loads
- Use remote write/read for long-term storage
- Enable TSDB compression

### Loki
- Increase `ingester` replicas for high log volume
- Use object storage (S3/GCS) instead of filesystem
- Enable log compression

### Grafana
- Enable caching (`--cache`)
- Use PostgreSQL for dashboard storage
- Increase `--max-concurrent-queries`

## Maintenance

### Daily
- Review alert summary
- Check health score dashboard
- Verify SLO compliance

### Weekly
- Review error budget consumption
- Analyze slow queries
- Update dashboards

### Monthly
- Audit alert rules
- Review retention policies
- Capacity planning

## References

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
- [Loki Documentation](https://grafana.com/docs/loki/latest/)
- [Jaeger Documentation](https://www.jaegertracing.io/docs/)
- [OpenTelemetry Documentation](https://opentelemetry.io/docs/)
- [Alertmanager Documentation](https://prometheus.io/docs/alerting/latest/alertmanager/)