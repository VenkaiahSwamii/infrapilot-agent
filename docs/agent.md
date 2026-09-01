# InfraPilot Enterprise Agent Guide

## Overview

The InfraPilot Agent is a lightweight, high-performance monitoring agent written in Go. It runs on target systems and collects metrics, logs, and events to send to the InfraPilot backend.

## Features

- **System Monitoring:** CPU, memory, disk, network metrics
- **Process Monitoring:** Process list, resource usage, lifecycle events
- **Docker Monitoring:** Container status, resource usage, events
- **Kubernetes Monitoring:** Cluster health, pod status, node metrics
- **Log Collection:** System logs, application logs, log tailing
- **Remote Access:** Secure shell (SSH-like) terminal access
- **File Management:** Remote file browsing and operations
- **Plugin Architecture:** Extensible plugin system

## Architecture

```
┌─────────────────────────────────────┐
│         InfraPilot Agent            │
│  ┌───────────────────────────────┐  │
│  │      Core Engine              │  │
│  │  - Plugin Manager             │  │
│  │  - Connection Handler         │  │
│  │  - Event Collector            │  │
│  └───────────────────────────────┘  │
│              │                      │
│  ┌───────────┼───────────────┐     │
│  │           │               │     │
│  ▼           ▼               ▼     │
│ ┌─────┐ ┌──────────┐ ┌─────────┐  │
│ │Linux│ │  Docker  │ │   K8s   │  │
│ │Plugin│ │ Plugin   │ │ Plugin  │  │
│ └─────┘ └──────────┘ └─────────┘  │
│                                      │
│  ┌──────────────────────────────┐   │
│  │   Plugin Manager             │   │
│  │   - Plugin Registry          │   │
│  │   - Lifecycle Management     │   │
│  │   - Event Dispatch           │   │
│  └──────────────────────────────┘   │
└─────────────────────────────────────┘
```

## Installation

### Linux (x86_64/amd64)

```bash
# Download latest agent
curl -fsSL https://infrapilot.io/agent/latest/linux-amd64 -o infrapilot-agent
chmod +x infrapilot-agent

# Move to PATH
sudo mv infrapilot-agent /usr/local/bin/
```

### Linux (ARM64)

```bash
curl -fsSL https://infrapilot.io/agent/latest/linux-arm64 -o infrapilot-agent
chmod +x infrapilot-agent
sudo mv infrapilot-agent /usr/local/bin/
```

### Windows

**Option 1: Chocolatey**
```powershell
choco install infrapilot-agent
```

**Option 2: Direct Download**
```powershell
# Download
Invoke-WebRequest -Uri "https://infrapilot.io/agent/latest/windows-amd64.exe" -OutFile "infrapilot-agent.exe"

# Install to Program Files
New-Item -ItemType Directory -Force -Path "C:\Program Files\InfraPilot\agent"
Move-Item infrapilot-agent.exe "C:\Program Files\InfraPilot\agent\"
```

### Docker

```bash
docker run -d \
  --name infrapilot-agent \
  --restart unless-stopped \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /sys:/sys:ro \
  infrapilot/agent:latest \
  -backend http://infrapilot-backend:8080
```

## Configuration

### Configuration File

Location: `/etc/infrapilot/config.yaml` or `~/.infrapilot/config.yaml`

