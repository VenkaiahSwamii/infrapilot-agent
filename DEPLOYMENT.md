# 🚀 Production Deployment & High Availability Guide

## Overview

InfraPilot Enterprise supports High Availability (HA) deployments across AWS, Azure, GCP, and self-hosted Kubernetes or Docker Swarm clusters.

---

## Production Deployment Architecture

```
                      AWS ALB / Nginx Load Balancer
                                   │
             ┌─────────────────────┴─────────────────────┐
             ▼                                           ▼
  Backend Node 1 (Go API)                     Backend Node 2 (Go API)
             │                                           │
             └─────────────────────┬─────────────────────┘
                                   │
                   ┌───────────────┴───────────────┐
                   ▼                               ▼
        RDS PostgreSQL (Multi-AZ)       ElastiCache Redis Cluster
```

---

## AWS Deployment Steps (EKS / RDS)

1. Provision PostgreSQL RDS (Multi-AZ) and ElastiCache Redis Cluster.
2. Update Kubernetes secrets in `k8s/secret.yaml` with DB DSN credentials.
3. Apply ingress and deployment manifests:
   ```bash
   kubectl apply -f k8s/secret.yaml
   kubectl apply -f k8s/configmap.yaml
   kubectl apply -f k8s/deployment.yaml
   kubectl apply -f k8s/service.yaml
   kubectl apply -f k8s/ingress.yaml
   kubectl apply -f k8s/hpa.yaml
   ```

---

## Environment Variables Configuration

| Parameter | Default | Purpose |
|---|---|---|
| `PORT` | `8080` | Backend API HTTP port |
| `DB_HOST` | `postgres` | PostgreSQL hostname |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_NAME` | `infrapilot_enterprise` | Database name |
| `JWT_SECRET` | *(Random)* | Secret key for signing JWT tokens |
| `REDIS_HOST` | `redis` | Redis server hostname |
