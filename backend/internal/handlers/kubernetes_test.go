package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestKubernetesClusterManagement_Handlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Register legacy endpoints
	r.GET("/api/v1/kubernetes/overview/:id", GetKubernetesOverview)
	r.GET("/api/v1/machines/:id/kubernetes", GetLinuxKubernetes)
	r.GET("/api/v1/kubernetes/nodes", GetKubernetesNodes)
	r.GET("/api/v1/kubernetes/pods", GetKubernetesPods)
	r.GET("/api/v1/kubernetes/deployments", GetKubernetesDeployments)
	r.GET("/api/v1/kubernetes/services", GetKubernetesServices)
	r.GET("/api/v1/kubernetes/pods/logs", GetKubernetesPodLogs)
	r.GET("/api/v1/kubernetes/yaml", GetKubernetesResourceYAML)

	// Register Sprint 10.8 endpoints
	r.GET("/api/v1/kubernetes/clusters", GetKubernetesClusters)
	r.GET("/api/v1/kubernetes/nodes/:clusterId", GetKubernetesNodesForCluster)
	r.GET("/api/v1/kubernetes/pods/:clusterId", GetKubernetesPodsForCluster)
	r.GET("/api/v1/kubernetes/deployments/:clusterId", GetKubernetesDeploymentsForCluster)
	r.GET("/api/v1/kubernetes/statefulsets/:clusterId", GetKubernetesStatefulSetsForCluster)
	r.GET("/api/v1/kubernetes/daemonsets/:clusterId", GetKubernetesDaemonSetsForCluster)
	r.GET("/api/v1/kubernetes/services/:clusterId", GetKubernetesServicesForCluster)
	r.GET("/api/v1/kubernetes/namespaces/:clusterId", GetKubernetesNamespacesForCluster)
	r.GET("/api/v1/kubernetes/storage/:clusterId", GetKubernetesStorageForCluster)
	r.GET("/api/v1/kubernetes/events/:clusterId", GetKubernetesEventsForCluster)
	r.GET("/api/v1/kubernetes/logs/:podId", GetKubernetesPodLogsByPodID)

	testMachineID := uuid.New().String()
	testClusterID := uuid.New().String()

	// 1. Test GetKubernetesOverview
	req1, _ := http.NewRequest("GET", "/api/v1/kubernetes/overview/"+testMachineID, nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	if w1.Code != http.StatusOK {
		t.Errorf("Expected GetKubernetesOverview status 200, got %d", w1.Code)
	}

	// 2. Test GetLinuxKubernetes
	req2, _ := http.NewRequest("GET", "/api/v1/machines/"+testMachineID+"/kubernetes", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Errorf("Expected GetLinuxKubernetes status 200, got %d", w2.Code)
	}

	// 3. Test GetKubernetesNodes
	reqNodes, _ := http.NewRequest("GET", "/api/v1/kubernetes/nodes", nil)
	wNodes := httptest.NewRecorder()
	r.ServeHTTP(wNodes, reqNodes)
	if wNodes.Code != http.StatusOK {
		t.Errorf("Expected GetKubernetesNodes status 200, got %d", wNodes.Code)
	}

	// 4. Test GetKubernetesPods
	reqPods, _ := http.NewRequest("GET", "/api/v1/kubernetes/pods", nil)
	wPods := httptest.NewRecorder()
	r.ServeHTTP(wPods, reqPods)
	if wPods.Code != http.StatusOK {
		t.Errorf("Expected GetKubernetesPods status 200, got %d", wPods.Code)
	}

	// 5. Test GetKubernetesDeployments
	reqDep, _ := http.NewRequest("GET", "/api/v1/kubernetes/deployments", nil)
	wDep := httptest.NewRecorder()
	r.ServeHTTP(wDep, reqDep)
	if wDep.Code != http.StatusOK {
		t.Errorf("Expected GetKubernetesDeployments status 200, got %d", wDep.Code)
	}

	// 6. Test GetKubernetesServices
	reqSvc, _ := http.NewRequest("GET", "/api/v1/kubernetes/services", nil)
	wSvc := httptest.NewRecorder()
	r.ServeHTTP(wSvc, reqSvc)
	if wSvc.Code != http.StatusOK {
		t.Errorf("Expected GetKubernetesServices status 200, got %d", wSvc.Code)
	}

	// 7. Test GetKubernetesPodLogs
	reqLogs, _ := http.NewRequest("GET", "/api/v1/kubernetes/pods/logs?pod=api-gateway-6f987c88b9-x2p8q&namespace=prod&tail=50", nil)
	wLogs := httptest.NewRecorder()
	r.ServeHTTP(wLogs, reqLogs)
	if wLogs.Code != http.StatusOK {
		t.Errorf("Expected GetKubernetesPodLogs status 200, got %d", wLogs.Code)
	}

	// 8. Test GetKubernetesResourceYAML
	reqYaml, _ := http.NewRequest("GET", "/api/v1/kubernetes/yaml?kind=deployment&name=api-gateway&namespace=prod", nil)
	wYaml := httptest.NewRecorder()
	r.ServeHTTP(wYaml, reqYaml)
	if wYaml.Code != http.StatusOK {
		t.Errorf("Expected GetKubernetesResourceYAML status 200, got %d", wYaml.Code)
	}

	// 9. Test Sprint 10.8: GetKubernetesClusters
	reqClusters, _ := http.NewRequest("GET", "/api/v1/kubernetes/clusters", nil)
	wClusters := httptest.NewRecorder()
	r.ServeHTTP(wClusters, reqClusters)
	if wClusters.Code != http.StatusOK {
		t.Errorf("Expected GetKubernetesClusters status 200, got %d", wClusters.Code)
	}

	// 10. Test Sprint 10.8: GetKubernetesNodesForCluster
	reqClusterNodes, _ := http.NewRequest("GET", "/api/v1/kubernetes/nodes/"+testClusterID, nil)
	wClusterNodes := httptest.NewRecorder()
	r.ServeHTTP(wClusterNodes, reqClusterNodes)
	if wClusterNodes.Code != http.StatusOK {
		t.Errorf("Expected GetKubernetesNodesForCluster status 200, got %d", wClusterNodes.Code)
	}

	// 11. Test Sprint 10.8: GetKubernetesPodsForCluster
	reqClusterPods, _ := http.NewRequest("GET", "/api/v1/kubernetes/pods/"+testClusterID, nil)
	wClusterPods := httptest.NewRecorder()
	r.ServeHTTP(wClusterPods, reqClusterPods)
	if wClusterPods.Code != http.StatusOK {
		t.Errorf("Expected GetKubernetesPodsForCluster status 200, got %d", wClusterPods.Code)
	}

	// 12. Test Sprint 10.8: GetKubernetesDeploymentsForCluster
	reqClusterDeps, _ := http.NewRequest("GET", "/api/v1/kubernetes/deployments/"+testClusterID, nil)
	wClusterDeps := httptest.NewRecorder()
	r.ServeHTTP(wClusterDeps, reqClusterDeps)
	if wClusterDeps.Code != http.StatusOK {
		t.Errorf("Expected GetKubernetesDeploymentsForCluster status 200, got %d", wClusterDeps.Code)
	}

	// 13. Test Sprint 10.8: GetKubernetesStatefulSetsForCluster
	reqClusterSS, _ := http.NewRequest("GET", "/api/v1/kubernetes/statefulsets/"+testClusterID, nil)
	wClusterSS := httptest.NewRecorder()
	r.ServeHTTP(wClusterSS, reqClusterSS)
	if wClusterSS.Code != http.StatusOK {
		t.Errorf("Expected GetKubernetesStatefulSetsForCluster status 200, got %d", wClusterSS.Code)
	}

	// 14. Test Sprint 10.8: GetKubernetesDaemonSetsForCluster
	reqClusterDS, _ := http.NewRequest("GET", "/api/v1/kubernetes/daemonsets/"+testClusterID, nil)
	wClusterDS := httptest.NewRecorder()
	r.ServeHTTP(wClusterDS, reqClusterDS)
	if wClusterDS.Code != http.StatusOK {
		t.Errorf("Expected GetKubernetesDaemonSetsForCluster status 200, got %d", wClusterDS.Code)
	}

	// 15. Test Sprint 10.8: GetKubernetesServicesForCluster
	reqClusterSvcs, _ := http.NewRequest("GET", "/api/v1/kubernetes/services/"+testClusterID, nil)
	wClusterSvcs := httptest.NewRecorder()
	r.ServeHTTP(wClusterSvcs, reqClusterSvcs)
	if wClusterSvcs.Code != http.StatusOK {
		t.Errorf("Expected GetKubernetesServicesForCluster status 200, got %d", wClusterSvcs.Code)
	}

	// 16. Test Sprint 10.8: GetKubernetesNamespacesForCluster
	reqClusterNs, _ := http.NewRequest("GET", "/api/v1/kubernetes/namespaces/"+testClusterID, nil)
	wClusterNs := httptest.NewRecorder()
	r.ServeHTTP(wClusterNs, reqClusterNs)
	if wClusterNs.Code != http.StatusOK {
		t.Errorf("Expected GetKubernetesNamespacesForCluster status 200, got %d", wClusterNs.Code)
	}

	// 17. Test Sprint 10.8: GetKubernetesStorageForCluster
	reqClusterStore, _ := http.NewRequest("GET", "/api/v1/kubernetes/storage/"+testClusterID, nil)
	wClusterStore := httptest.NewRecorder()
	r.ServeHTTP(wClusterStore, reqClusterStore)
	if wClusterStore.Code != http.StatusOK {
		t.Errorf("Expected GetKubernetesStorageForCluster status 200, got %d", wClusterStore.Code)
	}

	// 18. Test Sprint 10.8: GetKubernetesEventsForCluster
	reqClusterEvts, _ := http.NewRequest("GET", "/api/v1/kubernetes/events/"+testClusterID, nil)
	wClusterEvts := httptest.NewRecorder()
	r.ServeHTTP(wClusterEvts, reqClusterEvts)
	if wClusterEvts.Code != http.StatusOK {
		t.Errorf("Expected GetKubernetesEventsForCluster status 200, got %d", wClusterEvts.Code)
	}

	// 19. Test Sprint 10.8: GetKubernetesPodLogsByPodID
	reqPodLogs, _ := http.NewRequest("GET", "/api/v1/kubernetes/logs/api-gateway-6f987c88b9-x2p8q?tail=10", nil)
	wPodLogs := httptest.NewRecorder()
	r.ServeHTTP(wPodLogs, reqPodLogs)
	if wPodLogs.Code != http.StatusOK {
		t.Errorf("Expected GetKubernetesPodLogsByPodID status 200, got %d", wPodLogs.Code)
	}

	var logsResp map[string]interface{}
	_ = json.Unmarshal(wPodLogs.Body.Bytes(), &logsResp)
	if logsStr, ok := logsResp["logs"].(string); !ok || logsStr == "" {
		t.Errorf("Expected non-empty logs payload")
	}
}
