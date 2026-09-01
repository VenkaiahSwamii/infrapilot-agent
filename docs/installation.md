# InfraPilot Enterprise Installation Guide

This guide covers installation methods for InfraPilot Enterprise.

## Prerequisites

### Minimum Requirements

**Server:**
- CPU: 2 cores
- RAM: 4 GB
- Disk: 20 GB SSD
- OS: Linux (Ubuntu 20.04+, CentOS 8+, RHEL 8+)

**Supported Platforms:**
- Ubuntu 20.04 / 22.04 / 24.04 LTS
- CentOS Stream 8 / 9
- RHEL 8 / 9
- Debian 11 / 12
- Alpine Linux 3.15+

**Client:**
- Browser: Chrome 90+, Firefox 88+, Safari 14+, Edge 90+
- Network: HTTPS access to InfraPilot server

### Required Software

- Docker Engine 24.0+
- Docker Compose 2.20+
- Git
- Curl

## Installation Methods

### Method 1: Docker Compose (Recommended)

This is the fastest way to get InfraPilot running.

**Step 1: Clone Repository**

```bash
git clone https://github.com/yourusername/infrapilot-enterprise.git
cd infrapilot-enterprise
```

**Step 2: Configure Environment**

```bash
cp backend/.env.example backend/.env
```

Edit `backend/.env` and set secure passwords:

```env
DB_PASSWORD=your-secure-password-here
JWT_SECRET=your-jwt-secret-key-min-32-characters
```

**Step 3: Start Services**

```bash
# Start core services
docker compose up -d

# Start with monitoring (optional)
docker compose --profile observability up -d

# Start with agent (optional)
docker compose --profile agent up -d
```

**Step 4: Verify Installation**

```bash
# Check backend health
curl http://localhost:8080/healthz

# Access frontend
# Open browser: http://localhost
```

**Step 5: Login**

Default credentials (change immediately):
- Email: `admin@infrapilot.io`
- Password: `admin123`

### Method 2: Kubernetes

**Step 1: Create Namespace**

```bash
kubectl create namespace infrapilot
```

**Step 2: Configure Secrets**

```bash
# Edit secrets with your values
kubectl create secret generic infrapilot-secrets \
  --from-literal=DB_PASSWORD='your-password' \
  --from-literal=JWT_SECRET='your-jwt-secret' \
  -n infrapilot
```

**Step 3: Apply Configurations**

```bash
# Apply in order
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secret.yaml
kubectl apply -f k8s/postgres.yaml
kubectl apply -f k8s/redis.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/ingress.yaml
```

**Step 4: Verify Deployment**

```bash
kubectl get pods -n infrapilot
kubectl get services -n infrapilot
```

### Method 3: Install Script (Linux)

**Step 1: Download Installer**

```bash
curl -fsSL https://install.infrapilot.io | bash
```

**Step 2: Follow Prompts**

The installer will:
1. Check system requirements
2. Install Docker
3. Deploy InfraPilot
4. Configure firewall
5. Start services

### Method 4: Windows Installer

Download `InfraPilotSetup.exe` from the downloads page.

**Installation Steps:**
1. Run installer as Administrator
2. Accept license agreement
3. Choose installation directory
4. Set admin password
5. Complete installation

The installer will:
- Install Docker Desktop
- Deploy InfraPilot containers
- Configure Windows Firewall
- Create desktop shortcut
- Start services automatically

## Post-Installation

### 1. Configure SSL/TLS

**Using Let's Encrypt (Recommended)**

```bash
# Install cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.14.0/cert-manager.yaml

# Create ClusterIssuer
kubectl apply -f - <<EOF
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: your-email@example.com
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - http01:
        ingress:
          class: nginx
EOF
```

### 2. Configure Firewall

```bash
# Ubuntu/Debian
ufw allow 80/tcp
ufw allow 443/tcp
ufw allow 8080/tcp  # API
ufw allow 9090/tcp  # Metrics
ufw enable
```

