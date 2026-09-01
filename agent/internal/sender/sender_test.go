package sender

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"infrapilot/agent/internal/cache"
	"infrapilot/agent/internal/state"
)

func TestSendMetrics_SuccessAndReplay(t *testing.T) {
	cacheFile := "test_offline_queue_1.json"
	cache.Init(cacheFile)
	defer os.Remove(cacheFile)

	// Mock backend server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	}))
	defer server.Close()

	machineID := "0efc7e0f-27f9-4f1f-846a-7a3d2df9cfd6"
	apiKey := "test-api-key"

	payload := map[string]interface{}{
		"cpu":    25.5,
		"memory": 60.2,
		"disk":   45.0,
	}

	err := SendMetrics(server.URL, machineID, apiKey, payload)
	if err != nil {
		t.Fatalf("Expected SendMetrics success, got error: %v", err)
	}

	if !state.Connected {
		t.Errorf("Expected state.Connected to be true after successful metrics send")
	}
}

func TestOfflineQueue_AndReplay(t *testing.T) {
	cacheFile := "test_offline_queue_2.json"
	cache.Init(cacheFile)
	defer os.Remove(cacheFile)

	// Queue a metric while offline
	offlineMetric := map[string]interface{}{
		"cpu":     99.9,
		"memory":  88.8,
		"offline": true,
	}

	cache.QueueMetric(offlineMetric)
	if cache.QueueSize() == 0 {
		t.Fatalf("Expected non-zero queue size after QueueMetric")
	}

	// Mock backend server accepting replayed metrics
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	ReplayCachedMetrics(server.URL, "test-machine-id")

	if cache.QueueSize() != 0 {
		t.Errorf("Expected queue to be drained after replay, got size: %d", cache.QueueSize())
	}
}
