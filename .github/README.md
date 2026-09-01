# CI/CD Pipeline Documentation

This directory contains the complete CI/CD pipeline configuration for InfraPilot Enterprise using GitHub Actions.

## Architecture

```
Developer
    │
    ▼
Git Push
    │
    ▼
GitHub Repository
    │
    ▼
GitHub Actions
    │
 ┌───┴───────────────────────────────┐
 │                                   │
Code Quality                     Security Scan
 │                                   │
 ▼                                   ▼
Unit Tests                     Dependency Scan
 │                                   │
 └──────────────┬────────────────────┘
                ▼
Docker Image Build
                │
                ▼
Push to Docker Registry
                │
                ▼
Helm Upgrade
                │
                ▼
Kubernetes Cluster
                │
                ▼
InfraPilot Production
```

## Workflows

### 1. Backend CI (`backend.yml`)

**Trigger**: Push to `main`/`develop`, Pull requests to `main`

**Pipeline Stages**:
- ✅ Checkout code
- ✅ Set up Go 1.24
- ✅ Download dependencies
- ✅ Verify dependencies
- ✅ Run `go fmt`
- ✅ Run `go vet`
- ✅ Static analysis (golangci-lint)
- ✅ Unit tests with race detection
- ✅ Coverage check (minimum 80%)
- ✅ Build binary
- ✅ Upload coverage report

**Quality Gates**:
- All tests pass
- Code coverage ≥ 80%
- No lint errors
- No build errors

### 2. Frontend CI (`frontend.yml`)

**Trigger**: Push to `main`/`develop`, Pull requests to `main`

**Pipeline Stages**:
- ✅ Checkout code
- ✅ Set up Node.js 22
- ✅ Install dependencies (npm ci)
- ✅ Run ESLint
- ✅ Type checking (TypeScript)
- ✅ Run tests
- ✅ Production build
- ✅ Upload build artifacts

**Quality Gates**:
- All tests pass
- No lint errors
- Successful production build

### 3. Security Scan (`security.yml`)

**Trigger**: Push to `main`/`develop`, Pull requests to `main`, Weekly (Mondays)

**Security Scans**:
- ✅ Trivy vulnerability scanner (filesystem & Docker images)
- ✅ Gosec security scanner (Go code)
- ✅ npm audit (Node.js dependencies)
- ✅ govulncheck (Go vulnerabilities)
- ✅ Hardcoded secrets detection (git-secrets)
- ✅ Upload results to GitHub Security

**Fail Conditions**:
- Critical/High severity vulnerabilities
- Hardcoded secrets detected
- Failed security scans

### 4. Docker Build (`docker.yml`)

**Trigger**: Push to `main`/`develop`, Pull requests to `main`

**Pipeline Stages**:
- ✅ Build and push backend image
- ✅ Build and push frontend image
- ✅ Build and push agent image
- ✅ Generate SBOM (Software Bill of Materials)
- ✅ Upload SBOMs

**Image Tags**:
- `latest` - Latest build
- `<sha_short>` - Git commit hash
- `<version>` - Git tag (e.g., v1.0.0)

**Registry**: GitHub Container Registry (GHCR)
- URL: `ghcr.io/<org>/<repo>/<component>:<tag>`

### 5. Deploy (`deploy.yml`)

**Trigger**: Push to `main`/`develop`, Pull requests to `main`

**Deployment Environments**:
- **Production**: `main` branch → `infrapilot` namespace
- **Staging**: `develop` branch → `infrapilot-staging` namespace

**Pipeline Stages**:
- ✅ Set up Helm
- ✅ Configure kubectl (production/staging)
- ✅ Determine environment
- ✅ Run database migrations
- ✅ Deploy with Helm (upgrade --install)
- ✅ Verify rollout status (backend & frontend)
- ✅ Health checks (backend & frontend APIs)
- ✅ Display deployment info

**Rollback Strategy**:
- Automatic rollback on failure using `helm rollback`
- Triggered if any step fails after deployment

**Quality Gates**:
- Database migrations succeed
- Pods are running
- Health checks pass
- Rollout completes within timeout

### 6. Release (`release.yml`)

**Trigger**: Push of version tags (e.g., `v1.0.0`)

**Pipeline Stages**:
- ✅ Generate changelog from last tag
- ✅ Generate release notes with installation instructions
- ✅ Create GitHub Release
- ✅ Package Helm chart
- ✅ Upload release artifacts
- ✅ Deploy to production with version tag
- ✅ Verify production deployment
- ✅ Send Slack notification

