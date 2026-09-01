# Technical Blog Posts

In-depth technical articles demonstrating architecture decisions and implementation details.

---

## Blog Post 1: How I Built an Enterprise Monitoring Platform in Go

**Title:** How I Built an Enterprise Monitoring Platform in Go  
**Target Publication:** DEV.to, Medium, Personal Blog  
**Estimated Read Time:** 15 minutes  
**Tags:** `go`, `monitoring`, `backend`, `architecture`, `postgresql`, `redis`

### Introduction

After years of using fragmented monitoring tools, I decided to build **InfraPilot Enterprise** — a unified, real-time infrastructure monitoring platform. This article walks through the architecture, design decisions, and implementation details.

### The Problem

Modern DevOps teams typically use 5+ tools for monitoring, logging, alerts, and management:
- **Prometheus + Grafana** for metrics
- **ELK Stack** for logs
- **AlertManager** for alerts
- **SSH + scripts** for remote management
- **Custom dashboards** for analytics

This creates tool sprawl, increased costs, and cognitive overload.

### Architecture Overview

```
┌─────────────┐    WebSocket     ┌────────────────┐
│   Agents    │ ────────────────▶│   Go Backend   │
│ (Go)        │  REST API        │   (Gin)        │
└─────────────┘                  └────────┬───────┘
                                            │
                        ┌───────────────────┼────────────────────┐
                        │                   │                    │
                   ┌────▼─────┐        ┌───▼────┐         ┌──────▼──────┐
                   │PostgreSQL│        │  Redis │         │   Workers   │
                   │  (GORM)  │        │ Streams│         │  (Pool)     │
                   └──────────┘        └────────┘         └─────────────┘
```

### Key Design Decisions

#### 1. Why Go?

I chose Go for three reasons:
- **Concurrency**: Goroutines handle thousands of concurrent agent connections
- **Performance**: Compiled binary starts in milliseconds, low memory footprint (28MB)
- **Simplicity**: Easy to deploy single binary, cross-compilation for multiple platforms

**Code Example - Handling 10K WebSocket Connections:**

```go
type WebSocketHub struct {
    clients    map[*WebSocketClient]bool
    broadcast  chan []byte
    register   chan *WebSocketClient
    unregister chan *WebSocketClient
}

func (h *WebSocketHub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
            log.Printf("Client connected. Total: %d", len(h.clients))
        
        case client := <-h.unregister:
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
        
        case message := <-h.broadcast:
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
        }
    }
}
```

#### 2. Redis Streams as Event Bus

Instead of Kafka (which is overkill), I used **Redis Streams**:

```go
type EventBus struct {
    redisClient *redis.Client
    streamName  string
    consumerGroup string
}

func (eb *EventBus) Publish(event Event) error {
    data, _ := json.Marshal(event)
    return eb.redisClient.XAdd(ctx, &redis.XAddArgs{
        Stream: eb.streamName,
        Values: map[string]interface{}{"data": data},
    }).Err()
}

func (eb *EventBus) Subscribe(handler EventHandler) {
    for {
        streams, err := eb.redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
            Group:    eb.consumerGroup,
            Consumer: "metrics-worker",
            Streams:  []string{eb.streamName, ">"},
            Block:    5 * time.Second,
        }).Result()
        
        for _, stream := range streams {
            for _, message := range stream.Messages {
                eb.processMessage(message, handler)
            }
        }
    }
}
```

**Benefits:** No external dependencies, persistent queue, consumer groups, sub-millisecond latency.

#### 3. Worker Pool Pattern

To handle 100K+ metrics/second:

```go
type WorkerPool struct {
    workerCount int
    jobQueue    chan Job
    workers     []*Worker
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.workerCount; i++ {
        worker := &Worker{
            id:       i,
            jobQueue: make(chan Job, 100),
        }
        wp.workers = append(wp.workers, worker)
        go worker.Start()
    }
    
    // Dispatcher
    go func() {
        for job := range wp.jobQueue {
            wp.dispatch(job)
        }
    }()
}

func (w *Worker) Start() {
    for job := range w.jobQueue {
        job.Execute()
    }
}
```

