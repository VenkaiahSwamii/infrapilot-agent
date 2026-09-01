package cli

import (
	"fmt"
	"os"
	"runtime"

	"infrapilot/agent/internal/collector"
	"infrapilot/agent/internal/config"
	"infrapilot/agent/internal/register"
	"infrapilot/agent/internal/system"

	"github.com/spf13/cobra"
)

var (
	enrollServer string
	enrollToken  string
	enrollForce  bool
)

var enrollCmd = &cobra.Command{
	Use:   "enroll",
	Short: "Enroll this machine with an InfraPilot Enterprise backend",
	RunE: func(cmd *cobra.Command, args []string) error {

		if enrollServer == "" {
			return fmt.Errorf("--server is required (e.g. --server http://192.168.1.7:8080)")
		}
		if enrollToken == "" {
			return fmt.Errorf("--token is required (enrollment token from the InfraPilot dashboard)")
		}

		configPath := "config.json"
		if _, err := os.Stat(configPath); err == nil && !enrollForce {
			return fmt.Errorf("config.json already exists in this directory (already enrolled). Use --force to re-enroll and overwrite it")
		}

		fmt.Println("Enrolling with backend:", enrollServer)

		metrics, err := collector.GetMetrics()
		if err != nil {
			return fmt.Errorf("failed to collect system metrics for enrollment: %w", err)
		}

		regPayload := map[string]interface{}{
			"hostname":        metrics.Hostname,
			"os":              metrics.OS,
			"platform":        metrics.Platform,
			"ip":              metrics.IPAddress,
			"agent_version":   system.GetVersion(),
			"architecture":    runtime.GOARCH,
			"resource_type":   metrics.OS,
			"total_memory_gb": metrics.MemoryTotal / 1024 / 1024 / 1024,
			"total_disk_gb":   metrics.DiskTotal / 1024 / 1024 / 1024,
		}

		result, err := register.RegisterAgent(enrollServer, regPayload, enrollToken)
		if err != nil {
			return fmt.Errorf("enrollment failed: %w", err)
		}

		cfg := &config.Config{
			MachineID:  result["machine_id"],
			APIKey:     result["api_key"],
			KeyVersion: 1,
			BackendURL: enrollServer,
			Interval:   5,
		}

		if err := config.SaveConfigJSON(configPath, cfg); err != nil {
			return fmt.Errorf("failed to save config.json: %w", err)
		}

		fmt.Println("Enrollment successful. Machine ID:", cfg.MachineID)
		fmt.Println("Backend saved:", cfg.BackendURL)
		fmt.Println("Run 'agent.exe start' to begin monitoring.")

		return nil
	},
}

func init() {
	enrollCmd.Flags().StringVar(&enrollServer, "server", "", "Backend server URL, e.g. http://192.168.1.7:8080")
	enrollCmd.Flags().StringVar(&enrollToken, "token", "", "Enrollment token from the InfraPilot dashboard")
	enrollCmd.Flags().BoolVar(&enrollForce, "force", false, "Overwrite existing config.json if already enrolled")
	rootCmd.AddCommand(enrollCmd)
}
