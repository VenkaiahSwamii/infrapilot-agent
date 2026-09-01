package repository

import (
	"fmt"
	"strings"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SearchRepository handles search index operations
type SearchRepository struct {
	db *gorm.DB
}

// NewSearchRepository creates a new search repository
func NewSearchRepository() *SearchRepository {
	return &SearchRepository{
		db: database.DB,
	}
}

// SearchIndex represents search index entry for repository operations
type SearchIndex struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	ResourceType   string         `gorm:"size:50;not null;index" json:"resource_type"`
	ResourceID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"resource_id"`
	Title          string         `gorm:"type:text;not null" json:"title"`
	Description    string         `gorm:"type:text" json:"description"`
	Keywords       []string       `gorm:"type:text[]" json:"keywords"`
	OrganizationID uuid.UUID      `gorm:"type:uuid;index;not null" json:"organization_id"`
	Category       string         `gorm:"size:50;index" json:"category"`
	Severity       string         `gorm:"size:20" json:"severity"`
	Status         string         `gorm:"size:50" json:"status"`
	URL            string         `gorm:"type:text" json:"url"`
	Metadata       []byte         `gorm:"type:jsonb" json:"metadata"`
	RankScore      float64        `gorm:"default:0" json:"rank_score"`
	UpdatedAt      time.Time      `gorm:"index" json:"updated_at"`
	CreatedAt      time.Time      `json:"created_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// SearchResult represents a search result with highlighting
type SearchResult struct {
	SearchIndex
	HighlightTitle string `json:"highlight_title"`
	HighlightDesc  string `json:"highlight_description"`
}

// SearchOptions represents search query options
type SearchOptions struct {
	Query          string
	OrganizationID uuid.UUID
	Category       string
	Severity       string
	Status         string
	ResourceType   string
	DateFrom       *time.Time
	DateTo         *time.Time
	Page           int
	PageSize       int
	SortBy         string
	SortOrder      string
}

// SearchResultWithFacets represents search results with facet counts
type SearchResultWithFacets struct {
	Results  []SearchResult          `json:"results"`
	Total    int64                   `json:"total"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"page_size"`
	Facets   map[string][]FacetCount `json:"facets"`
	TookMs   int64                   `json:"took_ms"`
}

// FacetCount represents a facet aggregation
type FacetCount struct {
	Value string `json:"value"`
	Count int64  `json:"count"`
}

// Index indexes a resource for search
func (r *SearchRepository) Index(resourceType string, resourceID uuid.UUID, title string, description string, keywords []string, orgID uuid.UUID, category string, severity string, status string, url string, metadata []byte) error {
	now := time.Now()

	// Upsert: update if exists, insert if not
	var existing models.SearchIndex
	err := r.db.Where("resource_type = ? AND resource_id = ? AND organization_id = ? AND deleted_at IS NULL",
		resourceType, resourceID, orgID).First(&existing).Error

	if err == nil {
		// Update existing
		existing.Title = title
		existing.Description = description
		existing.Keywords = keywords
		existing.Category = category
		existing.Severity = severity
		existing.Status = status
		existing.URL = url
		existing.Metadata = metadata
		existing.UpdatedAt = now
		return r.db.Save(&existing).Error
	} else if err == gorm.ErrRecordNotFound {
		// Create new
		entry := models.SearchIndex{
			ResourceType:   models.SearchResultType(resourceType),
			ResourceID:     resourceID,
			Title:          title,
			Description:    description,
			Keywords:       keywords,
			OrganizationID: orgID,
			Category:       category,
			Severity:       severity,
			Status:         status,
			URL:            url,
			Metadata:       metadata,
			UpdatedAt:      now,
		}
		return r.db.Create(&entry).Error
	}
	return err
}

// RemoveFromIndex removes a resource from the search index
func (r *SearchRepository) RemoveFromIndex(resourceType string, resourceID uuid.UUID, orgID uuid.UUID) error {
	return r.db.Where("resource_type = ? AND resource_id = ? AND organization_id = ?",
		resourceType, resourceID, orgID).Delete(&models.SearchIndex{}).Error
}

// Search performs a full-text search
func (r *SearchRepository) Search(opts SearchOptions) (*SearchResultWithFacets, error) {
	start := time.Now()

	query := r.db.Table("search_index").Where("deleted_at IS NULL")

	// Organization filter (required for multi-tenancy)
	if opts.OrganizationID != uuid.Nil {
		query = query.Where("organization_id = ?", opts.OrganizationID)
	}

	// Full-text search
	if opts.Query != "" {
		searchQuery := strings.ToLower(opts.Query)
		searchTerms := strings.Fields(searchQuery)

		for _, term := range searchTerms {
			pattern := "%" + term + "%"
			query = query.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ? OR EXISTS (SELECT 1 FROM unnest(keywords) AS k WHERE LOWER(k) LIKE ?)",
				pattern, pattern, pattern)
		}
	}

	// Filters
	if opts.Category != "" {
		query = query.Where("category = ?", opts.Category)
	}
	if opts.Severity != "" {
		query = query.Where("severity = ?", opts.Severity)
	}
	if opts.Status != "" {
		query = query.Where("status = ?", opts.Status)
	}
	if opts.ResourceType != "" {
		query = query.Where("resource_type = ?", opts.ResourceType)
	}
	if opts.DateFrom != nil {
		query = query.Where("updated_at >= ?", *opts.DateFrom)
	}
	if opts.DateTo != nil {
		query = query.Where("updated_at <= ?", *opts.DateTo)
	}

	// Get total count
	var total int64
	countQuery := query.Session(&gorm.Session{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	// Facets
	facets := make(map[string][]FacetCount)

	// Category facets
	var categoryFacets []struct {
		Category string `json:"category"`
		Count    int64  `json:"count"`
	}
	if err := query.Select("category, COUNT(*) as count").Group("category").Scan(&categoryFacets).Error; err == nil {
		facets["category"] = make([]FacetCount, 0, len(categoryFacets))
		for _, f := range categoryFacets {
			if f.Category != "" {
				facets["category"] = append(facets["category"], FacetCount{Value: f.Category, Count: f.Count})
			}
		}
	}

	// Severity facets
	var severityFacets []struct {
		Severity string `json:"severity"`
		Count    int64  `json:"count"`
	}
	if err := query.Select("severity, COUNT(*) as count").Group("severity").Scan(&severityFacets).Error; err == nil {
		facets["severity"] = make([]FacetCount, 0, len(severityFacets))
		for _, f := range severityFacets {
			if f.Severity != "" {
				facets["severity"] = append(facets["severity"], FacetCount{Value: f.Severity, Count: f.Count})
			}
		}
	}

	// Sort
	sortBy := "updated_at"
	sortOrder := "DESC"
	if opts.SortBy != "" {
		sortBy = opts.SortBy
	}
	if opts.SortOrder != "" {
		sortOrder = strings.ToUpper(opts.SortOrder)
	}
	query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))

	// Pagination
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.PageSize < 1 {
		opts.PageSize = 20
	}
	offset := (opts.Page - 1) * opts.PageSize
	query = query.Limit(opts.PageSize).Offset(offset)

	// Execute search
	var results []SearchResult
	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	tookMs := time.Since(start).Milliseconds()

	return &SearchResultWithFacets{
		Results:  results,
		Total:    total,
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Facets:   facets,
		TookMs:   tookMs,
	}, nil
}

