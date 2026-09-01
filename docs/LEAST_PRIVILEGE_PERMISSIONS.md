# InfraPilot Enterprise - Native Non-Root & Least-Privilege Guide

## 1. Running Completely Without Docker (Bare-Metal)

InfraPilot Enterprise can be run 100% natively without Docker containers.

### Starting the Platform Natively

```bash
# 1. Run the native startup script
./scripts/install.sh
```

Or start components individually:

#### Backend API Server:
```bash
cd backend
go build -o infrapilot-server ./cmd/server/main.go
./infrapilot-server
```
*(Listens on port `8080`)*

#### Frontend Web Dashboard:
```bash
cd frontend
npm install
npm run dev -- --port 5173 --host 0.0.0.0
```
*(Accessible at `http://localhost:5173`)*

---

## 2. Non-Root / Standard User Agent Installation (No Sudo)

The InfraPilot monitoring agent is designed to run directly as a standard, unprivileged user.

### 1-Line Non-Root Installer:
```bash
curl -fsSL http://<BACKEND_HOST>:8080/downloads/install.sh | bash -s -- --server "http://<BACKEND_HOST>:8080" --token "<ENROLLMENT_TOKEN>"
```

### What happens when run as a non-root user:
- Installs the standalone native binary to `$HOME/.infrapilot/infrapilot-agent`
- Automatically registers a **systemd user service** (`systemctl --user enable --now infrapilot-agent`) or runs as a background daemon
- Collects all system telemetry without requiring root:
  - **CPU, Memory, Swap usage**
  - **Disk space and I/O statistics**
  - **Network bandwidth and socket latency**
  - **Running processes and resource consumption**
  - **Application logs and health metrics**

---

## 3. Granular Least-Privilege Permissions (Optional)

If you wish to allow the agent to perform specific administrative tasks (such as restarting services or checking disk SMART health), you **do NOT** need to grant full sudo.

Grant permissions **only for the specific commands needed** by creating `/etc/sudoers.d/infrapilot-agent`:

```sudoers
# /etc/sudoers.d/infrapilot-agent
# Grant execution rights ONLY to specific whitelisted commands

# 1. Service restart/status permissions only
your_username ALL=(ALL) NOPASSWD: /bin/systemctl restart nginx
your_username ALL=(ALL) NOPASSWD: /bin/systemctl restart docker
your_username ALL=(ALL) NOPASSWD: /bin/systemctl status *

# 2. Hardware SMART disk diagnostics only
your_username ALL=(ALL) NOPASSWD: /usr/sbin/smartctl -H *
your_username ALL=(ALL) NOPASSWD: /usr/sbin/smartctl -A *

# 3. Network packet capture capability (Alternative to sudo)
# sudo setcap cap_net_raw,cap_dac_read_search+ep ~/.infrapilot/infrapilot-agent
```

To apply:
```bash
sudo chmod 440 /etc/sudoers.d/infrapilot-agent
```
