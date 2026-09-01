package apm

import "infrapilot/backend/internal/models"

type OverviewStats struct {
	APIHealth      float64 `json:"api_health"`         // e.g. 99.98%
	AverageLatency float64 `json:"average_latency_ms"` // e.g. 58 ms
	RequestsPerSec float64 `json:"requests_per_sec"`   // e.g. 42 rps
	ErrorRate      float64 `json:"error_rate"`         // e.g. 0.8%
	SlowAPICount   int64   `json:"slow_api_count"`     // e.g. 3
	TotalRequests  int64   `json:"total_requests"`
	TotalErrors    int64   `json:"total_errors"`
}

type EndpointStat struct {
	Method       string  `json:"method"`
	Endpoint     string  `json:"endpoint"`
	AvgLatencyMs float64 `json:"avg_latency_ms"`
	RequestCount int64   `json:"request_count"`
	Percentage   float64 `json:"percentage"`
	IsSlow       bool    `json:"is_slow"`
}

type AIRecommendation struct {
	Issue          string   `json:"issue"`
	PossibleCauses []string `json:"possible_causes"`
	Recommendation string   `json:"recommendation"`
	TargetEndpoint string   `json:"target_endpoint"`
}

type APMDataResponse struct {
	Overview        OverviewStats          `json:"overview"`
	Endpoints       []EndpointStat         `json:"endpoints"`
	SlowEndpoints   []EndpointStat         `json:"slow_endpoints"`
	TopAPIs         []EndpointStat         `json:"top_apis"`
	Recommendations []AIRecommendation     `json:"recommendations"`
	Services        []models.ServiceMetric `json:"services"`
}
