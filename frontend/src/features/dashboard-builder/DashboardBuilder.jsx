import React, { useState, useEffect } from 'react';
import GridCanvas from './GridCanvas';
import WidgetLibrary from './WidgetLibrary';
import WidgetSettings from './WidgetSettings';
import DashboardManager from './DashboardManager';
import dashboardsApi from '../../api/dashboards';

export default function DashboardBuilder() {
  const [currentDashboard, setCurrentDashboard] = useState(null);
  const [widgets, setWidgets] = useState([]);
  const [selectedWidgetId, setSelectedWidgetId] = useState(null);
  const [showWidgetLibrary, setShowWidgetLibrary] = useState(false);
  const [showDashboardManager, setShowDashboardManager] = useState(false);
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);

  const selectedWidget = widgets.find(w => w.id === selectedWidgetId);

  useEffect(() => {
    loadDefaultDashboard();
  }, []);

  useEffect(() => {
    if (currentDashboard) {
      loadWidgets(currentDashboard.id);
    }
  }, [currentDashboard]);

  const loadDefaultDashboard = async () => {
    try {
      setLoading(true);
      const response = await dashboardsApi.getDefaultDashboard();
      if (response.success && response.data) {
        setCurrentDashboard(response.data);
      }
    } catch (err) {
      console.error('Failed to load default dashboard:', err);
    } finally {
      setLoading(false);
    }
  };

  const loadWidgets = async (dashboardId) => {
    try {
      const response = await dashboardsApi.getWidgets(dashboardId);
      if (response.success) {
        setWidgets(response.data || []);
      }
    } catch (err) {
      console.error('Failed to load widgets:', err);
    }
  };

  const handleSelectDashboard = (dashboard) => {
    setCurrentDashboard(dashboard);
    setShowDashboardManager(false);
    setSelectedWidgetId(null);
  };

  const handleAddWidget = async (widgetType, label, x = 0, y = 0) => {
    if (!currentDashboard) {
      alert('Please create or select a dashboard first');
      return;
    }

    try {
      const response = await dashboardsApi.createWidget(currentDashboard.id, {
        widget_type: widgetType,
        title: label,
        x,
        y,
        width: 4,
        height: 3,
      });

      if (response.success) {
        setWidgets([...widgets, response.data]);
        setShowWidgetLibrary(false);
      }
    } catch (err) {
      console.error('Failed to add widget:', err);
    }
  };

  const handleUpdateWidget = async (widgetId, updates) => {
    try {
      const response = await dashboardsApi.updateWidget(widgetId, updates);
      if (response.success) {
        setWidgets(widgets.map(w => w.id === widgetId ? { ...w, ...updates } : w));
        setSelectedWidgetId(null);
      }
    } catch (err) {
      console.error('Failed to update widget:', err);
    }
  };

  const handleDeleteWidget = async (widgetId) => {
    if (!window.confirm('Are you sure you want to delete this widget?')) {
      return;
    }

    try {
      await dashboardsApi.deleteWidget(widgetId);
      setWidgets(widgets.filter(w => w.id !== widgetId));
      setSelectedWidgetId(null);
    } catch (err) {
      console.error('Failed to delete widget:', err);
    }
  };

  const handleSaveDashboard = async () => {
    if (!currentDashboard) return;

    try {
      setSaving(true);
      await dashboardsApi.updateDashboard(currentDashboard.id, {
        name: currentDashboard.name,
        widgets: widgets.map(w => ({
          id: w.id,
          x: w.x,
          y: w.y,
          width: w.width,
          height: w.height,
          config: w.config,
        })),
      });
      alert('Dashboard saved successfully!');
    } catch (err) {
      console.error('Failed to save dashboard:', err);
      alert('Failed to save dashboard');
    } finally {
      setSaving(false);
    }
  };

  const handleCreateNewDashboard = async (name) => {
    try {
      const response = await dashboardsApi.createDashboard({ name });
      if (response.success) {
        setCurrentDashboard(response.data);
        setWidgets([]);
        setShowDashboardManager(false);
      }
    } catch (err) {
      console.error('Failed to create dashboard:', err);
    }
  };

  return (
    <div className="dashboard-builder">
      <div className="dashboard-builder-header">
        <div className="header-left">
          <button
            className="btn-icon"
            onClick={() => setShowDashboardManager(true)}
            title="Dashboard Manager"
          >
            ☰
          </button>
          <h1>{currentDashboard?.name || 'Dashboard Builder'}</h1>
        </div>
        <div className="header-actions">
          <button
            className="btn-secondary"
            onClick={() => setShowWidgetLibrary(true)}
          >
            + Add Widget
          </button>
          <button
            className="btn-primary"
            onClick={handleSaveDashboard}
            disabled={saving || !currentDashboard}
          >
            {saving ? 'Saving...' : 'Save Dashboard'}
          </button>
        </div>
      </div>

      {loading ? (
        <div className="loading">Loading dashboard...</div>
      ) : (
        <div className="dashboard-builder-content">
          {showDashboardManager ? (
            <DashboardManager
              onSelectDashboard={handleSelectDashboard}
              onCreateNew={handleCreateNewDashboard}
            />
          ) : (
            <>
              {!currentDashboard ? (
                <div className="empty-state">
                  <h2>No Dashboard Selected</h2>
                  <p>Select a dashboard from the manager or create a new one to get started.</p>
                  <button
                    className="btn-primary"
                    onClick={() => setShowDashboardManager(true)}
                  >
                    Open Dashboard Manager
                  </button>
                </div>
              ) : (
                <GridCanvas
                  widgets={widgets}
                  onUpdateWidget={handleUpdateWidget}
                  onDeleteWidget={handleDeleteWidget}
                  onAddWidget={handleAddWidget}
                  onSelectWidget={setSelectedWidgetId}
                  selectedWidgetId={selectedWidgetId}
                />
              )}
            </>
          )}
        </div>
      )}

      {showWidgetLibrary && (
        <WidgetLibrary
          onAddWidget={handleAddWidget}
          onClose={() => setShowWidgetLibrary(false)}
        />
      )}

      {selectedWidget && (
        <WidgetSettings
          widget={selectedWidget}
          onUpdate={(updates) => handleUpdateWidget(selectedWidget.id, updates)}
          onClose={() => setSelectedWidgetId(null)}
        />
      )}
    </div>
  );
}