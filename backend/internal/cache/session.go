package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Session represents a user session stored in Redis
type Session struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	OrgID     string    `json:"org_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

const (
	// SessionPrefix is the Redis key prefix for sessions
	SessionPrefix = "session:"
	// SessionTTL is the default session time-to-live
	SessionTTL = 24 * time.Hour
)

// SetSession stores a session in Redis
func SetSession(sessionID string, session *Session) error {
	key := SessionPrefix + sessionID
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	client := GetRedisClient()
	return client.Set(context.Background(), key, data, SessionTTL).Err()
}

// GetSession retrieves a session from Redis
func GetSession(sessionID string) (*Session, error) {
	key := SessionPrefix + sessionID
	client := GetRedisClient()
	data, err := client.Get(context.Background(), key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	var session Session
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, nil
}

// DeleteSession removes a session from Redis
func DeleteSession(sessionID string) error {
	key := SessionPrefix + sessionID
	client := GetRedisClient()
	return client.Del(context.Background(), key).Err()
}

// RefreshSession extends session TTL
func RefreshSession(sessionID string) error {
	key := SessionPrefix + sessionID
	client := GetRedisClient()
	return client.Expire(context.Background(), key, SessionTTL).Err()
}

// SetCache stores a value in Redis cache with TTL
func SetCache(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache value: %w", err)
	}

	client := GetRedisClient()
	return client.Set(context.Background(), key, data, ttl).Err()
}

// GetCache retrieves a value from Redis cache
func GetCache(key string, dest interface{}) error {
	client := GetRedisClient()
	data, err := client.Get(context.Background(), key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

// DeleteCache removes a key from Redis
func DeleteCache(key string) error {
	client := GetRedisClient()
	return client.Del(context.Background(), key).Err()
}

// RateLimiter checks and enforces rate limits
type RateLimiter struct {
	prefix string
	limit  int64
	window time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(prefix string, limit int64, window time.Duration) *RateLimiter {
	return &RateLimiter{
		prefix: prefix,
		limit:  limit,
		window: window,
	}
}

// Allow checks if a request is allowed under the rate limit
func (rl *RateLimiter) Allow(key string) (bool, error) {
	redisKey := fmt.Sprintf("ratelimit:%s:%s", rl.prefix, key)

	client := GetRedisClient()
	count, err := client.Incr(context.Background(), redisKey).Result()
	if err != nil {
		return false, fmt.Errorf("failed to increment rate limit counter: %w", err)
	}

	if count == 1 {
		client.Expire(context.Background(), redisKey, rl.window)
	}

	return count <= rl.limit, nil
}

// Reset resets the rate limit counter
func (rl *RateLimiter) Reset(key string) error {
	redisKey := fmt.Sprintf("ratelimit:%s:%s", rl.prefix, key)
	client := GetRedisClient()
	if client == nil {
		return nil
	}
	return client.Del(context.Background(), redisKey).Err()
}

// RevokeToken adds a token to the Redis revocation blacklist
func RevokeToken(token string, ttl time.Duration) error {
	client := GetRedisClient()
	if client == nil {
		return nil
	}
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	return client.Set(context.Background(), "revoked:"+token, "true", ttl).Err()
}

// IsTokenRevoked checks if a token is in the Redis revocation blacklist
func IsTokenRevoked(token string) bool {
	client := GetRedisClient()
	if client == nil {
		return false
	}
	val, err := client.Get(context.Background(), "revoked:"+token).Result()
	return err == nil && val == "true"
}

// RecordFailedLogin increments failed login count in Redis
func RecordFailedLogin(email string) (int64, error) {
	client := GetRedisClient()
	if client == nil {
		return 0, nil
	}
	ctx := context.Background()
	key := "failed_logins:" + email
	count, err := client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if count == 1 {
		client.Expire(ctx, key, 15*time.Minute)
	}
	if count >= 5 {
		client.Set(ctx, "account_locked:"+email, "true", 15*time.Minute)
	}
	return count, nil
}

// IsAccountLocked checks if an account is locked in Redis
func IsAccountLocked(email string) bool {
	client := GetRedisClient()
	if client == nil {
		return false
	}
	val, err := client.Get(context.Background(), "account_locked:"+email).Result()
	return err == nil && val == "true"
}

// ResetFailedLogins resets the failed login attempt counter for an email
func ResetFailedLogins(email string) {
	client := GetRedisClient()
	if client == nil {
		return
	}
	ctx := context.Background()
	client.Del(ctx, "failed_logins:"+email)
	client.Del(ctx, "account_locked:"+email)
}
