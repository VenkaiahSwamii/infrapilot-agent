package backup

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"infrapilot/backend/internal/logger"
)

// EncryptionService handles backup encryption/decryption
type EncryptionService struct {
	logger    *logger.Logger
	keyPath   string
	masterKey []byte
}

// NewEncryptionService creates a new encryption service
func NewEncryptionService(keyPath string) *EncryptionService {
	service := &EncryptionService{
		logger:  logger.Get(),
		keyPath: keyPath,
	}

	// Load or generate master key
	service.loadOrGenerateKey()

	return service
}

// loadOrGenerateKey loads existing key or generates a new one
func (e *EncryptionService) loadOrGenerateKey() {
	// Try to load existing key
	if keyBytes, err := os.ReadFile(e.keyPath); err == nil {
		e.masterKey = keyBytes
		e.logger.Info("Encryption key loaded", "path", e.keyPath)
		return
	}

	// Generate new key
	e.masterKey = make([]byte, 32) // AES-256
	if _, err := rand.Read(e.masterKey); err != nil {
		e.logger.Error("Failed to generate encryption key", "error", err)
		return
	}

	// Save key to file
	if err := os.WriteFile(e.keyPath, e.masterKey, 0600); err != nil {
		e.logger.Error("Failed to save encryption key", "error", err)
	} else {
		e.logger.Info("New encryption key generated", "path", e.keyPath)
	}
}

// EncryptFile encrypts a file using AES-256-GCM
func (e *EncryptionService) EncryptFile(inputPath, outputPath string) error {
	if len(e.masterKey) == 0 {
		return fmt.Errorf("encryption key not initialized")
	}

	// Read input file
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Create cipher
	block, err := aes.NewCipher(e.masterKey)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}

	// Create nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("failed to create nonce: %w", err)
	}

	// Encrypt
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)

	// Write encrypted file
	if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
		return fmt.Errorf("failed to write encrypted file: %w", err)
	}

	e.logger.Info("File encrypted", "input", inputPath, "output", outputPath)
	return nil
}

// DecryptFile decrypts a file using AES-256-GCM
func (e *EncryptionService) DecryptFile(inputPath, outputPath string) error {
	if len(e.masterKey) == 0 {
		return fmt.Errorf("encryption key not initialized")
	}

	// Read encrypted file
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read encrypted file: %w", err)
	}

	// Create cipher
	block, err := aes.NewCipher(e.masterKey)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}

	// Extract nonce
	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("failed to decrypt: %w", err)
	}

	// Write decrypted file
	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return fmt.Errorf("failed to write decrypted file: %w", err)
	}

	e.logger.Info("File decrypted", "input", inputPath, "output", outputPath)
	return nil
}

// RotateKey rotates the encryption key
func (e *EncryptionService) RotateKey() error {
	e.logger.Info("Rotating encryption key")

	// Backup old key
	oldKey := e.masterKey
	backupPath := e.keyPath + ".old"
	if err := os.WriteFile(backupPath, oldKey, 0600); err != nil {
		return fmt.Errorf("failed to backup old key: %w", err)
	}

	// Generate new key
	e.masterKey = make([]byte, 32)
	if _, err := rand.Read(e.masterKey); err != nil {
		// Restore old key
		e.masterKey = oldKey
		os.Remove(backupPath)
		return fmt.Errorf("failed to generate new key: %w", err)
	}

	// Save new key
	if err := os.WriteFile(e.keyPath, e.masterKey, 0600); err != nil {
		e.masterKey = oldKey
		os.Remove(e.keyPath)
		os.Rename(backupPath, e.keyPath)
		return fmt.Errorf("failed to save new key: %w", err)
	}

	// Remove backup
	os.Remove(backupPath)

	e.logger.Info("Encryption key rotated successfully")
	return nil
}

// GetKeyFingerprint returns a fingerprint of the current key
func (e *EncryptionService) GetKeyFingerprint() string {
	if len(e.masterKey) == 0 {
		return ""
	}

	hash := sha256.Sum256(e.masterKey)
	return hex.EncodeToString(hash[:8])
}

// IsInitialized checks if encryption is initialized
func (e *EncryptionService) IsInitialized() bool {
	return len(e.masterKey) > 0
}
