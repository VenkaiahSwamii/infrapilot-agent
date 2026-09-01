package handlers

import (
	"net/http"
	"strconv"

	"infrapilot/backend/internal/models"
	"infrapilot/backend/internal/services"
	"infrapilot/backend/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DeployHandler struct {
	deployService *services.DeployService
	hub           *websocket.Hub
}

func NewDeployHandler(deployService *services.DeployService, hub *websocket.Hub) *DeployHandler {
	return &DeployHandler{
		deployService: deployService,
		hub:           hub,
	}
}

// Git Repository Handlers

// CreateGitRepository handles POST /api/v1/deploy/repositories
func (h *DeployHandler) CreateGitRepository(c *gin.Context) {
	userID := c.GetString("user_id")
	org := c.GetString("organization")

	var req struct {
		Name         string `json:"name" binding:"required"`
		Provider     string `json:"provider" binding:"required,oneof=github gitlab bitbucket gitea azure-devops"`
		URL          string `json:"url" binding:"required"`
		CloneURL     string `json:"clone_url"`
		Branch       string `json:"branch"`
		AuthType     string `json:"auth_type" binding:"required,oneof=oauth ssh token basic"`
		AuthToken    string `json:"auth_token"`
		SSHKey       string `json:"ssh_key"`
		WebhookURL   string `json:"webhook_url"`
		WebhookToken string `json:"webhook_token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid := uuid.MustParse(userID)
	repo, err := h.deployService.CreateGitRepository(
		uid, org, req.Name, req.Provider, req.URL, req.CloneURL,
		req.Branch, req.AuthType, req.AuthToken, req.SSHKey,
		req.WebhookURL, req.WebhookToken,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, repo)
}

// GetGitRepository handles GET /api/v1/deploy/repositories/:id
func (h *DeployHandler) GetGitRepository(c *gin.Context) {
	id := c.Param("id")
	repoID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid repository id"})
		return
	}

	repo, err := h.deployService.GetGitRepository(repoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "repository not found"})
		return
	}

	c.JSON(http.StatusOK, repo)
}

// ListGitRepositories handles GET /api/v1/deploy/repositories
func (h *DeployHandler) ListGitRepositories(c *gin.Context) {
	org := c.GetString("organization")

	repos, err := h.deployService.ListGitRepositories(org)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"repositories": repos})
}

// UpdateGitRepository handles PATCH /api/v1/deploy/repositories/:id
func (h *DeployHandler) UpdateGitRepository(c *gin.Context) {
	id := c.Param("id")
	repoID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid repository id"})
		return
	}

	var req struct {
		Name       string `json:"name"`
		Branch     string `json:"branch"`
		WebhookURL string `json:"webhook_url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repo, err := h.deployService.UpdateGitRepository(repoID, req.Name, req.Branch, req.WebhookURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, repo)
}

// DeleteGitRepository handles DELETE /api/v1/deploy/repositories/:id
func (h *DeployHandler) DeleteGitRepository(c *gin.Context) {
	id := c.Param("id")
	repoID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid repository id"})
		return
	}

	if err := h.deployService.DeleteGitRepository(repoID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Pipeline Handlers

// CreatePipeline handles POST /api/v1/deploy/pipelines
func (h *DeployHandler) CreatePipeline(c *gin.Context) {
	userID := c.GetString("user_id")
	org := c.GetString("organization")

	var req struct {
		Name         string                 `json:"name" binding:"required"`
		Description  string                 `json:"description"`
		RepositoryID string                 `json:"repository_id" binding:"required"`
		Config       map[string]interface{} `json:"config"`
		TriggerType  string                 `json:"trigger_type" binding:"required,oneof=webhook schedule manual api"`
		Schedule     string                 `json:"schedule"`
		Environment  string                 `json:"environment" binding:"required,oneof=development staging production canary"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repoID := uuid.MustParse(req.RepositoryID)
	uid := uuid.MustParse(userID)

	pipeline, err := h.deployService.CreatePipeline(
		uid, org, req.Name, req.Description, repoID,
		req.Config, req.TriggerType, req.Schedule, req.Environment,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, pipeline)
}

// GetPipeline handles GET /api/v1/deploy/pipelines/:id
func (h *DeployHandler) GetPipeline(c *gin.Context) {
	id := c.Param("id")
	pipelineID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pipeline id"})
		return
	}

	pipeline, err := h.deployService.GetPipeline(pipelineID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "pipeline not found"})
		return
	}

	c.JSON(http.StatusOK, pipeline)
}

// ListPipelines handles GET /api/v1/deploy/pipelines
func (h *DeployHandler) ListPipelines(c *gin.Context) {
	org := c.GetString("organization")
	environment := c.Query("environment")

	pipelines, err := h.deployService.ListPipelines(org, environment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pipelines": pipelines})
}

// UpdatePipeline handles PATCH /api/v1/deploy/pipelines/:id
func (h *DeployHandler) UpdatePipeline(c *gin.Context) {
	id := c.Param("id")
	pipelineID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pipeline id"})
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Config      string `json:"config"`
		Schedule    string `json:"schedule"`
		Enabled     *bool  `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	pipeline, err := h.deployService.UpdatePipeline(pipelineID, req.Name, req.Description, req.Config, req.Schedule, enabled)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pipeline)
}

// DeletePipeline handles DELETE /api/v1/deploy/pipelines/:id
func (h *DeployHandler) DeletePipeline(c *gin.Context) {
	id := c.Param("id")
	pipelineID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pipeline id"})
		return
	}

	if err := h.deployService.DeletePipeline(pipelineID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Build Handlers

// GetBuild handles GET /api/v1/deploy/builds/:id
func (h *DeployHandler) GetBuild(c *gin.Context) {
	id := c.Param("id")
	buildID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid build id"})
		return
	}

	build, err := h.deployService.GetBuild(buildID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "build not found"})
		return
	}

	c.JSON(http.StatusOK, build)
}

