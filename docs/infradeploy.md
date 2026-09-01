# InfraDeploy - Enterprise CI/CD Platform

## Overview

InfraDeploy is the **Enterprise CI/CD Platform** within the InfraPilot ecosystem, providing continuous integration and deployment capabilities with deep integration into your existing infrastructure.

## Features

### Git Repository Management
- **Multi-Provider Support**: GitHub, GitLab, Bitbucket, Gitea, Azure DevOps
- **Secure Authentication**: OAuth, SSH, Token, Basic auth
- **Webhook Integration**: Real-time event triggers
- **Branch Management**: Configurable default branches
- **Repository Sync**: Automatic synchronization

### CI/CD Pipelines
- **Flexible Triggers**: Webhook, Schedule, Manual, API
- **Environment Support**: Development, Staging, Production, Canary
- **Pipeline Configuration**: JSON/YAML-based pipeline configs
- **Scheduled Builds**: Cron-based automated builds
- **Pipeline Enable/Disable**: Quick toggle for pipelines

### Build Management
- **Build Numbering**: Automatic incrementing build numbers
- **Status Tracking**: Pending, Running, Success, Failed, Cancelled
- **Build Logs**: Complete build logs stored and accessible
- **Artifact Storage**: JSON-based artifact URLs
- **Duration Tracking**: Precise build timing
- **Commit Tracking**: Full commit hash, message, author, branch

### Deployment Strategies
- **Rolling Updates**: Zero-d downtime rollouts
- **Blue-Green Deployments**: Instant failover capability
- **Canary Deployments**: Gradual traffic shifting
- **Recreate**: Clean slate deployments
- **Auto Rollback**: Automatic rollback on failure
- **Health Checks**: Configurable health endpoints

### Deployment Features
- **Multi-Environment**: Dev, Staging, Production, Canary
- **Kubernetes Integration**: Native K8s deployments
- **Helm Support**: Helm chart deployment
- **ArgoCD Integration**: GitOps workflows
- **Namespace Management**: K8s namespace support
- **Cluster Support**: Multi-cluster deployments
- **Replica Management**: Scale up/down

### Webhook Integration
- **Git Provider Webhooks**: GitHub, GitLab, Bitbucket, Gitea
- **Event Types**: Push, Pull Request, Tag, Release, Ping
- **Signature Verification**: Webhook signature validation
- **Payload Storage**: JSON payload archival
- **Processing Status**: Track webhook processing

### Deployment History
- **Complete Audit Trail**: All deployment actions logged
- **Rollback History**: Track all rollback operations
- **Change Tracking**: Image and replica changes
- **User Attribution**: Who triggered each action
- **Timeline View**: Chronological deployment history

### Statistics & Monitoring
- **Pipeline Count**: Total pipelines per organization
- **Build Metrics**: Success/failure counts
- **Deployment Count**: Total deployments
- **Real-time Updates**: Live statistics

## API Endpoints

### Git Repositories
```
POST   /api/v1/deploy/repositories          - Create repository
GET    /api/v1/deploy/repositories          - List repositories
GET    /api/v1/deploy/repositories/:id      - Get repository
PATCH  /api/v1/deploy/repositories/:id      - Update repository
DELETE /api/v1/deploy/repositories/:id      - Delete repository
```

### Pipelines
```
POST   /api/v1/deploy/pipelines             - Create pipeline
GET    /api/v1/deploy/pipelines             - List pipelines
GET    /api/v1/deploy/pipelines/:id         - Get pipeline
PATCH  /api/v1/deploy/pipelines/:id         - Update pipeline
DELETE /api/v1/deploy/pipelines/:id         - Delete pipeline
```

### Builds
```
POST   /api/v1/deploy/builds                - Create build
GET    /api/v1/deploy/builds                - List builds
GET    /api/v1/deploy/builds/:id            - Get build
PATCH  /api/v1/deploy/builds/:id/status     - Update build status
```

### Deployments
```
POST   /api/v1/deploy/deployments           - Create deployment
GET    /api/v1/deploy/deployments           - List deployments
GET    /api/v1/deploy/deployments/:id       - Get deployment
POST   /api/v1/deploy/deployments/:id/rollback - Rollback deployment
```

### Webhooks
```
POST   /api/v1/deploy/webhooks/:provider     - Handle webhook
```

### Statistics
```
GET    /api/v1/deploy/stats                 - Get deployment statistics
```

## Database Schema

### Tables
- `git_repositories` - Git repository configurations
- `pipelines` - CI/CD pipeline definitions
- `builds` - Build execution records
- `deployments` - Deployment records
- `deployment_history` - Deployment action history
- `webhook_events` - Incoming webhook events
- `deployment_environments` - Environment configurations

