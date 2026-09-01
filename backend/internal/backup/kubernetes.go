package backup

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"infrapilot/backend/internal/logger"
)

// KubernetesBackup handles Kubernetes resource backups
type KubernetesBackup struct {
	logger *logger.Logger
}

// NewKubernetesBackup creates a new Kubernetes backup handler
func NewKubernetesBackup() *KubernetesBackup {
	return &KubernetesBackup{
		logger: logger.Get(),
	}
}

// Backup performs a Kubernetes resources backup
func (k *KubernetesBackup) Backup() (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	backupDir := "./backups/kubernetes"

	// Ensure directory exists
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	backupFile := fmt.Sprintf("%s/backup_%s_kubernetes.tar.gz", backupDir, timestamp)

	// Create manifest file
	manifest := fmt.Sprintf(`Kubernetes Resources Backup
Timestamp: %s
Cluster: %s
`, timestamp, getEnvOrDefault("KUBE_CLUSTER_NAME", "default"))

	manifestFile := fmt.Sprintf("%s/manifest_%s.txt", backupDir, timestamp)
	if err := os.WriteFile(manifestFile, []byte(manifest), 0644); err != nil {
		return "", fmt.Errorf("failed to write manifest: %w", err)
	}

	// Export all resources using kubectl
	// We'll create individual YAML files for each resource type
	resources := []struct {
		name string
		args []string
	}{
		{"deployments", []string{"get", "deployments", "--all-namespaces", "-o", "yaml"}},
		{"statefulsets", []string{"get", "statefulsets", "--all-namespaces", "-o", "yaml"}},
		{"daemonsets", []string{"get", "daemonsets", "--all-namespaces", "-o", "yaml"}},
		{"services", []string{"get", "services", "--all-namespaces", "-o", "yaml"}},
		{"configmaps", []string{"get", "configmaps", "--all-namespaces", "-o", "yaml"}},
		{"secrets", []string{"get", "secrets", "--all-namespaces", "-o", "yaml"}},
		{"persistentvolumeclaims", []string{"get", "persistentvolumeclaims", "--all-namespaces", "-o", "yaml"}},
		{"persistentvolumes", []string{"get", "persistentvolumes", "-o", "yaml"}},
		{"ingresses", []string{"get", "ingresses", "--all-namespaces", "-o", "yaml"}},
		{"serviceaccounts", []string{"get", "serviceaccounts", "--all-namespaces", "-o", "yaml"}},
		{"roles", []string{"get", "roles", "--all-namespaces", "-o", "yaml"}},
		{"rolebindings", []string{"get", "rolebindings", "--all-namespaces", "-o", "yaml"}},
		{"horizontalpodautoscalers", []string{"get", "horizontalpodautoscalers", "--all-namespaces", "-o", "yaml"}},
	}

	for _, resource := range resources {
		resourceFile := fmt.Sprintf("%s/%s.yaml", backupDir, resource.name)
		cmd := exec.Command("kubectl", resource.args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// Capture output to file
		output, err := cmd.Output()
		if err != nil {
			k.logger.Warn("Failed to backup kubernetes resource",
				"resource", resource.name,
				"error", err)
			continue
		}

		if len(output) > 0 {
			if err := os.WriteFile(resourceFile, output, 0644); err != nil {
				k.logger.Warn("Failed to write resource file",
					"resource", resource.name,
					"error", err)
			}
		}
	}

	// Create tar.gz archive
	cmd := exec.Command("tar", "-czf", backupFile, "-C", backupDir, ".")
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to create kubernetes backup archive: %w", err)
	}

	k.logger.Info("Kubernetes backup completed", "file", backupFile)
	return backupFile, nil
}

// Restore restores Kubernetes resources from backup
func (k *KubernetesBackup) Restore(backupFile string) error {
	k.logger.Info("Starting Kubernetes restore", "file", backupFile)

	// Create temporary extraction directory
	extractDir := fmt.Sprintf("./backups/kubernetes/restore_%d", time.Now().Unix())
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		return fmt.Errorf("failed to create extraction directory: %w", err)
	}
	defer os.RemoveAll(extractDir)

	// Extract tar.gz
	cmd := exec.Command("tar", "-xzf", backupFile, "-C", extractDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to extract backup: %w", err)
	}

	// Apply each resource file
	files, err := os.ReadDir(extractDir)
	if err != nil {
		return fmt.Errorf("failed to read extraction directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || file.Name() == "manifest.txt" {
			continue
		}

		resourcePath := fmt.Sprintf("%s/%s", extractDir, file.Name())
		applyCmd := exec.Command("kubectl", "apply", "-f", resourcePath)
		applyCmd.Stdout = os.Stdout
		applyCmd.Stderr = os.Stderr

		if err := applyCmd.Run(); err != nil {
			k.logger.Warn("Failed to apply kubernetes resource",
				"file", file.Name(),
				"error", err)
		}
	}

	k.logger.Info("Kubernetes restore completed")
	return nil
}

// BackupNamespace backs up a specific namespace
func (k *KubernetesBackup) BackupNamespace(namespace string) (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	backupDir := "./backups/kubernetes"

	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	backupFile := fmt.Sprintf("%s/backup_%s_namespace_%s.tar.gz", backupDir, namespace, timestamp)

	// Export all resources in the namespace
	cmd := exec.Command("kubectl", "get", "all", "-n", namespace, "-o", "yaml")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to export namespace resources: %w", err)
	}

	// Write namespace resources
	nsFile := fmt.Sprintf("%s/%s_all.yaml", backupDir, namespace)
	if err := os.WriteFile(nsFile, output, 0644); err != nil {
		return "", fmt.Errorf("failed to write namespace file: %w", err)
	}

	// Create archive
	tarCmd := exec.Command("tar", "-czf", backupFile, "-C", backupDir, fmt.Sprintf("%s_all.yaml", namespace))
	if err := tarCmd.Run(); err != nil {
		return "", fmt.Errorf("failed to create namespace backup archive: %w", err)
	}

	k.logger.Info("Kubernetes namespace backup completed",
		"namespace", namespace,
		"file", backupFile)

	return backupFile, nil
}

// GetClusterInfo retrieves cluster information
func (k *KubernetesBackup) GetClusterInfo() (map[string]string, error) {
	info := make(map[string]string)

	// Get cluster name
	cmd := exec.Command("kubectl", "config", "current-context")
	if output, err := cmd.Output(); err == nil {
		info["context"] = string(output)
	}

	// Get nodes
	cmd = exec.Command("kubectl", "get", "nodes", "-o", "wide")
	if output, err := cmd.Output(); err == nil {
		info["nodes"] = string(output)
	}

	// Get namespaces
	cmd = exec.Command("kubectl", "get", "namespaces")
	if output, err := cmd.Output(); err == nil {
		info["namespaces"] = string(output)
	}

	return info, nil
}
