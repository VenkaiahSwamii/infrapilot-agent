package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"infrapilot/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func TestRateLimiter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("allows requests within limit", func(t *testing.T) {
		rl := middleware.NewRateLimiter(5, time.Minute)
		router := gin.New()
		router.GET("/test", rl.RateLimit(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		for i := 0; i < 5; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "192.168.1.1:12345"
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Request %d expected 200, got %d", i+1, w.Code)
			}
		}
	})

	t.Run("blocks requests over limit", func(t *testing.T) {
		rl := middleware.NewRateLimiter(3, time.Minute)
		router := gin.New()
		router.GET("/test", rl.RateLimit(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// First 3 requests should succeed
		for i := 0; i < 3; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "10.0.0.1:12345"
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Request %d expected 200, got %d", i+1, w.Code)
			}
		}

		// 4th request should be rate limited
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "10.0.0.1:12345"
		router.ServeHTTP(w, req)

		if w.Code != http.StatusTooManyRequests {
			t.Errorf("Expected 429, got %d", w.Code)
		}
	})

	t.Run("different IPs have separate limits", func(t *testing.T) {
		rl := middleware.NewRateLimiter(2, time.Minute)
		router := gin.New()
		router.GET("/test", rl.RateLimit(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// IP 1: 2 requests
		for i := 0; i < 2; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "192.168.1.1:12345"
			router.ServeHTTP(w, req)
			if w.Code != http.StatusOK {
				t.Errorf("IP1 request %d expected 200, got %d", i+1, w.Code)
			}
		}

		// IP 2: 2 requests (should work independently)
		for i := 0; i < 2; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "192.168.1.2:54321"
			router.ServeHTTP(w, req)
			if w.Code != http.StatusOK {
				t.Errorf("IP2 request %d expected 200, got %d", i+1, w.Code)
			}
		}

		// IP 1: 3rd request should fail
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		router.ServeHTTP(w, req)
		if w.Code != http.StatusTooManyRequests {
			t.Errorf("IP1 3rd request expected 429, got %d", w.Code)
		}
	})

	t.Run("resets after window expires", func(t *testing.T) {
		rl := middleware.NewRateLimiter(1, 100*time.Millisecond)
		router := gin.New()
		router.GET("/test", rl.RateLimit(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// First request succeeds
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "10.0.0.1:12345"
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("First request expected 200, got %d", w.Code)
		}

		// Second request (within window) fails
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "10.0.0.1:12345"
		router.ServeHTTP(w, req)
		if w.Code != http.StatusTooManyRequests {
			t.Errorf("Second request expected 429, got %d", w.Code)
		}

		// Wait for window to expire
		time.Sleep(150 * time.Millisecond)

		// Third request (new window) succeeds
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "10.0.0.1:12345"
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Third request expected 200, got %d", w.Code)
		}
	})
}
