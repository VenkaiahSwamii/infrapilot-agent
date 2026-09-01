package middleware

import (
	"crypto/x509"
	"net/http"

	"infrapilot/backend/internal/config"
	"infrapilot/backend/internal/security"

	"github.com/gin-gonic/gin"
)

// MTLSMiddleware verifies client certificate from mTLS connection
func MTLSMiddleware() gin.HandlerFunc {
	cfg := config.Get()

	var caPool *x509.CertPool
	if cfg.MTLSEnabled {
		pool, err := security.LoadCACert(cfg.TLSCAPath)
		if err != nil {
			panic(err)
		}
		caPool = pool
	}

	return func(c *gin.Context) {
		if !cfg.MTLSEnabled || caPool == nil {
			c.Next()
			return
		}

		// Verify peer certificate exists
		if c.Request.TLS == nil || len(c.Request.TLS.PeerCertificates) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Client certificate required",
			})
			c.Abort()
			return
		}

		// Validate certificate against CA
		opts := x509.VerifyOptions{
			Roots: caPool,
		}

		_, err := c.Request.TLS.PeerCertificates[0].Verify(opts)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid client certificate",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
