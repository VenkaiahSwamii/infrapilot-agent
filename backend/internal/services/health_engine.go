package services

import (
	"fmt"
)

type HealthEngine struct{}

func NewHealthEngine() *HealthEngine {
	return &HealthEngine{}
}

func CalculateHealthScore(
	cpu float64,
	memory float64,
	disk float64,
	latency float64,
	packetLoss float64,
) float64 {
	score := 100.0

	// CPU Deduction Rules
	if cpu > 90 {
		score -= 20
	} else if cpu > 75 {
		score -= 10
	}

	// Memory Deduction Rules
	if memory > 90 {
		score -= 20
	} else if memory > 75 {
		score -= 10
	}

	// Disk Deduction Rules
	if disk > 90 {
		score -= 20
	} else if disk > 75 {
		score -= 10
	}

	// Latency Deduction Rule
	if latency > 200 {
		score -= 15
	}

	// Packet Loss Deduction Rule
	if packetLoss > 2 {
		score -= 15
	}

	if score < 0 {
		score = 0
	}

	return score
}

func GetHealthScoreRating(score float64) (rating string, category string, stars string) {
	if score >= 90 {
		return "Green", "Healthy", "★★★★★ (Excellent)"
	} else if score >= 70 {
		return "Yellow", "Warning", "★★★☆☆ (Needs Attention)"
	} else {
		return "Red", "Critical", "★☆☆☆☆ (Action Required)"
	}
}

func GenerateHealthRecommendations(
	cpu float64,
	memory float64,
	disk float64,
	latency float64,
	packetLoss float64,
) []string {
	var recs []string

	if cpu > 90 {
		recs = append(recs, "CPU > 90% → Scale CPU resources or investigate high-load runaway processes.")
	} else if cpu > 75 {
		recs = append(recs, "CPU Elevated > 75% → Review active process CPU utilization.")
	}

	if memory > 90 {
		recs = append(recs, "Memory > 90% → Increase memory allocation or optimize running applications.")
	} else if memory > 75 {
		recs = append(recs, "Memory Elevated > 75% → Free application cache or expand RAM.")
	}

	if disk > 90 {
		recs = append(recs, "Disk > 90% → Clean up disk space (/tmp, logs) or expand storage capacity.")
	} else if disk > 75 {
		recs = append(recs, "Disk Approaching Capacity > 75% → Purge archived log files.")
	}

	if latency > 200 {
		recs = append(recs, "Latency > 200 ms → Investigate network connectivity or routing issues.")
	}

	if packetLoss > 2 {
		recs = append(recs, "Packet loss > 2% → Check network interfaces and physical/virtual connections.")
	}

	if len(recs) == 0 {
		recs = append(recs, fmt.Sprintf("System Operating Normally: All metrics are within optimal thresholds."))
	}

	return recs
}
