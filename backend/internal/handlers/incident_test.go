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

func TestIncidentHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/api/v1/incidents", GetIncidents)
	r.GET("/api/v1/incidents/analytics", GetIncidentAnalytics)
	r.GET("/api/v1/incidents/:id", GetIncidentByID)
	r.GET("/api/v1/incidents/:id/timeline", GetIncidentTimeline)
	r.GET("/api/v1/incidents/:id/alerts", GetIncidentAlerts)
	r.GET("/api/v1/incidents/:id/summary", GetAIIncidentSummary)
	r.POST("/api/v1/incidents/:id/resolve", ResolveIncident)
	r.POST("/api/v1/incidents/:id/reopen", ReopenIncident)

	testID := uuid.New().String()

	// 1. Test GetIncidents
	req1, _ := http.NewRequest("GET", "/api/v1/incidents", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("Expected GetIncidents status 200, got %d", w1.Code)
	}

	// 2. Test GetIncidentAnalytics
	reqAnalytics, _ := http.NewRequest("GET", "/api/v1/incidents/analytics", nil)
	wAnalytics := httptest.NewRecorder()
	r.ServeHTTP(wAnalytics, reqAnalytics)

	if wAnalytics.Code != http.StatusOK {
		t.Errorf("Expected GetIncidentAnalytics status 200, got %d", wAnalytics.Code)
	}

	// 3. Test GetIncidentByID
	req2, _ := http.NewRequest("GET", "/api/v1/incidents/"+testID, nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("Expected GetIncidentByID status 200, got %d", w2.Code)
	}

	// 4. Test GetIncidentTimeline
	reqTimeline, _ := http.NewRequest("GET", "/api/v1/incidents/"+testID+"/timeline", nil)
	wTimeline := httptest.NewRecorder()
	r.ServeHTTP(wTimeline, reqTimeline)

	if wTimeline.Code != http.StatusOK {
		t.Errorf("Expected GetIncidentTimeline status 200, got %d", wTimeline.Code)
	}

	// 5. Test GetIncidentAlerts
	reqAlerts, _ := http.NewRequest("GET", "/api/v1/incidents/"+testID+"/alerts", nil)
	wAlerts := httptest.NewRecorder()
	r.ServeHTTP(wAlerts, reqAlerts)

	if wAlerts.Code != http.StatusOK {
		t.Errorf("Expected GetIncidentAlerts status 200, got %d", wAlerts.Code)
	}

	// 6. Test ResolveIncident
	resolveBody, _ := json.Marshal(map[string]string{
		"resolution_note": "Service restarted and CPU load normalized",
	})
	reqResolve, _ := http.NewRequest("POST", "/api/v1/incidents/"+testID+"/resolve", bytes.NewBuffer(resolveBody))
	reqResolve.Header.Set("Content-Type", "application/json")
	wResolve := httptest.NewRecorder()
	r.ServeHTTP(wResolve, reqResolve)

	if wResolve.Code != http.StatusOK {
		t.Errorf("Expected ResolveIncident status 200, got %d", wResolve.Code)
	}

	// 7. Test GetAIIncidentSummary
	reqSummary, _ := http.NewRequest("GET", "/api/v1/incidents/"+testID+"/summary", nil)
	wSummary := httptest.NewRecorder()
	r.ServeHTTP(wSummary, reqSummary)

	if wSummary.Code != http.StatusOK {
		t.Errorf("Expected GetAIIncidentSummary status 200, got %d", wSummary.Code)
	}
}
