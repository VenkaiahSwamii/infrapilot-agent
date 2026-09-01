package cache

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// SetAPICache stores arbitrary data in Redis cache with expiration
func SetAPICache(ctx context.Context, key string, val string, ttl time.Duration) error {
	client := GetRedisClient()
	if client == nil {
		return nil
	}
	return client.Set(ctx, "apicache:"+key, val, ttl).Err()
}

// GetAPICache fetches string value from Redis cache
func GetAPICache(ctx context.Context, key string) (string, error) {
	client := GetRedisClient()
	if client == nil {
		return "", nil
	}
	return client.Get(ctx, "apicache:"+key).Result()
}

// CacheMiddleware provides automatic Redis response caching for GET endpoints
func CacheMiddleware(ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		client := GetRedisClient()
		if client == nil {
			c.Next()
			return
		}

		cacheKey := "apicache:http:" + c.Request.URL.String() + ":" + c.GetString("organizationId")
		ctx := c.Request.Context()

		val, err := client.Get(ctx, cacheKey).Result()
		if err == nil && val != "" {
			c.Header("X-Cache", "HIT")
			c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(val))
			c.Abort()
			return
		}

		c.Header("X-Cache", "MISS")
		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		c.Next()

		if c.Writer.Status() == http.StatusOK && w.body.Len() > 0 {
			_ = client.Set(ctx, cacheKey, w.body.String(), ttl).Err()
		}
	}
}
