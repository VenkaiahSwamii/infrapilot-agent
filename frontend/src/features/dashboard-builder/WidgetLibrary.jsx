import React, { useState } from 'react';

const WIDGET_CATEGORIES = {
  infrastructure: {
    label: 'Infrastructure',
    icon: '🖥️',
    widgets: [
      { type: 'cpu', label: 'CPU Usage', icon: '⚡' },
      { type: 'memory', label: 'Memory Usage', icon: '💾' },
      { type: 'disk', label: 'Disk Usage', icon: '💿' },
      { type: 'network', label: 'Network', icon: '🌐' },
      { type: 'uptime', label: 'Uptime', icon: '⏱️' },
      { type: 'availability', label: 'Availability', icon: '✓' },
    ],
  },
  monitoring: {
    label: 'Monitoring',
    icon: '📊',
    widgets: [
      { type: 'alerts', label: 'Alerts', icon: '🔔' },
      { type: 'active_incidents', label: 'Active Incidents', icon: '🚨' },
      { type: 'ai_recommendations', label: 'AI Recommendations', icon: '🤖' },
      { type: 'health_score', label: 'Health Score', icon: '❤️' },
    ],
  },
  kubernetes: {
    label: 'Kubernetes',
    icon: '☸️',
    widgets: [
      { type: 'cluster_health', label: 'Cluster Health', icon: '🏥' },
      { type: 'kubernetes_nodes', label: 'Nodes', icon: '🖥️' },
      { type: 'kubernetes_pods', label: 'Pods', icon: '📦' },
      { type: 'kubernetes_deployments', label: 'Deployments', icon: '🚀' },
      { type: 'kubernetes_services', label: 'Services', icon: '🔗' },
    ],
  },
  docker: {
    label: 'Docker',
    icon: '🐳',
    widgets: [
      { type: 'docker_containers', label: 'Containers', icon: '📦' },
      { type: 'docker_images', label: 'Images', icon: '🖼️' },
      { type: 'docker_volumes', label: 'Volumes', icon: '💾' },
      { type: 'docker_networks', label: 'Networks', icon: '🌐' },
    ],
  },
  ai: {
    label: 'AI Insights',
    icon: '🧠',
    widgets: [
      { type: 'ai_insights', label: 'AI Insights', icon: '💡' },
      { type: 'root_cause', label: 'Root Cause Analysis', icon: '🔍' },
      { type: 'predictions', label: 'Predictions', icon: '🔮' },
    ],
  },
  logs: {
    label: 'Logs',
    icon: '📝',
    widgets: [
      { type: 'live_logs', label: 'Live Logs', icon: '📡' },
      { type: 'error_logs', label: 'Error Logs', icon: '❌' },
      { type: 'log_search', label: 'Search', icon: '🔎' },
    ],
  },
  apm: {
    label: 'APM',
    icon: '⚡',
    widgets: [
      { type: 'latency', label: 'Latency', icon: '⏱️' },
      { type: 'rps', label: 'RPS', icon: '📈' },
      { type: 'error_rate', label: 'Error Rate', icon: '📉' },
      { type: 'top_apis', label: 'Top APIs', icon: '🏆' },
    ],
  },
};

export default function WidgetLibrary({ onAddWidget, onClose }) {
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedCategory, setSelectedCategory] = useState(null);

  const filteredCategories = Object.entries(WIDGET_CATEGORIES).reduce((acc, [key, category]) => {
    const filteredWidgets = category.widgets.filter(widget =>
      widget.label.toLowerCase().includes(searchTerm.toLowerCase()) ||
      widget.type.toLowerCase().includes(searchTerm.toLowerCase())
    );
    
    if (filteredWidgets.length > 0) {
      acc[key] = { ...category, widgets: filteredWidgets };
    }
    return acc;
  }, {});

  const handleDragStart = (e, widgetType, label) => {
    e.dataTransfer.setData('application/json', JSON.stringify({ widgetType, label }));
    e.dataTransfer.effectAllowed = 'copy';
  };

  return (
    <div className="widget-library">
      <div className="widget-library-header">
        <h3>Widget Library</h3>
        <button className="close-btn" onClick={onClose}>×</button>
      </div>

      <div className="widget-library-search">
        <input
          type="text"
          placeholder="Search widgets..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="search-input"
        />
      </div>

      <div className="widget-library-content">
        {Object.entries(filteredCategories).map(([key, category]) => (
          <div key={key} className="widget-category">
            <div 
              className="category-header"
              onClick={() => setSelectedCategory(
                selectedCategory === key ? null : key
              )}
            >
              <span className="category-icon">{category.icon}</span>
              <span className="category-label">{category.label}</span>
              <span className="category-arrow">
                {selectedCategory === key ? '▼' : '▶'}
              </span>
            </div>

            {(selectedCategory === key || searchTerm) && (
              <div className="widget-grid">
                {category.widgets.map((widget) => (
                  <div
                    key={widget.type}
                    className="widget-item"
                    draggable
                    onDragStart={(e) => handleDragStart(e, widget.type, widget.label)}
                    onClick={() => onAddWidget(widget.type, widget.label)}
                  >
                    <span className="widget-icon">{widget.icon}</span>
                    <span className="widget-label">{widget.label}</span>
                  </div>
                ))}
              </div>
            )}
          </div>
        ))}
      </div>

      <div className="widget-library-footer">
        <p>Drag widgets to the canvas or click to add</p>
      </div>
    </div>
  );
}