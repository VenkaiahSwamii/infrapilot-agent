package ai

type IntelligentRecommendation struct {
	ID                  string   `json:"id"`
	Problem             string   `json:"problem"`
	Actions             []string `json:"actions"`
	TargetComponent     string   `json:"target_component"`
	ExpectedImprovement string   `json:"expected_improvement"`
	Priority            string   `json:"priority"`
}

func GenerateRecommendations() []IntelligentRecommendation {
	return []IntelligentRecommendation{
		{
			ID:                  "rec-101",
			Problem:             "High Storage Usage (/var/log at 95%)",
			Actions:             []string{"Delete old archived logs", "Configure logrotate compression", "Expand EBS storage volume", "Archive backups to AWS S3 Glacier"},
			TargetComponent:     "Linux Storage Subsystem",
			ExpectedImprovement: "Reduces disk usage by 35% and prevents disk full crash",
			Priority:            "CRITICAL",
		},
		{
			ID:                  "rec-102",
			Problem:             "PostgreSQL Connection Pool Exhaustion (98/100 connections)",
			Actions:             []string{"Increase max_connections in postgresql.conf to 200", "Deploy PgBouncer connection pooler proxy", "Tune idle transaction timeout"},
			TargetComponent:     "PostgreSQL Database Engine",
			ExpectedImprovement: "Eliminates API 500 error spikes during traffic spikes",
			Priority:            "HIGH",
		},
	}
}
