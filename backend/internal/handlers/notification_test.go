package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNotificationHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/api/v1/notifications", GetNotifications)
	r.POST("/api/v1/notifications/test", SendTestNotification)
	r.GET("/api/v1/notification-policies", GetNotificationPolicies)
	r.POST("/api/v1/notification-policies", CreateNotificationPolicy)
	r.PATCH("/api/v1/notification-policies/:id", UpdateNotificationPolicy)
	r.DELETE("/api/v1/notification-policies/:id", DeleteNotificationPolicy)

	// 1. Test GetNotifications
	req1, _ := http.NewRequest("GET", "/api/v1/notifications", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("Expected GetNotifications status 200, got %d", w1.Code)
	}

	// 2. Test SendTestNotification
	testBody, _ := json.Marshal(map[string]string{
		"channel":   "Slack",
		"recipient": "#infrastructure-alerts",
	})
	req2, _ := http.NewRequest("POST", "/api/v1/notifications/test", bytes.NewBuffer(testBody))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("Expected SendTestNotification status 200, got %d: %s", w2.Code, w2.Body.String())
	}

	// 3. Test GetNotificationPolicies
	req3, _ := http.NewRequest("GET", "/api/v1/notification-policies", nil)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)

	if w3.Code != http.StatusOK {
		t.Errorf("Expected GetNotificationPolicies status 200, got %d", w3.Code)
	}

	// 4. Test CreateNotificationPolicy
	policyBody, _ := json.Marshal(map[string]string{
		"severity":  "Critical",
		"channel":   "Telegram",
		"recipient": "@InfraPilotAlertBot",
	})
	req4, _ := http.NewRequest("POST", "/api/v1/notification-policies", bytes.NewBuffer(policyBody))
	req4.Header.Set("Content-Type", "application/json")
	w4 := httptest.NewRecorder()
	r.ServeHTTP(w4, req4)

	if w4.Code != http.StatusCreated && w4.Code != http.StatusOK {
		t.Errorf("Expected CreateNotificationPolicy status 201/200, got %d: %s", w4.Code, w4.Body.String())
	}
}
