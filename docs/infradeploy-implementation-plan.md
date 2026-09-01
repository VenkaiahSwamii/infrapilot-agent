# InfraDeploy Implementation Plan

Complete implementation guide for building InfraDeploy - Enterprise CI/CD Platform.

---

## Executive Summary

**Project:** InfraDeploy  
**Timeline:** 8-10 weeks  
**Team Size:** 1-2 developers  
**Goal:** Build production-ready CI/CD platform with Git integration, build pipelines, and deployment strategies

### Success Criteria
- [ ] Support 5+ Git providers (GitHub, GitLab, Bitbucket, Gitea, Azure DevOps)
- [ ] Handle 1000+ builds per day
- [ ] Deployment success rate > 99%
- [ ] Rollback time < 30 seconds
- [ ] API latency < 100ms (p95)
- [ ] 80%+ test coverage

---

## Phase 1: Foundation & Git Integration (Weeks 1-2)

### Week 1: Project Setup & Database

#### Tasks

1. **Project Structure**
   ```
   infradeploy/
   ├── backend/
   │   ├── cmd/api/
   │   ├── internal/
   │   │   ├── handlers/
   │   │   ├── models/
   │   │   ├── repository/
   │   │   ├── services/
   │   │   ├── events/
   │   │   └── subscribers/
   │   ├── migrations/
   │   └── tests/
   ├── frontend/
   │   └── src/
   │       ├── features/
   │       │   ├── overview/
   │       │   ├── repositories/
   │       │   ├── pipelines/
   │       │   ├── builds/
   │       │   └── deployments/
   │       ├── api/
   │       └── components/
   ├── docs/
   ├── k8s/
   ├── docker-compose.yml
   └── README.md
   ```

2. **Database Schema**
   ```sql
   -- Git Repositories
   CREATE TABLE git_repositories (
       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       organization_id UUID NOT NULL,
       name VARCHAR(255) NOT NULL,
       provider VARCHAR(50) NOT NULL, -- github, gitlab, bitbucket, gitea
       url VARCHAR(500) NOT NULL,
       webhook_url VARCHAR(500),
       auth_type VARCHAR(50) NOT NULL, -- oauth, ssh, token, basic
       encrypted_credentials TEXT,
       default_branch VARCHAR(100) DEFAULT 'main',
       is_active BOOLEAN DEFAULT true,
       created_at TIMESTAMP DEFAULT NOW(),
       updated_at TIMESTAMP DEFAULT NOW()
   );

   -- Pipelines
   CREATE TABLE pipelines (
       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       repository_id UUID REFERENCES git_repositories(id),
       name VARCHAR(255) NOT NULL,
       description TEXT,
       trigger_type VARCHAR(50), -- webhook, schedule, manual, api
       cron_schedule VARCHAR(100),
       config JSONB NOT NULL,
       environment VARCHAR(50), -- dev, staging, prod, canary
       is_enabled BOOLEAN DEFAULT true,
       last_build_id UUID,
       created_at TIMESTAMP DEFAULT NOW(),
       updated_at TIMESTAMP DEFAULT NOW()
   );

   -- Builds
   CREATE TABLE builds (
       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       pipeline_id UUID REFERENCES pipelines(id),
       build_number INTEGER NOT NULL,
       status VARCHAR(50), -- pending, running, success, failed, cancelled
       trigger VARCHAR(50), -- webhook, schedule, manual, api
       commit_hash VARCHAR(40),
       commit_message TEXT,
       branch VARCHAR(100),
       author VARCHAR(255),
       started_at TIMESTAMP,
       finished_at TIMESTAMP,
       duration_ms INTEGER,
       logs TEXT,
       artifacts JSONB,
       created_at TIMESTAMP DEFAULT NOW()
   );

   -- Deployments
   CREATE TABLE deployments (
       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       build_id UUID REFERENCES builds(id),
       pipeline_id UUID REFERENCES pipelines(id),
       environment VARCHAR(50) NOT NULL,
       strategy VARCHAR(50), -- rolling, blue-green, canary, recreate
       namespace VARCHAR(100),
       cluster VARCHAR(100),
       image VARCHAR(500),
       replicas INTEGER,
       health_check VARCHAR(255),
       auto_rollback BOOLEAN DEFAULT true,
       status VARCHAR(50), -- pending, deploying, success, failed, rolled_back
       created_by UUID,
       created_at TIMESTAMP DEFAULT NOW(),
       deployed_at TIMESTAMP
   );

   -- Deployment History
   CREATE TABLE deployment_history (
       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       deployment_id UUID REFERENCES deployments(id),
       action VARCHAR(50), -- deploy, rollback, scale
       status VARCHAR(50),
       message TEXT,
       replicas_before INTEGER,
       replicas_after INTEGER,
       performed_by UUID,
       performed_at TIMESTAMP DEFAULT NOW()
   );

   -- Webhook Events
   CREATE TABLE webhook_events (
       id UUID PRIMARY KEY DEFAULT gen_random_uuid),
       repository_id UUID REFERENCES git_repositories(id),
       provider VARCHAR(50) NOT NULL,
       event_type VARCHAR(50), -- push, pull_request, tag, release
       payload JSONB NOT NULL,
       signature VARCHAR(255),
       processed BOOLEAN DEFAULT false,
       error_message TEXT,
       received_at TIMESTAMP DEFAULT NOW(),
       processed_at TIMESTAMP
   );

   -- Deployment Environments
   CREATE TABLE deployment_environments (
       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       organization_id UUID NOT NULL,
       name VARCHAR(100) NOT NULL,
       cluster VARCHAR(100),
       namespace VARCHAR(100),
       config JSONB,
       is_default BOOLEAN DEFAULT false,
       created_at TIMESTAMP DEFAULT NOW()
   );

   -- Indexes
   CREATE INDEX idx_repositories_org ON git_repositories(organization_id);
   CREATE INDEX idx_pipelines_repo ON pipelines(repository_id);
   CREATE INDEX idx_builds_pipeline ON builds(pipeline_id);
   CREATE INDEX idx_builds_status ON builds(status);
   CREATE INDEX idx_deployments_build ON deployments(build_id);
   CREATE INDEX idx_webhooks_repo ON webhook_events(repository_id);
   ```

