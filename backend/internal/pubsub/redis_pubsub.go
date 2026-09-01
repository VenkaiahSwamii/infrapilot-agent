package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"infrapilot/backend/internal/cache"
	"infrapilot/backend/internal/websocket"

	"github.com/redis/go-redis/v9"
)

const (
	// Channel names
	MetricsChannel   = "metrics"
	AlertsChannel    = "alerts"
	EventsChannel    = "events"
	WebSocketChannel = "websocket"
)

// Message represents a pub/sub message
type Message struct {
	Channel string      `json:"channel"`
	Payload interface{} `json:"payload"`
	Source  string      `json:"source,omitempty"`
	Time    int64       `json:"timestamp"`
}

// PubSub manages Redis Pub/Sub for multi-instance communication
type PubSub struct {
	redisClient *redis.Client
}

// NewPubSub creates a new PubSub instance
func NewPubSub() *PubSub {
	return &PubSub{}
}

// Publish publishes a message to a Redis channel
func (p *PubSub) Publish(channel string, payload interface{}) error {
	msg := Message{
		Channel: channel,
		Payload: payload,
		Time:    time.Now().Unix(),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	client := cache.GetRedisClient()
	return client.Publish(context.Background(), channel, data).Err()
}

// Subscribe subscribes to a Redis channel
func (p *PubSub) Subscribe(channel string, handler func(Message) error) error {
	client := cache.GetRedisClient()
	pubsub := client.Subscribe(context.Background(), channel)

	// Receive confirmation
	_, err := pubsub.Receive(context.Background())
	if err != nil {
		return fmt.Errorf("failed to subscribe to channel %s: %w", channel, err)
	}

	// Start processing messages in a goroutine
	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			var message Message
			if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
				log.Printf("Failed to unmarshal pubsub message: %v", err)
				continue
			}

			if err := handler(message); err != nil {
				log.Printf("Error handling message on channel %s: %v", channel, err)
			}
		}
	}()

	log.Printf("Subscribed to Redis channel: %s", channel)
	return nil
}

// Initialize initializes pub/sub subscriptions for backend
func Initialize() error {
	ps := NewPubSub()

	// Subscribe to websocket channel for cross-instance broadcasting
	if err := ps.Subscribe(WebSocketChannel, func(msg Message) error {
		return BroadcastWebSocket(msg.Payload)
	}); err != nil {
		return err
	}

	// Subscribe to metrics channel
	if err := ps.Subscribe(MetricsChannel, func(msg Message) error {
		log.Printf("Received metrics message from source: %s", msg.Source)
		return nil
	}); err != nil {
		return err
	}

	// Subscribe to alerts channel
	if err := ps.Subscribe(AlertsChannel, func(msg Message) error {
		return BroadcastWebSocket(msg.Payload)
	}); err != nil {
		return err
	}

	log.Println("Redis Pub/Sub initialized successfully")
	return nil
}

// BroadcastWebSocket broadcasts a message to all connected WebSocket clients
func BroadcastWebSocket(payload interface{}) error {
	log.Printf("Broadcasting WebSocket message from Redis Pub/Sub to connected clients")
	if websocket.WS != nil {
		websocket.WS.BroadcastFromPubSub(payload)
	}
	return nil
}

// PublishWebSocket publishes a WebSocket message to all backend instances
func PublishWebSocket(payload interface{}) error {
	ps := NewPubSub()
	return ps.Publish(WebSocketChannel, payload)
}

// PublishMetrics publishes a metrics message
func PublishMetrics(payload interface{}) error {
	ps := NewPubSub()
	return ps.Publish(MetricsChannel, payload)
}

// PublishAlert publishes an alert message
func PublishAlert(payload interface{}) error {
	ps := NewPubSub()
	return ps.Publish(AlertsChannel, payload)
}
