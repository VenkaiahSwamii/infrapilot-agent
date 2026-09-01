# InfraPilot Public API - SDK Examples

Official SDK examples and guides for integrating with InfraPilot Enterprise.

## Table of Contents

- [Overview](#overview)
- [Authentication](#authentication)
- [Go SDK](#go-sdk)
- [Python SDK](#python-sdk)
- [JavaScript/TypeScript SDK](#javascripttypescript-sdk)
- [REST API Examples](#rest-api-examples)
- [WebSocket API](#websocket-api)
- [Best Practices](#best-practices)

## Overview

The InfraPilot Public API provides REST and WebSocket interfaces for:

- Machine management
- Metrics collection and retrieval
- Alert management
- Remote operations (terminal, files, services)
- Docker and Kubernetes management
- Report generation

### Base URL

```bash
http://localhost:8080/api/v1
```

### Interactive Documentation

- Swagger UI: `http://localhost:8080/swagger/index.html`
- OpenAPI Spec: `http://localhost:8080/swagger/doc.json`
- ReDoc: `http://localhost:8080/redoc`

## Authentication

All API requests require authentication via JWT Bearer token.

### Getting a Token

```bash
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "admin@infrapilot.io",
  "password": "your-password"
}

Response:
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "uuid",
    "email": "admin@infrapilot.io",
    "role": "admin"
  }
}
```

### Using the Token

```bash
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Refreshing Tokens

```bash
POST /api/v1/auth/refresh
Authorization: Bearer <refresh-token>

Response:
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

## Go SDK

### Installation

```bash
go get github.com/venkaiswami/infrapilot-enterprise/sdk/go
```

### Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/venkaiswami/infrapilot-enterprise/sdk/go"
)

func main() {
    // Initialize client
    client := infrapilot.NewClient(&infrapilot.ClientConfig{
        BaseURL:  "http://localhost:8080/api/v1",
        Username: "admin@infrapilot.io",
        Password: "your-password",
    })

    // Test connection
    ctx := context.Background()
    health, err := client.Health.Check(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Status: %s\n", health.Status)

    // List machines
    machines, err := client.Machines.List(ctx, &infrapilot.MachineListOptions{
        Page:  1,
        Limit: 20,
    })
    if err != nil {
        log.Fatal(err)
    }

    for _, machine := range machines.Data {
        fmt.Printf("Machine: %s (%s) - %s\n", 
            machine.Name, 
            machine.IPAddress, 
            machine.Status)
    }
}
```

### Complete Examples

#### List Machines

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/venkaiswami/infrapilot-enterprise/sdk/go"
)

func main() {
    client := infrapilot.NewClient(&infrapilot.ClientConfig{
        BaseURL:  "http://localhost:8080/api/v1",
        Username: "admin@infrapilot.io",
        Password: "your-password",
    })

    ctx := context.Background()

    // List online machines
    machines, err := client.Machines.List(ctx, &infrapilot.MachineListOptions{
        Status: "online",
        Page:   1,
        Limit:  50,
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Total machines: %d\n", machines.Total)
    for _, machine := range machines.Data {
        fmt.Printf("- %s (CPU: %.1f%%, Memory: %.1f%%)\n",
            machine.Name,
            machine.Metrics.CPU,
            machine.Metrics.Memory)
    }
}
```

#### Get Metrics

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/venkaiswami/infrapilot-enterprise/sdk/go"
)

func main() {
    client := infrapilot.NewClient(&infrapilot.ClientConfig{
        BaseURL:  "http://localhost:8080/api/v1",
        Username: "admin@infrapilot.io",
        Password: "your-password",
    })

    ctx := context.Background()

    // Get machine ID (from list)
    machineID := "machine-uuid-here"

    // Get metrics for last 24 hours
    metrics, err := client.Machines.GetMetrics(ctx, machineID, &infrapilot.MetricsOptions{
        Period: "24h",
        Metrics: []string{"cpu", "memory", "disk", "network"},
    })
    if err != nil {
        log.Fatal(err)
    }

    // Print CPU history
    fmt.Println("CPU History:")
    for _, point := range metrics.CPU {
        timestamp, _ := time.Parse(time.RFC3339, point.Timestamp)
        fmt.Printf("  %s: %.1f%%\n", 
            timestamp.Format("15:04:05"), 
            point.Value)
    }
}
```

#### Create and Manage Alerts

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/venkaiswami/infrapilot-enterprise/sdk/go"
)

func main() {
    client := infrapilot.NewClient(&infrapilot.ClientConfig{
        BaseURL:  "http://localhost:8080/api/v1",
        Username: "admin@infrapilot.io",
        Password: "your-password",
    })

    ctx := context.Background()

    // Create alert rule
    alertRule, err := client.Alerts.CreateRule(ctx, &infrapilot.AlertRuleRequest{
        Name:        "High CPU Usage",
        MachineID:   "machine-uuid",
        Metric:      "cpu",
        Condition:   ">",
        Threshold:   90.0,
        Severity:    "warning",
        Description: "CPU usage exceeded 90%",
        Duration:    "5m",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Created alert rule: %s\n", alertRule.ID)

    // List active alerts
    alerts, err := client.Alerts.List(ctx, &infrapilot.AlertListOptions{
        Status:   "active",
        Severity: "critical",
    })
    if err != nil {
        log.Fatal(err)
    }

    for _, alert := range alerts.Data {
        fmt.Printf("Alert: %s - %s\n", alert.Name, alert.Message)
    }

    // Acknowledge alert
    err = client.Alerts.Acknowledge(ctx, alert.ID)
    if err != nil {
        log.Fatal(err)
    }

    // Resolve alert
    err = client.Alerts.Resolve(ctx, alert.ID)
    if err != nil {
        log.Fatal(err)
    }
}
```

#### Execute Remote Commands

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/venkaiswami/infrapilot-enterprise/sdk/go"
)

func main() {
    client := infrapilot.NewClient(&infrapilot.ClientConfig{
        BaseURL:  "http://localhost:8080/api/v1",
        Username: "admin@infrapilot.io",
        Password: "your-password",
    })

    ctx := context.Background()

    // Create terminal session
    session, err := client.Terminal.CreateSession(ctx, &infrapilot.TerminalSessionRequest{
        MachineID: "machine-uuid",
        Shell:     "/bin/bash",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Terminal session: %s\n", session.SessionID)
    fmt.Printf("Connect via: %s\n", session.WebSocketURL)

    // Execute single command
    result, err := client.Machines.ExecuteCommand(ctx, "machine-uuid", &infrapilot.CommandRequest{
        Command: "systemctl status nginx",
        Timeout: 30,
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Output:\n%s\n", result.Output)
    fmt.Printf("Exit code: %d\n", result.ExitCode)
}
```

#### Manage Files

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/venkaiswami/infrapilot-enterprise/sdk/go"
)

func main() {
    client := infrapilot.NewClient(&infrapilot.ClientConfig{
        BaseURL:  "http://localhost:8080/api/v1",
        Username: "admin@infrapilot.io",
        Password: "your-password",
    })

    ctx := context.Background()

    // List files
    files, err := client.Files.List(ctx, "machine-uuid", &infrapilot.FileListOptions{
        Path: "/var/log",
    })
    if err != nil {
        log.Fatal(err)
    }

    for _, file := range files {
        fmt.Printf("%s\t%s\t%d bytes\t%s\n",
            file.Type,
            file.Name,
            file.Size,
            file.Modified)
    }

    // Download file
    fileContent, err := client.Files.Download(ctx, "machine-uuid", "/var/log/syslog")
    if err != nil {
        log.Fatal(err)
    }

    os.WriteFile("syslog.log", fileContent, 0644)

    // Upload file
    err = client.Files.Upload(ctx, "machine-uuid", "/tmp/upload.txt", []byte("Hello, InfraPilot!"))
    if err != nil {
        log.Fatal(err)
    }
}
```

#### Manage Docker Containers

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/venkaiswami/infrapilot-enterprise/sdk/go"
)

func main() {
    client := infrapilot.NewClient(&infrapilot.ClientConfig{
        BaseURL:  "http://localhost:8080/api/v1",
        Username: "admin@infrapilot.io",
        Password: "your-password",
    })

    ctx := context.Background()

    // List containers
    containers, err := client.Docker.ListContainers(ctx, "machine-uuid")
    if err != nil {
        log.Fatal(err)
    }

    for _, container := range containers {
        fmt.Printf("Container: %s (%s)\n", container.Name, container.Image)
        fmt.Printf("  Status: %s, CPU: %.1f%%, Memory: %.0fMB\n",
            container.Status,
            container.CPU,
            container.Memory)
    }

    // Restart container
    err = client.Docker.ContainerAction(ctx, "machine-uuid", "container-id", "restart")
    if err != nil {
        log.Fatal(err)
    }
}
```

## Python SDK

### Installation

```bash
pip install infrapilot
```

### Quick Start

```python
import infrapilot

# Initialize client
client = infrapilot.Client(
    base_url="http://localhost:8080/api/v1",
    username="admin@infrapilot.io",
    password="your-password"
)

# Test connection
health = client.health.check()
print(f"Status: {health['status']}")

# List machines
machines = client.machines.list(limit=50)

for machine in machines['data']:
    cpu = machine['metrics']['cpu']
    memory = machine['metrics']['memory']
    print(f"{machine['name']} - CPU: {cpu}%, Memory: {memory}%")

# Get metrics
metrics = client.machines.get_metrics(
    machine_id="machine-uuid",
    period="24h",
    metrics=["cpu", "memory"]
)

for point in metrics['cpu']:
    print(f"{point['timestamp']}: {point['value']}%")

# Create alert
alert_rule = client.alerts.create_rule(
    machine_id="machine-uuid",
    metric="cpu",
    condition=">",
    threshold=90.0,
    severity="warning",
    duration="5m"
)

print(f"Created alert: {alert_rule['id']}")

# Execute command
result = client.machines.execute_command(
    machine_id="machine-uuid",
    command="df -h"
)

print(f"Output:\n{result['output']}")
print(f"Exit code: {result['exit_code']}")

# List files
files = client.files.list(machine_id="machine-uuid", path="/var/log")
for file in files:
    print(f"{file['name']} - {file['type']} - {file['size']} bytes")

# Upload file
client.files.upload(
    machine_id="machine-uuid",
    remote_path="/tmp/test.txt",
    local_path="./test.txt"
)

# Download file
client.files.download(
    machine_id="machine-uuid",
    remote_path="/var/log/syslog",
    local_path="./syslog.log"
)

# Docker operations
containers = client.docker.list(machine_id="machine-uuid")
for container in containers:
    print(f"{container['name']} - {container['status']}")

# Restart container
client.docker.action(
    machine_id="machine-uuid",
    container_id="container-name",
    action="restart"
)
```

### Advanced Usage

```python
import infrapilot
import asyncio

# Context manager
with infrapilot.Client(
    base_url="http://localhost:8080/api/v1",
    username="admin@infrapilot.io",
    password="your-password"
) as client:

    # List all alerts
    alerts = client.alerts.list(status="active")

    # Bulk acknowledge
    for alert in alerts['data']:
        client.alerts.acknowledge(alert['id'])

    # Generate report
    report = client.reports.generate(
        type="performance",
        machine_ids=["uuid1", "uuid2"],
        start_date="2025-01-01T00:00:00Z",
        end_date="2025-01-31T23:59:59Z",
        format="pdf"
    )

    # Download report
    with open("report.pdf", "wb") as f:
        f.write(report['content'])
```

## JavaScript/TypeScript SDK

### Installation

```bash
npm install @infrapilot/sdk
# or
yarn add @infrapilot/sdk
```

### Quick Start

```javascript
import { InfraPilotClient } from '@infrapilot/sdk';

// Initialize client
const client = new InfraPilotClient({
  baseURL: 'http://localhost:8080/api/v1',
  username: 'admin@infrapilot.io',
  password: 'your-password'
});

async function main() {
  // Test connection
  const health = await client.health.check();
  console.log('Status:', health.status);

  // List machines
  const machines = await client.machines.list({ limit: 50 });

  machines.data.forEach(machine => {
    const cpu = machine.metrics.cpu;
    const memory = machine.metrics.memory;
    console.log(`${machine.name} - CPU: ${cpu}%, Memory: ${memory}%`);
  });

  // Get metrics
  const metrics = await client.machines.getMetrics('machine-uuid', {
    period: '24h',
    metrics: ['cpu', 'memory', 'disk']
  });

  metrics.cpu.forEach(point => {
    console.log(`${point.timestamp}: ${point.value}%`);
  });

  // Create alert
  const alertRule = await client.alerts.createRule({
    machine_id: 'machine-uuid',
    metric: 'cpu',
    condition: '>',
    threshold: 90.0,
    severity: 'warning',
    duration: '5m'
  });

  console.log('Created alert:', alertRule.id);

  // Execute command
  const result = await client.machines.executeCommand('machine-uuid', {
    command: 'df -h'
  });

  console.log('Output:\n', result.output);
  console.log('Exit code:', result.exit_code);

  // File operations
  const files = await client.files.list('machine-uuid', { path: '/var/log' });
  files.forEach(file => {
    console.log(`${file.name} - ${file.type} - ${file.size} bytes`);
  });

  // Upload file
  const fs = require('fs');
  const fileContent = fs.readFileSync('./test.txt');
  await client.files.upload('machine-uuid', '/tmp/test.txt', fileContent);

  // Docker operations
  const containers = await client.docker.list('machine-uuid');
  containers.forEach(container => {
    console.log(`${container.name} - ${container.status}`);
  });

  // Restart container
  await client.docker.action('machine-uuid', 'container-name', 'restart');
}

main().catch(console.error);
```

### TypeScript with Type Definitions

```typescript
import { InfraPilotClient, Machine, Alert, Metrics } from '@infrapilot/sdk';

const client = new InfraPilotClient({
  baseURL: 'http://localhost:8080/api/v1',
  username: 'admin@infrapilot.io',
  password: 'your-password'
});

async function monitorMachines(): Promise<void> {
  // List machines with type safety
  const response = await client.machines.list({ limit: 50 });
  const machines: Machine[] = response.data;

  machines.forEach(machine => {
    const cpu: number = machine.metrics.cpu;
    const memory: number = machine.metrics.memory;
    console.log(`${machine.name}: CPU ${cpu}%, Memory ${memory}%`);
  });

  // Create alert rule
  const alertRule: Alert = await client.alerts.createRule({
    machine_id: 'machine-uuid',
    metric: 'cpu',
    condition: '>',
    threshold: 90.0,
    severity: 'warning',
    duration: '5m'
  });

  console.log(`Alert created: ${alertRule.id}`);

  // Get metrics
  const metrics: Metrics = await client.machines.getMetrics('machine-uuid', {
    period: '24h'
  });

  // Access typed metrics
  metrics.cpu.forEach(point => {
    console.log(`${point.timestamp}: ${point.value}%`);
  });
}
```

## REST API Examples

### cURL Examples

#### List Machines

```bash
# Login
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@infrapilot.io","password":"password"}' \
  | jq -r '.access_token')

# List machines
curl -X GET "http://localhost:8080/api/v1/machines?limit=50&status=online" \
  -H "Authorization: Bearer $TOKEN" \
  | jq .
```

#### Get Metrics

```bash
curl -X GET "http://localhost:8080/api/v1/machines/machine-uuid/metrics?period=24h&metrics=cpu,memory" \
  -H "Authorization: Bearer $TOKEN" \
  | jq .
```

#### Execute Command

```bash
curl -X POST "http://localhost:8080/api/v1/machines/machine-uuid/execute" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"command": "df -h"}'
```

### Python with requests

```python
import requests

BASE_URL = "http://localhost:8080/api/v1"

# Login
response = requests.post(f"{BASE_URL}/auth/login", json={
    "email": "admin@infrapilot.io",
    "password": "your-password"
})
data = response.json()
token = data['access_token']

headers = {"Authorization": f"Bearer {token}"}

# List machines
response = requests.get(f"{BASE_URL}/machines", headers=headers)
machines = response.json()

# Get metrics
response = requests.get(
    f"{BASE_URL}/machines/machine-uuid/metrics",
    params={"period": "24h", "metrics": "cpu,memory"},
    headers=headers
)
metrics = response.json()

# Execute command
response = requests.post(
    f"{BASE_URL}/machines/machine-uuid/execute",
    json={"command": "df -h"},
    headers=headers
)
result = response.json()

print(f"Output: {result['output']}")
print(f"Exit code: {result['exit_code']}")
```

## WebSocket API

### Connection

```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = () => {
  console.log('Connected to WebSocket');

  // Authenticate
  ws.send(JSON.stringify({
    type: 'auth',
    token: 'your-jwt-token'
  }));
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  
  switch (message.type) {
    case 'metric':
      console.log('Metric update:', message.data);
      break;
    case 'alert':
      console.log('Alert:', message.data);
      break;
    case 'terminal_output':
      console.log('Terminal:', message.data);
      break;
    case 'error':
      console.error('Error:', message.message);
      break;
  }
};
```

### Subscribe to Metrics

```javascript
ws.onopen = () => {
  // Subscribe to metrics channel
  ws.send(JSON.stringify({
    type: 'subscribe',
    channel: 'metrics',
    machine_id: 'machine-uuid'
  }));
};

// Receive real-time metrics
ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  
  if (message.type === 'metric') {
    console.log('CPU:', message.data.cpu);
    console.log('Memory:', message.data.memory);
    console.log('Disk:', message.data.disk);
  }
};
```

### Terminal Session

```javascript
ws.onopen = () => {
  // Create terminal session
  ws.send(JSON.stringify({
    type: 'create_terminal',
    machine_id: 'machine-uuid',
    shell: '/bin/bash'
  }));
};

// Send input to terminal
function sendInput(input) {
  ws.send(JSON.stringify({
    type: 'terminal_input',
    data: input
  }));
}

// Resize terminal
function resizeTerminal(cols, rows) {
  ws.send(JSON.stringify({
    type: 'terminal_resize',
    cols: cols,
    rows: rows
  }));
}
```

### Python WebSocket Client

```python
import asyncio
import websockets
import json

async def connect():
    uri = "ws://localhost:8080/ws"
    
    async with websockets.connect(uri) as websocket:
        # Authenticate
        await websocket.send(json.dumps({
            "type": "auth",
            "token": "your-jwt-token"
        }))
        
        # Subscribe to metrics
        await websocket.send(json.dumps({
            "type": "subscribe",
            "channel": "metrics",
            "machine_id": "machine-uuid"
        }))
        
        # Receive messages
        while True:
            message = await websocket.recv()
            data = json.loads(message)
            
            if data['type'] == 'metric':
                print(f"CPU: {data['data']['cpu']}%")
                print(f"Memory: {data['data']['memory']}%")
            elif data['type'] == 'alert':
                print(f"Alert: {data['data']}")

asyncio.run(connect())
```

## Best Practices

### 1. Token Management

```javascript
// Store token securely
const token = response.data.access_token;
sessionStorage.setItem('infrapilot_token', token);

// Refresh token before expiry
async function refreshToken() {
  const refreshToken = sessionStorage.getItem('infrapilot_refresh_token');
  
  const response = await fetch('http://localhost:8080/api/v1/auth/refresh', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${refreshToken}`
    }
  });
  
  const data = await response.json();
  sessionStorage.setItem('infrapilot_token', data.access_token);
}

// Auto-refresh before expiry
setInterval(refreshToken, 5 * 60 * 1000); // Refresh every 5 minutes
```

### 2. Error Handling

```javascript
async function listMachines() {
  try {
    const response = await client.machines.list();
    return response.data;
  } catch (error) {
    if (error.response?.status === 401) {
      // Token expired, refresh and retry
      await refreshToken();
      return listMachines();
    }
    
    console.error('Failed to list machines:', error);
    throw error;
  }
}
```

### 3. Rate Limiting

```javascript
class RateLimiter {
  constructor(limit, window) {
    this.limit = limit;
    this.window = window;
    this.requests = [];
  }

  async wait() {
    const now = Date.now();
    this.requests = this.requests.filter(time => now - time < this.window);
    
    if (this.requests.length >= this.limit) {
      const oldestRequest = this.requests[0];
      const waitTime = this.window - (now - oldestRequest);
      await new Promise(resolve => setTimeout(resolve, waitTime));
    }
    
    this.requests.push(now);
  }
}

const rateLimiter = new RateLimiter(100, 60000); // 100 requests per minute

async function makeRequest() {
  await rateLimiter.wait();
  return client.machines.list();
}
```

### 4. Retry Logic

```javascript
async function withRetry(operation, maxRetries = 3) {
  for (let i = 0; i < maxRetries; i++) {
    try {
      return await operation();
    } catch (error) {
      if (i === maxRetries - 1) throw error;
      
      if (error.response?.status >= 500) {
        // Server error, retry with backoff
        await new Promise(resolve => setTimeout(resolve, Math.pow(2, i) * 1000));
        continue;
      }
      
      throw error;
    }
  }
}

// Usage
const machines = await withRetry(() => client.machines.list());
```

### 5. Pagination

```javascript
async function listAllMachines() {
  let allMachines = [];
  let page = 1;
  let hasMore = true;

  while (hasMore) {
    const response = await client.machines.list({ page, limit: 100 });
    allMachines.push(...response.data);
    hasMore = response.pagination.page * response.pagination.limit < response.pagination.total;
    page++;
  }

  return allMachines;
}
```

## Support

- **Documentation**: https://github.com/venkaiswami/infrapilot-enterprise/tree/main/docs
- **API Reference**: `docs/api.md`
- **Issues**: https://github.com/venkaiswami/infrapilot-enterprise/issues
- **Discussions**: https://github.com/venkaiswami/infrapilot-enterprise/discussions

## Contributing

Contributions welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.