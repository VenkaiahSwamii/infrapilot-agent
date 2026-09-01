package services

import (
	"encoding/json"
	"fmt"
	"time"

	"infrapilot/backend/internal/events"
	"infrapilot/backend/internal/models"
	"infrapilot/backend/internal/repository"

	"github.com/google/uuid"
)

type DeployService struct {
	deployRepo *repository.DeployRepository
	eventBus   *events.EventBus
}

func NewDeployService(deployRepo *repository.DeployRepository, eventBus *events.EventBus) *DeployService {
	return &DeployService{
		deployRepo: deployRepo,
		eventBus:   eventBus,
	}
}

// Git Repository Operations

func (s *DeployService) CreateGitRepository(userID uuid.UUID, org string, name, provider, url, cloneURL, branch, authType, authToken, sshKey, webhookURL, webhookToken string) (*models.GitRepository, error) {
	repo := &models.GitRepository{
		Name:         name,
		Provider:     provider,
		URL:          url,
		CloneURL:     cloneURL,
		Branch:       branch,
		AuthType:     authType,
		AuthToken:    authToken,
		SSHKey:       sshKey,
		WebhookURL:   webhookURL,
		WebhookToken: webhookToken,
		Organization: org,
		CreatedBy:    userID,
		LastSyncAt:   time.Now(),
	}

	if err := s.deployRepo.CreateGitRepository(repo); err != nil {
		return nil, fmt.Errorf("failed to create git repository: %w", err)
	}

	// Emit event
	s.eventBus.Publish(GitRepositoryCreatedEvent{
		RepositoryID: repo.ID.String(),
		RepoName:     repo.Name,
		Provider:     repo.Provider,
		Organization: org,
		CreatedAt:    time.Now(),
	})

	return repo, nil
}

func (s *DeployService) GetGitRepository(id uuid.UUID) (*models.GitRepository, error) {
	return s.deployRepo.FindGitRepositoryByID(id)
}

func (s *DeployService) ListGitRepositories(org string) ([]models.GitRepository, error) {
	return s.deployRepo.FindGitRepositoriesByOrg(org)
}

func (s *DeployService) UpdateGitRepository(id uuid.UUID, name, branch, webhookURL string) (*models.GitRepository, error) {
	repo, err := s.deployRepo.FindGitRepositoryByID(id)
	if err != nil {
		return nil, err
	}

	repo.Name = name
	repo.Branch = branch
	repo.WebhookURL = webhookURL

	if err := s.deployRepo.UpdateGitRepository(repo); err != nil {
		return nil, fmt.Errorf("failed to update git repository: %w", err)
	}

	return repo, nil
}

func (s *DeployService) DeleteGitRepository(id uuid.UUID) error {
	return s.deployRepo.DeleteGitRepository(id)
}

// Pipeline Operations

func (s *DeployService) CreatePipeline(userID uuid.UUID, org, name, description string, repoID uuid.UUID, config interface{}, triggerType, schedule, environment string) (*models.Pipeline, error) {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal pipeline config: %w", err)
	}

	pipeline := &models.Pipeline{
		Name:         name,
		Description:  description,
		RepositoryID: repoID,
		Config:       string(configJSON),
		TriggerType:  triggerType,
		Schedule:     schedule,
		Enabled:      true,
		Environment:  environment,
		CreatedBy:    userID,
		Organization: org,
	}

	if err := s.deployRepo.CreatePipeline(pipeline); err != nil {
		return nil, fmt.Errorf("failed to create pipeline: %w", err)
	}

	s.eventBus.Publish(PipelineCreatedEvent{
		PipelineID:   pipeline.ID.String(),
		PipelineName: pipeline.Name,
		RepositoryID: repoID.String(),
		Environment:  environment,
		CreatedAt:    time.Now(),
	})

	return pipeline, nil
}

func (s *DeployService) GetPipeline(id uuid.UUID) (*models.Pipeline, error) {
	return s.deployRepo.FindPipelineByID(id)
}

func (s *DeployService) ListPipelines(org, environment string) ([]models.Pipeline, error) {
	return s.deployRepo.FindPipelinesByOrg(org, environment)
}

func (s *DeployService) UpdatePipeline(id uuid.UUID, name, description, configStr, schedule string, enabled bool) (*models.Pipeline, error) {
	pipeline, err := s.deployRepo.FindPipelineByID(id)
	if err != nil {
		return nil, err
	}

	pipeline.Name = name
	pipeline.Description = description
	pipeline.Config = configStr
	pipeline.Schedule = schedule
	pipeline.Enabled = enabled

	if err := s.deployRepo.UpdatePipeline(pipeline); err != nil {
		return nil, fmt.Errorf("failed to update pipeline: %w", err)
	}

	return pipeline, nil
}

