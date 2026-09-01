import React from 'react';
import { Clock, TrendingUp } from 'lucide-react';

export default function LatencyChart({ avgLatency, p95Latency }) {
  const points = [45, 48, 52, 60, 58, 65, 54, 58.4];

  return (
    <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '12px', padding: '20px' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
        <h4 style={{ margin: 0, color: '#f0f6fc', display: 'flex', alignItems: 'center', gap: '8px', fontSize: '15px' }}>
          <Clock size={18} color="#58a6ff" /> API Response Time Latency Trend (ms)
        </h4>
        <span style={{ fontSize: '12px', color: '#3fb950', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '4px' }}>
          <TrendingUp size={14} /> Avg: {avgLatency || 58.4} ms
        </span>
      </div>

      <div style={{ height: '140px', display: 'flex', alignItems: 'flex-end', gap: '12px', padding: '10px 0', borderBottom: '1px solid #30363d' }}>
        {points.map((val, idx) => (
          <div key={idx} style={{ flex: 1, display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '6px' }}>
            <div style={{
              width: '100%',
              height: `${(val / 80) * 100}%`,
              backgroundColor: val > 60 ? '#ffa657' : '#58a6ff',
              borderRadius: '4px 4px 0 0',
              transition: 'height 0.3s ease'
            }} />
            <span style={{ fontSize: '10px', color: '#8b949e' }}>t-{8 - idx}m</span>
          </div>
        ))}
      </div>
    </div>
  );
}
