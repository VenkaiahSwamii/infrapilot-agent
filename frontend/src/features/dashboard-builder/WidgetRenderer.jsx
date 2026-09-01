import React from 'react';

// Placeholder widget renderers - in production, these would render actual charts and data
const WidgetRenderers = {
  // Infrastructure widgets
  cpu: ({ config }) => (
    <div className="widget-chart">
      <div className="widget-metric">
        <span className="metric-value">45.5%</span>
        <span className="metric-label">CPU Usage</span>
      </div>
      <div className="widget-chart-placeholder">
        <div className="chart-bar" style={{ height: '60%' }}></div>
        <div className="chart-bar" style={{ height: '75%' }}></div>
        <div className="chart-bar" style={{ height: '45%' }}></div>
        <div className="chart-bar" style={{ height: '80%' }}></div>
        <div className="chart-bar" style={{ height: '55%' }}></div>
      </div>
    </div>
  ),

  memory: ({ config }) => (
    <div className="widget-chart">
      <div className="widget-metric">
        <span className="metric-value">62.3%</span>
        <span className="metric-label">Memory Usage</span>
      </div>
      <div className="widget-chart-placeholder">
        <div className="chart-bar" style={{ height: '70%' }}></div>
        <div className="chart-bar" style={{ height: '65%' }}></div>
        <div className="chart-bar" style={{ height: '75%' }}></div>
        <div className="chart-bar" style={{ height: '60%' }}></div>
        <div className="chart-bar" style={{ height: '68%' }}></div>
      </div>
    </div>
  ),

  disk: ({ config }) => (
    <div className="widget-chart">
      <div className="widget-metric">
        <span className="metric-value">78.2%</span>
        <span className="metric-label">Disk Usage</span>
      </div>
      <div className="progress-bar">
        <div className="progress-fill" style={{ width: '78.2%' }}></div>
      </div>
    </div>
  ),

  network: ({ config }) => (
    <div className="widget-chart">
      <div className="widget-metric">
        <span className="metric-value">1.2 Gbps</span>
        <span className="metric-label">Network Traffic</span>
      </div>
    </div>
  ),

  // Monitoring widgets
  alerts: ({ config }) => (
    <div className="widget-list">
      <div className="alert-item critical">
        <span className="alert-badge critical">Critical</span>
        <span className="alert-message">High CPU usage on server-01</span>
      </div>
      <div className="alert-item warning">
        <span className="alert-badge warning">Warning</span>
        <span className="alert-message">Memory usage above 80%</span>
      </div>
      <div className="alert-item info">
        <span className="alert-badge info">Info</span>
        <span className="alert-message">Backup completed successfully</span>
      </div>
    </div>
  ),

  active_incidents: ({ config }) => (
    <div className="widget-metric">
      <span className="metric-value">3</span>
      <span className="metric-label">Active Incidents</span>
    </div>
  ),

  health_score: ({ config }) => (
    <div className="widget-health-score">
      <div className="health-circle">
        <svg viewBox="0 0 100 100">
          <circle cx="50" cy="50" r="45" className="health-bg" />
          <circle cx="50" cy="50" r="45" className="health-fill" style={{ strokeDashoffset: 25 }} />
        </svg>
        <div className="health-text">
          <span className="health-grade">B+</span>
        </div>
      </div>
    </div>
  ),

  // Kubernetes widgets
  kubernetes_nodes: ({ config }) => (
    <div className="widget-metric">
      <span className="metric-value">5</span>
      <span className="metric-label">Nodes</span>
    </div>
  ),

  kubernetes_pods: ({ config }) => (
    <div className="widget-metric">
      <span className="metric-value">24</span>
      <span className="metric-label">Pods Running</span>
    </div>
  ),

  kubernetes_deployments: ({ config }) => (
    <div className="widget-metric">
      <span className="metric-value">8</span>
      <span className="metric-label">Deployments</span>
    </div>
  ),

  // Docker widgets
  docker_containers: ({ config }) => (
    <div className="widget-metric">
      <span className="metric-value">12</span>
      <span className="metric-label">Containers</span>
    </div>
  ),

  // Logs widgets
  live_logs: ({ config }) => (
    <div className="widget-logs">
      <div className="log-line">
        <span className="log-time">10:23:45</span>
        <span className="log-level info">INFO</span>
        <span className="log-message">Server started successfully</span>
      </div>
      <div className="log-line">
        <span className="log-time">10:23:46</span>
        <span className="log-level info">INFO</span>
        <span className="log-message">Connected to database</span>
      </div>
    </div>
  ),

  error_logs: ({ config }) => (
    <div className="widget-logs">
      <div className="log-line">
        <span className="log-time">10:23:45</span>
        <span className="log-level error">ERROR</span>
        <span className="log-message">Failed to connect to external API</span>
      </div>
    </div>
  ),

  // Default fallback
  default: ({ widget }) => (
    <div className="widget-placeholder">
      <span className="widget-type">{widget.widget_type}</span>
      <p>Widget content will be rendered here</p>
    </div>
  ),
};

export default function WidgetRenderer({ widget }) {
  const renderer = WidgetRenderers[widget.widget_type] || WidgetRenderers.default;
  
  // Parse config if available
  let config = {};
  try {
    if (widget.config) {
      config = typeof widget.config === 'string' ? JSON.parse(widget.config) : widget.config;
    }
  } catch (err) {
    console.error('Failed to parse widget config:', err);
  }

  return renderer({ widget, config });
}