**Performance achieved:** 150K metrics/second with 10 workers on a 4-core machine.

#### 4. Database Schema Design

```sql
-- Machines with flexible JSONB metadata
CREATE TABLE machines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hostname VARCHAR(255) UNIQUE NOT NULL,
    ip_address INET,
    os JSONB,  -- {"type": "linux", "version": "22.04"}
    status VARCHAR(50) DEFAULT 'online',
    last_heartbeat TIMESTAMP DEFAULT NOW(),
    enrolled_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_machines_status ON machines(status);
CREATE INDEX idx_machines_heartbeat ON machines(last_heartbeat);

-- Metrics with time-series optimization
CREATE TABLE metrics (
    time TIMESTAMP NOT NULL DEFAULT NOW(),
    machine_id UUID REFERENCES machines(id),
    cpu_usage FLOAT,
    memory_usage FLOAT,
    disk_usage FLOAT,
    network_in BIGINT,
    network_out BIGINT
);

SELECT create_hypertable('metrics', 'time', chunk_time_interval => INTERVAL '1 day');
CREATE INDEX idx_metrics_machine_time ON metrics(machine_id, time DESC);
```

**Key optimizations:**
- UUID keys for distributed systems
- JSONB for flexible metadata
- Hypertable partitioning for TimescaleDB
- Composite indexes for common queries

### Frontend Architecture

React with Context API for state:

```jsx
// contexts/WebSocketContext.js
import { createContext, useContext, useEffect, useState } from 'react';

export const WebSocketContext = createContext();

export function WebSocketProvider({ children, url }) {
    const [metrics, setMetrics] = useState({});
    const [socket, setSocket] = useState(null);
    
    useEffect(() => {
        const ws = new WebSocket(url);
        ws.onmessage = (event) => {
            const data = JSON.parse(event.data);
            setMetrics(prev => ({...prev, [data.machine_id]: data}));
        };
        setSocket(ws);
        return () => ws.close();
    }, [url]);
    
    return (
        <WebSocketContext.Provider value={{ metrics, socket }}>
            {children}
        </WebSocketContext.Provider>
    );
}

// Usage in components
function MetricsCard({ machineId }) {
    const { metrics } = useContext(WebSocketContext);
    const data = metrics[machineId];
    
    if (!data) return <div>Loading...</div>;
    
    return (
        <div className="metric-card">
            <h3>CPU</h3>
            <p>{data.cpu}%</p>
        </div>
    );
}
```

### Deployment Strategy

**Docker Multi-Stage Build:**

```dockerfile
# Stage 1: Build
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Stage 2: Production
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/configs/ ./configs/
EXPOSE 8080 9090
CMD ["./main"]
```

**Result:** 1.2GB → 12MB image!

### Lessons Learned

1. **Start with the event log first** - It becomes your debugging lifeline
2. **Implement rate limiting early** - You'll thank yourself when you have bugs
3. **Use connection pooling** - Database connections are expensive
4. **WebSocket heartbeats are mandatory** - Detect dead connections within 60s
5. **Structured logging from day one** - Correlation IDs save hours of debugging

### Performance Results

After optimization:
- **API Latency (p95):** 32ms
- **Metrics Ingestion:** 150K/sec
- **WebSocket Connections:** 12,500 concurrent
- **Dashboard Load:** 145ms
- **Agent Memory:** 28MB (including Docker/K8s plugins)

### Conclusion

Building this platform taught me that:
- **Go is perfect** for I/O-bound monitoring workloads
- **Event-driven architecture** provides excellent decoupling
- **Redis Streams** are powerful enough for most use cases without Kafka
- **WebSocket backpressure** handling is critical for production systems

