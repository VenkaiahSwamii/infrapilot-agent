package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"infrapilot/backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RenameRequest struct {
	OldPath string `json:"old_path" binding:"required"`
	NewPath string `json:"new_path" binding:"required"`
}

type WriteFileRequest struct {
	Path    string `json:"path" binding:"required"`
	Content string `json:"content"`
}

type MkdirRequest struct {
	Path string `json:"path" binding:"required"`
}

// ListFiles lists directory contents
func ListFiles(c *gin.Context) {
	dirPath := c.Query("path")
	if dirPath == "" {
		dirPath = "/"
		if runtime.GOOS == "windows" {
			dirPath = "C:\\"
		}
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type FileItem struct {
		Name    string `json:"name"`
		Path    string `json:"path"`
		Type    string `json:"type"`
		Size    int64  `json:"size,omitempty"`
		ModTime string `json:"mod_time"`
	}

	var response []FileItem = []FileItem{}
	for _, f := range files {
		info, err := f.Info()
		var size int64
		var modTime string
		if err == nil {
			size = info.Size()
			modTime = info.ModTime().Format(time.RFC3339)
		}

		fileType := "file"
		if f.IsDir() {
			fileType = "directory"
		}

		fullPath := filepath.Join(dirPath, f.Name())

		response = append(response, FileItem{
			Name:    f.Name(),
			Path:    fullPath,
			Type:    fileType,
			Size:    size,
			ModTime: modTime,
		})
	}

	c.JSON(http.StatusOK, response)
}

// ReadFileContent reads text content of a remote file
func ReadFileContent(c *gin.Context) {
	filePath := c.Query("path")
	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path query parameter is required"})
		return
	}

	// Limit to max 5MB for safety
	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Failed to open file: %v", err)})
		return
	}
	defer file.Close()

	content, err := io.ReadAll(io.LimitReader(file, 5*1024*1024))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to read file: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"path":    filePath,
		"content": string(content),
		"size":    len(content),
	})
}

// WriteFileContent saves text content to a remote file
func WriteFileContent(c *gin.Context) {
	var req WriteFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := os.WriteFile(req.Path, []byte(req.Content), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to write file: %v", err)})
		return
	}

	logFileAction(c, "Write", req.Path)

	c.JSON(http.StatusOK, gin.H{
		"message": "File saved successfully",
		"path":    req.Path,
	})
}

// DownloadFile handles file download transfer
func DownloadFile(c *gin.Context) {
	filePath := c.Query("path")
	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path query parameter is required"})
		return
	}

	c.File(filePath)
}

// UploadFile handles file upload ingestion
func UploadFile(c *gin.Context) {
	targetDir := c.Query("path")
	if targetDir == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path query parameter is required"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to parse uploaded file: %v", err)})
		return
	}

	cleanDir := filepath.Clean(targetDir)
	if stat, err := os.Stat(cleanDir); err != nil {
		// Directory does not exist, try to create it
		if err := os.MkdirAll(cleanDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create target directory: %v", err)})
			return
		}
	} else if !stat.IsDir() {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Target path '%s' is a file, not a directory", cleanDir)})
		return
	}

	targetPath := filepath.Join(cleanDir, filepath.Base(file.Filename))

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to open source file: %v", err)})
		return
	}
	defer src.Close()

	out, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Permission denied writing '%s' (%v). If targeting C:\\ root, try writing inside a subfolder such as C:\\Users or D:\\", filepath.Base(targetPath), err)})
		return
	}
	defer out.Close()

	if _, err = io.Copy(out, src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to stream upload content: %v", err)})
		return
	}

	logFileAction(c, "Upload", targetPath)

	c.JSON(http.StatusOK, gin.H{"message": "file uploaded successfully", "path": targetPath})
}

// DeleteFile deletes a file or directory
func DeleteFile(c *gin.Context) {
	targetPath := c.Query("path")
	if targetPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path query parameter is required"})
		return
	}

	if err := os.RemoveAll(targetPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logFileAction(c, "Delete", targetPath)

	c.JSON(http.StatusOK, gin.H{"message": "file deleted successfully"})
}

// RenameFile renames a file or directory
func RenameFile(c *gin.Context) {
	var req RenameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := os.Rename(req.OldPath, req.NewPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logFileAction(c, "Rename", req.OldPath+" -> "+req.NewPath)

	c.JSON(http.StatusOK, gin.H{"message": "file renamed successfully"})
}

// MakeDirectory creates a new directory
func MakeDirectory(c *gin.Context) {
	var req MkdirRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cleanPath := filepath.Clean(req.Path)
	if stat, err := os.Stat(cleanPath); err == nil && stat.IsDir() {
		c.JSON(http.StatusOK, gin.H{"message": "Directory already exists", "path": cleanPath})
		return
	}

	if err := os.MkdirAll(cleanPath, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create directory '%s': %v. Note: Creating folders directly in C:\\ root may require administrator permissions. Try creating inside a subfolder such as C:\\Users or D:\\", cleanPath, err)})
		return
	}

	logFileAction(c, "Mkdir", cleanPath)

	c.JSON(http.StatusOK, gin.H{"message": "Directory created successfully", "path": cleanPath})
}

func logFileAction(c *gin.Context, action string, target string) {
	usernameVal, exists := c.Get("username")
	username := "admin"
	if exists {
		username = fmt.Sprintf("%v", usernameVal)
	}

	utils.LogAudit(username, uuid.Nil, fmt.Sprintf("%s file: %s", action, target), "Success")
}
