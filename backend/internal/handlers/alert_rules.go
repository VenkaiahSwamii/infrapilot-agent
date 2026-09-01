package handlers

import (
	"net/http"

	"infrapilot/backend/internal/models"
	"infrapilot/backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AlertRuleHandler struct {
	service *services.AlertRuleService
}

func NewAlertRuleHandler() *AlertRuleHandler {
	return &AlertRuleHandler{
		service: services.NewAlertRuleService(),
	}
}

// ListRules handles GET /organizations/:orgId/alert-rules
func (h *AlertRuleHandler) ListRules(c *gin.Context) {
	orgID := c.Param("orgId")
	rules, err := h.service.ListRules(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load alert rules"})
		return
	}
	c.JSON(http.StatusOK, rules)
}

// GetRule handles GET /organizations/:orgId/alert-rules/:id
func (h *AlertRuleHandler) GetRule(c *gin.Context) {
	id := c.Param("id")

	ruleID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule id"})
		return
	}

	rule, err := h.service.ListRules(c.Param("orgId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load alert rules"})
		return
	}

	for _, r := range rule {
		if r.ID == ruleID {
			c.JSON(http.StatusOK, r)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "alert rule not found"})
}

// CreateRule handles POST /organizations/:orgId/alert-rules
func (h *AlertRuleHandler) CreateRule(c *gin.Context) {
	orgID := c.Param("orgId")

	var req services.CreateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.OrganizationID = orgID
	rule, err := h.service.CreateRule(req)
	if err != nil {
		switch err {
		case services.ErrInvalidMetric, services.ErrInvalidOperator, services.ErrInvalidSeverity:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case services.ErrDuplicateName:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, rule)
}

// UpdateRule handles PATCH /organizations/:orgId/alert-rules/:id
func (h *AlertRuleHandler) UpdateRule(c *gin.Context) {
	_ = c.Param("orgId")
	id := c.Param("id")

	ruleID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule id"})
		return
	}

	var req services.UpdateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rule, err := h.service.UpdateRule(ruleID, req)
	if err != nil {
		switch err {
		case services.ErrRuleNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case services.ErrInvalidMetric, services.ErrInvalidOperator, services.ErrInvalidSeverity:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case services.ErrDuplicateName:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, rule)
}

// DeleteRule handles DELETE /organizations/:orgId/alert-rules/:id
func (h *AlertRuleHandler) DeleteRule(c *gin.Context) {
	id := c.Param("id")

	ruleID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule id"})
		return
	}

	if err := h.service.DeleteRule(ruleID); err != nil {
		if err == services.ErrRuleNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ToggleRule handles PATCH /organizations/:orgId/alert-rules/:id/toggle
func (h *AlertRuleHandler) ToggleRule(c *gin.Context) {
	_ = c.Param("orgId")
	id := c.Param("id")

	ruleID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule id"})
		return
	}

	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var rule *models.AlertRule
	if req.Enabled {
		rule, err = h.service.EnableRule(ruleID)
	} else {
		rule, err = h.service.DisableRule(ruleID)
	}

	if err != nil {
		if err == services.ErrRuleNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rule)
}
