package backup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"infrapilot/backend/internal/logger"
)

// StorageType represents the type of remote storage
type StorageType string

const (
	StorageTypeLocal StorageType = "local"
	StorageTypeS3    StorageType = "s3"
	StorageTypeMinIO StorageType = "minio"
	StorageTypeAzure StorageType = "azure"
	StorageTypeGCS   StorageType = "gcs"
	StorageTypeNAS   StorageType = "nas"
)

// StorageConfig holds storage configuration
type StorageConfig struct {
	Type      StorageType
	Endpoint  string
	Bucket    string
	Region    string
	AccessKey string
	SecretKey string
	Path      string // Local path for NAS or local storage
}

// StorageService handles backup storage operations
type StorageService struct {
	logger *logger.Logger
	config *StorageConfig
}

// NewStorageService creates a new storage service
func NewStorageService(config *StorageConfig) *StorageService {
	return &StorageService{
		logger: logger.Get(),
		config: config,
	}
}

// Upload uploads a backup file to remote storage
func (s *StorageService) Upload(localPath string) (string, error) {
	s.logger.Info("Uploading backup", "path", localPath, "storage", s.config.Type)

	switch s.config.Type {
	case StorageTypeLocal, StorageTypeNAS:
		return s.uploadLocal(localPath)
	case StorageTypeS3, StorageTypeMinIO:
		return s.uploadS3(localPath)
	case StorageTypeAzure:
		return s.uploadAzure(localPath)
	case StorageTypeGCS:
		return s.uploadGCS(localPath)
	default:
		return "", fmt.Errorf("unsupported storage type: %s", s.config.Type)
	}
}

// uploadLocal copies file to local/NAS storage
func (s *StorageService) uploadLocal(localPath string) (string, error) {
	destPath := filepath.Join(s.config.Path, filepath.Base(localPath))

	// Create destination directory if needed
	if err := os.MkdirAll(s.config.Path, 0755); err != nil {
		return "", fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Copy file
	cmd := exec.Command("cp", localPath, destPath)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	s.logger.Info("Backup uploaded to local storage", "path", destPath)
	return destPath, nil
}

// uploadS3 uploads to S3 or MinIO
func (s *StorageService) uploadS3(localPath string) (string, error) {
	endpoint := s.config.Endpoint
	if endpoint == "" && s.config.Type == StorageTypeMinIO {
		endpoint = "http://localhost:9000"
	}

	cmdArgs := []string{"s3", "cp", localPath,
		fmt.Sprintf("s3://%s/%s", s.config.Bucket, filepath.Base(localPath)),
	}

	if endpoint != "" && s.config.Type == StorageTypeMinIO {
		cmdArgs = append([]string{"--endpoint-url", endpoint}, cmdArgs...)
	}

	cmd := exec.Command("aws", cmdArgs...)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", s.config.AccessKey),
		fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", s.config.SecretKey),
		fmt.Sprintf("AWS_DEFAULT_REGION=%s", s.config.Region),
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("S3 upload failed: %s: %w", string(output), err)
	}

	s3Path := fmt.Sprintf("s3://%s/%s", s.config.Bucket, filepath.Base(localPath))
	s.logger.Info("Backup uploaded to S3", "path", s3Path)
	return s3Path, nil
}

// uploadAzure uploads to Azure Blob Storage
func (s *StorageService) uploadAzure(localPath string) (string, error) {
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s",
		s.config.Endpoint, s.config.Bucket)

	blobPath := fmt.Sprintf("%s/%s", containerURL, filepath.Base(localPath))

	args := []string{
		"storage", "blob", "upload",
		"--container-name", s.config.Bucket,
		"--name", filepath.Base(localPath),
		"--file", localPath,
		"--account-name", s.config.Endpoint,
	}
	cmd := exec.Command("az", args...)

	cmd.Env = append(os.Environ(),
		fmt.Sprintf("AZURE_STORAGE_KEY=%s", s.config.AccessKey),
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Azure upload failed: %s: %w", string(output), err)
	}

	s.logger.Info("Backup uploaded to Azure", "path", blobPath)
	return blobPath, nil
}

// uploadGCS uploads to Google Cloud Storage
func (s *StorageService) uploadGCS(localPath string) (string, error) {
	gcsPath := fmt.Sprintf("gs://%s/%s", s.config.Bucket, filepath.Base(localPath))

	cmd := exec.Command("gsutil", "cp", localPath, gcsPath)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GOOGLE_APPLICATION_CREDENTIALS=%s", s.config.AccessKey),
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("GCS upload failed: %s: %w", string(output), err)
	}

	s.logger.Info("Backup uploaded to GCS", "path", gcsPath)
	return gcsPath, nil
}