3. **Go Module Setup**
   ```bash
   go mod init github.com/venkaiswami/infradeploy
   go get github.com/gin-gonic/gin
   go get github.com/lib/pq
   go get github.com/go-redis/redis/v8
   go get github.com/gorilla/websocket
   go get github.com/xanzy/go-gitlab
   go get github.com/google/go-github/v50
   go get gopkg.in/yaml.v3
   ```

4. **Shared Infrastructure**
   - Connection pooling for PostgreSQL
   - Redis for caching and job queues
   - Structured logging (same as InfraPilot)
   - JWT authentication middleware
   - Rate limiting

### Week 2: Git Provider Integration

#### OAuth Flow

```go
// models/git_repository.go
type GitProvider string

const (
    GitHub GitProvider = "github"
    GitLab GitProvider = "gitlab"
    Bitbucket GitProvider = "bitbucket"
    Gitea GitProvider = "gitea"
)

type GitRepository struct {
    ID                  uuid.UUID
    OrganizationID      uuid.UUID
    Name                string
    Provider            GitProvider
    URL                 string
    WebhookURL          string
    AuthType            string
    EncryptedCredentials []byte
    DefaultBranch        string
    IsActive            bool
    CreatedAt           time.Time
    UpdatedAt           time.Time
}

type GitClient interface {
    ListRepositories(ctx context.Context) ([]*Repository, error)
    GetWebhookSecret(ctx context.Context) (string, error)
    VerifyWebhookSignature(payload []byte, signature string) bool
    CreateWebhook(ctx context.Context, repoID string, webhook WebhookConfig) error
}
```

#### GitHub Integration

```go
// clients/github.go
type GitHubClient struct {
    client *github.Client
    token  string
}

func (g *GitHubClient) ListRepositories(ctx context.Context) ([]*Repository, error) {
    repos, _, err := g.client.Repositories.List(ctx, "", &github.RepositoryListOptions{
        ListOptions: github.ListOptions{PerPage: 100},
    })
    
    if err != nil {
        return nil, err
    }
    
    var result []*Repository
    for _, repo := range repos {
        result = append(result, &Repository{
            ID:       repo.ID,
            Name:     repo.Name,
            FullName: repo.FullName,
            URL:      repo.CloneURL,
            Private:  repo.Private,
        })
    }
    
    return result, nil
}

func (g *GitHubClient) VerifyWebhookSignature(payload []byte, signature string) bool {
    h := hmac.New(sha256.New, []byte(g.webhookSecret))
    h.Write(payload)
    expectedSignature := "sha256=" + hex.EncodeToString(h.Sum(nil))
    
    return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
```

