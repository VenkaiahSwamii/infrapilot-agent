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

func TestRemediationHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/api/v1/remediation-policies", GetRemediationPolicies)
	r.POST("/api/v1/remediation-policies", CreateRemediationPolicy)
	r.PATCH("/api/v1/remediation-policies/:id", UpdateRemediationPolicy)
	r.DELETE("/api/v1/remediation-policies/:id", DeleteRemediationPolicy)

	r.GET("/api/v1/remediation-jobs", GetRemediationJobs)
	r.GET("/api/v1/remediation-jobs/analytics", GetRemediationAnalytics)
	r.GET("/api/v1/remediation-jobs/:id", GetRemediationJobByID)
	r.POST("/api/v1/remediation-jobs/:id/retry", RetryRemediationJobHandler)
	r.POST("/api/v1/remediation/test", TestRemediationHandler)

	testID := uuid.New().String()

	// 1. Test GetRemediationPolicies
	req1, _ := http.NewRequest("GET", "/api/v1/remediation-policies", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("Expected GetRemediationPolicies status 200, got %d", w1.Code)
	}

	// 2. Test GetRemediationAnalytics
	reqAnalytics, _ := http.NewRequest("GET", "/api/v1/remediation-jobs/analytics", nil)
	wAnalytics := httptest.NewRecorder()
	r.ServeHTTP(wAnalytics, reqAnalytics)

	if wAnalytics.Code != http.StatusOK {
		t.Errorf("Expected GetRemediationAnalytics status 200, got %d", wAnalytics.Code)
	}

	// 3. Test CreateRemediationPolicy
	policyBody, _ := json.Marshal(map[string]string{
		"name":        "Test Policy",
		"alert_type":  "Memory",
		"severity":    "Critical",
		"action_type": "restart_service",
	})
	req2, _ := http.NewRequest("POST", "/api/v1/remediation-policies", bytes.NewBuffer(policyBody))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusCreated && w2.Code != http.StatusOK {
		t.Errorf("Expected CreateRemediationPolicy status 201/200, got %d: %s", w2.Code, w2.Body.String())
	}

	// 4. Test GetRemediationJobs
	reqJobs, _ := http.NewRequest("GET", "/api/v1/remediation-jobs", nil)
	wJobs := httptest.NewRecorder()
	r.ServeHTTP(wJobs, reqJobs)

	if wJobs.Code != http.StatusOK {
		t.Errorf("Expected GetRemediationJobs status 200, got %d", wJobs.Code)
	}

	// 5. Test TestRemediationHandler
	testExecBody, _ := json.Marshal(map[string]string{
		"machine_id":  testID,
		"action_type": "cleanup_disk",
	})
	reqTest, _ := http.NewRequest("POST", "/api/v1/remediation/test", bytes.NewBuffer(testExecBody))
	reqTest.Header.Set("Content-Type", "application/json")
	wTest := httptest.NewRecorder()
	r.ServeHTTP(wTest, reqTest)

	if wTest.Code != http.StatusOK {
		t.Errorf("Expected TestRemediationHandler status 200, got %d: %s", wTest.Code, wTest.Body.String())
	}
}
