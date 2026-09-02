package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"infrapilot/backend/internal/cache"

	"github.com/redis/go-redis/v9"
)

const (
	// Stream prefixes
	MetricStream    = "infrapilot:metrics"
	AlertStream     = "infrapilot:alerts"
	LogStream       = "infrapilot:logs"
	DiscoveryStream = "infrapilot:discovery"
	AIStream        = "infrapilot:ai"
	InventoryStream = "infrapilot:inventory"
)

// Message represents a queue message
type Message struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Payload    map[string]interface{} `json:"payload"`
	Timestamp  time.Time              `json:"timestamp"`
	Retries    int                    `json:"retries"`
	MaxRetries int                    `json:"max_retries"`
}

// Queue manages message queuing with Redis Streams
type Queue struct {
	redis *redis.Client
}

// NewQueue creates a new message queue
func NewQueue() *Queue {
	return &Queue{
		redis: cache.GetRedisClient(),
	}
}

// Publish publishes a message to a stream
func (q *Queue) Publish(streamName string, msg *Message) error {
	if q.redis == nil {
		return nil
	}

	msg.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	msg.Timestamp = time.Now()
	if msg.MaxRetries == 0 {
		msg.MaxRetries = 3
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	ctx := context.Background()
	err = q.redis.XAdd(ctx, &redis.XAddArgs{
		Stream: streamName,
		Values: map[string]interface{}{
			"data": data,
		},
	}).Err()

	if err != nil {
		return fmt.Errorf("failed to publish to stream %s: %w", streamName, err)
	}

	return nil
}

// Consume consumes messages from a stream
func (q *Queue) Consume(streamName, consumerGroup, consumerName string, batchSize int) ([]*Message, error) {
	if q.redis == nil {
		time.Sleep(2 * time.Second)
		return nil, nil
	}

	ctx := context.Background()

	// Create consumer group if it doesn't exist
	err := q.redis.XGroupCreateMkStream(ctx, streamName, consumerGroup, "0").Err()
	if err != nil && !isRedisStreamError(err) {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	// Read messages
	results, err := q.redis.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    consumerGroup,
		Consumer: consumerName,
		Streams:  []string{streamName, ">"},
		Count:    int64(batchSize),
		Block:    5 * time.Second,
	}).Result()

	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read from stream: %w", err)
	}

	messages := make([]*Message, 0)
	for _, streamResult := range results {
		for _, message := range streamResult.Messages {
			dataBytes, ok := message.Values["data"].([]byte)
			if !ok {
				continue
			}

			var msg Message
			if err := json.Unmarshal(dataBytes, &msg); err != nil {
				continue
			}

			messages = append(messages, &msg)
		}
	}

	return messages, nil
}

// Ack acknowledges a message
func (q *Queue) Ack(streamName, consumerGroup, messageID string) error {
	if q.redis == nil {
		return nil
	}
	ctx := context.Background()
	return q.redis.XAck(ctx, streamName, consumerGroup, messageID).Err()
}

// Nack negatively acknowledges a message (no-op in Redis Streams - message will be redelivered after timeout)
func (q *Queue) Nack(streamName, consumerGroup, messageID string) error {
	if q.redis == nil {
		return nil
	}
	return nil
}

// PublishToDLQ publishes a message to the dead letter queue
func (q *Queue) PublishToDLQ(streamName string, msg *Message) error {
	if q.redis == nil {
		return nil
	}
	dlqStream := fmt.Sprintf("%s:dlq", streamName)

	dlqMsg := &Message{
		ID:         msg.ID,
		Type:       "dlq." + msg.Type,
		Payload:    msg.Payload,
		Timestamp:  time.Now(),
		Retries:    msg.Retries,
		MaxRetries: msg.MaxRetries,
	}

	return q.Publish(dlqStream, dlqMsg)
}

// GetStreamLength returns the length of a stream
func (q *Queue) GetStreamLength(streamName string) (int64, error) {
	ctx := context.Background()
	info, err := q.redis.XLen(ctx, streamName).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get stream length: %w", err)
	}
	return info, nil
}

// CreateConsumerGroup creates a consumer group for a stream
func (q *Queue) CreateConsumerGroup(streamName, groupName string) error {
	ctx := context.Background()
	return q.redis.XGroupCreateMkStream(ctx, streamName, groupName, "0").Err()
}

func isRedisStreamError(err error) bool {
	return err != nil && err.Error() == "BUSYGROUP Consumer Group name already exists"
}
