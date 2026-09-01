import { useState, useEffect } from 'react';
import { searchApi } from '../../api/search';
import './SearchResults.css';

const RESOURCE_ICONS = {
  machine: '🖥️',
  log: '📄',
  alert: '⚠️',
  incident: '🚨',
  trace: '🔍',
  docker: '🐳',
  kubernetes: '☸️',
  ai: '🤖',
  report: '📊',
  user: '👤',
  dashboard: '📈',
};

const RESOURCE_LABELS = {
  machine: 'Machine',
  log: 'Log',
  alert: 'Alert',
  incident: 'Incident',
  trace: 'Trace',
  docker: 'Docker',
  kubernetes: 'Kubernetes',
  ai: 'AI',
  report: 'Report',
  user: 'User',
  dashboard: 'Dashboard',
};

export default function SearchResults({ query, filters = {}, onResultClick }) {
  const [results, setResults] = useState([]);
  const [loading, setLoading] = useState(false);
  const [totalCount, setTotalCount] = useState(0);
  const [page, setPage] = useState(1);
  const [facets, setFacets] = useState({});

  useEffect(() => {
    const performSearch = async () => {
      if (!query.trim()) {
        setResults([]);
        setTotalCount(0);
        return;
      }

      setLoading(true);
      try {
        const response = await searchApi.search({
          q: query,
          ...filters,
          page,
          page_size: 20,
        });
        setResults(response.results || []);
        setTotalCount(response.total_count || 0);
        setFacets(response.facets || {});
      } catch (error) {
        console.error('Search failed:', error);
        setResults([]);
      } finally {
        setLoading(false);
      }
    };

    performSearch();
  }, [query, filters, page]);

  const handleResultClick = (result) => {
    if (onResultClick) {
      onResultClick(result);
    }
  };

  const getSeverityColor = (severity) => {
    switch (severity) {
      case 'P1': return '#ef4444';
      case 'P2': return '#f97316';
      case 'P3': return '#eab308';
      case 'P4': return '#3b82f6';
      case 'CRITICAL': return '#ef4444';
      case 'HIGH': return '#f97316';
      case 'MEDIUM': return '#eab308';
      case 'LOW': return '#3b82f6';
      default: return '#6b7280';
    }
  };

  const getStatusColor = (status) => {
    switch (status) {
      case 'ONLINE':
      case 'OPEN':
      case 'ACTIVE':
        return '#10b981';
      case 'OFFLINE':
      case 'CLOSED':
      case 'RESOLVED':
        return '#6b7280';
      case 'WARNING':
        return '#f59e0b';
      case 'CRITICAL':
        return '#ef4444';
      default:
        return '#6b7280';
    }
  };

  if (loading) {
    return <div className="search-loading-state">Searching...</div>;
  }

  if (results.length === 0 && query) {
    return <div className="search-empty-state">No results found for "{query}"</div>;
  }

  return (
    <div className="search-results">
      {totalCount > 0 && (
        <div className="search-results-header">
          <span className="search-results-count">
            Found {totalCount} result{totalCount !== 1 ? 's' : ''} for "{query}"
          </span>
        </div>
      )}

      <div className="search-results-list">
        {results.map((result) => (
          <div
            key={result.id}
            className="search-result-item"
            onClick={() => handleResultClick(result)}
          >
            <div className="search-result-icon">
              {RESOURCE_ICONS[result.resource_type] || '📄'}
            </div>
            <div className="search-result-content">
              <div className="search-result-title">
                {result.title}
                <span className="search-result-type">
                  {RESOURCE_LABELS[result.resource_type] || result.resource_type}
                </span>
              </div>
              {result.description && (
                <div className="search-result-description">
                  {result.description}
                </div>
              )}
              <div className="search-result-meta">
                {result.category && (
                  <span className="search-result-category">{result.category}</span>
                )}
                {result.severity && (
                  <span
                    className="search-result-severity"
                    style={{ color: getSeverityColor(result.severity) }}
                  >
                    {result.severity}
                  </span>
                )}
                {result.status && (
                  <span
                    className="search-result-status"
                    style={{ color: getStatusColor(result.status) }}
                  >
                    {result.status}
                  </span>
                )}
              </div>
            </div>
          </div>
        ))}
      </div>

      {totalCount > 20 && (
        <div className="search-pagination">
          <button
            onClick={() => setPage(p => Math.max(1, p - 1))}
            disabled={page === 1}
            className="search-pagination-button"
          >
            Previous
          </button>
          <span className="search-pagination-info">
            Page {page}
          </span>
          <button
            onClick={() => setPage(p => p + 1)}
            disabled={results.length < 20}
            className="search-pagination-button"
          >
            Next
          </button>
        </div>
      )}
    </div>
  );
}