**Release Artifacts**:
- GitHub Release with changelog
- Helm chart package (.tgz)
- Software Bill of Materials (SBOM)

## Branch Strategy

```
main (production)
├── develop (staging)
├── feature/*
├── hotfix/*
└── release/*
```

### Branch Protection Rules

1. **main branch**:
   - Require pull request before merging
   - Require passing CI checks
   - Require at least 1 approval
   - Require status checks: backend, frontend, security, docker, deploy

2. **develop branch**:
   - Require passing CI checks
   - Deploy to staging automatically

3. **feature branches**:
   - Run CI but do not deploy

4. **hotfix branches**:
   - Expedited review process
   - Fast-track to production

## Deployment Environments

| Environment | Branch | Namespace | Helm Values | Kubeconfig Secret |
|-------------|--------|-----------|-------------|-------------------|
| Production | main | infrapilot | values-prod.yaml | KUBE_CONFIG |
| Staging | develop | infrapilot-staging | values-stage.yaml | KUBE_CONFIG_STAGING |

## Secrets Management

All secrets are stored in GitHub Secrets (encrypted at rest).

### Required Secrets

See [SECRETS.md](./SECRETS.md) for complete list of required secrets.

**Quick Checklist**:
- [ ] GHCR_TOKEN
- [ ] KUBE_CONFIG (production)
- [ ] KUBE_CONFIG_STAGING
- [ ] JWT_SECRET
- [ ] POSTGRES_PASSWORD
- [ ] REDIS_PASSWORD
- [ ] QDRANT_API_KEY
- [ ] SLACK_WEBHOOK

## Code Quality Gates

All PRs must pass:
1. ✅ Backend tests with ≥80% coverage
2. ✅ Frontend build success
3. ✅ Security scans (no critical/high vulnerabilities)
4. ✅ Docker images build successfully
5. ✅ No lint errors

## Notifications

### Deployment Notifications

Sent to Slack on:
- Production deployment success/failure
- Staging deployment success/failure
- New release published

**Notification Content**:
- Commit ID
- Branch
- Author
- Deployment status
- Environment
- Changelog (for releases)

## Monitoring

Track these metrics to improve delivery reliability:
- Build duration
- Test duration
- Deployment duration
- Failure rate
- Success rate
- Rollback count

## Troubleshooting

### Pipeline Failures

1. **Backend tests failing**:
   ```bash
   go test ./... -v
   go tool cover -func=coverage.out
   ```

2. **Frontend build failing**:
   ```bash
   cd frontend && npm ci && npm run build
   ```

3. **Security scan failures**:
   - Check Trivy results in GitHub Security tab
   - Update dependencies with vulnerabilities
   - Run `npm audit fix` or `go get -u`

4. **Deployment failures**:
   - Check kubectl logs: `kubectl logs -n infrapilot deployment/infrapilot-backend`
   - Check Helm release: `helm list -n infrapilot`
   - Check pod status: `kubectl get pods -n infrapilot`

5. **Health check failures**:
   - Verify services are running
   - Check ingress configuration
   - Test API endpoints manually

### Common Issues

- **Secret not found**: Verify secret is created in correct repository
- **Kubernetes auth failed**: Regenerate kubeconfig with correct permissions
- **Docker push failed**: Check registry credentials and permissions
- **Helm upgrade failed**: Check values files and image tags

## Best Practices

1. **Never commit secrets** to the repository
2. **Always run CI locally** before pushing: `act -n .github/workflows/backend.yml`
3. **Use semantic versioning** for tags: `v1.2.3`
4. **Keep workflows DRY**: Use reusable workflows and composite actions
5. **Cache dependencies**: Reduce build times
6. **Set timeouts**: Prevent hanging jobs
7. **Review security alerts**: Address critical/high vulnerabilities promptly
8. **Test rollback strategy**: Ensure fast recovery from failed deployments

## Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Helm Documentation](https://helm.sh/docs/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Docker Documentation](https://docs.docker.com/)
- [Trivy Documentation](https://aquasecurity.github.io/trivy/)

## Support

For issues with the CI/CD pipeline:
1. Check workflow logs in GitHub Actions tab
2. Review this documentation
3. Check [SECRETS.md](./SECRETS.md) for secret configuration
4. Open an issue in the repository

---

**Last Updated**: 2025-07-21
**Pipeline Version**: 1.0.0
**Maintained by**: InfraPilot DevOps Team