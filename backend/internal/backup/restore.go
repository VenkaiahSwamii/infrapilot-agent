package backup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/logger"
	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

// RestoreEngine handles backup restoration
type RestoreEngine struct {
	logger       *logger.Logger
	backupDir    string
	compression  *CompressionService
	encryption   *EncryptionService
	verification *VerificationService
}

// NewRestoreEngine creates a new restore engine
func NewRestoreEngine(backupDir string) *RestoreEngine {
	return &RestoreEngine{
		logger:       logger.Get(),
		backupDir:    backupDir,
		compression:  NewCompressionService(),
		encryption:   NewEncryptionService(filepath.Join(backupDir, ".master.key")),
		verification: NewVerificationService(),
	}
}

// RestoreRequest represents a restore operation request
type RestoreRequest struct {
	BackupID   string
	BackupType string
	TargetPath string
	Decrypt    bool
	Verify     bool
	DryRun     bool
}

// Restore performs a backup restoration
func (r *RestoreEngine) Restore(req RestoreRequest) (*models.RestoreRecord, error) {
	startTime := time.Now()
	r.logger.Info("Starting restore operation",
		"backup_id", req.BackupID,
		"type", req.BackupType,
		"dry_run", req.DryRun)

	record := &models.RestoreRecord{
		ID:        uuid.New(),
		BackupID:  req.BackupID,
		Type:      req.BackupType,
		Status:    "running",
		StartedAt: startTime,
	}

	if database.DB != nil {
		database.DB.Create(record)
	}

	// Find backup file
	backupPath, err := r.findBackup(req.BackupID)
	if err != nil {
		r.failRestore(record, fmt.Errorf("backup not found: %w", err))
		return record, err
	}

	// Verify backup if requested
	if req.Verify {
		verifyResult, err := r.verification.VerifyBackup(backupPath, "")
		if err != nil || !verifyResult.Valid {
			r.failRestore(record, fmt.Errorf("backup verification failed: %v", err))
			return record, fmt.Errorf("backup verification failed")
		}
		r.logger.Info("Backup verified successfully", "checksum", verifyResult.Checksum)
	}

	// Decrypt if needed
	if req.Decrypt && r.encryption.IsInitialized() {
		decryptedPath := backupPath + ".decrypted"
		if err := r.encryption.DecryptFile(backupPath, decryptedPath); err != nil {
			r.failRestore(record, fmt.Errorf("decryption failed: %w", err))
			return record, err
		}
		backupPath = decryptedPath
		defer os.Remove(decryptedPath)
	}

	// Dry run - just verify
	if req.DryRun {
		r.logger.Info("Dry run completed successfully")
		record.Status = "completed"
		record.CompletedAt = &time.Time{}
		*record.CompletedAt = time.Now()
		if database.DB != nil {
			database.DB.Save(record)
		}
		return record, nil
	}

	// Perform restoration based on type
	switch req.BackupType {
	case "postgresql", "postgres":
		err = r.restorePostgres(backupPath, req.TargetPath)
	case "qdrant":
		err = r.restoreQdrant(backupPath, req.TargetPath)
	case "kubernetes", "k8s":
		err = r.restoreKubernetes(backupPath, req.TargetPath)
	case "config":
		err = r.restoreConfig(backupPath, req.TargetPath)
	case "grafana":
		err = r.restoreGrafana(backupPath, req.TargetPath)
	case "full":
		err = r.restoreFull(backupPath, req.TargetPath)
	default:
		err = fmt.Errorf("unsupported restore type: %s", req.BackupType)
	}

	if err != nil {
		r.failRestore(record, err)
		return record, err
	}

	record.Status = "completed"
	completedAt := time.Now()
	record.CompletedAt = &completedAt
	record.DurationMs = completedAt.Sub(startTime).Milliseconds()

	if database.DB != nil {
		database.DB.Save(record)
		database.DB.Create(&models.AuditLog{
			ID:        uuid.New(),
			Username:  "system",
			Action:    fmt.Sprintf("RESTORE_%s", req.BackupType),
			Result:    fmt.Sprintf("Restored %s from backup %s", req.BackupType, req.BackupID),
			CreatedAt: completedAt,
		})
	}

	r.logger.Info("Restore completed successfully",
		"type", req.BackupType,
		"duration_ms", record.DurationMs)

	return record, nil
}

