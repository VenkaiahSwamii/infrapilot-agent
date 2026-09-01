package backup

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/logger"
	"infrapilot/backend/internal/utils"

	"github.com/google/uuid"
)

const BackupDir = "./backups"

// RunBackup generates a SQL dump of the database tables and writes it to a file.
func RunBackup() (string, error) {
	if database.DB == nil {
		return "", fmt.Errorf("database connection not initialized")
	}

	if err := os.MkdirAll(BackupDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	fileName := fmt.Sprintf("backup_%s.sql", time.Now().Format("20060102_150405"))
	filePath := filepath.Join(BackupDir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create backup file: %w", err)
	}
	defer file.Close()

	// Write backup header
	file.WriteString("-- InfraPilot Enterprise Database Backup\n")
	file.WriteString(fmt.Sprintf("-- Generated: %s\n\n", time.Now().Format(time.RFC3339)))

	// Fetch all tables
	var tables []string
	query := "SELECT table_name FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA() AND table_type = 'BASE TABLE'"
	if err := database.DB.Raw(query).Scan(&tables).Error; err != nil {
		return "", fmt.Errorf("failed to list database tables: %w", err)
	}

	for _, table := range tables {
		// Skip temporary or GORM migration tables if needed, but backing them up is generally good
		if table == "gorm_migrations" {
			continue
		}

		file.WriteString(fmt.Sprintf("-- Table: %s\n", table))

		// Fetch table columns
		var columns []string
		colQuery := fmt.Sprintf("SELECT column_name FROM information_schema.columns WHERE table_name = '%s' AND table_schema = CURRENT_SCHEMA()", table)
		if err := database.DB.Raw(colQuery).Scan(&columns).Error; err != nil {
			logger.Error("Failed to fetch columns", "table", table, "error", err)
			continue
		}

		if len(columns) == 0 {
			continue
		}

		// Fetch all rows
		var rows []map[string]interface{}
		if err := database.DB.Table(table).Find(&rows).Error; err != nil {
			logger.Error("Failed to fetch rows", "table", table, "error", err)
			continue
		}

		for _, row := range rows {
			var vals []string
			for _, col := range columns {
				val := row[col]
				if val == nil {
					vals = append(vals, "NULL")
				} else {
					// Escape quotes and format SQL value
					valStr := fmt.Sprintf("%v", val)
					valStr = strings.ReplaceAll(valStr, "'", "''")
					vals = append(vals, fmt.Sprintf("'%s'", valStr))
				}
			}

			insertStmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);\n",
				table,
				strings.Join(columns, ", "),
				strings.Join(vals, ", "),
			)
			file.WriteString(insertStmt)
		}
		file.WriteString("\n")
	}

	logger.Info("Database backup completed successfully", "file", fileName)
	return fileName, nil
}

// RestoreBackup executes the SQL statements from a backup file to restore the database.
func RestoreBackup(fileName string) error {
	if database.DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	filePath := filepath.Join(BackupDir, fileName)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	// Simple SQL parser by statement
	statements := strings.Split(string(content), ";\n")
	tx := database.DB.Begin()
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}

		if err := tx.Exec(stmt).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute SQL statement: %w", err)
		}
	}
	tx.Commit()

	logger.Info("Database restored successfully", "file", fileName)
	return nil
}

// BackupJob implements the scheduler.Job interface for automated database backups.
type BackupJob struct{}

func (b *BackupJob) Name() string {
	return "DatabaseBackup"
}

func (b *BackupJob) Execute(ctx context.Context) error {
	fileName, err := RunBackup()
	if err != nil {
		utils.LogAudit("SYSTEM", uuid.Nil, "Scheduled Backup", "Failure: "+err.Error())
		return err
	}
	utils.LogAudit("SYSTEM", uuid.Nil, "Scheduled Backup: "+fileName, "Success")
	return nil
}

func (b *BackupJob) Schedule() string {
	// Daily backup at midnight
	return "0 0 * * *"
}
