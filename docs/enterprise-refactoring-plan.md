# InfraPilot Enterprise v1.0 - Architecture Review & Refactoring Plan

## Executive Summary

**Current State:** Functional MVP with clean folder structure but missing Clean Architecture boundaries, proper DDD layers, and enterprise-grade separation of concerns.

**Target State:** Production-ready enterprise platform with Clean Architecture, proper DDD, scalable worker pools, and AI-ready foundations.

**Estimated Effort:** 3-4 weeks for Sprint 1 (Clean Architecture) + ongoing feature development

---

## 1. Current Architecture Assessment

### ✅ What's Working Well

1. **Folder Structure** - Already follows Go best practices with `cmd/`, `internal/`, `pkg/` separation
2. **Technology Stack** - Modern choices: Go, Gin, GORM, PostgreSQL, WebSockets, React
3. **Multi-tenant Foundation** - Organization-scoped data model exists
4. **RBAC** - Role-based access control implemented
5. **WebSocket Hub** - Real-time updates infrastructure in place
6. **Agent Architecture** - Lightweight Go agent with collectors and sender

### ❌ Critical Issues Found

#### Backend Issues

1. **No Clean Architecture Boundaries**
   - Models directly import GORM (infrastructure coupling)
   - Handlers directly depend on Gin framework
   - No domain layer with pure business entities
   - No repository interfaces

2. **Monolithic Route Registration**
   - `routes.go` is 288 lines with everything inline
   - No route grouping by feature/domain
   - Duplicate route definitions (e.g., `/heartbeat` appears twice)

3. **Missing DDD Layers**
   - No `domain/` layer with entities and interfaces
   - No `dto/` or `request/` `response/` objects
   - Business logic mixed in handlers

4. **No Worker Pool Implementation**
   - Architecture doc mentions workers but no `cmd/worker/` implementation
   - No durable ingestion queue
   - No rollup generation

5. **Tight Coupling**
   - Handlers call `database.DB` directly (violates DIP)
   - No service interfaces for testing
   - WebSocket hub tightly coupled to Gin

6. **Missing Enterprise Features**
   - No graceful shutdown
   - No health checks
   - No structured logging
   - No distributed tracing
   - No metrics/monitoring

#### Frontend Issues

1. **No TypeScript** - Using plain JavaScript despite @types in devDependencies
2. **Missing State Management** - No Redux/Zustand for complex state
3. **No Query Caching** - TanStack Query not installed
4. **No UI Framework** - Tailwind not installed, inconsistent styling
5. **Feature Structure Unclear** - `features/` folder exists but organization unclear

#### Agent Issues

1. **No Offline Spooling** - Architecture mentions spooling but not implemented
2. **No Self-Updater** - Missing updater package
3. **No Packaging** - Missing systemd/service wrappers
4. **Collectors Not Platform-Separated** - No linux/ vs windows/ separation

---

## 2. Refactoring Roadmap

### Phase 1: Clean Architecture Foundation (Sprint 1)

**Goal:** Establish proper Clean Architecture boundaries without breaking existing functionality.

#### Week 1: Domain Layer & Interfaces

**Backend Structure:**
```
backend/
  cmd/
    api/                    # REST + WebSocket API server
    worker/                 # Metric/alert/rollup workers (NEW)
    scheduler/              # Offline detection, cleanup, AI jobs (NEW)
    migrate/                # Database migration runner (NEW)
  internal/
    domain/                 # NEW: Pure business entities
      organizations/
        entity.go
        repository.go       # Interface only
        service.go          # Interface only
      users/
        entity.go
        repository.go
        service.go
      machines/
        entity.go
        repository.go
        service.go
      metrics/
        entity.go
        repository.go
        service.go
      alerts/
        entity.go
        repository.go
        service.go
    app/                    # Dependency wiring
      api.go                # Rename from server.go
      worker.go             # NEW
      scheduler.go          # NEW
    config/                 # Configuration management
    infrastructure/         # NEW: External concerns
      postgres/
        connection.go
        transaction.go
        repository/         # GORM implementations
      websocket/
        hub.go
        client.go
      cache/
        redis.go
      queue/
        nats.go
    services/              # Business logic (keep, refactor)
    handlers/              # HTTP handlers (keep, refactor)
    middleware/            # Keep
    models/                # DEPRECATE: Move to domain/entities
    repository/            # DEPRECATE: Move to infrastructure/postgres/repository
    transport/             # NEW: HTTP/WebSocket adapters
      http/
        rest.go
        websocket.go
    pipeline/              # NEW: Ingestion queues, worker pools
      queue.go
      worker_pool.go
      fanout.go
    observability/         # NEW: Logging, metrics, tracing
      logger.go
      metrics.go
      tracer.go
```

