CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE organizations (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name TEXT NOT NULL,
  slug TEXT NOT NULL UNIQUE,
  plan TEXT NOT NULL DEFAULT 'enterprise',
  retention_days INTEGER NOT NULL DEFAULT 30 CHECK (retention_days > 0),
  machine_limit INTEGER,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  full_name TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'invited', 'disabled')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE organization_members (
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role TEXT NOT NULL CHECK (role IN ('owner', 'admin', 'operator', 'viewer')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (organization_id, user_id)
);

CREATE TABLE refresh_tokens (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_hash TEXT NOT NULL UNIQUE,
  user_agent TEXT,
  ip_address INET,
  expires_at TIMESTAMPTZ NOT NULL,
  revoked_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE enrollment_tokens (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  created_by UUID REFERENCES users(id) ON DELETE SET NULL,
  name TEXT NOT NULL,
  token_hash TEXT NOT NULL UNIQUE,
  scopes TEXT[] NOT NULL DEFAULT ARRAY['agent:enroll'],
  max_uses INTEGER,
  used_count INTEGER NOT NULL DEFAULT 0,
  expires_at TIMESTAMPTZ,
  revoked_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_enrollment_tokens_org
  ON enrollment_tokens (organization_id, created_at DESC);

CREATE TABLE api_tokens (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
  user_id UUID REFERENCES users(id) ON DELETE SET NULL,
  name TEXT NOT NULL,
  token_hash TEXT NOT NULL UNIQUE,
  scopes TEXT[] NOT NULL DEFAULT '{}',
  expires_at TIMESTAMPTZ,
  revoked_at TIMESTAMPTZ,
  last_used_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE machines (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  hostname TEXT NOT NULL,
  display_name TEXT,
  ip_address INET,
  os TEXT,
  os_version TEXT,
  kernel TEXT,
  architecture TEXT,
  platform TEXT,
  machine_type TEXT NOT NULL DEFAULT 'server',
  environment TEXT,
  region TEXT,
  agent_version TEXT,
  fingerprint TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'online' CHECK (status IN ('online', 'offline', 'degraded', 'maintenance')),
  last_seen_at TIMESTAMPTZ,
  enrolled_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (organization_id, fingerprint)
);

CREATE INDEX idx_machines_org_status ON machines (organization_id, status);
CREATE INDEX idx_machines_org_last_seen ON machines (organization_id, last_seen_at DESC);
CREATE INDEX idx_machines_org_type ON machines (organization_id, machine_type);

CREATE TABLE machine_api_keys (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  machine_id UUID NOT NULL REFERENCES machines(id) ON DELETE CASCADE,
  key_hash TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL DEFAULT 'default',
  last_used_at TIMESTAMPTZ,
  expires_at TIMESTAMPTZ,
  revoked_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_machine_api_keys_machine
  ON machine_api_keys (machine_id, revoked_at);

CREATE TABLE machine_tags (
  machine_id UUID NOT NULL REFERENCES machines(id) ON DELETE CASCADE,
  key TEXT NOT NULL,
  value TEXT NOT NULL,
  PRIMARY KEY (machine_id, key)
);

CREATE TABLE machine_latest_state (
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  machine_id UUID PRIMARY KEY REFERENCES machines(id) ON DELETE CASCADE,
  sampled_at TIMESTAMPTZ NOT NULL,
  cpu_usage DOUBLE PRECISION,
  memory_usage DOUBLE PRECISION,
  disk_usage DOUBLE PRECISION,
  swap_usage DOUBLE PRECISION,
  load_1 DOUBLE PRECISION,
  load_5 DOUBLE PRECISION,
  load_15 DOUBLE PRECISION,
  upload_mbps DOUBLE PRECISION,
  download_mbps DOUBLE PRECISION,
  bandwidth_mbps DOUBLE PRECISION,
  latency_ms DOUBLE PRECISION,
  packet_loss DOUBLE PRECISION,
  disk_read_bps DOUBLE PRECISION,
  disk_write_bps DOUBLE PRECISION,
  iops_read DOUBLE PRECISION,
  iops_write DOUBLE PRECISION,
  cpu_temperature DOUBLE PRECISION,
  uptime_seconds BIGINT,
  raw JSONB NOT NULL DEFAULT '{}'::jsonb
);

CREATE INDEX idx_machine_latest_org_sampled
  ON machine_latest_state (organization_id, sampled_at DESC);

CREATE TABLE metric_samples (
  id BIGSERIAL PRIMARY KEY,
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  machine_id UUID NOT NULL REFERENCES machines(id) ON DELETE CASCADE,
  batch_id UUID,
  sampled_at TIMESTAMPTZ NOT NULL,
  cpu_usage DOUBLE PRECISION,
  memory_usage DOUBLE PRECISION,
  disk_usage DOUBLE PRECISION,
  swap_usage DOUBLE PRECISION,
  load_1 DOUBLE PRECISION,
  load_5 DOUBLE PRECISION,
  load_15 DOUBLE PRECISION,
  upload_mbps DOUBLE PRECISION,
  download_mbps DOUBLE PRECISION,
  bandwidth_mbps DOUBLE PRECISION,
  latency_ms DOUBLE PRECISION,
  packet_loss DOUBLE PRECISION,
  disk_read_bps DOUBLE PRECISION,
  disk_write_bps DOUBLE PRECISION,
  iops_read DOUBLE PRECISION,
  iops_write DOUBLE PRECISION,
  cpu_temperature DOUBLE PRECISION,
  uptime_seconds BIGINT,
  raw JSONB NOT NULL DEFAULT '{}'::jsonb,
  received_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_metric_samples_org_machine_time
  ON metric_samples (organization_id, machine_id, sampled_at DESC);

CREATE INDEX idx_metric_samples_org_time
  ON metric_samples (organization_id, sampled_at DESC);

CREATE TABLE metric_rollups_1m (
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  machine_id UUID NOT NULL REFERENCES machines(id) ON DELETE CASCADE,
  bucket_at TIMESTAMPTZ NOT NULL,
  avg_cpu DOUBLE PRECISION,
  max_cpu DOUBLE PRECISION,
  avg_memory DOUBLE PRECISION,
  max_memory DOUBLE PRECISION,
  avg_disk DOUBLE PRECISION,
  max_disk DOUBLE PRECISION,
  avg_upload_mbps DOUBLE PRECISION,
  avg_download_mbps DOUBLE PRECISION,
  avg_disk_read_bps DOUBLE PRECISION,
  avg_disk_write_bps DOUBLE PRECISION,
  sample_count INTEGER NOT NULL DEFAULT 0,
  PRIMARY KEY (organization_id, machine_id, bucket_at)
);

CREATE TABLE metric_rollups_1h (
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  machine_id UUID NOT NULL REFERENCES machines(id) ON DELETE CASCADE,
  bucket_at TIMESTAMPTZ NOT NULL,
  avg_cpu DOUBLE PRECISION,
  max_cpu DOUBLE PRECISION,
  avg_memory DOUBLE PRECISION,
  max_memory DOUBLE PRECISION,
  avg_disk DOUBLE PRECISION,
  max_disk DOUBLE PRECISION,
  avg_upload_mbps DOUBLE PRECISION,
  avg_download_mbps DOUBLE PRECISION,
  avg_disk_read_bps DOUBLE PRECISION,
  avg_disk_write_bps DOUBLE PRECISION,
  sample_count INTEGER NOT NULL DEFAULT 0,
  PRIMARY KEY (organization_id, machine_id, bucket_at)
);

CREATE TABLE alert_rules (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  metric TEXT NOT NULL,
  operator TEXT NOT NULL CHECK (operator IN ('>', '>=', '<', '<=', '=', '!=')),
  threshold DOUBLE PRECISION,
  duration_seconds INTEGER NOT NULL DEFAULT 60,
  severity TEXT NOT NULL CHECK (severity IN ('info', 'warning', 'critical')),
  scope JSONB NOT NULL DEFAULT '{}'::jsonb,
  enabled BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_alert_rules_org_enabled
  ON alert_rules (organization_id, enabled);

CREATE TABLE alerts (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  machine_id UUID REFERENCES machines(id) ON DELETE CASCADE,
  rule_id UUID REFERENCES alert_rules(id) ON DELETE SET NULL,
  type TEXT NOT NULL,
  severity TEXT NOT NULL CHECK (severity IN ('info', 'warning', 'critical')),
  status TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'acknowledged', 'resolved')),
  title TEXT NOT NULL,
  message TEXT NOT NULL,
  source TEXT NOT NULL DEFAULT 'rule',
  opened_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  acknowledged_at TIMESTAMPTZ,
  acknowledged_by UUID REFERENCES users(id) ON DELETE SET NULL,
  resolved_at TIMESTAMPTZ,
  metadata JSONB NOT NULL DEFAULT '{}'::jsonb
);

CREATE INDEX idx_alerts_org_status
  ON alerts (organization_id, status, opened_at DESC);

CREATE INDEX idx_alerts_machine_status
  ON alerts (organization_id, machine_id, status, opened_at DESC);

CREATE TABLE alert_events (
  id BIGSERIAL PRIMARY KEY,
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  alert_id UUID NOT NULL REFERENCES alerts(id) ON DELETE CASCADE,
  event_type TEXT NOT NULL,
  message TEXT,
  created_by UUID REFERENCES users(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  metadata JSONB NOT NULL DEFAULT '{}'::jsonb
);

CREATE TABLE logs (
  id BIGSERIAL PRIMARY KEY,
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  machine_id UUID NOT NULL REFERENCES machines(id) ON DELETE CASCADE,
  collected_at TIMESTAMPTZ NOT NULL,
  source TEXT NOT NULL,
  level TEXT,
  message TEXT NOT NULL,
  metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
  received_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_logs_org_machine_time
  ON logs (organization_id, machine_id, collected_at DESC);

CREATE TABLE process_snapshots (
  id BIGSERIAL PRIMARY KEY,
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  machine_id UUID NOT NULL REFERENCES machines(id) ON DELETE CASCADE,
  sampled_at TIMESTAMPTZ NOT NULL,
  pid INTEGER NOT NULL,
  name TEXT NOT NULL,
  command TEXT,
  cpu_usage DOUBLE PRECISION,
  memory_bytes BIGINT,
  user_name TEXT
);

CREATE INDEX idx_process_snapshots_machine_time
  ON process_snapshots (organization_id, machine_id, sampled_at DESC);

CREATE TABLE filesystems (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  machine_id UUID NOT NULL REFERENCES machines(id) ON DELETE CASCADE,
  mount_point TEXT NOT NULL,
  filesystem_type TEXT,
  total_bytes BIGINT,
  used_bytes BIGINT,
  usage_percent DOUBLE PRECISION,
  last_seen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (machine_id, mount_point)
);

CREATE TABLE storage_devices (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  machine_id UUID NOT NULL REFERENCES machines(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  device_type TEXT,
  size_bytes BIGINT,
  health TEXT,
  smart_status TEXT,
  raid_status TEXT,
  lvm_group TEXT,
  read_bps DOUBLE PRECISION,
  write_bps DOUBLE PRECISION,
  iops_read DOUBLE PRECISION,
  iops_write DOUBLE PRECISION,
  last_seen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  raw JSONB NOT NULL DEFAULT '{}'::jsonb,
  UNIQUE (machine_id, name)
);

CREATE TABLE network_interfaces (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  machine_id UUID NOT NULL REFERENCES machines(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  mac_address TEXT,
  ip_addresses TEXT[] NOT NULL DEFAULT '{}',
  upload_mbps DOUBLE PRECISION,
  download_mbps DOUBLE PRECISION,
  packet_loss DOUBLE PRECISION,
  latency_ms DOUBLE PRECISION,
  last_seen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (machine_id, name)
);

CREATE TABLE network_connections (
  id BIGSERIAL PRIMARY KEY,
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  machine_id UUID NOT NULL REFERENCES machines(id) ON DELETE CASCADE,
  sampled_at TIMESTAMPTZ NOT NULL,
  protocol TEXT,
  local_address TEXT,
  local_port INTEGER,
  remote_address TEXT,
  remote_port INTEGER,
  state TEXT,
  process_name TEXT,
  pid INTEGER
);

CREATE TABLE open_ports (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  machine_id UUID NOT NULL REFERENCES machines(id) ON DELETE CASCADE,
  protocol TEXT NOT NULL,
  port INTEGER NOT NULL,
  process_name TEXT,
  pid INTEGER,
  last_seen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (machine_id, protocol, port)
);

CREATE TABLE docker_containers (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  machine_id UUID NOT NULL REFERENCES machines(id) ON DELETE CASCADE,
  container_id TEXT NOT NULL,
  name TEXT NOT NULL,
  image TEXT,
  status TEXT NOT NULL,
  cpu_usage DOUBLE PRECISION,
  memory_usage DOUBLE PRECISION,
  restart_count INTEGER NOT NULL DEFAULT 0,
  last_seen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  raw JSONB NOT NULL DEFAULT '{}'::jsonb,
  UNIQUE (machine_id, container_id)
);

CREATE TABLE docker_images (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  machine_id UUID NOT NULL REFERENCES machines(id) ON DELETE CASCADE,
  image_id TEXT NOT NULL,
  repository TEXT,
  tag TEXT,
  size_bytes BIGINT,
  created_at_source TIMESTAMPTZ,
  last_seen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (machine_id, image_id)
);

CREATE TABLE docker_volumes (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  machine_id UUID NOT NULL REFERENCES machines(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  driver TEXT,
  mountpoint TEXT,
  last_seen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  raw JSONB NOT NULL DEFAULT '{}'::jsonb,
  UNIQUE (machine_id, name)
);

CREATE TABLE kubernetes_clusters (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  version TEXT,
  status TEXT NOT NULL DEFAULT 'unknown',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (organization_id, name)
);

CREATE TABLE kubernetes_nodes (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  cluster_id UUID NOT NULL REFERENCES kubernetes_clusters(id) ON DELETE CASCADE,
  machine_id UUID REFERENCES machines(id) ON DELETE SET NULL,
  name TEXT NOT NULL,
  status TEXT NOT NULL,
  roles TEXT[] NOT NULL DEFAULT '{}',
  cpu_capacity DOUBLE PRECISION,
  memory_capacity BIGINT,
  last_seen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  raw JSONB NOT NULL DEFAULT '{}'::jsonb,
  UNIQUE (cluster_id, name)
);

CREATE TABLE kubernetes_pods (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  cluster_id UUID NOT NULL REFERENCES kubernetes_clusters(id) ON DELETE CASCADE,
  node_id UUID REFERENCES kubernetes_nodes(id) ON DELETE SET NULL,
  namespace TEXT NOT NULL,
  name TEXT NOT NULL,
  status TEXT NOT NULL,
  restarts INTEGER NOT NULL DEFAULT 0,
  cpu_usage DOUBLE PRECISION,
  memory_usage BIGINT,
  last_seen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  raw JSONB NOT NULL DEFAULT '{}'::jsonb,
  UNIQUE (cluster_id, namespace, name)
);

CREATE TABLE kubernetes_workloads (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  cluster_id UUID NOT NULL REFERENCES kubernetes_clusters(id) ON DELETE CASCADE,
  namespace TEXT NOT NULL,
  kind TEXT NOT NULL,
  name TEXT NOT NULL,
  desired_replicas INTEGER,
  ready_replicas INTEGER,
  status TEXT,
  last_seen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  raw JSONB NOT NULL DEFAULT '{}'::jsonb,
  UNIQUE (cluster_id, namespace, kind, name)
);

CREATE TABLE kubernetes_services (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  cluster_id UUID NOT NULL REFERENCES kubernetes_clusters(id) ON DELETE CASCADE,
  namespace TEXT NOT NULL,
  name TEXT NOT NULL,
  service_type TEXT,
  cluster_ip TEXT,
  ports JSONB NOT NULL DEFAULT '[]'::jsonb,
  last_seen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (cluster_id, namespace, name)
);

CREATE TABLE kubernetes_volumes (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  cluster_id UUID NOT NULL REFERENCES kubernetes_clusters(id) ON DELETE CASCADE,
  namespace TEXT,
  name TEXT NOT NULL,
  volume_type TEXT NOT NULL,
  status TEXT,
  capacity_bytes BIGINT,
  last_seen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  raw JSONB NOT NULL DEFAULT '{}'::jsonb,
  UNIQUE (cluster_id, namespace, name, volume_type)
);

CREATE TABLE audit_logs (
  id BIGSERIAL PRIMARY KEY,
  organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
  actor_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
  action TEXT NOT NULL,
  resource_type TEXT NOT NULL,
  resource_id TEXT,
  ip_address INET,
  user_agent TEXT,
  metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_audit_logs_org_time
  ON audit_logs (organization_id, created_at DESC);