```yaml
server:
  url: "http://infrapilot.io:8080"
  ws_url: "ws://infrapilot.io:8080/ws"
  timeout: 30s
  retry_interval: 5s

agent:
  id: ""  # Auto-generated if empty
  name: ""  # Auto-generated from hostname
  tags:
    - production
    - us-east-1
  update_interval: 10s

auth:
  token: ""  # Enrollment token
  certificate_path: "/etc/infrapilot/cert.pem"

plugins:
  linux:
    enabled: true
    collect_cpu: true
    collect_memory: true
    collect_disk: true
    collect_network: true
    collect_processes: true
    interval: 10s

  docker:
    enabled: true
    socket_path: "/var/run/docker.sock"
    collect_containers: true
    collect_images: false
    interval: 15s

  kubernetes:
    enabled: false
    kubeconfig: "~/.kube/config"
    in_cluster: false
    collect_nodes: true
    collect_pods: true
    collect_deployments: true
    interval: 30s

  logs:
    enabled: true
    paths:
      - /var/log/syslog
      - /var/log/messages
      - /var/log/app/*.log
    max_lines: 1000
    tail: true

  security:
    enabled: true
    collect_users: true
    collect_ssh_logins: true
    collect_sudo_usage: true
    interval: 60s

terminal:
  enabled: true
  shells:
    - /bin/bash
    - /bin/sh
  max_sessions: 5
  idle_timeout: 3600s

files:
  enabled: true
  allowed_paths:
    - /var/log
    - /home
  blocked_paths:
    - /etc/shadow
    - /etc/passwd
  max_file_size: 100MB
```

### Environment Variables

```bash
INFrapilot_BACKEND_URL=http://infrapilot.io:8080
INFrapilot_BACKEND_WS_URL=ws://infrapilot.io:8080/ws
INFrapilot_CONFIG_PATH=/etc/infrapilot/config.yaml
INFrapilot_LOG_LEVEL=info  # debug, info, warn, error
INFrapilot_LOG_FILE=/var/log/infrapilot/agent.log
```

## Usage

### Start Agent

```bash
# Start in foreground
infrapilot-agent -backend http://infrapilot.io:8080

# Start with config file
infrapilot-agent -config /etc/infrapilot/config.yaml
```

### Enrollment

When first connecting to the backend:

```bash
# Option 1: Use token (recommended)
infrapilot-agent enroll --token YOUR_ENROLLMENT_TOKEN

# Option 2: Interactive enrollment
infrapilot-agent enroll
# Follow prompts to:
# 1. Enter backend URL
# 2. Login with admin credentials
# 3. Set agent name and tags
```

The enrollment process:
1. Authenticates with the backend
2. Registers the agent
3. Downloads a certificate
4. Saves configuration

### Service Installation

#### Linux (systemd)

```bash
# Create service file
sudo tee /etc/systemd/system/infrapilot-agent.service << EOF
[Unit]
Description=InfraPilot Agent
After=network.target

[Service]
Type=simple
User=infrapilot
WorkingDirectory=/opt/infrapilot
ExecStart=/usr/local/bin/infrapilot-agent -config /etc/infrapilot/config.yaml
Restart=always
RestartSec=10
Environment="INFrapilot_LOG_LEVEL=info"

[Install]
WantedBy=multi-user.target
EOF

# Create user
sudo useradd -r -s /bin/false infrapilot

# Set permissions
sudo mkdir -p /etc/infrapilot /var/log/infrapilot
sudo chown -R infrapilot:infrapilot /etc/infrapilot /var/log/infrapilot

# Enable and start
sudo systemctl daemon-reload
sudo systemctl enable infrapilot-agent
sudo systemctl start infrapilot-agent
sudo systemctl status infrapilot-agent
```

#### Windows (Service)

```powershell
# Install service using NSSM
nssm install InfraPilotAgent "C:\Program Files\InfraPilot\agent\infrapilot-agent.exe"
nssm set InfraPilotAgent AppParameters "-config C:\Program Files\InfraPilot\agent\config.yaml"
nssm set InfraPilotAgent DisplayName "InfraPilot Agent"
nssm set InfraPilotAgent Start SERVICE_AUTO_START

# Start service
nssm start InfraPilotAgent
```

### Updating

```bash
# Linux
infrapilot-agent update

# Or manually
wget https://infrapilot.io/agent/latest/linux-amd64 -O infrapilot-agent
sudo systemctl restart infrapilot-agent

# Docker
docker pull infrapilot/agent:latest
docker compose up -d agent
```

## Plugin System

### Creating Custom Plugins

