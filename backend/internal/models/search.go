package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SearchResultType represents the category of search result
type SearchResultType string

const (
	SearchResultTypeMachine    SearchResultType = "machine"
	SearchResultTypeLog        SearchResultType = "log"
	SearchResultTypeAlert      SearchResultType = "alert"
	SearchResultTypeIncident   SearchResultType = "incident"
	SearchResultTypeTrace      SearchResultType = "trace"
	SearchResultTypeDocker     SearchResultType = "docker"
	SearchResultTypeKubernetes SearchResultType = "kubernetes"
	SearchResultTypeAI         SearchResultType = "ai"
	SearchResultTypeReport     SearchResultType = "report"
	SearchResultTypeUser       SearchResultType = "user"
	SearchResultTypeDashboard  SearchResultType = "dashboard"
)

// SearchIndex represents a unified search index entry
type SearchIndex struct {
	ID             uuid.UUID        `gorm:"type:uuid;primaryKey" json:"id"`
	ResourceType   SearchResultType `gorm:"size:50;not null;index" json:"resource_type"`
	ResourceID     uuid.UUID        `gorm:"type:uuid;not null;index" json:"resource_id"`
	Title          string           `gorm:"type:text;not null" json:"title"`
	Description    string           `gorm:"type:text" json:"description"`
	Keywords       []string         `gorm:"type:text[]" json:"keywords"`
	OrganizationID uuid.UUID        `gorm:"type:uuid;index;not null" json:"organization_id"`
	Category       string           `gorm:"size:50;index" json:"category"`
	Severity       string           `gorm:"size:20" json:"severity"`
	Status         string           `gorm:"size:50" json:"status"`
	URL            string           `gorm:"type:text" json:"url"`
	Metadata       []byte           `gorm:"type:jsonb" json:"metadata"`
	RankScore      float64          `gorm:"default:0" json:"rank_score"`
	UpdatedAt      time.Time        `gorm:"index" json:"updated_at"`
	CreatedAt      time.Time        `json:"created_at"`
	DeletedAt      gorm.DeletedAt   `gorm:"index" json:"-"`
}

// SearchIndexBeforeCreate hooks the ID generator
func (s *SearchIndex) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// SavedSearch represents a user's saved search
type SavedSearch struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index" json:"organization_id"`
	Name           string    `gorm:"size:255;not null" json:"name"`
	Query          string    `gorm:"type:text;not null" json:"query"`
	Filters        []byte    `gorm:"type:jsonb" json:"filters"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (s *SavedSearch) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// RecentSearch represents a recent search query
type RecentSearch struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index" json:"organization_id"`
	Query          string    `gorm:"type:text;not null" json:"query"`
	ResultCount    int       `json:"result_count"`
	SearchedAt     time.Time `gorm:"index" json:"searched_at"`
}

func (r *RecentSearch) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// SearchAnalytics represents search analytics tracking
type SearchAnalytics struct {
	ID                uuid.UUID        `gorm:"type:uuid;primaryKey" json:"id"`
	UserID            uuid.UUID        `gorm:"type:uuid;not null;index" json:"user_id"`
	OrganizationID    uuid.UUID        `gorm:"type:uuid;not null;index" json:"organization_id"`
	Query             string           `gorm:"type:text;not null" json:"query"`
	ResultCount       int              `json:"result_count"`
	Filters           []byte           `gorm:"type:jsonb" json:"filters"`
	ClickedResultID   uuid.UUID        `gorm:"type:uuid" json:"clicked_result_id"`
	ClickedResultType SearchResultType `gorm:"size:50" json:"clicked_result_type"`
	ResponseTimeMs    int              `json:"response_time_ms"`
	SearchedAt        time.Time        `gorm:"index" json:"searched_at"`
}

func (s *SearchAnalytics) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// SearchRequest represents a search API request
type SearchRequest struct {
	Query          string     `json:"query" binding:"required,min=1"`
	Category       string     `json:"category"`
	Severity       string     `json:"severity"`
	Status         string     `json:"status"`
	OrganizationID uuid.UUID  `json:"organization_id"`
	ResourceType   string     `json:"resource_type"`
	DateFrom       *time.Time `json:"date_from"`
	DateTo         *time.Time `json:"date_to"`
	Tags           []string   `json:"tags"`
	OwnerID        uuid.UUID  `json:"owner_id"`
	Page           int        `json:"page" binding:"min=1"`
	PageSize       int        `json:"page_size" binding:"min=1,max=100"`
}

// SearchResponse represents a search API response
type SearchResponse struct {
	Query       string                  `json:"query"`
	TotalCount  int64                   `json:"total_count"`
	Page        int                     `json:"page"`
	PageSize    int                     `json:"page_size"`
	TotalPages  int                     `json:"total_pages"`
	Results     []SearchIndex           `json:"results"`
	Suggestions []string                `json:"suggestions"`
	Facets      map[string][]FacetCount `json:"facets"`
	TookMs      int64                   `json:"took_ms"`
}

// FacetCount represents a facet aggregation result
type FacetCount struct {
	Value string `json:"value"`
	Count int64  `json:"count"`
}

// SearchSuggestion represents an autocomplete suggestion
type SearchSuggestion struct {
	Text        string           `json:"text"`
	Type        SearchResultType `json:"type"`
	Category    string           `json:"category"`
	Description string           `json:"description"`
}

// PopularSearch represents a popular search term
type PopularSearch struct {
	Query       string `json:"query"`
	SearchCount int    `json:"search_count"`
}

// AISearchRequest represents an AI-powered natural language search
type AISearchRequest struct {
	Query string `json:"query" binding:"required,min=1"`
}

// AISearchResponse represents AI search results with explanation
type AISearchResponse struct {
	Query          string        `json:"query"`
	Interpretation string        `json:"interpretation"`
	Results        []SearchIndex `json:"results"`
	Summary        string        `json:"summary"`
	TookMs         int64         `json:"took_ms"`
}
