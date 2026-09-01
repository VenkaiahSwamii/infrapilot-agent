package backup

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"infrapilot/backend/internal/logger"
)

// CompressionService handles backup compression
type CompressionService struct {
	logger *logger.Logger
}

// NewCompressionService creates a new compression service
func NewCompressionService() *CompressionService {
	return &CompressionService{
		logger: logger.Get(),
	}
}

// CompressionType represents the compression algorithm
type CompressionType string

const (
	CompressionTypeGzip  CompressionType = "gzip"
	CompressionTypeTarGz CompressionType = "tar.gz"
	CompressionTypeZip   CompressionType = "zip"
)

// CompressorConfig holds compression configuration
type CompressorConfig struct {
	Type     CompressionType
	Level    int  // 1-9 for gzip, 1-9 for zip
	Parallel bool // Use parallel compression
}

// CompressFile compresses a single file
func (c *CompressionService) CompressFile(inputPath, outputPath string, config CompressorConfig) error {
	c.logger.Info("Compressing file",
		"input", inputPath,
		"output", outputPath,
		"type", config.Type)

	switch config.Type {
	case CompressionTypeGzip:
		return c.compressGzip(inputPath, outputPath, config.Level)
	case CompressionTypeZip:
		return c.compressZip(inputPath, outputPath, config.Level)
	default:
		return fmt.Errorf("unsupported compression type: %s", config.Type)
	}
}

// CompressDirectory compresses an entire directory
func (c *CompressionService) CompressDirectory(inputDir, outputPath string, config CompressorConfig) error {
	c.logger.Info("Compressing directory",
		"input", inputDir,
		"output", outputPath,
		"type", config.Type)

	switch config.Type {
	case CompressionTypeTarGz:
		return c.compressTarGz(inputDir, outputPath, config.Level)
	default:
		return fmt.Errorf("unsupported compression type for directory: %s", config.Type)
	}
}

// compressGzip compresses a file using gzip
func (c *CompressionService) compressGzip(inputPath, outputPath string, level int) error {
	// Use system gzip command
	cmd := exec.Command("gzip", "-c", "-"+fmt.Sprintf("%d", level), inputPath)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	cmd.Stdout = outputFile
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("gzip compression failed: %w", err)
	}

	return nil
}

// compressZip compresses a file using zip
func (c *CompressionService) compressZip(inputPath, outputPath string, level int) error {
	// Use system zip command
	args := []string{"-j", "-" + fmt.Sprintf("%d", level), outputPath, inputPath}
	if level > 5 {
		args = []string{"-j", outputPath, inputPath}
	}

	cmd := exec.Command("zip", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("zip compression failed: %w", err)
	}

	return nil
}

// compressTarGz compresses a directory using tar.gz
func (c *CompressionService) compressTarGz(inputDir, outputPath string, level int) error {
	// Use system tar command with gzip compression
	cmd := exec.Command("tar", "-czf", outputPath, "-C", inputDir, ".")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tar.gz compression failed: %w", err)
	}

	return nil
}

// DecompressFile decompresses a file
func (c *CompressionService) DecompressFile(inputPath, outputPath string) error {
	c.logger.Info("Decompressing file", "input", inputPath, "output", outputPath)

	// Determine compression type from extension
	if len(inputPath) > 3 && inputPath[len(inputPath)-3:] == ".gz" {
		return c.decompressGzip(inputPath, outputPath)
	} else if len(inputPath) > 4 && inputPath[len(inputPath)-4:] == ".zip" {
		return c.decompressZip(inputPath, outputPath)
	}

	return fmt.Errorf("unsupported compression format")
}

// DecompressDirectory decompresses a tar.gz directory
func (c *CompressionService) DecompressDirectory(inputPath, outputDir string) error {
	c.logger.Info("Decompressing directory", "input", inputPath, "output", outputDir)

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Use tar to extract
	cmd := exec.Command("tar", "-xzf", inputPath, "-C", outputDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tar extraction failed: %w", err)
	}

	return nil
}

// decompressGzip decompresses a gzip file
func (c *CompressionService) decompressGzip(inputPath, outputPath string) error {
	// Determine output filename
	outputFile := outputPath
	if outputFile == "" {
		// Remove .gz extension
		if len(inputPath) > 3 {
			outputFile = inputPath[:len(inputPath)-3]
		}
	}

	// Use gunzip command
	cmd := exec.Command("gunzip", "-c", inputPath)

	out, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	cmd.Stdout = out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("gunzip decompression failed: %w", err)
	}

	return nil
}

// decompressZip decompresses a zip file
func (c *CompressionService) decompressZip(inputPath, outputDir string) error {
	// Create output directory if needed
	if outputDir == "" {
		outputDir = "."
	}

	// Use unzip command
	cmd := exec.Command("unzip", "-q", inputPath, "-d", outputDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("zip decompression failed: %w", err)
	}

	return nil
}

// GetCompressionRatio calculates the compression ratio
func (c *CompressionService) GetCompressionRatio(originalSize, compressedSize int64) float64 {
	if originalSize == 0 {
		return 0.0
	}
	return float64(compressedSize) / float64(originalSize)
}

// CompressWithTimestamp compresses a file/directory with timestamp in filename
func (c *CompressionService) CompressWithTimestamp(inputPath, baseOutputPath string, config CompressorConfig) (string, error) {
	timestamp := time.Now().Format("20060102_150405")

	// Determine extension based on compression type
	var ext string
	switch config.Type {
	case CompressionTypeGzip:
		ext = ".gz"
	case CompressionTypeTarGz:
		ext = ".tar.gz"
	case CompressionTypeZip:
		ext = ".zip"
	}

	outputPath := fmt.Sprintf("%s_%s%s", baseOutputPath, timestamp, ext)

	// Check if input is directory or file
	info, err := os.Stat(inputPath)
	if err != nil {
		return "", fmt.Errorf("failed to stat input: %w", err)
	}

	if info.IsDir() {
		if err := c.CompressDirectory(inputPath, outputPath, config); err != nil {
			return "", err
		}
	} else {
		if err := c.CompressFile(inputPath, outputPath, config); err != nil {
			return "", err
		}
	}

	c.logger.Info("Compression completed", "output", outputPath)
	return outputPath, nil
}
