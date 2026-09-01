# InfraPilot Enterprise Backend Guide

## Overview

The InfraPilot Backend is a high-performance, scalable API server and WebSocket hub built with Go. It handles authentication, real-time communication, event processing, and business logic.

## Technology Stack

- **Language:** Go 1.24
- **Web Framework:** Gin
- **WebSocket:** Gorilla WebSocket
- **Database Driver:** pgx (PostgreSQL)
- **Cache:** go-redis
- **Queue:** Redis Streams
- **Authentication:** JWT (golang-jwt)
- **Logging:** zap / logrus
- **Metrics:** Prometheus client
- **Tracing:** OpenTelemetry
- **Testing:** testify, golang/mock

## Architecture

### Directory Structure

```
backend/
├── cmd/
│   ├── server/
│   │   └── main.go           # Application entry point
│   └── agent/
│       └── main.go            # Agent entry point
├── internal/
│   ├── app/
│   │   └── server.go          # Server setup and configuration
│   ├── routes/
│   │   └── routes.go          # Route definitions
│   ├── handlers/
│   │   ├── auth/
│   │   │   └── handler.go
│   │   ├── machines/
│   │   │   └── handler.go
│   │   ├── metrics/
│   │   │   └── handler.go
│   │   ├── terminal/
│   │   │   └── handler.go
│   │   └── files/
│   │       └── handler.go
│   ├── models/
│   │   ├── user.go
│   │   ├── machine.go
│   │   ├── session.go
│   │   ├── alert.go
│   │   └── metric.go
│   ├── database/
│   │   ├── database.go        # Connection pool setup
│   │   └── migrations/        # SQL migrations
│   ├── cache/
│   │   ├── session.go         # Session management
│   │   └── metrics.go         # Metrics caching
│   ├── events/
│   │   ├── bus.go             # Event bus
│   │   ├── metric.go
│   │   ├── alert.go
│   │   └── machine.go
│   ├── subscribers/
│   │   ├── database.go        # Persist to PostgreSQL
│   │   ├── websocket.go       # Broadcast to WS clients
│   │   └── alerts.go          # Alert processing
│   ├── workers/
│   │   ├── pool.go
│   │   ├── metrics_worker.go
│   │   └── alert_worker.go
│   ├── queue/
│   │   └── queue.go           # Redis Streams queue
│   ├── scheduler/
│   │   └── scheduler.go       # Cron jobs
│   ├── observability/
│   │   ├── health.go          # Health checks
│   │   ├── metrics.go         # Prometheus metrics
│   │   └── tracing.go         # OpenTelemetry
│   ├── middleware/
│   │   ├── auth.go            # JWT authentication
│   │   ├── rbac.go            # Role-based access
│   │   ├── cors.go
│   │   ├── ratelimit.go
│   │   └── logging.go
│   └── config/
│       └── config.go          # Configuration loading
├── migrations/
│   ├── 001_initial_schema.sql
│   ├── 002_add_alerts.sql
│   └── ...
├── Dockerfile
├── go.mod
└── go.sum
```

### Component Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                        HTTP Server (Gin)                     │
│  ┌─────────────┐  ┌──────────────┐  ┌──────────────────┐  │
│  │   Auth      │  │   Machines   │  │   Terminal       │  │
│  │   Routes    │  │   Routes     │  │   Routes         │  │
│  └─────────────┘  └──────────────┘  └──────────────────┘  │
└──────────────────────────┬──────────────────────────────────┘
                           │
            ┌──────────────┼──────────────┐
            ▼              ▼              ▼
    ┌──────────────┐ ┌──────────┐ ┌──────────────┐
    │  Middleware   │ │ Handlers │ │  WebSocket   │
    │  - Auth       │ │ - Logic  │ │  Hub         │
    │  - RBAC       │ │ - Valid  │ │  - Broadcast │
    │  - Rate Limit │ │ - Process│ │  - Sessions  │
    └──────────────┘ └──────────┘ └──────────────┘
            │              │              │
            └──────────────┼──────────────┘
                           ▼
                  ┌──────────────────┐
                  │   Event Bus      │
                  │  (Redis Streams) │
                  └────────┬─────────┘
                           │
            ┌──────────────┼──────────────┐
            ▼              ▼              ▼
    ┌──────────────┐ ┌──────────┐ ┌──────────────┐
    │   Database   │ │  Cache   │ │   Workers    │
    │  Subscriber  │ │(Sessions)│ │  - Metrics   │
    │              │ │  - State │ │  - Alerts    │
    └──────────────┘ └──────────┘ │  - Reports   │
                                   └──────────────┘
