import React, { useState } from 'react';
import { TrendingUp, Clock } from 'lucide-react';
import HistoryChart from '../../charts/HistoryChart.jsx';

export default function Trends({ trends, range, onRangeChange }) {
  const points = trends || [];

  const timestamps = points.map(p => p.timestamp || '');
  const cpuData = points.map(p => p.cpu_usage || 0);
  const memData = points.map(p => p.memory_usage || 0);
  const alertData = points.map(p => p.alert_count || 0);

  const ranges = [
    { id: '1h', label: 'Last Hour' },
    { id: '1d', label: 'Last Day' },
    { id: '1w', label: 'Last Week' },
    { id: '1m', label: 'Last Month' },
    { id: '1y', label: 'Last Year' },
  ];

  return (
    <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '12px', padding: '24px' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
        <div>
          <h3 style={{ margin: 0, color: '#f0f6fc', fontSize: '18px', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '8px' }}>
            <TrendingUp size={20} color="#58a6ff" />
            Historical Trend Analysis
          </h3>
          <span style={{ fontSize: '13px', color: '#8b949e' }}>Cross-infrastructure metric progression over time</span>
        </div>

        {/* Time Range Selector */}
        <div style={{ display: 'flex', gap: '6px', background: '#0d1117', padding: '4px', borderRadius: '8px', border: '1px solid #30363d' }}>
          {ranges.map((r) => (
            <button
              key={r.id}
              onClick={() => onRangeChange(r.id)}
              style={{
                backgroundColor: range === r.id ? '#1f6feb' : 'transparent',
                color: range === r.id ? '#fff' : '#8b949e',
                border: 'none',
                padding: '6px 12px',
                borderRadius: '6px',
                cursor: 'pointer',
                fontSize: '12px',
                fontWeight: 600
              }}
            >
              {r.label}
            </button>
          ))}
        </div>
      </div>

      <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
        <div>
          <h4 style={{ color: '#f0f6fc', margin: '0 0 12px 0', fontSize: '14px' }}>📈 CPU Utilization Trend (%)</h4>
          <HistoryChart labels={timestamps} data={cpuData} title="CPU Usage" color="#58a6ff" />
        </div>

        <div>
          <h4 style={{ color: '#f0f6fc', margin: '0 0 12px 0', fontSize: '14px' }}>📈 Memory Utilization Trend (%)</h4>
          <HistoryChart labels={timestamps} data={memData} title="Memory Usage" color="#a855f7" />
        </div>

        <div>
          <h4 style={{ color: '#f0f6fc', margin: '0 0 12px 0', fontSize: '14px' }}>📈 Alert Frequency & Anomalies Timeline</h4>
          <HistoryChart labels={timestamps} data={alertData} title="Alert Count" color="#f78166" />
        </div>
      </div>
    </div>
  );
}
