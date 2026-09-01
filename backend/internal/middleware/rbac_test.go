package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"infrapilot/backend/internal/models"
)

// helper to perform a request with a given role and required roles
func performRequest(t *testing.T, userRole string, requiredRoles ...string) int {
	// set Gin to test mode
	gin.SetMode(gin.TestMode)
	r := gin.New()
	// middleware to set role in context
	r.Use(func(c *gin.Context) {
		c.Set("role", userRole)
		c.Next()
	})
	// apply RequireRoles middleware
	r.Use(RequireRoles(requiredRoles...))
	// dummy handler should not be called on deny
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)
	return w.Code
}

func TestRequireRolesAllowed(t *testing.T) {
	// Admin role should be allowed when Admin is among required roles
	code := performRequest(t, models.RoleAdmin, models.RoleAdmin, models.RoleViewer)
	if code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, code)
	}
}

func TestRequireRolesForbidden(t *testing.T) {
	// Viewer role should be denied when only Admin is required
	code := performRequest(t, models.RoleViewer, models.RoleAdmin)
	if code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, code)
	}
}

func TestRequireRolesUnauthorized(t *testing.T) {
	// No role set in context should result in 401
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequireRoles(models.RoleAdmin))
	r.GET("/test", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"msg": "ok"}) })
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}
