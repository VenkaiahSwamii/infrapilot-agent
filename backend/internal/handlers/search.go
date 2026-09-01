package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"infrapilot/backend/internal/models"
	"infrapilot/backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SearchHandler handles search-related HTTP requests
type SearchHandler struct {
	searchService *services.SearchService
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(searchService *services.SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

// Search searches across all platform resources
func (h *SearchHandler) Search(c *gin.Context) {
	start := time.Now()

	// Get authenticated user
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	orgID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Parse query parameters
	query := c.Query("q")
	if strings.TrimSpace(query) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	// Parse filters
	category := c.Query("category")
	severity := c.Query("severity")
	status := c.Query("status")
	resourceType := c.Query("resource_type")

	// Parse pagination
	page := 1
	pageSize := 20
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	// Parse date filters
	var dateFrom, dateTo *time.Time
	if df := c.Query("date_from"); df != "" {
		if parsed, err := time.Parse(time.RFC3339, df); err == nil {
			dateFrom = &parsed
		}
	}
	if dt := c.Query("date_to"); dt != "" {
		if parsed, err := time.Parse(time.RFC3339, dt); err == nil {
			dateTo = &parsed
		}
	}

	// Build search request
	req := services.SearchRequest{
		Query:          query,
		Category:       category,
		Severity:       severity,
		Status:         status,
		OrganizationID: orgID.(uuid.UUID),
		ResourceType:   resourceType,
		DateFrom:       dateFrom,
		DateTo:         dateTo,
		Page:           page,
		PageSize:       pageSize,
		UserID:         userID.(uuid.UUID),
	}

	// Execute search
	result, err := h.searchService.Search(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}

	responseTime := time.Since(start).Milliseconds()

	c.JSON(http.StatusOK, models.SearchResponse{
		Query:       result.Query,
		TotalCount:  result.TotalCount,
		Page:        result.Page,
		PageSize:    result.PageSize,
		TotalPages:  result.TotalPages,
		Results:     result.Results,
		Suggestions: result.Suggestions,
		Facets:      result.Facets,
		TookMs:      responseTime,
	})
}

// GetSuggestions returns search suggestions for autocomplete
func (h *SearchHandler) GetSuggestions(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	orgID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	query := c.Query("q")
	if strings.TrimSpace(query) == "" {
		c.JSON(http.StatusOK, []string{})
		return
	}

	suggestions, err := h.searchService.GetSuggestions(query, orgID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get suggestions"})
		return
	}

	c.JSON(http.StatusOK, suggestions)
}

// GetRecentSearches returns recent searches for the current user
func (h *SearchHandler) GetRecentSearches(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	orgID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	searches, err := h.searchService.GetRecentSearches(userID.(uuid.UUID), orgID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get recent searches"})
		return
	}

	c.JSON(http.StatusOK, searches)
}

// GetPopularSearches returns popular search terms
func (h *SearchHandler) GetPopularSearches(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	orgID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	searches, err := h.searchService.GetPopularSearches(orgID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get popular searches"})
		return
	}

	c.JSON(http.StatusOK, searches)
}

// SaveSearch saves a search query
func (h *SearchHandler) SaveSearch(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	orgID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		Name    string                 `json:"name" binding:"required"`
		Query   string                 `json:"query" binding:"required"`
		Filters map[string]interface{} `json:"filters"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.searchService.SaveSearch(userID.(uuid.UUID), orgID.(uuid.UUID), req.Name, req.Query, req.Filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save search"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "search saved"})
}

// GetSavedSearches returns saved searches for the current user
func (h *SearchHandler) GetSavedSearches(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	orgID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	searches, err := h.searchService.GetSavedSearches(userID.(uuid.UUID), orgID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get saved searches"})
		return
	}

	c.JSON(http.StatusOK, searches)
}

// DeleteSavedSearch deletes a saved search
func (h *SearchHandler) DeleteSavedSearch(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	orgID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	searchIDStr := c.Param("id")
	searchID, err := uuid.Parse(searchIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid search id"})
		return
	}

	err = h.searchService.DeleteSavedSearch(searchID, userID.(uuid.UUID), orgID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete saved search"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "search deleted"})
}

// AISearch performs AI-powered natural language search
func (h *SearchHandler) AISearch(c *gin.Context) {
	start := time.Now()

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	orgID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req models.AISearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.searchService.AISearch(userID.(uuid.UUID), orgID.(uuid.UUID), req.Query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI search failed"})
		return
	}

	c.JSON(http.StatusOK, models.AISearchResponse{
		Query:          result.Query,
		Interpretation: result.Query,
		Results:        result.Results,
		Summary:        result.Summary,
		TookMs:         time.Since(start).Milliseconds(),
	})
}
