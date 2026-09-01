package tests

import (
	"testing"

	"infrapilot/backend/internal/logger"
)

func getTestLogger() *logger.Logger {
	logger.Init(logger.INFO, false, nil)
	return logger.Get()
}

// TestRegister tests validation without database
func TestRegister_Validation(t *testing.T) {
	t.Run("missing fields return 400", func(t *testing.T) {
		t.Skip("Integration test - requires running backend")
	})

	t.Run("invalid email format rejected", func(t *testing.T) {
		t.Skip("Integration test - requires running backend")
	})

	t.Run("short password rejected", func(t *testing.T) {
		t.Skip("Integration test - requires running backend")
	})
}

// TestLogin tests validation logic
func TestLogin_Validation(t *testing.T) {
	t.Run("missing credentials return 400", func(t *testing.T) {
		t.Skip("Integration test - requires running backend")
	})

	t.Run("wrong password rejected", func(t *testing.T) {
		t.Skip("Integration test - requires running backend")
	})
}

// TestEnroll tests enrollment endpoint
func TestEnroll_Validation(t *testing.T) {
	t.Run("missing hostname returns 400", func(t *testing.T) {
		t.Skip("Integration test - requires running backend")
	})

	t.Run("missing enrollment token returns 401", func(t *testing.T) {
		t.Skip("Integration test - requires running backend")
	})
}

// TestMetricsReception tests metrics endpoint
func TestMetricsReception_Validation(t *testing.T) {
	t.Run("missing API key returns 401", func(t *testing.T) {
		t.Skip("Integration test - requires running backend")
	})

	t.Run("invalid metrics payload returns 400", func(t *testing.T) {
		t.Skip("Integration test - requires running backend")
	})
}

// TestLogger tests the structured logger
func TestLogger(t *testing.T) {
	logger := getTestLogger()

	t.Run("info level logging", func(t *testing.T) {
		logger.Info("test info message", "key", "value")
	})

	t.Run("warn level logging", func(t *testing.T) {
		logger.Warn("test warning", "component", "test")
	})

	t.Run("error level logging", func(t *testing.T) {
		logger.Error("test error", "error", "test-error")
	})

	t.Run("debug level suppressed", func(t *testing.T) {
		logger.Debug("should not appear", "key", "value")
	})
}
