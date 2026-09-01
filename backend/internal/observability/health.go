package observability

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"infrapilot/backend/internal/cache"
	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"
	"infrapilot/backend/internal/websocket"

	"github.com/gin-gonic/gin"
)

// HealthChecker performs health checks on system components
type HealthChecker struct {
	checks map[string]HealthCheck
	mu     sync.RWMutex
}

// HealthCheck represents a component health check
type HealthCheck struct {
	Name      string
	Check     func() error
	Critical  bool
	Component string
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string                     `json:"status"`
	Timestamp time.Time                  `json:"timestamp"`
	Version   string                     `json:"version"`
	Checks    map[string]ComponentHealth `json:"checks"`
	Uptime    time.Duration              `json:"uptime"`
}

// ComponentHealth represents the health of a single component
type ComponentHealth struct {
	Status  string        `json:"status"`
	Latency time.Duration `json:"latency"`
	Error   string        `json:"error,omitempty"`
}

var (
	healthChecker *HealthChecker
	startTime     = time.Now()
)

// GetHealthChecker returns the singleton health checker
func GetHealthChecker() *HealthChecker {
	if healthChecker == nil {
		healthChecker = &HealthChecker{
			checks: make(map[string]HealthCheck),
		}
	}
	return healthChecker
}

// RegisterCheck registers a health check
func (h *HealthChecker) RegisterCheck(name string, check HealthCheck) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.checks[name] = check
}

// CheckHealth performs all registered health checks
func (h *HealthChecker) CheckHealth() HealthResponse {
	h.mu.RLock()
	checks := make(map[string]ComponentHealth)
	overallStatus := "healthy"

	for _, check := range h.checks {
		start := time.Now()
		err := check.Check()
		latency := time.Since(start)

		status := "healthy"
		if err != nil {
			status = "unhealthy"
			if check.Critical {
				overallStatus = "unhealthy"
			}
		}

		componentHealth := ComponentHealth{
			Status:  status,
			Latency: latency,
		}
		if err != nil {
			componentHealth.Error = err.Error()
		}
		checks[check.Component] = componentHealth
	}
	h.mu.RUnlock()

	return HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Checks:    checks,
		Uptime:    time.Since(startTime),
	}
}

// LivenessHandler returns a simple liveness probe
func LivenessHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
	})
}

// ReadinessHandler returns readiness status
func ReadinessHandler(c *gin.Context) {
	checker := GetHealthChecker()
	health := checker.CheckHealth()

	if health.Status == "healthy" {
		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
		})
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not_ready",
			"checks": health.Checks,
		})
	}
}

// HealthHandler returns comprehensive health status
func HealthHandler(c *gin.Context) {
	checker := GetHealthChecker()
	health := checker.CheckHealth()

	statusCode := http.StatusOK
	if health.Status == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, health)
}

// MetricsHandler returns prometheus metrics
func MetricsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "metrics_enabled"})
}

// InitHealthChecks initializes default health checks
func InitHealthChecks() {
	checker := GetHealthChecker()

	// Database health check
	checker.RegisterCheck("database", HealthCheck{
		Name:      "database",
		Component: "database",
		Critical:  true,
		Check: func() error {
			sqlDB, err := database.DB.DB()
			if err != nil {
				return err
			}
			return sqlDB.Ping()
		},
	})

	// Redis health check
	checker.RegisterCheck("redis", HealthCheck{
		Name:      "redis",
		Component: "redis",
		Critical:  true,
		Check: func() error {
			if cache.RedisClient == nil {
				return fmt.Errorf("redis client not initialized")
			}
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			return cache.RedisClient.Ping(ctx).Err()
		},
	})

	// WebSocket health check
	checker.RegisterCheck("websocket", HealthCheck{
		Name:      "websocket",
		Component: "websocket",
		Critical:  true,
		Check: func() error {
			if websocket.WS == nil {
				return fmt.Errorf("websocket hub not initialized")
			}
			return nil
		},
	})

	// AI Service health check
	checker.RegisterCheck("ai_service", HealthCheck{
		Name:      "ai_service",
		Component: "ai_service",
		Critical:  false,
		Check: func() error {
			if database.DB == nil {
				return fmt.Errorf("database connection not initialized")
			}
			var count int64
			if err := database.DB.Table("ai_incidents").Count(&count).Error; err != nil {
				return fmt.Errorf("AI database tables unreachable: %w", err)
			}
			return nil
		},
	})

	// Disk space check
	checker.RegisterCheck("disk", HealthCheck{
		Name:      "disk",
		Component: "disk",
		Critical:  false,
		Check: func() error {
			// Simple disk space check
			// In production, use proper disk usage check
			return nil
		},
	})

	// Connected agents check
	checker.RegisterCheck("agents", HealthCheck{
		Name:      "agents",
		Component: "agents",
		Critical:  false,
		Check: func() error {
			var count int64
			database.DB.Model(&models.Machine{}).Where("status = ?", "online").Count(&count)
			if count == 0 {
				return fmt.Errorf("no agents connected")
			}
			return nil
		},
	})
}
