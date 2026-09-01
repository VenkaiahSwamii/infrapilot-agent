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

func TestDockerDeepManagement_Handlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/api/v1/docker/containers/:id", GetMachineContainers)
	r.POST("/api/v1/docker/container/start", ContainerStart)
	r.POST("/api/v1/docker/container/stop", ContainerStop)
	r.POST("/api/v1/docker/container/restart", ContainerRestart)
	r.POST("/api/v1/docker/container/remove", ContainerRemove)
	r.GET("/api/v1/docker/containers/logs", GetContainerLogs)
	r.GET("/api/v1/docker/images", GetDockerImages)
	r.GET("/api/v1/docker/networks", GetDockerNetworks)
	r.GET("/api/v1/docker/volumes", GetDockerVolumes)

	testMachineID := uuid.New().String()

	// 1. Test GetMachineContainers
	req1, _ := http.NewRequest("GET", "/api/v1/docker/containers/"+testMachineID, nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Fatalf("Expected GetMachineContainers status 200, got %d", w1.Code)
	}

	var containers []map[string]interface{}
	_ = json.Unmarshal(w1.Body.Bytes(), &containers)
	if len(containers) == 0 {
		t.Errorf("Expected non-empty container list")
	}

	// 2. Test Container Actions
	actions := []struct {
		endpoint string
		action   string
	}{
		{"/api/v1/docker/container/start", "start"},
		{"/api/v1/docker/container/stop", "stop"},
		{"/api/v1/docker/container/restart", "restart"},
		{"/api/v1/docker/container/remove", "remove"},
	}

	for _, act := range actions {
		body, _ := json.Marshal(map[string]string{
			"machine_id": testMachineID,
			"container":  "nginx-prod",
		})
		req, _ := http.NewRequest("POST", act.endpoint, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected %s status 200, got %d: %s", act.action, w.Code, w.Body.String())
		}
	}

	// 3. Test GetContainerLogs
	reqLogs, _ := http.NewRequest("GET", "/api/v1/docker/containers/logs?container=nginx-prod&tail=50", nil)
	wLogs := httptest.NewRecorder()
	r.ServeHTTP(wLogs, reqLogs)

	if wLogs.Code != http.StatusOK {
		t.Errorf("Expected GetContainerLogs status 200, got %d", wLogs.Code)
	}

	// 4. Test GetDockerImages
	reqImg, _ := http.NewRequest("GET", "/api/v1/docker/images", nil)
	wImg := httptest.NewRecorder()
	r.ServeHTTP(wImg, reqImg)

	if wImg.Code != http.StatusOK {
		t.Errorf("Expected GetDockerImages status 200, got %d", wImg.Code)
	}

	// 5. Test GetDockerNetworks
	reqNet, _ := http.NewRequest("GET", "/api/v1/docker/networks", nil)
	wNet := httptest.NewRecorder()
	r.ServeHTTP(wNet, reqNet)

	if wNet.Code != http.StatusOK {
		t.Errorf("Expected GetDockerNetworks status 200, got %d", wNet.Code)
	}

	// 6. Test GetDockerVolumes
	reqVol, _ := http.NewRequest("GET", "/api/v1/docker/volumes", nil)
	wVol := httptest.NewRecorder()
	r.ServeHTTP(wVol, reqVol)

	if wVol.Code != http.StatusOK {
		t.Errorf("Expected GetDockerVolumes status 200, got %d", wVol.Code)
	}
}