The full source code is available at [github.com/venkaiswami/infrapilot-enterprise](https://github.com/venkaiswami/infrapilot-enterprise).

---

## Blog Post 2: Building a Cross-Platform Monitoring Agent

**Title:** Building a Cross-Platform Monitoring Agent  
**Target:** Go blog, GolangWeekly  
**Estimated Read Time:** 12 minutes  
**Tags:** `go`, `cross-platform`, `agent`, `docker`, `kubernetes`, `plugin-system`

### Introduction

Learn how I built a monitoring agent that works on Linux, Windows, and macOS while collecting system metrics, Docker stats, and Kubernetes information.

### Architecture

```go
// Plugin interface for cross-platform support
type Plugin interface {
    Name() string
    Collect(ctx context.Context) (map[string]interface{}, error)
    Enabled() bool
}
```

### Platform Detection

```go
func GetPlatform() string {
    runtimeOS := runtime.GOOS
    switch runtimeOS {
    case "linux":
        return "linux"
    case "windows":
        return "windows"
    case "darwin":
        return "darwin"
    default:
        return "unknown"
    }
}
```

### Plugin System

```go
type PluginManager struct {
    plugins map[string]Plugin
}

func (pm *PluginManager) Register(plugin Plugin) {
    pm.plugins[plugin.Name()] = plugin
}

func (pm *PluginManager) CollectAll(ctx context.Context) map[string]map[string]interface{} {
    results := make(map[string]map[string]interface{})
    var wg sync.WaitGroup
    
    for name, plugin := range pm.plugins {
        if plugin.Enabled() {
            wg.Add(1)
            go func(name string, p Plugin) {
                defer wg.Done()
                data, _ := p.Collect(ctx)
                results[name] = data
            }(name, plugin)
        }
    }
    
    wg.Wait()
    return results
}
```

### Linux Plugin Example

```go
type LinuxPlugin struct{}

func (p *LinuxPlugin) Name() string { return "linux" }

func (p *LinuxPlugin) Collect(ctx context.Context) (map[string]interface{}, error) {
    // CPU
    cpu, _ := p.getCPUUsage()
    
    // Memory
    mem, _ := p.getMemoryUsage()
    
    // Disk
    disk, _ := p.getDiskUsage()
    
    return map[string]interface{}{
        "cpu":       cpu,
        "memory":    mem,
        "disk":      disk,
        "processes": p.getTopProcesses(),
        "services":  p.getSystemdServices(),
    }, nil
}
```

### Cross-Compilation

```bash
# Build for all platforms
GOOS=linux GOARCH=amd64 go build -o agent-linux-amd64
GOOS=linux GOARCH=arm64 go build -o agent-linux-arm64
GOOS=darwin GOARCH=amd64 go build -o agent-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o agent-darwin-arm64
GOOS=windows GOARCH=amd64 go build -o agent-windows-amd64.exe
```

---

## Blog Post 3: Implementing Real-Time Monitoring with WebSockets

**Title:** Implementing Real-Time Monitoring with WebSockets  
**Target:** [LogRocket](https://blog.logrocket.com), [Smashing Magazine](https://www.smashingmagazine.com)  
**Estimated Read Time:** 10 minutes  
**Tags:** `websocket`, `react`, `go`, `real-time`, `scaling`

### Introduction

Building a real-time dashboard that updates sub-second across thousands of machines requires careful WebSocket architecture.

### The Scaling Challenge

- 1,000 machines sending metrics every 5 seconds
- 500 concurrent dashboard users
- Average message size: 2KB
- Bandwidth: ~4MB/second

### Go Backend - WebSocket Hub

```go
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    mu         sync.RWMutex
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan []byte, 256),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client] = true
            h.mu.Unlock()
            
            // Subscribe to machine's metrics
            go client.subscribeToMachine(client.machineID)
        
        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.Send)
            }
            h.mu.Unlock()
        
        case message := <-h.broadcast:
            h.mu.RLock()
            for client := range h.clients {
                select {
                case client.Send <- message:
                default:
                    // Buffer full, drop message
                    log.Warn("Client buffer full, dropping message")
                }
            }
            h.mu.RUnlock()
        }
    }
}
```

### React WebSocket Hook

```jsx
function useMetrics(machineId) {
    const [metrics, setMetrics] = useState(null);
    const [ws, setWs] = useState(null);
    
    useEffect(() => {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const websocket = new WebSocket(`${protocol}//${window.location.host}/ws?machine=${machineId}`);
        
        websocket.onmessage = (event) => {
            const data = JSON.parse(event.data);
            setMetrics(data);
        };
        
        websocket.onerror = (error) => {
            console.error('WebSocket error:', error);
        };
        
        setWs(websocket);
        
        return () => {
            websocket.close();
        };
    }, [machineId]);
    
    return metrics;
}

