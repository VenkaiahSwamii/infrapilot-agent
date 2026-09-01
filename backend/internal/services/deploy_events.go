package services

import (
	"time"
)

type GitRepositoryCreatedEvent struct {
	RepositoryID string
	RepoName     string
	Provider     string
	Organization string
	CreatedAt    time.Time
}

func (e GitRepositoryCreatedEvent) Name() string {
	return "git.repository.created"
}

func (e GitRepositoryCreatedEvent) Timestamp() time.Time {
	return e.CreatedAt
}

type PipelineCreatedEvent struct {
	PipelineID   string
	PipelineName string
	RepositoryID string
	Environment  string
	CreatedAt    time.Time
}

func (e PipelineCreatedEvent) Name() string {
	return "pipeline.created"
}

func (e PipelineCreatedEvent) Timestamp() time.Time {
	return e.CreatedAt
}

type BuildCreatedEvent struct {
	BuildID     string
	PipelineID  string
	BuildNumber int
	Branch      string
	CommitHash  string
	Status      string
	CreatedAt   time.Time
}

func (e BuildCreatedEvent) Name() string {
	return "build.created"
}

func (e BuildCreatedEvent) Timestamp() time.Time {
	return e.CreatedAt
}

type BuildStatusChangedEvent struct {
	BuildID   string
	Status    string
	Duration  int64
	CreatedAt time.Time
}

func (e BuildStatusChangedEvent) Name() string {
	return "build.status.changed"
}

func (e BuildStatusChangedEvent) Timestamp() time.Time {
	return e.CreatedAt
}

type DeploymentCreatedEvent struct {
	DeploymentID string
	DeployName   string
	Environment  string
	Image        string
	CreatedAt    time.Time
}

func (e DeploymentCreatedEvent) Name() string {
	return "deployment.created"
}

func (e DeploymentCreatedEvent) Timestamp() time.Time {
	return e.CreatedAt
}

type DeploymentRolledBackEvent struct {
	DeploymentID string
	Environment  string
	Image        string
	RollbackTo   string
	Reason       string
	CreatedAt    time.Time
}

func (e DeploymentRolledBackEvent) Name() string {
	return "deployment.rolled_back"
}

func (e DeploymentRolledBackEvent) Timestamp() time.Time {
	return e.CreatedAt
}

type WebhookReceivedEvent struct {
	WebhookID  string
	Provider   string
	EventType  string
	ReceivedAt time.Time
}

func (e WebhookReceivedEvent) Name() string {
	return "webhook.received"
}

func (e WebhookReceivedEvent) Timestamp() time.Time {
	return e.ReceivedAt
}
