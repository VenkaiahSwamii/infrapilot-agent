package services

import (
	"encoding/json"
	"strings"
	"time"

	"infrapilot/backend/internal/models"
	"infrapilot/backend/internal/repository"

	"github.com/google/uuid"
)

// SearchService handles search business logic
type SearchService struct {
	searchRepo *repository.SearchRepository
}

// NewSearchService creates a new search service
func NewSearchService(searchRepo *repository.SearchRepository) *SearchService {
	return &SearchService{
		searchRepo: searchRepo,
	}
}

// SearchRequest represents a search request
type SearchRequest struct {
	Query          string
	Category       string
	Severity       string
	Status         string
	OrganizationID uuid.UUID
	ResourceType   string
	DateFrom       *time.Time
	DateTo         *time.Time
	Page           int
	PageSize       int
	UserID         uuid.UUID
}

// SearchResult represents search results
type SearchResult struct {
	Query       string                         `json:"query"`
	TotalCount  int64                          `json:"total_count"`
	Page        int                            `json:"page"`
	PageSize    int                            `json:"page_size"`
	TotalPages  int                            `json:"total_pages"`
	Results     []models.SearchIndex           `json:"results"`
	Suggestions []string                       `json:"suggestions"`
	Facets      map[string][]models.FacetCount `json:"facets"`
	TookMs      int64                          `json:"took_ms"`
	Summary     string                         `json:"summary,omitempty"`
}

// Search performs a search query
func (s *SearchService) Search(req SearchRequest) (*SearchResult, error) {
	start := time.Now()

	// Build search options
	opts := repository.SearchOptions{
		Query:          req.Query,
		OrganizationID: req.OrganizationID,
		Category:       req.Category,
		Severity:       req.Severity,
		Status:         req.Status,
		ResourceType:   req.ResourceType,
		DateFrom:       req.DateFrom,
		DateTo:         req.DateTo,
		Page:           req.Page,
		PageSize:       req.PageSize,
		SortBy:         "updated_at",
		SortOrder:      "DESC",
	}

	// Execute search
	result, err := s.searchRepo.Search(opts)
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := 0
	if result.PageSize > 0 {
		totalPages = int(result.Total) / result.PageSize
		if int(result.Total)%result.PageSize > 0 {
			totalPages++
		}
	}

	// Convert repository results to models
	results := make([]models.SearchIndex, 0, len(result.Results))
	for _, r := range result.Results {
		results = append(results, models.SearchIndex{
			ID:             r.ID,
			ResourceType:   models.SearchResultType(r.ResourceType),
			ResourceID:     r.ResourceID,
			Title:          r.Title,
			Description:    r.Description,
			Keywords:       r.Keywords,
			OrganizationID: r.OrganizationID,
			Category:       r.Category,
			Severity:       r.Severity,
			Status:         r.Status,
			URL:            r.URL,
			Metadata:       r.Metadata,
			RankScore:      r.RankScore,
			UpdatedAt:      r.UpdatedAt,
			CreatedAt:      r.CreatedAt,
		})
	}

	// Convert facets
	facets := make(map[string][]models.FacetCount)
	for k, v := range result.Facets {
		facets[k] = make([]models.FacetCount, 0, len(v))
		for _, f := range v {
			facets[k] = append(facets[k], models.FacetCount{Value: f.Value, Count: f.Count})
		}
	}

	// Track analytics asynchronously
	go func() {
		_ = s.searchRepo.TrackSearchAnalytics(req.UserID, req.OrganizationID, req.Query, int(result.Total), nil, int(time.Since(start).Milliseconds()))
	}()

	// Add to recent searches
	go func() {
		_ = s.searchRepo.AddRecentSearch(req.UserID, req.OrganizationID, req.Query, int(result.Total))
	}()

	return &SearchResult{
		Query:       req.Query,
		TotalCount:  result.Total,
		Page:        result.Page,
		PageSize:    result.PageSize,
		TotalPages:  totalPages,
		Results:     results,
		Suggestions: []string{},
		Facets:      facets,
		TookMs:      result.TookMs,
	}, nil
}

// GetSuggestions returns search suggestions
func (s *SearchService) GetSuggestions(query string, orgID uuid.UUID) ([]string, error) {
	return s.searchRepo.GetSuggestions(query, orgID, 10)
}

// GetRecentSearches returns recent searches for a user
func (s *SearchService) GetRecentSearches(userID, orgID uuid.UUID) ([]models.RecentSearch, error) {
	return s.searchRepo.GetRecentSearches(userID, orgID, 10)
}

// GetPopularSearches returns popular searches
func (s *SearchService) GetPopularSearches(orgID uuid.UUID) ([]models.PopularSearch, error) {
	return s.searchRepo.GetPopularSearches(orgID, 10)
}

// SaveSearch saves a search query
func (s *SearchService) SaveSearch(userID, orgID uuid.UUID, name string, query string, filters map[string]interface{}) error {
	filterData, err := json.Marshal(filters)
	if err != nil {
		return err
	}
	return s.searchRepo.SaveSearch(userID, orgID, name, query, filterData)
}

