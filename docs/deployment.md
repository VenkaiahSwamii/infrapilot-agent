# InfraPilot Enterprise Deployment Guide

This guide covers deployment architectures and best practices.

## Deployment Architectures

### Single Server (Small Scale)

Suitable for: Development, testing, small teams (< 50 monitored machines)

```
                    ┌─────────────────┐
                    │   Load Balancer │
                    │   (Nginx)       │
                    └────────┬────────┘
                             │
                    ┌────────▼────────┐
                    │   Backend       │
                    │   (1 replica)   │
                    └────────┬────────┘
                             │
          ┌──────────────────┼──────────────────┐
          │                  │                  │
    ┌─────▼─────┐     ┌─────▼─────┐     ┌─────▼─────┐
    │PostgreSQL │     │   Redis   │     │  Workers  │
    │  (1 node) │     │  (1 node) │     │ (1 proc)  │
    └───────────┘     └───────────┘     └───────────┘
```

### High Availability (Medium Scale)

Suitable for: Production, medium teams (50-500 monitored machines)

```
                    ┌─────────────────┐
                    │   Load Balancer │
                    │   (HAProxy)     │
                    └────────┬────────┘
                             │
              ┌──────────────┼──────────────┐
              │              │              │
        ┌─────▼─────┐  ┌─────▼─────┐  ┌─────▼─────┐
        │ Backend 1 │  │ Backend 2 │  │ Backend 3 │
        └─────┬─────┘  └─────┬─────┘  └─────┬─────┘
              │              │              │
              └──────────────┼──────────────┘
                             │
                    ┌────────▼────────┐
                    │  PostgreSQL HA  │
                    │  (Primary +     │
                    │   Replicas)     │
                    └────────┬────────┘
                             │
                    ┌────────▼────────┐
                    │  Redis Cluster  │
                    └────────┬────────┘
                             │
                    ┌────────▼────────┐
                    │   Workers Pool  │
                    │  (Multi-proc)   │
                    └─────────────────┘
```

### Distributed (Enterprise Scale)

Suitable for: Enterprise, large teams (500+ monitored machines)

```
                            ┌─────────────┐
                            │    CDN      │
                            │  (CloudFlare)│
                            └──────┬──────┘
                                   │
                            ┌──────▼──────┐
                            │   Frontend  │
                            │  (S3 + CF)  │
                            └──────┬──────┘
                                   │
                            ┌──────▼──────┐
                            │ API Gateway │
                            │  (Kong/TS)  │
                            └──────┬──────┘
                                   │
              ┌────────────────────┼────────────────────┐
              │                    │                    │
       ┌──────▼──────┐      ┌──────▼──────┐      ┌──────▼──────┐
       │ Backend -   │      │ Backend -   │      │ Backend -   │
       │  Region US  │      │  Region EU  │      │  Region ASIA│
       └──────┬──────┘      └──────┬──────┘      └──────┬──────┘
              │                    │                    │
              └────────────────────┼────────────────────┘
                                   │
                            ┌──────▼──────┐
                            │ PostgreSQL  │
                            │  (Primary + │
                            │  Read Repl.)│
                            └──────┬──────┘
                                   │
                            ┌──────▼──────┐
                            │ Redis       │
                            │  (Cluster)  │
                            └──────┬──────┘
                                   │
                            ┌──────▼────────┐
                            │  Message Bus  │
                            │  (Kafka)      │
                            └──────┬────────┘
                                   │
                            ┌──────▼────────┐
                            │   Workers     │
                            │  (Global)     │
                            └───────────────┘
```

## Docker Compose Deployment

### Development Environment

```bash
# Start everything
docker compose up -d

# View logs
docker compose logs -f backend

# Stop everything
docker compose down
```

### Production Environment

```bash
# Set production environment
export NODE_ENV=production
export DB_PASSWORD=secure-password

# Start with resource limits
docker compose up -d

# Verify all services are running
docker compose ps
```

## Kubernetes Deployment

### Prerequisites

- Kubernetes 1.24+
- kubectl configured
- Helm 3.x (optional)
- cert-manager (for TLS)
- nginx-ingress (for ingress)

### Deployment Steps

**1. Create Namespace**

```bash
kubectl create namespace infrapilot
kubectl config set-context --current --namespace=infrapilot
```

**2. Deploy Infrastructure**

```bash
# Deploy PostgreSQL
kubectl apply -f k8s/postgres.yaml

# Deploy Redis
kubectl apply -f k8s/redis.yaml

# Wait for databases to be ready
kubectl wait --for=condition=ready pod -l app=infrapilot-postgres --timeout=300s
kubectl wait --for=condition=ready pod -l app=infrapilot-redis --timeout=300s
```

**3. Configure Secrets**

```bash
# Generate secrets
kubectl create secret generic infrapilot-secrets \
  --from-literal=DB_PASSWORD=$(openssl rand -base64 32) \
  --from-literal=JWT_SECRET=$(openssl rand -base64 64) \
  --from-literal=GRAFANA_PASSWORD=$(openssl rand -base64 32)
```

**4. Deploy Application**

```bash
# Deploy config
kubectl apply -f k8s/configmap.yaml

# Deploy backend and frontend
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml

# Deploy ingress (optional)
kubectl apply -f k8s/ingress.yaml
```

**5. Verify Deployment**

```bash
# Check pods
kubectl get pods -n infrapilot

# Check services
kubectl get services -n infrapilot

# Check ingress
kubectl get ingress -n infrapilot
```

**6. Access Application**

```bash
# Port forward for local access
kubectl port-forward service/infrapilot-frontend 8080:80

# Or use ingress if configured
# https://infrapilot.io
```

## Cloud Deployments

### AWS EKS

