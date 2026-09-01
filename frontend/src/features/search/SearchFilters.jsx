import React from 'react';
import './SearchFilters.css';

export const CATEGORIES = [
  { id: '', label: 'All Categories' },
  { id: 'machine', label: '🖥️ Infrastructure' },
  { id: 'log', label: '📄 Logs' },
  { id: 'alert', label: '⚠️ Alerts' },
  { id: 'incident', label: '🚨 Incidents' },
  { id: 'trace', label: '🔍 Traces & APM' },
  { id: 'docker', label: '🐳 Docker' },
  { id: 'kubernetes', label: '☸️ Kubernetes' },
  { id: 'ai', label: '🤖 AI Insights' },
  { id: 'report', label: '📊 Reports' },
  { id: 'user', label: '👤 Users' },
  { id: 'dashboard', label: '📈 Dashboards' },
];

export const SEVERITIES = [
  { id: '', label: 'All Severities' },
  { id: 'P1', label: 'P1 - Critical' },
  { id: 'P2', label: 'P2 - High' },
  { id: 'P3', label: 'P3 - Medium' },
  { id: 'P4', label: 'P4 - Low' },
];

export const STATUSES = [
  { id: '', label: 'All Statuses' },
  { id: 'ONLINE', label: 'Online / Healthy' },
  { id: 'OFFLINE', label: 'Offline / Down' },
  { id: 'OPEN', label: 'Open' },
  { id: 'ACTIVE', label: 'Active' },
  { id: 'RESOLVED', label: 'Resolved / Closed' },
];

export const TIME_RANGES = [
  { id: '', label: 'Any Time' },
  { id: '1h', label: 'Last 1 Hour' },
  { id: '24h', label: 'Last 24 Hours' },
  { id: '7d', label: 'Last 7 Days' },
  { id: '30d', label: 'Last 30 Days' },
];

export default function SearchFilters({ filters, onFilterChange, onReset }) {
  const handleChange = (field, value) => {
    onFilterChange({
      ...filters,
      [field]: value,
    });
  };

  return (
    <div className="search-filters-sidebar">
      <div className="filters-header">
        <h3>🔍 Search Filters</h3>
        <button className="reset-filters-btn" onClick={onReset}>Clear All</button>
      </div>

      {/* Category Filter */}
      <div className="filter-group">
        <label className="filter-label">Resource Category</label>
        <div className="category-pill-group">
          {CATEGORIES.map((cat) => (
            <button
              key={cat.id || 'all'}
              className={`category-pill ${filters.resource_type === cat.id ? 'active' : ''}`}
              onClick={() => handleChange('resource_type', cat.id)}
            >
              {cat.label}
            </button>
          ))}
        </div>
      </div>

      {/* Severity Filter */}
      <div className="filter-group">
        <label className="filter-label">Severity Level</label>
        <select
          className="filter-select"
          value={filters.severity || ''}
          onChange={(e) => handleChange('severity', e.target.value)}
        >
          {SEVERITIES.map((sev) => (
            <option key={sev.id || 'all'} value={sev.id}>{sev.label}</option>
          ))}
        </select>
      </div>

      {/* Status Filter */}
      <div className="filter-group">
        <label className="filter-label">Operational Status</label>
        <select
          className="filter-select"
          value={filters.status || ''}
          onChange={(e) => handleChange('status', e.target.value)}
        >
          {STATUSES.map((st) => (
            <option key={st.id || 'all'} value={st.id}>{st.label}</option>
          ))}
        </select>
      </div>

      {/* Date Range Filter */}
      <div className="filter-group">
        <label className="filter-label">Time Window</label>
        <select
          className="filter-select"
          value={filters.time_range || ''}
          onChange={(e) => handleChange('time_range', e.target.value)}
        >
          {TIME_RANGES.map((tr) => (
            <option key={tr.id || 'all'} value={tr.id}>{tr.label}</option>
          ))}
        </select>
      </div>

      {/* Scope / Machine Keyword Filter */}
      <div className="filter-group">
        <label className="filter-label">Specific Machine / Host</label>
        <input
          type="text"
          className="filter-input"
          placeholder="e.g. server01"
          value={filters.machine || ''}
          onChange={(e) => handleChange('machine', e.target.value)}
        />
      </div>
    </div>
  );
}
