package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestGetMachineHealthHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/api/v1/machines/:id/health", GetMachineHealth)

	testID := uuid.New().String()

	req, _ := http.NewRequest("GET", "/api/v1/machines/"+testID+"/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected GetMachineHealth status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Errorf("Failed to parse JSON response: %v", err)
	}

	if _, ok := resp["health_score"]; !ok {
		t.Errorf("Expected 'health_score' key in JSON response")
	}

	if _, ok := resp["ai_recommendations"]; !ok {
		t.Errorf("Expected 'ai_recommendations' key in JSON response")
	}
}
