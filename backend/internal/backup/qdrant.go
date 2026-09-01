package backup

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

	"infrapilot/backend/internal/logger"
)

// QdrantBackup handles Qdrant-specific backup operations
type QdrantBackup struct {
	logger *logger.Logger
}

// NewQdrantBackup creates a new Qdrant backup handler
func NewQdrantBackup() *QdrantBackup {
	return &QdrantBackup{
		logger: logger.Get(),
	}
}

// QdrantCollection represents a Qdrant collection for backup
type QdrantCollection struct {
	Name    string `json:"name"`
	Vectors int64  `json:"vectors_count"`
}

// Backup performs a Qdrant backup using HTTP API
func (q *QdrantBackup) Backup() (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	backupDir := "./backups/qdrant"

	// Ensure directory exists
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	backupFile := fmt.Sprintf("%s/backup_%s_qdrant.tar.gz", backupDir, timestamp)

	// Create backup metadata
	metadata := map[string]interface{}{
		"timestamp": timestamp,
		"type":      "qdrant",
		"version":   "1.0",
	}

	metadataJSON, _ := json.MarshalIndent(metadata, "", "  ")
	metadataFile := fmt.Sprintf("%s/metadata_%s.json", backupDir, timestamp)
	if err := os.WriteFile(metadataFile, metadataJSON, 0644); err != nil {
		return "", fmt.Errorf("failed to write metadata: %w", err)
	}

	// Use curl to trigger Qdrant snapshot creation
	qdrantHost := getEnvOrDefault("QDRANT_HOST", "localhost")
	qdrantPort := getEnvOrDefault("QDRANT_PORT", "6333")

	// Create snapshot via Qdrant API
	snapshotCmd := exec.Command("curl", "-X", "POST",
		fmt.Sprintf("http://%s:%s/snapshots", qdrantHost, qdrantPort),
		"-H", "Content-Type: application/json",
		"-d", `{"wait": true}`)

	if err := snapshotCmd.Run(); err != nil {
		q.logger.Warn("Qdrant snapshot creation failed, using file-based backup", "error", err)
	}

	// Create tar.gz of snapshots directory if it exists
	// Qdrant stores snapshots in ./qdrant_snapshots by default
	snapshotsDir := "./qdrant_snapshots"
	if _, err := os.Stat(snapshotsDir); err == nil {
		cmd := exec.Command("tar", "-czf", backupFile, "-C", ".", "qdrant_snapshots/")
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to create qdrant backup archive: %w", err)
		}
	} else {
		// Create an empty tar.gz with metadata if no snapshots directory
		cmd := exec.Command("tar", "-czf", backupFile, "-C", backupDir, fmt.Sprintf("metadata_%s.json", timestamp))
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to create qdrant backup archive: %w", err)
		}
	}

	q.logger.Info("Qdrant backup completed", "file", backupFile)
	return backupFile, nil
}

// Restore restores Qdrant from backup
func (q *QdrantBackup) Restore(backupFile string) error {
	q.logger.Info("Starting Qdrant restore", "file", backupFile)

	// Create temporary extraction directory
	extractDir := fmt.Sprintf("./backups/qdrant/restore_%d", time.Now().Unix())
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		return fmt.Errorf("failed to create extraction directory: %w", err)
	}
	defer os.RemoveAll(extractDir)

	// Extract tar.gz
	cmd := exec.Command("tar", "-xzf", backupFile, "-C", extractDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to extract backup: %w", err)
	}

	// Restore snapshots to Qdrant
	qdrantHost := getEnvOrDefault("QDRANT_HOST", "localhost")
	qdrantPort := getEnvOrDefault("QDRANT_PORT", "6333")

	// Check if we have snapshot files and restore them
	snapshotFiles, err := os.ReadDir(extractDir)
	if err != nil {
		return fmt.Errorf("failed to read extraction directory: %w", err)
	}

	for _, file := range snapshotFiles {
		if file.Name() == "snapshots" || file.IsDir() {
			// Use Qdrant API to restore from snapshot
			restoreCmd := exec.Command("curl", "-X", "POST",
				fmt.Sprintf("http://%s:%s/snapshots/restore", qdrantHost, qdrantPort),
				"-H", "Content-Type: application/json",
				"-d", fmt.Sprintf(`{"location": "%s/%s"}`, extractDir, file.Name()))

			if err := restoreCmd.Run(); err != nil {
				q.logger.Warn("Qdrant restore command failed", "error", err)
			}
		}
	}

	q.logger.Info("Qdrant restore completed")
	return nil
}

// ListCollections lists all Qdrant collections
func (q *QdrantBackup) ListCollections() ([]QdrantCollection, error) {
	qdrantHost := getEnvOrDefault("QDRANT_HOST", "localhost")
	qdrantPort := getEnvOrDefault("QDRANT_PORT", "6333")

	cmd := exec.Command("curl", "-s",
		fmt.Sprintf("http://%s:%s/collections", qdrantHost, qdrantPort))

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list collections: %w", err)
	}

	var result struct {
		Result struct {
			Collections []QdrantCollection `json:"collections"`
		} `json:"result"`
	}

	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse collections response: %w", err)
	}

	return result.Result.Collections, nil
}
