package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GitRepository represents a connected git repository
type GitRepository struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name         string         `gorm:"not null" json:"name"`
	Provider     string         `gorm:"not null" json:"provider"` // github, gitlab, bitbucket, gitea
	URL          string         `gorm:"not null" json:"url"`
	CloneURL     string         `json:"-"` // Contains credentials, never send to client
	Branch       string         `gorm:"default:main" json:"branch"`
	AuthType     string         `gorm:"not null" json:"auth_type"` // oauth, ssh, token
	AuthToken    string         `json:"-"`                         // Encrypted credential
	SSHKey       string         `json:"-"`                         // Private SSH key if applicable
	WebhookURL   string         `json:"webhook_url,omitempty"`
	WebhookToken string         `json:"-"` // For webhook verification
	Organization string         `gorm:"default:Default Organization" json:"organization"`
	CreatedBy    uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	LastSyncAt   time.Time      `json:"last_sync_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (g *GitRepository) BeforeCreate(tx *gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return nil
}

// Build represents a CI pipeline build execution
type Build struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	PipelineID    uuid.UUID  `gorm:"type:uuid;not null" json:"pipeline_id"`
	Pipeline      Pipeline   `gorm:"foreignKey:PipelineID" json:"pipeline,omitempty"`
	BuildNumber   int        `gorm:"not null" json:"build_number"`
	CommitHash    string     `gorm:"not null" json:"commit_hash"`
	CommitMessage string     `gorm:"size:1024" json:"commit_message"`
	Author        string     `json:"author"`
	Branch        string     `json:"branch"`
	Tag           string     `json:"tag,omitempty"`
	Status        string     `gorm:"not null;default:pending" json:"status"` // pending, running, success, failed, cancelled
	Trigger       string     `gorm:"not null" json:"trigger"`                // manual, webhook, schedule, api
	StartedAt     *time.Time `json:"started_at"`
	FinishedAt    *time.Time `json:"finished_at"`
	DurationMs    int64      `json:"duration_ms"`
	Logs          string     `gorm:"type:text" json:"logs,omitempty"`
	Artifacts     []byte     `gorm:"type:jsonb" json:"artifacts,omitempty"` // JSON array of artifact URLs
	ErrorMessage  string     `gorm:"size:2048" json:"error_message,omitempty"`
	Organization  string     `gorm:"default:Default Organization" json:"organization"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (b *Build) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// Pipeline represents a CI/CD pipeline configuration
type Pipeline struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name         string         `gorm:"not null" json:"name"`
	Description  string         `json:"description"`
	RepositoryID uuid.UUID      `gorm:"type:uuid;not null" json:"repository_id"`
	Repository   GitRepository  `gorm:"foreignKey:RepositoryID" json:"repository,omitempty"`
	Config       string         `gorm:"type:jsonb;not null" json:"config"`            // Pipeline YAML/JSON config
	TriggerType  string         `gorm:"not null;default:webhook" json:"trigger_type"` // webhook, schedule, manual
	Schedule     string         `json:"schedule,omitempty"`                           // Cron expression for scheduled builds
	Enabled      bool           `gorm:"default:true" json:"enabled"`
	Environment  string         `gorm:"default:development" json:"environment"` // development, staging, production
	CreatedBy    uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	LastRunAt    *time.Time     `json:"last_run_at"`
	Organization string         `gorm:"default:Default Organization" json:"organization"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *Pipeline) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// Deployment represents a Kubernetes/Docker deployment
type Deployment struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Name         string     `gorm:"not null" json:"name"`
	Environment  string     `gorm:"not null" json:"environment"`     // development, staging, production, canary
	Namespace    string     `json:"namespace"`                       // K8s namespace
	Cluster      string     `json:"cluster"`                         // K8s cluster name
	Strategy     string     `gorm:"default:rolling" json:"strategy"` // rolling, blue-green, canary, recreate
	BuildID      *uuid.UUID `gorm:"type:uuid" json:"build_id,omitempty"`
	Build        *Build     `gorm:"foreignKey:BuildID" json:"build,omitempty"`
	Image        string     `gorm:"not null" json:"image"` // Docker image:tag
	Replicas     int        `gorm:"default:1" json:"replicas"`
	Status       string     `gorm:"not null;default:pending" json:"status"` // pending, deploying, running, failed, rolled-back
	HealthCheck  *string    `json:"health_check,omitempty"`                 // Health check endpoint
	AutoRollback bool       `gorm:"default:true" json:"auto_rollback"`
	Config       []byte     `gorm:"type:jsonb" json:"config,omitempty"`    // Helm values, ArgoCD config
	Variables    []byte     `gorm:"type:jsonb" json:"variables,omitempty"` // Environment variables
	DeployedBy   uuid.UUID  `gorm:"type:uuid;not null" json:"deployed_by"`
	DeployedAt   *time.Time `json:"deployed_at"`
	RolledBackAt *time.Time `json:"rolled_back_at"`
	ErrorMessage string     `gorm:"size:2048" json:"error_message,omitempty"`
	Organization string     `gorm:"default:Default Organization" json:"organization"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (d *Deployment) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

