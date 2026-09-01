import React from 'react';
import './SearchSuggestions.css';

export default function SearchSuggestions({ suggestions, onSelect, activeIndex }) {
  if (!suggestions || suggestions.length === 0) return null;

  return (
    <div className="search-suggestions-dropdown">
      <div className="suggestions-header">Quick Suggestions</div>
      {suggestions.map((item, index) => {
        const text = typeof item === 'string' ? item : item.text;
        const icon = typeof item === 'object' && item.type ? getIconForType(item.type) : '⚡';

        return (
          <div
            key={index}
            className={`suggestion-item ${index === activeIndex ? 'active' : ''}`}
            onClick={() => onSelect(text)}
          >
            <span className="suggestion-icon">{icon}</span>
            <span className="suggestion-text">{text}</span>
          </div>
        );
      })}
    </div>
  );
}

function getIconForType(type) {
  switch (type) {
    case 'machine': return '🖥️';
    case 'log': return '📄';
    case 'alert': return '⚠️';
    case 'incident': return '🚨';
    case 'docker': return '🐳';
    case 'kubernetes': return '☸️';
    case 'ai': return '🤖';
    default: return '🔍';
  }
}
