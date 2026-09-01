package plugins

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Manager manages all monitoring plugins
type Manager struct {
	plugins map[string]Plugin
	mu      sync.RWMutex
}

// NewManager creates a new plugin manager
func NewManager() *Manager {
	return &Manager{
		plugins: make(map[string]Plugin),
	}
}

// Register adds a plugin to the manager
func (m *Manager) Register(plugin Plugin) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.plugins[plugin.Name()] = plugin
}

// Unregister removes a plugin from the manager
func (m *Manager) Unregister(pluginName string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.plugins, pluginName)
}

// GetPlugin retrieves a plugin by name
func (m *Manager) GetPlugin(pluginName string) (Plugin, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	plugin, exists := m.plugins[pluginName]
	return plugin, exists
}

// GetAllPlugins returns all registered plugins
func (m *Manager) GetAllPlugins() []Plugin {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugins := make([]Plugin, 0, len(m.plugins))
	for _, plugin := range m.plugins {
		plugins = append(plugins, plugin)
	}
	return plugins
}

// GetEnabledPlugins returns only enabled plugins
func (m *Manager) GetEnabledPlugins() []Plugin {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugins := make([]Plugin, 0, len(m.plugins))
	for _, plugin := range m.plugins {
		if plugin.IsEnabled() {
			plugins = append(plugins, plugin)
		}
	}
	return plugins
}

// CollectAll gathers metrics from all enabled plugins
func (m *Manager) CollectAll() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	results := make(map[string]interface{})

	for name, plugin := range m.plugins {
		if !plugin.IsEnabled() {
			continue
		}

		data, err := plugin.Collect()
		if err != nil {
			// Update health status on error
			health := plugin.HealthCheck()
			health.Status = "unhealthy"
			health.Error = err.Error()
			continue
		}

		results[name] = data
	}

	return results
}

// CollectPlugin gathers metrics from a specific plugin
func (m *Manager) CollectPlugin(pluginName string) (interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugin, exists := m.plugins[pluginName]
	if !exists {
		return nil, fmt.Errorf("plugin not found: %s", pluginName)
	}

	if !plugin.IsEnabled() {
		return nil, fmt.Errorf("plugin is disabled: %s", pluginName)
	}

	return plugin.Collect()
}

// HealthCheckAll returns health status for all plugins
func (m *Manager) HealthCheckAll() []PluginHealth {
	m.mu.RLock()
	defer m.mu.RUnlock()

	healthList := make([]PluginHealth, 0, len(m.plugins))
	for _, plugin := range m.plugins {
		healthList = append(healthList, plugin.HealthCheck())
	}
	return healthList
}

// PluginCount returns the total number of registered plugins
func (m *Manager) PluginCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.plugins)
}

// EnabledPluginCount returns the number of enabled plugins
func (m *Manager) EnabledPluginCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := 0
	for _, plugin := range m.plugins {
		if plugin.IsEnabled() {
			count++
		}
	}
	return count
}

// ToJSON serializes all collected metrics to JSON
func (m *Manager) ToJSON() ([]byte, error) {
	results := m.CollectAll()
	return json.Marshal(results)
}

// MetricsSummary represents a summary of all plugin metrics
type MetricsSummary struct {
	Timestamp time.Time              `json:"timestamp"`
	Plugins   int                    `json:"plugins_count"`
	Enabled   int                    `json:"enabled_count"`
	Health    []PluginHealth         `json:"health"`
	Data      map[string]interface{} `json:"data"`
}

// GetSummary returns a comprehensive summary of all metrics
func (m *Manager) GetSummary() MetricsSummary {
	return MetricsSummary{
		Timestamp: time.Now(),
		Plugins:   m.PluginCount(),
		Enabled:   m.EnabledPluginCount(),
		Health:    m.HealthCheckAll(),
		Data:      m.CollectAll(),
	}
}
