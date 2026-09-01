package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetAPMStats(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/api/v1/observability/apm", GetAPMStats)

	req, _ := http.NewRequest("GET", "/api/v1/observability/apm", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if _, ok := resp["avg_response_time_ms"]; !ok {
		t.Errorf("Expected avg_response_time_ms in response")
	}
}

func TestGetBusinessKPIs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/api/v1/observability/kpis", GetBusinessKPIs)

	req, _ := http.NewRequest("GET", "/api/v1/observability/kpis", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if _, ok := resp["infrastructure_uptime_pct"].(float64); !ok {
		t.Errorf("Expected infrastructure_uptime_pct float64 in response")
	}
}
