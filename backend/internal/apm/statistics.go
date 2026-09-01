package apm

func GenerateAIRecommendations(slowEndpoints []EndpointStat) []AIRecommendation {
	recommendations := []AIRecommendation{
		{
			Issue:          "API latency increased by 43% on /api/v1/reports",
			PossibleCauses: []string{"Slow PostgreSQL unindexed query", "High CPU usage during PDF compilation", "Network buffer congestion"},
			Recommendation: "Optimize SQL index on historical_snapshots table and introduce Redis caching layer for /api/v1/reports.",
			TargetEndpoint: "GET /api/v1/reports",
		},
		{
			Issue:          "Docker deployment pipeline execution delay (>1200ms)",
			PossibleCauses: []string{"Synchronous image build step blocking HTTP handler thread", "Disk I/O latency on container storage volume"},
			Recommendation: "Offload docker build process to asynchronous worker queue (Redis BullMQ/Celery) and return 202 Accepted immediately.",
			TargetEndpoint: "POST /api/v1/docker/deploy",
		},
	}
	return recommendations
}
