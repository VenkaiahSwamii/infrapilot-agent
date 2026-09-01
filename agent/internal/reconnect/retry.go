package reconnect

import (
	"log"
	"time"
)

// RetryStrategy defines how retries should behave
type RetryStrategy struct {
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
}

// DefaultRetryStrategy returns sane defaults for exponential backoff
func DefaultRetryStrategy() *RetryStrategy {
	return &RetryStrategy{
		InitialInterval: 5 * time.Second,
		MaxInterval:     2 * time.Minute,
		Multiplier:      2.0,
	}
}

// ExecuteWithBackoff executes the provided operation with exponential backoff
func ExecuteWithBackoff(operation func() error, strategy *RetryStrategy, operationName string) {
	if operation == nil {
		return
	}
	if strategy == nil {
		strategy = DefaultRetryStrategy()
	}

	interval := strategy.InitialInterval
	maxRetries := 10
	retries := 0

	for {
		err := operation()
		if err == nil {
			if retries > 0 {
				log.Printf("[%s] Reconnected successfully after %d retries", operationName, retries)
			}
			return
		}

		retries++
		if retries > maxRetries {
			log.Printf("[%s] Max retries reached (%d), resetting backoff. Error: %v", operationName, maxRetries, err)
			retries = 0
			interval = strategy.InitialInterval
			continue
		}

		log.Printf("[%s] Operation failed (attempt %d/%d): %v. Retrying in %v", operationName, retries, maxRetries, err, interval)
		time.Sleep(interval)

		// Exponential backoff: double the interval each time
		interval = time.Duration(float64(interval) * strategy.Multiplier)
		if interval > strategy.MaxInterval {
			interval = strategy.MaxInterval
		}
	}
}
