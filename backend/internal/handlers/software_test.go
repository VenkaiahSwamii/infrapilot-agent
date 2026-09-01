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

func TestSoftwareManagement_Handlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/api/v1/software/:id", GetMachineSoftware)
	r.POST("/api/v1/software/update", UpdateSoftwarePackage)
	r.POST("/api/v1/software/install", InstallSoftwarePackage)
	r.POST("/api/v1/software/remove", RemoveSoftwarePackage)
	r.GET("/api/v1/software/patches", GetPendingPatches)
	r.POST("/api/v1/software/patch-all", PatchAllPackages)

	testMachineID := uuid.New().String()

	// 1. Test GetMachineSoftware
	req1, _ := http.NewRequest("GET", "/api/v1/software/"+testMachineID, nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Fatalf("Expected GetMachineSoftware status 200, got %d", w1.Code)
	}

	var packages []map[string]interface{}
	_ = json.Unmarshal(w1.Body.Bytes(), &packages)
	if len(packages) == 0 {
		t.Errorf("Expected non-empty software inventory list")
	}

	// 2. Test InstallSoftwarePackage
	installBody, _ := json.Marshal(map[string]string{
		"machine_id": testMachineID,
		"package":    "htop",
	})
	req2, _ := http.NewRequest("POST", "/api/v1/software/install", bytes.NewBuffer(installBody))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("Expected InstallSoftwarePackage status 200, got %d: %s", w2.Code, w2.Body.String())
	}

	// 3. Test UpdateSoftwarePackage
	updateBody, _ := json.Marshal(map[string]string{
		"machine_id": testMachineID,
		"package":    "openssl",
	})
	req3, _ := http.NewRequest("POST", "/api/v1/software/update", bytes.NewBuffer(updateBody))
	req3.Header.Set("Content-Type", "application/json")
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)

	if w3.Code != http.StatusOK {
		t.Errorf("Expected UpdateSoftwarePackage status 200, got %d: %s", w3.Code, w3.Body.String())
	}

	// 4. Test RemoveSoftwarePackage
	removeBody, _ := json.Marshal(map[string]string{
		"machine_id": testMachineID,
		"package":    "telnet",
	})
	req4, _ := http.NewRequest("POST", "/api/v1/software/remove", bytes.NewBuffer(removeBody))
	req4.Header.Set("Content-Type", "application/json")
	w4 := httptest.NewRecorder()
	r.ServeHTTP(w4, req4)

	if w4.Code != http.StatusOK {
		t.Errorf("Expected RemoveSoftwarePackage status 200, got %d: %s", w4.Code, w4.Body.String())
	}

	// 5. Test GetPendingPatches
	req5, _ := http.NewRequest("GET", "/api/v1/software/patches", nil)
	w5 := httptest.NewRecorder()
	r.ServeHTTP(w5, req5)

	if w5.Code != http.StatusOK {
		t.Errorf("Expected GetPendingPatches status 200, got %d: %s", w5.Code, w5.Body.String())
	}

	// 6. Test PatchAllPackages
	req6, _ := http.NewRequest("POST", "/api/v1/software/patch-all", nil)
	w6 := httptest.NewRecorder()
	r.ServeHTTP(w6, req6)

	if w6.Code != http.StatusOK {
		t.Errorf("Expected PatchAllPackages status 200, got %d: %s", w6.Code, w6.Body.String())
	}
}