// Download downloads a backup from remote storage
func (s *StorageService) Download(remotePath, localDir string) (string, error) {
	s.logger.Info("Downloading backup", "remote", remotePath, "local", localDir)

	if err := os.MkdirAll(localDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create local directory: %w", err)
	}

	filename := filepath.Base(remotePath)
	localPath := filepath.Join(localDir, filename)

	switch s.config.Type {
	case StorageTypeLocal, StorageTypeNAS:
		cmd := exec.Command("cp", remotePath, localPath)
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to copy file: %w", err)
		}
	case StorageTypeS3, StorageTypeMinIO:
		endpoint := s.config.Endpoint
		if endpoint == "" && s.config.Type == StorageTypeMinIO {
			endpoint = "http://localhost:9000"
		}

		cmdArgs := []string{"s3", "cp", remotePath, localPath}
		if endpoint != "" && s.config.Type == StorageTypeMinIO {
			cmdArgs = append([]string{"--endpoint-url", endpoint}, cmdArgs...)
		}

		cmd := exec.Command("aws", cmdArgs...)
		cmd.Env = append(os.Environ(),
			fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", s.config.AccessKey),
			fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", s.config.SecretKey),
			fmt.Sprintf("AWS_DEFAULT_REGION=%s", s.config.Region),
		)

		if output, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("S3 download failed: %s: %w", string(output), err)
		}
	default:
		return "", fmt.Errorf("unsupported storage type for download: %s", s.config.Type)
	}

	s.logger.Info("Backup downloaded", "path", localPath)
	return localPath, nil
}

// Delete deletes a backup from remote storage
func (s *StorageService) Delete(remotePath string) error {
	s.logger.Info("Deleting backup", "path", remotePath)

	switch s.config.Type {
	case StorageTypeLocal, StorageTypeNAS:
		return os.Remove(remotePath)
	case StorageTypeS3, StorageTypeMinIO:
		endpoint := s.config.Endpoint
		if endpoint == "" && s.config.Type == StorageTypeMinIO {
			endpoint = "http://localhost:9000"
		}

		cmdArgs := []string{"s3", "rm", remotePath}
		if endpoint != "" && s.config.Type == StorageTypeMinIO {
			cmdArgs = append([]string{"--endpoint-url", endpoint}, cmdArgs...)
		}
		cmd := exec.Command("aws", cmdArgs...)
		cmd.Env = append(os.Environ(),
			fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", s.config.AccessKey),
			fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", s.config.SecretKey),
			fmt.Sprintf("AWS_DEFAULT_REGION=%s", s.config.Region),
		)

		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("S3 delete failed: %s: %w", string(output), err)
		}
	default:
		return fmt.Errorf("unsupported storage type for delete: %s", s.config.Type)
	}

	s.logger.Info("Backup deleted", "path", remotePath)
	return nil
}

// List lists backups in remote storage
func (s *StorageService) List() ([]string, error) {
	s.logger.Info("Listing backups", "storage", s.config.Type)

	var backups []string

	switch s.config.Type {
	case StorageTypeLocal, StorageTypeNAS:
		files, err := os.ReadDir(s.config.Path)
		if err != nil {
			return nil, fmt.Errorf("failed to list files: %w", err)
		}
		for _, file := range files {
			if !file.IsDir() {
				backups = append(backups, file.Name())
			}
		}
	case StorageTypeS3, StorageTypeMinIO:
		endpoint := s.config.Endpoint
		if endpoint == "" && s.config.Type == StorageTypeMinIO {
			endpoint = "http://localhost:9000"
		}

		cmdArgs := []string{"s3", "ls", fmt.Sprintf("s3://%s", s.config.Bucket)}
		if endpoint != "" && s.config.Type == StorageTypeMinIO {
			cmdArgs = append([]string{"--endpoint-url", endpoint}, cmdArgs...)
		}

		cmd := exec.Command("aws", cmdArgs...)
		cmd.Env = append(os.Environ(),
			fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", s.config.AccessKey),
			fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", s.config.SecretKey),
			fmt.Sprintf("AWS_DEFAULT_REGION=%s", s.config.Region),
		)

		output, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("S3 list failed: %w", err)
		}
		backups = append(backups, string(output))
	default:
		return nil, fmt.Errorf("unsupported storage type for list: %s", s.config.Type)
	}

	return backups, nil
}

// GetStorageUsage returns storage usage statistics
func (s *StorageService) GetStorageUsage() (int64, error) {
	switch s.config.Type {
	case StorageTypeLocal, StorageTypeNAS:
		var totalSize int64
		files, err := os.ReadDir(s.config.Path)
		if err != nil {
			return 0, fmt.Errorf("failed to read storage directory: %w", err)
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			info, err := file.Info()
			if err != nil {
				continue
			}
			totalSize += info.Size()
		}
		return totalSize, nil
	default:
		return 0, fmt.Errorf("storage usage not supported for type: %s", s.config.Type)
	}
}
