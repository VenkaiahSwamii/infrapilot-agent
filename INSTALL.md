# ⚙️ InfraPilot Enterprise Installation Guide

This guide covers deployment and installation across Windows and Linux (Ubuntu/Debian) environments natively without Docker.

---

## Prerequisites

- **Go**: 1.22+ (for running/building backend & agent)
- **Node.js**: 20+ & npm (for running/building frontend)
- **PostgreSQL**: 15+ (or SQLite for local dev)
- **Redis**: 7+ (optional for stream processing)

---

## 1. Quick Native Start (No Docker Required)

### 1-Click Launchers

**On Windows (PowerShell):**
```powershell
.\start-dev.ps1
```

**On Linux / macOS (Bash):**
```bash
chmod +x ./start-dev.sh
./start-dev.sh
```

### Manual Service Launch

```bash
# 1. Launch Backend API
cd backend
go run cmd/server/main.go

# 2. Launch Frontend UI
cd frontend
npm install
npm run dev

# 3. Launch Agent
cd agent
go run cmd/agent/main.go
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
