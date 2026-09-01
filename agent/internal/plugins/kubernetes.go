package plugins

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"infrapilot/agent/internal/metrics"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// KubernetesNodeMetric represents a Kubernetes node metric
type KubernetesNodeMetric struct {
	Name             string  `json:"name"`
	Status           string  `json:"status"`
	Role             string  `json:"role"`
	Version          string  `json:"version"`
	CPUUsagePercent  float64 `json:"cpu_usage_percent"`
	MemoryUsageBytes int64   `json:"memory_usage_bytes"`
}

// KubernetesPlugin collects Kubernetes cluster metrics
type KubernetesPlugin struct {
	enabled   bool
	mu        sync.RWMutex
	health    PluginHealth
	clientset *kubernetes.Clientset
}

// NewKubernetesPlugin creates a new Kubernetes plugin instance
func NewKubernetesPlugin() *KubernetesPlugin {
	plugin := &KubernetesPlugin{
		enabled: true,
		health:  NewPluginHealth("kubernetes"),
	}

	// Try to initialize Kubernetes client
	home, err := os.UserHomeDir()
	if err != nil {
		plugin.health.Status = "degraded"
		plugin.health.Error = "Failed to get home directory: " + err.Error()
		plugin.health.Metadata["mode"] = "unavailable"
		return plugin
	}

	kubeconfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		plugin.health.Status = "degraded"
		plugin.health.Error = "Failed to load kubeconfig: " + err.Error()
		plugin.health.Metadata["mode"] = "unavailable"
		return plugin
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		plugin.health.Status = "degraded"
		plugin.health.Error = "Failed to create Kubernetes client: " + err.Error()
		plugin.health.Metadata["mode"] = "unavailable"
		return plugin
	}

	plugin.clientset = clientset
	plugin.health.Status = "healthy"
	plugin.health.Metadata["mode"] = "live"

	return plugin
}

// Name returns the plugin name
func (p *KubernetesPlugin) Name() string {
	return "kubernetes"
}

// IsEnabled returns whether the plugin is enabled
func (p *KubernetesPlugin) IsEnabled() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.enabled
}

// SetEnabled enables or disables the plugin
func (p *KubernetesPlugin) SetEnabled(enabled bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.enabled = enabled
}

// HealthCheck returns the health status of the plugin
func (p *KubernetesPlugin) HealthCheck() PluginHealth {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.health
}

// Collect gathers Kubernetes cluster metrics
func (p *KubernetesPlugin) Collect() (interface{}, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	startTime := time.Now()

	result := make(map[string]interface{})

	// Collect pods
	pods, err := p.collectPods()
	if err != nil {
		p.health.Status = "degraded"
		p.health.Error = err.Error()
		result["error"] = err.Error()
		result["mode"] = "unavailable"
		result["nodes"] = []KubernetesNodeMetric{}
		result["pods"] = []metrics.KubernetesPodMetric{}
	} else {
		p.health.Status = "healthy"
		p.health.Error = ""
		result["pods"] = pods
		result["mode"] = p.health.Metadata["mode"]

		nodes, nodeErr := p.collectNodes()
		if nodeErr != nil {
			return nil, nodeErr
		}
		result["nodes"] = nodes
	}

	p.health.LastRun = time.Now()
	p.health.Metadata["duration_ms"] = time.Since(startTime).String()

	return result, nil
}

func (p *KubernetesPlugin) collectNodes() ([]KubernetesNodeMetric, error) {
	if p.clientset == nil {
		return nil, fmt.Errorf("Kubernetes client not initialized")
	}
	nodeList, err := p.clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	nodes := make([]KubernetesNodeMetric, 0, len(nodeList.Items))
	for _, node := range nodeList.Items {
		status := "NotReady"
		for _, condition := range node.Status.Conditions {
			if condition.Type == v1.NodeReady && condition.Status == v1.ConditionTrue {
				status = "Ready"
				break
			}
		}
		nodes = append(nodes, KubernetesNodeMetric{Name: node.Name, Status: status, Version: node.Status.NodeInfo.KubeletVersion})
	}
	return nodes, nil
}

// collectPods gathers Kubernetes pod information
func (p *KubernetesPlugin) collectPods() ([]metrics.KubernetesPodMetric, error) {
	if p.clientset == nil {
		return nil, fmt.Errorf("Kubernetes client not initialized")
	}

	podList, err := p.clientset.CoreV1().Pods("").List(
		context.Background(),
		metav1.ListOptions{},
	)
	if err != nil {
		return nil, err
	}

	pods := make([]metrics.KubernetesPodMetric, 0, len(podList.Items))
	for _, pod := range podList.Items {
		pods = append(pods, metrics.KubernetesPodMetric{
			Name:         pod.Name,
			Namespace:    pod.Namespace,
			Node:         pod.Spec.NodeName,
			Status:       string(pod.Status.Phase),
			CPUUsage:     0,
			MemoryUsage:  0,
			RestartCount: totalRestarts(pod),
		})
	}

	return pods, nil
}

// totalRestarts calculates total restart count for a pod
func totalRestarts(pod v1.Pod) int {
	total := 0
	for _, c := range pod.Status.ContainerStatuses {
		total += int(c.RestartCount)
	}
	return total
}
