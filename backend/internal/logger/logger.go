package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Level represents log severity
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger is a structured logger for InfraPilot
type Logger struct {
	mu       sync.Mutex
	level    Level
	output   io.Writer
	jsonMode bool
	prefix   string
}

var (
	defaultLogger *Logger
	once          sync.Once
)

// Init initializes the default logger
func Init(level Level, jsonMode bool, output io.Writer) {
	once.Do(func() {
		if output == nil {
			output = os.Stdout
		}
		defaultLogger = &Logger{
			level:    level,
			output:   output,
			jsonMode: jsonMode,
		}
		log.SetOutput(output)
		log.SetFlags(0)
	})
}

// Get returns the default logger
func Get() *Logger {
	if defaultLogger == nil {
		Init(INFO, false, os.Stdout)
	}
	return defaultLogger
}

// SetLevel changes the log level at runtime
func SetLevel(level Level) {
	if defaultLogger != nil {
		defaultLogger.mu.Lock()
		defer defaultLogger.mu.Unlock()
		defaultLogger.level = level
	}
}

// log is the internal logging method
func (l *Logger) log(level Level, msg string, args ...interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now().UTC().Format(time.RFC3339)

	// Get caller info (skip 2 frames: log -> Info/Debug/etc -> caller)
	_, file, line, ok := runtime.Caller(2)
	caller := "?"
	if ok {
		// Short file path
		short := file
		if idx := strings.LastIndex(file, "/"); idx >= 0 {
			short = file[idx+1:]
		}
		caller = fmt.Sprintf("%s:%d", short, line)
	}

	if l.jsonMode {
		l.writeJSON(now, level, msg, caller, args...)
	} else {
		l.writeText(now, level, msg, caller, args...)
	}
}

func (l *Logger) writeJSON(now string, level Level, msg string, caller string, args ...interface{}) {
	// Build key-value pairs
	fields := fmt.Sprintf(`"time":"%s","level":"%s","caller":"%s","msg":"%s"`, now, level, caller, escapeJSON(msg))
	for i := 0; i < len(args); i += 2 {
		key := fmt.Sprintf("%v", args[i])
		var val string
		if i+1 < len(args) {
			val = fmt.Sprintf("%v", args[i+1])
		} else {
			val = ""
		}
		fields += fmt.Sprintf(`,"%s":"%s"`, escapeJSON(key), escapeJSON(val))
	}
	fmt.Fprintf(l.output, "{%s}\n", fields)
}

func (l *Logger) writeText(now string, level Level, msg string, caller string, args ...interface{}) {
	prefix := l.prefix
	// Build key=value pairs
	var pairs []string
	for i := 0; i < len(args); i += 2 {
		key := fmt.Sprintf("%v", args[i])
		var val string
		if i+1 < len(args) {
			val = fmt.Sprintf("%v", args[i+1])
		} else {
			val = ""
		}
		pairs = append(pairs, fmt.Sprintf("%s=%s", key, val))
	}

	kv := ""
	if len(pairs) > 0 {
		kv = " " + strings.Join(pairs, " ")
	}

	fmt.Fprintf(l.output, "%s [%s] %s %s%s%s\n", now, level, caller, prefix, msg, kv)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.log(DEBUG, msg, args...)
}

// Info logs an info message
func (l *Logger) Info(msg string, args ...interface{}) {
	l.log(INFO, msg, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.log(WARN, msg, args...)
}

// Error logs an error message
func (l *Logger) Error(msg string, args ...interface{}) {
	l.log(ERROR, msg, args...)
}

// Fatal logs an error message and exits
func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.log(ERROR, msg, args...)
	os.Exit(1)
}

// WithPrefix returns a logger with a prefix
func (l *Logger) WithPrefix(prefix string) *Logger {
	return &Logger{
		level:    l.level,
		output:   l.output,
		jsonMode: l.jsonMode,
		prefix:   l.prefix + "[" + prefix + "] ",
	}
}

// Package-level convenience functions
func Debug(msg string, args ...interface{}) { Get().Debug(msg, args...) }
func Info(msg string, args ...interface{})  { Get().Info(msg, args...) }
func Warn(msg string, args ...interface{})  { Get().Warn(msg, args...) }
func Error(msg string, args ...interface{}) { Get().Error(msg, args...) }
func Fatal(msg string, args ...interface{}) { Get().Fatal(msg, args...) }

func escapeJSON(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}