// ListBuilds handles GET /api/v1/deploy/builds
func (h *DeployHandler) ListBuilds(c *gin.Context) {
	org := c.GetString("organization")
	pipelineID := c.Query("pipeline_id")
	limitStr := c.DefaultQuery("limit", "20")

	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	var builds []models.Build
	var err error

	if pipelineID != "" {
		pid, err := uuid.Parse(pipelineID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pipeline id"})
			return
		}
		builds, err = h.deployService.GetBuildsByPipeline(pid, limit)
	} else {
		builds, err = h.deployService.GetBuildsByOrg(org, limit)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"builds": builds})
}

// CreateBuild handles POST /api/v1/deploy/builds
func (h *DeployHandler) CreateBuild(c *gin.Context) {
	var req struct {
		PipelineID    string `json:"pipeline_id" binding:"required"`
		CommitHash    string `json:"commit_hash" binding:"required"`
		CommitMessage string `json:"commit_message"`
		Author        string `json:"author"`
		Branch        string `json:"branch"`
		Tag           string `json:"tag"`
		Trigger       string `json:"trigger" binding:"required,oneof=manual webhook schedule api"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pipelineID := uuid.MustParse(req.PipelineID)
	build, err := h.deployService.CreateBuild(pipelineID, req.CommitHash, req.CommitMessage, req.Author, req.Branch, req.Tag, req.Trigger)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, build)
}

// UpdateBuildStatus handles PATCH /api/v1/deploy/builds/:id/status
func (h *DeployHandler) UpdateBuildStatus(c *gin.Context) {
	id := c.Param("id")
	buildID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid build id"})
		return
	}

	var req struct {
		Status       string `json:"status" binding:"required,oneof=pending running success failed cancelled"`
		ErrorMessage string `json:"error_message"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	build, err := h.deployService.UpdateBuildStatus(buildID, req.Status, req.ErrorMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, build)
}

// Deployment Handlers

// CreateDeployment handles POST /api/v1/deploy/deployments
func (h *DeployHandler) CreateDeployment(c *gin.Context) {
	userID := c.GetString("user_id")
	org := c.GetString("organization")

	var req struct {
		Name         string                 `json:"name" binding:"required"`
		Environment  string                 `json:"environment" binding:"required,oneof=development staging production canary"`
		Namespace    string                 `json:"namespace"`
		Cluster      string                 `json:"cluster"`
		Strategy     string                 `json:"strategy" binding:"required,oneof=rolling blue-green canary recreate"`
		BuildID      string                 `json:"build_id"`
		Image        string                 `json:"image" binding:"required"`
		Replicas     int                    `json:"replicas"`
		HealthCheck  string                 `json:"health_check"`
		AutoRollback bool                   `json:"auto_rollback"`
		Config       map[string]interface{} `json:"config"`
		Variables    map[string]interface{} `json:"variables"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid := uuid.MustParse(userID)
	var buildID *uuid.UUID
	if req.BuildID != "" {
		bid := uuid.MustParse(req.BuildID)
		buildID = &bid
	}

	deployment, err := h.deployService.CreateDeployment(
		uid, org, req.Name, req.Environment, req.Namespace, req.Cluster,
		req.Strategy, req.Image, req.Replicas, buildID, req.AutoRollback,
		req.Config, req.Variables,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, deployment)
}

// GetDeployment handles GET /api/v1/deploy/deployments/:id
func (h *DeployHandler) GetDeployment(c *gin.Context) {
	id := c.Param("id")
	deployID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid deployment id"})
		return
	}

	deployment, err := h.deployService.GetDeployment(deployID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "deployment not found"})
		return
	}

	c.JSON(http.StatusOK, deployment)
}

// ListDeployments handles GET /api/v1/deploy/deployments
func (h *DeployHandler) ListDeployments(c *gin.Context) {
	org := c.GetString("organization")
	environment := c.Query("environment")
	limitStr := c.DefaultQuery("limit", "20")

	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	deployments, err := h.deployService.ListDeployments(org, environment, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deployments": deployments})
}

// RollbackDeployment handles POST /api/v1/deploy/deployments/:id/rollback
func (h *DeployHandler) RollbackDeployment(c *gin.Context) {
	id := c.Param("id")
	deployID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid deployment id"})
		return
	}

	userID := c.GetString("user_id")
	uid := uuid.MustParse(userID)

	var req struct {
		Reason string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deployment, err := h.deployService.RollbackDeployment(deployID, uid, req.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, deployment)
}

// Webhook Handlers

// HandleWebhook handles POST /api/v1/deploy/webhooks/:provider
func (h *DeployHandler) HandleWebhook(c *gin.Context) {
	provider := c.Param("provider")
	eventType := c.GetHeader("X-GitHub-Event") // or X-GitLab-Event, etc.
	signature := c.GetHeader("X-Hub-Signature-256")

	payload, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read payload"})
		return
	}

	repoID := c.Query("repository_id")
	var repoUUID *uuid.UUID
	if repoID != "" {
		id := uuid.MustParse(repoID)
		repoUUID = &id
	}

	if err := h.deployService.ProcessWebhookEvent(provider, eventType, signature, payload, repoUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "accepted"})
}

// Statistics Handler

// GetDeployStats handles GET /api/v1/deploy/stats
func (h *DeployHandler) GetDeployStats(c *gin.Context) {
	org := c.GetString("organization")

	stats, err := h.deployService.GetDeployStats(org)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
