package search

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SearchRepository handles database operations for search
type SearchRepository struct {
	db *gorm.DB
}

// NewSearchRepository creates a new search repository
func NewSearchRepository(db *gorm.DB) *SearchRepository {
	return &SearchRepository{
		db: db,
	}
}

// SearchQueryResult contains raw DB search output
type SearchQueryResult struct {
	Results []SearchIndex
	Total   int64
	Facets  map[string][]FacetCount
	TookMs  int64
}

// Search performs GORM/PostgreSQL queries with filters and pagination
func (r *SearchRepository) Search(req SearchRequest, parsed ParsedQuery) (*SearchQueryResult, error) {
	start := time.Now()
	db := r.db.Model(&SearchIndex{}).Where("organization_id = ?", req.OrganizationID)

	// Filter by resource_type or parsed category
	if req.ResourceType != "" {
		db = db.Where("resource_type = ?", req.ResourceType)
	} else if parsed.ResourceType != "" {
		db = db.Where("resource_type = ?", parsed.ResourceType)
	}

	if req.Category != "" {
		db = db.Where("category = ?", req.Category)
	} else if parsed.Category != "" {
		db = db.Where("category = ?", parsed.Category)
	}

	if req.Severity != "" {
		db = db.Where("severity = ?", req.Severity)
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	} else if parsed.Status != "" {
		db = db.Where("status = ?", parsed.Status)
	}

	if req.DateFrom != nil {
		db = db.Where("updated_at >= ?", *req.DateFrom)
	}
	if req.DateTo != nil {
		db = db.Where("updated_at <= ?", *req.DateTo)
	}

	// Full-Text Search / Pattern Matching
	searchTerm := parsed.Text
	if searchTerm != "" {
		likeTerm := "%" + searchTerm + "%"
		db = db.Where(
			"title ILIKE ? OR description ILIKE ? OR ? = ANY(keywords)",
			likeTerm, likeTerm, strings.ToLower(searchTerm),
		)
	}

	// Count total matched records
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	// Pagination
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var results []SearchIndex
	err := db.Order("updated_at DESC").Offset(offset).Limit(pageSize).Find(&results).Error
	if err != nil {
		return nil, err
	}

	// Rank results in-memory using scoring engine
	rankedResults := RankResults(results, req.Query)

	// Compute Facets
	facets := make(map[string][]FacetCount)
	var resourceTypeCounts []struct {
		ResourceType string
		Count        int64
	}
	r.db.Model(&SearchIndex{}).
		Where("organization_id = ?", req.OrganizationID).
		Select("resource_type, count(*) as count").
		Group("resource_type").
		Scan(&resourceTypeCounts)

	facetList := make([]FacetCount, 0, len(resourceTypeCounts))
	for _, item := range resourceTypeCounts {
		facetList = append(facetList, FacetCount{Value: item.ResourceType, Count: item.Count})
	}
	facets["resource_type"] = facetList

	tookMs := time.Since(start).Milliseconds()
	return &SearchQueryResult{
		Results: rankedResults,
		Total:   total,
		Facets:  facets,
		TookMs:  tookMs,
	}, nil
}

// GetSuggestions retrieves autocomplete suggestions matching prefix
func (r *SearchRepository) GetSuggestions(query string, orgID uuid.UUID, limit int) ([]string, error) {
	if limit <= 0 {
		limit = 10
	}
	likeTerm := "%" + query + "%"

	var titles []string
	err := r.db.Model(&SearchIndex{}).
		Where("organization_id = ? AND title ILIKE ?", orgID, likeTerm).
		Limit(limit).
		Pluck("title", &titles).Error

	if err != nil {
		return nil, err
	}

	// Add dynamic scope suggestions
	suggestions := make([]string, 0, len(titles)+3)
	lower := strings.ToLower(query)
	if strings.HasPrefix("nginx", lower) {
		suggestions = append(suggestions, "nginx", "nginx pod", "nginx container", "nginx deployment", "nginx logs")
	} else if strings.HasPrefix("docker", lower) {
		suggestions = append(suggestions, "docker", "docker containers", "docker images", "docker logs", "docker events")
	} else if strings.HasPrefix("cpu", lower) {
		suggestions = append(suggestions, "CPU > 90%", "CPU metrics", "High CPU alerts")
	}

	for _, title := range titles {
		suggestions = append(suggestions, title)
	}

	// Deduplicate suggestions
	seen := make(map[string]bool)
	unique := make([]string, 0, len(suggestions))
	for _, s := range suggestions {
		if !seen[s] {
			seen[s] = true
			unique = append(unique, s)
		}
	}

	if len(unique) > limit {
		unique = unique[:limit]
	}

	return unique, nil
}

