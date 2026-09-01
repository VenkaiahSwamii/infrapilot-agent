# InfraPilot Enterprise — Demo Video Script

**Target length:** 10–15 minutes
**Audience:** Recruiters, hiring managers, engineers evaluating the project
**Tone:** Confident, concise, technical — narrate what you're doing and *why* it matters.

---

## 0:00 – 0:45 | Intro

> "Hi, I'm [Your Name], and this is InfraPilot Enterprise — a production-style infrastructure monitoring platform I built from scratch using Go and React. It monitors servers, Docker containers, and Kubernetes clusters in real time, with remote terminal access, file management, alerting, and enterprise reporting. Let me walk you through it."

Show: Landing page / website (`website/index.html`) or architecture diagram for 5 seconds.

## 0:45 – 2:00 | Login & Authentication

- Show the login page (`frontend/src/features/auth/LoginPage.jsx`)
- Log in with an admin account
- Mention: JWT access/refresh tokens, RBAC (SuperAdmin, Admin, Operator, Viewer), rate limiting, bcrypt password hashing

> "Authentication is JWT-based with automatic refresh token rotation, and every endpoint is protected by role-based access control with four tiers."

## 2:00 – 4:00 | Dashboard Overview

- Show `EnterpriseDashboard.jsx` — KPI cards (machines, CPU, memory, disk, alerts, containers, pods)
- Point out live-updating charts
- Explain the WebSocket connection powering real-time updates without polling

> "This dashboard updates in real time over WebSocket — no polling. Every KPI card and chart reflects live data streamed from the agents through the backend's event bus."

## 4:00 – 5:30 | Agent Enrollment & Live Metrics

- Show enrolling a new machine (enrollment token flow)
- Show the agent starting up in a terminal (`go run cmd/agent/main.go --server ... --token ...`)
- Switch to the machine detail page, show the Overview and Metrics tabs populating live

> "The agent is a lightweight, cross-platform Go binary. It enrolls with a token, then starts streaming system metrics — CPU, memory, disk, network, and processes — within seconds."

## 5:30 – 7:00 | Docker Monitoring

- Open the Docker tab on a machine running containers
- Show container list, CPU/memory per container, start/stop/restart actions

> "InfraPilot also monitors Docker directly through the Engine API — container status, resource usage, and lifecycle management, all from the dashboard."

## 7:00 – 8:30 | Kubernetes Monitoring

- Open the Kubernetes tab
- Show pods, nodes, and cluster health

> "For teams running Kubernetes, the agent's K8s plugin talks to the cluster API to surface pod and node status directly alongside your other infrastructure."

## 8:30 – 10:00 | Remote Terminal & File Manager

- Open Terminal tab, run a couple of commands (e.g., `top`, `df -h`)
- Open Files tab, browse a directory, download a file

> "For day-to-day operations, there's a full remote terminal built on xterm.js and WebSocket, and a file manager for browsing, uploading, and downloading files — both with full audit logging."

## 10:00 – 11:30 | Alerts & Reports

- Show the Alerts page — create/acknowledge/resolve a threshold alert
- Show the Reports page — generate a PDF/Excel report

> "Alerts are threshold-based with severity levels, and can be acknowledged and resolved with a full history. Reports can be generated on demand or scheduled, in PDF, CSV, or Excel format."

## 11:30 – 13:00 | Architecture Overview

- Show `docs/architecture.md` diagram or draw it live
- Walk through: Agent → Backend (Gin) → Event Bus → PostgreSQL / Redis / WebSocket Hub → React frontend
- Mention observability: Prometheus metrics, Grafana dashboards, OpenTelemetry tracing

> "Under the hood, agents push metrics to a Go backend built on Gin. An internal event bus fans events out to subscribers — the database, the WebSocket hub, the alert engine, and an audit logger. Everything is observable via Prometheus, Grafana, and distributed tracing with OpenTelemetry."

## 13:00 – 14:00 | Deployment

- Show `docker-compose up -d` bringing the stack up
- Show `kubectl apply -f k8s/` and `kubectl get pods -n infrapilot`

> "It deploys with a single `docker-compose up`, or to Kubernetes with the manifests in the `k8s/` folder. There's also a one-click cloud deployment script for a fresh Ubuntu VM."

## 14:00 – 15:00 | Closing

> "That's InfraPilot Enterprise — real-time infrastructure monitoring, Docker and Kubernetes visibility, remote management, and enterprise reporting, built with Go, React, PostgreSQL, and Redis. The full source is on GitHub — link in the description. Thanks for watching!"

Show: GitHub repo URL and portfolio site URL as an end card.

---

## Shot List Summary

| # | Segment | Duration |
|---|---------|----------|
| 1 | Intro | 45s |
| 2 | Login | 1m15s |
| 3 | Dashboard | 2m |
| 4 | Agent enrollment & live metrics | 1m30s |
| 5 | Docker monitoring | 1m30s |
| 6 | Kubernetes monitoring | 1m30s |
| 7 | Terminal & file manager | 1m30s |
| 8 | Alerts & reports | 1m30s |
| 9 | Architecture overview | 1m30s |
| 10 | Deployment | 1m |
| 11 | Closing | 1m |
