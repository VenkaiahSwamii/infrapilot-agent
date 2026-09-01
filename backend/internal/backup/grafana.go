package backup

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

	"infrapilot/backend/internal/logger"
)

// GrafanaBackup handles Grafana backup operations
type GrafanaBackup struct {
	logger *logger.Logger
}

// NewGrafanaBackup creates a new Grafana backup handler
func NewGrafanaBackup() *GrafanaBackup {
	return &GrafanaBackup{
		logger: logger.Get(),
	}
}

// GrafanaDashboard represents a Grafana dashboard
type GrafanaDashboard struct {
	UID   string `json:"uid"`
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

// Backup performs a Grafana backup
func (g *GrafanaBackup) Backup() (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	backupDir := "./backups/grafana"

	// Ensure directory exists
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	backupFile := fmt.Sprintf("%s/backup_%s_grafana.tar.gz", backupDir, timestamp)

	// Get Grafana configuration from environment
	grafanaURL := getEnvOrDefault("GRAFANA_URL", "http://localhost:3000")
	grafanaAPIKey := getEnvOrDefault("GRAFANA_API_KEY", "")

	// Backup dashboards
	if err := g.backupDashboards(grafanaURL, grafanaAPIKey, backupDir, timestamp); err != nil {
		g.logger.Warn("Failed to backup dashboards", "error", err)
	}

	// Backup datasources
	if err := g.backupDatasources(grafanaURL, grafanaAPIKey, backupDir, timestamp); err != nil {
		g.logger.Warn("Failed to backup datasources", "error", err)
	}

	// Backup alert rules
	if err := g.backupAlertRules(grafanaURL, grafanaAPIKey, backupDir, timestamp); err != nil {
		g.logger.Warn("Failed to backup alert rules", "error", err)
	}

	// Backup notification channels
	if err := g.backupNotificationChannels(grafanaURL, grafanaAPIKey, backupDir, timestamp); err != nil {
		g.logger.Warn("Failed to backup notification channels", "error", err)
	}

	// Create tar.gz archive
	cmd := exec.Command("tar", "-czf", backupFile, "-C", backupDir, ".")
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to create grafana backup archive: %w", err)
	}

	g.logger.Info("Grafana backup completed", "file", backupFile)
	return backupFile, nil
}

// backupDashbacks backs up all Grafana dashboards
func (g *GrafanaBackup) backupDashboards(grafanaURL, apiKey, backupDir, timestamp string) error {
	g.logger.Info("Backing up Grafana dashboards")

	// Create dashboards directory
	dashboardsDir := fmt.Sprintf("%s/dashboards", backupDir)
	if err := os.MkdirAll(dashboardsDir, 0755); err != nil {
		return err
	}

	// Search for all dashboards
	searchCmd := exec.Command("curl", "-s",
		fmt.Sprintf("%s/api/search?query=&limit=1000", grafanaURL),
		"-H", "Authorization: Bearer "+apiKey,
	)

	output, err := searchCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to search dashboards: %w", err)
	}

	var searchResults []struct {
		UID   string `json:"uid"`
		Slug  string `json:"slug"`
		Title string `json:"title"`
	}

	if err := json.Unmarshal(output, &searchResults); err != nil {
		return fmt.Errorf("failed to parse search results: %w", err)
	}

	// Download each dashboard
	for _, result := range searchResults {
		if result.UID == "" {
			continue
		}

		// Get dashboard JSON
		dashboardCmd := exec.Command("curl", "-s",
			fmt.Sprintf("%s/api/dashboards/uid/%s", grafanaURL, result.UID),
			"-H", "Authorization: Bearer "+apiKey)

		dashboardOutput, err := dashboardCmd.Output()
		if err != nil {
			g.logger.Warn("Failed to backup dashboard", "uid", result.UID, "error", err)
			continue
		}

		// Parse and save dashboard without metadata
		var dashboardData struct {
			Dashboard map[string]interface{} `json:"dashboard"`
		}

		if err := json.Unmarshal(dashboardOutput, &dashboardData); err == nil {
			cleanDashboard, _ := json.MarshalIndent(dashboardData.Dashboard, "", "  ")
			dashboardFile := fmt.Sprintf("%s/%s.json", dashboardsDir, result.Slug)
			if err := os.WriteFile(dashboardFile, cleanDashboard, 0644); err != nil {
				g.logger.Warn("Failed to write dashboard", "slug", result.Slug, "error", err)
			}
		}
	}

	g.logger.Info("Grafana dashboards backed up", "count", len(searchResults))
	return nil
}

// backupDatasources backs up all Grafana datasources
func (g *GrafanaBackup) backupDatasources(grafanaURL, apiKey, backupDir, timestamp string) error {
	g.logger.Info("Backing up Grafana datasources")

	datasourcesDir := fmt.Sprintf("%s/datasources", backupDir)
	if err := os.MkdirAll(datasourcesDir, 0755); err != nil {
		return err
	}

	// Get all datasources
	datasourcesCmd := exec.Command("curl", "-s",
		fmt.Sprintf("%s/api/datasources", grafanaURL),
		"-H", "Authorization: Bearer "+apiKey)

	output, err := datasourcesCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get datasources: %w", err)
	}

	// Parse datasources
	var datasources []map[string]interface{}
	if err := json.Unmarshal(output, &datasources); err != nil {
		return fmt.Errorf("failed to parse datasources: %w", err)
	}

	// Save each datasource
	for _, ds := range datasources {
		dsJSON, _ := json.MarshalIndent(ds, "", "  ")
		dsName := ds["name"].(string)
		dsFile := fmt.Sprintf("%s/%s.json", datasourcesDir, dsName)
		if err := os.WriteFile(dsFile, dsJSON, 0644); err != nil {
			g.logger.Warn("Failed to write datasource", "name", dsName, "error", err)
		}
	}

	g.logger.Info("Grafana datasources backed up", "count", len(datasources))
	return nil
}

