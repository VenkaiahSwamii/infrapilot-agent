-- InfraDeploy: CI/CD Pipeline Tables
-- Part of the InfraPilot Enterprise ecosystem

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Git Repositories table
CREATE TABLE git_repositories (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name TEXT NOT NULL,
  provider TEXT NOT NULL CHECK (provider IN ('github', 'gitlab', 'bitbucket', 'gitea', 'azure-devops')),
  url TEXT NOT NULL,
  clone_url TEXT,
  branch TEXT NOT NULL DEFAULT 'main',
  auth_type TEXT NOT NULL CHECK (auth_type IN ('oauth', 'ssh', 'token', 'basic')),
  auth_token TEXT,
  ssh_key TEXT,
  webhook_url TEXT,
  webhook_token TEXT,
  organization TEXT NOT NULL DEFAULT 'Default Organization',
  created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  last_sync_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_git_repositories_org
  ON git_repositories (organization, created_at DESC);

-- Pipelines table
CREATE TABLE pipelines (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name TEXT NOT NULL,
  description TEXT,
  repository_id UUID NOT NULL REFERENCES git_repositories(id) ON DELETE CASCADE,
  config JSONB NOT NULL DEFAULT '{}'::jsonb,
  trigger_type TEXT NOT NULL DEFAULT 'webhook' CHECK (trigger_type IN ('webhook', 'schedule', 'manual', 'api')),
  schedule TEXT,
  enabled BOOLEAN NOT NULL DEFAULT true,
  environment TEXT NOT NULL DEFAULT 'development' CHECK (environment IN ('development', 'staging', 'production', 'canary')),
  created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  last_run_at TIMESTAMPTZ,
  organization TEXT NOT NULL DEFAULT 'Default Organization',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_pipelines_repo
  ON pipelines (repository_id, created_at DESC);

CREATE INDEX idx_pipelines_org
  ON pipelines (organization, environment, created_at DESC);

-- Builds table
CREATE TABLE builds (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  pipeline_id UUID NOT NULL REFERENCES pipelines(id) ON DELETE CASCADE,
  build_number INTEGER NOT NULL,
  commit_hash TEXT NOT NULL,
  commit_message TEXT,
  author TEXT,
  branch TEXT NOT NULL,
  tag TEXT,
  status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'running', 'success', 'failed', 'cancelled')),
  trigger TEXT NOT NULL CHECK (trigger IN ('manual', 'webhook', 'schedule', 'api')),
  started_at TIMESTAMPTZ,
  finished_at TIMESTAMPTZ,
  duration_ms BIGINT DEFAULT 0,
  logs TEXT,
  artifacts JSONB DEFAULT '[]'::jsonb,
  error_message TEXT,
  organization TEXT NOT NULL DEFAULT 'Default Organization',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_builds_pipeline
  ON builds (pipeline_id, created_at DESC);

CREATE INDEX idx_builds_org_status
  ON builds (organization, status, created_at DESC);

CREATE INDEX idx_builds_commit
  ON builds (commit_hash);

-- Deployments table
CREATE TABLE deployments (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name TEXT NOT NULL,
  environment TEXT NOT NULL CHECK (environment IN ('development', 'staging', 'production', 'canary')),
  namespace TEXT,
  cluster TEXT,
  strategy TEXT NOT NULL DEFAULT 'rolling' CHECK (strategy IN ('rolling', 'blue-green', 'canary', 'recreate')),
  build_id UUID REFERENCES builds(id) ON DELETE SET NULL,
  image TEXT NOT NULL,
  replicas INTEGER NOT NULL DEFAULT 1,
  status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'deploying', 'running', 'failed', 'rolled-back')),
  health_check TEXT,
  auto_rollback BOOLEAN NOT NULL DEFAULT true,
  config JSONB DEFAULT '{}'::jsonb,
  variables JSONB DEFAULT '{}'::jsonb,
  deployed_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  deployed_at TIMESTAMPTZ,
  rolled_back_at TIMESTAMPTZ,
  error_message TEXT,
  organization TEXT NOT NULL DEFAULT 'Default Organization',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_deployments_org_env
  ON deployments (organization, environment, created_at DESC);

CREATE INDEX idx_deployments_build
  ON deployments (build_id);

-- Deployment History table
CREATE TABLE deployment_history (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  deployment_id UUID NOT NULL REFERENCES deployments(id) ON DELETE CASCADE,
  action TEXT NOT NULL CHECK (action IN ('deploy', 'rollback', 'scale', 'restart', 'promote')),
  from_image TEXT,
  to_image TEXT,
  from_replicas INTEGER,
  to_replicas INTEGER,
  reason TEXT,
  triggered_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  triggered_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  organization TEXT NOT NULL DEFAULT 'Default Organization'
);

CREATE INDEX idx_deployment_history_deployment
  ON deployment_history (deployment_id, triggered_at DESC);

-- Webhook Events table
CREATE TABLE webhook_events (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  provider TEXT NOT NULL CHECK (provider IN ('github', 'gitlab', 'bitbucket', 'gitea', 'azure-devops')),
  event_type TEXT NOT NULL CHECK (event_type IN ('push', 'pull_request', 'tag', 'release', 'ping')),
  signature TEXT,
  payload JSONB NOT NULL DEFAULT '{}'::jsonb,
  processed BOOLEAN NOT NULL DEFAULT false,
  error TEXT,
  received_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  processed_at TIMESTAMPTZ,
  repository_id UUID REFERENCES git_repositories(id) ON DELETE SET NULL
);

CREATE INDEX idx_webhook_events_repo
  ON webhook_events (repository_id, received_at DESC);

CREATE INDEX idx_webhook_events_processed
  ON webhook_events (processed, received_at DESC) WHERE processed = false;

-- Deployment Environments table
CREATE TABLE deployment_environments (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name TEXT NOT NULL,
  slug TEXT NOT NULL,
  organization TEXT NOT NULL DEFAULT 'Default Organization',
  cluster TEXT NOT NULL,
  namespace TEXT NOT NULL,
  auto_deploy BOOLEAN NOT NULL DEFAULT false,
  requires_approval BOOLEAN NOT NULL DEFAULT false,
  allowed_branches TEXT[] DEFAULT '{}',
  created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (organization, slug)
);

CREATE INDEX idx_deployment_envs_org
  ON deployment_environments (organization, name);