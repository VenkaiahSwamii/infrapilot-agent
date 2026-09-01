package ai

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestAIService_Handlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/api/v1/ai/analyze", AnalyzeHandler)
	r.POST("/api/v1/ai/chat", ChatHandler)
	r.GET("/api/v1/ai/incidents", GetIncidentsHandler)
	r.GET("/api/v1/ai/recommendations", GetRecommendationsHandler)
	r.GET("/api/v1/ai/health-score", GetHealthScoreHandler)
	r.GET("/api/v1/ai/predictions", GetPredictionsHandler)

	// 1. Test AnalyzeHandler
	analyzeReq := AnalyzeRequest{
		AlertID:  uuid.New().String(),
		ServerID: uuid.New().String(),
	}
	reqBytes, _ := json.Marshal(analyzeReq)
	req1, _ := http.NewRequest("POST", "/api/v1/ai/analyze", bytes.NewBuffer(reqBytes))
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	if w1.Code != http.StatusOK {
		t.Errorf("Expected AnalyzeHandler status 200, got %d", w1.Code)
	}

	// 2. Test ChatHandler
	chatReq := ChatRequest{
		SessionID: "sess-abc",
		Query:     "Why is Production-01 slow?",
	}
	reqBytes2, _ := json.Marshal(chatReq)
	req2, _ := http.NewRequest("POST", "/api/v1/ai/chat", bytes.NewBuffer(reqBytes2))
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Errorf("Expected ChatHandler status 200, got %d", w2.Code)
	}

	// 3. Test GetIncidentsHandler
	req3, _ := http.NewRequest("GET", "/api/v1/ai/incidents", nil)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	if w3.Code != http.StatusOK {
		t.Errorf("Expected GetIncidentsHandler status 200, got %d", w3.Code)
	}

	// 4. Test GetRecommendationsHandler
	req4, _ := http.NewRequest("GET", "/api/v1/ai/recommendations", nil)
	w4 := httptest.NewRecorder()
	r.ServeHTTP(w4, req4)
	if w4.Code != http.StatusOK {
		t.Errorf("Expected GetRecommendationsHandler status 200, got %d", w4.Code)
	}

	// 5. Test GetHealthScoreHandler
	req5, _ := http.NewRequest("GET", "/api/v1/ai/health-score", nil)
	w5 := httptest.NewRecorder()
	r.ServeHTTP(w5, req5)
	if w5.Code != http.StatusOK {
		t.Errorf("Expected GetHealthScoreHandler status 200, got %d", w5.Code)
	}

	// 6. Test GetPredictionsHandler
	req6, _ := http.NewRequest("GET", "/api/v1/ai/predictions", nil)
	w6 := httptest.NewRecorder()
	r.ServeHTTP(w6, req6)
	if w6.Code != http.StatusOK {
		t.Errorf("Expected GetPredictionsHandler status 200, got %d", w6.Code)
	}
}
