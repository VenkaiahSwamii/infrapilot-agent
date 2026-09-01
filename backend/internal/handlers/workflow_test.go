package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestWorkflowHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/api/v1/workflows", GetWorkflows)
	r.POST("/api/v1/workflows", CreateWorkflow)
	r.GET("/api/v1/workflows/:id", GetWorkflowByID)
	r.PATCH("/api/v1/workflows/:id", UpdateWorkflow)
	r.DELETE("/api/v1/workflows/:id", DeleteWorkflow)
	r.POST("/api/v1/workflows/:id/run", RunWorkflow)
	r.POST("/api/v1/workflows/generate-ai", GenerateAIWorkflowHandler)
	r.GET("/api/v1/workflows/executions", GetWorkflowExecutions)
	r.GET("/api/v1/workflows/executions/analytics", GetWorkflowAnalytics)
	r.GET("/api/v1/workflows/executions/:id", GetWorkflowExecutionByID)
	r.POST("/api/v1/workflows/executions/:id/cancel", CancelWorkflowHandler)

	testID := uuid.New().String()

	// 1. Test GetWorkflows
	req1, _ := http.NewRequest("GET", "/api/v1/workflows", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("Expected GetWorkflows status 200, got %d", w1.Code)
	}

	// 2. Test GetWorkflowAnalytics
	reqAnalytics, _ := http.NewRequest("GET", "/api/v1/workflows/executions/analytics", nil)
	wAnalytics := httptest.NewRecorder()
	r.ServeHTTP(wAnalytics, reqAnalytics)

	if wAnalytics.Code != http.StatusOK {
		t.Errorf("Expected GetWorkflowAnalytics status 200, got %d", wAnalytics.Code)
	}

	// 3. Test CreateWorkflow
	createBody, _ := json.Marshal(map[string]interface{}{
		"name":         "Test Multi-step Runbook",
		"trigger_type": "Alert",
		"steps": []map[string]interface{}{
			{"step_order": 1, "name": "Restart Service", "action_type": "restart_service", "action_value": "nginx"},
			{"step_order": 2, "name": "Wait 30s", "action_type": "wait", "action_value": "30"},
		},
	})
	reqCreate, _ := http.NewRequest("POST", "/api/v1/workflows", bytes.NewBuffer(createBody))
	reqCreate.Header.Set("Content-Type", "application/json")
	wCreate := httptest.NewRecorder()
	r.ServeHTTP(wCreate, reqCreate)

	if wCreate.Code != http.StatusCreated && wCreate.Code != http.StatusOK {
		t.Errorf("Expected CreateWorkflow status 201/200, got %d: %s", wCreate.Code, wCreate.Body.String())
	}

	// 4. Test GenerateAIWorkflowHandler
	aiBody, _ := json.Marshal(map[string]string{
		"prompt": "When CPU is over 90%, restart nginx and send Slack alert",
	})
	reqAI, _ := http.NewRequest("POST", "/api/v1/workflows/generate-ai", bytes.NewBuffer(aiBody))
	reqAI.Header.Set("Content-Type", "application/json")
	wAI := httptest.NewRecorder()
	r.ServeHTTP(wAI, reqAI)

	if wAI.Code != http.StatusCreated && wAI.Code != http.StatusOK {
		t.Errorf("Expected GenerateAIWorkflowHandler status 201/200, got %d: %s", wAI.Code, wAI.Body.String())
	}

	// 5. Test GetWorkflowExecutions
	reqExecs, _ := http.NewRequest("GET", "/api/v1/workflows/executions", nil)
	wExecs := httptest.NewRecorder()
	r.ServeHTTP(wExecs, reqExecs)

	if wExecs.Code != http.StatusOK {
		t.Errorf("Expected GetWorkflowExecutions status 200, got %d", wExecs.Code)
	}

	// 6. Test RunWorkflow
	runBody, _ := json.Marshal(map[string]string{
		"machine_id": testID,
	})
	reqRun, _ := http.NewRequest("POST", "/api/v1/workflows/"+testID+"/run", bytes.NewBuffer(runBody))
	reqRun.Header.Set("Content-Type", "application/json")
	wRun := httptest.NewRecorder()
	r.ServeHTTP(wRun, reqRun)

	if wRun.Code != http.StatusOK {
		t.Errorf("Expected RunWorkflow status 200, got %d: %s", wRun.Code, wRun.Body.String())
	}
}
