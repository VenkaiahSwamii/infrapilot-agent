# InfraPilot Enterprise - Kubernetes Production Manifests

This directory contains production-grade Kubernetes manifests for deploying InfraPilot Enterprise.

## Directory Structure

```
deploy/kubernetes/
├── namespace.yaml              # Kubernetes namespace with pod security standards
├── configmap.yaml              # Non-sensitive configuration
├── secrets.yaml                # Sensitive credentials
├── DEPLOYMENT_ORDER.md         # Production deployment sequence
├── README.md                   # This file
│
├── backend/
│   ├── deployment.yaml         # Backend API deployment (3 replicas)
│   └── service.yaml            # ClusterIP service for backend
│
├── frontend/
│   ├── deployment.yaml         # Frontend deployment (2 replicas)
│   └── service.yaml            # ClusterIP service for frontend
│
├── postgres/
│   └── statefulset.yaml        # PostgreSQL StatefulSet with PVC (50Gi)
│
├── redis/
│   ├── deployment.yaml         # Redis StatefulSet (2 replicas, 5Gi storage)
│   └── service.yaml            # Headless service for Redis
│
├── qdrant/
│   └── statefulset.yaml        # Qdrant StatefulSet with PVC (100Gi)
│
├── ingress/
│   └── ingress.yaml            # NGINX ingress with TLS and WebSocket support
│
├── autoscaling/
│   └── hpa.yaml                # Horizontal Pod Autoscaler for backend (3-10 replicas)
│
├── network/
│   └── network-policies.yaml   # Network policies for zero-trust networking
│
└── resource-quotas/
    └── quota.yaml              # Resource quotas and limit ranges
```

## Components

### Backend (Go + Gin)
- **Replicas**: 3 (minimum)
- **Ports**: 8080 (HTTP), 9090 (metrics)
- **Resources**: 500m-2 CPU, 512Mi-2Gi memory
- **Health**: `/healthz` endpoint
- **Autoscaling**: CPU 70%, Memory 80%

### Frontend (React + NGINX)
- **Replicas**: 2
- **Ports**: 80 (HTTP), 443 (HTTPS)
- **Resources**: 200m-1 CPU, 256Mi-512Mi memory
- **Health**: `/healthz` endpoint

### PostgreSQL (v15)
- **Type**: StatefulSet
- **Replicas**: 1 (primary)
- **Storage**: 50Gi persistent volume
- **Port**: 5432
- **Health**: `pg_isready` check

### Redis (v7)
- **Type**: StatefulSet
- **Replicas**: 2 (Sentinel-ready)
- **Storage**: 5Gi persistent volume
- **Port**: 6379
- **Health**: `redis-cli ping` check

### Qdrant
- **Type**: StatefulSet
- **Replicas**: 1
- **Storage**: 100Gi persistent volume
- **Ports**: 6333 (HTTP), 6334 (gRPC)
- **Health**: `/healthz` endpoint

## Storage Requirements

| Component   | Storage  | Access Mode      |
|-------------|----------|------------------|
| PostgreSQL  | 50Gi     | ReadWriteOnce    |
| Redis       | 5Gi      | ReadWriteOnce    |
| Qdrant      | 100Gi    | ReadWriteOnce    |
| Grafana     | 10Gi     | ReadWriteOnce    |
| Loki        | 50Gi     | ReadWriteOnce    |
| Prometheus  | 100Gi    | ReadWriteOnce    |

## Networking

### Ingress Routes
- `/` → Frontend (React app)
- `/api` → Backend API
- `/ws` → WebSocket service

### TLS/HTTPS
- Managed by cert-manager
- ClusterIssuer: `letsencrypt-prod`
- Secret: `infrapilot-tls`

### Network Policies
- Default deny-all ingress/egress
- Frontend → Backend (port 8080)
- Backend → PostgreSQL (port 5432)
- Backend → Redis (port 6379)
- Backend → Qdrant (port 6333)
- DNS egress allowed for all pods

## Resource Quotas

### CPU
- Requests: 8 cores
- Limits: 16 cores

### Memory
- Requests: 16Gi
- Limits: 32Gi

### Other
- Pods: 20
- PVCs: 10
- Services: 10

## Security

- All containers run as non-root
- Pod Security Standards: restricted
- Secrets stored in Kubernetes Secrets
- Sensitive values injected via environment variables
- No hardcoded credentials

## Monitoring

Prometheus scrapes:
- Backend: `infrapilot-backend:9090/metrics`
- PostgreSQL: `infrapilot-postgres-exporter:9187`
- Redis: `infrapilot-redis-exporter:9121`
- Qdrant: `infrapilot-qdrant:6333/metrics`

## Deployment

See [DEPLOYMENT_ORDER.md](./DEPLOYMENT_ORDER.md) for step-by-step deployment instructions.

## Validation

```bash
# Validate all resources
kubectl get pods -n infrapilot
kubectl get svc -n infrapilot
kubectl get ingress -n infrapilot
kubectl get pvc -n infrapilot
kubectl get hpa -n infrapilot

# Check for issues
kubectl get events -n infrapilot --sort-by='.lastTimestamp'
```

## Troubleshooting

### Pods Pending
```bash
kubectl describe pod <pod-name> -n infrapilot
```
Check if PVC is bound and node has resources.

### CrashLoopBackOff
```bash
kubectl logs <pod-name> -n infrapilot --previous
kubectl describe pod <pod-name> -n infrapilot
```

### Ingress No Address
```bash
kubectl get pods -n ingress-nginx
kubectl describe ingress infrapilot-ingress -n infrapilot
```

## Production Readiness

- ✅ High availability (multiple replicas)
- ✅ Rolling updates (zero-downtime)
- ✅ Health probes (liveness, readiness, startup)
- ✅ Resource limits and requests
- ✅ Persistent storage for stateful services
- ✅ TLS/HTTPS encryption
- ✅ WebSocket support
- ✅ Network segmentation
- ✅ Horizontal autoscaling
- ✅ Pod security standards
- ✅ Resource quotas
- ✅ Secrets management