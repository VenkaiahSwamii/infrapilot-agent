# API and WebSocket Design

Base path: `/api/v1`

Authentication:

- User REST APIs: `Authorization: Bearer <jwt>`.
- WebSocket APIs: JWT in `Authorization` header where possible, otherwise a short-lived socket token.
- Agent APIs: enrollment token for first enrollment, then `X-Machine-API-Key`.

Every tenant resource must be scoped by `organization_id`. Handlers should never trust a client-provided organization ID without checking membership or machine ownership.

## Auth Endpoints

| Method | Path | Purpose |
| --- | --- | --- |
| `POST` | `/auth/register` | Create first user or invite-accepted user. |
| `POST` | `/auth/login` | Return access JWT and refresh token. |
| `POST` | `/auth/refresh` | Rotate refresh token and return new JWT. |
| `POST` | `/auth/logout` | Revoke refresh token. |
| `GET` | `/auth/me` | Return user profile, organizations, and roles. |

## Organization Endpoints

| Method | Path | Purpose |
| --- | --- | --- |
| `GET` | `/organizations` | List organizations visible to the user. |
| `POST` | `/organizations` | Create organization. |
| `GET` | `/organizations/:orgId` | Get organization details. |
| `PATCH` | `/organizations/:orgId` | Update organization settings. |
| `GET` | `/organizations/:orgId/overview` | Aggregate live dashboard summary. |
| `GET` | `/organizations/:orgId/users` | List members. |
| `POST` | `/organizations/:orgId/users` | Invite/add member. |
| `PATCH` | `/organizations/:orgId/users/:userId` | Change role or status. |

## Enrollment and Agent Endpoints

| Method | Path | Auth | Purpose |
| --- | --- | --- | --- |
| `POST` | `/organizations/:orgId/enrollment-tokens` | JWT admin | Create enrollment token. |
| `GET` | `/organizations/:orgId/enrollment-tokens` | JWT admin | List tokens without raw secret. |
| `DELETE` | `/organizations/:orgId/enrollment-tokens/:tokenId` | JWT admin | Revoke token. |
| `POST` | `/agent/enroll` | Enrollment token | Register or re-enroll machine. |
| `POST` | `/agent/heartbeat` | Machine key | Update last seen and status. |
| `POST` | `/agent/metrics` | Machine key | Submit metric batch. |
| `POST` | `/agent/logs` | Machine key | Submit log batch. |
| `GET` | `/agent/config` | Machine key | Pull remote config. |
| `POST` | `/agent/key-rotation` | Machine key | Rotate machine API key. |

Enrollment request:

```json
{
  "enrollment_token": "ip_enroll_xxx",
  "hostname": "db-prod-01",
  "fingerprint": "sha256:...",
  "os": "linux",
  "os_version": "Ubuntu 24.04",
  "kernel": "6.8.0",
  "architecture": "amd64",
  "machine_type": "server",
  "agent_version": "1.0.0",
  "interfaces": [
    { "name": "eth0", "mac_address": "00:11:22:33:44:55", "ip_addresses": ["10.0.2.15"] }
  ]
}
```

Enrollment response:

```json
{
  "machine_id": "3a5347e4-7c17-4ae1-9b0a-0e8ee8bd7c2b",
  "organization_id": "9428495a-099a-4be7-9313-8736da18d734",
  "machine_api_key": "ip_machine_xxx",
  "config": {
    "metrics_interval_seconds": 5,
    "heartbeat_interval_seconds": 15,
    "log_collection_enabled": true
  }
}
```

Metrics batch request:

```json
{
  "batch_id": "ca0f63a3-d082-4421-ad32-272a191dbdb2",
  "schema_version": "2026-07-04",
  "sent_at": "2026-07-04T08:00:00Z",
  "samples": [
    {
      "sampled_at": "2026-07-04T07:59:58Z",
      "cpu_usage": 64.2,
      "memory_usage": 72.1,
      "disk_usage": 81.4,
      "swap_usage": 8.2,
      "load": { "one": 2.1, "five": 1.8, "fifteen": 1.5 },
      "network": { "upload_mbps": 21.5, "download_mbps": 48.9, "latency_ms": 12.4, "packet_loss": 0.1 },
      "disk_io": { "read_bps": 2400000, "write_bps": 880000, "iops_read": 120, "iops_write": 44 },
      "uptime_seconds": 928334
    }
  ],
  "inventory": {
    "filesystems": [],
    "network_interfaces": [],
    "docker_containers": [],
    "kubernetes_pods": [],
    "storage_devices": []
  }
}
```

