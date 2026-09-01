package plugins

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

// LogsPlugin collects system and application logs
type LogsPlugin struct {
	enabled bool
	mu      sync.RWMutex
	health  PluginHealth
}

// NewLogsPlugin creates a new Logs plugin instance
func NewLogsPlugin() *LogsPlugin {
	return &LogsPlugin{
		enabled: true,
		health:  NewPluginHealth("logs"),
	}
}

// Name returns the plugin name
func (p *LogsPlugin) Name() string {
	return "logs"
}

// IsEnabled returns whether the plugin is enabled
func (p *LogsPlugin) IsEnabled() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.enabled
}

// SetEnabled enables or disables the plugin
func (p *LogsPlugin) SetEnabled(enabled bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.enabled = enabled
}

// HealthCheck returns the health status of the plugin
func (p *LogsPlugin) HealthCheck() PluginHealth {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.health
}

// Collect gathers log metrics
func (p *LogsPlugin) Collect() (interface{}, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	startTime := time.Now()

	result := make(map[string]interface{})

	// Collect log files
	logFiles, err := p.collectLogFiles()
	if err != nil {
		p.health.Status = "degraded"
		p.health.Error = err.Error()
		result["error"] = err.Error()
	} else {
		p.health.Status = "healthy"
		p.health.Error = ""
		result["log_files"] = logFiles
	}

	// Collect recent log entries
	logEntries, err := p.collectRecentLogs()
	if err != nil {
		result["log_entries_error"] = err.Error()
	} else {
		result["recent_logs"] = logEntries
	}

	// Collect log statistics
	logStats, err := p.collectLogStatistics()
	if err != nil {
		result["log_stats_error"] = err.Error()
	} else {
		result["statistics"] = logStats
	}

	p.health.LastRun = time.Now()
	p.health.Metadata["duration_ms"] = time.Since(startTime).String()

	return result, nil
}

// LogFile represents a log file with metadata
type LogFile struct {
	Path     string    `json:"path"`
	Size     int64     `json:"size_bytes"`
	Modified time.Time `json:"modified"`
	Lines    int       `json:"lines"`
	Service  string    `json:"service,omitempty"`
	LogLevel string    `json:"log_level,omitempty"`
}

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"` // INFO, WARN, ERROR, DEBUG
	Service   string    `json:"service"`
	Message   string    `json:"message"`
	Source    string    `json:"source"`
}

// LogStatistics represents aggregated log statistics
type LogStatistics struct {
	TotalLogs     int            `json:"total_logs"`
	ErrorCount    int            `json:"error_count"`
	WarningCount  int            `json:"warning_count"`
	InfoCount     int            `json:"info_count"`
	DebugCount    int            `json:"debug_count"`
	LogsByService map[string]int `json:"logs_by_service"`
	LogsByHour    map[int]int    `json:"logs_by_hour"`
	TopErrors     []string       `json:"top_errors"`
}

// collectLogFiles gathers information about log files
func (p *LogsPlugin) collectLogFiles() ([]LogFile, error) {
	var logFiles []LogFile

	// Common log directories
	logDirs := []string{
		"/var/log",
		"/var/log/syslog",
		"/var/log/messages",
	}

	if runtime.GOOS == "windows" {
		logDirs = []string{
			"C:\\Windows\\Logs",
			"C:\\Windows\\System32\\winevt\\Logs",
		}
	}

	for _, logDir := range logDirs {
		if _, err := os.Stat(logDir); os.IsNotExist(err) {
			continue
		}

		// Walk through log directory
		_ = filepath.Walk(logDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			// Only process regular files
			if info.IsDir() {
				return nil
			}

			// Check if it's a log file
			if !p.isLogFile(path) {
				return nil
			}

			// Count lines
			lines := 0
			if content, err := os.ReadFile(path); err == nil {
				lines = len(strings.Split(string(content), "\n"))
			}

			logFiles = append(logFiles, LogFile{
				Path:     path,
				Size:     info.Size(),
				Modified: info.ModTime(),
				Lines:    lines,
				Service:  p.extractServiceName(path),
				LogLevel: p.detectLogLevel(path),
			})

			return nil
		})
	}

	// Sort by modification time (most recent first)
	sort.Slice(logFiles, func(i, j int) bool {
		return logFiles[i].Modified.After(logFiles[j].Modified)
	})

	// Limit to top 20 log files
	if len(logFiles) > 20 {
		logFiles = logFiles[:20]
	}

	return logFiles, nil
}

// collectRecentLogs gathers recent log entries
func (p *LogsPlugin) collectRecentLogs() ([]LogEntry, error) {
	return []LogEntry{}, nil
}

// collectLogStatistics gathers log statistics
func (p *LogsPlugin) collectLogStatistics() (LogStatistics, error) {
	return LogStatistics{LogsByService: map[string]int{}, LogsByHour: map[int]int{}, TopErrors: []string{}}, nil
}

// isLogFile checks if a file is a log file
func (p *LogsPlugin) isLogFile(path string) bool {
	ext := filepath.Ext(path)
	logExtensions := []string{".log", ".out", ".err", ".txt"}

	for _, logExt := range logExtensions {
		if ext == logExt {
			return true
		}
	}

	// Check if filename contains 'log'
	base := filepath.Base(path)
	return regexp.MustCompile(`(?i)log`).MatchString(base)
}

// extractServiceName extracts service name from log file path
func (p *LogsPlugin) extractServiceName(path string) string {
	base := filepath.Base(path)

	// Remove extension
	base = regexp.MustCompile(`\.(log|out|err|txt)$`).ReplaceAllString(base, "")

	// Remove common suffixes
	base = regexp.MustCompile(`[-_]?\d{4}$`).ReplaceAllString(base, "")
	base = regexp.MustCompile(`[-_]?\d{2}$`).ReplaceAllString(base, "")

	return base
}

// detectLogLevel detects the log level from file path
func (p *LogsPlugin) detectLogLevel(path string) string {
	base := strings.ToLower(filepath.Base(path))

	if strings.Contains(base, "error") {
		return "ERROR"
	} else if strings.Contains(base, "warn") {
		return "WARN"
	} else if strings.Contains(base, "debug") {
		return "DEBUG"
	}

	return "INFO"
}
