package services

import (
	"testing"
)

func TestHealthEngine_CalculateHealthScore(t *testing.T) {
	// 1. Perfect Health Score Test
	score1 := CalculateHealthScore(25.0, 40.0, 50.0, 30.0, 0.0)
	if score1 != 100.0 {
		t.Errorf("Expected health score 100.0 for low metrics, got %.1f", score1)
	}

	color1, category1, stars1 := GetHealthScoreRating(score1)
	if color1 != "Green" || category1 != "Healthy" {
		t.Errorf("Expected Green/Healthy for 100.0 score, got %s/%s (%s)", color1, category1, stars1)
	}

	// 2. High CPU & High Memory Deductions Test (>90 => -20 each)
	score2 := CalculateHealthScore(95.0, 92.0, 50.0, 30.0, 0.0)
	if score2 != 60.0 {
		t.Errorf("Expected health score 60.0 for high CPU & RAM, got %.1f", score2)
	}

	color2, category2, _ := GetHealthScoreRating(score2)
	if color2 != "Red" || category2 != "Critical" {
		t.Errorf("Expected Red/Critical for 60.0 score, got %s/%s", color2, category2)
	}

	// 3. AI Health Recommendations Test
	recs := GenerateHealthRecommendations(95.0, 92.0, 95.0, 250.0, 5.0)
	if len(recs) < 4 {
		t.Errorf("Expected at least 4 health recommendations for high metrics, got %d", len(recs))
	}
}
