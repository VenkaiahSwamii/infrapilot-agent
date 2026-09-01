package ai

type CapacityPrediction struct {
	Resource           string  `json:"resource"`
	CurrentUsagePct    float64 `json:"current_usage_pct"`
	GrowthRatePerDay   float64 `json:"growth_rate_per_day"`
	DaysUntilDepletion int     `json:"days_until_depletion"`
	AlertSeverity      string  `json:"alert_severity"`
	Recommendation     string  `json:"recommendation"`
}

func GenerateCapacityPredictions() []CapacityPrediction {
	return []CapacityPrediction{
		{
			Resource:           "Storage Capacity (/var/log)",
			CurrentUsagePct:    89.2,
			GrowthRatePerDay:   2.7,
			DaysUntilDepletion: 4,
			AlertSeverity:      "HIGH",
			Recommendation:     "Rotate system logs and purge backups older than 7 days.",
		},
		{
			Resource:           "PostgreSQL Database Disk (/var/lib/postgresql)",
			CurrentUsagePct:    74.5,
			GrowthRatePerDay:   1.2,
			DaysUntilDepletion: 21,
			AlertSeverity:      "MEDIUM",
			Recommendation:     "Provision additional 50GB EBS storage volume.",
		},
	}
}
