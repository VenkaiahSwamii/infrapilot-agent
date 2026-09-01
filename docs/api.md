# InfraPilot Enterprise API Documentation

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

All API endpoints (except `/auth/*`) require a Bearer token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Endpoints

### Authentication

#### POST /auth/login

Login with email and password.

**Request:**
```json
{
  "email": "admin@infrapilot.io",
  "password": "secure-password"
}
```

**Response:**
```json
{
  "access_token": "eyJ...",
  "refresh_token": "eyJ...",
  "user": {
    "id": "uuid",
    "email": "admin@infrapilot.io",
    "role": "admin"
  }
}
```

#### POST /auth/refresh

Refresh access token using refresh token.

**Headers:**
```
Authorization: Bearer <refresh-token>
```

**Response:**
```json
{
  "access_token": "eyJ..."
}
```

#### POST /auth/logout

Invalidate current session.

### Machines

#### GET /machines

List all registered machines.

**Query Parameters:**
- `page` (int): Page number (default: 1)
- `limit` (int): Items per page (default: 20)
- `search` (string): Search by name or IP
- `status` (string): Filter by status (online, offline, warning)

**Response:**
```json
{
  "machines": [
    {
      "id": "uuid",
      "name": "web-server-01",
      "ip_address": "192.168.1.100",
      "os": "linux",
      "status": "online",
      "last_seen": "2025-01-01T12:00:00Z",
      "metrics": {
        "cpu": 45.2,
        "memory": 62.1,
        "disk": 34.5
      }
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100
  }
}
```

#### GET /machines/:id

Get machine details.

**Response:**
```json
{
  "id": "uuid",
  "name": "web-server-01",
  "ip_address": "192.168.1.100",
  "os": "linux",
  "status": "online",
  "metadata": {
    "cpu_cores": 8,
    "memory_total": 16384,
    "disk_total": 500
  },
  "last_seen": "2025-01-01T12:00:00Z"
}
```

#### DELETE /machines/:id

Unregister a machine.

### Metrics

#### GET /machines/:id/metrics

Get current metrics for a machine.

**Query Parameters:**
- `metric` (string): Specific metric (cpu, memory, disk, network)
- `period` (string): Time period (1h, 6h, 24h, 7d, 30d)

**Response:**
```json
{
  "machine_id": "uuid",
  "metrics": {
    "cpu": [
      {"timestamp": "2025-01-01T12:00:00Z", "value": 45.2},
      {"timestamp": "2025-01-01T11:00:00Z", "value": 42.1}
    ],
    "memory": [
      {"timestamp": "2025-01-01T12:00:00Z", "value": 62.1}
    ]
  }
}
```

### Agents

#### GET /agents

List all connected agents.

#### GET /agents/:id

Get agent details.

#### POST /agents/:id/actions

Execute action on agent.

**Request:**
```json
{
  "action": "restart_service",
  "params": {
    "service": "nginx"
  }
}
```

### Terminal

#### POST /terminal/sessions

Create terminal session.

**Request:**
```json
{
  "machine_id": "uuid",
  "shell": "/bin/bash"
}
```

**Response:**
```json
{
  "session_id": "uuid",
  "ws_url": "ws://localhost:8080/ws/terminal/session_id"
}
```

#### WebSocket /ws/terminal/:session_id

WebSocket connection for terminal I/O.

**Messages (Client → Server):**
```json
{"type": "input", "data": "ls -la"}
{"type": "resize", "cols": 120, "rows": 40}
```

**Messages (Server → Client):**
```json
{"type": "output", "data": "total 48\n"}
{"type": "error", "data": "Permission denied"}
```

### Files

#### GET /machines/:id/files

List files in directory.

**Query Parameters:**
- `path` (string): Directory path (default: /)

**Response:**
```json
{
  "files": [
    {
      "name": "Documents",
      "path": "/home/user/Documents",
      "type": "directory",
      "size": 4096,
      "modified": "2025-01-01T10:00:00Z"
    }
  ]
}
```

#### GET /machines/:id/files/content

Get file content.

**Query Parameters:**
- `path` (string): File path

**Response:**
```json
{
  "content": "file contents...",
  "encoding": "utf-8"
}
```

### Logs

#### GET /machines/:id/logs

Get machine logs.

**Query Parameters:**
- `service` (string): Filter by service
- `level` (string): Filter by level (info, warn, error)
- `since` (string): ISO timestamp
- `limit` (int): Max log lines (default: 100)

**Response:**
```json
{
  "logs": [
    {
      "timestamp": "2025-01-01T12:00:00Z",
      "level": "info",
      "service": "nginx",
      "message": "Server started"
    }
  ]
}
```

### Alerts

#### GET /alerts

List alerts.

**Query Parameters:**
- `status` (string): active, resolved, acknowledged
- `severity` (string): critical, warning, info

**Response:**
```json
{
  "alerts": [
    {
      "id": "uuid",
      "machine_id": "uuid",
      "name": "High CPU Usage",
      "severity": "warning",
      "status": "active",
      "message": "CPU usage exceeded 90%",
      "created_at": "2025-01-01T12:00:00Z"
    }
  ]
}
```

#### POST /alerts/:id/acknowledge

Acknowledge alert.

#### POST /alerts/:id/resolve

Resolve alert.

### Reports

#### GET /reports

List generated reports.

#### POST /reports/generate

Generate new report.

**Request:**
```json
{
  "type": "performance",
  "machine_ids": ["uuid"],
  "start_date": "2025-01-01T00:00:00Z",
  "end_date": "2025-01-01T23:59:59Z",
  "format": "pdf"
}
```

#### GET /reports/:id/download

Download report file.

### Docker

#### GET /machines/:id/docker/containers

List Docker containers.

**Response:**
```json
{
  "containers": [
    {
      "id": "abc123",
      "name": "web-server",
      "image": "nginx:latest",
      "status": "running",
      "ports": ["80:80", "443:443"],
      "cpu": 0.5,
      "memory": 256
    }
  ]
}
```

#### POST /machines/:id/docker/containers/:id/actions

Execute Docker action.

**Request:**
```json
{
  "action": "restart"
}
```

### Kubernetes

#### GET /machines/:id/kubernetes/clusters

List Kubernetes clusters.

#### GET /machines/:id/kubernetes/clusters/:cluster/pods

List pods in cluster.

## Error Responses

All errors follow this format:

```json
{
  "error": "ERROR_CODE",
  "message": "Human-readable error message",
  "details": {}
}
```

### HTTP Status Codes

- `200 OK` - Success
- `201 Created` - Resource created
- `400 Bad Request` - Invalid request
- `401 Unauthorized` - Missing or invalid auth
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error
- `503 Service Unavailable` - Service temporarily unavailable

## Rate Limiting

API is rate limited to:
- 100 requests per minute for authenticated users
- 10 requests per minute for unauthenticated endpoints

Rate limit headers are included in responses:

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1704067200
```

## WebSocket Events

### Client Events Send

```json
{"type": "subscribe", "channel": "metrics"}
{"type": "unsubscribe", "channel": "metrics"}
{"type": "terminal_input", "data": "ls -la"}
```

### Server Events Receive

```json
{"type": "metric", "data": {"cpu": 45.2, "memory": 62.1}}
{"type": "alert", "data": {"id": "uuid", "severity": "critical"}}
{"type": "terminal_output", "data": "file listing"}
{"type": "error", "message": "Connection lost"}