# InfraPilot Enterprise - Project Synopsis & Report

---

## 1. Cover Page
* **Title**: InfraPilot Enterprise v1.0
* **Domain**: Cloud Infrastructure Orchestration & Real-time Metrics Observability
* **Date**: July 2026

---

## 2. Certificate
Certified that this project work is a record of genuine development carried out under industry best-practices and system validation cycles.

---

## 3. Acknowledgement
We express gratitude to the open-source community contributors, systems engineers, and technical mentors who provided frameworks and architectures to bring this platform to reality.

---

## 4. Abstract
InfraPilot Enterprise is an end-to-end multi-tenant infrastructure monitoring platform. By utilizing a lightweight Go-based host agent, a concurrent REST API engine in Go, PostgreSQL skipped-locks queue processing, Redis caches, and a React dark-theme client, it delivers real-time telemetry streaming, remote service execution controls, full-text log search, and automated self-healing action alerts.

---

## 5. Introduction
Observability is a cornerstone of modern reliability engineering. InfraPilot addresses the challenge of managing servers, containers, logs, and remediation tasks from a single, unified, performant panel.

---

## 6. Problem Statement
Existing infrastructure tools are either resource-heavy (Java-based agents consuming substantial host RAM) or lack native execution channels, forcing administrators to navigate between monitoring charts and terminal SSH connections to patch anomalies.

---

## 7. Objectives
* Design a lightweight Go agent consuming <15MB RAM.
* Build a scalable metrics ingestion pipeline handling high-frequency loads without database write bottle-necks.
* Secure operations via stateless JWT role-based access.
* Provide remote command execution, service controls, and container observability.

---

## 8. Existing System
* Heavy Java/Python monitoring agents.
* Static alert thresholds resulting in page fatigue.
* No built-in, audited command execution channels.

---

## 9. Proposed System
* Lightweight, single-binary Go agent.
* Asynchronous metrics enqueuing via PostgreSQL skipped-locks transactional loops.
* Native Windows/Linux execution tools and container manager.
* Audited single-click incident fixes recommendations.

---

## 10. System Architecture
Telemetry is gathered by the Go agent every 5s, sent over HTTP, stored in the `queued_metrics` table, processed by the database queue worker, evaluated for alert thresholds, and broadcasted to browser clients via Gorilla WebSockets.

---

## 11. Technology Stack
* **Backend**: Go 1.26+, Gin, GORM, Gorilla WebSockets.
* **Database/Cache**: PostgreSQL, Redis.
* **Frontend**: React, Vite, TailwindCSS.

---

## 12. Database Design
Entity relationship mappings include:
* `users` — Admin and Viewer identity schemas.
* `machines` — Monitored host details.
* `metrics` / `historical_metrics` — Telemetry statistics.
* `linux_dockers` — Container runtime statistics.
* `commands` / `audit_logs` — Audited execution logs.

---

## 13. API Design
* `POST /api/v1/auth/register` — Session user registration.
* `POST /api/v1/auth/login` — JWT credentials validation.
* `GET /api/v1/healthz` — System connectivity health parameters.
* `GET /api/v1/docs` — Developer documentation portal.
* `GET /api/v1/docker/containers/:id` — Remote container statistics.

---

## 14. Frontend Design
Builds upon React Single Page Application utilizing a sleek dark-themed console with glowing cyan/violet telemetry nodes and charts.

---

## 15. Backend Design
Engineered around Gin routes setups and asynchronous processing threads to handle telemetry ingestion separate from persistent disk write workloads.

---

## 16. Monitoring Agent
Compiled single binary querying local filesystems, CPU counters, network interface statistics, and running docker containers.

---

## 17. Testing
Verified using native testing frameworks inside `backend/internal/handlers/observability_test.go` checking APM and KPI endpoints.

---

## 18. Results
All unit and integration assertions pass cleanly, and real-time telemetry updates correctly.

---

## 19. Future Scope
* Kubernetes cluster monitoring.
* Multi-cloud AWS/Azure dashboards.
* Predictive baselines anomaly tracking using machine learning algorithms.

---

## 20. Conclusion
InfraPilot Enterprise bridges the gap between passive telemetry monitoring and active system orchestration.

---

## 21. References
* Go documentation and net/http packages.
* PostgreSQL transaction and locking strategies.
* Vite React asset bundle compiler specifications.