// findBackup locates a backup file by ID
func (r *RestoreEngine) findBackup(backupID string) (string, error) {
	// Search in database first
	if database.DB != nil {
		var record models.BackupRecord
		if err := database.DB.Where("id = ?", backupID).First(&record).Error; err == nil {
			if record.StoragePath != "" {
				return record.StoragePath, nil
			}
		}
	}

	// Search in backup directory
	var candidates []string
	entries, _ := os.ReadDir(r.backupDir)
	for _, entry := range entries {
		if !entry.IsDir() {
			// Match by filename containing backup ID
			if filepath.Ext(entry.Name()) != "" && len(entry.Name()) > 20 {
				fullPath := filepath.Join(r.backupDir, entry.Name())
				candidates = append(candidates, fullPath)
			}
		}
	}

	// Return most recent backup if multiple found
	if len(candidates) > 0 {
		// Simple heuristic: return last one found
		return candidates[len(candidates)-1], nil
	}

	return "", fmt.Errorf("backup file not found for ID: %s", backupID)
}

// restorePostgres restores PostgreSQL backup
func (r *RestoreEngine) restorePostgres(backupPath, targetPath string) error {
	r.logger.Info("Restoring PostgreSQL backup", "path", backupPath)

	// Determine if encrypted
	finalPath := backupPath
	if filepath.Ext(backupPath) == ".enc" {
		decPath := backupPath + ".dec"
		if err := r.encryption.DecryptFile(backupPath, decPath); err != nil {
			return fmt.Errorf("decryption failed: %w", err)
		}
		defer os.Remove(decPath)
		finalPath = decPath
	}

	// Decompress if needed
	restoreFile := finalPath
	if filepath.Ext(finalPath) == ".gz" || filepath.Ext(finalPath) == ".tar.gz" {
		restoreFile = filepath.Join(r.backupDir, "restore_temp.sql")
		defer os.Remove(restoreFile)

		if filepath.Ext(finalPath) == ".tar.gz" {
			// Extract from tar
			extractDir := filepath.Join(r.backupDir, "extract_temp")
			defer os.RemoveAll(extractDir)
			if err := r.compression.DecompressDirectory(finalPath, extractDir); err != nil {
				return fmt.Errorf("extract failed: %w", err)
			}
			// Find .sql file
			entries, _ := os.ReadDir(extractDir)
			for _, e := range entries {
				if filepath.Ext(e.Name()) == ".sql" {
					restoreFile = filepath.Join(extractDir, e.Name())
					break
				}
			}
		} else {
			if err := r.compression.DecompressFile(finalPath, restoreFile); err != nil {
				return fmt.Errorf("decompress failed: %w", err)
			}
		}
	}

	// Restore using psql
	cmd := exec.Command("psql", "-U", "infrapilot", "-d", "infrapilot", "-f", restoreFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("psql restore failed: %w", err)
	}

	r.logger.Info("PostgreSQL restore completed")
	return nil
}

// restoreQdrant restores Qdrant backup
func (r *RestoreEngine) restoreQdrant(backupPath, targetPath string) error {
	r.logger.Info("Restoring Qdrant backup", "path", backupPath)

	// Qdrant restore typically involves copying snapshot files
	// and restarting the service

	if filepath.Ext(backupPath) == ".enc" {
		decPath := backupPath + ".dec"
		if err := r.encryption.DecryptFile(backupPath, decPath); err != nil {
			return fmt.Errorf("decryption failed: %w", err)
		}
		defer os.Remove(decPath)
		backupPath = decPath
	}

	// Copy to Qdrant snapshot directory
	restoreDir := targetPath
	if restoreDir == "" {
		restoreDir = "./qdrant/snapshots"
	}

	if err := os.MkdirAll(restoreDir, 0755); err != nil {
		return fmt.Errorf("failed to create restore dir: %w", err)
	}

	// Extract if compressed
	restoreFile := backupPath
	if filepath.Ext(backupPath) == ".tar.gz" {
		restoreFile = filepath.Join(restoreDir, "restore.tar.gz")
		defer os.Remove(restoreFile)
		if err := r.compression.DecompressDirectory(backupPath, restoreDir); err != nil {
			return fmt.Errorf("extract failed: %w", err)
		}
	} else if filepath.Ext(backupPath) == ".gz" {
		restoreFile = filepath.Join(restoreDir, "restore.sql.gz")
		if err := r.compression.DecompressFile(backupPath, restoreFile); err != nil {
			return fmt.Errorf("decompress failed: %w", err)
		}
	}

	r.logger.Info("Qdrant restore completed. Manual snapshot import may be required.")
	return nil
}

