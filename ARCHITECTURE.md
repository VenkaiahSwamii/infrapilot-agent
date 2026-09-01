# 🏗️ InfraPilot Enterprise Architecture Guide (v1.1.0)

## Overview

InfraPilot Enterprise is built on a decoupled, asynchronous micro-architecture designed for low telemetry latency, horizontal scale, and multi-tenant isolation.

```
                     Search & Control Plane
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│                      Go Backend API                         │
│   (Gin Router, JWT Auth, RBAC, Rate Limiting, CORS, Gzip)   │
└───────┬──────────────────────┬──────────────────────┬───────┘
        │                      │                      │
        ▼                      ▼                      ▼
 ┌──────────────┐       ┌──────────────┐       ┌──────────────┐
 │ PostgreSQL   │       │ Redis Stream │       │ WebSockets   │
 │ (Data Store) │       │ & PubSub     │       │ Broadcast    │
 └──────────────┘       └──────────────┘       └──────┬───────┘
                                                      │
                                                      ▼
                                              React Dashboard
```

## Key Architectural Components

### 1. Cross-Platform Agents (`/agent`)
- Headless daemon running on Linux (systemd) or Windows (Windows Service).
- Collects OS system metrics, process states, Docker container telemetry, and Kubernetes cluster health.
- Connects securely via TLS 1.3 to Go Backend API using encrypted API Keys and JWT tokens.
- Offline Queue: Buffer metrics locally during network partitioning and replays upon reconnection.

### 2. High-Performance Go Backend (`/backend`)
- **Web Layer**: Built on Gin web framework with middleware for Request ID tracing, Helmet security headers, CORS, Gzip compression, and Rate limiting.
- **Unified Search Engine (`internal/search`)**: In-memory query parser, scoring ranking engine, full-text tsvector search, and AI natural language query synthesis.
- **Observability Probes**: Liveness (`/health`), readiness (`/ready`), and Prometheus exposition (`/metrics`).

### 3. Database Layer (PostgreSQL & Redis)
- **PostgreSQL 15**: Primary relational data store for servers, metrics, alerts, incidents, traces, logs, search index, and organization audit logs.
- **Redis 7**: Distributed cache and Redis Streams queue for asynchronous event processing.

### 4. Real-time Communication (WebSocket Hub)
- Concurrent WebSocket Hub (`internal/websocket`) managing client sessions.
- Pushes real-time metric streams, search index updates, alert triggers, and incident status changes directly to active browser sessions.

### 5. Frontend Control Plane (`/frontend`)
- Built on React 18, Vite, and Vanilla CSS.
- Features: Dashboard Builder, APM & Distributed Tracing, AI Center, Alerting, Logs, Cloud Workloads, and Unified Enterprise Search.