### 3. Backup Configuration

```bash
# Create backup script
cat > /opt/infrapilot/backup.sh << 'EOF'
#!/bin/bash
BACKUP_DIR="/opt/infrapilot/backups"
DATE=$(date +%Y%m%d_%H%M%S)

# Backup PostgreSQL
docker exec infrapilot-postgres pg_dump -U postgres infrapilot_enterprise > $BACKUP_DIR/db_$DATE.sql

# Backup Redis
docker exec infrapilot-redis redis-cli BGSAVE

# Compress old backups
find $BACKUP_DIR -name "*.sql" -mtime +7 | gzip
EOF

chmod +x /opt/infrapilot/backup.sh

# Add to cron
(crontab -l 2>/dev/null; "0 2 * * * /opt/infrapilot/backup.sh") | crontab -
```

### 4. Configure Monitoring

Access monitoring dashboards:
- Grafana: http://your-server:3000
- Prometheus: http://your-server:9091
- Jaeger: http://your-server:16686

Default Grafana credentials:
- Username: `admin`
- Password: (set via `GRAFANA_PASSWORD` env var)

## Verification

### Health Checks

```bash
# Backend health
curl http://localhost:8080/healthz
# Expected: {"status": "healthy"}

# API test
curl http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@infrapilot.io","password":"admin123"}'
```

### Database Connection

```bash
docker exec -it infrapilot-postgres psql -U postgres -d infrapilot_enterprise -c "SELECT version();"
```

### Redis Connection

```bash
docker exec -it infrapilot-redis redis-cli ping
# Expected: PONG
```

## Configuration

### Environment Variables

**Backend (.env)**

| Variable | Description | Default |
|----------|-------------|---------|
| DB_HOST | PostgreSQL host | postgres |
| DB_PORT | PostgreSQL port | 5432 |
| DB_USER | PostgreSQL user | postgres |
| DB_PASSWORD | PostgreSQL password | postgres |
| DB_NAME | Database name | infrapilot_enterprise |
| REDIS_HOST | Redis host | redis |
| REDIS_PORT | Redis port | 6379 |
| JWT_SECRET | JWT signing secret | (required) |
| PROMETHEUS_ENABLED | Enable Prometheus metrics | true |
| OTEL_EXPORTER_OTLP_ENDPOINT | Jaeger endpoint | jaeger:4317 |

### Resource Limits

Edit `docker-compose.yml` to adjust resource allocation:

```yaml
deploy:
  resources:
    limits:
      cpus: '2'
      memory: 2G
    reservations:
      cpus: '1'
      memory: 1G
```

## Upgrading

### Docker Compose

```bash
git pull
docker compose down
docker compose up -d --build
```

### Kubernetes

```bash
kubectl apply -f k8s/ -n infrapilot
```

## Uninstallation

### Docker Compose

```bash
docker compose down -v
rm -rf backend/.env
```

### Kubernetes

```bash
kubectl delete namespace infrapilot
```

## Troubleshooting

### Backend won't start

```bash
# Check logs
docker logs infrapilot-backend

# Common issues:
# 1. Database not ready - wait for postgres health check
# 2. Port conflict - check if 8080 is in use
# 3. Invalid .env - verify configuration
```

### Frontend shows blank page

```bash
# Check nginx logs
docker logs infrapilot-frontend

# Verify API_URL in nginx.conf
# Rebuild frontend
docker compose build frontend
docker compose up -d frontend
```

### Agent won't connect

```bash
# Verify agent config
cat agent/config.yaml

# Check backend WebSocket endpoint
curl -I http://localhost:8080/ws

# Check firewall allows port 8080
```

## Support

- Documentation: https://docs.infrapilot.io
- GitHub Issues: https://github.com/yourusername/infrapilot-enterprise/issues
- Email: support@infrapilot.io