```

## Core Components

### Server Setup

**File:** `internal/app/server.go`

```go
package app

import (
    "log"
    "os"
    "github.com/gin-gonic/gin"
    "github.com/infrapilot/backend/internal/routes"
    "github.com/infrapilot/backend/internal/database"
    "github.com/infrapilot/backend/internal/cache"
    "github.com/infrapilot/backend/internal/events"
    "github.com/infrapilot/backend/internal/workers"
)

func NewServer(cfg *config.Config) *gin.Engine {
    // Initialize database
    db, err := database.NewConnection(cfg.DatabaseURL)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Initialize Redis
    redisClient := cache.NewRedisClient(cfg.RedisURL)

    // Initialize event bus
    eventBus := events.NewBus(redisClient)

    // Initialize worker pool
    workerPool := workers.NewPool(cfg.WorkerCount, eventBus)

    // Setup routes
    r := gin.Default()
    routes.SetupRoutes(r, db, redisClient, eventBus, workerPool)

    // Start workers
    workerPool.Start()

    return r
}
```

### Event Bus

**File:** `internal/events/bus.go`

```go
package events

type Bus interface {
    Publish(event Event) error
    Subscribe(eventType string, subscriber Subscriber) error
    Start()
    Stop()
}

type Event struct {
    ID        string
    Type      string
    Payload   interface{}
    Timestamp time.Time
}

type Subscriber interface {
    Handle(event Event) error
    Name() string
}
```

### Workers

**File:** `internal/workers/pool.go`

```go
package workers

type Pool struct {
    workers []Worker
    queue   chan Event
    wg      sync.WaitGroup
}

func NewPool(count int, eventBus events.Bus) *Pool {
    pool := &Pool{
        workers: make([]Worker, count),
        queue:   make(chan Event, 1000),
    }

    for i := 0; i < count; i++ {
        pool.workers[i] = NewMetricsWorker(eventBus)
    }

    return pool
}

func (p *Pool) Start() {
    for _, worker := range p.workers {
        p.wg.Add(1)
        go func(w Worker) {
            defer p.wg.Done()
            w.Start()
        }(worker)
    }
}
```

### Middleware

**File:** `internal/middleware/auth.go`

```go
package middleware

import (
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secret []byte) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
            c.Abort()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return secret, nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }

        claims := token.Claims.(jwt.MapClaims)
        c.Set("user_id", claims["user_id"])
        c.Set("role", claims["role"])
        c.Next()
    }
}
```

### Database Layer

**File:** `internal/database/database.go`

```go
package database

import (
    "context"
    "time"
    "github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
    Pool *pgxpool.Pool
}

func NewConnection(databaseURL string) (*DB, error) {
    pool, err := pgxpool.New(context.Background(), databaseURL)
    if err != nil {
        return nil, err
    }

    // Verify connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := pool.Ping(ctx); err != nil {
        return nil, err
    }

    return &DB{Pool: pool}, nil
}
```

## API Design

### Route Structure

```go
// internal/routes/routes.go
func SetupRoutes(r *gin.Engine, db *database.DB, redis *redis.Client, bus *events.Bus, pool *workers.Pool) {
    // Health check
    r.GET("/healthz", health.Check)

    // Auth routes (no auth required)
    auth := r.Group("/api/v1/auth")
    {
        auth.POST("/login", authHandler.Login)
        auth.POST("/refresh", authHandler.Refresh)
        auth.POST("/logout", authHandler.Logout)
    }

    // Protected routes
    api := r.Group("/api/v1")
    api.Use(middleware.AuthMiddleware())
    {
        // Machines
        machines := api.Group("/machines")
        {
            machines.GET("", machineHandler.List)
            machines.GET("/:id", machineHandler.Get)
            machines.DELETE("/:id", machineHandler.Delete)
            machines.GET("/:id/metrics", metricsHandler.Get)
            machines.GET("/:id/docker/containers", dockerHandler.ListContainers)
            machines.GET("/:id/kubernetes/clusters", k8sHandler.ListClusters)
        }

        // Terminal
        api.POST("/terminal/sessions", terminalHandler.CreateSession)
        api.GET("/ws/terminal/:session_id", terminalHandler.HandleWebSocket)

        // Files
        api.GET("/machines/:id/files", fileHandler.List)
        api.GET("/machines/:id/files/content", fileHandler.GetContent)

        // Logs
        api.GET("/machines/:id/logs", logHandler.Get)

        // Alerts
        alerts := api.Group("/alerts")
        {
            alerts.GET("", alertHandler.List)
            alerts.POST("/:id/acknowledge", alertHandler.Acknowledge)
            alerts.POST("/:id/resolve", alertHandler.Resolve)
        }

        // Reports
        reports := api.Group("/reports")
        {
            reports.GET("", reportHandler.List)
            reports.POST("/generate", reportHandler.Generate)
            reports.GET("/:id/download", reportHandler.Download)
        }
    }

    // Metrics endpoint (Prometheus)
    r.GET("/metrics", metricsHandler.Prometheus)
}
```

## Configuration

**File:** `internal/config/config.go`

```go
package config