func (s *DeployService) DeletePipeline(id uuid.UUID) error {
	return s.deployRepo.DeletePipeline(id)
}

// Build Operations

func (s *DeployService) CreateBuild(pipelineID uuid.UUID, commitHash, commitMessage, author, branch, tag, trigger string) (*models.Build, error) {
	buildNumber, err := s.deployRepo.GetNextBuildNumber(pipelineID)
	if err != nil {
		return nil, fmt.Errorf("failed to get next build number: %w", err)
	}

	pipeline, _ := s.deployRepo.FindPipelineByID(pipelineID)
	org := "Default Organization"
	if pipeline != nil {
		org = pipeline.Organization
	}

	build := &models.Build{
		PipelineID:    pipelineID,
		BuildNumber:   buildNumber,
		CommitHash:    commitHash,
		CommitMessage: commitMessage,
		Author:        author,
		Branch:        branch,
		Tag:           tag,
		Status:        models.BuildStatusPending,
		Trigger:       trigger,
		Organization:  org,
	}

	if err := s.deployRepo.CreateBuild(build); err != nil {
		return nil, fmt.Errorf("failed to create build: %w", err)
	}

	// Emit build created event
	s.eventBus.Publish(BuildCreatedEvent{
		BuildID:     build.ID.String(),
		PipelineID:  pipelineID.String(),
		BuildNumber: buildNumber,
		Status:      build.Status,
		CreatedAt:   time.Now(),
	})

	return build, nil
}

func (s *DeployService) GetBuild(id uuid.UUID) (*models.Build, error) {
	return s.deployRepo.FindBuildByID(id)
}

func (s *DeployService) GetBuildsByPipeline(pipelineID uuid.UUID, limit int) ([]models.Build, error) {
	return s.deployRepo.FindBuildsByPipeline(pipelineID, limit)
}

func (s *DeployService) GetBuildsByOrg(org string, limit int) ([]models.Build, error) {
	return s.deployRepo.FindBuildsByOrg(org, limit)
}

func (s *DeployService) UpdateBuildStatus(id uuid.UUID, status, errorMessage string) (*models.Build, error) {
	build, err := s.deployRepo.FindBuildByID(id)
	if err != nil {
		return nil, err
	}

	build.Status = status
	build.ErrorMessage = errorMessage

	if status == models.BuildStatusRunning {
		now := time.Now()
		build.StartedAt = &now
	} else if status == models.BuildStatusSuccess || status == models.BuildStatusFailed || status == models.BuildStatusCancelled {
		now := time.Now()
		build.FinishedAt = &now
		if build.StartedAt != nil {
			build.DurationMs = now.Sub(*build.StartedAt).Milliseconds()
		}
	}

	if err := s.deployRepo.UpdateBuild(build); err != nil {
		return nil, fmt.Errorf("failed to update build status: %w", err)
	}

	s.eventBus.Publish(BuildStatusChangedEvent{
		BuildID:   build.ID.String(),
		Status:    status,
		Duration:  build.DurationMs,
		CreatedAt: time.Now(),
	})

	return build, nil
}

// Deployment Operations

func (s *DeployService) CreateDeployment(userID uuid.UUID, org, name, environment, namespace, cluster, strategy, image string, replicas int, buildID *uuid.UUID, autoRollback bool, config, variables map[string]interface{}) (*models.Deployment, error) {
	configJSON, _ := json.Marshal(config)
	variablesJSON, _ := json.Marshal(variables)

	deploy := &models.Deployment{
		Name:         name,
		Environment:  environment,
		Namespace:    namespace,
		Cluster:      cluster,
		Strategy:     strategy,
		BuildID:      buildID,
		Image:        image,
		Replicas:     replicas,
		Status:       models.DeploymentStatusPending,
		AutoRollback: autoRollback,
		Config:       configJSON,
		Variables:    variablesJSON,
		DeployedBy:   userID,
		Organization: org,
	}

	if err := s.deployRepo.CreateDeployment(deploy); err != nil {
		return nil, fmt.Errorf("failed to create deployment: %w", err)
	}

	// Create deployment history
	history := &models.DeploymentHistory{
		DeploymentID: deploy.ID,
		Action:       "deploy",
		ToImage:      image,
		ToReplicas:   replicas,
		Reason:       "Initial deployment",
		TriggeredBy:  userID,
		Organization: org,
	}
	s.deployRepo.CreateDeploymentHistory(history)

	s.eventBus.Publish(DeploymentCreatedEvent{
		DeploymentID: deploy.ID.String(),
		DeployName:   deploy.Name,
		Environment:  environment,
		Image:        image,
		CreatedAt:    time.Now(),
	})

	return deploy, nil
}

