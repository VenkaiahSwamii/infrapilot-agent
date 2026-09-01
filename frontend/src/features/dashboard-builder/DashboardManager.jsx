import React, { useState, useEffect } from 'react';
import dashboardsApi from '../../api/dashboards';

export default function DashboardManager({ onSelectDashboard, onCreateNew }) {
  const [dashboards, setDashboards] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [newDashboardName, setNewDashboardName] = useState('');
  const [templates, setTemplates] = useState([]);

  useEffect(() => {
    loadDashboards();
    loadTemplates();
  }, []);

  const loadDashboards = async () => {
    try {
      const response = await dashboardsApi.getDashboards();
      if (response.success) {
        setDashboards(response.data.dashboards || []);
      }
    } catch (err) {
      console.error('Failed to load dashboards:', err);
    } finally {
      setLoading(false);
    }
  };

  const loadTemplates = async () => {
    try {
      const response = await dashboardsApi.getTemplates();
      if (response.success) {
        setTemplates(response.data || []);
      }
    } catch (err) {
      console.error('Failed to load templates:', err);
    }
  };

  const handleCreateDashboard = async (e) => {
    e.preventDefault();
    if (!newDashboardName.trim()) return;

    try {
      const response = await dashboardsApi.createDashboard({
        name: newDashboardName,
        description: '',
      });
      
      if (response.success) {
        setShowCreateModal(false);
        setNewDashboardName('');
        loadDashboards();
        if (onSelectDashboard) {
          onSelectDashboard(response.data);
        }
      }
    } catch (err) {
      console.error('Failed to create dashboard:', err);
    }
  };

  const handleDeleteDashboard = async (id, e) => {
    e.stopPropagation();
    if (!window.confirm('Are you sure you want to delete this dashboard?')) {
      return;
    }

    try {
      await dashboardsApi.deleteDashboard(id);
      loadDashboards();
    } catch (err) {
      console.error('Failed to delete dashboard:', err);
    }
  };

  const handleDuplicateDashboard = async (id, e) => {
    e.stopPropagation();
    const name = prompt('Enter name for duplicated dashboard:');
    if (!name) return;

    try {
      const response = await dashboardsApi.duplicateDashboard(id, name);
      if (response.success) {
        loadDashboards();
      }
    } catch (err) {
      console.error('Failed to duplicate dashboard:', err);
    }
  };

  const handleSetDefault = async (id, e) => {
    e.stopPropagation();
    try {
      await dashboardsApi.setDefaultDashboard(id);
      loadDashboards();
    } catch (err) {
      console.error('Failed to set default dashboard:', err);
    }
  };

  return (
    <div className="dashboard-manager">
      <div className="dashboard-manager-header">
        <h2>Dashboards</h2>
        <button className="btn-primary" onClick={() => setShowCreateModal(true)}>
          + New Dashboard
        </button>
      </div>

      {loading ? (
        <div className="loading">Loading dashboards...</div>
      ) : (
        <div className="dashboard-list">
          {dashboards.length === 0 ? (
            <div className="empty-state">
              <p>No dashboards yet. Create your first dashboard to get started!</p>
            </div>
          ) : (
            dashboards.map(dashboard => (
              <div
                key={dashboard.id}
                className="dashboard-item"
                onClick={() => onSelectDashboard && onSelectDashboard(dashboard)}
              >
                <div className="dashboard-item-header">
                  <h3>{dashboard.name}</h3>
                  <div className="dashboard-item-actions">
                    {!dashboard.is_default && (
                      <button
                        className="btn-icon"
                        onClick={(e) => handleSetDefault(dashboard.id, e)}
                        title="Set as default"
                      >
                        ⭐
                      </button>
                    )}
                    <button
                      className="btn-icon"
                      onClick={(e) => handleDuplicateDashboard(dashboard.id, e)}
                      title="Duplicate"
                    >
                      📋
                    </button>
                    <button
                      className="btn-icon delete"
                      onClick={(e) => handleDeleteDashboard(dashboard.id, e)}
                      title="Delete"
                    >
                      🗑️
                    </button>
                  </div>
                </div>
                <div className="dashboard-item-meta">
                  <span className={`badge ${dashboard.is_default ? 'primary' : 'secondary'}`}>
                    {dashboard.is_default ? 'Default' : 'Custom'}
                  </span>
                  <span className="dashboard-item-date">
                    {new Date(dashboard.updated_at).toLocaleDateString()}
                  </span>
                </div>
              </div>
            ))
          )}
        </div>
      )}

      {showCreateModal && (
        <div className="modal-overlay" onClick={() => setShowCreateModal(false)}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <div className="modal-header">
              <h3>Create New Dashboard</h3>
              <button className="close-btn" onClick={() => setShowCreateModal(false)}>×</button>
            </div>
            <form onSubmit={handleCreateDashboard}>
              <div className="form-group">
                <label htmlFor="name">Dashboard Name</label>
                <input
                  type="text"
                  id="name"
                  value={newDashboardName}
                  onChange={(e) => setNewDashboardName(e.target.value)}
                  placeholder="My Dashboard"
                  required
                  autoFocus
                />
              </div>
              <div className="modal-actions">
                <button type="button" className="btn-secondary" onClick={() => setShowCreateModal(false)}>
                  Cancel
                </button>
                <button type="submit" className="btn-primary">
                  Create Dashboard
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}