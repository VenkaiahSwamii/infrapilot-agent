# Docker Assets

InfraPilot Enterprise ships a `Dockerfile` for each service, colocated with its source so build contexts stay small and CI caching stays effective. This folder documents how they fit together — it does not duplicate the Dockerfiles themselves.

## Service Images

| Service | Dockerfile | Build Context | Base Image | Final Size |
|---------|-----------|----------------|-------------|------------|
| Backend  | [`backend/Dockerfile.backend`](../backend/Dockerfile.backend) | `./backend` | `golang:1.24-alpine` → `alpine:3.19` | ~25 MB |
| Frontend | [`frontend/Dockerfile.frontend`](../frontend/Dockerfile.frontend) | `./frontend` | `node:22-alpine` → `nginx:1.25-alpine` | ~40 MB |
| Agent    | [`agent/Dockerfile.agent`](../agent/Dockerfile.agent) | `./agent` | `golang:1.24-alpine` → `alpine:3.19` | ~18 MB |

All three use **multi-stage builds** — a build stage compiles the Go binary or Vite bundle, and a minimal runtime stage copies only the final artifact. This keeps production images small, reduces the attack surface, and speeds up `docker pull` in Kubernetes.

## Local Development

```bash
# Build & start everything (backend, frontend, postgres, redis)
docker-compose up -d

# Include observability stack (Prometheus, Grafana, Jaeger)
docker-compose --profile observability up -d

# Include a monitoring agent instance
docker-compose --profile agent up -d
```

See the root [`docker-compose.yml`](../docker-compose.yml) for the full service definitions, health checks, and resource limits.

## Building Images Manually

```bash
# Backend
docker build -f backend/Dockerfile.backend -t infrapilot/backend:latest ./backend

# Frontend
docker build -f frontend/Dockerfile.frontend -t infrapilot/frontend:latest ./frontend

# Agent
docker build -f agent/Dockerfile.agent -t infrapilot/agent:latest ./agent
```

## Registry & CI/CD

Images are built and pushed automatically by [`.github/workflows/ci.yml`](../.github/workflows/ci.yml) (on every push to `main`) and [`.github/workflows/cd.yml`](../.github/workflows/cd.yml) (on tagged releases), publishing to GitHub Container Registry (`ghcr.io`):

```
ghcr.io/venkaiswami/infrapilot-enterprise/backend:latest
ghcr.io/venkaiswami/infrapilot-enterprise/frontend:latest
ghcr.io/venkaiswami/infrapilot-enterprise/agent:latest
```

## Kubernetes

For production Kubernetes deployment, see the [`k8s/`](../k8s/) directory, which references these same images. Update `k8s/deployment.yaml` with the image tag you want to roll out, then:

```bash
kubectl apply -f k8s/
kubectl rollout status deployment/infrapilot-backend -n infrapilot
```