// DeploymentHistory tracks deployment events and rollbacks
type DeploymentHistory struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	DeploymentID uuid.UUID `gorm:"type:uuid;not null" json:"deployment_id"`
	Action       string    `gorm:"not null" json:"action"` // deploy, rollback, scale, restart
	FromImage    string    `json:"from_image,omitempty"`
	ToImage      string    `json:"to_image,omitempty"`
	FromReplicas int       `json:"from_replicas,omitempty"`
	ToReplicas   int       `json:"to_replicas,omitempty"`
	Reason       string    `json:"reason,omitempty"`
	TriggeredBy  uuid.UUID `gorm:"type:uuid;not null" json:"triggered_by"`
	TriggeredAt  time.Time `gorm:"not null" json:"triggered_at"`
	Organization string    `gorm:"default:Default Organization" json:"organization"`
}

func (h *DeploymentHistory) BeforeCreate(tx *gorm.DB) error {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	return nil
}

// WebhookEvent represents an incoming webhook from a Git provider
type WebhookEvent struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Provider     string     `gorm:"not null" json:"provider"`   // github, gitlab, bitbucket, gitea
	EventType    string     `gorm:"not null" json:"event_type"` // push, pull_request, tag, release
	Signature    string     `gorm:"size:512" json:"-"`          // Webhook signature for verification
	Payload      []byte     `gorm:"type:jsonb;not null" json:"payload"`
	Processed    bool       `gorm:"default:false" json:"processed"`
	Error        string     `gorm:"size:1024" json:"error,omitempty"`
	ReceivedAt   time.Time  `gorm:"not null" json:"received_at"`
	ProcessedAt  *time.Time `json:"processed_at,omitempty"`
	RepositoryID *uuid.UUID `gorm:"type:uuid" json:"repository_id,omitempty"`
}

func (w *WebhookEvent) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}

// DeploymentStatus constants
const (
	DeploymentStatusPending    = "pending"
	DeploymentStatusDeploying  = "deploying"
	DeploymentStatusRunning    = "running"
	DeploymentStatusFailed     = "failed"
	DeploymentStatusRolledBack = "rolled-back"
)

// BuildStatus constants
const (
	BuildStatusPending   = "pending"
	BuildStatusRunning   = "running"
	BuildStatusSuccess   = "success"
	BuildStatusFailed    = "failed"
	BuildStatusCancelled = "cancelled"
)

// DeploymentStrategy constants
const (
	DeploymentStrategyRolling   = "rolling"
	DeploymentStrategyBlueGreen = "blue-green"
	DeploymentStrategyCanary    = "canary"
	DeploymentStrategyRecreate  = "recreate"
)

// GitProvider constants
const (
	GitProviderGitHub      = "github"
	GitProviderGitLab      = "gitlab"
	GitProviderBitbucket   = "bitbucket"
	GitProviderGitea       = "gitea"
	GitProviderAzureDevOps = "azure-devops"
)

// AuthType constants
const (
	AuthTypeOAuth = "oauth"
	AuthTypeSSH   = "ssh"
	AuthTypeToken = "token"
	AuthTypeBasic = "basic"
)

// TriggerType constants
const (
	TriggerTypeWebhook  = "webhook"
	TriggerTypeSchedule = "schedule"
	TriggerTypeManual   = "manual"
	TriggerTypeAPI      = "api"
)
