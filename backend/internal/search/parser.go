package search

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ParsedQuery represents the structured interpretation of a search string
type ParsedQuery struct {
	RawQuery       string
	Text           string
	ResourceType   string
	Category       string
	Severity       string
	Status         string
	MetricName     string
	MetricOperator string
	MetricValue    float64
	IsMetricQuery  bool
	IsAIQuery      bool
	Keywords       []string
}

var (
	metricRegex = regexp.MustCompile(`(?i)(CPU|Memory|Disk|Network)\s*([><=]+)\s*(\d+(?:\.\d+)?)\%?`)
	scopeRegex  = regexp.MustCompile(`(?i)^(pod|container|image|node|cluster|deployment|service|namespace|log|incident|alert)\s+(.+)`)
)

// ParseQuery parses a raw query string into a structured ParsedQuery
func ParseQuery(raw string) ParsedQuery {
	trimmed := strings.TrimSpace(raw)
	parsed := ParsedQuery{
		RawQuery: raw,
		Text:     trimmed,
		Keywords: strings.Fields(strings.ToLower(trimmed)),
	}

	if trimmed == "" {
		return parsed
	}

	// 1. Check for Metric Query (e.g. CPU > 90%)
	if matches := metricRegex.FindStringSubmatch(trimmed); len(matches) == 4 {
		metricName := strings.ToUpper(matches[1])
		operator := matches[2]
		val, err := strconv.ParseFloat(matches[3], 64)
		if err == nil {
			parsed.IsMetricQuery = true
			parsed.MetricName = metricName
			parsed.MetricOperator = operator
			parsed.MetricValue = val
			parsed.Text = metricName
			return parsed
		}
	}

	// 2. Check for Scope Query (e.g. "pod nginx", "container app-api")
	if matches := scopeRegex.FindStringSubmatch(trimmed); len(matches) == 3 {
		scope := strings.ToLower(matches[1])
		term := strings.TrimSpace(matches[2])
		parsed.Text = term

		switch scope {
		case "pod", "node", "cluster", "deployment", "namespace":
			parsed.ResourceType = string(SearchResultTypeKubernetes)
			parsed.Category = scope
		case "container", "image":
			parsed.ResourceType = string(SearchResultTypeDocker)
			parsed.Category = scope
		case "log":
			parsed.ResourceType = string(SearchResultTypeLog)
		case "incident":
			parsed.ResourceType = string(SearchResultTypeIncident)
		case "alert":
			parsed.ResourceType = string(SearchResultTypeAlert)
		}
		return parsed
	}

	// 3. Keyword / Categorical Mapping
	lower := strings.ToLower(trimmed)
	switch lower {
	case "docker":
		parsed.ResourceType = string(SearchResultTypeDocker)
	case "kubernetes", "k8s":
		parsed.ResourceType = string(SearchResultTypeKubernetes)
	case "incidents", "incident":
		parsed.ResourceType = string(SearchResultTypeIncident)
	case "alerts", "alert":
		parsed.ResourceType = string(SearchResultTypeAlert)
	case "logs", "log":
		parsed.ResourceType = string(SearchResultTypeLog)
	case "machines", "machine", "server", "servers":
		parsed.ResourceType = string(SearchResultTypeMachine)
	}

	// 4. Check for Natural Language AI search queries
	if isNaturalLanguage(lower) {
		parsed.IsAIQuery = true
		interpretAIQuery(lower, &parsed)
	}

	return parsed
}

func isNaturalLanguage(lower string) bool {
	aiPrefixes := []string{
		"which ", "why is", "show ", "what ", "how many", "list ", "find ", "where ",
	}
	for _, prefix := range aiPrefixes {
		if strings.HasPrefix(lower, prefix) {
			return true
		}
	}
	return strings.Contains(lower, "unhealthy") || strings.Contains(lower, "failing") || strings.Contains(lower, "error") || strings.Contains(lower, "slow")
}

func interpretAIQuery(lower string, p *ParsedQuery) {
	if strings.Contains(lower, "unhealthy") || strings.Contains(lower, "offline") || strings.Contains(lower, "failing") {
		p.Status = "CRITICAL"
		if strings.Contains(lower, "server") || strings.Contains(lower, "machine") {
			p.ResourceType = string(SearchResultTypeMachine)
		} else if strings.Contains(lower, "pod") {
			p.ResourceType = string(SearchResultTypeKubernetes)
			p.Category = "pod"
		}
	}
	if strings.Contains(lower, "incident") {
		p.ResourceType = string(SearchResultTypeIncident)
	}
	if strings.Contains(lower, "alert") {
		p.ResourceType = string(SearchResultTypeAlert)
	}
	if strings.Contains(lower, "docker") {
		p.ResourceType = string(SearchResultTypeDocker)
	}
}

// GenerateAISummary creates a concise explanation summary based on search results
func GenerateAISummary(query string, results []SearchIndex) string {
	if len(results) == 0 {
		return fmt.Sprintf("No relevant resources found matching natural language query: '%s'.", query)
	}

	counts := make(map[SearchResultType]int)
	criticalCount := 0
	for _, r := range results {
		counts[r.ResourceType]++
		if r.Severity == "P1" || r.Severity == "CRITICAL" || r.Status == "CRITICAL" || r.Status == "OFFLINE" {
			criticalCount++
		}
	}

	var parts []string
	for rType, count := range counts {
		parts = append(parts, fmt.Sprintf("%d %s(s)", count, rType))
	}

	summary := fmt.Sprintf("Found %d total matching entity(ies) [%s] for query '%s'.", len(results), strings.Join(parts, ", "), query)
	if criticalCount > 0 {
		summary += fmt.Sprintf(" ⚠️ Attention: %d resource(s) require immediate remediation.", criticalCount)
	}
	return summary
}
