package apm

import (
	"time"

	"infrapilot/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func APMMiddleware(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := float64(time.Since(start).Microseconds()) / 1000.0 // ms float64
		statusCode := c.Writer.Status()

		traceIDVal, _ := c.Get("trace_id")
		var traceUUID uuid.UUID
		if str, ok := traceIDVal.(string); ok {
			traceUUID, _ = uuid.Parse(str)
		}

		reqRecord := &models.APIRequest{
			ID:           uuid.New(),
			TraceID:      traceUUID,
			Method:       c.Request.Method,
			Endpoint:     c.Request.URL.Path,
			StatusCode:   statusCode,
			DurationMs:   duration,
			RequestSize:  c.Request.ContentLength,
			ResponseSize: int64(c.Writer.Size()),
			ClientIP:     c.ClientIP(),
			IsSlow:       duration > 500.0,
			IsError:      statusCode >= 400,
			CreatedAt:    start,
		}

		if service != nil {
			_ = service.RecordAPIRequest(reqRecord)
		}
	}
}