#### Webhook Handler

```go
// handlers/webhook.go
func (h *Handlers) HandleGitHubWebhook(c *gin.Context) {
    repositoryID := c.Param("repositoryId")
    
    payload, err := c.GetRawData()
    if err != nil {
        c.JSON(400, gin.H{"error": "Invalid payload"})
        return
    }
    
    signature := c.GetHeader("X-Hub-Signature-256")
    
    // Verify signature
    repo, err := h.repositoryRepo.GetByID(c, repositoryID)
    if err != nil {
        c.JSON(404, gin.H{"error": "Repository not found"})
        return
    }
    
    client := h.gitClients[repo.Provider]
    if !client.VerifyWebhookSignature(payload, signature) {
        c.JSON(401, gin.H{"error": "Invalid signature"})
        return
    }
    
    // Parse webhook event
    event, err := h.parseGitHubEvent(payload)
    if err != nil {
        c.JSON(400, gin.H{"error": "Invalid event"})
        return
    }
    
    // Store webhook event
    webhookEvent := &WebhookEvent{
        RepositoryID:  repositoryID,
        Provider:      string(repo.Provider),
        EventType:     event.Type,
        Payload:       payload,
        Signature:     signature,
        Processed:     false,
        ReceivedAt:    time.Now(),
    }
    
    h.webhookRepo.Create(c, webhookEvent)
    
    // Publish to event bus
    h.eventBus.Publish(Event{
        Type: "webhook.received",
        Data: webhookEvent,
    })
    
    c.JSON(200, gin.H{"status": "accepted"})
}
```

---

## Phase 2: Pipeline & Build System (Weeks 3-4)

### Week 3: Pipeline Engine

#### Pipeline Configuration

```go
// models/pipeline.go
type Pipeline struct {
    ID          uuid.UUID
    RepositoryID uuid.UUID
    Name        string
    Description string
    TriggerType string // webhook, schedule, manual, api
    CronSchedule string
    Config      PipelineConfig
    Environment string
    IsEnabled   bool
    LastBuildID *uuid.UUID
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type PipelineConfig struct {
    Steps []PipelineStep `json:"steps"`
    Environment map[string]string `json:"environment"`
    Cache []CacheConfig `json:"cache"`
    Timeout int `json:"timeout"` // minutes
}

type PipelineStep struct {
    Name        string            `json:"name"`
    Image       string            `json:"image"`
    Commands    []string          `json:"commands"`
    Environment map[string]string `json:"environment"`
    DependsOn   []string          `json:"depends_on"`
    Timeout     int               `json:"timeout"`
}
```

#### Pipeline Executor

