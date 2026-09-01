package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"infrapilot/agent/internal/metrics"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// DockerPlugin collects Docker container metrics
type DockerPlugin struct {
	enabled bool
	mu      sync.RWMutex
	health  PluginHealth
	client  *client.Client
}

// NewDockerPlugin creates a new Docker plugin instance
func NewDockerPlugin() *DockerPlugin {
	plugin := &DockerPlugin{
		enabled: true,
		health:  NewPluginHealth("docker"),
	}

	// Initialize Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		plugin.health.Status = "unhealthy"
		plugin.health.Error = "Failed to create Docker client: " + err.Error()
	} else {
		plugin.client = cli
		plugin.health.Status = "healthy"
	}

	return plugin
}

// Name returns the plugin name
func (p *DockerPlugin) Name() string {
	return "docker"
}

// IsEnabled returns whether the plugin is enabled
func (p *DockerPlugin) IsEnabled() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.enabled
}

// SetEnabled enables or disables the plugin
func (p *DockerPlugin) SetEnabled(enabled bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.enabled = enabled
}

// HealthCheck returns the health status of the plugin
func (p *DockerPlugin) HealthCheck() PluginHealth {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.health
}

// Collect gathers Docker container metrics
func (p *DockerPlugin) Collect() (interface{}, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	startTime := time.Now()

	if p.client == nil {
		p.health.Status = "unhealthy"
		p.health.Error = "Docker client not initialized"
		return []metrics.DockerContainerMetric{}, nil
	}

	result := make([]metrics.DockerContainerMetric, 0)

	containers, err := p.client.ContainerList(
		context.Background(),
		container.ListOptions{All: true},
	)

	if err != nil {
		p.health.Status = "degraded"
		p.health.Error = "Failed to list containers: " + err.Error()
		return []metrics.DockerContainerMetric{}, nil
	}

	for _, c := range containers {
		name := ""
		if len(c.Names) > 0 {
			name = c.Names[0]
		}

		cpuPercent := 0.0
		memoryUsed := uint64(0)
		memoryLimit := uint64(0)

		// Get container stats
		stats, err := p.client.ContainerStats(context.Background(), c.ID, false)
		if err == nil {
			body, readErr := io.ReadAll(stats.Body)
			stats.Body.Close()

			if readErr == nil {
				var response struct {
					MemoryStats struct {
						Usage uint64 `json:"usage"`
						Limit uint64 `json:"limit"`
					} `json:"memory_stats"`

					CPUStats struct {
						CPUUsage struct {
							TotalUsage uint64 `json:"total_usage"`
						} `json:"cpu_usage"`

						SystemUsage uint64 `json:"system_cpu_usage"`
					} `json:"cpu_stats"`
				}

				json.Unmarshal(body, &response)

				memoryUsed = response.MemoryStats.Usage
				memoryLimit = response.MemoryStats.Limit

				if memoryLimit > 0 {
					cpuPercent = float64(response.CPUStats.CPUUsage.TotalUsage) /
						float64(response.CPUStats.SystemUsage+1) * 100
				}
			}
		}

		result = append(result, metrics.DockerContainerMetric{
			ID:           c.ID[:12],
			Name:         name,
			Image:        c.Image,
			Status:       c.Status,
			State:        c.State,
			CPUPercent:   cpuPercent,
			MemoryUsed:   memoryUsed,
			MemoryLimit:  memoryLimit,
			RestartCount: 0,
		})
	}

	p.health.Status = "healthy"
	p.health.Error = ""
	p.health.LastRun = time.Now()
	p.health.Metadata["container_count"] = string(rune(len(result)))
	p.health.Metadata["duration_ms"] = time.Since(startTime).String()

	return result, nil
}

// Ping checks if the Docker daemon is accessible
func (p *DockerPlugin) Ping() error {
	if p.client == nil {
		return fmt.Errorf("Docker client not initialized")
	}
	_, err := p.client.Ping(context.Background())
	return err
}
