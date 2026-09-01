# 🛡️ InfraPilot Enterprise Security Policy

## Security Model Overview

InfraPilot Enterprise is designed following zero-trust principles and Defense in Depth architecture.

---

## 1. Authentication & Session Security
- **JWT Expiration & Refresh Rotation**: Access tokens expire after short durations; refresh tokens are rotated and stored securely.
- **Password Protection**: Passwords are hashed using bcrypt with custom cost factors.
- **MFA Support**: Multi-Factor Authentication via TOTP.

---

## 2. API & Network Security
- **Rate Limiting**: Integrated token bucket rate limiter to prevent denial of service (DoS).
- **Helmet Security Headers**: Strictly configured Content-Security-Policy (CSP), HSTS, X-Frame-Options, X-Content-Type-Options, and Referrer-Policy.
- **Agent TLS 1.3 Encryption**: All agent telemetry is encrypted in transit over TLS 1.3.

---

## 3. Role-Based Access Control (RBAC)

| Role | Access Scope |
|---|---|
| **SuperAdmin** | Full system control, billing, tenant administration, compliance settings |
| **Admin** | Organization user management, alert configuration, automation runbooks |
| **DevOps Engineer**| Read/Write infrastructure, incident resolution, remote actions |
| **Auditor** | Read-only access to audit logs, compliance compliance dashboards, reports |
| **Viewer** | Read-only metric and dashboard access |

---

## Reporting Vulnerabilities

To report a vulnerability, please email `security@infrapilot.io`. Reports will be acknowledged within 24 hours.