**Actions:**
1. Create `internal/domain/` with pure entities (no GORM tags)
2. Define repository interfaces in `domain/*/repository.go`
3. Define service interfaces in `domain/*/service.go`
4. Move GORM models to `infrastructure/postgres/models/` (or keep as implementation detail)
5. Create `infrastructure/postgres/repository/` with GORM implementations

**Example:**
```go
// internal/domain/machines/entity.go
package machines

import "time"

type Machine struct {
    ID          string
    Hostname    string
    IPAddress   string
    OS          string
    Status      MachineStatus
    LastSeen    time.Time
    Metadata    map[string]string
}

type MachineStatus string

const (
    StatusOnline  MachineStatus = "ONLINE"
    StatusOffline MachineStatus = "OFFLINE"
)

// internal/domain/machines/repository.go
package machines

type Repository interface {
    Create(ctx context.Context, machine *Machine) error
    Update(ctx context.Context, machine *Machine) error
    FindByID(ctx context.Context, id string) (*Machine, error)
    FindByOrganization(ctx context.Context, orgID string) ([]Machine, error)
    UpdateStatus(ctx context.Context, id string, status MachineStatus) error
}

// internal/domain/machines/service.go
package machines

type Service interface {
    Register(ctx context.Context, input RegisterInput) (*Machine, error)
    UpdateStatus(ctx context.Context, id string, status MachineStatus) error
    GetMachines(ctx context.Context, orgID string) ([]Machine, error)
}
```

#### Week 2: Infrastructure Layer

**Actions:**
1. Create `infrastructure/postgres/repository/` with GORM implementations
2. Create `infrastructure/websocket/` with hub implementation
3. Create `infrastructure/cache/` with Redis client
4. Create `infrastructure/queue/` with NATS/RabbitMQ client
5. Implement repository interfaces from Week 1

**Example:**
```go
// internal/infrastructure/postgres/repository/machine_repository.go
package repository

import (
    "context"
    "infrapilot/backend/domain/machines"
    "infrapilot/backend/internal/models"
    "infrapilot/backend/internal/database"
)

type MachineRepository struct {
    db *database.DB
}

func NewMachineRepository(db *database.DB) machines.Repository {
    return &MachineRepository{db: db}
}

func (r *MachineRepository) Create(ctx context.Context, machine *machines.Machine) error {
    m := &models.Machine{
        ID:        uuid.MustParse(machine.ID),
        Hostname:  machine.Hostname,
        IPAddress: machine.IPAddress,
        // ... map domain entity to GORM model
    }
    return r.db.WithContext(ctx).Create(m).Error
}

func (r *MachineRepository) FindByOrganization(ctx context.Context, orgID string) ([]machines.Machine, error) {
    // ... implementation
}
```

#### Week 3: Transport Layer & Dependency Injection

**Actions:**
1. Create `transport/http/` with REST handlers
2. Create `transport/websocket/` with WebSocket handlers
3. Refactor `handlers/` to use DTOs and domain services
4. Implement dependency injection in `app/api.go`
5. Remove direct `database.DB` calls from handlers

**Example:**
```go
// internal/transport/http/machine_handler.go
package http

type MachineHandler struct {
    machineService machines.Service
}

func (h *MachineHandler) Register(c *gin.Context) {
    var req RegisterMachineRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, err)
        return
    }

    input := machines.RegisterInput{
        Hostname:  req.Hostname,
        IPAddress: req.IPAddress,
        // ...
    }

    machine, err := h.machineService.Register(c.Request.Context(), input)
    if err != nil {
        c.JSON(http.StatusInternalServerError, err)
        return
    }

    c.JSON(http.StatusCreated, machine)
}

// internal/app/api.go
func RunAPI() error {
    // Load config
    // Connect to database
    // Connect to Redis
    // Connect to NATS
    
    // Initialize repositories
    machineRepo := postgres.NewMachineRepository(db)
    
    // Initialize services
    machineService := machines.NewService(machineRepo)
    
    // Initialize handlers
    machineHandler := http.NewMachineHandler(machineService)
    
    // Setup routes
    router := gin.Default()
    api := router.Group("/api/v1")
    {
        machines := api.Group("/machines")
        {
            machines.POST("/register", machineHandler.Register)
            machines.GET("/:id", machineHandler.GetByID)
        }
    }
    
    return router.Run(":8080")
}
```

#### Week 4: Worker Pool & Pipeline

**Actions:**
1. Implement `cmd/worker/main.go`
2. Create `pipeline/` with ingestion queue
3. Create worker pools for metrics, alerts, rollups
4. Implement fan-out to WebSocket hub
5. Add graceful shutdown to all processes

