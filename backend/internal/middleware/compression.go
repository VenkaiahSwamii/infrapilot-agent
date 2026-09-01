package middleware

import (
	"compress/gzip"
	"strings"

	"github.com/gin-gonic/gin"
)

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipWriter) Write(b []byte) (int, error) {
	return g.writer.Write(b)
}

func (g *gzipWriter) WriteString(s string) (int, error) {
	return g.writer.Write([]byte(s))
}

// GzipCompressionMiddleware provides Gzip response compression for API payloads
func GzipCompressionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bypass Gzip compression for WebSockets, downloads, installer scripts, and binary payloads
		path := c.Request.URL.Path
		if strings.EqualFold(c.GetHeader("Upgrade"), "websocket") ||
			strings.Contains(strings.ToLower(c.GetHeader("Connection")), "upgrade") ||
			path == "/ws" ||
			strings.Contains(path, "/ws/") ||
			strings.Contains(path, "/terminal/ws") ||
			strings.HasPrefix(path, "/downloads") ||
			strings.HasPrefix(path, "/api/v1/install") ||
			strings.HasPrefix(path, "/api/v1/agent/install") ||
			strings.HasSuffix(path, ".exe") ||
			strings.HasSuffix(path, ".sh") ||
			strings.HasSuffix(path, ".ps1") ||
			strings.Contains(path, "infrapilot-agent") {
			c.Next()
			return
		}

		if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")

		gz := gzip.NewWriter(c.Writer)
		defer func() {
			_ = gz.Close()
		}()

		c.Writer = &gzipWriter{ResponseWriter: c.Writer, writer: gz}
		c.Next()
	}
}