func (s *DeployService) GetDeployment(id uuid.UUID) (*models.Deployment, error) {
	return s.deployRepo.FindDeploymentByID(id)
}

func (s *DeployService) ListDeployments(org, environment string, limit int) ([]models.Deployment, error) {
	if environment != "" {
		return s.deployRepo.FindDeploymentsByEnv(org, environment, limit)
	}
	return s.deployRepo.FindDeploymentsByOrg(org, limit)
}

func (s *DeployService) UpdateDeploymentStatus(id uuid.UUID, status, errorMessage string) (*models.Deployment, error) {
	deploy, err := s.deployRepo.FindDeploymentByID(id)
	if err != nil {
		return nil, err
	}

	deploy.Status = status
	deploy.ErrorMessage = errorMessage

	if status == models.DeploymentStatusRunning {
		now := time.Now()
		deploy.DeployedAt = &now
	} else if status == models.DeploymentStatusRolledBack {
		now := time.Now()
		deploy.RolledBackAt = &now
	}

	if err := s.deployRepo.UpdateDeployment(deploy); err != nil {
		return nil, fmt.Errorf("failed to update deployment status: %w", err)
	}

	return deploy, nil
}

func (s *DeployService) RollbackDeployment(id uuid.UUID, userID uuid.UUID, reason string) (*models.Deployment, error) {
	deploy, err := s.deployRepo.FindDeploymentByID(id)
	if err != nil {
		return nil, err
	}

	// Record rollback
	history := &models.DeploymentHistory{
		DeploymentID: deploy.ID,
		Action:       "rollback",
		FromImage:    deploy.Image,
		Reason:       reason,
		TriggeredBy:  userID,
		Organization: deploy.Organization,
	}
	s.deployRepo.CreateDeploymentHistory(history)

	deploy.Status = models.DeploymentStatusRolledBack
	now := time.Now()
	deploy.RolledBackAt = &now

	if err := s.deployRepo.UpdateDeployment(deploy); err != nil {
		return nil, fmt.Errorf("failed to rollback deployment: %w", err)
	}

	s.eventBus.Publish(DeploymentRolledBackEvent{
		DeploymentID: deploy.ID.String(),
		Environment:  deploy.Environment,
		Image:        deploy.Image,
		CreatedAt:    time.Now(),
	})

	return deploy, nil
}

// Webhook Operations

func (s *DeployService) ProcessWebhookEvent(provider, eventType, signature string, payload []byte, repoID *uuid.UUID) error {
	event := &models.WebhookEvent{
		Provider:     provider,
		EventType:    eventType,
		Signature:    signature,
		Payload:      payload,
		RepositoryID: repoID,
	}

	if err := s.deployRepo.CreateWebhookEvent(event); err != nil {
		return fmt.Errorf("failed to create webhook event: %w", err)
	}

	s.eventBus.Publish(WebhookReceivedEvent{
		WebhookID:  event.ID.String(),
		Provider:   provider,
		EventType:  eventType,
		ReceivedAt: time.Now(),
	})

	return nil
}

func (s *DeployService) GetUnprocessedWebhookEvents(limit int) ([]models.WebhookEvent, error) {
	return s.deployRepo.FindUnprocessedWebhookEvents(limit)
}

func (s *DeployService) MarkWebhookProcessed(eventID uuid.UUID) error {
	return s.deployRepo.MarkWebhookProcessed(eventID)
}

// Statistics

func (s *DeployService) GetDeployStats(org string) (map[string]interface{}, error) {
	pipelineCount, _ := s.deployRepo.CountPipelinesByOrg(org)
	successBuilds, _ := s.deployRepo.CountBuildsByOrgAndStatus(org, models.BuildStatusSuccess)
	failedBuilds, _ := s.deployRepo.CountBuildsByOrgAndStatus(org, models.BuildStatusFailed)
	deploymentCount, _ := s.deployRepo.CountDeploymentsByOrg(org)

	return map[string]interface{}{
		"pipelines":      pipelineCount,
		"builds_success": successBuilds,
		"builds_failed":  failedBuilds,
		"deployments":    deploymentCount,
	}, nil
}
