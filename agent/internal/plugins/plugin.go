package plugins

import "time"

// Plugin represents a monitoring plugin that can collect metrics
type Plugin interface {
	// Name returns the unique identifier for this plugin
	Name() string

	// Collect gathers metrics from the plugin's source
	// Returns the collected data as an interface{} and any error encountered
	Collect() (interface{}, error)

	// HealthCheck returns the current health status of the plugin
	HealthCheck() PluginHealth

	// IsEnabled returns whether the plugin is currently enabled
	IsEnabled() bool

	// SetEnabled enables or disables the plugin
	SetEnabled(enabled bool)
}

// PluginHealth represents the health status of a plugin
type PluginHealth struct {
	Plugin   string    `json:"plugin"`
	Status   string    `json:"status"` // "healthy", "degraded", "unhealthy"
	LastRun  time.Time `json:"last_run"`
	Error    string    `json:"error,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// NewPluginHealth creates a new PluginHealth instance
func NewPluginHealth(pluginName string) PluginHealth {
	return PluginHealth{
		Plugin:   pluginName,
		Status:   "healthy",
		LastRun:  time.Now(),
		Metadata: make(map[string]string),
	}
}