# High Availability Deployment Guide

## Prerequisites

- Docker 20.10+ and Docker Compose 2.0+ (for Docker Swarm)
- Kubernetes 1.24+ (for K8s deployment)
- Helm 3.0+ (optional, for K8s)
- 4vCPU, 8GB RAM minimum per node
- Load balancer with SSL certificate

## Quick Start

### Option 1: Docker Compose (Development/Testing)

```bash
# Start all services with 3 backend replicas
docker compose up --scale backend=3

# Verify all backends are running
docker compose ps

# Check NGINX load balancing
docker logs infrapilot-nginx | grep "upstream"
```

### Option 2: Docker Swarm (Production)

```bash
# Initialize Swarm (if not already done)
docker swarm init

# Deploy stack
docker stack deploy -c docker-compose.yml infrapilot

# Scale backend to 3 instances
docker service scale infrapilot_backend=3

# Verify services
docker service ls
```

### Option 3: Kubernetes (Production)

```bash
# Create namespace
kubectl create namespace infrapilot

# Apply configurations
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secret.yaml
kubectl apply -f k8s/postgres.yaml
kubectl apply -f k8s/redis.yaml
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/ingress.yaml

# Verify deployments
kubectl get pods -n infrapilot
kubectl get services -n infrapilot

# Check backend replicas
kubectl get deployment infrapilot-backend -n infrapilot
```

## Configuration

### 1. Environment Variables

Backend instances read from environment:

```yaml
# docker-compose.yml
backend:
  environment:
    - DB_HOST=postgres
    - DB_PORT=5432
    - REDIS_HOST=redis
    - REDIS_PORT=6379
    - QUEUE_TYPE=redis_streams
```

### 2. NGINX Configuration

The load balancer routes traffic:

```nginx
upstream backend {
    server backend1:8080 max_fails=3 fail_timeout=30s;
    server backend2:8080 max_fails=3 fail_timeout=30s;
    server backend3:8080 max_fails=3 fail_timeout=30s;
    
    least_conn;  # Load balancing algorithm
}
```

### 3. Redis Pub/Sub

Backends communicate via Redis:

```go
// Publish events to all backends
pubsub.PublishWebSocket(event)

// Subscribe to cross-backend events
pubsub.Subscribe(pubsub.WebSocketChannel, handler)
```

## Scaling

### Horizontal Scaling

**Docker Compose**:
```bash
# Scale backend instances
docker compose up --scale backend=5

# Scale frontend instances
docker compose up --scale frontend=3
```

**Kubernetes**:
```bash
# Scale backend
kubectl scale deployment infrapilot-backend --replicas=5 -n infrapilot

# Scale frontend
kubectl scale deployment infrapilot-frontend --replicas=3 -n infrapilot

# Auto-scaling based on CPU
kubectl autoscale deployment infrapilot-backend --cpu-percent=70 --min=3 --max=10 -n infrapilot
```

### Database Scaling

PostgreSQL read replicas:

```bash
# Add read replica (Kubernetes)
kubectl apply -f k8s/postgres-replica.yaml

# Update backend to use read replica for SELECT queries
kubectl set env deployment/infrapilot-backend \
  DB_READ_HOST=infrapilot-postgres-read \
  -n infrapilot
```

## Health Monitoring

### Check Service Health

```bash
# Backend health
curl http://localhost:8080/healthz

# Database health
docker exec infrapilot-postgres pg_isready -U postgres

# Redis health
docker exec infrapilot-redis redis-cli ping

# NGINX health
docker exec infrapilot-nginx nginx -t
```

### Monitor Logs

```bash
# Backend logs
docker compose logs -f backend

# NGINX logs
docker compose logs -f nginx

# PostgreSQL logs
docker compose logs -f postgres

# Kubernetes logs
kubectl logs -f deployment/infrapilot-backend -n infrapilot
```

## Failover Testing

### Test Backend Failover

```bash
# Kill one backend
docker stop infrapilot-backend.1

# Verify traffic routes to other backends
curl http://localhost:8080/healthz

# Check NGINX detects failure
docker logs infrapilot-nginx | grep "mark"
```

### Test Database Failover

```bash
# Promote replica to primary
docker exec infrapilot-postgres-replica pg_ctl promote -D /var/lib/postgresql/data

# Update connection string
kubectl set env deployment/infrapilot-backend DB_HOST=infrapilot-postgres-replica -n infrapilot

# Verify backend reconnects
kubectl rollout status deployment/infrapilot-backend -n infrapilot
```

### Test WebSocket Failover

```bash
# Connect dashboard to backend
wscat -c ws://localhost/ws

# Kill that backend
docker stop infrapilot-backend.2

# Dashboard should auto-reconnect via NGINX
# Verify in dashboard UI
```

## Load Balancer Configuration

### SSL/TLS Setup