## Machine and Analytics Endpoints

| Method | Path | Purpose |
| --- | --- | --- |
| `GET` | `/organizations/:orgId/machines` | Search/filter machine inventory. |
| `GET` | `/organizations/:orgId/machines/:machineId` | Machine profile and latest state. |
| `PATCH` | `/organizations/:orgId/machines/:machineId` | Rename, tag, maintenance mode. |
| `GET` | `/organizations/:orgId/machines/:machineId/metrics` | Historical CPU/RAM/storage/network metrics. |
| `GET` | `/organizations/:orgId/machines/:machineId/processes` | Latest/top process snapshots. |
| `GET` | `/organizations/:orgId/machines/:machineId/filesystems` | Filesystem usage. |
| `GET` | `/organizations/:orgId/machines/:machineId/network` | Interfaces, open ports, connections. |
| `GET` | `/organizations/:orgId/machines/:machineId/storage` | SMART, RAID/LVM, disk IO, IOPS. |
| `GET` | `/organizations/:orgId/machines/:machineId/docker` | Containers, images, volumes, networks. |
| `GET` | `/organizations/:orgId/machines/:machineId/kubernetes` | Nodes, pods, workloads, services, PV/PVC. |
| `GET` | `/organizations/:orgId/machines/:machineId/logs` | Filtered machine logs. |

Metrics query parameters:

- `range`: `1h`, `24h`, `7d`, `30d`.
- `resolution`: `raw`, `1m`, `1h`, `auto`.
- `metrics`: comma-separated fields, for example `cpu_usage,memory_usage,upload_mbps`.

## Alert Endpoints

| Method | Path | Purpose |
| --- | --- | --- |
| `GET` | `/organizations/:orgId/alerts` | List alerts by status, severity, machine, type. |
| `POST` | `/organizations/:orgId/alerts/:alertId/ack` | Acknowledge alert. |
| `POST` | `/organizations/:orgId/alerts/:alertId/resolve` | Manually resolve alert. |
| `GET` | `/organizations/:orgId/alert-rules` | List rules. |
| `POST` | `/organizations/:orgId/alert-rules` | Create rule. |
| `PATCH` | `/organizations/:orgId/alert-rules/:ruleId` | Update rule. |
| `DELETE` | `/organizations/:orgId/alert-rules/:ruleId` | Disable/delete rule. |

## WebSocket Endpoint

Preferred endpoint:

```text
GET /ws/v1/organizations/:orgId
Authorization: Bearer <jwt>
```

Fallback for browser limitations:

```text
GET /ws/v1/organizations/:orgId?socket_token=<short_lived_token>
```

Connection rules:

- Validate JWT or socket token.
- Verify user membership in `orgId`.
- Join an organization room.
- Allow optional machine subscriptions after connect.
- Send ping/pong heartbeats every 25-30 seconds.
- Drop slow consumers after bounded buffering is exceeded.

Client subscribe message:

```json
{
  "type": "subscribe",
  "channels": [
    "organization.overview",
    "alerts",
    "machine.3a5347e4-7c17-4ae1-9b0a-0e8ee8bd7c2b"
  ]
}
```

Server event envelope:

```json
{
  "type": "metric.snapshot",
  "organization_id": "9428495a-099a-4be7-9313-8736da18d734",
  "machine_id": "3a5347e4-7c17-4ae1-9b0a-0e8ee8bd7c2b",
  "timestamp": "2026-07-04T08:00:00Z",
  "payload": {}
}
```

Core events:

- `organization.overview`
- `machine.registered`
- `machine.online`
- `machine.offline`
- `metric.snapshot`
- `metric.rollup`
- `alert.created`
- `alert.acknowledged`
- `alert.resolved`
- `docker.container_changed`
- `kubernetes.pod_changed`
- `kubernetes.node_changed`
- `storage.health_changed`

Scaling guidance:

- Use Redis/NATS pub-sub when multiple API replicas host WebSocket connections.
- Workers publish events to the bus; API replicas fan out only to locally connected sockets.
- Coalesce metric updates per organization and machine before broadcasting.
- Use REST for initial state and backfill after reconnect.