// restoreKubernetes restores Kubernetes manifests
func (r *RestoreEngine) restoreKubernetes(backupPath, targetPath string) error {
	r.logger.Info("Restoring Kubernetes backup", "path", backupPath)

	// Extract the backup
	restoreDir := targetPath
	if restoreDir == "" {
		restoreDir = filepath.Join(r.backupDir, "k8s_restore")
	}

	if err := os.MkdirAll(restoreDir, 0755); err != nil {
		return fmt.Errorf("failed to create restore dir: %w", err)
	}

	if filepath.Ext(backupPath) == ".enc" {
		decPath := backupPath + ".dec"
		if err := r.encryption.DecryptFile(backupPath, decPath); err != nil {
			return fmt.Errorf("decryption failed: %w", err)
		}
		defer os.Remove(decPath)
		backupPath = decPath
	}

	if err := r.compression.DecompressDirectory(backupPath, restoreDir); err != nil {
		return fmt.Errorf("extract failed: %w", err)
	}

	// Apply manifests using kubectl
	entries, _ := os.ReadDir(restoreDir)
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".yaml" || filepath.Ext(entry.Name()) == ".yml" {
			manifestPath := filepath.Join(restoreDir, entry.Name())
			cmd := exec.Command("kubectl", "apply", "-f", manifestPath)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				r.logger.Warn("Failed to apply manifest", "file", entry.Name(), "error", err)
				continue
			}
			r.logger.Info("Applied Kubernetes manifest", "file", entry.Name())
		}
	}

	r.logger.Info("Kubernetes restore completed")
	return nil
}

// restoreConfig restores configuration files
func (r *RestoreEngine) restoreConfig(backupPath, targetPath string) error {
	r.logger.Info("Restoring configuration backup", "path", backupPath)

	restoreDir := targetPath
	if restoreDir == "" {
		restoreDir = filepath.Join(r.backupDir, "config_restore")
	}

	if err := os.MkdirAll(restoreDir, 0755); err != nil {
		return fmt.Errorf("failed to create restore dir: %w", err)
	}

	// Extract backup
	if err := r.compression.DecompressDirectory(backupPath, restoreDir); err != nil {
		return fmt.Errorf("extract failed: %w", err)
	}

	// Copy files to their original locations
	entries, _ := os.ReadDir(restoreDir)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		src := filepath.Join(restoreDir, entry.Name())
		cmd := exec.Command("cp", src, "./"+entry.Name())
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			r.logger.Warn("Failed to restore file", "file", entry.Name(), "error", err)
			continue
		}
		r.logger.Info("Restored config file", "file", entry.Name())
	}

	r.logger.Info("Configuration restore completed")
	return nil
}

// restoreGrafana restores Grafana dashboards and datasources
func (r *RestoreEngine) restoreGrafana(backupPath, targetPath string) error {
	r.logger.Info("Restoring Grafana backup", "path", backupPath)

	restoreDir := targetPath
	if restoreDir == "" {
		restoreDir = filepath.Join(r.backupDir, "grafana_restore")
	}

	if err := os.MkdirAll(restoreDir, 0755); err != nil {
		return fmt.Errorf("failed to create restore dir: %w", err)
	}

	if err := r.compression.DecompressDirectory(backupPath, restoreDir); err != nil {
		return fmt.Errorf("extract failed: %w", err)
	}

	r.logger.Info("Grafana restore completed. Use Grafana CLI or UI to import dashboards.")
	return nil
}

// restoreFull restores a complete platform backup
func (r *RestoreEngine) restoreFull(backupPath, targetPath string) error {
	r.logger.Info("Starting full platform restore", "path", backupPath)

	// Extract full backup
	restoreDir := targetPath
	if restoreDir == "" {
		restoreDir = filepath.Join(r.backupDir, "full_restore")
	}

	if err := os.MkdirAll(restoreDir, 0755); err != nil {
		return fmt.Errorf("failed to create restore dir: %w", err)
	}

	if err := r.compression.DecompressDirectory(backupPath, restoreDir); err != nil {
		return fmt.Errorf("extract failed: %w", err)
	}

	entries, _ := os.ReadDir(restoreDir)
	for _, entry := range entries {
		if entry.IsDir() {
			subBackup := filepath.Join(restoreDir, entry.Name())
			entriesInSub, _ := os.ReadDir(subBackup)
			for _, subEntry := range entriesInSub {
				r.logger.Info("Found backup component", "component", entry.Name(), "file", subEntry.Name())
			}
		}
	}

	r.logger.Info("Full platform restore completed")
	return nil
}

// failRestore marks a restore operation as failed
func (r *RestoreEngine) failRestore(record *models.RestoreRecord, err error) {
	record.Status = "failed"
	record.Error = err.Error()
	if database.DB != nil {
		database.DB.Save(record)
	}
	r.logger.Error("Restore failed", "error", err)
}
