package search

import (
	"time"

	"github.com/google/uuid"
)

// SearchService handles business logic for enterprise search
type SearchService struct {
	repo    *SearchRepository
	indexer *SearchIndexer
}

// NewSearchService creates a new search service
func NewSearchService(repo *SearchRepository, indexer *SearchIndexer) *SearchService {
	return &SearchService{
		repo:    repo,
		indexer: indexer,
	}
}

// Search executes structured search with query parsing, ranking, and analytics tracking
func (s *SearchService) Search(req SearchRequest) (*SearchResponse, error) {
	start := time.Now()

	// Ensure seed/reindexed data exists for org
	_ = s.indexer.ReindexAll(req.OrganizationID)

	// Parse query
	parsed := ParseQuery(req.Query)

	// Execute DB Search
	dbResult, err := s.repo.Search(req, parsed)
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if req.PageSize > 0 {
		totalPages = int(dbResult.Total) / req.PageSize
		if int(dbResult.Total)%req.PageSize > 0 {
			totalPages++
		}
	} else {
		req.PageSize = 20
	}

	// Fetch suggestions for autocomplete
	suggestions, _ := s.repo.GetSuggestions(req.Query, req.OrganizationID, 5)

	// Track analytics & recent searches asynchronously
	go func() {
		_ = s.repo.TrackSearchAnalytics(req.UserID, req.OrganizationID, req.Query, int(dbResult.Total), nil, int(time.Since(start).Milliseconds()))
		_ = s.repo.AddRecentSearch(req.UserID, req.OrganizationID, req.Query, int(dbResult.Total))
	}()

	return &SearchResponse{
		Query:       req.Query,
		TotalCount:  dbResult.Total,
		Page:        req.Page,
		PageSize:    req.PageSize,
		TotalPages:  totalPages,
		Results:     dbResult.Results,
		Suggestions: suggestions,
		Facets:      dbResult.Facets,
		TookMs:      time.Since(start).Milliseconds(),
	}, nil
}

// AISearch performs natural language query parsing and returns AI summarized response
func (s *SearchService) AISearch(userID, orgID uuid.UUID, query string) (*AISearchResponse, error) {
	start := time.Now()

	// Ensure index exists
	_ = s.indexer.ReindexAll(orgID)

	req := SearchRequest{
		Query:          query,
		OrganizationID: orgID,
		Page:           1,
		PageSize:       20,
		UserID:         userID,
	}

	res, err := s.Search(req)
	if err != nil {
		return nil, err
	}

	summary := GenerateAISummary(query, res.Results)

	return &AISearchResponse{
		Query:          query,
		Interpretation: fmtInterpretation(query, res.Results),
		Results:        res.Results,
		Summary:        summary,
		TookMs:         time.Since(start).Milliseconds(),
	}, nil
}

func fmtInterpretation(query string, results []SearchIndex) string {
	parsed := ParseQuery(query)
	if parsed.IsMetricQuery {
		return "Metric constraint filter"
	}
	if parsed.ResourceType != "" {
		return "Scoped resource lookup for " + parsed.ResourceType
	}
	return "Natural language platform discovery query"
}

// GetSuggestions delegates suggestion fetch
func (s *SearchService) GetSuggestions(query string, orgID uuid.UUID) ([]string, error) {
	return s.repo.GetSuggestions(query, orgID, 10)
}

// GetRecentSearches retrieves recent user queries
func (s *SearchService) GetRecentSearches(userID, orgID uuid.UUID) ([]RecentSearch, error) {
	return s.repo.GetRecentSearches(userID, orgID)
}

// GetPopularSearches retrieves top organization queries
func (s *SearchService) GetPopularSearches(orgID uuid.UUID) ([]PopularSearch, error) {
	return s.repo.GetPopularSearches(orgID)
}

// SaveSearch delegates saving a search query
func (s *SearchService) SaveSearch(userID, orgID uuid.UUID, name, query string, filters map[string]interface{}) error {
	return s.repo.SaveSearch(userID, orgID, name, query, filters)
}

// GetSavedSearches retrieves saved user searches
func (s *SearchService) GetSavedSearches(userID, orgID uuid.UUID) ([]SavedSearch, error) {
	return s.repo.GetSavedSearches(userID, orgID)
}

// DeleteSavedSearch deletes saved search by ID
func (s *SearchService) DeleteSavedSearch(id, userID, orgID uuid.UUID) error {
	return s.repo.DeleteSavedSearch(id, userID, orgID)
}
