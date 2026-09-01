# GitHub Secrets Configuration

This document describes all the secrets required for the CI/CD pipeline to function correctly.

## Secrets Setup

All secrets must be configured in your GitHub repository settings:
**Settings → Secrets and variables → Actions → New repository secret**

## Required Secrets

### Container Registry

#### `GHCR_TOKEN`
- **Description**: GitHub Container Registry authentication token
- **How to create**: 
  1. Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
  2. Generate new token with `write:packages`, `read:packages`, and `repo` scopes
  3. Copy the token and store it as `GHCR_TOKEN`

#### `DOCKER_USERNAME` (Alternative)
- **Description**: Docker Hub username (if using Docker Hub instead of GHCR)
- **Used in**: docker.yml

#### `DOCKER_PASSWORD` (Alternative)
- **Description**: Docker Hub password or access token
- **Used in**: docker.yml

### Kubernetes Authentication

#### `KUBE_CONFIG` (Production)
- **Description**: Base64-encoded kubeconfig file for production cluster
- **How to create**:
  ```bash
  kubectl config view --raw --flatten | base64 -w 0
  ```
- **Used in**: deploy.yml (production deployments)

#### `KUBE_CONFIG_STAGING`
- **Description**: Base64-encoded kubeconfig file for staging cluster
- **How to create**: Same as above but for staging cluster
- **Used in**: deploy.yml (staging deployments)

### Application Secrets

These secrets are used in Kubernetes deployments via Helm values and ConfigMaps.

#### `JWT_SECRET`
- **Description**: Secret key for JWT token signing
- **Format**: 64+ character random string
- **Example generation**:
  ```bash
  openssl rand -hex 32
  ```
- **Used in**: Backend authentication middleware

#### `POSTGRES_PASSWORD`
- **Description**: PostgreSQL database password
- **Format**: Strong password (minimum 16 characters)
- **Used in**: Database initialization and backend configuration

#### `POSTGRES_USER` (Optional)
- **Description**: PostgreSQL username
- **Default**: `postgres`

#### `POSTGRES_DB` (Optional)
- **Description**: PostgreSQL database name
- **Default**: `infrapilot`

#### `REDIS_PASSWORD`
- **Description**: Redis authentication password
- **Format**: Strong password (minimum 16 characters)
- **Used in**: Redis configuration and backend cache

#### `QDRANT_API_KEY`
- **Description**: API key for Qdrant vector database
- **Format**: 64+ character random string
- **Used in**: AI/ML features and vector search

### Notifications

#### `SLACK_WEBHOOK`
- **Description**: Slack incoming webhook URL for deployment notifications
- **How to create**:
  1. Go to Slack → Settings → Apps → Incoming Webhooks
  2. Activate incoming webhooks
  3. Add webhook to your workspace and channel
  4. Copy the webhook URL
- **Used in**: release.yml, deploy.yml (optional)

#### `TEAMS_WEBHOOK` (Optional)
- **Description**: Microsoft Teams incoming webhook URL
- **Used in**: Alternative notification channel

#### `EMAIL_USERNAME` (Optional)
- **Description**: SMTP email username
- **Used in**: Email notifications

#### `EMAIL_PASSWORD` (Optional)
- **Description**: SMTP email password or app-specific password
- **Used in**: Email notifications

### Helm (Optional)

#### `HELM_TOKEN`
- **Description**: GitHub token for pushing Helm charts to GitHub Pages
- **Scope**: `repo` permissions required
- **Used in**: release.yml

## Environment-Specific Secrets

### Production Environment
- `KUBE_CONFIG` - Production Kubernetes cluster
- `SLACK_WEBHOOK_PROD` - Production notification channel
- All application secrets with production values

### Staging Environment
- `KUBE_CONFIG_STAGING` - Staging Kubernetes cluster
- `SLACK_WEBHOOK_STAGING` - Staging notification channel
- Application secrets can be shared or have staging-specific values

## Secret Rotation Policy

1. **Critical secrets** (JWT, database passwords, API keys): Rotate every 90 days
2. **Registry tokens**: Rotate every 180 days
3. **Notification webhooks**: Review annually or when team access changes

## Security Best Practices

1. **Never commit secrets** to the repository
2. **Use environment-specific secrets** for different deployment stages
3. **Enable secret scanning** in GitHub repository settings
4. **Limit secret access** to required workflows only
5. **Audit secret usage** regularly via GitHub Actions logs
6. **Use GitHub Environments** for additional secret scoping (production vs staging)

## Verification

To verify secrets are properly configured:

1. Go to **Settings → Secrets and variables → Actions**
2. Ensure all required secrets are listed
3. Check that secrets are encrypted (values are hidden)
4. Verify secret names match workflow references

## Troubleshooting

### Secret not found error
- Check for typos in secret names (case-sensitive)
- Ensure secret is in the correct repository
- Verify the workflow has access to the secret

### Authentication failures
- Verify token hasn't expired
- Check token scopes/permissions
- Regenerate token if necessary

### Deployment failures
- Validate kubeconfig is correct for the target environment
- Ensure Kubernetes cluster is accessible
- Check that Helm values reference correct secrets

## Additional Resources

- [GitHub Actions Secrets Documentation](https://docs.github.com/en/actions/security-guides/encrypted-secrets)
- [GitHub Container Registry Documentation](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- [Kubernetes Authentication Documentation](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/)