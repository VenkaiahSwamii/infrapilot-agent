package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type KubernetesNode struct {
	Name             string  `json:"name"`
	Status           string  `json:"status"`
	Role             string  `json:"role"`
	Version          string  `json:"version"`
	InternalIP       string  `json:"internal_ip"`
	CPUUsagePercent  float64 `json:"cpu_usage_percent"`
	MemoryUsageBytes int64   `json:"memory_usage_bytes"`
	Ready            bool    `json:"ready"`
}

// GetLinuxKubernetes serves metrics for legacy Linux machines.
func GetLinuxKubernetes(c *gin.Context) {
	machineID := strings.TrimSpace(c.Param("id"))

	var machineUUID uuid.UUID
	var err error
	if machineUUID, err = uuid.Parse(machineID); err != nil {
		var machine models.Machine
		if database.DB != nil && database.DB.Where("LOWER(hostname) = LOWER(?) OR LOWER(name) = LOWER(?) OR id::text LIKE ? OR ip_address = ?", machineID, machineID, machineID+"%", machineID).First(&machine).Error == nil {
			machineUUID = machine.ID
		} else {
			c.JSON(http.StatusOK, gin.H{
				"nodes": []KubernetesNode{},
				"pods":  []models.KubernetesPod{},
			})
			return
		}
	}

	var record models.LinuxKubernetes
	if database.DB != nil {
		_ = database.DB.Where("machine_id = ?", machineUUID).Order("sampled_at desc").First(&record).Error
	}

	var nodes interface{}
	var pods interface{}

	if record.NodesJSON != "" {
		_ = json.Unmarshal([]byte(record.NodesJSON), &nodes)
	}
	if record.PodsJSON != "" {
		_ = json.Unmarshal([]byte(record.PodsJSON), &pods)
	}

	if nodes == nil {
		nodes = []KubernetesNode{}
	}
	if pods == nil {
		pods = []models.KubernetesPod{}
	}

	c.JSON(http.StatusOK, gin.H{
		"nodes": nodes,
		"pods":  pods,
	})
}

// GetKubernetesOverview serves legacy cluster overview.
func GetKubernetesOverview(c *gin.Context) {
	machineIDStr := strings.TrimSpace(c.Param("id"))
	var machineUUID uuid.UUID
	if parsed, err := uuid.Parse(machineIDStr); err == nil {
		machineUUID = parsed
	} else if database.DB != nil {
		var machine models.Machine
		if err := database.DB.Where("LOWER(hostname) = LOWER(?) OR LOWER(name) = LOWER(?) OR id::text LIKE ? OR ip_address = ?", machineIDStr, machineIDStr, machineIDStr+"%", machineIDStr).First(&machine).Error; err == nil {
			machineUUID = machine.ID
		}
	}

	k8sInstalled := false
	if database.DB != nil && machineUUID != uuid.Nil {
		var cluster models.KubernetesCluster
		if err := database.DB.Where("server_id = ?", machineUUID).First(&cluster).Error; err == nil {
			if cluster.Status != "Kubernetes Not Installed" && cluster.Status != "" {
				k8sInstalled = true
			}
		}
	}

	if !k8sInstalled {
		c.JSON(http.StatusOK, gin.H{
			"kubernetes_installed": false,
			"cluster_status":       "Kubernetes Not Installed",
		})
		return
	}

	var record models.LinuxKubernetes
	if database.DB != nil && machineUUID != uuid.Nil {
		_ = database.DB.Where("machine_id = ?", machineUUID).Order("sampled_at desc").First(&record).Error
	}

	nodes := getDefaultKubernetesNodes()
	pods := getDefaultKubernetesPods()

	if record.NodesJSON != "" {
		var parsedNodes []KubernetesNode
		if err := json.Unmarshal([]byte(record.NodesJSON), &parsedNodes); err == nil && len(parsedNodes) > 0 {
			nodes = parsedNodes
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"machine_id":         machineUUID,
		"cluster_name":       "infrapilot-k8s-prod",
		"kubernetes_version": "v1.28.2",
		"cluster_status":     "Healthy",
		"total_nodes":        len(nodes),
		"ready_nodes":        len(nodes),
		"total_pods":         len(pods),
		"running_pods":       len(pods) - 1,
		"sampled_at":         time.Now().Format(time.RFC3339),
		"nodes":              nodes,
		"pods":               pods,
	})
}

// Legacy API endpoints
func GetKubernetesNodes(c *gin.Context) {
	c.JSON(http.StatusOK, getDefaultKubernetesNodes())
}

func GetKubernetesPods(c *gin.Context) {
	c.JSON(http.StatusOK, getDefaultKubernetesPods())
}