// Usage
function MachineMetrics({ machineId }) {
    const metrics = useMetrics(machineId);
    
    if (!metrics) return <div>Connecting...</div>;
    
    return (
        <div>
            <h3>{metrics.hostname}</h3>
            <p>CPU: {metrics.cpu}%</p>
            <p>Memory: {metrics.memory}%</p>
            <p>Status: {metrics.status}</p>
        </div>
    );
}
```

### Handling Reconnection

```jsx
function useReconnectingWebSocket(url, retries = 5) {
    const [ws, setWs] = useState(null);
    const [status, setStatus] = useState('connecting');
    
    useEffect(() => {
        let wsRef;
        let retryCount = 0;
        
        const connect = () => {
            const websocket = new WebSocket(url);
            wsRef = websocket;
            
            websocket.onopen = () => {
                setStatus('connected');
                retryCount = 0;
            };
            
            websocket.onclose = () => {
                setStatus('disconnected');
                
                // Exponential backoff
                if (retryCount < retries) {
                    const delay = Math.min(1000 * Math.pow(2, retryCount), 30000);
                    setTimeout(connect, delay);
                    retryCount++;
                }
            };
            
            setWs(websocket);
        };
        
        connect();
        
        return () => wsRef?.close();
    }, [url, retries]);
    
    return { ws, status };
}
```

### Performance Optimizations

1. **Message Batching**: Send multiple metrics in a single message
2. **Compression**: Enable permessage-deflate
3. **Binary Protocol**: Use Protocol Buffers instead of JSON (3x smaller)
4. **Heartbeat**: Detect dead connections within 60s

---

## Blog Post 4: Monitoring Docker and Kubernetes from Scratch

**Title:** Monitoring Docker and Kubernetes from Scratch  
**Target:** [Container Journal](https://containerjournal.com), [CNCF Blog](https://www.cncf.io)  
**Estimated Read Time:** 14 minutes  
**Tags:** `docker`, `kubernetes`, `container`, `monitoring`, `devops`

### Introduction

Learn how I implemented Docker and Kubernetes monitoring without relying on cAdvisor or Heapster.

### Docker Engine API Integration

```go
type DockerPlugin struct {
    client *docker.Client
}

func (p *DockerPlugin) Name() string { return "docker" }

func (p *DockerPlugin) Collect(ctx context.Context) (map[string]interface{}, error) {
    containers, _, err := p.client.ContainerList(ctx, 
        docker.ContainerListOptions{All: true})
    
    if err != nil {
        return nil, err
    }
    
    stats := make(map[string]interface{})
    
    for _, container := range containers {
        jsonStats, err := p.client.ContainerStats(ctx, container.ID, false)
        if err != nil {
            continue
        }
        
        var stats docker.Stats
        json.NewDecoder(jsonStats.Body).Decode(&stats)
        
        cpuPercent := calculateCPUPercent(stats)
        memPercent := calculateMemoryPercent(stats)
        
        stats[container.ID] = map[string]interface{}{
            "name":    container.Names[0],
            "image":   container.Image,
            "status":  container.State,
            "cpu":     cpuPercent,
            "memory":  memPercent,
        }
    }
    
    return stats, nil
}

