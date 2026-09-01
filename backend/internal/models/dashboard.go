package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WidgetType represents the type of dashboard widget
type WidgetType string

const (
	// Infrastructure widgets
	WidgetTypeCPU          WidgetType = "cpu"
	WidgetTypeMemory       WidgetType = "memory"
	WidgetTypeDisk         WidgetType = "disk"
	WidgetTypeNetwork      WidgetType = "network"
	WidgetTypeUptime       WidgetType = "uptime"
	WidgetTypeAvailability WidgetType = "availability"

	// Monitoring widgets
	WidgetTypeAlerts            WidgetType = "alerts"
	WidgetTypeActiveIncidents   WidgetType = "active_incidents"
	WidgetTypeAIRecommendations WidgetType = "ai_recommendations"
	WidgetTypeHealthScore       WidgetType = "health_score"

	// Kubernetes widgets
	WidgetTypeClusterHealth  WidgetType = "cluster_health"
	WidgetTypeK8sNodes       WidgetType = "kubernetes_nodes"
	WidgetTypeK8sPods        WidgetType = "kubernetes_pods"
	WidgetTypeK8sDeployments WidgetType = "kubernetes_deployments"
	WidgetTypeK8sServices    WidgetType = "kubernetes_services"

	// Docker widgets
	WidgetTypeContainers WidgetType = "docker_containers"
	WidgetTypeImages     WidgetType = "docker_images"
	WidgetTypeVolumes    WidgetType = "docker_volumes"
	WidgetTypeDNetworks  WidgetType = "docker_networks"

	// AI widgets
	WidgetTypeAIInsights  WidgetType = "ai_insights"
	WidgetTypeRootCause   WidgetType = "root_cause"
	WidgetTypePredictions WidgetType = "predictions"

	// Logs widgets
	WidgetTypeLiveLogs  WidgetType = "live_logs"
	WidgetTypeErrorLogs WidgetType = "error_logs"
	WidgetTypeLogSearch WidgetType = "log_search"

	// Tracing widgets
	WidgetTypeTraceTimeline WidgetType = "trace_timeline"
	WidgetTypeSlowRequests  WidgetType = "slow_requests"
	WidgetTypeServiceMap    WidgetType = "service_map"

	// APM widgets
	WidgetTypeLatency   WidgetType = "latency"
	WidgetTypeRPS       WidgetType = "rps"
	WidgetTypeErrorRate WidgetType = "error_rate"
	WidgetTypeTopAPIs   WidgetType = "top_apis"
)