type Config struct {
    DatabaseURL      string
    RedisURL         string
    JWTSecret        string
    ServerPort       int
    MetricsPort      int
    
    // Worker configuration
    WorkerCount      int
    QueueSize        int
    
    // Observability
    PrometheusEnabled bool
    JaegerEndpoint    string
}

func Load() (*Config, error) {
    cfg := &Config{}
    
    cfg.DatabaseURL = os.Getenv("DATABASE_URL")
    cfg.RedisURL = os.Getenv("REDIS_URL")
    cfg.JWTSecret = os.Getenv("JWT_SECRET")
    cfg.ServerPort = parseIntEnv("SERVER_PORT", 8080)
    cfg.WorkerCount = parseIntEnv("WORKER_COUNT", 4)
    
    return cfg, nil
}
```

## Observability

### Prometheus Metrics

```go
package observability

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    RequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total HTTP requests",
        },
        []string{"method", "path", "status"},
    )

    RequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path"},
    )

    WebSocketConnections = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "websocket_connections",
            Help: "Active WebSocket connections",
        },
    )
)

func init() {
    prometheus.MustRegister(RequestsTotal)
    prometheus.MustRegister(RequestDuration)
    prometheus.MustRegister(WebSocketConnections)
}
```

### Health Checks

```go
package health

func Check(c *gin.Context) {
    checks := map[string]string{
        "status": "healthy",
        "checks": map[string]string{
            "database": checkDatabase(),
            "redis":    checkRedis(),
            "queue":    checkQueue(),
        },
    }
    
    c.JSON(http.StatusOK, checks)
}
```

## Testing

### Unit Tests

```go
// internal/handlers/auth/handler_test.go
func TestLogin(t *testing.T) {
    // Setup
    db := setupTestDB()
    handler := NewAuthHandler(db)

    // Test
    c, _ := gin.CreateTestContext(httptest.NewRecorder())
    c.Request = httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(
        `{"email":"test@example.com","password":"password"}`))
    c.Request.Header.Set("Content-Type", "application/json")

    handler.Login(c)

    // Assert
    assert.Equal(t, http.StatusOK, c.Writer.Status())
}
```

### Integration Tests

```go
// Tests/integration_test.go
func TestAPIEndpoints(t *testing.T) {
    // Start test server
    go main.Run()
    defer main.Stop()

    baseURL := "http://localhost:8080/api/v1"
    
    // Test login
    token := getAuthToken(t, baseURL)
    assert.NotEmpty(t, token)
    
    // Test machines list
    machines := getMachines(t, baseURL, token)
    assert.NotNil(t, machines)
}
```

## Performance Optimization

### Connection Pooling

```go
// Database connection pool
poolConfig, _ := pgxpool.ParseConfig(databaseURL)
poolConfig.MaxConnections = 50
poolConfig.MinConnections = 10
poolConfig.MaxConnLifetime = 1 * time.Hour
poolConfig.MaxConnIdleTime = 30 * time.Minute
```

### Caching Strategy

```go
// Cache frequently accessed data
func GetMachineCached(ctx context.Context, machineID string) (*Machine, error) {
    cacheKey := fmt.Sprintf("machine:%s", machineID)
    
    // Try cache first
    if cached, found := redis.Get(ctx, cacheKey); found {
        return cached.(*Machine), nil
    }
    
    // Fetch from database
    machine, err := fetchMachineFromDB(ctx, machineID)
    if err != nil {
        return nil, err
    }
    
    // Cache for 5 minutes
    redis.Set(ctx, cacheKey, machine, 5*time.Minute)
    
    return machine, nil
}
```

### Request Batching

```go
// Batch metrics for better throughput
func BatchMetrics(metrics []Metric) error {
    if len(metrics) == 0 {
        return nil
    }
    
    // Use COPY for bulk insert
    _, err := db.Pool.CopyFrom(
        context.Background(),
        pgx.Identifier{"metrics"},
        []string{"machine_id", "timestamp", "value"},
        pgx.CopyFromSlice(len(metrics), func(i int) ([]interface{}, error) {
            return []interface{}{
                metrics[i].MachineID,
                metrics[i].Timestamp,
                metrics[i].Value,
            }, nil
        }),
    )
    
    return err
}
```

## Deployment

### Environment Variables

```env
# Database
DATABASE_URL=postgres://user:pass@postgres:5432/infrapilot?sslmode=require
DATABASE_POOL_SIZE=20