```go
package main

import (
    "context"
    "time"
    "github.com/infrapilot/agent/pkg/plugin"
)

type CustomPlugin struct {
    name        string
    interval    time.Duration
    running     bool
}

func NewCustomPlugin() plugin.Plugin {
    return &CustomPlugin{
        name:     "custom",
        interval: 10 * time.Second,
    }
}

func (p *CustomPlugin) Name() string {
    return p.name
}

func (p *CustomPlugin) Interval() time.Duration {
    return p.interval
}

func (p *CustomPlugin) Collect(ctx context.Context) (*plugin.Metrics, error) {
    metrics := &plugin.Metrics{}
    
    // Collect your custom metrics
    metrics.Custom.Gauge("my_metric", 42.0)
    
    return metrics, nil
}

func (p *CustomPlugin) Start(ctx context.Context) error {
    p.running = true
    return nil
}

func (p *CustomPlugin) Stop() error {
    p.running = false
    return nil
}
```

### Plugin Registration

```go
// Register plugin in agent
plugin.Register("custom", NewCustomPlugin)
```

## Security

### Authentication

Agents authenticate using:
1. **Enrollment Token:** One-time secret from admin
2. **TLS Certificate:** Auto-issued during enrollment
3. **JWT Token:** Derived from certificate

### Network Security

- All connections use **TLS 1.3**
- Certificate pinning for backend
- MTLs for agent-to-backend communication
- WebSocket over WSS

### Permissions

The agent requires these Linux capabilities:
- `CAP_SYS_PTRACE` (for process monitoring)
- `CAP_DAC_READ_SEARCH` (for file access)
- `CAP_NET_ADMIN` (for network monitoring)

### Filesystem Access

The agent respects these rules:
- Cannot access files outside `allowed_paths`
- Automatically blocks sensitive files (`/etc/shadow`, etc.)
- Logs all file access attempts
- Read-only by default

## Monitoring

### Agent Status

Check agent status:
```bash
systemctl status infrapilot-agent
journalctl -u infrapilot-agent -f
```

### Metrics

Agent exposes metrics via Prometheus:
```
http://localhost:9090/metrics
```

Key metrics:
- `agent_collect_duration_seconds`
- `agent_events_sent_total`
- `agent_connection_status`
- `agent_plugin_errors_total`

### Logs

Log levels: `debug`, `info`, `warn`, `error`

```bash
# View logs
journalctl -u infrapilot-agent -f

# Change log level
sudo systemctl edit infrapilot-agent
# Add: Environment="INFrapilot_LOG_LEVEL=debug"
sudo systemctl restart infrapilot-agent
```

## Troubleshooting

### Agent won't enroll

```bash
# Check connectivity
curl http://infrapilot.io:8080/healthz

# Check logs
journalctl -u infrapilot-agent -n 100

# Verify token
infrapilot-agent verify-token YOUR_TOKEN
```

### High resource usage

```bash
# Reduce collection intervals in config.yaml
# Disable unused plugins
plugins:
  kubernetes:
    enabled: false

# Limit process collection
plugins:
  linux:
    max_processes: 100
```

### Certificate expired

```bash
# Renew certificate
infrapilot-agent renew-cert

# Or re-enroll
infrapilot-agent enroll --force
```

## Best Practices

1. **Run as non-root user** with minimal capabilities
2. **Enable only required plugins**
3. **Use tags** to organize agents (env, location, role)
4. **Regular updates** for security patches
5. **Monitor agent itself** (health, resource usage)
6. **Secure config files** (600 permissions)
7. **Use firewall** to restrict agent outbound connections

## API Reference

### Agent Endpoints

The agent exposes these internal endpoints:

- `GET /healthz` - Liveness check
- `GET /metrics` - Prometheus metrics
- `POST /api/v1/enroll` - Enroll agent
- `POST /api/v1/heartbeat` - Send heartbeat
- `WS /ws` - WebSocket connection