// Dashboard represents a user-created dashboard
type Dashboard struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	OrganizationID uuid.UUID      `gorm:"type:uuid;index" json:"organization_id"`
	OwnerID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"owner_id"`
	Name           string         `gorm:"not null" json:"name"`
	Description    string         `gorm:"type:text" json:"description"`
	IsDefault      bool           `gorm:"default:false" json:"is_default"`
	IsShared       bool           `gorm:"default:false" json:"is_shared"`
	SharedWith     []byte         `gorm:"type:jsonb;default:'[]'" json:"shared_with"` // Array of user/team IDs
	Config         []byte         `gorm:"type:jsonb;default:'{}'" json:"config"`      // Dashboard-level config
	TemplateID     *uuid.UUID     `gorm:"type:uuid" json:"template_id,omitempty"`     // Reference to template if applicable
	CreatedBy      uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Widgets []DashboardWidget `gorm:"foreignKey:DashboardID;constraint:OnDelete:CASCADE" json:"widgets,omitempty"`
}

// DashboardWidget represents a single widget in a dashboard
type DashboardWidget struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	DashboardID uuid.UUID      `gorm:"type:uuid;not null;index" json:"dashboard_id"`
	WidgetType  WidgetType     `gorm:"not null" json:"widget_type"`
	Title       string         `gorm:"not null" json:"title"`
	X           int            `gorm:"not null;default:0" json:"x"`           // Grid column position
	Y           int            `gorm:"not null;default:0" json:"y"`           // Grid row position
	Width       int            `gorm:"not null;default:4" json:"width"`       // Width in grid columns
	Height      int            `gorm:"not null;default:3" json:"height"`      // Height in grid rows
	Config      []byte         `gorm:"type:jsonb;default:'{}'" json:"config"` // Widget-specific configuration
	Order       int            `gorm:"default:0" json:"order"`                // Display order within dashboard
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate hook to generate UUIDs
func (d *Dashboard) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

func (w *DashboardWidget) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}

// DashboardConfig represents dashboard-level configuration
type DashboardConfig struct {
	Theme           string `json:"theme"`           // light, dark, auto
	RefreshInterval int    `json:"refreshInterval"` // Auto-refresh in seconds
	TimeRange       string `json:"timeRange"`       // 5m, 15m, 1h, 6h, 24h, 7d
	GridColumns     int    `json:"gridColumns"`     // Number of columns in grid (default: 12)
	RowHeight       int    `json:"rowHeight"`       // Height of each row in pixels
	Margin          int    `json:"margin"`          // Margin between widgets in pixels
	ShowTitle       bool   `json:"showTitle"`       // Show dashboard title
	ShowRefreshBtn  bool   `json:"showRefreshBtn"`  // Show manual refresh button
	ShowTimeRange   bool   `json:"showTimeRange"`   // Show time range selector
	CompactMode     bool   `json:"compactMode"`     // Compact widget padding
}

// DefaultDashboardConfig returns default dashboard configuration
func DefaultDashboardConfig() DashboardConfig {
	return DashboardConfig{
		Theme:           "dark",
		RefreshInterval: 30,
		TimeRange:       "1h",
		GridColumns:     12,
		RowHeight:       80,
		Margin:          8,
		ShowTitle:       true,
		ShowRefreshBtn:  true,
		ShowTimeRange:   true,
		CompactMode:     false,
	}
}

// WidgetConfig represents widget-specific configuration
type WidgetConfig struct {
	// Common config
	Machine   string `json:"machine,omitempty"`   // Target machine/server ID or "all"
	Refresh   int    `json:"refresh,omitempty"`   // Refresh interval in seconds
	Chart     string `json:"chart,omitempty"`     // Chart type: line, area, bar, gauge, number, table
	Period    string `json:"period,omitempty"`    // Time period: 5m, 15m, 1h, 6h, 24h, 7d
	Threshold int    `json:"threshold,omitempty"` // Warning threshold percentage

	// Alerts config
	Severity string `json:"severity,omitempty"` // Filter by severity
	Limit    int    `json:"limit,omitempty"`    // Max items to show

	// Logs config
	LogLevel  string `json:"logLevel,omitempty"`  // Filter log level
	Search    string `json:"search,omitempty"`    // Search query
	TailLines int    `json:"tailLines,omitempty"` // Number of log lines

	// Kubernetes/Docker config
	Namespace string `json:"namespace,omitempty"` // K8s namespace
	Cluster   string `json:"cluster,omitempty"`   // Cluster name

	// Display options
	ShowLegend       bool `json:"showLegend,omitempty"`       // Show chart legend
	ShowAxis         bool `json:"showAxis,omitempty"`         // Show axis labels
	Stacked          bool `json:"stacked,omitempty"`          // Stacked chart
	AnomalyDetection bool `json:"anomalyDetection,omitempty"` // Enable anomaly detection

	// AI-specific
	ShowDetails bool `json:"showDetails,omitempty"` // Show detailed information

	// Service map
	Depth int `json:"depth,omitempty"` // Trace depth for service map
}

// DefaultWidgetConfig returns default widget configuration
func DefaultWidgetConfig() WidgetConfig {
	return WidgetConfig{
		Refresh:    5,
		Chart:      "line",
		Period:     "1h",
		Threshold:  80,
		Limit:      20,
		ShowLegend: true,
		ShowAxis:   true,
	}
}

// DashboardTemplate represents a pre-built dashboard template
type DashboardTemplate struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Category    string           `json:"category"` // infrastructure, kubernetes, docker, security, executive
	Icon        string           `json:"icon"`     // Icon name
	Widgets     []TemplateWidget `json:"widgets"`
	IsDefault   bool             `json:"is_default"`
}

// TemplateWidget represents a widget in a template
type TemplateWidget struct {
	WidgetType WidgetType             `json:"widget_type"`
	Title      string                 `json:"title"`
	X          int                    `json:"x"`
	Y          int                    `json:"y"`
	Width      int                    `json:"width"`
	Height     int                    `json:"height"`
	Config     map[string]interface{} `json:"config"`
}

// DashboardShare represents a shared dashboard entry
type DashboardShare struct {
	ID          uuid.UUID  `json:"id"`
	DashboardID uuid.UUID  `json:"dashboard_id"`
	SharedBy    uuid.UUID  `json:"shared_by"`
	SharedWith  uuid.UUID  `json:"shared_with"` // User or Team ID
	ShareType   string     `json:"share_type"`  // user, team, organization
	Permission  string     `json:"permission"`  // view, edit, admin
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}
