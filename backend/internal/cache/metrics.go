package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"infrapilot/backend/internal/models"

	"github.com/redis/go-redis/v9"
)

const (
	MetricPrefix = "metric:"
	MetricTTL    = 15 * time.Minute
)

// CachedMetric represents a metric with metadata for caching
type CachedMetric struct {
	MachineID   string                 `json:"machine_id"`
	Hostname    string                 `json:"hostname"`
	IPAddress   string                 `json:"ip_address"`
	CPUUsage    float64                `json:"cpu_usage"`
	MemoryUsage float64                `json:"memory_usage"`
	DiskUsage   float64                `json:"disk_usage"`
	Tags        map[string]string      `json:"tags"`
	Timestamp   time.Time              `json:"timestamp"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

// SetMetric caches a single metric
func SetMetric(machineID string, metric *models.Metric) error {
	key := fmt.Sprintf("%s%s:%d", MetricPrefix, machineID, metric.CreatedAt.UnixNano())

	cached := CachedMetric{
		MachineID:   metric.MachineID.String(),
		Hostname:    metric.Hostname,
		IPAddress:   metric.IPAddress,
		CPUUsage:    metric.CPUUsage,
		MemoryUsage: metric.MemoryUsage,
		DiskUsage:   metric.DiskUsage,
		Timestamp:   metric.CreatedAt,
	}

	data, err := json.Marshal(cached)
	if err != nil {
		return fmt.Errorf("failed to marshal metric: %w", err)
	}

	client := GetRedisClient()
	return client.Set(context.Background(), key, data, MetricTTL).Err()
}

// GetMetricsByMachine retrieves all cached metrics for a machine
func GetMetricsByMachine(machineID string) ([]CachedMetric, error) {
	pattern := fmt.Sprintf("%s%s:*", MetricPrefix, machineID)

	client := GetRedisClient()
	keys, err := client.Keys(context.Background(), pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get metric keys: %w", err)
	}

	if len(keys) == 0 {
		return nil, nil
	}

	pipe := client.Pipeline()
	for _, key := range keys {
		pipe.Get(context.Background(), key)
	}

	results, err := pipe.Exec(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to pipeline get metrics: %w", err)
	}

	metrics := make([]CachedMetric, 0, len(results))
	for _, result := range results {
		if result.Err() != nil {
			continue
		}

		cmd, ok := result.(*redis.StringCmd)
		if !ok {
			continue
		}

		data, err := cmd.Result()
		if err != nil {
			continue
		}

		var metric CachedMetric
		if err := json.Unmarshal([]byte(data), &metric); err != nil {
			continue
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

// SetMetricBatch caches multiple metrics at once
func SetMetricBatch(metrics []models.Metric) error {
	if len(metrics) == 0 {
		return nil
	}

	client := GetRedisClient()
	pipe := client.Pipeline()

	for _, metric := range metrics {
		key := fmt.Sprintf("%s%s:%d", MetricPrefix, metric.MachineID.String(), metric.CreatedAt.UnixNano())

		cached := CachedMetric{
			MachineID:   metric.MachineID.String(),
			Hostname:    metric.Hostname,
			IPAddress:   metric.IPAddress,
			CPUUsage:    metric.CPUUsage,
			MemoryUsage: metric.MemoryUsage,
			DiskUsage:   metric.DiskUsage,
			Timestamp:   metric.CreatedAt,
		}

		data, err := json.Marshal(cached)
		if err != nil {
			continue
		}

		pipe.Set(context.Background(), key, data, MetricTTL)
	}

	_, err := pipe.Exec(context.Background())
	if err != nil {
		return fmt.Errorf("failed to batch cache metrics: %w", err)
	}

	return nil
}

// DeleteMetricsByMachine removes all cached metrics for a machine
func DeleteMetricsByMachine(machineID string) error {
	pattern := fmt.Sprintf("%s%s:*", MetricPrefix, machineID)

	client := GetRedisClient()
	keys, err := client.Keys(context.Background(), pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get metric keys: %w", err)
	}

	if len(keys) == 0 {
		return nil
	}

	return client.Del(context.Background(), keys...).Err()
}

// InvalidateMetricCache clears all metric cache entries
func InvalidateMetricCache() error {
	pattern := MetricPrefix + "*"

	client := GetRedisClient()
	keys, err := client.Keys(context.Background(), pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get metric keys: %w", err)
	}

	if len(keys) == 0 {
		return nil
	}

	return client.Del(context.Background(), keys...).Err()
}