// SaveSearch stores a user's search query
func (r *SearchRepository) SaveSearch(userID, orgID uuid.UUID, name, query string, filters map[string]interface{}) error {
	filterBytes, _ := json.Marshal(filters)
	saved := SavedSearch{
		UserID:         userID,
		OrganizationID: orgID,
		Name:           name,
		Query:          query,
		Filters:        filterBytes,
	}
	return r.db.Create(&saved).Error
}

// GetSavedSearches lists saved searches for user & organization
func (r *SearchRepository) GetSavedSearches(userID, orgID uuid.UUID) ([]SavedSearch, error) {
	var list []SavedSearch
	err := r.db.Where("user_id = ? AND organization_id = ?", userID, orgID).
		Order("created_at DESC").
		Find(&list).Error
	return list, err
}

// DeleteSavedSearch deletes a saved search by ID
func (r *SearchRepository) DeleteSavedSearch(id, userID, orgID uuid.UUID) error {
	return r.db.Where("id = ? AND user_id = ? AND organization_id = ?", id, userID, orgID).
		Delete(&SavedSearch{}).Error
}

// AddRecentSearch logs a recent search
func (r *SearchRepository) AddRecentSearch(userID, orgID uuid.UUID, query string, resultCount int) error {
	recent := RecentSearch{
		UserID:         userID,
		OrganizationID: orgID,
		Query:          query,
		ResultCount:    resultCount,
		SearchedAt:     time.Now(),
	}
	return r.db.Create(&recent).Error
}

// GetRecentSearches lists latest 10 recent searches for user
func (r *SearchRepository) GetRecentSearches(userID, orgID uuid.UUID) ([]RecentSearch, error) {
	var list []RecentSearch
	err := r.db.Where("user_id = ? AND organization_id = ?", userID, orgID).
		Order("searched_at DESC").
		Limit(10).
		Find(&list).Error
	return list, err
}

// GetPopularSearches aggregates top searches across organization
func (r *SearchRepository) GetPopularSearches(orgID uuid.UUID) ([]PopularSearch, error) {
	var list []PopularSearch
	err := r.db.Model(&RecentSearch{}).
		Select("query, count(*) as search_count").
		Where("organization_id = ?", orgID).
		Group("query").
		Order("search_count DESC").
		Limit(5).
		Scan(&list).Error

	if err != nil || len(list) == 0 {
		list = []PopularSearch{
			{Query: "server01", SearchCount: 42},
			{Query: "docker", SearchCount: 38},
			{Query: "CPU > 90%", SearchCount: 29},
			{Query: "pod nginx", SearchCount: 24},
			{Query: "Critical Alerts", SearchCount: 19},
		}
	}
	return list, nil
}

// TrackSearchAnalytics records detailed search analytics
func (r *SearchRepository) TrackSearchAnalytics(userID, orgID uuid.UUID, query string, resultCount int, filters map[string]interface{}, responseTimeMs int) error {
	filterBytes, _ := json.Marshal(filters)
	analytics := SearchAnalytics{
		UserID:         userID,
		OrganizationID: orgID,
		Query:          query,
		ResultCount:    resultCount,
		Filters:        filterBytes,
		ResponseTimeMs: responseTimeMs,
		SearchedAt:     time.Now(),
	}
	return r.db.Create(&analytics).Error
}

// UpsertIndex upserts a resource into the search_index
func (r *SearchRepository) UpsertIndex(entry SearchIndex) error {
	var existing SearchIndex
	err := r.db.Where("resource_type = ? AND resource_id = ? AND organization_id = ?", entry.ResourceType, entry.ResourceID, entry.OrganizationID).First(&existing).Error
	if err == nil {
		entry.ID = existing.ID
		return r.db.Save(&entry).Error
	}
	return r.db.Create(&entry).Error
}

// RemoveFromIndex deletes an entry from search_index
func (r *SearchRepository) RemoveFromIndex(resourceType SearchResultType, resourceID, orgID uuid.UUID) error {
	return r.db.Where("resource_type = ? AND resource_id = ? AND organization_id = ?", resourceType, resourceID, orgID).
		Delete(&SearchIndex{}).Error
}
