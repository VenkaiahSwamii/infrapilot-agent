package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestSecurityHeadersMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(SecurityHeadersMiddleware())
	r.GET("/test-security", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req, _ := http.NewRequest("GET", "/test-security", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Header().Get("Strict-Transport-Security") == "" {
		t.Errorf("Missing HSTS header")
	}
	if w.Header().Get("X-Frame-Options") != "DENY" {
		t.Errorf("Expected X-Frame-Options DENY, got %s", w.Header().Get("X-Frame-Options"))
	}
	if w.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Errorf("Expected X-Content-Type-Options nosniff, got %s", w.Header().Get("X-Content-Type-Options"))
	}
	if w.Header().Get("X-XSS-Protection") == "" {
		t.Errorf("Missing X-XSS-Protection header")
	}
}

func TestRateLimiterMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	limiter := NewRateLimiter(3, 1*time.Minute)

	r := gin.New()
	r.Use(limiter.RateLimit())
	r.GET("/test-rate", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// First 3 requests should pass
	for i := 0; i < 3; i++ {
		req, _ := http.NewRequest("GET", "/test-rate", nil)
		req.RemoteAddr = "192.168.1.100:12345"
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Request %d expected 200, got %d", i+1, w.Code)
		}
	}

	// 4th request should return 429 Too Many Requests
	req, _ := http.NewRequest("GET", "/test-rate", nil)
	req.RemoteAddr = "192.168.1.100:12345"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status 429 Too Many Requests, got %d", w.Code)
	}
}

func TestCORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(CORSMiddleware())
	r.OPTIONS("/api/v1/test", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req, _ := http.NewRequest("OPTIONS", "/api/v1/test", nil)
	req.Header.Set("Origin", "https://app.infrapilot.io")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent && w.Code != http.StatusOK {
		t.Errorf("Expected OPTIONS status 204/200, got %d", w.Code)
	}
}
