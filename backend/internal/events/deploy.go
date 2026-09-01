package events

import "time"

// GitRepositoryCreatedEvent is published when a git repository is created
type GitRepositoryCreatedEvent struct {
	RepositoryID string
	RepoName     string
	Provider     string
	Organization string
	CreatedAt    time.Time
}

func (e GitRepositoryCreatedEvent) Name() string {
	return "deploy.git_repository.created"
}

func (e GitRepositoryCreatedEvent) Timestamp() time.Time {
	return e.CreatedAt
}

// PipelineCreatedEvent is published when a pipeline is created
type PipelineCreatedEvent struct {
	PipelineID   string
	PipelineName string
	RepositoryID string
	Environment  string
	CreatedAt    time.Time
}

func (e PipelineCreatedEvent) Name() string {
	return "deploy.pipeline.created"
}

func (e PipelineCreatedEvent) Timestamp() time.Time {
	return e.CreatedAt
}

// BuildCreatedEvent is published when a build is created
type BuildCreatedEvent struct {
	BuildID     string
	PipelineID  string
	BuildNumber int
	Status      string
	CreatedAt   time.Time
}

func (e BuildCreatedEvent) Name() string {
	return "deploy.build.created"
}

func (e BuildCreatedEvent) Timestamp() time.Time {
	return e.CreatedAt
}

// BuildStatusChangedEvent is published when a build status changes
type BuildStatusChangedEvent struct {
	BuildID   string
	Status    string
	Duration  int64
	CreatedAt time.Time
}

func (e BuildStatusChangedEvent) Name() string {
	return "deploy.build.status_changed"
}

func (e BuildStatusChangedEvent) Timestamp() time.Time {
	return e.CreatedAt
}

// DeploymentCreatedEvent is published when a deployment is created
type DeploymentCreatedEvent struct {
	DeploymentID string
	DeployName   string
	Environment  string
	Image        string
	CreatedAt    time.Time
}

func (e DeploymentCreatedEvent) Name() string {
	return "deploy.deployment.created"
}

func (e DeploymentCreatedEvent) Timestamp() time.Time {
	return e.CreatedAt
}

// DeploymentRolledBackEvent is published when a deployment is rolled back
type DeploymentRolledBackEvent struct {
	DeploymentID string
	Environment  string
	Image        string
	CreatedAt    time.Time
}

func (e DeploymentRolledBackEvent) Name() string {
	return "deploy.deployment.rolled_back"
}

func (e DeploymentRolledBackEvent) Timestamp() time.Time {
	return e.CreatedAt
}

// WebhookReceivedEvent is published when a webhook is received
type WebhookReceivedEvent struct {
	WebhookID  string
	Provider   string
	EventType  string
	ReceivedAt time.Time
}

func (e WebhookReceivedEvent) Name() string {
	return "deploy.webhook.received"
}

func (e WebhookReceivedEvent) Timestamp() time.Time {
	return e.ReceivedAt
}