// backupAlertRules backs up all Grafana alert rules
func (g *GrafanaBackup) backupAlertRules(grafanaURL, apiKey, backupDir, timestamp string) error {
	g.logger.Info("Backing up Grafana alert rules")

	alertRulesDir := fmt.Sprintf("%s/alert-rules", backupDir)
	if err := os.MkdirAll(alertRulesDir, 0755); err != nil {
		return err
	}

	// Get all alert rules
	alertRulesCmd := exec.Command("curl", "-s",
		fmt.Sprintf("%s/api/alert-rules", grafanaURL),
		"-H", "Authorization: Bearer "+apiKey)

	output, err := alertRulesCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get alert rules: %w", err)
	}

	// Save alert rules
	alertRulesFile := fmt.Sprintf("%s/alert-rules.json", alertRulesDir)
	if err := os.WriteFile(alertRulesFile, output, 0644); err != nil {
		return fmt.Errorf("failed to write alert rules: %w", err)
	}

	g.logger.Info("Grafana alert rules backed up")
	return nil
}

// backupNotificationChannels backs up all Grafana notification channels
func (g *GrafanaBackup) backupNotificationChannels(grafanaURL, apiKey, backupDir, timestamp string) error {
	g.logger.Info("Backing up Grafana notification channels")

	notificationChannelsDir := fmt.Sprintf("%s/notification-channels", backupDir)
	if err := os.MkdirAll(notificationChannelsDir, 0755); err != nil {
		return err
	}

	// Get all notification channels
	channelsCmd := exec.Command("curl", "-s",
		fmt.Sprintf("%s/api/alert-notifications", grafanaURL),
		"-H", "Authorization: Bearer "+apiKey)

	output, err := channelsCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get notification channels: %w", err)
	}

	// Save notification channels
	channelsFile := fmt.Sprintf("%s/notification-channels.json", notificationChannelsDir)
	if err := os.WriteFile(channelsFile, output, 0644); err != nil {
		return fmt.Errorf("failed to write notification channels: %w", err)
	}

	g.logger.Info("Grafana notification channels backed up")
	return nil
}

// Restore restores Grafana from backup
func (g *GrafanaBackup) Restore(backupFile string) error {
	g.logger.Info("Starting Grafana restore", "file", backupFile)

	// Create temporary extraction directory
	extractDir := fmt.Sprintf("./backups/grafana/restore_%d", time.Now().Unix())
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		return fmt.Errorf("failed to create extraction directory: %w", err)
	}
	defer os.RemoveAll(extractDir)

	// Extract tar.gz
	cmd := exec.Command("tar", "-xzf", backupFile, "-C", extractDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to extract backup: %w", err)
	}

	// Get Grafana configuration
	grafanaURL := getEnvOrDefault("GRAFANA_URL", "http://localhost:3000")
	apiKey := getEnvOrDefault("GRAFANA_API_KEY", "")

	// Restore datasources first
	datasourcesDir := fmt.Sprintf("%s/datasources", extractDir)
	if files, err := os.ReadDir(datasourcesDir); err == nil {
		for _, file := range files {
			datasourceFile := fmt.Sprintf("%s/%s", datasourcesDir, file.Name())
			dsJSON, err := os.ReadFile(datasourceFile)
			if err != nil {
				continue
			}

			// Parse datasource
			var datasource map[string]interface{}
			if err := json.Unmarshal(dsJSON, &datasource); err != nil {
				continue
			}

			// Remove id field to let Grafana assign new id
			delete(datasource, "id")
			delete(datasource, "orgId")

			// Create datasource
			dsJSON, _ = json.Marshal(datasource)
			createCmd := exec.Command("curl", "-X", "POST",
				fmt.Sprintf("%s/api/datasources", grafanaURL),
				"-H", "Authorization: Bearer "+apiKey,
				"-H", "Content-Type: application/json",
				"-d", string(dsJSON))

			if err := createCmd.Run(); err != nil {
				g.logger.Warn("Failed to create datasource", "name", file.Name(), "error", err)
			}
		}
	}

	// Restore dashboards
	dashboardsDir := fmt.Sprintf("%s/dashboards", extractDir)
	if files, err := os.ReadDir(dashboardsDir); err == nil {
		for _, file := range files {
			dashboardFile := fmt.Sprintf("%s/%s", dashboardsDir, file.Name())
			dashboardJSON, err := os.ReadFile(dashboardFile)
			if err != nil {
				continue
			}

			// Wrap in dashboard object
			dashboardWrapper := map[string]interface{}{
				"dashboard": json.RawMessage(dashboardJSON),
				"overwrite": true,
			}

			wrapperJSON, _ := json.Marshal(dashboardWrapper)
			importCmd := exec.Command("curl", "-X", "POST",
				fmt.Sprintf("%s/api/dashboards/import", grafanaURL),
				"-H", "Authorization: Bearer "+apiKey,
				"-H", "Content-Type: application/json",
				"-d", string(wrapperJSON))

			if err := importCmd.Run(); err != nil {
				g.logger.Warn("Failed to import dashboard", "name", file.Name(), "error", err)
			}
		}
	}

	g.logger.Info("Grafana restore completed")
	return nil
}
