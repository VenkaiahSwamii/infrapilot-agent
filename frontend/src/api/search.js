import { apiClient } from './client';

export const searchApi = {
  // Search across all resources
  search: async (params = {}) => {
    const query = new URLSearchParams();
    if (params.q) query.append('q', params.q);
    if (params.category) query.append('category', params.category);
    if (params.severity) query.append('severity', params.severity);
    if (params.status) query.append('status', params.status);
    if (params.resource_type) query.append('resource_type', params.resource_type);
    if (params.page) query.append('page', params.page);
    if (params.page_size) query.append('page_size', params.page_size);
    if (params.date_from) query.append('date_from', params.date_from);
    if (params.date_to) query.append('date_to', params.date_to);

    const response = await apiClient.get(`/search?${query.toString()}`);
    return response.data;
  },

  // Get search suggestions
  getSuggestions: async (query) => {
    const response = await apiClient.get(`/search/suggestions?q=${encodeURIComponent(query)}`);
    return response.data;
  },

  // Get recent searches
  getRecentSearches: async () => {
    const response = await apiClient.get('/search/recent');
    return response.data;
  },

  // Get popular searches
  getPopularSearches: async () => {
    const response = await apiClient.get('/search/popular');
    return response.data;
  },

  // Save a search
  saveSearch: async (name, query, filters = {}) => {
    const response = await apiClient.post('/search/save', {
      name,
      query,
      filters,
    });
    return response.data;
  },

  // Get saved searches
  getSavedSearches: async () => {
    const response = await apiClient.get('/search/saved');
    return response.data;
  },

  // Delete saved search
  deleteSavedSearch: async (id) => {
    const response = await apiClient.delete(`/search/saved/${id}`);
    return response.data;
  },

  // AI-powered search
  aiSearch: async (query) => {
    const response = await apiClient.post('/search/ai', {
      query,
    });
    return response.data;
  },
};

export default searchApi;