```go
// services/pipeline_executor.go
type PipelineExecutor struct {
    dockerClient *docker.Client
    eventBus     EventBus
    logger       *slog.Logger
}

func (pe *PipelineExecutor) Execute(ctx context.Context, pipeline *Pipeline, build *Build) error {
    // Update build status
    build.Status = "running"
    build.StartedAt = time.Now()
    pe.buildRepo.Update(ctx, build)
    
    // Publish event
    pe.eventBus.Publish(Event{
        Type: "build.started",
        Data: build,
    })
    
    // Execute steps
    for _, step := range pipeline.Config.Steps {
        if err := pe.executeStep(ctx, step, build); err != nil {
            build.Status = "failed"
            build.FinishedAt = time.Now()
            build.DurationMs = int(time.Since(build.StartedAt).Milliseconds())
            pe.buildRepo.Update(ctx, build)
            
            pe.eventBus.Publish(Event{
                Type: "build.failed",
                Data: build,
            })
            
            return err
        }
    }
    
    // Build success
    build.Status = "success"
    build.FinishedAt = time.Now()
    build.DurationMs = int(time.Since(build.StartedAt).Milliseconds())
    pe.buildRepo.Update(ctx, build)
    
    pe.eventBus.Publish(Event{
        Type: "build.completed",
        Data: build,
    })
    
    return nil
}

func (pe *PipelineExecutor) executeStep(ctx context.Context, step PipelineStep, build *Build) error {
    pe.logger.Info("Executing step", "step", step.Name, "build", build.ID)
    
    // Create Docker container for step
    container, err := pe.dockerClient.ContainerCreate(ctx, &container.Config{
        Image: step.Image,
        Cmd:   step.Commands,
        Env:   pe.buildEnvVars(step),
    }, nil, nil, nil, "")
    
    if err != nil {
        return err
    }
    
    // Start container
    if err := pe.dockerClient.ContainerStart(ctx, container.ID, container.StartOptions); err != nil {
        return err
    }
    
    // Wait for completion
    statusCh, errCh := pe.dockerClient.ContainerWait(ctx, container.ID, container.WaitCondition)
    select {
    case err := <-errCh:
        return err
    case status := <-statusCh:
        if status.StatusCode != 0 {
            return fmt.Errorf("step failed with status: %d", status.StatusCode)
        }
    }
    
    // Collect logs
    logs, _ := pe.dockerClient.ContainerLogs(ctx, container.ID, container.LogsOptions)
    build.Logs += string(logs)
    
    // Cleanup
    pe.dockerClient.ContainerRemove(ctx, container.ID, container.RemoveOptions)
    
    return nil
}
```

### Week 4: Docker Image Building

#### Docker Builder

```go
// services/docker_builder.go
type DockerBuilder struct {
    dockerClient *docker.Client
    registry     *RegistryClient
}

func (db *DockerBuilder) Build(ctx context.Context, build *Build, step PipelineStep) (string, error) {
    dockerfile := fmt.Sprintf("./builds/%s/Dockerfile", build.ID)
    
    // Write Dockerfile
    if err := os.WriteFile(dockerfile, []byte(step.Config["dockerfile"]), 0644); err != nil {
        return "", err
    }
    
    // Build image
    imageName := fmt.Sprintf("registry.example.com/%s:%d", 
        strings.ToLower(build.Repository), build.BuildNumber)
    
    resp, err := db.dockerClient.ImageBuild(ctx, nil, types.ImageBuildOptions{
        Dockerfile: dockerfile,
        Context:    open(dockerfile),
        Tags:       []string{imageName},
    })
    
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    // Push to registry
    if err := db.registry.Push(ctx, imageName); err != nil {
        return "", err
    }
    
    // Store artifact
    build.Artifacts = []Artifact{
        {
            Name: imageName,
            Type: "docker_image",
            URL:  imageName,
        },
    }
    
    return imageName, nil
}
```

---

## Phase 3: Deployment Engine (Weeks 5-6)

### Week 5: Kubernetes Deployment

#### Deployment Service

