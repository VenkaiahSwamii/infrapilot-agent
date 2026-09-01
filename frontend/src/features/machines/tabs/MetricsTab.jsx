import React from 'react';

export default function MetricsTab({ metrics }) {
  if (!metrics) return <div className="tab-empty">No metrics data available</div>;
  return (
    <div className="tab-container">
      <div className="metrics-grid">
        {[
          { label: 'CPU Usage', value: metrics.cpu_usage, unit: '%', color: '#06b6d4' },
          { label: 'Memory', value: metrics.memory_usage, unit: '%', color: '#22c55e' },
          { label: 'Disk', value: metrics.disk_usage, unit: '%', color: '#eab308' },
          { label: 'Upload', value: metrics.upload_mbps, unit: 'Mbps', color: '#a78bfa' },
          { label: 'Download', value: metrics.download_mbps, unit: 'Mbps', color: '#f97316' },
          { label: 'Disk Read', value: metrics.disk_read_bps, unit: 'B/s', color: '#38bdf8' },
          { label: 'Disk Write', value: metrics.disk_write_bps, unit: 'B/s', color: '#fb7185' },
          { label: 'CPU Frequency', value: metrics.cpu_frequency_mhz, unit: 'MHz', color: '#818cf8' },
          { label: 'CPU Temperature', value: metrics.cpu_temperature, unit: '°C', color: '#ef4444' },
        ].filter(m => m.value != null).map(m => (
          <div key={m.label} className="metric-card">
            <h3>{m.label}</h3>
            <div className="metric-value" style={{ color: m.color }}>
              {typeof m.value === 'number' ? m.value.toFixed(1) : m.value}
              <span className="metric-unit">{m.unit}</span>
            </div>
            <div className="metric-bar"><div className="metric-fill" style={{ width: `${Math.min(m.value, 100)}%`, background: m.color }} /></div>
          </div>
        ))}
      </div>
      <style>{`
        .tab-container { padding: 20px; }
        .tab-empty { padding: 40px; text-align: center; color: #94a3b8; }
        .metrics-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 16px; }
        .metric-card { background: #1e293b; border: 1px solid #334155; border-radius: 10px; padding: 16px; }
        .metric-card h3 { color: #94a3b8; font-size: 13px; margin: 0 0 8px; text-transform: uppercase; letter-spacing: 0.5px; }
        .metric-value { font-size: 28px; font-weight: 700; margin-bottom: 12px; }
        .metric-unit { font-size: 14px; font-weight: 400; margin-left: 4px; opacity: 0.7; }
        .metric-bar { height: 4px; background: #334155; border-radius: 2px; overflow: hidden; }
        .metric-fill { height: 100%; border-radius: 2px; transition: width 0.5s; }
      `}</style>
    </div>
  );
}
