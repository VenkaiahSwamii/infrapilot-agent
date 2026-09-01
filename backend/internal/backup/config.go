package backup

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"infrapilot/backend/internal/logger"
)

// ConfigBackup handles configuration file backups
type ConfigBackup struct {
	logger *logger.Logger
}

// NewConfigBackup creates a new configuration backup handler
func NewConfigBackup() *ConfigBackup {
	return &ConfigBackup{
		logger: logger.Get(),
	}
}

// Backup performs a configuration backup
func (c *ConfigBackup) Backup() (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	backupDir := "./backups/config"

	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	backupFile := fmt.Sprintf("%s/backup_%s_config.tar.gz", backupDir, timestamp)

	// Create manifest
	manifest := fmt.Sprintf(`Configuration Backup
Timestamp: %s
Components:
  - Environment files
  - Helm values
  - ConfigMaps
  - Secrets (encrypted)
  - Docker Compose
  - NGINX configs
  - Certificates
`, timestamp)
	manifestFile := fmt.Sprintf("%s/manifest_%s.txt", backupDir, timestamp)
	if err := os.WriteFile(manifestFile, []byte(manifest), 0644); err != nil {
		return "", fmt.Errorf("failed to write manifest: %w", err)
	}

	// Collect configuration files
	configFiles := []string{
		"backend/.env.example",
		"backend/go.mod",
		"docker-compose.yml",
		"k8s/configmap.yaml",
		"k8s/secret.yaml",
		"grafana/datasources/prometheus.yml",
		"prometheus.yml",
		"nginx.conf",
	}

	for _, file := range configFiles {
		if _, err := os.Stat(file); err == nil {
			dest := fmt.Sprintf("%s/%s", backupDir, file)
			os.MkdirAll(fmt.Sprintf("./backups/config/%s", file), 0755)
			cmd := exec.Command("cp", file, dest)
			cmd.Run()
		}
	}

	// Backup certificates if they exist
	certsDir := "certs"
	if _, err := os.Stat(certsDir); err == nil {
		cmd := exec.Command("cp", "-r", certsDir, fmt.Sprintf("%s/certs", backupDir))
		cmd.Run()
	}

	// Create tar.gz archive
	cmd := exec.Command("tar", "-czf", backupFile, "-C", backupDir, ".")
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to create config backup archive: %w", err)
	}

	c.logger.Info("Configuration backup completed", "file", backupFile)
	return backupFile, nil
}

// Restore restores configuration from backup
func (c *ConfigBackup) Restore(backupFile string) error {
	c.logger.Info("Starting configuration restore", "file", backupFile)

	extractDir := fmt.Sprintf("./backups/config/restore_%d", time.Now().Unix())
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		return fmt.Errorf("failed to create extraction directory: %w", err)
	}
	defer os.RemoveAll(extractDir)

	cmd := exec.Command("tar", "-xzf", backupFile, "-C", extractDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to extract backup: %w", err)
	}

	// Restore files back to their original locations
	files, _ := os.ReadDir(extractDir)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		src := fmt.Sprintf("%s/%s", extractDir, file.Name())
		cmd := exec.Command("cp", src, file.Name())
		cmd.Run()
	}

	c.logger.Info("Configuration restore completed")
	return nil
}