// GetSuggestions returns search suggestions for autocomplete
func (r *SearchRepository) GetSuggestions(query string, orgID uuid.UUID, limit int) ([]string, error) {
	if limit == 0 {
		limit = 10
	}

	searchQuery := strings.ToLower(query)
	searchTerms := strings.Fields(searchQuery)

	var suggestions []string
	suggestionMap := make(map[string]bool)

	// Get matching titles
	var titles []string
	q := r.db.Table("search_index").
		Select("DISTINCT title").
		Where("deleted_at IS NULL AND organization_id = ?", orgID).
		Limit(limit * 2)

	for _, term := range searchTerms {
		pattern := "%" + term + "%"
		q = q.Where("LOWER(title) LIKE ?", pattern)
	}

	if err := q.Find(&titles).Error; err != nil {
		return nil, err
	}

	for _, title := range titles {
		suggestionMap[strings.ToLower(title)] = true
	}

	// Get matching keywords
	var keywords []string
	q = r.db.Table("search_index").
		Select("DISTINCT unnest(keywords) as keyword").
		Where("deleted_at IS NULL AND organization_id = ?", orgID).
		Limit(limit * 2)

	for _, term := range searchTerms {
		pattern := "%" + term + "%"
		q = q.Where("LOWER(keyword) LIKE ?", pattern)
	}

	if err := q.Find(&keywords).Error; err != nil {
		return nil, err
	}

	for _, keyword := range keywords {
		suggestionMap[strings.ToLower(keyword)] = true
	}

	// Convert to slice
	for suggestion := range suggestionMap {
		suggestions = append(suggestions, suggestion)
	}

	// Sort by relevance (simple length match)
	sortOrder := func(i, j int) bool {
		return len(suggestions[i]) < len(suggestions[j])
	}

	// Simple bubble sort for small slices
	for i := 0; i < len(suggestions); i++ {
		for j := i + 1; j < len(suggestions); j++ {
			if sortOrder(j, i) {
				suggestions[i], suggestions[j] = suggestions[j], suggestions[i]
			}
		}
	}

	if len(suggestions) > limit {
		suggestions = suggestions[:limit]
	}

	return suggestions, nil
}

