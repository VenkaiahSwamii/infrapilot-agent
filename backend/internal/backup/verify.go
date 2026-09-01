package backup

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"

	"infrapilot/backend/internal/logger"
)

// VerificationService handles backup verification
type VerificationService struct {
	logger *logger.Logger
}

// NewVerificationService creates a new verification service
func NewVerificationService() *VerificationService {
	return &VerificationService{
		logger: logger.Get(),
	}
}

// VerifyResult contains the result of a backup verification
type VerifyResult struct {
	BackupID string
	Valid    bool
	Checksum string
	Expected string
	Size     int64
	Error    string
}

// VerifyBackup verifies a backup file integrity
func (v *VerificationService) VerifyBackup(backupFile string, expectedChecksum string) (*VerifyResult, error) {
	v.logger.Info("Verifying backup", "file", backupFile)

	result := &VerifyResult{
		BackupID: backupFile,
	}

	// Check file exists
	fileInfo, err := os.Stat(backupFile)
	if err != nil {
		result.Error = fmt.Sprintf("file not found: %v", err)
		return result, nil
	}

	result.Size = fileInfo.Size()

	// Calculate checksum
	checksum, err := v.calculateSHA256(backupFile)
	if err != nil {
		result.Error = fmt.Sprintf("failed to calculate checksum: %v", err)
		return result, nil
	}

	result.Checksum = checksum

	// Verify checksum
	if expectedChecksum != "" {
		result.Expected = expectedChecksum
		if checksum == expectedChecksum {
			result.Valid = true
		} else {
			result.Valid = false
			result.Error = "checksum mismatch"
		}
	} else {
		// No expected checksum, just report the calculated one
		result.Valid = true
	}

	v.logger.Info("Backup verification completed",
		"file", backupFile,
		"valid", result.Valid,
		"checksum", checksum)

	return result, nil
}

// VerifyArchive verifies a compressed archive integrity
func (v *VerificationService) VerifyArchive(archiveFile string) (*VerifyResult, error) {
	v.logger.Info("Verifying archive", "file", archiveFile)

	result := &VerifyResult{
		BackupID: archiveFile,
	}

	// Check file exists
	fileInfo, err := os.Stat(archiveFile)
	if err != nil {
		result.Error = fmt.Sprintf("file not found: %v", err)
		return result, nil
	}

	result.Size = fileInfo.Size()

	// Calculate checksum
	checksum, err := v.calculateSHA256(archiveFile)
	if err != nil {
		result.Error = fmt.Sprintf("failed to calculate checksum: %v", err)
		return result, nil
	}

	result.Checksum = checksum

	// Verify tar.gz archive
	if len(archiveFile) > 7 && archiveFile[len(archiveFile)-7:] == ".tar.gz" {
		valid, err := v.verifyTarGz(archiveFile)
		if err != nil {
			result.Error = fmt.Sprintf("archive verification failed: %v", err)
			result.Valid = false
			return result, nil
		}
		result.Valid = valid
	}

	v.logger.Info("Archive verification completed",
		"file", archiveFile,
		"valid", result.Valid)

	return result, nil
}

// verifyTarGz verifies a tar.gz archive can be read
func (v *VerificationService) verifyTarGz(archiveFile string) (bool, error) {
	// Try to list contents without extracting
	cmd := exec.Command("tar", "-tzf", archiveFile)
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("tar list failed: %w", err)
	}

	// Check if we got any output
	if len(output) == 0 {
		return false, fmt.Errorf("archive appears to be empty")
	}

	return true, nil
}

// CalculateChecksum calculates SHA256 checksum of a file
func (v *VerificationService) CalculateChecksum(filePath string) (string, error) {
	return v.calculateSHA256(filePath)
}

// CalculateMD5 calculates MD5 checksum of a file
func (v *VerificationService) CalculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to calculate MD5: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// calculateSHA256 calculates SHA256 checksum of a file
func (v *VerificationService) calculateSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to calculate SHA256: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// VerifyBackupChain verifies a chain of backups
func (v *VerificationService) VerifyBackupChain(backupFiles []string) []VerifyResult {
	results := make([]VerifyResult, 0, len(backupFiles))

	for _, backupFile := range backupFiles {
		result, err := v.VerifyBackup(backupFile, "")
		if err != nil {
			v.logger.Error("Failed to verify backup", "file", backupFile, "error", err)
			continue
		}
		results = append(results, *result)
	}

	return results
}
