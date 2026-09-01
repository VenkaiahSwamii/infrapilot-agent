import React, { useState, useEffect } from 'react';
import { searchApi } from '../../api/search';
import { getWebSocketUrl } from '../../websocket/liveEvents.js';
import SearchBar from './SearchBar';
import SearchResults from './SearchResults';
import SearchFilters from './SearchFilters';
import RecentSearches from './RecentSearches';
import './GlobalSearch.css';

export default function GlobalSearch() {
  const [query, setQuery] = useState('');
  const [isAIMode, setIsAIMode] = useState(false);
  const [aiSummary, setAiSummary] = useState('');
  const [aiInterpretation, setAiInterpretation] = useState('');

  const [filters, setFilters] = useState({
    resource_type: '',
    severity: '',
    status: '',
    time_range: '',
    machine: '',
  });

  const [results, setResults] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(1);

  const [recentSearches, setRecentSearches] = useState([]);
  const [popularSearches, setPopularSearches] = useState([]);
  const [savedSearches, setSavedSearches] = useState([]);

  const [saveModalOpen, setSaveModalOpen] = useState(false);
  const [saveName, setSaveName] = useState('');

  // Load initial search metadata (recent, popular, saved)
  useEffect(() => {
    loadMetadata();
  }, []);

  const loadMetadata = async () => {
    try {
      const [recent, popular, saved] = await Promise.all([
        searchApi.getRecentSearches().catch(() => []),
        searchApi.getPopularSearches().catch(() => []),
        searchApi.getSavedSearches().catch(() => []),
      ]);
      setRecentSearches(recent || []);
      setPopularSearches(popular || []);
      setSavedSearches(saved || []);
    } catch (err) {
      console.error('Failed to load search metadata', err);
    }
  };

  // Perform search on query or filter changes
  useEffect(() => {
    if (!query.trim()) {
      setResults([]);
      setTotalCount(0);
      setAiSummary('');
      return;
    }

    const executeSearch = async () => {
      setLoading(true);
      try {
        if (isAIMode) {
          const resp = await searchApi.aiSearch(query);
          setResults(resp.results || []);
          setTotalCount((resp.results || []).length);
          setAiSummary(resp.summary || '');
          setAiInterpretation(resp.interpretation || '');
        } else {
          const resp = await searchApi.search({
            q: query,
            ...filters,
            page,
            page_size: 20,
          });
          setResults(resp.results || []);
          setTotalCount(resp.total_count || 0);
          setAiSummary('');
        }
      } catch (err) {
        console.error('Search error:', err);
      } finally {
        setLoading(false);
      }
    };

    const timer = setTimeout(executeSearch, 250);
    return () => clearTimeout(timer);
  }, [query, filters, isAIMode, page]);

  // Setup WebSocket live update listener for search updates
  useEffect(() => {
    const wsUrl = getWebSocketUrl();
    let socket;
    try {
      socket = new WebSocket(wsUrl);
      socket.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          if (data.type === 'search_index_updated' && query.trim()) {
            // Trigger refresh
            setPage(p => p);
          }
        } catch (e) {
          // ignore
        }
      };
    } catch (e) {
      // WS connection fallback
    }

    return () => {
      if (socket && socket.readyState === WebSocket.OPEN) {
        socket.close();
      }
    };
  }, [query]);

  const handleSearchSubmit = (q) => {
    setQuery(q);
    setPage(1);
  };

  const handleResetFilters = () => {
    setFilters({
      resource_type: '',
      severity: '',
      status: '',
      time_range: '',
      machine: '',
    });
  };

  const handleSaveCurrentSearch = async () => {
    if (!saveName.trim() || !query.trim()) return;
    try {
      await searchApi.saveSearch(saveName, query, filters);
      setSaveModalOpen(false);
      setSaveName('');
      loadMetadata();
    } catch (err) {
      alert('Failed to save search');
    }
  };

  const handleDeleteSavedSearch = async (id) => {
    try {
      await searchApi.deleteSavedSearch(id);
      loadMetadata();
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <div className="global-search-container">
      {/* Header Banner */}
      <div className="global-search-header">
        <div className="header-title-row">
          <h1>🚀 Unified Enterprise Search</h1>
          <div className="search-mode-toggle">
            <button
              className={`mode-btn ${!isAIMode ? 'active' : ''}`}
              onClick={() => setIsAIMode(false)}
            >
              Standard Full-Text
            </button>
            <button
              className={`mode-btn ai ${isAIMode ? 'active' : ''}`}
              onClick={() => setIsAIMode(true)}
            >
              🤖 Natural Language AI
            </button>
          </div>
        </div>

        <p className="header-subtitle">
          Search across Machines, Logs, Alerts, Incidents, Kubernetes, Docker, Tracing, AI Insights, and Reports in one place.
        </p>

        {/* Global Search Bar */}
        <div className="search-bar-hero-wrapper">
          <SearchBar
            onSearch={handleSearchSubmit}
            placeholder={
              isAIMode
                ? "Ask natural language questions (e.g. 'Which servers are unhealthy?', 'Why is server01 slow?')..."
                : "Search machines, logs, incidents, Docker, K8s, CPU > 90%..."
            }
            autoFocus
          />
          {query.trim() && (
            <button
              className="save-search-trigger-btn"
              onClick={() => setSaveModalOpen(true)}
            >
              ⭐ Save Search
            </button>
          )}
        </div>

        {/* AI Explanation Banner */}
        {isAIMode && aiSummary && (
          <div className="ai-summary-card">
            <div className="ai-summary-header">
              <span className="ai-badge">🤖 AI Interpretation: {aiInterpretation}</span>
            </div>
            <p className="ai-summary-text">{aiSummary}</p>
          </div>
        )}
      </div>

      {/* Main Workspace Layout */}
      <div className="global-search-body">
        {/* Left Filters Sidebar */}
        <aside className="search-filters-column">
          <SearchFilters
            filters={filters}
            onFilterChange={setFilters}
            onReset={handleResetFilters}
          />
        </aside>

        {/* Right Search Content */}
        <main className="search-content-column">
          {query.trim() === '' ? (
            <RecentSearches
              recentSearches={recentSearches}
              popularSearches={popularSearches}
              savedSearches={savedSearches}
              onSelectSearch={handleSearchSubmit}
              onDeleteSaved={handleDeleteSavedSearch}
            />
          ) : (
            <SearchResults
              results={results}
              loading={loading}
              totalCount={totalCount}
              page={page}
              onPageChange={setPage}
              query={query}
            />
          )}
        </main>
      </div>

      {/* Save Search Modal */}
      {saveModalOpen && (
        <div className="modal-backdrop">
          <div className="save-modal-card">
            <h3>⭐ Save Search Query</h3>
            <p>Save query <code>"{query}"</code> for quick access in your dashboard.</p>
            <input
              type="text"
              className="save-modal-input"
              placeholder="Enter search name (e.g. Critical Alerts)"
              value={saveName}
              onChange={(e) => setSaveName(e.target.value)}
              autoFocus
            />
            <div className="modal-actions">
              <button className="cancel-btn" onClick={() => setSaveModalOpen(false)}>Cancel</button>
              <button className="confirm-btn" onClick={handleSaveCurrentSearch}>Save Query</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
