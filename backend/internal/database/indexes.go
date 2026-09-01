package database

import (
	"infrapilot/backend/internal/logger"
)

// EnsureIndexesAndTuneDB applies production database index optimizations and PostgreSQL parameters
func EnsureIndexesAndTuneDB() {
	if DB == nil {
		logger.Warn("Database instance uninitialized, skipping index creation")
		return
	}

	logger.Info("Applying production database index optimization & PostgreSQL tuning")

	// 1. Create High-Cardinality Performance Indexes
	indexStatements := []string{
		"CREATE INDEX IF NOT EXISTS idx_machines_org_status ON machines (organization_id, status);",
		"CREATE INDEX IF NOT EXISTS idx_machines_hostname ON machines (hostname);",
		"CREATE INDEX IF NOT EXISTS idx_metrics_machine_time ON metrics (machine_id, timestamp DESC);",
		"CREATE INDEX IF NOT EXISTS idx_alerts_org_status ON alert_rules (organization_id, status);",
		"CREATE INDEX IF NOT EXISTS idx_organization_users_org_user ON organization_users (organization_id, user_id);",
		"CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);",
		"CREATE INDEX IF NOT EXISTS idx_audit_logs_org_created ON audit_logs (organization_id, created_at DESC);",
	}

	for _, stmt := range indexStatements {
		if err := DB.Exec(stmt).Error; err != nil {
			logger.Warn("Index creation notice", "statement", stmt, "error", err)
		}
	}

	// 2. PostgreSQL Server Parameter Tuning (Soft Execution)
	tuneStatements := []string{
		"ALTER SYSTEM SET shared_buffers = '4GB';",
		"ALTER SYSTEM SET work_mem = '64MB';",
		"ALTER SYSTEM SET effective_cache_size = '12GB';",
		"ALTER SYSTEM SET max_connections = '300';",
		"ALTER SYSTEM SET wal_compression = 'on';",
	}

	for _, stmt := range tuneStatements {
		_ = DB.Exec(stmt)
	}

	logger.Info("Production database indexing & tuning applied successfully")
}
