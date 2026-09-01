# Changelog

All notable changes to InfraPilot Enterprise are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/).

## [1.1.0] - 2026-07-28

### Added - InfraPilot Enterprise v1.1.0 Production Release
- **Unified Enterprise Search (Sprint v1.1.9)**: Single search bar across all platform resources (Machines, Logs, Alerts, Incidents, Traces, Docker, Kubernetes, AI Insights, Reports, Users).
- **Natural Language AI Search**: Conversational query parser ("Which servers are unhealthy?", "Why is server01 slow?").
- **AI Center & Incident Management**: Root cause recommendations, incident timelines, and anomaly prediction records.
- **Application Performance Monitoring (APM) & Tracing**: Distributed tracing with OpenTelemetry, span waterfall flamegraphs, request throughput (RPS), and latency distribution.
- **Drag-and-Drop Dashboard Builder**: Custom tile layout builder with interactive chart widgets and layout persistence.
- **Observability Probes**: Self-monitoring health (`/health`), Kubernetes readiness (`/ready`), and Prometheus metrics exposition (`/metrics`).
- **Production Containerization**: Multi-stage `Dockerfile.backend`, `Dockerfile.frontend`, `Dockerfile.agent`, and complete Kubernetes manifests (`k8s/`).

### Improved
- Security hardening with Helmet security headers, CORS, rate limiting, and RBAC enforcement across all endpoints.
- DB connection pooling, indexing, and memory allocation optimization.
- Live WebSocket auto-reconnection and event streaming.

---

## [1.0.0-Enterprise] - 2026-07-27

### Added - Enterprise Platform Release (Sprints 11.0 - 12.0)
- **Multi-Tenancy & Org Management (Sprint 11.1)**: Full tenant data isolation, organization quotas, audit logging, role delegation, API keys.
- **Enterprise Analytics & SLA Suite (Sprint 11.2)**: Capacity forecasting, trend analysis, 99.99% SLA compliance tracking, automated PDF/Excel/CSV exports.
- **Hybrid Multi-Cloud Observability (Sprint 11.3)**: AWS (EC2, EKS, RDS, Lambda, S3), Azure (VMs, AKS, SQL), and GCP (Compute, GKE, Cloud SQL) integration, FinOps cost analysis, and CSPM posture findings.
- **AIOps Auto-Remediation Engine (Sprint 11.4)**: Rule-based auto-remediation, SSH/WinRM/K8s/Docker execution, 1-Click operational runbooks library, approval workflows, multi-channel notifications.
- **Platform Production Launch (Sprint 12.0)**: Self-monitoring health endpoints (`/api/v1/platform/health`), k6 performance load testing suite, 1-command backup/restore scripts (`scripts/backup.sh`), automated installer scripts (`scripts/install.sh`, `scripts/Install-InfraPilot.ps1`), and comprehensive User & Admin guides (`docs/`).

---

## [1.0.0] - 2026-07-08

### Added
- Initial release with telemetry collection, alert engine, JWT auth, and basic UI.