```bash
# Copy SSL certificates
cp server.crt server.key infra/nginx/ssl/

# Update nginx.conf with domain name
server_name your-domain.com;

# Restart NGINX
docker compose restart nginx
```

### Custom Load Balancing

Edit `infra/nginx/nginx.conf`:

```nginx
upstream backend {
    # Round-robin (default)
    server backend1:8080;
    server backend2:8080;
    server backend3:8080;
    
    # Or use least_conn for least connections
    least_conn;
    
    # Or ip_hash for sticky sessions
    ip_hash;
}
```

## Monitoring & Alerting

### Access Dashboards

```bash
# Prometheus
open http://localhost:9091

# Grafana
open http://localhost:3000
# Login: admin / admin

# Jaeger (tracing)
open http://localhost:16686
```

### Key Metrics to Monitor

- Backend request latency (p95, p99)
- Error rate (5xx responses)
- WebSocket connection count
- Database connection pool usage
- Redis memory usage
- Pod restart count

### Alerts

Configure in `grafana/dashboards/` or Prometheus rules:

```yaml
# Example alert rule
groups:
  - name: infrapilot_alerts
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
```

## Backup & Restore

### Database Backup

```bash
# Manual backup
docker exec infrapilot-postgres pg_dump -U postgres infrapilot_enterprise > backup.sql

# Automated daily backup (add to crontab)
0 2 * * * docker exec infrapilot-postgres pg_dump -U postgres infrapilot_enterprise | gzip > /backups/backup_$(date +\%Y\%m\%d).sql.gz
```

### Restore Database

```bash
# Stop backend services
docker compose stop backend worker

# Restore backup
docker exec -i infrapilot-postgres psql -U postgres infrapilot_enterprise < backup.sql

# Restart services
docker compose start backend worker
```

## Troubleshooting

### Backend Not Receiving Traffic

```bash
# Check NGINX upstream status
docker exec infrapilot-nginx curl http://backend/healthz

# Verify backend health endpoint
curl http://localhost:8080/healthz

# Check backend logs
docker compose logs backend
```

### WebSocket Connections Failing

```bash
# Verify NGINX WebSocket proxy config
docker exec infrapilot-nginx cat /etc/nginx/nginx.conf | grep -A 10 "location /ws"

# Check backend WebSocket handler
docker compose logs backend | grep "WebSocket"

# Test WebSocket connection
wscat -c ws://localhost/ws
```

### Database Connection Issues

```bash
# Check PostgreSQL is accepting connections
docker exec infrapilot-postgres netstat -an | grep 5432

# Verify credentials
docker exec infrapilot-postgres psql -U postgres -c "SELECT 1"

# Check connection pool
curl http://localhost:9090/metrics | grep db_connections
```

## Performance Tuning

### Backend Tuning

```bash
# Increase worker processes
export GIN_MODE=release
export GIN_WORKER_PROCESSES=4

# Database connection pool
export DB_MAX_OPEN_CONNS=50
export DB_MAX_IDLE_CONNS=10
```

### Redis Tuning

```bash
# Increase memory
docker compose exec redis redis-cli CONFIG SET maxmemory 1gb

# Enable keyspace notifications for Pub/Sub
docker compose exec redis redis-cli CONFIG SET notify-keyspace-events KEA
```

### PostgreSQL Tuning

```bash
# Update postgresql.conf
shared_buffers = 256MB
work_mem = 16MB
maintenance_work_mem = 128MB
effective_cache_size = 1GB
```

## Security Hardening

### Network Isolation

```bash
# Verify Docker network is private
docker network inspect infrapilot

# Ensure no ports exposed except 80/443
docker compose ps
```

### mTLS Configuration

```bash
# Generate certificates
openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365

# Update NGINX config
ssl_certificate /etc/nginx/ssl/server.crt;
ssl_certificate_key /etc/nginx/ssl/server.key;
```

### Secrets Management

```bash
# Use Docker Secrets (Swarm)
echo "postgres_password" | docker secret create db_password -

# Use Kubernetes Secrets
kubectl create secret generic infrapilot-secrets \
  --from-literal=db_password=postgres \
  --from-literal=redis_password=
```

## Maintenance

### Rolling Update

```bash
# Docker Swarm
docker service update --image infrapilot/backend:1.0.1 infrapilot_backend

# Kubernetes
kubectl set image deployment/infrapilot-backend backend=infrapilot/backend:1.0.1 -n infrapilot
kubectl rollout status deployment/infrapilot-backend -n infrapilot
```

### Zero-Downtime Restart

```bash
# Restart one backend at a time
docker compose up -d --scale backend=4
docker compose stop backend.1
docker compose rm -f backend.1
docker compose up -d --scale backend=3
```

## Support

For issues or questions:
- Documentation: https://docs.infrapilot.io
- GitHub Issues: https://github.com/your-org/infrapilot-enterprise/issues
- Slack: #infrapilot-support