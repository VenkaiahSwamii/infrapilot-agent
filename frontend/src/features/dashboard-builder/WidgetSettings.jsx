import React, { useState, useEffect } from 'react';

export default function WidgetSettings({ widget, onUpdate, onClose }) {
  const [formData, setFormData] = useState({
    title: widget?.title || '',
    width: widget?.width || 4,
    height: widget?.height || 3,
    refresh: 5,
    chart: 'line',
    period: '1h',
    threshold: 80,
    machine: 'all',
  });

  useEffect(() => {
    if (widget?.config) {
      try {
        const config = typeof widget.config === 'string' ? JSON.parse(widget.config) : widget.config;
        setFormData(prev => ({
          ...prev,
          ...config,
          title: widget.title || prev.title,
          width: widget.width || prev.width,
          height: widget.height || prev.height,
        }));
      } catch (err) {
        console.error('Failed to parse widget config:', err);
      }
    }
  }, [widget]);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    
    const config = {
      refresh: parseInt(formData.refresh),
      chart: formData.chart,
      period: formData.period,
      threshold: parseInt(formData.threshold),
      machine: formData.machine,
    };

    onUpdate({
      title: formData.title,
      width: parseInt(formData.width),
      height: parseInt(formData.height),
      config,
    });
  };

  if (!widget) return null;

  return (
    <div className="widget-settings-overlay" onClick={onClose}>
      <div className="widget-settings-modal" onClick={e => e.stopPropagation()}>
        <div className="widget-settings-header">
          <h3>Widget Settings</h3>
          <button className="close-btn" onClick={onClose}>×</button>
        </div>

        <form onSubmit={handleSubmit} className="widget-settings-form">
          <div className="form-group">
            <label htmlFor="title">Title</label>
            <input
              type="text"
              id="title"
              name="title"
              value={formData.title}
              onChange={handleChange}
              placeholder="Widget title"
            />
          </div>

          <div className="form-row">
            <div className="form-group">
              <label htmlFor="width">Width (columns)</label>
              <input
                type="number"
                id="width"
                name="width"
                min="2"
                max="12"
                value={formData.width}
                onChange={handleChange}
              />
            </div>

            <div className="form-group">
              <label htmlFor="height">Height (rows)</label>
              <input
                type="number"
                id="height"
                name="height"
                min="2"
                max="10"
                value={formData.height}
                onChange={handleChange}
              />
            </div>
          </div>

          <div className="form-group">
            <label htmlFor="machine">Machine/Server</label>
            <select
              id="machine"
              name="machine"
              value={formData.machine}
              onChange={handleChange}
            >
              <option value="all">All Machines</option>
              <option value="server-01">server-01</option>
              <option value="server-02">server-02</option>
              <option value="server-03">server-03</option>
            </select>
          </div>

          <div className="form-group">
            <label htmlFor="refresh">Refresh Interval (seconds)</label>
            <input
              type="number"
              id="refresh"
              name="refresh"
              min="5"
              max="300"
              step="5"
              value={formData.refresh}
              onChange={handleChange}
            />
          </div>

          <div className="form-group">
            <label htmlFor="chart">Chart Type</label>
            <select
              id="chart"
              name="chart"
              value={formData.chart}
              onChange={handleChange}
            >
              <option value="line">Line Chart</option>
              <option value="area">Area Chart</option>
              <option value="bar">Bar Chart</option>
              <option value="gauge">Gauge</option>
              <option value="number">Number</option>
              <option value="table">Table</option>
            </select>
          </div>

          <div className="form-group">
            <label htmlFor="period">Time Period</label>
            <select
              id="period"
              name="period"
              value={formData.period}
              onChange={handleChange}
            >
              <option value="5m">Last 5 minutes</option>
              <option value="15m">Last 15 minutes</option>
              <option value="1h">Last 1 hour</option>
              <option value="6h">Last 6 hours</option>
              <option value="24h">Last 24 hours</option>
              <option value="7d">Last 7 days</option>
            </select>
          </div>

          <div className="form-group">
            <label htmlFor="threshold">Warning Threshold (%)</label>
            <input
              type="number"
              id="threshold"
              name="threshold"
              min="0"
              max="100"
              value={formData.threshold}
              onChange={handleChange}
            />
          </div>

          <div className="form-actions">
            <button type="button" className="btn-secondary" onClick={onClose}>
              Cancel
            </button>
            <button type="submit" className="btn-primary">
              Save Changes
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}