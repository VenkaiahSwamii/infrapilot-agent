# InfraPilot Enterprise Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of InfraPilot Enterprise seriously. If you believe you have found a security vulnerability, please report it to us as described below.

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report them via email to **venkaiswami@pm.me** (or the project maintainer's security contact).

You should receive a response within 48 hours. If for some reason you do not, please follow up via email to ensure we received your original message.

Please include the following information:
- Type of issue (e.g., buffer overflow, SQL injection, cross-site scripting, etc.)
- Full paths of source file(s) related to the manifestation of the issue
- The location of the affected source code (tag/branch/commit or direct URL)
- Any special configuration required to reproduce the issue
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the issue, including how an attacker might exploit it

## Security Features

### Authentication & Authorization
- **JWT-based authentication** with access/refresh token rotation
- **Password hashing** using bcrypt with configurable cost factor
- **API key authentication** for agent-to-server communication
- **Role-Based Access Control (RBAC)** with four tiers:
  - `SuperAdmin`: Full system access
  - `Admin`: Administrative operations
  - `Operator`: Operational commands
  - `Viewer`: Read-only access

### Network Security
- **Rate limiting** on API endpoints (configurable requests/minute/IP)
- **CORS middleware** with configurable allowed origins
- **HTTPS enforcement** via Nginx reverse proxy
- **WebSocket connections** authenticated via JWT tokens
- **Internal port binding** (PostgreSQL, Redis bound to 127.0.0.1)

### Data Security
- **SQL injection prevention** via GORM parameterized queries
- **Sensitive data redaction** in structured logs
- **Environment-based configuration** (no hardcoded secrets)
- **Encryption key rotation** support for agent API keys
- **Database connection pooling** with configurable limits

### Infrastructure Security
- **Docker containers** run as non-root users where possible
- **Health checks** on all services for availability monitoring
- **Graceful shutdown** for clean resource cleanup
- **Audit logging** for all sensitive operations
- **Prometheus metrics** restricted to internal network

## Security Configuration Checklist

### Production Deployment
- [ ] Change all default passwords and secrets
- [ ] Enable HTTPS with valid SSL certificates
- [ ] Configure CORS to allow only trusted origins
- [ ] Set rate limits appropriate for your use case
- [ ] Enable structured JSON logging for log aggregation
- [ ] Restrict database and Redis ports to localhost
- [ ] Configure firewall rules for your cloud provider
- [ ] Enable regular database backups
- [ ] Set up log shipping to a centralized system (e.g., Loki, ELK)
- [ ] Configure alerting on security-relevant audit events

### Environment Variables
```bash
# Required changes from defaults
JWT_SECRET=<generate-a-strong-random-secret>
DB_PASSWORD=<generate-a-strong-random-password>
REDIS_PASSWORD=<generate-a-strong-random-password>
ENROLLMENT_TOKEN=<generate-a-unique-enrollment-key>
API_ENCRYPTION_KEY=<generate-a-random-32-byte-hex-key>

# Security hardening
CORS_ORIGINS=https://your-domain.com
SERVER_MODE=release
LOG_LEVEL=info
LOG_JSON=true
RATE_LIMIT_PER_MINUTE=60
```

## Dependency Security

We use automated dependency scanning via:
- **GitHub Dependabot** for automated dependency updates
- **Trivy** vulnerability scanner in CI pipeline
- **Go vulnerability checker** (`govulncheck`) for Go dependencies
- **npm audit** for JavaScript dependencies

## Incident Response

In case of a security incident:
1. **Isolate** affected systems immediately
2. **Rotate** all credentials and API keys
3. **Analyze** logs to determine scope and impact
4. **Patch** vulnerabilities and deploy fixes
5. **Communicate** with affected users/stakeholders
6. **Document** lessons learned and update procedures

## Responsible Disclosure

We kindly ask that:
- You give us reasonable time to fix the issue before public disclosure
- You make a good faith effort to avoid privacy violations and data destruction
- You do not access or modify user data without explicit permission

We strive to:
- Respond promptly to all valid reports
- Provide clear and timely updates on the resolution process
- Give credit to researchers who report valid vulnerabilities (with their permission)

## License

This security policy is provided under the same license as InfraPilot Enterprise (MIT).