// GetSavedSearches returns saved searches
func (s *SearchService) GetSavedSearches(userID, orgID uuid.UUID) ([]models.SavedSearch, error) {
	return s.searchRepo.GetSavedSearches(userID, orgID)
}

// DeleteSavedSearch deletes a saved search
func (s *SearchService) DeleteSavedSearch(id, userID, orgID uuid.UUID) error {
	return s.searchRepo.DeleteSavedSearch(id, userID, orgID)
}

// AISearch performs AI-powered search
func (s *SearchService) AISearch(userID, orgID uuid.UUID, query string) (*SearchResult, error) {
	start := time.Now()

	// Parse natural language query
	searchTerms := s.parseNaturalLanguage(query)

	// Build search options
	opts := repository.SearchOptions{
		Query:          searchTerms,
		OrganizationID: orgID,
		Page:           1,
		PageSize:       20,
		SortBy:         "updated_at",
		SortOrder:      "DESC",
	}

	// Execute search
	result, err := s.searchRepo.Search(opts)
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := 0
	if result.PageSize > 0 {
		totalPages = int(result.Total) / result.PageSize
		if int(result.Total)%result.PageSize > 0 {
			totalPages++
		}
	}

	// Convert results
	results := make([]models.SearchIndex, 0, len(result.Results))
	for _, r := range result.Results {
		results = append(results, models.SearchIndex{
			ID:             r.ID,
			ResourceType:   models.SearchResultType(r.ResourceType),
			ResourceID:     r.ResourceID,
			Title:          r.Title,
			Description:    r.Description,
			Keywords:       r.Keywords,
			OrganizationID: r.OrganizationID,
			Category:       r.Category,
			Severity:       r.Severity,
			Status:         r.Status,
			URL:            r.URL,
			Metadata:       r.Metadata,
			RankScore:      r.RankScore,
			UpdatedAt:      r.UpdatedAt,
			CreatedAt:      r.CreatedAt,
		})
	}

	// Convert facets
	facets := make(map[string][]models.FacetCount)
	for k, v := range result.Facets {
		facets[k] = make([]models.FacetCount, 0, len(v))
		for _, f := range v {
			facets[k] = append(facets[k], models.FacetCount{Value: f.Value, Count: f.Count})
		}
	}

	// Generate summary
	summary := s.generateSummary(query, results)

	return &SearchResult{
		Query:       query,
		TotalCount:  result.Total,
		Page:        result.Page,
		PageSize:    result.PageSize,
		TotalPages:  totalPages,
		Results:     results,
		Suggestions: []string{},
		Facets:      facets,
		TookMs:      time.Since(start).Milliseconds(),
		Summary:     summary,
	}, nil
}

// parseNaturalLanguage parses natural language queries
func (s *SearchService) parseNaturalLanguage(query string) string {
	// Simple NLU - extract key terms
	// In production, this would use the AI engine
	query = strings.ToLower(query)

	// Remove common words
	stopWords := map[string]bool{
		"what": true, "which": true, "where": true, "when": true, "why": true, "how": true,
		"is": true, "are": true, "the": true, "a": true, "an": true, "show": true, "get": true,
		"find": true, "list": true, "all": true, "of": true, "for": true, "in": true, "on": true,
	}

	terms := strings.Fields(query)
	var filtered []string
	for _, term := range terms {
		if !stopWords[term] {
			filtered = append(filtered, term)
		}
	}

	return strings.Join(filtered, " ")
}

// generateSummary generates a summary of search results
func (s *SearchService) generateSummary(query string, results []models.SearchIndex) string {
	if len(results) == 0 {
		return "No results found for your query."
	}

	// Group by resource type
	typeCounts := make(map[models.SearchResultType]int)
	for _, r := range results {
		typeCounts[r.ResourceType]++
	}

	// Build summary
	var parts []string
	for rType := range typeCounts {
		parts = append(parts, string(rType))
	}

	return "Found " + string(rune(len(results))) + " results: " + strings.Join(parts, ", ")
}

// IndexResource indexes a resource for search
func (s *SearchService) IndexResource(resourceType string, resourceID uuid.UUID, title string, description string, keywords []string, orgID uuid.UUID, category string, severity string, status string, url string, metadata map[string]interface{}) error {
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		metadataBytes = []byte("{}")
	}

	return s.searchRepo.Index(resourceType, resourceID, title, description, keywords, orgID, category, severity, status, url, metadataBytes)
}

// RemoveResource removes a resource from search index
func (s *SearchService) RemoveResource(resourceType string, resourceID uuid.UUID, orgID uuid.UUID) error {
	return s.searchRepo.RemoveFromIndex(resourceType, resourceID, orgID)
}

// RefreshIndex refreshes the search index for a resource
func (s *SearchService) RefreshIndex(resourceType string, resourceID uuid.UUID, orgID uuid.UUID) error {
	// In production, this would re-fetch the resource data and update the index
	// For now, just mark as updated
	return nil
}