```go
// services/deployment_service.go
type DeploymentService struct {
    clientset     *kubernetes.Clientset
    helmClient    *helmclient.Client
    eventBus      EventBus
}

func (ds *DeploymentService) Deploy(ctx context.Context, deployment *Deployment) error {
    // Update status
    deployment.Status = "deploying"
    ds.deploymentRepo.Update(ctx, deployment)
    
    ds.eventBus.Publish(Event{
        Type: "deployment.started",
        Data: deployment,
    })
    
    switch deployment.Strategy {
    case "rolling":
        return ds.rollingDeploy(ctx, deployment)
    case "blue-green":
        return ds.blueGreenDeploy(ctx, deployment)
    case "canary":
        return ds.canaryDeploy(ctx, deployment)
    case "recreate":
        return ds.recreateDeploy(ctx, deployment)
    default:
        return fmt.Errorf("unknown strategy: %s", deployment.Strategy)
    }
}

func (ds *DeploymentService) rollingDeploy(ctx context.Context, deployment *Deployment) error {
    // Get current deployment
    deploy, err := ds.clientset.AppsV1().Deployments(deployment.Namespace).Get(
        ctx, deployment.Name, metav1.GetOptions{})
    
    if err != nil {
        // Create new deployment
        deploy = &appsv1.Deployment{
            ObjectMeta: metav1.ObjectMeta{
                Name:      deployment.Name,
                Namespace: deployment.Namespace,
            },
            Spec: appsv1.DeploymentSpec{
                Replicas: &deployment.Replicas,
                Selector: &metav1.LabelSelector{
                    MatchLabels: map[string]string{"app": deployment.Name},
                },
                Template: corev1.PodTemplateSpec{
                    ObjectMeta: metav1.ObjectMeta{
                        Labels: map[string]string{"app": deployment.Name},
                    },
                    Spec: corev1.PodSpec{
                        Containers: []corev1.Container{
                            {
                                Name:  deployment.Name,
                                Image: deployment.Image,
                                Ports: []corev1.ContainerPort{
                                    {ContainerPort: 8080},
                                },
                                LivenessProbe: &corev1.Probe{
                                    ProbeHandler: corev1.Handler{
                                        HTTPGet: &corev1.HTTPGetAction{
                                            Path: deployment.HealthCheck,
                                            Port: intstr.FromInt(8080),
                                        },
                                    },
                                },
                            },
                        },
                    },
                },
            },
        }
    } else {
        // Update existing deployment
        deploy.Spec.Template.Spec.Containers[0].Image = deployment.Image
        deploy.Spec.Replicas = &deployment.Replicas
    }
    
    // Apply deployment
    _, err = ds.clientset.AppsV1().Deployments(deployment.Namespace).Apply(ctx, deploy, metav1.ApplyOptions{})
    if err != nil {
        return err
    }
    
    // Wait for rollout
    if err := ds.waitForRollout(ctx, deployment); err != nil {
        // Auto rollback
        if deployment.AutoRollback {
            return ds.rollback(ctx, deployment)
        }
        return err
    }
    
    deployment.Status = "success"
    deployment.DeployedAt = time.Now()
    ds.deploymentRepo.Update(ctx, deployment)
    
    ds.eventBus.Publish(Event{
        Type: "deployment.completed",
        Data: deployment,
    })
    
    return nil
}
```

### Week 6: Rollback & Advanced Strategies

#### Rollback Implementation

```go
func (ds *DeploymentService) rollback(ctx context.Context, deployment *Deployment) error {
    // Get deployment history
    history, _ := ds.clientset.AppsV1().Deployments(deployment.Namespace).GetRollback(
        ctx, deployment.Name, metav1.GetOptions{})
    
    // Rollback to previous revision
    rollback := &appsv1.DeploymentRollback{
        Name: deployment.Name,
        RollbackTo: &appsv1.RollbackConfig{
            Revision: history.Revision - 1,
        },
    }
    
    _, err := ds.clientset.AppsV1().Deployments(deployment.Namespace).Rollback(ctx, rollback, metav1.CreateOptions{})
    if err != nil {
        return err
    }
    
    deployment.Status = "rolled_back"
    ds.deploymentRepo.Update(ctx, deployment)
    
    ds.eventBus.Publish(Event{
        Type: "deployment.rolled_back",
        Data: deployment,
    })
    
    // Record in history
    historyRecord := &DeploymentHistory{
        DeploymentID: deployment.ID,
        Action:       "rollback",
        Status:       "success",
        Message:      "Auto-rollback due to failed deployment",
    }
    ds.historyRepo.Create(ctx, historyRecord)
    
    return nil
}

func (ds *DeploymentService) blueGreenDeploy(ctx context.Context, deployment *Deployment) error {
    // Deploy to green environment
    greenDeployment := deployment.Name + "-green"
    
    // Create/update green deployment
    err := ds.createDeployment(ctx, greenDeployment, deployment.Image, deployment.Replicas)
    if err != nil {
        return err
    }
    
    // Wait for health check
    if err := ds.waitForRollout(ctx, deployment); err != nil {
        return err
    }
    
    // Switch traffic (update service selector)
    service, _ := ds.clientset.CoreV1().Services(deployment.Namespace).Get(ctx, deployment.Name, metav1.GetOptions{})
    service.Spec.Selector["version"] = "green"
    _, err = ds.clientset.CoreV1().Services(deployment.Namespace).Update(ctx, service, metav1.UpdateOptions{})
    if err != nil {
        return err
    }
    
    // Delete blue deployment
    ds.clientset.AppsV1().Deployments(deployment.Namespace).Delete(ctx, deployment.Name+"-blue", metav1.DeleteOptions{})
    
    // Rename green to blue
    ds.clientset.AppsV1().Deployments(deployment.Namespace).Delete(ctx, deployment.Name, metav1.DeleteOptions{})
    greenDeploy, _ := ds.clientset.AppsV1().Deployments(deployment.Namespace).Get(ctx, greenDeployment, metav1.GetOptions{})
    greenDeploy.Name = deployment.Name
    _, err = ds.clientset.AppsV1().Deployments(deployment.Namespace).Create(ctx, greenDeploy, metav1.CreateOptions{})
    
    return err
}
```

