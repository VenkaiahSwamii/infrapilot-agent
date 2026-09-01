package workers

import (
	"context"
	"encoding/json"
	"fmt"

	"infrapilot/backend/internal/cache"
	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"
	"infrapilot/backend/internal/queue"
)

// ProcessMetric processes a metric message
func ProcessMetric(ctx context.Context, msg *queue.Message) error {
	if msg.Payload == nil {
		return fmt.Errorf("metric payload is nil")
	}

	data, err := json.Marshal(msg.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal metric payload: %w", err)
	}

	var metric models.Metric
	if err := json.Unmarshal(data, &metric); err != nil {
		return fmt.Errorf("failed to unmarshal metric: %w", err)
	}

	// Store in database
	if err := database.DB.Create(&metric).Error; err != nil {
		return fmt.Errorf("failed to store metric in database: %w", err)
	}

	// Update machine registry with latest live data
	err = database.DB.Model(&models.Machine{}).
		Where("id = ?", metric.MachineID).
		Updates(map[string]interface{}{
			"hostname":   metric.Hostname,
			"ip_address": metric.IPAddress,
			"os":         metric.OS,
			"status":     "ONLINE",
			"last_seen":  metric.CreatedAt,
			"updated_at": metric.CreatedAt,
		}).Error

	if err != nil {
		return fmt.Errorf("failed to update machine: %w", err)
	}

	// Cache in Redis for fast access
	if err := cache.SetMetric(metric.MachineID.String(), &metric); err != nil {
		fmt.Printf("Warning: Failed to cache metric: %v\n", err)
	}

	return nil
}

// ProcessAlert processes an alert message
func ProcessAlert(ctx context.Context, msg *queue.Message) error {
	if msg.Payload == nil {
		return fmt.Errorf("alert payload is nil")
	}

	data, err := json.Marshal(msg.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal alert payload: %w", err)
	}

	var alert models.AlertRule
	if err := json.Unmarshal(data, &alert); err != nil {
		return fmt.Errorf("failed to unmarshal alert: %w", err)
	}

	// Store in database
	if err := database.DB.Create(&alert).Error; err != nil {
		return fmt.Errorf("failed to store alert in database: %w", err)
	}

	return nil
}

// ProcessLog processes a log message
func ProcessLog(ctx context.Context, msg *queue.Message) error {
	if msg.Payload == nil {
		return fmt.Errorf("log payload is nil")
	}

	data, err := json.Marshal(msg.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal log payload: %w", err)
	}

	var log models.LinuxLog
	if err := json.Unmarshal(data, &log); err != nil {
		return fmt.Errorf("failed to unmarshal log: %w", err)
	}

	// Store in database
	if err := database.DB.Create(&log).Error; err != nil {
		return fmt.Errorf("failed to store log in database: %w", err)
	}

	return nil
}

// ProcessDiscovery processes a discovery message
func ProcessDiscovery(ctx context.Context, msg *queue.Message) error {
	if msg.Payload == nil {
		return fmt.Errorf("discovery payload is nil")
	}

	data, err := json.Marshal(msg.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal discovery payload: %w", err)
	}

	var discovery map[string]interface{}
	if err := json.Unmarshal(data, &discovery); err != nil {
		return fmt.Errorf("failed to unmarshal discovery: %w", err)
	}

	// Process discovery data based on type
	discoveryType, ok := discovery["type"].(string)
	if !ok {
		return fmt.Errorf("discovery type not found")
	}

	switch discoveryType {
	case "machine":
		// Process machine discovery
		fmt.Printf("Processing machine discovery: %v\n", discovery)
	case "service":
		// Process service discovery
		fmt.Printf("Processing service discovery: %v\n", discovery)
	case "network":
		// Process network discovery
		fmt.Printf("Processing network discovery: %v\n", discovery)
	default:
		fmt.Printf("Unknown discovery type: %s\n", discoveryType)
	}

	return nil
}

// ProcessAI processes an AI message
func ProcessAI(ctx context.Context, msg *queue.Message) error {
	if msg.Payload == nil {
		return fmt.Errorf("AI payload is nil")
	}

	data, err := json.Marshal(msg.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal AI payload: %w", err)
	}

	var aiData map[string]interface{}
	if err := json.Unmarshal(data, &aiData); err != nil {
		return fmt.Errorf("failed to unmarshal AI: %w", err)
	}

	// Process AI insights
	fmt.Printf("Processing AI insight: %v\n", aiData)

	return nil
}

// ProcessInventory processes an inventory message
func ProcessInventory(ctx context.Context, msg *queue.Message) error {
	if msg.Payload == nil {
		return fmt.Errorf("inventory payload is nil")
	}

	data, err := json.Marshal(msg.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal inventory payload: %w", err)
	}

	var inventory map[string]interface{}
	if err := json.Unmarshal(data, &inventory); err != nil {
		return fmt.Errorf("failed to unmarshal inventory: %w", err)
	}

	// Process inventory data
	fmt.Printf("Processing inventory update: %v\n", inventory)

	return nil
}
