package app

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	"infrapilot/agent/internal/cache"
	"infrapilot/agent/internal/client"
	"infrapilot/agent/internal/collector"
	logscollector "infrapilot/agent/internal/collector/logs"
	"infrapilot/agent/internal/commands"
	"infrapilot/agent/internal/config"
	"infrapilot/agent/internal/heartbeat"
	"infrapilot/agent/internal/reconnect"
	"infrapilot/agent/internal/register"
	"infrapilot/agent/internal/sender"
	"infrapilot/agent/internal/state"
)

func RunAgent() {
	log.Println("Starting InfraPilot Monitoring Agent...")

	cache.Init("offline_queue.json")

	store := config.NewConfigStore(config.DefaultConfigPath())

	var cfg *config.EnterpriseConfig

	if store.Exists() {
		// config.json exists, attempt to load it
		var loadErr error
		cfg, loadErr = store.Load()
		if loadErr != nil {
			log.Printf("Warning: Failed to parse existing config.json: %v. Will re-enroll.", loadErr)
			cfg = nil
		} else if cfg.MetricsInterval == 0 {
			cfg.MetricsInterval = 5
		}
		if cfg != nil && cfg.HeartbeatInterval == 0 {
			cfg.HeartbeatInterval = 15
		}
	}

	if cfg == nil || cfg.MachineID == "" || cfg.APIKey == "" {
		// config.json does not exist or lacks credentials, run agent self-registration/enrollment
		log.Println("Agent credentials missing or incomplete. Enrolling machine with backend...")

		backendURL := "http://localhost:8080"
		if cfg != nil && cfg.Server != "" {
			backendURL = cfg.Server
		}
		enrollmentToken := ""
		// Check config.yaml to get the server URL if it exists
		if yamlCfg, yamlErr := config.LoadConfig("config.yaml"); yamlErr == nil {
			if yamlCfg.BackendURL != "" {
				backendURL = yamlCfg.BackendURL
			}
			enrollmentToken = yamlCfg.EnrollmentToken
		}

		metrics, err := collector.GetMetrics()
		if err != nil {
			log.Fatalf("Failed to fetch initial system metrics for enrollment: %v", err)
		}

		// Sprint 6.5 Enrollment Request Payload
		regPayload := map[string]interface{}{
			"enrollment_token": enrollmentToken,
			"hostname":         metrics.Hostname,
			"os":               metrics.OS,
			"platform":         metrics.Platform,
			"ip_address":       metrics.IPAddress,
			"kernel":           metrics.Kernel,
			"architecture":     runtime.GOARCH,
			"cpu_model":        metrics.CPUModel,
			"cpu_cores":        metrics.CPUCores,
			"total_memory_gb":  metrics.MemoryTotal / 1024 / 1024 / 1024,
			"total_disk_gb":    metrics.DiskTotal / 1024 / 1024 / 1024,
			"memory":           metrics.MemoryTotal / 1024 / 1024 / 1024,
		}

		// Perform enrollment call using the isolated register package
		result, regErr := register.RegisterAgent(backendURL, regPayload, enrollmentToken)
		if regErr != nil {
			log.Fatalf("Enrollment failed: %v", regErr)
		}

		cfg = &config.EnterpriseConfig{
			Server:            backendURL,
			MachineID:         result["machine_id"],
			APIKey:            result["api_key"],
			Organization:      result["organization_id"],
			MetricsInterval:   5,
			HeartbeatInterval: 15,
		}

		if err := store.Save(cfg); err != nil {
			log.Fatalf("Failed to save config.json: %v", err)
		}

		log.Println("Enrollment successful. Machine ID:", cfg.MachineID)
		log.Println("Configuration saved in config.json")
	}

	// Initialize AppConfig in config package
	config.AppConfig = &config.Config{
		BackendURL: cfg.Server,
		MachineID:  cfg.MachineID,
		APIKey:     cfg.APIKey,
		KeyVersion: cfg.KeyVersion,
		Interval:   cfg.MetricsInterval,
	}

	// Initialize the shared client
	if err := client.InitHTTPClient(); err != nil {
		log.Printf("Failed to initialize HTTP client: %v", err)
	}

	log.Printf("Agent initialized. Backend: %s, Machine ID: %s, Org: %s", cfg.Server, cfg.MachineID, cfg.Organization)

	// API key rotation poller
	go commands.PollKeyRotation(cfg.Server, cfg.MachineID, cfg.APIKey, store)

	// Heartbeat loop with reconnection
	go func() {
		reconnect.ExecuteWithBackoff(nil, reconnect.DefaultRetryStrategy(), "Heartbeat")
		for {
			if state.Connected {
				fmt.Println("Connected")
			} else {
				fmt.Println("Disconnected")
			}

			err := heartbeat.Send(cfg.Server, cfg.MachineID, cfg.APIKey)
			if err != nil {
				log.Printf("Heartbeat Sender Error: %v", err)
			}
			time.Sleep(time.Duration(cfg.HeartbeatInterval) * time.Second)
		}
	}()

	// Poll and execute remote commands in a background thread
	go commands.PollAndExecute(cfg.Server, cfg.MachineID, cfg.APIKey)

	// Initialize Agent v2 Modular Collector Manager
	ctx := context.Background()
	mgr := collector.NewManager(5 * time.Second)
	mgr.Register(logscollector.NewCollector(cfg.MachineID, cfg.Server))
	go mgr.Start(ctx)

	// Metrics collection loop with reconnection
	reconnect.ExecuteWithBackoff(nil, reconnect.DefaultRetryStrategy(), "Metrics")
	for {
		if state.Connected {
			fmt.Println("Connected")
		} else {
			fmt.Println("Disconnected")
		}

		metrics, err := collector.GetMetrics()
		if err != nil {
			log.Printf("Metrics Collection Error: %v", err)
			time.Sleep(time.Duration(cfg.MetricsInterval) * time.Second)
			continue
		}

		err = sender.SendMetrics(cfg.Server, cfg.MachineID, cfg.APIKey, metrics)
		if err != nil {
			log.Printf("Sender Error: %v", err)
		} else {
			log.Println("Metrics successfully delivered.")
		}

		time.Sleep(time.Duration(cfg.MetricsInterval) * time.Second)
	}
}