---

## Phase 4: Frontend Development (Weeks 7-8)

### Week 7: Core UI Components

#### Overview Dashboard

```jsx
// features/overview/OverviewPage.jsx
import {useState, useEffect} from 'react';
import {statsAPI, buildsAPI, deploymentsAPI} from '../../api';

export function OverviewPage() {
    const [stats, setStats] = useState(null);
    const [recentBuilds, setRecentBuilds] = useState([]);
    const [recentDeployments, setRecentDeployments] = useState([]);
    
    useEffect(() => {
        loadData();
        setupWebSocket();
    }, []);
    
    const loadData = async () => {
        const [statsRes, buildsRes, deploymentsRes] = await Promise.all([
            statsAPI.getStats(),
            buildsAPI.getRecent(10),
            deploymentsAPI.getRecent(10),
        ]);
        
        setStats(statsRes.data);
        setRecentBuilds(buildsRes.data);
        setRecentDeployments(deploymentsRes.data);
    };
    
    const setupWebSocket = () => {
        const ws = new WebSocket(WS_URL);
        
        ws.onmessage = (event) => {
            const data = JSON.parse(event.data);
            switch (data.type) {
                case 'build.completed':
                    loadData(); // Refresh
                    break;
                case 'deployment.completed':
                    loadData(); // Refresh
                    break;
            }
        };
        
        return () => ws.close();
    };
    
    return (
        <div className="overview-page">
            <div className="stats-grid">
                <StatCard title="Total Pipelines" value={stats?.pipelines} />
                <StatCard title="Total Builds" value={stats?.builds} />
                <StatCard title="Success Rate" value={`${stats?.successRate}%`} />
                <StatCard title="Deployments" value={stats?.deployments} />
            </div>
            
            <div className="recent-activity">
                <BuildsTable builds={recentBuilds} />
                <DeploymentsTable deployments={recentDeployments} />
            </div>
        </div>
    );
}
```

#### Pipeline Configuration UI

```jsx
// features/pipelines/PipelineConfig.jsx
export function PipelineConfig({pipeline, onSubmit}) {
    const [config, setConfig] = useState({
        name: pipeline?.name || '',
        trigger_type: pipeline?.trigger_type || 'webhook',
        cron_schedule: pipeline?.cron_schedule || '',
        environment: pipeline?.environment || 'production',
        steps: pipeline?.config?.steps || [],
    });
    
    const addStep = () => {
        setConfig({
            ...config,
            steps: [...config.steps, {
                name: '',
                image: '',
                commands: [],
            }],
        });
    };
    
    const updateStep = (index, field, value) => {
        const steps = [...config.steps];
        steps[index][field] = value;
        setConfig({...config, steps});
    };
    
    return (
        <form onSubmit={handleSubmit}>
            <Input
                label="Pipeline Name"
                value={config.name}
                onChange={e => setConfig({...config, name: e.target.value})}
            />
            
            <Select
                label="Trigger Type"
                value={config.trigger_type}
                onChange={e => setConfig({...config, trigger_type: e.target.value})}
            >
                <option value="webhook">Webhook</option>
                <option value="schedule">Schedule</option>
                <option value="manual">Manual</option>
                <option value="api">API</option>
            </Select>
            
            {config.trigger_type === 'schedule' && (
                <Input
                    label="Cron Schedule"
                    value={config.cron_schedule}
                    onChange={e => setConfig({...config, cron_schedule: e.target.value})}
                    placeholder="0 0 * * *"
                />
            )}
            
            <div className="steps">
                <h3>Build Steps</h3>
                {config.steps.map((step, index) => (
                    <StepEditor
                        key={index}
                        step={step}
                        onChange={(field, value) => updateStep(index, field, value)}
                    />
                ))}
                <button type="button" onClick={addStep}>Add Step</button>
            </div>
            
            <button type="submit">Save Pipeline</button>
        </form>
    );
}
```

