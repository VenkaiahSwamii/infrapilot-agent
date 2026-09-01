package analytics

import (
	"time"
)

type TrendPoint struct {
	Timestamp   string  `json:"timestamp"`
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskUsage   float64 `json:"disk_usage"`
	NetworkIn   int64   `json:"network_in"`
	NetworkOut  int64   `json:"network_out"`
	AlertCount  int     `json:"alert_count"`
}

type TrendEngine struct {
	repo Repository
}

func NewTrendEngine(repo Repository) *TrendEngine {
	return &TrendEngine{repo: repo}
}

func (t *TrendEngine) GetTrends(orgID string, rangeStr string) ([]TrendPoint, error) {
	since := time.Now().Add(-24 * time.Hour)
	switch rangeStr {
	case "1h":
		since = time.Now().Add(-1 * time.Hour)
	case "1w":
		since = time.Now().Add(-7 * 24 * time.Hour)
	case "1m":
		since = time.Now().Add(-30 * 24 * time.Hour)
	case "1y":
		since = time.Now().Add(-365 * 24 * time.Hour)
	}

	history, err := t.repo.GetRecentPerformanceHistory(orgID, since)
	if err != nil {
		return nil, err
	}

	var points []TrendPoint
	for _, h := range history {
		points = append(points, TrendPoint{
			Timestamp:   h.Timestamp.Format("15:04:05"),
			CPUUsage:    h.CPUUsage,
			MemoryUsage: h.MemoryUsage,
			DiskUsage:   h.DiskUsage,
			NetworkIn:   h.NetworkIn,
			NetworkOut:  h.NetworkOut,
			AlertCount:  int(h.CPUUsage/20.0) + int(h.MemoryUsage/30.0),
		})
	}

	return points, nil
}
