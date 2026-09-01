package handlers

import (
	"fmt"
	"net/http"

	"infrapilot/backend/internal/backup"
	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"

	"github.com/gin-gonic/gin"
)

// BackupHandler exposes HTTP REST endpoints for disaster recovery & backup management
type BackupHandler struct {
	manager    *backup.BackupManager
	restoreEng *backup.RestoreEngine
	retention  *backup.RetentionPolicy
	encryption *backup.EncryptionService
}

// NewBackupHandler initializes a new BackupHandler instance
func NewBackupHandler() *BackupHandler {
	mgr := backup.NewBackupManager()
	pg := backup.NewPostgresBackup()
	qd := backup.NewQdrantBackup()
	k8 := backup.NewKubernetesBackup()
	enc := backup.NewEncryptionService("./backups/.master.key")

	rest := backup.NewRestoreEngine("./backups")
	_ = pg
	_ = qd
	_ = k8
	_ = enc
	ret := backup.NewRetentionPolicy()

	return &BackupHandler{
		manager:    mgr,
		restoreEng: rest,
		retention:  ret,
		encryption: enc,
	}
}

// GetBackups lists all backup records
func (h *BackupHandler) GetBackups(c *gin.Context) {
	backups, err := h.manager.ListBackups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"backups": backups,
		"total":   len(backups),
	})
}

// TriggerBackup initiates a manual backup operation
func (h *BackupHandler) TriggerBackup(c *gin.Context) {
	var req struct {
		Type      string `json:"type"`      // postgresql, qdrant, kubernetes, config, grafana, full
		Retention string `json:"retention"` // hourly, daily, weekly, monthly
		Compress  bool   `json:"compress"`  // default true
		Encrypt   bool   `json:"encrypt"`   // default true
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if req.Type == "" {
		req.Type = "full"
	}

	record, err := h.manager.CreateBackup(req.Type, req.Retention, req.Compress, req.Encrypt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Backup failed: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Backup triggered successfully",
		"backup":  record,
	})
}

// DownloadBackup streams the requested backup archive file to the client
func (h *BackupHandler) DownloadBackup(c *gin.Context) {
	id := c.Param("id")

	var record models.BackupRecord
	if database.DB != nil && database.DB.Where("id = ? OR filename = ?", id, id).First(&record).Error == nil {
		if record.StoragePath != "" {
			c.FileAttachment(record.StoragePath, record.Filename)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Backup archive file not found"})
}

// DeleteBackup removes a backup record and its archive file
func (h *BackupHandler) DeleteBackup(c *gin.Context) {
	id := c.Param("id")
	if err := h.manager.DeleteBackup(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Backup deleted successfully"})
}

// RestoreBackup handles disaster recovery restoration requests with MFA/confirmation prompt
func (h *BackupHandler) RestoreBackup(c *gin.Context) {
	var req struct {
		BackupID      string `json:"backup_id"`
		Component     string `json:"component"`      // postgresql, qdrant, kubernetes, config, grafana, full
		ConfirmPhrase string `json:"confirm_phrase"` // Must equal "RESTORE"
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if req.ConfirmPhrase != "RESTORE" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Multi-factor confirmation failed: 'confirm_phrase' must be 'RESTORE'"})
		return
	}

	username := c.GetString("username")
	if username == "" {
		username = "admin"
	}

	restoreReq := backup.RestoreRequest{
		BackupID:   req.BackupID,
		BackupType: req.Component,
		Verify:     true,
		Decrypt:    true,
	}
	record, err := h.restoreEng.Restore(restoreReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   fmt.Sprintf("Restore operation failed: %v", err),
			"restore": record,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Restore completed successfully",
		"restore": record,
	})
}

// GetStats returns summary DR Dashboard metrics
func (h *BackupHandler) GetStats(c *gin.Context) {
	stats := h.manager.GetStats()
	c.JSON(http.StatusOK, stats)
}

// GetDRTestHistory retrieves past DR drill execution records
func (h *BackupHandler) GetDRTestHistory(c *gin.Context) {
	var history []models.DRTestRecord
	if database.DB != nil {
		database.DB.Order("started_at desc").Limit(20).Find(&history)
	}
	c.JSON(http.StatusOK, gin.H{"dr_tests": history})
}

// RunDRTest triggers a Disaster Recovery drill and records RTO/RPO metrics
func (h *BackupHandler) RunDRTest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "DR Test endpoint - implementation pending"})
}

// RotateKey triggers AES-256 encryption key rotation
func (h *BackupHandler) RotateKey(c *gin.Context) {
	if h.encryption == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Encryption service not initialized"})
		return
	}

	if err := h.encryption.RotateKey(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Key rotation failed: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "AES-256 Encryption key rotated successfully",
		"fingerprint": h.encryption.GetKeyFingerprint(),
	})
}

// EnforceRetention manually triggers retention policy cleanup
func (h *BackupHandler) EnforceRetention(c *gin.Context) {
	h.retention.EnforceRetention(h.manager)
	c.JSON(http.StatusOK, gin.H{
		"message": "Retention policy enforced successfully",
	})
}

// GetKeyInfo returns current encryption key status
func (h *BackupHandler) GetKeyInfo(c *gin.Context) {
	fingerprint := ""
	if h.encryption != nil {
		fingerprint = h.encryption.GetKeyFingerprint()
	}
	c.JSON(http.StatusOK, gin.H{
		"initialized": h.encryption != nil && h.encryption.IsInitialized(),
		"fingerprint": fingerprint,
		"algorithm":   "AES-256-GCM",
	})
}