**Example:**
```go
// cmd/worker/main.go
package main

func main() {
    // Load config
    // Connect to database
    // Connect to NATS
    // Connect to WebSocket hub
    
    // Create ingestion queue consumer
    metricQueue := pipeline.NewMetricQueue(natsClient)
    
    // Create worker pools
    metricWorkerPool := pipeline.NewWorkerPool(10, metricQueue.Consume)
    alertWorkerPool := pipeline.NewWorkerPool(5, alertQueue.Consume)
    rollupWorkerPool := pipeline.NewWorkerPool(2, rollupQueue.Consume)
    
    // Start workers
    metricWorkerPool.Start()
    alertWorkerPool.Start()
    rollupWorkerPool.Start()
    
    // Wait for shutdown signal
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan
    
    // Graceful shutdown
    metricWorkerPool.Stop()
    alertWorkerPool.Stop()
    rollupWorkerPool.Stop()
}
```

---

### Phase 2: AI Engine Foundation (Sprint 2)

**Goal:** Build the AI brain of InfraPilot.

**Structure:**
```
backend/
  internal/
    ai/                          # NEW: AI Engine
      anomaly.go                 # Anomaly detection
      prediction.go              # Capacity planning predictions
      recommendation.go          # AI-generated recommendations
      rootcause.go               # Root cause analysis
      summary.go                 # Incident summaries
      engine.go                  # Main AI engine orchestrator
      models.go                  # AI-specific models
      trainers/                  # Model training (future)
        anomaly_trainer.go
        prediction_trainer.go
```

**Integration Points:**
1. Workers feed historical metrics to AI engine
2. AI engine publishes predictions to WebSocket hub
3. Frontend displays AI insights in dashboard
4. Alerts enhanced with AI confidence scores

---

### Phase 3: Production Readiness (Sprint 3)

**Goal:** Make the platform production-ready.

**Checklist:**
- [ ] Structured logging (Zap/Zerolog)
- [ ] Distributed tracing (OpenTelemetry)
- [ ] Metrics (Prometheus client)
- [ ] Health checks (ready, live, startup)
- [ ] Graceful shutdown
- [ ] Configuration validation
- [ ] API rate limiting
- [ ] Request/response logging
- [ ] Error tracking (Sentry)
- [ ] Database connection pooling
- [ ] Redis connection pooling
- [ ] Circuit breakers
- [ ] Timeout configuration
- [ ] Retry policies with exponential backoff

---

### Phase 4: Frontend Modernization (Sprint 4)

**Goal:** Modernize frontend with TypeScript and proper state management.

**Actions:**
1. Migrate to TypeScript
2. Install TanStack Query for REST caching
3. Install Tailwind CSS for styling
4. Implement Zustand for state management
5. Create feature-based folder structure
6. Implement proper error boundaries
7. Add React Query DevTools
8. Add proper loading/error states

**Target Structure:**
```
frontend/src/
  app/
    store/                    # Zustand stores
    providers/                # Context providers
    hooks/                    # Custom hooks
  api/
    client.ts                 # Axios/Fetch client
    endpoints/                # API endpoint definitions
    types/                    # TypeScript types
  features/
    auth/
      components/
      hooks/
      store/
      types/
    machines/
      components/
      hooks/
      store/
      types/
    alerts/
      components/
      hooks/
      store/
      types/
    dashboard/
      components/
      hooks/
      store/
      types/
  components/
    ui/                       # Reusable UI components
    layout/                   # Layout components
  routes/
    index.tsx
  websocket/
    client.ts
    hooks/
  charts/
    components/
  utils/
    helpers.ts
    constants.ts
```

---

### Phase 5: Agent Enhancements (Sprint 5)

**Goal:** Make the agent production-ready.

**Actions:**
1. Implement offline spooling (local disk buffer)
2. Implement self-updater
3. Add platform-specific collectors (linux/ vs windows/)
4. Add systemd/service packaging
5. Implement exponential backoff retry
6. Add idempotency keys
7. Implement payload compression
8. Add schema versioning

---

### Phase 6: Deployment & DevOps (Sprint 6)

**Goal:** Production deployment infrastructure.

**Actions:**
1. Docker Compose for development
2. Kubernetes manifests for production
3. Helm charts
4. CI/CD pipelines (GitHub Actions)
5. Terraform for infrastructure
6. Monitoring stack (Prometheus + Grafana)
7. Log aggregation (Loki/ELK)
8. Alertmanager for platform alerts

---

## 3. Immediate Next Steps (This Week)

