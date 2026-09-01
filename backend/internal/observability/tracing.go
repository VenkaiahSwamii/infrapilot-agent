package observability

import (
	"context"
	"fmt"

	"infrapilot/backend/internal/config"
	"infrapilot/backend/internal/logger"
)

// InitTracing initializes OpenTelemetry tracing if configured
// Returns nil if tracing is not configured or unavailable
func InitTracing() (interface{}, error) {
	serviceName := config.GetEnv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		logger.Debug("OpenTelemetry tracing not configured (OTEL_SERVICE_NAME not set)")
		return nil, nil
	}

	logger.Info("OpenTelemetry tracing configured",
		"service_name", serviceName,
	)

	// OpenTelemetry is optional - requires additional dependencies
	// Install with: go get go.opentelemetry.io/otel/...
	return nil, nil
}

// ShutdownTracing shuts down tracing provider
func ShutdownTracing(tp interface{}) {
	if tp != nil {
		logger.Debug("Shutting down tracing provider")
	}
}

// StartSpan creates a new tracing span (no-op with placeholder)
func StartSpan(ctx context.Context, name string) (context.Context, interface{}) {
	return ctx, nil
}

// SpanFromContext gets span from context
func SpanFromContext(ctx context.Context) interface{} {
	return nil
}

// AddEvent adds an event to the current span
func AddEvent(ctx context.Context, name string, attributes ...interface{}) {
	// No-op when tracing is not available
}

// RecordError records an error in the current span
func RecordError(ctx context.Context, err error) {
	// No-op when tracing is not available
}

// SetAttributes sets attributes on the current span
func SetAttributes(ctx context.Context, attributes ...interface{}) {
	// No-op when tracing is not available
}

// GetTracer returns the global tracer (unused, kept for API compatibility)
func GetTracer() {}

// Deprecated: fmtPrintln is a bridge to the structured logger
func fmtPrintln(v ...interface{}) {
	msg := fmt.Sprint(v...)
	logger.Debug(msg)
}

// Deprecated: fmtPrintf is a bridge to the structured logger
func fmtPrintf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	logger.Debug(msg)
}