# Redis
REDIS_URL=redis://redis:6379/0

# JWT
JWT_SECRET=your-256-bit-secret-key-here
JWT_EXPIRY=1h
REFRESH_TOKEN_EXPIRY=7d

# Server
SERVER_PORT=8080
METRICS_PORT=9090
WORKER_COUNT=4

# CORS
CORS_ORIGIN=https://infrapilot.io

# Observability
PROMETHEUS_ENABLED=true
OTEL_EXPORTER_OTLP_ENDPOINT=jaeger:4317
```

### Docker Build

```bash
# Multi-stage build
docker build -t infrapilot-backend:latest -f Dockerfile.backend .
```

### Process Management

```bash
# Systemd service
# /etc/systemd/system/infrapilot-backend.service
[Unit]
Description=InfraPilot Backend
After=network.target postgresql.service redis.service

[Service]
Type=simple
User=infrapilot
WorkingDirectory=/opt/infrapilot/backend
ExecStart=/opt/infrapilot/backend/infrapilot
Restart=always
RestartSec=10
EnvironmentFile=/etc/infrapilot/.env

[Install]
WantedBy=multi-user.target
```

## Monitoring

### Metrics Endpoint

```
GET /metrics
```

Key metrics exposed:
- `http_requests_total` - Request counter by method/path/status
- `http_request_duration_seconds` - Request latency histogram
- `websocket_connections` - Active WebSocket connections
- `db_connection_pool` - Database connection pool stats
- `redis_hits_total` - Cache hit counter
- `worker_queue_size` - Current queue length
- `events_processed_total` - Event processing rate

### Logging

Structured JSON logging with correlation IDs:

```json
{
  "timestamp": "2025-01-01T12:00:00Z",
  "level": "info",
  "service": "backend",
  "trace_id": "abc123",
  "message": "Request processed",
  "method": "GET",
  "path": "/api/v1/machines",
  "status": 200,
  "duration_ms": 45
}
```

## Security

### Authentication

- JWT with RS256 asymmetric keys
- Refresh tokens with rotation
- Password hashing with bcrypt (cost 12)

### Authorization

- Role-Based Access Control (RBAC)
- Roles: admin, operator, viewer
- Resource-level permissions

### Input Validation

- Request validation with go-playground/validator
- SQL injection prevention with parameterized queries
- XSS prevention with proper escaping

### Rate Limiting

```go
rateLimiter := tollbooth.NewLimiter(
    time.Minute, 
    100, // 100 requests per minute
)
```

## Best Practices

1. **Error Handling:** Always return proper HTTP status codes
2. **Logging:** Log all errors with context
3. **Metrics:** Instrument all critical paths
4. **Testing:** Maintain >80% test coverage
5. **Documentation:** Update API docs for every change
6. **Graceful Shutdown:** Handle SIGTERM properly
7. **Database:** Use transactions for multi-step operations
8. **Security:** Validate all inputs, never trust client data