### Priority 1: Domain Layer (Days 1-2)

1. Create `internal/domain/` structure
2. Define pure entities for: User, Organization, Machine, Metric, Alert
3. Define repository interfaces
4. Define service interfaces

### Priority 2: Infrastructure Layer (Days 3-4)

1. Create `internal/infrastructure/postgres/repository/`
2. Implement repository interfaces with GORM
3. Create `internal/infrastructure/websocket/`
4. Move WebSocket hub to infrastructure layer

### Priority 3: Dependency Injection (Day 5)

1. Refactor `app/server.go` to `app/api.go`
2. Implement dependency injection
3. Wire up all layers
4. Test that everything still works

### Priority 4: Testing (Day 5)

1. Write unit tests for domain layer
2. Write integration tests for repositories
3. Write handler tests with mocked services

---

## 4. Code Quality Standards

### Testing Requirements
- **Domain Layer:** 100% unit test coverage
- **Services:** 80% unit test coverage
- **Handlers:** 70% integration test coverage
- **Repositories:** 60% integration test coverage

### Code Review Checklist
- [ ] No direct `database.DB` calls in handlers
- [ ] All dependencies injected via interfaces
- [ ] No framework imports in domain layer
- [ ] All errors handled explicitly
- [ ] All functions have godoc comments
- [ ] All exported types have godoc comments
- [ ] No magic numbers (use constants)
- [ ] All HTTP handlers validate input
- [ ] All HTTP handlers return proper status codes
- [ ] All sensitive data redacted in logs

### Naming Conventions
- **Packages:** lowercase, single word if possible
- **Interfaces:** `-er` suffix (e.g., `Repository`, `Service`)
- **Functions:** camelCase, verb-first (e.g., `GetMachineByID`)
- **Types:** PascalCase
- **Constants:** PascalCase or UPPER_SNAKE_CASE
- **Files:** lowercase, descriptive (e.g., `machine_repository.go`)

---

## 5. Migration Strategy

### Step 1: Parallel Implementation
- Create new structure alongside existing code
- Implement new layers without breaking existing functionality
- Run both old and new code in parallel

### Step 2: Gradual Migration
- Migrate one feature at a time
- Start with read-only operations (easier to test)
- Move to write operations
- Keep old code as fallback

### Step 3: Cutover
- Switch to new architecture
- Remove old code
- Update documentation

### Step 4: Validation
- Run full test suite
- Load test with production-like traffic
- Monitor for 1 week in production

---

## 6. Success Metrics

### Architecture Metrics
- [ ] 0 direct `database.DB` calls in handlers
- [ ] 100% of business logic in service layer
- [ ] 0 GORM imports in domain layer
- [ ] 0 Gin imports in domain layer
- [ ] 100% of dependencies injected via interfaces

### Performance Metrics
- [ ] API p99 latency < 200ms
- [ ] Worker ingestion throughput > 10,000 metrics/sec
- [ ] WebSocket broadcast latency < 100ms
- [ ] Dashboard load time < 2s

### Quality Metrics
- [ ] Test coverage > 70%
- [ ] 0 critical security vulnerabilities
- [ ] 0 lint errors
- [ ] All godoc comments present

---

## 7. Risk Mitigation

### Risk 1: Breaking Changes During Refactoring
**Mitigation:** Feature flags, parallel implementation, comprehensive testing

### Risk 2: Performance Degradation
**Mitigation:** Load testing at each phase, profiling, benchmarking

### Risk 3: Team Productivity Drop
**Mitigation:** Incremental migration, clear documentation, pair programming

### Risk 4: Production Incidents
**Mitigation:** Staging environment, canary deployments, rollback plan

---

## 8. Timeline Summary

| Sprint | Duration | Goal | Deliverable |
|--------|----------|------|-------------|
| 1 | 2 weeks | Clean Architecture | Domain layer, infrastructure layer, DI |
| 2 | 2 weeks | AI Engine | Anomaly detection, predictions, recommendations |
| 3 | 1 week | Production Readiness | Logging, tracing, metrics, health checks |
| 4 | 1 week | Frontend Modernization | TypeScript, TanStack Query, Tailwind |
| 5 | 1 week | Agent Enhancements | Spooling, updater, packaging |
| 6 | 2 weeks | Deployment & DevOps | K8s, CI/CD, monitoring |

**Total:** 9 weeks to production-ready enterprise platform

---

## 9. Conclusion

The current codebase is a solid foundation but needs architectural improvements to become production-ready. This plan provides a clear, incremental path to a scalable, maintainable, enterprise-grade platform.

**Next Action:** Begin Sprint 1, Week 1 - Domain Layer implementation.