### Week 8: Advanced UI Features

#### Build Logs Viewer

```jsx
// features/builds/BuildLogs.jsx
import {useState, useEffect, useRef} from 'react';

export function BuildLogs({buildId}) {
    const [logs, setLogs] = useState([]);
    const [isLive, setIsLive] = useState(false);
    const logsEndRef = useRef(null);
    
    useEffect(() => {
        // Fetch initial logs
        fetchLogs();
        
        // Connect to WebSocket for live updates
        if (isLive) {
            const ws = new WebSocket(`${WS_URL}/builds/${buildId}/logs`);
            ws.onmessage = (event) => {
                const log = JSON.parse(event.data);
                setLogs(prev => [...prev, log]);
            };
            return () => ws.close();
        }
    }, [buildId, isLive]);
    
    useEffect(() => {
        logsEndRef.current?.scrollIntoView({behavior: 'smooth'});
    }, [logs]);
    
    return (
        <div className="build-logs">
            <div className="toolbar">
                <button onClick={() => setIsLive(!isLive)}>
                    {isLive ? '⏸ Pause' : '▶ Live'}
                </button>
                <button onClick={() => setLogs([])}>Clear</button>
            </div>
            
            <div className="logs-container">
                {logs.map((log, index) => (
                    <div key={index} className={`log-line log-${log.level}`}>
                        <span className="timestamp">{log.timestamp}</span>
                        <span className="message">{log.message}</span>
                    </div>
                ))}
                <div ref={logsEndRef} />
            </div>
        </div>
    );
}
```

---

## Phase 5: Testing & Documentation (Week 9)

### Backend Tests

```go
// tests/pipeline_test.go
func TestPipelineExecution(t *testing.T) {
    // Setup
    ctx := context.Background()
    pipeline := &Pipeline{
        Name: "Test Pipeline",
        Config: PipelineConfig{
            Steps: []PipelineStep{
                {
                    Name:    "Build",
                    Image:   "golang:1.21",
                    Commands: []string{"go build", "go test"},
                },
            },
        },
    }
    
    // Execute
    build := &Build{
        PipelineID: pipeline.ID,
        Status:     "pending",
    }
    
    executor := NewPipelineExecutor(dockerClient, eventBus, logger)
    err := executor.Execute(ctx, pipeline, build)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "success", build.Status)
    assert.NotNil(t, build.FinishedAt)
}

// Test webhook handling
func TestGitHubWebhook(t *testing.T) {
    handler := NewWebhookHandler(repo, gitClient)
    
    // Mock payload
    payload := []byte(`{"ref": "refs/heads/main"}`)
    signature := hmac.New(sha256.New, []byte(secret))
    
    req := httptest.NewRequest("POST", "/webhooks/github", bytes.NewBuffer(payload))
    req.Header.Set("X-Hub-Signature-256", signature)
    
    w := httptest.NewRecorder()
    handler.HandleGitHubWebhook(w, req)
    
    assert.Equal(t, 200, w.Code)
}
```

### Integration Tests

```go
// tests/integration_test.go
func TestFullPipeline(t *testing.T) {
    // 1. Create repository
    repo := createTestRepository(t)
    
    // 2. Create pipeline
    pipeline := createTestPipeline(t, repo.ID)
    
    // 3. Trigger build via webhook
    build := triggerBuild(t, pipeline.ID, "main")
    
    // 4. Wait for completion
    build = waitForBuildCompletion(t, build.ID)
    
    assert.Equal(t, "success", build.Status)
    
    // 5. Deploy
    deployment := deployBuild(t, build.ID, "production")
    
    // 6. Verify deployment
    deployment = waitForDeployment(t, deployment.ID)
    assert.Equal(t, "success", deployment.Status)
}
```

---

## Phase 6: Production Readiness (Week 10)

### Monitoring Setup

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'infradeploy'
    static_configs:
      - targets: ['localhost:9090']