func calculateCPUPercent(stats docker.Stats) float64 {
    cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage - stats.PreCPUStats.CPUUsage.TotalUsage)
    systemDelta := float64(stats.CPUStats.SystemUsage - stats.PreCPUStats.SystemUsage)
    
    if systemDelta > 0.0 && cpuDelta > 0.0 {
        return (cpuDelta / systemDelta) * float64(len(stats.CPUStats.CPUUsage.PercpuUsage)) * 100.0
    }
    return 0.0
}
```

### Kubernetes Client Integration

```go
type KubernetesPlugin struct {
    config    *rest.Config
    clientset *kubernetes.Clientset
}

func (p *KubernetesPlugin) Name() string { return "kubernetes" }

func (p *KubernetesPlugin) Collect(ctx context.Context) (map[string]interface{}, error) {
    nodes, err := p.clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
    if err != nil {
        return nil, err
    }
    
    pods, err := p.clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
    if err != nil {
        return nil, err
    }
    
    result := map[string]interface{}{
        "nodes": p.parseNodes(nodes.Items),
        "pods":  p.parsePods(pods.Items),
    }
    
    return result, nil
}

func (p *KubernetesPlugin) parseNodes(nodes []v1.Node) []NodeInfo {
    var nodeInfos []NodeInfo
    
    for _, node := range nodes {
        cpuCapacity := node.Status.Capacity.Cpu().MilliValue()
        cpuAllocatable := node.Status.Allocatable.Cpu().MilliValue()
        memCapacity := node.Status.Capacity.Memory().Value()
        memAllocatable := node.Status.Allocatable.Memory().Value()
        
        nodeInfo := NodeInfo{
            Name:           node.Name,
            Status:         string(node.Status.Conditions[len(node.Status.Conditions)-1].Type),
            CPU:            float64(cpuAllocatable) / float64(cpuCapacity) * 100,
            Memory:         float64(memAllocatable) / float64(memCapacity) * 100,
            KernelVersion:  node.Status.NodeInfo.KernelVersion,
            OSImage:        node.Status.NodeInfo.OSImage,
            ContainerRuntime: node.Status.NodeInfo.ContainerRuntimeVersion,
        }
        
        nodeInfos = append(nodeInfos, nodeInfo)
    }
    
    return nodeInfos
}
```

### Pod Status Tracking

```go
func (p *KubernetesPlugin) getPodStatus(pod v1.Pod) string {
    for _, condition := range pod.Status.Conditions {
        if condition.Type == v1.PodReady {
            if condition.Status == "True" {
                return "Running"
            }
            return "Not Ready"
        }
    }
    
    for _, containerStatus := range pod.Status.ContainerStatuses {
        if containerStatus.State.Waiting != nil {
            return containerStatus.State.Waiting.Reason
        }
        if containerStatus.State.Terminated != nil {
            return "Terminated"
        }
    }
    
    return "Unknown"
}
```

### Service Discovery

```go
func (p *KubernetesPlugin) watchServices(ctx context.Context) {
    watcher, err := p.clientset.CoreV1().Services("").Watch(ctx, metav1.ListOptions{})
    if err != nil {
        log.Fatal(err)
    }
    
    for event := range watcher.ResultChan() {
        service := event.Object.(*v1.Service)
        
        // Push to metrics stream
        event := Event{
            Type: "service",
            Data: map[string]interface{}{
                "name":      service.Name,
                "namespace": service.Namespace,
                "type":      string(service.Spec.Type),
                "cluster_ip": service.Spec.ClusterIP,
                "ports":    p.extractPorts(service.Spec.Ports),
            },
        }
        
        p.eventBus.Publish(event)
    }
}
```

### Results

- **Docker stats**: 100ms latency
- **Kubernetes nodes**: Updated every 5s
- **Pod status**: Real-time via watch API
- **Resource overhead**: < 2% CPU, 30MB memory

### Key takeaways

1. Use Docker Engine API for container stats (don't parse `/sys/fs/cgroup`)
2. Use Kubernetes informers (not polling) for resource efficiency
3. Implement RBAC with minimal permissions
4. Cache resource lists to reduce API calls

---

## Blog Post 5: Event-Driven Architecture with Redis Streams

**Title:** Event-Driven Architecture with Redis Streams  
**Target:** [InfoQ](https://www.infoq.com), [HighScalability](http://highscalability.com)  
**Estimated Read Time:** 11 minutes  
**Tags:** `redis`, `event-driven`, `architecture`, `microservices`, `go`

### Introduction

How Redis Streams provides a simple yet powerful message bus without Kafka's operational complexity.

### Why Event-Driven?

- **Decoupling**: Services communicate without knowing about each other
- **Scalability**: Easy to add consumers
- **Reliability**: Persisted events survive crashes
- **Replayability**: Reprocess events for debugging

### Event Bus Implementation

```go
type EventBus struct {
    client     *redis.Client
    streamName string
    consumers  map[string]EventHandler
    mu         sync.RWMutex
}

