package tracing

import (
	"fmt"
	"time"

	"infrapilot/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TracingMiddleware(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
			c.Header("X-Trace-ID", traceID)
		}

		start := time.Now()
		c.Set("trace_id", traceID)

		c.Next()

		duration := time.Since(start).Milliseconds()
		statusCode := c.Writer.Status()
		isSlow := duration > 500
		hasError := statusCode >= 500

		traceRecord := &models.Trace{
			ID:          uuid.New(),
			TraceID:     traceID,
			Name:        fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path),
			ServiceName: "api-server",
			HTTPMethod:  c.Request.Method,
			URLPath:     c.Request.URL.Path,
			StatusCode:  statusCode,
			DurationMs:  duration,
			HasError:    hasError,
			IsSlow:      isSlow,
			Timestamp:   start,
			CreatedAt:   time.Now(),
		}

		if service != nil {
			_ = service.RecordTrace(traceRecord)
		}
	}
}