```

### CI/CD Pipeline

```yaml
# .github/workflows/ci.yml
name: CI/CD

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.24'
      - run: go test -v -race -coverprofile=coverage.out ./...
      - uses: codecov/codecov-action@v3
  
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: golangci/golangci-lint-action@v3
  
  build:
    needs: [test, lint]
    runs-on: ubuntu-latest
    steps:
      - uses: docker/build-push-action@v4
        with:
          push: true
          tags: infradeploy:${{ github.sha }}
```

### Kubernetes Deployment

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: infradeploy-backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: infradeploy-backend
  template:
    spec:
      containers:
      - name: backend
        image: infrapilot/infradeploy:v1.0.0
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: infradeploy-secrets
              key: database-url
        resources:
          requests:
            cpu: "250m"
            memory: "512Mi"
          limits:
            cpu: "500m"
            memory: "1Gi"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 5
```

---

## Testing Strategy

### Unit Tests (80%+ Coverage)
- [ ] Git provider clients
- [ ] Pipeline executor
- [ ] Deployment strategies
- [ ] Webhook handlers
- [ ] Authentication/Authorization

### Integration Tests
- [ ] End-to-end pipeline execution
- [ ] Webhook → Build → Deploy flow
- [ ] Rollback scenarios
- [ ] Multi-environment deployments

### Load Tests
- [ ] 1000 concurrent builds
- [ ] 100 webhook events/second
- [ ] API latency under load

---

## Documentation Plan

### User Documentation
- [ ] Installation guide
- [ ] Quick start tutorial
- [ ] Pipeline configuration guide
- [ ] Deployment strategies guide
- [ ] Troubleshooting guide

### API Documentation
- [ ] OpenAPI/Swagger spec
- [ ] Authentication guide
- [ ] Webhook integration guide
- [ ] SDK examples (Go, Python, JavaScript)

### Developer Documentation
- [ ] Architecture overview
- [ ] Contributing guide
- [ ] Testing guide
- [ ] Deployment guide

---

## Marketing & Launch

### Launch Checklist
- [ ] GitHub repository public
- [ ] Demo deployed and accessible
- [ ] Blog posts published
- [ ] Social media announcements
- [ ] Product Hunt launch
- [ ] Hacker News submission
- [ ] Reddit posts (r/golang, r/devops, r/selfhosted)

### Content Marketing
- [ ] "How I Built InfraDeploy" blog series
- [ ] Demo video (10 minutes)
- [ ] Architecture deep-dive article
- [ ] Comparison with Jenkins/GitLab CI

---

## Post-Launch

### Week 1-2: Bug Fixes
- [ ] Monitor GitHub issues
- [ ] Fix critical bugs
- [ ] Release v1.0.1 hotfix

### Month 1-2: Feature Updates
- [ ] Add Jenkins integration
- [ ] Add GitHub Actions support
- [ ] Implement canary analysis
- [ ] Add deployment approvals
- [ ] Telegram/Discord notifications

### Month 3-6: Scale & Polish
- [ ] Multi-cluster support
- [ ] Advanced RBAC
- [ ] Plugin marketplace
- [ ] Mobile-responsive UI
- [ ] Performance optimizations

---

## Success Metrics

### Technical
- [ ] API latency < 100ms
- [ ] Build execution time < 5 minutes (average)
- [ ] Deployment time < 30 seconds
- [ ] 99.9% uptime
- [ ] < 1% failed builds due to platform issues

### Community
- [ ] 500+ GitHub stars
- [ ] 50+ forks
- [ ] 20+ contributors
- [ ] Active Discord/Slack community
- [ ] Blog posts with 10K+ views

---

## Team & Resources

### Required Skills
- **Backend:** Go, PostgreSQL, Redis, WebSocket
- **Frontend:** React, Axios, WebSocket
- **DevOps:** Docker, Kubernetes, Helm, CI/CD
- **Git Integration:** GitHub/GitLab/Bitbucket APIs

### Tools & Services
- GitHub (code hosting, CI/CD)
- DigitalOcean/AWS (hosting)
- Docker Hub (image registry)
- Namecheap (domain)

---

<p align="center">
  Ready to build InfraDeploy! 🚀
</p>