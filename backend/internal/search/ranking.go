package search

import (
	"sort"
	"strings"
	"time"
)

// CalculateRankScore computes a relevance score for a search result given a query string
func CalculateRankScore(item SearchIndex, query string) float64 {
	score := 0.0
	lowerQuery := strings.ToLower(strings.TrimSpace(query))
	if lowerQuery == "" {
		return item.RankScore
	}

	lowerTitle := strings.ToLower(item.Title)
	lowerDesc := strings.ToLower(item.Description)

	// 1. Exact Match Boost
	if lowerTitle == lowerQuery {
		score += 100.0
	} else if strings.HasPrefix(lowerTitle, lowerQuery) {
		score += 50.0
	} else if strings.Contains(lowerTitle, lowerQuery) {
		score += 30.0
	}

	if strings.Contains(lowerDesc, lowerQuery) {
		score += 10.0
	}

	// 2. Keyword Match Boost
	for _, kw := range item.Keywords {
		lowerKW := strings.ToLower(kw)
		if lowerKW == lowerQuery {
			score += 25.0
		} else if strings.Contains(lowerKW, lowerQuery) {
			score += 10.0
		}
	}

	// 3. Active Incidents and Critical Alerts Priority Boost
	if item.ResourceType == SearchResultTypeIncident && (item.Status == "OPEN" || item.Status == "ACTIVE") {
		score += 35.0
	}
	if item.ResourceType == SearchResultTypeAlert && (item.Severity == "CRITICAL" || item.Severity == "P1") {
		score += 25.0
	}

	// 4. Recency Boost
	if !item.UpdatedAt.IsZero() {
		hoursOld := time.Since(item.UpdatedAt).Hours()
		if hoursOld <= 1 {
			score += 15.0
		} else if hoursOld <= 24 {
			score += 10.0
		} else if hoursOld <= 168 { // 7 days
			score += 5.0
		}
	}

	return score
}

// RankResults sorts a slice of SearchIndex entries in descending order of calculated rank score
func RankResults(results []SearchIndex, query string) []SearchIndex {
	for i := range results {
		results[i].RankScore = CalculateRankScore(results[i], query)
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].RankScore != results[j].RankScore {
			return results[i].RankScore > results[j].RankScore
		}
		return results[i].UpdatedAt.After(results[j].UpdatedAt)
	})

	return results
}