func NewEventBus(redisAddr string) (*EventBus, error) {
    client := redis.NewClient(&redis.Options{
        Addr: redisAddr,
    })
    
    // Create consumer group
    err := client.XGroupCreateMkStream(ctx, 
        "infrapilot:events",
        "metrics-consumers",
        "0").Err()
    
    return &EventBus{
        client:     client,
        streamName: "infrapilot:events",
        consumers:  make(map[string]EventHandler),
    }, nil
}

func (eb *EventBus) Publish(event Event) error {
    payload, _ := json.Marshal(event)
    
    return eb.client.XAdd(ctx, &redis.XAddArgs{
        Stream: eb.streamName,
        Values: map[string]interface{}{
            "type":    event.Type,
            "payload": string(payload),
            "ts":      time.Now().UnixMilli(),
        },
        MaxLen: 10000, // Keep last 10K events
    }).Err()
}

type EventHandler interface {
    Handle(ctx context.Context, event Event) error
    Name() string
}
```

### Subscribers

```go
// Database subscriber - persists to PostgreSQL
type DatabaseSubscriber struct {
    db *gorm.DB
}

func (s *DatabaseSubscriber) Handle(ctx context.Context, event Event) error {
    switch event.Type {
    case "metric":
        return s.saveMetric(event)
    case "alert":
        return s.saveAlert(event)
    case "machine":
        return s.saveMachine(event)
    }
    return nil
}

func (s *DatabaseSubscriber) Name() string {
    return "database-subscriber"
}

// WebSocket subscriber - pushes to connected clients
type WebSocketSubscriber struct {
    hub *WebSocketHub
}

func (s *WebSocketSubscriber) Handle(ctx context.Context, event Event) error {
    message, _ := json.Marshal(event)
    s.hub.Broadcast(message)
    return nil
}
```

### Running Multiple Consumers

```go
func main() {
    eventBus, _ := NewEventBus("localhost:6379")
    
    // Register consumers
    eventBus.Register(&DatabaseSubscriber{db: db})
    eventBus.Register(&WebSocketSubscriber{hub: hub})
    eventBus.Register(&AlertSubscriber{})
    eventBus.Register(&AISubscriber{})
    eventBus.Register(&AuditSubscriber{})
    
    // Start consuming (fan-out)
    eventBus.StartConsumers(ctx)
}
```

### Benefits Over Kafka

| Feature | Redis Streams | Kafka |
|---------|--------------|-------|
| Setup complexity | Low (1 binary) | High (ZooKeeper + brokers) |
| Memory | Low | High |
| Learning curve | Easy | Steep |
| Persistence | Yes | Yes |
| Consumer groups | Yes | Yes |
| Use case fit | Perfect (<100K msg/s) | Overkill here |

---

## Blog Post 6: Implementing RBAC with JWT Tokens in Go

**Title:** Implementing RBAC with JWT Tokens in Go  
**Target:** [Go.dev Blog](https://go.dev/blog), [Auth0 Blog](https://auth0.com/blog)  
**Estimated Read Time:** 9 minutes  
**Tags:** `jwt`, `authentication`, `authorization`, `go`, `security`

### Introduction

Building enterprise-grade authentication with role-based access control.

### JWT Middleware

```go
type Claims struct {
    UserID   string `json:"user_id"`
    Email    string `json:"email"`
    Role     string `json:"role"`
    Perms    []string `json:"permissions"`
    jwt.RegisteredClaims
}

