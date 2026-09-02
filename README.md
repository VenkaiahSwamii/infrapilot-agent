# 🚀 InfraPilot Enterprise (v1.1.0 Production Release)

> **Production-ready unified enterprise observability, infrastructure telemetry monitoring, APM, and AIOps management platform**

[![Version](https://img.shields.io/badge/Release-v1.1.0--Production-success?logo=github)](CHANGELOG.md)
[![Go Version](https://img.shields.io/badge/Go-1.22-blue?logo=go)](https://go.dev/)
[![React](https://img.shields.io/badge/React-18-61dafb?logo=react)](https://react.dev/)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)
[![CI Status](https://img.shields.io/badge/CI-Passing-brightgreen?logo=githubactions)](.github/workflows/ci-cd.yml)
[![Kubernetes](https://img.shields.io/badge/K8s-Ready-326ce5?logo=kubernetes)](k8s/)
[![Security](https://img.shields.io/badge/Security-Passed-brightgreen?logo=trivy)](SECURITY.md)

InfraPilot Enterprise is a production-grade, self-hosted observability and infrastructure management platform. It aggregates live metrics, distributed traces, system logs, container workloads, Kubernetes clusters, incident alerts, and AI-driven root cause insights into a unified control plane.

---

## 🏛️ System Architecture

```
React Dashboard
        │
        ▼
Go REST API
        │
 ┌──────┼─────────┐
 │      │         │
 ▼      ▼         ▼
AI   PostgreSQL WebSocket
 │
 ▼
Agents
 │
 ├── Windows
 ├── Linux
 ├── Docker
 └── Kubernetes
```

---

## ✨ Enterprise Features Matrix (v1.1.0)

### 🔍 Unified Enterprise Search
- **Single Search Bar**: Global search across Machines, Logs, Alerts, Incidents, Docker, Kubernetes, APM, AI Insights, and Reports.
- **Natural Language AI Search**: Query platform state using conversational English (e.g. *"Which servers are unhealthy?"*, *"Why is server01 slow?"*).
- **Saved & Popular Searches**: Save frequent searches with 1-click execution and auto-complete suggestions.

### 🤖 Autonomous AI Operations & AI Center
- **Incident Analysis & Correlation**: Automated root-cause detection and remediations.
- **Predictive Analytics**: Anomaly detection and capacity forecasting before outages occur.

### ⚡ APM & Distributed Tracing
- **OpenTelemetry & Jaeger Tracing**: Real-time span visualizer, waterfall flamegraphs, and service map topology.
- **Application Metrics**: Request throughput (RPS), P95/P99 latency distribution, and error rates.

### 📊 Drag-and-Drop Dashboard Builder
- **Custom Tile Builder**: Create customized dashboards with line charts, gauge metrics, stat cards, and log widgets.
- **Layout Persistence**: Save, export, and load custom enterprise layouts.

### 🐳 Container & Kubernetes Observability
- **Docker Monitoring**: Live container status, memory limits, CPU metrics, and streaming logs.
- **Kubernetes Monitoring**: Pods, Deployments, Nodes, StatefulSets, Namespaces, and Cluster Events.

### 🛡️ Production Hardening & Security
- **RBAC Enforcement**: Fine-grained role permissions (SuperAdmin, Admin, Operator, Auditor, Viewer).
- **Security Headers & TLS 1.3**: Helmet security headers, CORS protection, rate limiting, and encrypted agent communications.
- **Immutable Audit Logging**: Full chain of audit events recorded for compliance.

---

## 🛠️ Technologies Used

| Layer | Technology |
|---|---|
| **Backend Core** | Go 1.22, Gin Web Framework, GORM, PostgreSQL 15 |
| **Frontend UI** | React 18, Vite, Lucide Icons, Recharts, Vanilla CSS |
| **Real-time Stream** | WebSockets (gorilla/websocket) |
| **Caching & Stream** | Redis 7, Redis Streams |
| **Observability** | Prometheus, OpenTelemetry, Jaeger |
| **Deployment** | Native Execution, Systemd, Windows Service, Kubernetes, Helm |

---

## ⚡ Quick Start & Native Execution (No Docker Required)

### 1-Click Dev Launcher

#### On Windows (PowerShell):
```powershell
.\start-dev.ps1
```

#### On Linux / macOS (Bash):
```bash
chmod +x ./start-dev.sh
./start-dev.sh
```

### Manual Execution

```bash
# 1. Start Backend API
cd backend
go run cmd/server/main.go

# 2. Start Frontend UI (in a new terminal)
cd frontend
npm install
npm run dev

# 3. Start InfraPilot Agent (optional, in a new terminal)
cd agent
go run cmd/agent/main.go
```

---

## 📚 Complete Production Documentation

- 📄 **[ARCHITECTURE.md](ARCHITECTURE.md)** - Technical design & component interactions
- ⚙️ **[INSTALL.md](INSTALL.md)** - Installation on Windows, Ubuntu, AWS, Azure
- 🚀 **[DEPLOYMENT.md](DEPLOYMENT.md)** - Production HA, Docker & Kubernetes guide
- 🛡️ **[SECURITY.md](SECURITY.md)** - RBAC, encryption, audit compliance
- 📡 **[API.md](API.md)** - REST & WebSocket API specification
- 🤖 **[AGENT.md](AGENT.md)** - Windows Service & Linux systemd Agent setup
- 🤝 **[CONTRIBUTING.md](CONTRIBUTING.md)** - Developer contribution guidelines
- 📜 **[CHANGELOG.md](CHANGELOG.md)** - Version release history

---

## 🌐 Observability & Self-Monitoring

InfraPilot monitors its own health out of the box:

- **Health Probe**: `GET http://localhost:8080/health`
- **Readiness Probe**: `GET http://localhost:8080/ready`
- **Prometheus Metrics**: `GET http://localhost:8080/metrics`

---

## 📜 License

Distributed under the [MIT License](LICENSE). Copyright © 2026 InfraPilot Enterprise.