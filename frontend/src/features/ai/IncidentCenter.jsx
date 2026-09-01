import React, { useState, useEffect } from 'react';
import { incidentApi } from '../../api/incidents.js';
import './IncidentCenter.css';

export default function IncidentCenter() {
  const [incidents, setIncidents] = useState([]);
  const [selectedIncident, setSelectedIncident] = useState(null);
  const [timeline, setTimeline] = useState([]);
  const [alerts, setAlerts] = useState([]);
  const [analytics, setAnalytics] = useState(null);
  const [loading, setLoading] = useState(true);
  const [view, setView] = useState('list'); // list, detail, analytics
  const [filters, setFilters] = useState({ status: '', severity: '' });
  const [newComment, setNewComment] = useState('');
  const [resolutionNote, setResolutionNote] = useState('');

  useEffect(() => {
    loadIncidents();
    loadAnalytics();
  }, [filters]);

  const loadIncidents = async () => {
    setLoading(true);
    try {
      const data = await incidentApi.getIncidents(filters);
      setIncidents(data);
    } catch (error) {
      console.error('Failed to load incidents:', error);
    } finally {
      setLoading(false);
    }
  };

  const loadAnalytics = async () => {
    try {
      const data = await incidentApi.getAnalytics();
      setAnalytics(data);
    } catch (error) {
      console.error('Failed to load analytics:', error);
    }
  };

  const selectIncident = async (incident) => {
    setSelectedIncident(incident);
    setView('detail');
    try {
      const [timelineData, alertsData] = await Promise.all([
        incidentApi.getTimeline(incident.id),
        incidentApi.getAlerts(incident.id),
      ]);
      setTimeline(timelineData);
      setAlerts(alertsData);
    } catch (error) {
      console.error('Failed to load incident details:', error);
    }
  };

  const handleResolve = async () => {
    if (!selectedIncident) return;
    try {
      await incidentApi.resolveIncident(selectedIncident.id, resolutionNote);
      setResolutionNote('');
      await selectIncident({ ...selectedIncident, status: 'RESOLVED' });
      loadIncidents();
    } catch (error) {
      console.error('Failed to resolve incident:', error);
    }
  };

  const handleClose = async () => {
    if (!selectedIncident) return;
    try {
      await incidentApi.closeIncident(selectedIncident.id);
      await selectIncident({ ...selectedIncident, status: 'CLOSED' });
      loadIncidents();
    } catch (error) {
      console.error('Failed to close incident:', error);
    }
  };

  const handleAddComment = async () => {
    if (!selectedIncident || !newComment.trim()) return;
    try {
      await incidentApi.addComment(selectedIncident.id, newComment);
      setNewComment('');
      // Refresh timeline
      const timelineData = await incidentApi.getTimeline(selectedIncident.id);
      setTimeline(timelineData);
    } catch (error) {
      console.error('Failed to add comment:', error);
    }
  };

  const getSeverityColor = (severity) => {
    const colors = {
      P1: '#dc2626', // Critical - Red
      P2: '#ea580c', // High - Orange
      P3: '#d97706', // Medium - Yellow
      P4: '#65a30d', // Low - Green
      P5: '#6b7280', // Info - Gray
    };
    return colors[severity] || '#6b7280';
  };

  const getStatusColor = (status) => {
    const colors = {
      OPEN: '#3b82f6',
      ACKNOWLEDGED: '#f59e0b',
      INVESTIGATING: '#8b5cf6',
      MITIGATED: '#10b981',
      RESOLVED: '#059669',
      CLOSED: '#6b7280',
    };
    return colors[status] || '#6b7280';
  };

  return (
    <div className="incident-center">
      <div className="incident-header">
        <h1>Incident Management Center</h1>
        <div className="incident-stats">
          {analytics && (
            <>
              <div className="stat-card">
                <div className="stat-value">{analytics.open_incidents}</div>
                <div className="stat-label">Open</div>
              </div>
              <div className="stat-card critical">
                <div className="stat-value">{analytics.critical_incidents}</div>
                <div className="stat-label">Critical</div>
              </div>
              <div className="stat-card resolved">
                <div className="stat-value">{analytics.resolved_incidents}</div>
                <div className="stat-label">Resolved</div>
              </div>
              <div className="stat-card">
                <div className="stat-value">{analytics.mttr_human}</div>
                <div className="stat-label">MTTR</div>
              </div>
            </>
          )}
        </div>
      </div>

      {view === 'list' && (
        <div className="incident-list-view">
          <div className="incident-filters">
            <select
              value={filters.status}
              onChange={(e) => setFilters({ ...filters, status: e.target.value })}
            >
              <option value="">All Statuses</option>
              <option value="OPEN">Open</option>
              <option value="ACKNOWLEDGED">Acknowledged</option>
              <option value="INVESTIGATING">Investigating</option>
              <option value="RESOLVED">Resolved</option>
              <option value="CLOSED">Closed</option>
            </select>
            <select
              value={filters.severity}
              onChange={(e) => setFilters({ ...filters, severity: e.target.value })}
            >
              <option value="">All Severities</option>
              <option value="P1">P1 - Critical</option>
              <option value="P2">P2 - High</option>
              <option value="P3">P3 - Medium</option>
              <option value="P4">P4 - Low</option>
            </select>
          </div>

          {loading ? (
            <div className="loading">Loading incidents...</div>
          ) : (
            <div className="incident-table">
              <table>
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>Title</th>
                    <th>Severity</th>
                    <th>Status</th>
                    <th>Alerts</th>
                    <th>Started</th>
                    <th>Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {incidents.length === 0 ? (
                    <tr>
                      <td colSpan="7" className="no-data">No incidents found</td>
                    </tr>
                  ) : (
                    incidents.map((incident) => (
                      <tr
                        key={incident.id}
                        onClick={() => selectIncident(incident)}
                        className={selectedIncident?.id === incident.id ? 'selected' : ''}
                      >
                        <td className="incident-id">{incident.id?.slice(0, 8)}...</td>
                        <td className="incident-title">{incident.title}</td>
                        <td>
                          <span
                            className="severity-badge"
                            style={{ backgroundColor: getSeverityColor(incident.severity) }}
                          >
                            {incident.severity}
                          </span>
                        </td>
                        <td>
                          <span
                            className="status-badge"
                            style={{ backgroundColor: getStatusColor(incident.status) }}
                          >
                            {incident.status}
                          </span>
                        </td>
                        <td>{incident.alert_count}</td>
                        <td>{new Date(incident.started_at).toLocaleString()}</td>
                        <td>
                          <button
                            className="btn-view"
                            onClick={(e) => {
                              e.stopPropagation();
                              selectIncident(incident);
                            }}
                          >
                            View
                          </button>
                        </td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>
          )}
        </div>
      )}

      {view === 'detail' && selectedIncident && (
        <div className="incident-detail-view">
          <button className="btn-back" onClick={() => setView('list')}>
            ← Back to List
          </button>

          <div className="incident-detail-header">
            <div className="incident-info">
              <h2>{selectedIncident.title}</h2>
              <div className="incident-meta">
                <span
                  className="severity-badge"
                  style={{ backgroundColor: getSeverityColor(selectedIncident.severity) }}
                >
                  {selectedIncident.severity}
                </span>
                <span
                  className="status-badge"
                  style={{ backgroundColor: getStatusColor(selectedIncident.status) }}
                >
                  {selectedIncident.status}
                </span>
                <span>Started: {new Date(selectedIncident.started_at).toLocaleString()}</span>
                {selectedIncident.resolved_at && (
                  <span>Resolved: {new Date(selectedIncident.resolved_at).toLocaleString()}</span>
                )}
              </div>
              {selectedIncident.root_cause && (
                <div className="root-cause">
                  <strong>Root Cause:</strong> {selectedIncident.root_cause}
                </div>
              )}
            </div>

            <div className="incident-actions">
              {selectedIncident.status !== 'RESOLVED' && selectedIncident.status !== 'CLOSED' && (
                <div className="action-group">
                  <input
                    type="text"
                    placeholder="Resolution note..."
                    value={resolutionNote}
                    onChange={(e) => setResolutionNote(e.target.value)}
                  />
                  <button className="btn-resolve" onClick={handleResolve}>
                    Resolve
                  </button>
                </div>
              )}
              {selectedIncident.status === 'RESOLVED' && (
                <button className="btn-close" onClick={handleClose}>
                  Close Incident
                </button>
              )}
            </div>
          </div>

          <div className="incident-detail-content">
            <div className="incident-sidebar">
              <div className="detail-section">
                <h3>Linked Alerts ({alerts.length})</h3>
                <div className="alerts-list">
                  {alerts.map((alert) => (
                    <div key={alert.id} className="alert-item">
                      <div className="alert-title">{alert.title}</div>
                      <div className="alert-meta">
                        <span className={`alert-severity ${alert.severity?.toLowerCase()}`}>
                          {alert.severity}
                        </span>
                        <span>{alert.category}</span>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </div>

            <div className="incident-main">
              <div className="detail-section">
                <h3>Timeline</h3>
                <div className="timeline">
                  {timeline.map((event) => (
                    <div key={event.id} className="timeline-event">
                      <div className="timeline-time">
                        {new Date(event.created_at).toLocaleString()}
                      </div>
                      <div className="timeline-content">
                        <div className="timeline-type">{event.event_type}</div>
                        <div className="timeline-message">{event.message}</div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              <div className="detail-section">
                <h3>Add Comment</h3>
                <div className="comment-form">
                  <textarea
                    placeholder="Add a comment..."
                    value={newComment}
                    onChange={(e) => setNewComment(e.target.value)}
                    rows={3}
                  />
                  <button className="btn-comment" onClick={handleAddComment}>
                    Add Comment
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {view === 'analytics' && analytics && (
        <div className="incident-analytics-view">
          <button className="btn-back" onClick={() => setView('list')}>
            ← Back to List
          </button>
          <h2>Incident Analytics</h2>
          <div className="analytics-grid">
            <div className="analytics-card">
              <h3>Root Cause Distribution</h3>
              <div className="distribution-list">
                {Object.entries(analytics.root_cause_distribution || {}).map(([cause, count]) => (
                  <div key={cause} className="distribution-item">
                    <span className="cause-name">{cause}</span>
                    <span className="cause-count">{count}</span>
                  </div>
                ))}
              </div>
            </div>
            <div className="analytics-card">
              <h3>Performance Metrics</h3>
              <div className="metrics-list">
                <div>MTTD: {analytics.mttd_human}</div>
                <div>MTTR: {analytics.mttr_human}</div>
              </div>
            </div>
            <div className="analytics-card">
              <h3>Top Affected Machines</h3>
              <div className="machine-list">
                {(analytics.top_affected_machines || []).map((machine) => (
                  <div key={machine.hostname} className="machine-item">
                    <span>{machine.hostname}</span>
                    <span>{machine.incidents} incidents</span>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}