**Prerequisites:**
- AWS CLI configured
- eksctl installed

**1. Create EKS Cluster**

```bash
eksctl create cluster \
  --name infrapilot \
  --version 1.28 \
  --nodegroup-name workers \
  --node-type m5.xlarge \
  --nodes 3 \
  --nodes-min 2 \
  --nodes-max 5
```

**2. Install Addons**

```bash
# Install nginx ingress controller
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/aws/deploy.yaml

# Install cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.14.0/cert-manager.yaml
```

**3. Deploy Application**

```bash
# Update image references in k8s/deployment.yaml to use ECR
# Example: 123456789012.dkr.ecr.us-east-1.amazonaws.com/infrapilot/backend:latest

kubectl apply -f k8s/
```

### Google GKE

```bash
# Create cluster
gcloud container clusters create infrapilot \
  --num-nodes=3 \
  --machine-type=e2-medium \
  --zone=us-central1-a

# Get credentials
gcloud container clusters get-credentials infrapilot --zone=us-central1-a

# Deploy
kubectl apply -f k8s/
```

### Azure AKS

```bash
# Create cluster
az aks create \
  --resource-group infrapilot-rg \
  --name infrapilot \
  --node-count 3 \
  --enable-managed-identity

# Get credentials
az aks get-credentials --resource-group infrapilot-rg --name infrapilot

# Deploy
kubectl apply -f k8s/
```

## Helm Charts (Optional)

Create Helm chart for easier deployment:

```bash
helm create infrapilot
```

Structure:
```
infrapilot/
├── Chart.yaml
├── values.yaml
├── templates/
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── configmap.yaml
│   └── secret.yaml
└── values/
    ├── production.yaml
    └── staging.yaml
```

## Performance Tuning

### Database Optimization

```sql
-- Increase shared buffers
ALTER SYSTEM SET shared_buffers = '2GB';
ALTER SYSTEM SET effective_cache_size = '6GB';
ALTER SYSTEM SET maintenance_work_mem = '512MB';
ALTER SYSTEM SET checkpoint_completion_target = 0.9;
ALTER SYSTEM SET wal_buffers = '16MB';
ALTER SYSTEM SET default_statistics_target = 100;
ALTER SYSTEM SET random_page_cost = 1.1;
ALTER SYSTEM SET effective_io_concurrency = 200;
ALTER SYSTEM SET work_mem = '10486kB';
ALTER SYSTEM SET min_wal_size = '1GB';
ALTER SYSTEM SET max_wal_size = '4GB';

-- Restart PostgreSQL
```

### Redis Optimization

```conf
# redis.conf
maxmemory 2gb
maxmemory-policy allkeys-lru
save 900 1
save 300 10
save 60 10000
```

### Application Tuning

```env
# Backend
WORKER_PROCESSES=4
WORKER_QUEUE_SIZE=1000
CACHE_TTL=300
METRICS_BATCH_SIZE=100
```

## Monitoring Setup

### Prometheus Configuration

Create `prometheus.yml`:

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'infrapilot-backend'
    static_configs:
      - targets: ['backend:9090']
  
  - job_name: 'postgres'
    static_configs:
      - targets: ['postgres-exporter:9187']
  
  - job_name: 'redis'
    static_configs:
      - targets: ['redis-exporter:9121']
```

### Grafana Dashboards

Import pre-configured dashboards:
1. InfraPilot Overview
2. Backend Metrics
3. Database Performance
4. Queue Health
5. Agent Status

## Backup & Recovery

### Database Backup

```bash
# Automated backup script
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
docker exec infrapilot-postgres pg_dump -U postgres infrapilot_enterprise | gzip > backup_$DATE.sql.gz

# Upload to S3
aws s3 cp backup_$DATE.sql.gz s3://infrapilot-backups/
```

### Recovery

```bash
# Restore database
gunzip < backup_20250101_020000.sql.gz | docker exec -i infrapilot-postgres psql -U postgres infrapilot_enterprise
```

## Disaster Recovery

**Recovery Time Objective (RTO):** 15 minutes  
**Recovery Point Objective (RPO):** 1 hour

### Procedures

1. **Database Failure**: Promote read replica to primary
2. **Server Failure**: Kubernetes reschedules pods automatically
3. **Data Center Failure**: Deploy to secondary region

## Scaling Guidelines

### When to Scale

- **Backend**: CPU > 70% for 5 minutes
- **Database**: Connections > 80% of max_connections
- **Redis**: Memory > 80%
- **Queue**: Pending jobs > 1000

### Scaling Commands

```bash
# Kubernetes
kubectl scale deployment infrapilot-backend --replicas=5

# Docker Compose
docker compose up -d --scale backend=5
```

## Security Hardening

### Network Policies

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: infrapilot-backend-policy
spec:
  podSelector:
    matchLabels:
      app: infrapilot-backend
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: infrapilot-frontend
    ports:
    - protocol: TCP
      port: 8080
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: infrapilot-postgres
    ports:
    - protocol: TCP
      port: 5432
```

### RBAC

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: infrapilot-backend
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: infrapilot-backend-role
rules:
- apiGroups: [""]
  resources: ["configmaps", "secrets"]
  verbs: ["get", "list"]
```

## Maintenance

### Rolling Updates

```bash
# Kubernetes (zero downtime)
kubectl rollout restart deployment/infrapilot-backend

# Docker Compose
docker compose up -d --no-deps --build backend
```

### Certificate Renewal

```bash
# cert-manager auto-renews, but trigger manually if needed
kubectl certificate renew infrapilot-tls
```

### Log Rotation

```bash
# Docker log rotation
# Add to /etc/docker/daemon.json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m",
    "max-file": "3"
  }
}