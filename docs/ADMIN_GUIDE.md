# 🛠️ InfraPilot Enterprise v1.0 Administrator Guide

This guide covers deployment, cluster scaling, multi-tenant RBAC configuration, backup & disaster recovery procedures, and security hardening for system administrators.

---

## 💻 Quickstart Production Installation

### Linux / macOS
```bash
curl -fsSL https://install.infrapilot.io | bash
```

### Windows PowerShell
```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; .\scripts\Install-InfraPilot.ps1
```

### Docker Compose
```bash
docker compose up -d
```

---

## 🔒 Security Hardening & RBAC

1. **Role-Based Access Control (RBAC)**:
   - `SuperAdmin`: Full platform, multi-tenant org creation, security posture, and root controls.
   - `Admin`: Organization level administration, approval processing, and rule management.
   - `Operator`: Infrastructure actions, container control, and runbook execution.
   - `Viewer`: Read-only telemetry access.

2. **mTLS & Secrets Management**:
   - AES-256 GCM encryption for stored SSH keys, API keys, and cloud credentials.
   - Audit logging recorded for every administrative action.

---

## 📦 Backup & Disaster Recovery

### One-Command Backup
```bash
./scripts/backup.sh
```

### One-Command Recovery
```bash
./scripts/restore.sh ./backups/20260727_120000
```