func JWTMiddleware(secret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.AbortWithStatusJSON(401, gin.H{"error": "Missing token"})
            return
        }
        
        token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
            return []byte(secret), nil
        })
        
        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
            return
        }
        
        claims := token.Claims.(*Claims)
        c.Set("user_id", claims.UserID)
        c.Set("role", claims.Role)
        c.Set("permissions", claims.Perms)
        c.Next()
    }
}
```

### Role-Based Authorization

```go
type Role int

const (
    SuperAdmin Role = iota
    Admin
    Operator
    Viewer
)

var (
    roleHierarchy = map[Role]int{
        SuperAdmin:  4,
        Admin:       3,
        Operator:    2,
        Viewer:      1,
    }
    
    rolePermissions = map[Role][]string{
        SuperAdmin: []string{"*"},
        Admin:      []string{"read:*", "write:*", "delete:machines"},
        Operator:   []string{"read:*", "write:machines"},
        Viewer:     []string{"read:machines"},
    }
)

func RequirePermission(permission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        role := c.MustGet("role").(Role)
        perms := rolePermissions[role]
        
        for _, p := range perms {
            if p == "*" || p == permission {
                c.Next()
                return
            }
            
            // Wildcard matching
            if strings.HasSuffix(p, ":*") && strings.HasPrefix(permission, strings.TrimSuffix(p, "*")) {
                c.Next()
                return
            }
        }
        
        c.AbortWithStatusJSON(403, gin.H{"error": "Forbidden"})
    }
}
```

### Usage in Routes

```go
r := gin.Default()

api := r.Group("/api/v1")
{
    // SuperAdmin only
    admin := api.Group("/admin")
    admin.Use(JWTMiddleware(secret))
    admin.Use(RequirePermission("write:*"))
    {
        admin.POST("/users", createUser)
        admin.DELETE("/users/:id", deleteUser)
    }
    
    // Operators and above
    machines := api.Group("/machines")
    machines.Use(JWTMiddleware(secret))
    machines.Use(RequirePermission("write:machines"))
    {
        machines.POST("/", createMachine)
        machines.PUT("/:id", updateMachine)
    }
    
    // Viewers and above
    api.GET("/machines", listMachines)
}
```

### Login Handler

```go
func Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    user, err := db.GetUserByEmail(req.Email)
    if err != nil || !bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) {
        c.JSON(401, gin.H{"error": "Invalid credentials"})
        return
    }
    
    claims := Claims{
        UserID: user.ID,
        Email:  user.Email,
        Role:   Role(user.Role),
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, _ := token.SignedString([]byte(secret))
    
    c.JSON(200, gin.H{"token": tokenString})
}
```

---

## LinkedIn Posts

See [docs/linkedin-posts.md](docs/linkedin-posts.md) for social media content.

---

## Writing Tips

### Structure

1. **Hook** (1 paragraph) - Problem or surprising fact
2. **Context** (2-3 paragraphs) - Why this matters
3. **Technical Deep Dive** (main content) - Code examples, diagrams
4. **Results** (1-2 paragraphs) - Performance metrics, outcomes
5. **Lessons Learned** (bulleted list) - Takeaways
6. **Call to Action** - Star repo, follow, feedback

### Code Formatting

- Use syntax highlighting with language tags
- Keep examples under 50 lines
- Add comments for complex logic
- Link to full source on GitHub

### Promotion

- Post on Reddit: r/golang, r/devops, r/kubernetes
- Share on Hacker News
- Tweet with screenshots/diagrams
- Cross-post to DEV.to, Medium
- Add to your portfolio

---

<p align="center">
  Contribute your own blog post! See [CONTRIBUTING.md](CONTRIBUTING.md)
</p>