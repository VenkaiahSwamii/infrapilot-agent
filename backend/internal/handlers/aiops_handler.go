package handlers

import (
	"net/http"

	"infrapilot/backend/internal/aiops"

	"github.com/gin-gonic/gin"
)

// AIOpsHandler exposes REST endpoints for predictive analytics and autonomous AIOps
type AIOpsHandler struct {
	anomalyEng     *aiops.AnomalyEngine
	predictionEng  *aiops.PredictionEngine
	forecastEng    *aiops.ForecastEngine
	rootCauseEng   *aiops.RootCauseEngine
	recommendEng   *aiops.RecommendationEngine
	remediationEng *aiops.RemediationEngine
	healthScoreEng *aiops.HealthScoreEngine
	chatAssistant  *aiops.AIChatAssistant
	costSlaSecEng  *aiops.CostSLASecurityEngine
}

// NewAIOpsHandler initializes a new AIOpsHandler
func NewAIOpsHandler() *AIOpsHandler {
	anom := aiops.NewAnomalyEngine()
	pred := aiops.NewPredictionEngine()
	forc := aiops.NewForecastEngine()
	rca := aiops.NewRootCauseEngine()
	recom := aiops.NewRecommendationEngine()
	rem := aiops.NewRemediationEngine()
	hs := aiops.NewHealthScoreEngine()
	chat := aiops.NewAIChatAssistant()
	css := aiops.NewCostSLASecurityEngine()

	// Start background AIOps scheduler
	scheduler := aiops.NewAIOpsScheduler(anom, pred, forc, recom, hs)
	scheduler.Start()

	return &AIOpsHandler{
		anomalyEng:     anom,
		predictionEng:  pred,
		forecastEng:    forc,
		rootCauseEng:   rca,
		recommendEng:   recom,
		remediationEng: rem,
		healthScoreEng: hs,
		chatAssistant:  chat,
		costSlaSecEng:  css,
	}
}

// GetDashboard returns executive AIOps summary metrics
func (h *AIOpsHandler) GetDashboard(c *gin.Context) {
	orgID := c.GetString("organizationId")
	health := h.healthScoreEng.CalculateAIHealthScore(orgID)
	preds, _ := h.predictionEng.PredictFailures(orgID)
	anoms, _ := h.anomalyEng.DetectAnomaliesForMachine("server-01", orgID)
	recoms, _ := h.recommendEng.GetRecommendations(orgID)

	c.JSON(http.StatusOK, gin.H{
		"health_score":    health,
		"predictions":     preds,
		"anomalies":       anoms,
		"recommendations": recoms,
	})
}

// GetAnomalies lists detected telemetry anomalies
func (h *AIOpsHandler) GetAnomalies(c *gin.Context) {
	orgID := c.GetString("organizationId")
	machineID := c.DefaultQuery("machine_id", "server-01")

	anomalies, err := h.anomalyEng.DetectAnomaliesForMachine(machineID, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"anomalies": anomalies})
}

// GetPredictions lists predicted failure events
func (h *AIOpsHandler) GetPredictions(c *gin.Context) {
	orgID := c.GetString("organizationId")
	predictions, err := h.predictionEng.PredictFailures(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"predictions": predictions})
}

// GetForecasts retrieves 7d/30d/90d capacity forecasts
func (h *AIOpsHandler) GetForecasts(c *gin.Context) {
	orgID := c.GetString("organizationId")
	forecasts, err := h.forecastEng.GetCapacityForecasts(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"forecasts": forecasts})
}

// GetRootCause retrieves incident root cause causality analysis
func (h *AIOpsHandler) GetRootCause(c *gin.Context) {
	incidentID := c.Param("incidentId")
	orgID := c.GetString("organizationId")

	rca, err := h.rootCauseEng.AnalyzeRootCause(incidentID, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rca)
}

// GetRecommendations returns actionable AI optimization recommendations
func (h *AIOpsHandler) GetRecommendations(c *gin.Context) {
	orgID := c.GetString("organizationId")
	recoms, err := h.recommendEng.GetRecommendations(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"recommendations": recoms})
}

// ExecuteRemediation triggers auto-remediation execution
func (h *AIOpsHandler) ExecuteRemediation(c *gin.Context) {
	var req struct {
		RecommendationID string `json:"recommendation_id" binding:"required"`
		Mode             string `json:"mode"` // manual, semi_automatic, fully_automatic
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Mode == "" {
		req.Mode = "semi_automatic"
	}

	username := c.GetString("username")
	if username == "" {
		username = "admin"
	}

	if err := h.remediationEng.ExecuteRemediation(req.RecommendationID, req.Mode, username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Remediation executed successfully"})
}

// GetHealthScore retrieves overall composite AI health score
func (h *AIOpsHandler) GetHealthScore(c *gin.Context) {
	orgID := c.GetString("organizationId")
	health := h.healthScoreEng.CalculateAIHealthScore(orgID)
	c.JSON(http.StatusOK, health)
}

// QueryChat answers natural language questions using live telemetry & RAG
func (h *AIOpsHandler) QueryChat(c *gin.Context) {
	var req struct {
		Query string `json:"query" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orgID := c.GetString("organizationId")
	resp, err := h.chatAssistant.Query(req.Query, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetCostOptimization returns cost saving recommendations
func (h *AIOpsHandler) GetCostOptimization(c *gin.Context) {
	orgID := c.GetString("organizationId")
	costs, err := h.costSlaSecEng.GetCostOptimizations(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"cost_optimizations": costs})
}

// GetSLAMetrics retrieves SLA availability & MTTR metrics
func (h *AIOpsHandler) GetSLAMetrics(c *gin.Context) {
	orgID := c.GetString("organizationId")
	sla, err := h.costSlaSecEng.GetSLAMetrics(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sla)
}

// GetSecurityAnalytics returns security threat events
func (h *AIOpsHandler) GetSecurityAnalytics(c *gin.Context) {
	orgID := c.GetString("organizationId")
	sec, err := h.costSlaSecEng.GetSecurityAnalytics(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"security_events": sec})
}
