package search

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SearchHandler handles HTTP search endpoints
type SearchHandler struct {
	service *SearchService
}

// NewSearchHandler creates a new SearchHandler instance
func NewSearchHandler(service *SearchService) *SearchHandler {
	return &SearchHandler{
		service: service,
	}
}

// Search handles search requests
func (h *SearchHandler) Search(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	orgIDVal, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	query := c.Query("q")
	if strings.TrimSpace(query) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	page := 1
	if p := c.Query("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}
	pageSize := 20
	if ps := c.Query("page_size"); ps != "" {
		if val, err := strconv.Atoi(ps); err == nil && val > 0 && val <= 100 {
			pageSize = val
		}
	}

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

	req := SearchRequest{
		Query:          query,
		Category:       c.Query("category"),
		Severity:       c.Query("severity"),
		Status:         c.Query("status"),
		ResourceType:   c.Query("resource_type"),
		Machine:        c.Query("machine"),
		DateFrom:       dateFrom,
		DateTo:         dateTo,
		Page:           page,
		PageSize:       pageSize,
		OrganizationID: orgIDVal.(uuid.UUID),
		UserID:         userIDVal.(uuid.UUID),
	}

	resp, err := h.service.Search(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetSuggestions handles autocomplete suggestions
func (h *SearchHandler) GetSuggestions(c *gin.Context) {
	orgIDVal, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	query := c.Query("q")
	if strings.TrimSpace(query) == "" {
		c.JSON(http.StatusOK, []string{})
		return
	}

	suggestions, err := h.service.GetSuggestions(query, orgIDVal.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get suggestions"})
		return
	}

	c.JSON(http.StatusOK, suggestions)
}

// GetRecentSearches handles recent search queries for user
func (h *SearchHandler) GetRecentSearches(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	orgIDVal, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	searches, err := h.service.GetRecentSearches(userIDVal.(uuid.UUID), orgIDVal.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get recent searches"})
		return
	}

	c.JSON(http.StatusOK, searches)
}

// GetPopularSearches handles popular searches for org
func (h *SearchHandler) GetPopularSearches(c *gin.Context) {
	orgIDVal, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	searches, err := h.service.GetPopularSearches(orgIDVal.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get popular searches"})
		return
	}

	c.JSON(http.StatusOK, searches)
}

// SaveSearch saves a user query
func (h *SearchHandler) SaveSearch(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	orgIDVal, exists := c.Get("organization_id")
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

	err := h.service.SaveSearch(userIDVal.(uuid.UUID), orgIDVal.(uuid.UUID), req.Name, req.Query, req.Filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save search"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "search saved"})
}

// GetSavedSearches lists saved searches
func (h *SearchHandler) GetSavedSearches(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	orgIDVal, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	searches, err := h.service.GetSavedSearches(userIDVal.(uuid.UUID), orgIDVal.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get saved searches"})
		return
	}

	c.JSON(http.StatusOK, searches)
}

// DeleteSavedSearch deletes saved search by ID
func (h *SearchHandler) DeleteSavedSearch(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	orgIDVal, exists := c.Get("organization_id")
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

	err = h.service.DeleteSavedSearch(searchID, userIDVal.(uuid.UUID), orgIDVal.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete saved search"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "search deleted"})
}

// AISearch executes natural language AI search
func (h *SearchHandler) AISearch(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	orgIDVal, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req AISearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.AISearch(userIDVal.(uuid.UUID), orgIDVal.(uuid.UUID), req.Query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI search failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