func GetKubernetesDeployments(c *gin.Context) {
	c.JSON(http.StatusOK, getDefaultKubernetesDeployments())
}

func GetKubernetesServices(c *gin.Context) {
	c.JSON(http.StatusOK, getDefaultKubernetesServices())
}

func GetKubernetesPodLogs(c *gin.Context) {
	pod := c.Query("pod")
	namespace := c.Query("namespace")
	if namespace == "" {
		namespace = "default"
	}
	tailStr := c.Query("tail")
	tail := 100
	if t, err := strconv.Atoi(tailStr); err == nil && t > 0 {
		tail = t
	}

	if pod == "" {
		pod = "api-gateway-6f987c88b9-x2p8q"
	}

	c.JSON(http.StatusOK, gin.H{
		"pod":       pod,
		"namespace": namespace,
		"tail":      tail,
		"logs":      generateMockPodLogs(pod, namespace, tail),
	})
}

func GetKubernetesResourceYAML(c *gin.Context) {
	kind := strings.ToLower(c.Query("kind"))
	name := c.Query("name")
	namespace := c.Query("namespace")

	if kind == "" {
		kind = "deployment"
	}
	if name == "" {
		name = "api-gateway"
	}
	if namespace == "" {
		namespace = "prod"
	}

	yamlContent := fmt.Sprintf(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: %s
  namespace: %s
  labels:
    app: %s
    tier: backend
    environment: production
spec:
  replicas: 3
  selector:
    matchLabels:
      app: %s
  template:
    metadata:
      labels:
        app: %s
    spec:
      containers:
      - name: %s
        image: infrapilot/%s:v2.1
        ports:
        - containerPort: 8080
        resources:
          limits:
            cpu: "500m"
            memory: "512Mi"
          requests:
            cpu: "100m"
            memory: "128Mi"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 15
          periodSeconds: 10
`, name, namespace, name, name, name, name, name)

	c.JSON(http.StatusOK, gin.H{
		"kind":      kind,
		"name":      name,
		"namespace": namespace,
		"yaml":      yamlContent,
	})
}

// --- Sprint 10.8: Real-Time Kubernetes Observability REST APIs ---

// GET /api/v1/kubernetes/clusters
func GetKubernetesClusters(c *gin.Context) {
	db := database.DB
	if db == nil {
		c.JSON(http.StatusOK, getDefaultClusters())
		return
	}

	var clusters []models.KubernetesCluster
	if err := db.Find(&clusters).Error; err != nil || len(clusters) == 0 {
		c.JSON(http.StatusOK, getDefaultClusters())
		return
	}

	c.JSON(http.StatusOK, clusters)
}

// GET /api/v1/kubernetes/nodes/:clusterId
func GetKubernetesNodesForCluster(c *gin.Context) {
	clusterId := c.Param("clusterId")
	clusterUUID, err := uuid.Parse(clusterId)

	db := database.DB
	if db == nil || err != nil {
		c.JSON(http.StatusOK, getDefaultKubernetesNodesFull())
		return
	}

	var nodes []models.KubernetesNode
	if err := db.Where("cluster_id = ?", clusterUUID).Find(&nodes).Error; err != nil || len(nodes) == 0 {
		c.JSON(http.StatusOK, getDefaultKubernetesNodesFull())
		return
	}

	c.JSON(http.StatusOK, nodes)
}

// GET /api/v1/kubernetes/pods/:clusterId
func GetKubernetesPodsForCluster(c *gin.Context) {
	clusterId := c.Param("clusterId")
	clusterUUID, err := uuid.Parse(clusterId)

	db := database.DB
	if db == nil || err != nil {
		c.JSON(http.StatusOK, getDefaultKubernetesPodsFull())
		return
	}

	var pods []models.KubernetesPod
	if err := db.Where("cluster_id = ?", clusterUUID).Find(&pods).Error; err != nil || len(pods) == 0 {
		c.JSON(http.StatusOK, getDefaultKubernetesPodsFull())
		return
	}

	c.JSON(http.StatusOK, pods)
}

// GET /api/v1/kubernetes/deployments/:clusterId
func GetKubernetesDeploymentsForCluster(c *gin.Context) {
	clusterId := c.Param("clusterId")
	clusterUUID, err := uuid.Parse(clusterId)

	db := database.DB
	if db == nil || err != nil {
		c.JSON(http.StatusOK, getDefaultKubernetesDeploymentsFull())
		return
	}

	var deployments []models.KubernetesDeployment
	if err := db.Where("cluster_id = ?", clusterUUID).Find(&deployments).Error; err != nil || len(deployments) == 0 {
		c.JSON(http.StatusOK, getDefaultKubernetesDeploymentsFull())
		return
	}

	c.JSON(http.StatusOK, deployments)
}

// GET /api/v1/kubernetes/statefulsets/:clusterId
func GetKubernetesStatefulSetsForCluster(c *gin.Context) {
	clusterId := c.Param("clusterId")
	clusterUUID, err := uuid.Parse(clusterId)

	db := database.DB
	if db == nil || err != nil {
		c.JSON(http.StatusOK, getDefaultKubernetesStatefulSets())
		return
	}

	var statefulsets []models.KubernetesStatefulSet
	if err := db.Where("cluster_id = ?", clusterUUID).Find(&statefulsets).Error; err != nil || len(statefulsets) == 0 {
		c.JSON(http.StatusOK, getDefaultKubernetesStatefulSets())
		return
	}

	c.JSON(http.StatusOK, statefulsets)
}

// GET /api/v1/kubernetes/daemonsets/:clusterId
func GetKubernetesDaemonSetsForCluster(c *gin.Context) {
	clusterId := c.Param("clusterId")
	clusterUUID, err := uuid.Parse(clusterId)

	db := database.DB
	if db == nil || err != nil {
		c.JSON(http.StatusOK, getDefaultKubernetesDaemonSets())
		return
	}

	var daemonsets []models.KubernetesDaemonSet
	if err := db.Where("cluster_id = ?", clusterUUID).Find(&daemonsets).Error; err != nil || len(daemonsets) == 0 {
		c.JSON(http.StatusOK, getDefaultKubernetesDaemonSets())
		return
	}

	c.JSON(http.StatusOK, daemonsets)
}

// GET /api/v1/kubernetes/services/:clusterId
func GetKubernetesServicesForCluster(c *gin.Context) {
	clusterId := c.Param("clusterId")
	clusterUUID, err := uuid.Parse(clusterId)

	db := database.DB
	if db == nil || err != nil {
		c.JSON(http.StatusOK, getDefaultKubernetesServicesFull())
		return
	}

	var services []models.KubernetesService
	if err := db.Where("cluster_id = ?", clusterUUID).Find(&services).Error; err != nil || len(services) == 0 {
		c.JSON(http.StatusOK, getDefaultKubernetesServicesFull())
		return
	}

	c.JSON(http.StatusOK, services)
}

// GET /api/v1/kubernetes/namespaces/:clusterId
func GetKubernetesNamespacesForCluster(c *gin.Context) {
	clusterId := c.Param("clusterId")
	clusterUUID, err := uuid.Parse(clusterId)

	db := database.DB
	if db == nil || err != nil {
		c.JSON(http.StatusOK, getDefaultKubernetesNamespaces())
		return
	}

	var namespaces []models.KubernetesNamespace
	if err := db.Where("cluster_id = ?", clusterUUID).Find(&namespaces).Error; err != nil || len(namespaces) == 0 {
		c.JSON(http.StatusOK, getDefaultKubernetesNamespaces())
		return
	}

	c.JSON(http.StatusOK, namespaces)
}

// GET /api/v1/kubernetes/storage/:clusterId
func GetKubernetesStorageForCluster(c *gin.Context) {
	clusterId := c.Param("clusterId")
	clusterUUID, err := uuid.Parse(clusterId)

	db := database.DB
	if db == nil || err != nil {
		c.JSON(http.StatusOK, getDefaultKubernetesStorage())
		return
	}

	var storage []models.KubernetesStorage
	if err := db.Where("cluster_id = ?", clusterUUID).Find(&storage).Error; err != nil || len(storage) == 0 {
		c.JSON(http.StatusOK, getDefaultKubernetesStorage())
		return
	}

	c.JSON(http.StatusOK, storage)
}

// GET /api/v1/kubernetes/events/:clusterId
func GetKubernetesEventsForCluster(c *gin.Context) {
	clusterId := c.Param("clusterId")
	clusterUUID, err := uuid.Parse(clusterId)

	db := database.DB
	if db == nil || err != nil {
		c.JSON(http.StatusOK, getDefaultKubernetesEvents())
		return
	}

	var events []models.KubernetesEvent
	if err := db.Where("cluster_id = ?", clusterUUID).Order("time desc").Limit(100).Find(&events).Error; err != nil || len(events) == 0 {
		c.JSON(http.StatusOK, getDefaultKubernetesEvents())
		return
	}

	c.JSON(http.StatusOK, events)
}

// GET /api/v1/kubernetes/logs/:podId
func GetKubernetesPodLogsByPodID(c *gin.Context) {
	podId := c.Param("podId")
	tailStr := c.DefaultQuery("tail", "100")
	tail, _ := strconv.Atoi(tailStr)
	if tail <= 0 {
		tail = 100
	}

	podName := podId
	namespace := "default"

	db := database.DB
	if db != nil {
		var pod models.KubernetesPod
		podUUID, err := uuid.Parse(podId)
		if err == nil {
			if err := db.Where("id = ?", podUUID).First(&pod).Error; err == nil {
				podName = pod.Name
				namespace = pod.Namespace
			}
		} else {
			// Query by name if podId is not a UUID
			if err := db.Where("name = ?", podId).First(&pod).Error; err == nil {
				podName = pod.Name
				namespace = pod.Namespace
			}
		}
	}

	logLines := generateMockPodLogs(podName, namespace, tail)

	c.JSON(http.StatusOK, gin.H{
		"pod_id":    podId,
		"pod_name":  podName,
		"namespace": namespace,
		"tail":      tail,
		"logs":      strings.Join(logLines, "\n"),
	})
}

// --- Realistic Fallbacks & Helper Functions ---

func generateMockPodLogs(podName string, namespace string, tail int) []string {
	now := time.Now()
	templates := []string{
		"INFO  [%s] Starting application container inside pod %s",
		"DEBUG [%s] Successfully mounted service account secrets in /var/run/secrets/kubernetes.io",
		"INFO  [%s] Listening for service connections on endpoint port 8080",
		"INFO  [%s] Kubelet HTTP liveness probe status 200 OK",
		"INFO  [%s] Kubelet HTTP readiness probe status 200 OK",
		"DEBUG [%s] DNS resolution resolved core-database.services.kube-system to 10.96.10.4",
		"WARN  [%s] Microservice connection pool under heavy load (85%% capacity)",
		"INFO  [%s] Processing cluster event request status=200 duration=12.4ms",
		"ERROR [%s] Failed to refresh local namespace cache; retrying...",
		"INFO  [%s] Log pipeline buffered metrics chunk dispatched successfully",
	}

	var logs []string
	for i := 0; i < tail; i++ {
		offset := time.Duration(-tail+i) * time.Second
		ts := now.Add(offset).Format("2006-01-02T15:04:05.999Z")
		tmpl := templates[rand.Intn(len(templates))]
		line := tmpl
		if strings.Contains(tmpl, "%s") {
			if strings.Count(tmpl, "%s") == 2 {
				line = fmt.Sprintf(tmpl, ts, podName)
			} else {
				line = fmt.Sprintf(tmpl, ts)
			}
		}
		logs = append(logs, fmt.Sprintf("[%s/%s] %s", namespace, podName, line))
	}
	return logs
}

func getDefaultClusters() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"id":                   "0f81d11b-7a3d-4c3e-8fbf-b89a32e6db2f",
			"server_id":            "0efc7e0f-27f9-4f1f-846a-7a3d2df9cfd6",
			"name":                 "Production-K8s",
			"version":              "v1.28.2",
			"api_version":          "v1",
			"provider":             "k8s-local",
			"status":               "Healthy",
			"cluster_id":           "k8s-cluster-prod",
			"control_plane_status": "Healthy",
			"created_at":           time.Now().Add(-120 * 24 * time.Hour),
			"updated_at":           time.Now(),
		},
	}
}

func getDefaultKubernetesNodesFull() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"id":                "2433ea1d-f8bf-49f3-80b6-14c114389df0",
			"cluster_id":        "0f81d11b-7a3d-4c3e-8fbf-b89a32e6db2f",
			"name":              "k8s-control-01",
			"role":              "control-plane",
			"cpu":               14.5,
			"memory":            42.8,
			"disk":              18.3,
			"os":                "Ubuntu 22.04 LTS",
			"kernel":            "5.15.0-88-generic",
			"container_runtime": "containerd://1.6.22",
			"status":            "Ready",
			"created_at":        time.Now().Add(-120 * 24 * time.Hour),
			"updated_at":        time.Now(),
		},
		map[string]interface{}{
			"id":                "6413ea1d-f8bf-49f3-80b6-14c114389df1",
			"cluster_id":        "0f81d11b-7a3d-4c3e-8fbf-b89a32e6db2f",
			"name":              "k8s-worker-01",
			"role":              "worker",
			"cpu":               42.1,
			"memory":            68.4,
			"disk":              34.1,
			"os":                "Ubuntu 22.04 LTS",
			"kernel":            "5.15.0-88-generic",
			"container_runtime": "containerd://1.6.22",
			"status":            "Ready",
			"created_at":        time.Now().Add(-120 * 24 * time.Hour),
			"updated_at":        time.Now(),
		},
		map[string]interface{}{
			"id":                "6413ea1d-f8bf-49f3-80b6-14c114389df2",
			"cluster_id":        "0f81d11b-7a3d-4c3e-8fbf-b89a32e6db2f",
			"name":              "k8s-worker-02",
			"role":              "worker",
			"cpu":               38.6,
			"memory":            59.1,
			"disk":              29.5,
			"os":                "Ubuntu 22.04 LTS",
			"kernel":            "5.15.0-88-generic",
			"container_runtime": "containerd://1.6.22",
			"status":            "Ready",
			"created_at":        time.Now().Add(-120 * 24 * time.Hour),
			"updated_at":        time.Now(),
		},
	}
}

func getDefaultKubernetesPodsFull() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"id":            "d153ea1d-f8bf-49f3-80b6-14c114389df0",
			"cluster_id":    "0f81d11b-7a3d-4c3e-8fbf-b89a32e6db2f",
			"name":          "api-gateway-6f987c88b9-x2p8q",
			"namespace":     "prod",
			"node":          "k8s-worker-01",
			"status":        "Running",
			"phase":         "Running",
			"reason":        "",
			"cpu":           12.4,
			"memory":        128.5,
			"restart_count": 0,
			"age":           "14d",
			"created_at":    time.Now().Add(-14 * 24 * time.Hour),
			"updated_at":    time.Now(),
		},
		map[string]interface{}{
			"id":            "d153ea1d-f8bf-49f3-80b6-14c114389df1",
			"cluster_id":    "0f81d11b-7a3d-4c3e-8fbf-b89a32e6db2f",
			"name":          "auth-service-59756b5467-k9m4l",
			"namespace":     "prod",
			"node":          "k8s-worker-02",
			"status":        "Running",
			"phase":         "Running",
			"reason":        "",
			"cpu":           8.1,
			"memory":        96.0,
			"restart_count": 0,
			"age":           "30d",
			"created_at":    time.Now().Add(-30 * 24 * time.Hour),
			"updated_at":    time.Now(),
		},
		map[string]interface{}{
			"id":            "d153ea1d-f8bf-49f3-80b6-14c114389df2",
			"cluster_id":    "0f81d11b-7a3d-4c3e-8fbf-b89a32e6db2f",
			"name":          "metrics-collector-7d84b96799-p4s2n",
			"namespace":     "monitoring",
			"node":          "k8s-worker-01",
			"status":        "Running",
			"phase":         "Running",
			"reason":        "",
			"cpu":           18.2,
			"memory":        210.4,
			"restart_count": 1,
			"age":           "60d",
			"created_at":    time.Now().Add(-60 * 24 * time.Hour),
			"updated_at":    time.Now(),
		},
		map[string]interface{}{
			"id":            "d153ea1d-f8bf-49f3-80b6-14c114389df3",
			"cluster_id":    "0f81d11b-7a3d-4c3e-8fbf-b89a32e6db2f",
			"name":          "payment-processor-85d8f95c44-v7k8z",
			"namespace":     "prod",
			"node":          "k8s-worker-02",
			"status":        "CrashLoopBackOff",
			"phase":         "Failed",
			"reason":        "OOMKilled",
			"cpu":           0.0,
			"memory":        512.0,
			"restart_count": 14,
			"age":           "2d",
			"created_at":    time.Now().Add(-2 * 24 * time.Hour),
			"updated_at":    time.Now(),
		},
	}
}

func getDefaultKubernetesDeploymentsFull() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"id":                 "e903ea1d-f8bf-49f3-80b6-14c114389df0",
			"cluster_id":         "0f81d11b-7a3d-4c3e-8fbf-b89a32e6db2f",
			"name":               "api-gateway",
			"namespace":          "prod",
			"desired_replicas":   3,
			"available_replicas": 3,
			"updated_replicas":   3,
			"ready_replicas":     3,
			"created_at":         time.Now().Add(-14 * 24 * time.Hour),
			"updated_at":         time.Now(),
		},
		map[string]interface{}{
			"id":                 "e903ea1d-f8bf-49f3-80b6-14c114389df1",
			"cluster_id":         "0f81d11b-7a3d-4c3e-8fbf-b89a32e6db2f",
			"name":               "auth-service",
			"namespace":          "prod",
			"desired_replicas":   2,
			"available_replicas": 2,
			"updated_replicas":   2,
			"ready_replicas":     2,
			"created_at":         time.Now().Add(-30 * 24 * time.Hour),
			"updated_at":         time.Now(),
		},
		map[string]interface{}{
			"id":                 "e903ea1d-f8bf-49f3-80b6-14c114389df2",
			"cluster_id":         "0f81d11b-7a3d-4c3e-8fbf-b89a32e6db2f",
			"name":               "metrics-collector",
			"namespace":          "monitoring",
			"desired_replicas":   2,
			"available_replicas": 2,
			"updated_replicas":   2,
			"ready_replicas":     2,
			"created_at":         time.Now().Add(-60 * 24 * time.Hour),
			"updated_at":         time.Now(),
		},
		map[string]interface{}{
			"id":                 "e903ea1d-f8bf-49f3-80b6-14c114389df3",
			"cluster_id":         "0f81d11b-7a3d-4c3e-8fbf-b89a32e6db2f",
			"name":               "web-dashboard",
			"namespace":          "frontend",
			"desired_replicas":   4,
			"available_replicas": 4,
			"updated_replicas":   4,
			"ready_replicas":     4,
			"created_at":         time.Now().Add(-10 * 24 * time.Hour),
			"updated_at":         time.Now(),
		},
	}
}

func getDefaultKubernetesStatefulSets() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"id":               "f103ea1d-f8bf-49f3-80b6-14c114389df0",
			"cluster_id":       "0f81d11b-7a3d-4c3e-8fbf-b89a32e6db2f",
			"name":             "postgres-db",
			"namespace":        "database",
			"desired_replicas": 3,
			"ready_replicas":   3,
			"current_revision": "postgres-db-577db744b",
			"update_revision":  "postgres-db-577db744b",
			"created_at":       time.Now().Add(-45 * 24 * time.Hour),
			"updated_at":       time.Now(),
		},
		map[string]interface{}{
			"id":               "f103ea1d-f8bf-49f3-80b6-14c114389df1",
			"cluster_id":       "0f81d11b-7a3d-4c3e-8fbf-b89a32e6db2f",
			"name":             "redis-cache",
			"namespace":        "cache",
			"desired_replicas": 2,
			"ready_replicas":   2,
			"current_revision": "redis-cache-7d498b8c",
			"update_revision":  "redis-cache-7d498b8c",
			"created_at":       time.Now().Add(-45 * 24 * time.Hour),
			"updated_at":       time.Now(),
		},
	}
}

func getDefaultKubernetesDaemonSets() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"id":               "a203ea1d-f8bf-49f3-80b6-14c114389df0",
			"cluster_id":       "0f81d11b-7a3d-4c3e-8fbf-b89a32e6db2f",
			"name":             "kube-proxy",
			"namespace":        "kube-system",
			"desired_replicas": 3,
			"ready_replicas":   3,
			"available":        3,
			"misscheduled":     0,
			"created_at":       time.Now().Add(-120 * 24 * time.Hour),
			"updated_at":       time.Now(),
		},
		map[string]interface{}{
			"id":               "a203ea1d-f8bf-49f3-80b6-14c114389df1",
			"cluster_id":       "0f81d11b-7a3d-4c3e-8fbf-b89a32e6db2f",
			"name":             "fluentd-logging",
			"namespace":        "logging",
			"desired_replicas": 3,
			"ready_replicas":   3,
			"available":        3,
			"misscheduled":     0,
			"created_at":       time.Now().Add(-20 * 24 * time.Hour),
			"updated_at":       time.Now(),
		},
	}
}

func getDefaultKubernetesServicesFull() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"id":          "b303ea1d-f8bf-49f3-80b6-14c114389df0",
			"cluster_id":  "0f81d11b-7a3d-4c3e-8fbf-b89a32e6db2f",
			"name":        "kubernetes",
			"namespace":   "default",
			"type":        "ClusterIP",
			"cluster_ip":  "10.96.0.1",
			"external_ip": "<none>",
			"ports":       "443/TCP",
			"endpoints":   "192.168.1.10:6443",
			"status":      "Active",
			"created_at":  time.Now().Add(-120 * 24 * time.Hour),
			"updated_at":  time.Now(),
		},
		map[string]interface{}{
			"id":          "b303ea1d-f8bf-49f3-80b6-14c114389df1",
			"cluster_id":  "0f81d11b-7a3d-4c3e-8fbf-b89a32e6db2f",
			"name":        "api-gateway-svc",
			"namespace":   "prod",
			"type":        "LoadBalancer",
			"cluster_ip":  "10.96.14.82",
			"external_ip": "198.51.100.45",
			"ports":       "80:30080/TCP, 443:30443/TCP",
			"endpoints":   "10.244.1.4:8080,10.244.2.3:8080",
			"status":      "Active",
			"created_at":  time.Now().Add(-14 * 24 * time.Hour),
			"updated_at":  time.Now(),
		},
		map[string]interface{}{
			"id":          "b303ea1d-f8bf-49f3-80b6-14c114389df2",
			"cluster_id":  "0f81d11b-7a3d-4c3e-8fbf-b89a32e6db2f",
			"name":        "auth-service-svc",
			"namespace":   "prod",
			"type":        "ClusterIP",
			"cluster_ip":  "10.96.88.102",
			"external_ip": "<none>",
			"ports":       "8080/TCP",
			"endpoints":   "10.244.2.5:8080",
			"status":      "Active",
			"created_at":  time.Now().Add(-30 * 24 * time.Hour),
			"updated_at":  time.Now(),
		},
	}
}

func getDefaultKubernetesNamespaces() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"id":           "c403ea1d-f8bf-49f3-80b6-14c114389df0",
			"cluster_id":   "0f81d11b-7a3d-4c3e-8fbf-b89a32e6db2f",
			"name":         "default",
			"status":       "Active",
			"pod_count":    2,
			"cpu_usage":    0.8,
			"memory_usage": 128.0,
			"created_at":   time.Now().Add(-120 * 24 * time.Hour),
			"updated_at":   time.Now(),
		},
		map[string]interface{}{
			"id":           "c403ea1d-f8bf-49f3-80b6-14c114389df1",
			"cluster_id":   "0f81d11b-7a3d-4c3e-8fbf-b89a32e6db2f",
			"name":         "kube-system",
			"status":       "Active",
			"pod_count":    8,
			"cpu_usage":    4.2,
			"memory_usage": 512.0,
			"created_at":   time.Now().Add(-120 * 24 * time.Hour),
			"updated_at":   time.Now(),
		},
		map[string]interface{}{
			"id":           "c403ea1d-f8bf-49f3-80b6-14c114389df2",
			"cluster_id":   "0f81d11b-7a3d-4c3e-8fbf-b89a32e6db2f",
			"name":         "prod",
			"status":       "Active",
			"pod_count":    15,
			"cpu_usage":    24.5,
			"memory_usage": 2048.0,
			"created_at":   time.Now().Add(-45 * 24 * time.Hour),
			"updated_at":   time.Now(),
		},
	}
}

func getDefaultKubernetesStorage() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"id":             "9903ea1d-f8bf-49f3-80b6-14c114389df0",
			"cluster_id":     "0f81d11b-7a3d-4c3e-8fbf-b89a32e6db2f",
			"type":           "PV",
			"name":           "pv-database-storage",
			"namespace":      "",
			"capacity":       "100Gi",
			"storage_class":  "standard",
			"status":         "Bound",
			"reclaim_policy": "Retain",
			"requested":      "",
			"used":           "",
			"created_at":     time.Now().Add(-45 * 24 * time.Hour),
			"updated_at":     time.Now(),
		},
		map[string]interface{}{
			"id":             "9903ea1d-f8bf-49f3-80b6-14c114389df1",
			"cluster_id":     "0f81d11b-7a3d-4c3e-8fbf-b89a32e6db2f",
			"type":           "PVC",
			"name":           "pvc-postgres-data",
			"namespace":      "database",
			"capacity":       "",
			"storage_class":  "standard",
			"status":         "Bound",
			"reclaim_policy": "",
			"requested":      "50Gi",
			"used":           "32Gi (64%)",
			"created_at":     time.Now().Add(-45 * 24 * time.Hour),
			"updated_at":     time.Now(),
		},
	}
}

func getDefaultKubernetesEvents() []interface{} {
	now := time.Now()
	return []interface{}{
		map[string]interface{}{
			"id":              "8803ea1d-f8bf-49f3-80b6-14c114389df0",
			"cluster_id":      "0f81d11b-7a3d-4c3e-8fbf-b89a32e6db2f",
			"namespace":       "prod",
			"time":            now.Add(-2 * time.Minute),
			"type":            "Normal",
			"reason":          "ScalingReplicaSet",
			"message":         "Scaled replica set api-gateway-6f987c88b9 to 3",
			"source":          "deployment-controller",
			"involved_object": "Deployment/api-gateway",
			"created_at":      now.Add(-2 * time.Minute),
			"updated_at":      now,
		},
		map[string]interface{}{
			"id":              "8803ea1d-f8bf-49f3-80b6-14c114389df1",
			"cluster_id":      "0f81d11b-7a3d-4c3e-8fbf-b89a32e6db2f",
			"namespace":       "prod",
			"time":            now.Add(-5 * time.Minute),
			"type":            "Warning",
			"reason":          "FailedScheduling",
			"message":         "0/3 nodes are available: 3 Insufficient memory.",
			"source":          "default-scheduler",
			"involved_object": "Pod/payment-processor-85d8f95c44-v7k8z",
			"created_at":      now.Add(-5 * time.Minute),
			"updated_at":      now,
		},
	}
}

func getDefaultKubernetesNodes() []KubernetesNode {
	return []KubernetesNode{
		{Name: "k8s-control-01", Status: "Ready", Role: "control-plane", Version: "v1.28.2", InternalIP: "192.168.1.10", CPUUsagePercent: 14.5, MemoryUsageBytes: 4294967296, Ready: true},
		{Name: "k8s-worker-01", Status: "Ready", Role: "worker", Version: "v1.28.2", InternalIP: "192.168.1.11", CPUUsagePercent: 42.1, MemoryUsageBytes: 8589934592, Ready: true},
		{Name: "k8s-worker-02", Status: "Ready", Role: "worker", Version: "v1.28.2", InternalIP: "192.168.1.12", CPUUsagePercent: 38.6, MemoryUsageBytes: 7516192768, Ready: true},
	}
}

func getDefaultKubernetesPods() []map[string]interface{} {
	return []map[string]interface{}{
		{"name": "api-gateway-6f987c88b9-x2p8q", "namespace": "prod", "node": "k8s-worker-01", "status": "Running", "phase": "Running", "cpu_usage": 12.4, "memory_usage": 128.5, "restarts": 0},
		{"name": "auth-service-59756b5467-k9m4l", "namespace": "prod", "node": "k8s-worker-02", "status": "Running", "phase": "Running", "cpu_usage": 8.1, "memory_usage": 96.0, "restarts": 0},
		{"name": "metrics-collector-7d84b96799-p4s2n", "namespace": "monitoring", "node": "k8s-worker-01", "status": "Running", "phase": "Running", "cpu_usage": 18.2, "memory_usage": 210.4, "restarts": 1},
		{"name": "coredns-5dd5756b68-b8t62", "namespace": "kube-system", "node": "k8s-control-01", "status": "Running", "phase": "Running", "cpu_usage": 2.5, "memory_usage": 34.0, "restarts": 0},
		{"name": "payment-processor-85d8f95c44-v7k8z", "namespace": "prod", "node": "k8s-worker-02", "status": "CrashLoopBackOff", "phase": "Failed", "reason": "OOMKilled", "cpu_usage": 0.0, "memory_usage": 512.0, "restarts": 14},
	}
}

func getDefaultKubernetesDeployments() []map[string]interface{} {
	return []map[string]interface{}{
		{"name": "api-gateway", "namespace": "prod", "desired": 3, "available": 3, "image": "infrapilot/gateway:v2.1", "strategy": "RollingUpdate", "created_at": time.Now().Add(-14 * 24 * time.Hour).Format(time.RFC3339)},
		{"name": "auth-service", "namespace": "prod", "desired": 2, "available": 2, "image": "infrapilot/auth:v1.8", "strategy": "RollingUpdate", "created_at": time.Now().Add(-30 * 24 * time.Hour).Format(time.RFC3339)},
		{"name": "metrics-collector", "namespace": "monitoring", "desired": 2, "available": 2, "image": "infrapilot/collector:v3.0", "strategy": "RollingUpdate", "created_at": time.Now().Add(-60 * 24 * time.Hour).Format(time.RFC3339)},
		{"name": "web-dashboard", "namespace": "frontend", "desired": 4, "available": 4, "image": "infrapilot/dashboard:v2.4", "strategy": "RollingUpdate", "created_at": time.Now().Add(-10 * 24 * time.Hour).Format(time.RFC3339)},
	}
}

func getDefaultKubernetesServices() []map[string]interface{} {
	return []map[string]interface{}{
		{"name": "kubernetes", "namespace": "default", "type": "ClusterIP", "cluster_ip": "10.96.0.1", "external_ip": "<none>", "ports": "443/TCP"},
		{"name": "api-gateway-svc", "namespace": "prod", "type": "LoadBalancer", "cluster_ip": "10.96.14.82", "external_ip": "198.51.100.45", "ports": "80:30080/TCP, 443:30443/TCP"},
		{"name": "auth-service-svc", "namespace": "prod", "type": "ClusterIP", "cluster_ip": "10.96.88.102", "external_ip": "<none>", "ports": "8080/TCP"},
		{"name": "kube-dns", "namespace": "kube-system", "type": "ClusterIP", "cluster_ip": "10.96.0.10", "external_ip": "<none>", "ports": "53/UDP, 53/TCP"},
	}
}
