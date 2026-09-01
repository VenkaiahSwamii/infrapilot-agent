package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"infrapilot/backend/internal/backup"
	"infrapilot/backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RestoreRequest struct {
	FileName string `json:"file_name" binding:"required"`
}

// RunBackupHandler runs a database backup.
func RunBackupHandler(c *gin.Context) {
	fileName, err := backup.RunBackup()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to run backup: " + err.Error()})
		return
	}

	actor := c.GetString("username")
	if actor == "" {
		actor = "admin"
	}
	utils.LogAudit(actor, uuid.Nil, "Manual Backup: "+fileName, "Success")

	c.JSON(http.StatusOK, gin.H{
		"message":   "Database backup completed successfully",
		"file_name": fileName,
	})
}

// ListBackupsHandler returns a list of database backup files.
func ListBackupsHandler(c *gin.Context) {
	files, err := os.ReadDir(backup.BackupDir)
	if err != nil {
		// If directory does not exist or empty, return empty list
		c.JSON(http.StatusOK, gin.H{"backups": []string{}})
		return
	}

	var backups []gin.H
	for _, f := range files {
		if f.IsDir() || !filepath.HasPrefix(f.Name(), "backup_") {
			continue
		}
		info, err := f.Info()
		if err != nil {
			continue
		}
		backups = append(backups, gin.H{
			"file_name":  f.Name(),
			"size_bytes": info.Size(),
			"created_at": info.ModTime(),
		})
	}

	c.JSON(http.StatusOK, gin.H{"backups": backups})
}

// RestoreBackupHandler restores the database from a backup file.
func RestoreBackupHandler(c *gin.Context) {
	var req RestoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := backup.RestoreBackup(req.FileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore backup: " + err.Error()})
		return
	}

	actor := c.GetString("username")
	if actor == "" {
		actor = "admin"
	}
	utils.LogAudit(actor, uuid.Nil, "Database Restored from: "+req.FileName, "Success")

	c.JSON(http.StatusOK, gin.H{"message": "Database restore completed successfully"})
}