## Frontend Features

### Overview Dashboard
- **Statistics Cards**: Key metrics at a glance
- **Recent Builds**: Latest build activity
- **Recent Deployments**: Latest deployment activity
- **Tab Navigation**: Easy switching between views

### Git Repositories View
- **Card-based Layout**: Visual repository cards
- **Provider Badges**: Visual provider identification
- **Quick Actions**: Configure and delete actions
- **Repository Info**: URL, branch, auth type display

### Pipelines View
- **Table View**: Comprehensive pipeline listing
- **Environment Badges**: Visual environment indicators
- **Status Indicators**: Enabled/disabled status
- **Repository Linking**: Direct repository references

### Builds View
- **Build History**: Complete build history
- **Status Badges**: Visual status indicators
- **Commit Details**: Hash, message, author, branch
- **Duration Tracking**: Build execution time

### Deployments View
- **Deployment Table**: All deployments listed
- **Strategy Display**: Deployment strategy shown
- **Replica Count**: Number of replicas
- **Rollback Action**: Quick rollback button

## Integration with InfraPilot

InfraDeploy is part of the InfraPilot Enterprise ecosystem:

- **Unified Authentication**: JWT-based auth shared across products
- **Event Bus Integration**: Real-time event publishing
- **WebSocket Support**: Live updates via WebSocket
- **Common Database**: Shared PostgreSQL database
- **Single Sign-On**: Integrated with InfraPortal

## Configuration

### Environment Variables
```env
DATABASE_URL=postgresql://user:pass@localhost:5432/infrapilot
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-secret-key
SERVER_PORT=8080
```

### Pipeline Configuration Example
```json
{
  "name": "Backend API Pipeline",
  "trigger_type": "webhook",
  "environment": "production",
  "config": {
    "steps": [
      {
        "name": "Build",
        "image": "golang:1.21",
        "commands": ["go build", "go test"]
      },
      {
        "name": "Docker Build",
        "image": "docker:latest",
        "commands": ["docker build -t app:${BUILD_NUMBER} ."]
      }
    ]
  }
}
```

### Deployment Configuration Example
```json
{
  "name": "Production API",
  "environment": "production",
  "namespace": "production",
  "cluster": "us-west-2",
  "strategy": "rolling",
  "image": "registry.example.com/api:v1.2.3",
  "replicas": 3,
  "health_check": "/healthz",
  "auto_rollback": true
}
```

## Tech Stack

### Backend
- **Go 1.21**: Core backend language
- **Gin**: Web framework
- **GORM**: Database ORM
- **PostgreSQL**: Primary database
- **Redis**: Caching layer
- **WebSocket**: Real-time communication

### Frontend
- **React 18**: UI framework
- **JavaScript**: Primary language
- **Axios**: HTTP client
- **CSS3**: Styling
- **React Router**: Navigation

## Event Types

InfraDeploy publishes the following events:

- `deploy.git_repository.created` - New git repository added
- `deploy.pipeline.created` - New pipeline created
- `deploy.build.created` - New build initiated
- `deploy.build.status_changed` - Build status updated
- `deploy.deployment.created` - New deployment created
- `deploy.deployment.rolled_back` - Deployment rolled back
- `deploy.webhook.received` - Webhook received

## Security Features

- **JWT Authentication**: All endpoints secured
- **Role-Based Access Control**: Admin, Operator, Viewer roles
- **Webhook Signatures**: GitHub/GitLab signature verification
- **Credential Encryption**: Tokens and SSH keys encrypted
- **Audit Logging**: All actions logged

## Deployment

### Docker Compose
```yaml
services:
  backend:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgresql://user:pass@postgres:5432/infrapilot
      - REDIS_URL=redis://redis:6379
```

### Kubernetes
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: infradeploy-backend
spec:
  replicas: 3
  template:
    spec:
      containers:
        - name: backend
          image: infrapilot/infradeploy:latest
          ports:
            - containerPort: 8080
```

## Roadmap

### Phase 1: Core CI/CD (Complete)
- ✅ Git Repository Management
- ✅ Pipeline Configuration
- ✅ Build Execution
- ✅ Deployment Management
- ✅ Webhook Integration

### Phase 2: Advanced Features (Next)
- Multi-cluster deployments
- Deployment approvals workflow
- Advanced rollback strategies
- Integration with Jenkins, GitHub Actions, GitLab CI
- Terraform integration

### Phase 3: Enterprise Features
- Advanced RBAC for deployments
- Deployment templates
- Canary analysis
- A/B testing support
- Advanced monitoring integrations

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md) for contribution guidelines.

## License

Part of InfraPilot Enterprise - See [LICENSE](../LICENSE) for details.