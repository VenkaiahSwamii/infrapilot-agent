package repository

import (
	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeployRepository struct{}

func NewDeployRepository() *DeployRepository {
	return &DeployRepository{}
}

// Git Repository methods
func (r *DeployRepository) CreateGitRepository(repo *models.GitRepository) error {
	return database.DB.Create(repo).Error
}

func (r *DeployRepository) FindGitRepositoryByID(id uuid.UUID) (*models.GitRepository, error) {
	var repo models.GitRepository
	err := database.DB.First(&repo, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &repo, nil
}

func (r *DeployRepository) FindGitRepositoriesByOrg(org string) ([]models.GitRepository, error) {
	var repos []models.GitRepository
	err := database.DB.Where("organization = ?", org).
		Order("created_at DESC").
		Find(&repos).Error
	return repos, err
}

func (r *DeployRepository) FindGitRepositoriesByProvider(provider string) ([]models.GitRepository, error) {
	var repos []models.GitRepository
	err := database.DB.Where("provider = ?", provider).
		Order("created_at DESC").
		Find(&repos).Error
	return repos, err
}

func (r *DeployRepository) UpdateGitRepository(repo *models.GitRepository) error {
	return database.DB.Save(repo).Error
}

func (r *DeployRepository) DeleteGitRepository(id uuid.UUID) error {
	return database.DB.Delete(&models.GitRepository{}, id).Error
}

// Pipeline methods
func (r *DeployRepository) CreatePipeline(pipeline *models.Pipeline) error {
	return database.DB.Create(pipeline).Error
}

func (r *DeployRepository) FindPipelineByID(id uuid.UUID) (*models.Pipeline, error) {
	var pipeline models.Pipeline
	err := database.DB.Preload("Repository").First(&pipeline, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &pipeline, nil
}

func (r *DeployRepository) FindPipelinesByRepo(repoID uuid.UUID) ([]models.Pipeline, error) {
	var pipelines []models.Pipeline
	err := database.DB.Preload("Repository").
		Where("repository_id = ?", repoID).
		Order("created_at DESC").
		Find(&pipelines).Error
	return pipelines, err
}

func (r *DeployRepository) FindPipelinesByOrg(org string, env string) ([]models.Pipeline, error) {
	var pipelines []models.Pipeline
	query := database.DB.Preload("Repository").Where("organization = ?", org)
	if env != "" {
		query = query.Where("environment = ?", env)
	}
	err := query.Order("created_at DESC").Find(&pipelines).Error
	return pipelines, err
}

func (r *DeployRepository) UpdatePipeline(pipeline *models.Pipeline) error {
	return database.DB.Save(pipeline).Error
}

func (r *DeployRepository) DeletePipeline(id uuid.UUID) error {
	return database.DB.Delete(&models.Pipeline{}, id).Error
}

// Build methods
func (r *DeployRepository) CreateBuild(build *models.Build) error {
	return database.DB.Create(build).Error
}

func (r *DeployRepository) FindBuildByID(id uuid.UUID) (*models.Build, error) {
	var build models.Build
	err := database.DB.Preload("Pipeline").First(&build, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &build, nil
}

func (r *DeployRepository) FindBuildsByPipeline(pipelineID uuid.UUID, limit int) ([]models.Build, error) {
	var builds []models.Build
	query := database.DB.Preload("Pipeline").
		Where("pipeline_id = ?", pipelineID).
		Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&builds).Error
	return builds, err
}

func (r *DeployRepository) FindBuildsByOrg(org string, limit int) ([]models.Build, error) {
	var builds []models.Build
	query := database.DB.Preload("Pipeline").
		Where("organization = ?", org).
		Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&builds).Error
	return builds, err
}

func (r *DeployRepository) FindBuildsByCommit(commitHash string) ([]models.Build, error) {
	var builds []models.Build
	err := database.DB.Preload("Pipeline").
		Where("commit_hash = ?", commitHash).
		Order("created_at DESC").
		Find(&builds).Error
	return builds, err
}

func (r *DeployRepository) UpdateBuild(build *models.Build) error {
	return database.DB.Save(build).Error
}

func (r *DeployRepository) GetNextBuildNumber(pipelineID uuid.UUID) (int, error) {
	var maxBuild int
	err := database.DB.Model(&models.Build{}).
		Where("pipeline_id = ?", pipelineID).
		Select("COALESCE(MAX(build_number), 0)").
		Scan(&maxBuild).Error
	if err != nil {
		return 0, err
	}
	return maxBuild + 1, nil
}

// Deployment methods
func (r *DeployRepository) CreateDeployment(deploy *models.Deployment) error {
	return database.DB.Create(deploy).Error
}

func (r *DeployRepository) FindDeploymentByID(id uuid.UUID) (*models.Deployment, error) {
	var deploy models.Deployment
	err := database.DB.Preload("Build").First(&deploy, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &deploy, nil
}

func (r *DeployRepository) FindDeploymentsByEnv(org, environment string, limit int) ([]models.Deployment, error) {
	var deploys []models.Deployment
	query := database.DB.Preload("Build").
		Where("organization = ? AND environment = ?", org, environment).
		Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&deploys).Error
	return deploys, err
}

func (r *DeployRepository) FindDeploymentsByOrg(org string, limit int) ([]models.Deployment, error) {
	var deploys []models.Deployment
	query := database.DB.Preload("Build").
		Where("organization = ?", org).
		Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&deploys).Error
	return deploys, err
}

func (r *DeployRepository) UpdateDeployment(deploy *models.Deployment) error {
	return database.DB.Save(deploy).Error
}

func (r *DeployRepository) DeleteDeployment(id uuid.UUID) error {
	return database.DB.Delete(&models.Deployment{}, id).Error
}

// Deployment History methods
func (r *DeployRepository) CreateDeploymentHistory(history *models.DeploymentHistory) error {
	return database.DB.Create(history).Error
}

func (r *DeployRepository) FindDeploymentHistory(deploymentID uuid.UUID, limit int) ([]models.DeploymentHistory, error) {
	var history []models.DeploymentHistory
	query := database.DB.
		Where("deployment_id = ?", deploymentID).
		Order("triggered_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&history).Error
	return history, err
}

// Webhook methods
func (r *DeployRepository) CreateWebhookEvent(event *models.WebhookEvent) error {
	return database.DB.Create(event).Error
}

func (r *DeployRepository) FindUnprocessedWebhookEvents(limit int) ([]models.WebhookEvent, error) {
	var events []models.WebhookEvent
	err := database.DB.
		Where("processed = ?", false).
		Order("received_at ASC").
		Limit(limit).
		Find(&events).Error
	return events, err
}

func (r *DeployRepository) MarkWebhookProcessed(eventID uuid.UUID) error {
	now := gorm.DeletedAt{}
	return database.DB.Model(&models.WebhookEvent{}).
		Where("id = ?", eventID).
		Updates(map[string]interface{}{
			"processed":    true,
			"processed_at": now,
		}).Error
}

// Stats methods
func (r *DeployRepository) CountPipelinesByOrg(org string) (int64, error) {
	var count int64
	err := database.DB.Model(&models.Pipeline{}).Where("organization = ?", org).Count(&count).Error
	return count, err
}

func (r *DeployRepository) CountBuildsByOrgAndStatus(org, status string) (int64, error) {
	var count int64
	err := database.DB.Model(&models.Build{}).
		Where("organization = ? AND status = ?", org, status).
		Count(&count).Error
	return count, err
}

func (r *DeployRepository) CountDeploymentsByOrg(org string) (int64, error) {
	var count int64
	err := database.DB.Model(&models.Deployment{}).Where("organization = ?", org).Count(&count).Error
	return count, err
}
