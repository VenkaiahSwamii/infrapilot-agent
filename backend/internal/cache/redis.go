package cache

import (
	"context"
	"fmt"
	"time"

	"infrapilot/backend/internal/config"
	"infrapilot/backend/internal/logger"

	"github.com/redis/go-redis/v9"
)

// RedisClient is the global Redis client
var RedisClient *redis.Client

// ConnectRedis initializes Redis connection
func ConnectRedis() error {
	cfg := config.Get()

	RedisClient = redis.NewClient(&redis.Options{
		Addr:         cfg.GetRedisAddr(),
		Password:     cfg.RedisPassword,
		DB:           cfg.RedisDB,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     cfg.RedisPoolSize,
		MinIdleConns: 5,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Redis connected successfully",
		"addr", cfg.GetRedisAddr(),
		"db", cfg.RedisDB,
	)
	return nil
}

// GetRedisClient returns the global Redis client
func GetRedisClient() *redis.Client {
	return RedisClient
}

// CloseRedis closes Redis connection
func CloseRedis() {
	if RedisClient != nil {
		if err := RedisClient.Close(); err != nil {
			logger.Error("Error closing Redis connection",
				"error", err,
			)
		}
	}
}
