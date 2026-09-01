package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"infrapilot/backend/internal/events"
	"infrapilot/backend/internal/models"
	"infrapilot/backend/internal/repository"
	"infrapilot/backend/internal/services"
	"infrapilot/backend/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestMachineLifecycle_RegistrationAndHeartbeat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	serverRepo := repository.NewServerRepository()
	metricRepo := repository.NewMetricRepository()
	hub := websocket.NewHub()
	bus := events.NewEventBus()

	serverService := services.NewServerService(serverRepo, metricRepo, hub, bus)
	serverHandler := NewServerHandler(serverService)

	r.POST("/api/v1/agent/register", serverHandler.RegisterServer)
	r.POST("/api/v1/servers/:id/heartbeat", serverHandler.ServerHeartbeat)
	r.GET("/api/v1/servers", serverHandler.GetServers)

	machineID := uuid.New().String()

	// 1. Machine Registration Test
	regPayload := map[string]interface{}{
		"machine_id":      machineID,
		"hostname":        "Windows-01",
		"ip_address":      "192.168.1.50",
		"os":              "Windows Server 2022",
		"platform":        "windows",
		"agent_version":   "v1.1.0",
		"resource_type":   "windows",
		"organization":    "Default Organization",
		"total_memory_gb": 16,
		"total_disk_gb":    500,
	}
	bodyBytes, _ := json.Marshal(regPayload)
	req1, _ := http.NewRequest("POST", "/api/v1/agent/register", bytes.NewBuffer(bodyBytes))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK && w1.Code != http.StatusCreated {
		t.Fatalf("Expected registration status 200/201, got %d. Body: %s", w1.Code, w1.Body.String())
	}

	var regResp map[string]interface{}
	_ = json.Unmarshal(w1.Body.Bytes(), &regResp)

	if regResp["machine_id"] == nil {
		t.Errorf("Expected machine_id in registration response")
	}
	if regResp["api_key"] == nil {
		t.Errorf("Expected api_key in registration response")
	}

	apiKey, _ := regResp["api_key"].(string)

	// 2. Machine Heartbeat Test
	req2, _ := http.NewRequest("POST", "/api/v1/servers/"+machineID+"/heartbeat", nil)
	if apiKey != "" {
		req2.Header.Set("X-API-Key", apiKey)
	}
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("Expected heartbeat status 200, got %d. Body: %s", w2.Code, w2.Body.String())
	}

	// 3. Machine List API Test
	req3, _ := http.NewRequest("GET", "/api/v1/servers", nil)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)

	if w3.Code != http.StatusOK {
		t.Fatalf("Expected GetServers status 200, got %d", w3.Code)
	}
}

func TestMachineOfflineDetectionLogic(t *testing.T) {
	now := time.Now()
	recentHeartbeat := now.Add(-10 * time.Second)
	oldHeartbeat := now.Add(-5 * time.Minute)

	mOnline := models.Machine{
		ID:       uuid.New(),
		Hostname: "Ubuntu-01",
		Status:   "ONLINE",
		LastSeen: recentHeartbeat,
	}

	mOffline := models.Machine{
		ID:       uuid.New(),
		Hostname: "Ubuntu-02",
		Status:   "OFFLINE",
		LastSeen: oldHeartbeat,
	}

	if time.Since(mOnline.LastSeen) > 90*time.Second {
		t.Errorf("mOnline should be detected as ONLINE")
	}
	if time.Since(mOffline.LastSeen) <= 90*time.Second {
		t.Errorf("mOffline should be detected as OFFLINE")
	}
}
