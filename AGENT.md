# 🤖 InfraPilot Agent Deployment & Operations

## Overview

The InfraPilot Agent is a lightweight cross-platform binary that gathers hardware metrics, process telemetry, Docker container stats, and Kubernetes pod metrics.

---

## Supported Operating Systems
- **Linux**: RHEL 8+, Ubuntu 20.04+, Debian 11+, Alpine 3.16+
- **Windows**: Windows Server 2019+, Windows 10/11
- **Containers**: Docker 20.10+, Kubernetes 1.24+

---

## Key Features
- **Low Overhead**: < 1.5% CPU utilization and < 35MB RAM footprint.
- **Offline Buffering**: Stores metrics in a local SQLite/memory ring buffer if network connectivity drops.
- **Auto Reconnect**: Exponential backoff reconnect logic for WebSockets and HTTP APIs.
- **API Key Rotation**: Periodically verifies and rotates tenant API keys over TLS 1.3.
