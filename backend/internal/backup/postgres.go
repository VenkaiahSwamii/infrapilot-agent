package backup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"infrapilot/backend/internal/logger"
)

// PostgresBackup handles PostgreSQL-specific backup operations
type PostgresBackup struct {
	logger *logger.Logger
}

// NewPostgresBackup creates a new PostgreSQL backup handler
func NewPostgresBackup() *PostgresBackup {
	return &PostgresBackup{
		logger: logger.Get(),
	}
}

// getEnvOrDefault retrieves an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Backup performs a PostgreSQL backup (full, schema, data)
func (p *PostgresBackup) Backup(backupType string) (string, error) {
	timestamp := time.Now().Format("2006-01-02-150405")
	backupDir := "./backups/postgres"

	// Ensure directory exists
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	var backupFile string
	var cmd *exec.Cmd

	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "5432")
	dbUser := getEnvOrDefault("DB_USER", "infrapilot")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "infrapilot")
	dbName := getEnvOrDefault("DB_NAME", "infrapilot")

	switch backupType {
	case "full", "":
		backupFile = filepath.Join(backupDir, fmt.Sprintf("%s.sql.gz", timestamp))
		cmd = exec.Command("pg_dump",
			"-h", dbHost,
			"-p", dbPort,
			"-U", dbUser,
			"-d", dbName,
			"-F", "c",
			"-Z", "9",
		)
	case "schema":
		backupFile = filepath.Join(backupDir, fmt.Sprintf("%s-schema.sql.gz", timestamp))
		cmd = exec.Command("pg_dump",
			"-h", dbHost,
			"-p", dbPort,
			"-U", dbUser,
			"-d", dbName,
			"-F", "c",
			"-Z", "9",
			"--schema-only",
		)
	case "data":
		backupFile = filepath.Join(backupDir, fmt.Sprintf("%s-data.sql.gz", timestamp))
		cmd = exec.Command("pg_dump",
			"-h", dbHost,
			"-p", dbPort,
			"-U", dbUser,
			"-d", dbName,
			"-F", "c",
			"-Z", "9",
			"--data-only",
		)
	default:
		return "", fmt.Errorf("unsupported backup type: %s", backupType)
	}

	// Set password
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", dbPassword))

	file, err := os.Create(backupFile)
	if err != nil {
		return "", fmt.Errorf("failed to create backup file: %w", err)
	}
	defer file.Close()

	cmd.Stdout = file
	cmd.Stderr = os.Stderr

	p.logger.Info("Starting PostgreSQL backup", "type", backupType, "file", backupFile)

	if err := cmd.Run(); err != nil {
		os.Remove(backupFile)
		return "", fmt.Errorf("pg_dump failed: %w", err)
	}

	p.logger.Info("PostgreSQL backup completed", "file", backupFile)
	return backupFile, nil
}

// Restore restores a PostgreSQL backup using pg_restore
func (p *PostgresBackup) Restore(backupFile string) error {
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "5432")
	dbUser := getEnvOrDefault("DB_USER", "infrapilot")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "infrapilot")
	dbName := getEnvOrDefault("DB_NAME", "infrapilot")

	p.logger.Info("Starting PostgreSQL restore", "file", backupFile)

	cmd := exec.Command("pg_restore",
		"-h", dbHost,
		"-p", dbPort,
		"-U", dbUser,
		"-d", dbName,
		"-c", // Clean before restore
		"--if-exists",
		backupFile,
	)

	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", dbPassword))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pg_restore failed: %w", err)
	}

	p.logger.Info("PostgreSQL restore completed")
	return nil
}

// ListBackups lists all local PostgreSQL backup files
func (p *PostgresBackup) ListBackups() ([]string, error) {
	backupDir := "./backups/postgres"
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, nil
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

// BackupSchemaOnly performs a schema-only backup
func (p *PostgresBackup) BackupSchemaOnly() (string, error) {
	return p.Backup("schema")
}

// BackupDataOnly performs a data-only backup
func (p *PostgresBackup) BackupDataOnly() (string, error) {
	return p.Backup("data")
}