// GetRecentSearches returns recent searches for a user
func (r *SearchRepository) GetRecentSearches(userID, orgID uuid.UUID, limit int) ([]models.RecentSearch, error) {
	if limit == 0 {
		limit = 10
	}

	var searches []models.RecentSearch
	err := r.db.Where("user_id = ? AND organization_id = ?", userID, orgID).
		Order("searched_at DESC").
		Limit(limit).
		Find(&searches).Error

	return searches, err
}

// AddRecentSearch adds a search to recent searches
func (r *SearchRepository) AddRecentSearch(userID, orgID uuid.UUID, query string, resultCount int) error {
	search := models.RecentSearch{
		UserID:         userID,
		OrganizationID: orgID,
		Query:          query,
		ResultCount:    resultCount,
		SearchedAt:     time.Now(),
	}
	return r.db.Create(&search).Error
}

// GetPopularSearches returns popular search terms
func (r *SearchRepository) GetPopularSearches(orgID uuid.UUID, limit int) ([]models.PopularSearch, error) {
	if limit == 0 {
		limit = 10
	}

	var results []models.PopularSearch
	err := r.db.Table("recent_searches").
		Select("query, COUNT(*) as search_count").
		Where("organization_id = ?", orgID).
		Group("query").
		Order("search_count DESC").
		Limit(limit).
		Scan(&results).Error

	return results, err
}

// SaveSearch saves a search query
func (r *SearchRepository) SaveSearch(userID, orgID uuid.UUID, name string, query string, filters []byte) error {
	search := models.SavedSearch{
		UserID:         userID,
		OrganizationID: orgID,
		Name:           name,
		Query:          query,
		Filters:        filters,
		UpdatedAt:      time.Now(),
	}
	return r.db.Create(&search).Error
}

// GetSavedSearches returns saved searches for a user
func (r *SearchRepository) GetSavedSearches(userID, orgID uuid.UUID) ([]models.SavedSearch, error) {
	var searches []models.SavedSearch
	err := r.db.Where("user_id = ? AND organization_id = ?", userID, orgID).
		Order("updated_at DESC").
		Find(&searches).Error
	return searches, err
}

// DeleteSavedSearch deletes a saved search
func (r *SearchRepository) DeleteSavedSearch(id, userID, orgID uuid.UUID) error {
	return r.db.Where("id = ? AND user_id = ? AND organization_id = ?", id, userID, orgID).
		Delete(&models.SavedSearch{}).Error
}

// TrackSearchAnalytics tracks search analytics
func (r *SearchRepository) TrackSearchAnalytics(userID, orgID uuid.UUID, query string, resultCount int, filters []byte, responseTimeMs int) error {
	analytics := models.SearchAnalytics{
		UserID:         userID,
		OrganizationID: orgID,
		Query:          query,
		ResultCount:    resultCount,
		Filters:        filters,
		ResponseTimeMs: responseTimeMs,
		SearchedAt:     time.Now(),
	}
	return r.db.Create(&analytics).Error
}

// GetByID gets a search index entry by ID
func (r *SearchRepository) GetByID(id uuid.UUID) (*models.SearchIndex, error) {
	var entry models.SearchIndex
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&entry).Error
	if err != nil {
		return nil, err
	}
	return &entry, nil
}
