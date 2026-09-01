# ⚙️ InfraPilot Enterprise Installation Guide

This guide covers deployment and installation across Windows, Linux (Ubuntu/Debian), Docker, and Kubernetes environments.

---

## Prerequisites

- **Go**: 1.22+ (for building backend from source)
- **Node.js**: 20+ & npm (for building frontend)
- **PostgreSQL**: 15+
- **Redis**: 7+
- **Docker & Docker Compose** (Optional, recommended)

---

## 1. Quick Installation via Docker Compose

```bash
git clone https://github.com/venkaiswami/infrapilot-enterprise.git
cd infrapilot-enterprise

# Launch all infrastructure services
docker compose up -d

# Check service health
curl http://localhost:8080/health
```

---

## 2. Linux Agent Installation (systemd)

```bash
# Download & compile agent
cd agent
go build -o infrapilot-agent ./cmd/agent

# Register systemd service
sudo cp scripts/infrapilot-agent.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable infrapilot-agent
sudo systemctl start infrapilot-agent
```

---

## 3. Windows Agent Installation (Service)

```powershell
# Run PowerShell as Administrator
cd agent
go build -o infrapilot-agent.exe .\cmd\agent

# Install Windows Service
New-Service -Name "InfraPilotAgent" -BinaryPathName "$PWD\infrapilot-agent.exe" -StartupType Automatic
Start-Service -Name "InfraPilotAgent"
```

---

## 4. Kubernetes Deployment

```bash
# Deploy all Kubernetes objects into 'infrapilot' namespace
kubectl apply -f k8s/

# Verify status
kubectl get pods -n infrapilot
```
