package alerts

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"infrapilot/backend/internal/models"
	"infrapilot/backend/internal/services"
)

type K8sPodPayload struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Status    string `json:"status"`
	Phase     string `json:"phase"`
	Reason    string `json:"reason"`
	Restarts  int    `json:"restarts"`
}

type K8sNodePayload struct {
	Name   string  `json:"name"`
	Status string  `json:"status"`
	CPU    float64 `json:"cpu"`
	Memory float64 `json:"memory"`
}

type K8sClusterPayload struct {
	Name               string `json:"name"`
	Status             string `json:"status"`
	ControlPlaneStatus string `json:"control_plane_status"`
}

type K8sDeploymentPayload struct {
	Name              string `json:"name"`
	Namespace         string `json:"namespace"`
	DesiredReplicas   int    `json:"desired_replicas"`
	AvailableReplicas int    `json:"available_replicas"`
}

type K8sStoragePayload struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Status    string `json:"status"`
	Used      string `json:"used"` // e.g., "92%"
}

// CheckKubernetesHealth evaluates Kubernetes cluster health, node resources, workloads, storage, and event states.
func CheckKubernetesHealth(machine models.Machine, metric models.Metric, input services.SaveMetricInput) []models.LinuxAlert {
	var alerts []models.LinuxAlert

	// 1. API Server / Control Plane Unreachable Check
	if input.K8sInstalled && input.K8sClusterJSON != "" {
		var cluster K8sClusterPayload
		if err := json.Unmarshal([]byte(input.K8sClusterJSON), &cluster); err == nil {
			statusStr := strings.ToLower(cluster.Status + " " + cluster.ControlPlaneStatus)
			if strings.Contains(statusStr, "unreachable") || strings.Contains(statusStr, "unhealthy") || strings.Contains(statusStr, "unavailable") {
				alerts = append(alerts, models.LinuxAlert{
					Title:              "Kubernetes Control Plane Unhealthy",
					Description:        fmt.Sprintf("Kubernetes cluster '%s' control plane is unhealthy: %s (Status: %s)", cluster.Name, cluster.ControlPlaneStatus, cluster.Status),
					Category:           "Kubernetes",
					Component:          "kubernetes_control_plane",
					Source:             "KubernetesChecker",
					Type:               "k8s_control_plane_unhealthy",
					Severity:           "Critical",
					Priority:           models.MapSeverityToPriority("Critical"),
					Message:            fmt.Sprintf("Cluster %s Control Plane Unhealthy", cluster.Name),
					Status:             "OPEN",
					RecoverySuggestion: "Verify API server connectivity, check etcd status, and inspect kube-apiserver system logs.",
				})
			}
		}
	}

	// 2. Kubernetes Pods Check (CrashLoopBackOff, Image Pull Failure, Pending, Evicted)
	if input.K8sPodsJSON != "" && input.K8sPodsJSON != "[]" {
		var pods []K8sPodPayload
		if err := json.Unmarshal([]byte(input.K8sPodsJSON), &pods); err == nil {
			for _, pod := range pods {
				statusStr := strings.ToLower(pod.Status + " " + pod.Reason + " " + pod.Phase)

				if strings.Contains(statusStr, "crashloopbackoff") || strings.Contains(statusStr, "oomkilled") || strings.Contains(statusStr, "imagepullbackoff") || strings.Contains(statusStr, "errimagepull") {
					alerts = append(alerts, models.LinuxAlert{
						Title:              "Kubernetes Pod CrashLoopBackOff",
						Description:        fmt.Sprintf("Kubernetes pod '%s/%s' in failing state: %s (Reason: %s, Restarts: %d)", pod.Namespace, pod.Name, pod.Status, pod.Reason, pod.Restarts),
						Category:           "Kubernetes",
						Component:          fmt.Sprintf("pod_%s", pod.Name),
						Source:             "KubernetesChecker",
						Type:               "k8s_pod_crash",
						Severity:           "Critical",
						Priority:           models.MapSeverityToPriority("Critical"),
						Message:            fmt.Sprintf("Pod %s/%s CrashLoopBackOff", pod.Namespace, pod.Name),
						Status:             "OPEN",
						RecoverySuggestion: fmt.Sprintf("Check pod logs (`kubectl logs %s -n %s`) and describe pod (`kubectl describe pod %s -n %s`).", pod.Name, pod.Namespace, pod.Name, pod.Namespace),
					})
				} else if strings.Contains(statusStr, "pending") || strings.Contains(statusStr, "evicted") || strings.Contains(statusStr, "failed") {
					alerts = append(alerts, models.LinuxAlert{
						Title:              "Kubernetes Pod Pending/Evicted",
						Description:        fmt.Sprintf("Kubernetes pod '%s/%s' in non-running phase: %s", pod.Namespace, pod.Name, pod.Status),
						Category:           "Kubernetes",
						Component:          fmt.Sprintf("pod_%s", pod.Name),
						Source:             "KubernetesChecker",
						Type:               "k8s_pod_pending",
						Severity:           "Warning",
						Priority:           models.MapSeverityToPriority("Warning"),
						Message:            fmt.Sprintf("Pod %s/%s Pending or Evicted", pod.Namespace, pod.Name),
						Status:             "OPEN",
						RecoverySuggestion: fmt.Sprintf("Check node resource capacity (`kubectl top nodes`) or pod event warnings (`kubectl get events -n %s`).", pod.Namespace),
					})
				}
			}
		}
	}

	// 3. Kubernetes Nodes Check (NotReady, CPU & Memory Usage > 90%)
	if input.K8sNodesJSON != "" && input.K8sNodesJSON != "[]" {
		var nodes []K8sNodePayload
		if err := json.Unmarshal([]byte(input.K8sNodesJSON), &nodes); err == nil {
			for _, node := range nodes {
				statusStr := strings.ToLower(node.Status)
				if strings.Contains(statusStr, "notready") || strings.Contains(statusStr, "unknown") || strings.Contains(statusStr, "offline") {
					alerts = append(alerts, models.LinuxAlert{
						Title:              "Kubernetes Node NotReady",
						Description:        fmt.Sprintf("Kubernetes cluster node '%s' is NotReady or disconnected", node.Name),
						Category:           "Kubernetes",
						Component:          fmt.Sprintf("k8s_node_%s", node.Name),
						Source:             "KubernetesChecker",
						Type:               "k8s_node_notready",
						Severity:           "Critical",
						Priority:           models.MapSeverityToPriority("Critical"),
						Message:            fmt.Sprintf("Node %s NotReady", node.Name),
						Status:             "OPEN",
						RecoverySuggestion: fmt.Sprintf("Verify kubelet service status on node '%s' (`systemctl status kubelet`) and node network connectivity.", node.Name),
					})
				}

				if node.CPU > 90.0 {
					alerts = append(alerts, models.LinuxAlert{
						Title:              "Kubernetes Node High CPU Usage",
						Description:        fmt.Sprintf("Kubernetes node '%s' CPU usage is critical: %.2f%%", node.Name, node.CPU),
						Category:           "Kubernetes",
						Component:          fmt.Sprintf("k8s_node_%s", node.Name),
						Source:             "KubernetesChecker",
						Type:               "k8s_node_cpu_high",
						Severity:           "Warning",
						Priority:           models.MapSeverityToPriority("Warning"),
						Message:            fmt.Sprintf("Node %s High CPU: %.2f%%", node.Name, node.CPU),
						Status:             "OPEN",
						RecoverySuggestion: fmt.Sprintf("Identify resource-intensive pods on node '%s' using `kubectl top pods --all-namespaces` and consider scaling workloads.", node.Name),
					})
				}

				if node.Memory > 90.0 {
					alerts = append(alerts, models.LinuxAlert{
						Title:              "Kubernetes Node High Memory Usage",
						Description:        fmt.Sprintf("Kubernetes node '%s' memory usage is critical: %.2f%%", node.Name, node.Memory),
						Category:           "Kubernetes",
						Component:          fmt.Sprintf("k8s_node_%s", node.Name),
						Source:             "KubernetesChecker",
						Type:               "k8s_node_memory_high",
						Severity:           "Warning",
						Priority:           models.MapSeverityToPriority("Warning"),
						Message:            fmt.Sprintf("Node %s High Memory: %.2f%%", node.Name, node.Memory),
						Status:             "OPEN",
						RecoverySuggestion: fmt.Sprintf("Check node '%s' memory consumers and adjust pod resource requests or memory limits.", node.Name),
					})
				}
			}
		}
	}

	// 4. Deployment Replica Mismatch Check
	if input.K8sDeploymentsJSON != "" && input.K8sDeploymentsJSON != "[]" {
		var deployments []K8sDeploymentPayload
		if err := json.Unmarshal([]byte(input.K8sDeploymentsJSON), &deployments); err == nil {
			for _, dep := range deployments {
				if dep.DesiredReplicas != dep.AvailableReplicas {
					alerts = append(alerts, models.LinuxAlert{
						Title:              "Kubernetes Deployment Replica Mismatch",
						Description:        fmt.Sprintf("Deployment '%s/%s' replica counts mismatch. Desired: %d, Available: %d", dep.Namespace, dep.Name, dep.DesiredReplicas, dep.AvailableReplicas),
						Category:           "Kubernetes",
						Component:          fmt.Sprintf("deployment_%s", dep.Name),
						Source:             "KubernetesChecker",
						Type:               "k8s_deployment_mismatch",
						Severity:           "Warning",
						Priority:           models.MapSeverityToPriority("Warning"),
						Message:            fmt.Sprintf("Deployment %s/%s Replica Mismatch", dep.Namespace, dep.Name),
						Status:             "OPEN",
						RecoverySuggestion: fmt.Sprintf("Run `kubectl describe deployment %s -n %s` to check controller rollouts or events.", dep.Name, dep.Namespace),
					})
				}
			}
		}
	}

	// 5. PVC Storage Full Check (> 90%)
	if input.K8sStorageJSON != "" && input.K8sStorageJSON != "[]" {
		var storageList []K8sStoragePayload
		if err := json.Unmarshal([]byte(input.K8sStorageJSON), &storageList); err == nil {
			for _, st := range storageList {
				if st.Type == "PVC" && st.Used != "" {
					usedClean := strings.TrimSuffix(st.Used, "%")
					if usedPct, err := strconv.ParseFloat(usedClean, 64); err == nil && usedPct > 90.0 {
						alerts = append(alerts, models.LinuxAlert{
							Title:              "Kubernetes PVC Storage Full",
							Description:        fmt.Sprintf("Kubernetes Persistent Volume Claim '%s/%s' storage is almost full: %s usage", st.Namespace, st.Name, st.Used),
							Category:           "Kubernetes",
							Component:          fmt.Sprintf("pvc_%s", st.Name),
							Source:             "KubernetesChecker",
							Type:               "k8s_pvc_full",
							Severity:           "Warning",
							Priority:           models.MapSeverityToPriority("Warning"),
							Message:            fmt.Sprintf("PVC %s/%s Storage Full: %s", st.Namespace, st.Name, st.Used),
							Status:             "OPEN",
							RecoverySuggestion: fmt.Sprintf("Extend the PVC size using storage-class resize capabilities or clean up temporary files in PVC '%s'.", st.Name),
						})
					}
				}
			}
		}
	}

	return alerts
}
