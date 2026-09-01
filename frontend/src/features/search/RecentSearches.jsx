import React from 'react';
import './RecentSearches.css';

export default function RecentSearches({ recentSearches = [], popularSearches = [], savedSearches = [], onSelectSearch, onDeleteSaved }) {
  return (
    <div className="recent-searches-container">
      {/* Saved Searches */}
      {savedSearches.length > 0 && (
        <div className="search-section">
          <div className="section-title">
            <span>⭐ Saved Searches</span>
          </div>
          <div className="search-chips">
            {savedSearches.map((item) => (
              <div key={item.id} className="search-chip saved" onClick={() => onSelectSearch(item.query)}>
                <span className="chip-name">{item.name}</span>
                <span className="chip-query">({item.query})</span>
                {onDeleteSaved && (
                  <button
                    className="delete-chip-btn"
                    onClick={(e) => {
                      e.stopPropagation();
                      onDeleteSaved(item.id);
                    }}
                    title="Remove saved search"
                  >
                    ×
                  </button>
                )}
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Popular Platform Searches */}
      <div className="search-section">
        <div className="section-title">
          <span>🔥 Popular Platform Searches</span>
        </div>
        <div className="search-chips">
          {popularSearches.map((item, idx) => (
            <div key={idx} className="search-chip popular" onClick={() => onSelectSearch(item.query)}>
              <span className="chip-query">🔍 {item.query}</span>
              {item.search_count > 0 && <span className="chip-badge">{item.search_count} searches</span>}
            </div>
          ))}
        </div>
      </div>

      {/* Recent User Searches */}
      {recentSearches.length > 0 && (
        <div className="search-section">
          <div className="section-title">
            <span>🕒 Recent Searches</span>
          </div>
          <div className="search-chips">
            {recentSearches.map((item, idx) => (
              <div key={idx} className="search-chip recent" onClick={() => onSelectSearch(item.query)}>
                <span className="chip-query">🕒 {item.query}</span>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
