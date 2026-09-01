package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestExecuteTerminalCommand_Blocked(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/v1/terminal/execute", ExecuteTerminalCommand)

	body, _ := json.Marshal(map[string]string{
		"machine_id": "00000000-0000-0000-0000-000000000001",
		"session_id": "00000000-0000-0000-0000-000000000002",
		"command":    "",
	})

	req, _ := http.NewRequest("POST", "/api/v1/terminal/execute", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for empty command, got %d", w.Code)
	}
}

func TestRemoteFileManager_Operations(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/api/v1/files", ListFiles)
	r.GET("/api/v1/files/content", ReadFileContent)
	r.POST("/api/v1/files/content", WriteFileContent)
	r.POST("/api/v1/files/mkdir", MakeDirectory)
	r.DELETE("/api/v1/files", DeleteFile)

	// Create temp directory for testing
	tmpDir, err := os.MkdirTemp("", "infrapilot_remote_ops_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFilePath := filepath.Join(tmpDir, "sample.txt")

	// 1. Write file content
	writeBody, _ := json.Marshal(map[string]string{
		"path":    testFilePath,
		"content": "Hello InfraPilot Enterprise Remote Operations",
	})
	req1, _ := http.NewRequest("POST", "/api/v1/files/content", bytes.NewBuffer(writeBody))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("Expected WriteFileContent status 200, got %d: %s", w1.Code, w1.Body.String())
	}

	// 2. Read file content
	req2, _ := http.NewRequest("GET", "/api/v1/files/content?path="+testFilePath, nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("Expected ReadFileContent status 200, got %d: %s", w2.Code, w2.Body.String())
	}

	var readResp map[string]interface{}
	_ = json.Unmarshal(w2.Body.Bytes(), &readResp)
	if content, ok := readResp["content"].(string); !ok || content != "Hello InfraPilot Enterprise Remote Operations" {
		t.Errorf("Expected file content match, got %v", readResp["content"])
	}

	// 3. Make Directory
	subDir := filepath.Join(tmpDir, "subfolder")
	mkdirBody, _ := json.Marshal(map[string]string{
		"path": subDir,
	})
	req3, _ := http.NewRequest("POST", "/api/v1/files/mkdir", bytes.NewBuffer(mkdirBody))
	req3.Header.Set("Content-Type", "application/json")
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)

	if w3.Code != http.StatusOK {
		t.Errorf("Expected MakeDirectory status 200, got %d", w3.Code)
	}

	// 4. List Directory
	req4, _ := http.NewRequest("GET", "/api/v1/files?path="+tmpDir, nil)
	w4 := httptest.NewRecorder()
	r.ServeHTTP(w4, req4)

	if w4.Code != http.StatusOK {
		t.Errorf("Expected ListFiles status 200, got %d", w4.Code)
	}

	// 5. Delete File
	req5, _ := http.NewRequest("DELETE", "/api/v1/files?path="+testFilePath, nil)
	w5 := httptest.NewRecorder()
	r.ServeHTTP(w5, req5)

	if w5.Code != http.StatusOK {
		t.Errorf("Expected DeleteFile status 200, got %d", w5.Code)
	}
}
