package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestDashboardHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/api/v1/dashboard/stats", Dashboard)
	r.GET("/api/v1/overview", GetOverview)

	// 1. Test GET /api/v1/dashboard/stats
	req1, _ := http.NewRequest("GET", "/api/v1/dashboard/stats", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("Expected Dashboard status 200, got %d", w1.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w1.Body.Bytes(), &resp); err != nil {
		t.Errorf("Failed to parse JSON response: %v", err)
	}

	if _, ok := resp["totalMachines"]; !ok {
		t.Errorf("Expected 'totalMachines' key in response JSON")
	}

	// 2. Test GET /api/v1/overview
	req2, _ := http.NewRequest("GET", "/api/v1/overview", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("Expected GetOverview status 200, got %d", w2.Code)
	}
}
