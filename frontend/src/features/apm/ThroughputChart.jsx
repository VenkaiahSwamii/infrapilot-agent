import React from 'react';
import { Activity, Zap } from 'lucide-react';

export default function ThroughputChart({ rps }) {
  const points = [32, 38, 41, 44, 42, 45, 40, 42];

  return (
    <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '12px', padding: '20px' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
        <h4 style={{ margin: 0, color: '#f0f6fc', display: 'flex', alignItems: 'center', gap: '8px', fontSize: '15px' }}>
          <Zap size={18} color="#a855f7" /> Throughput &amp; Requests Per Second (RPS)
        </h4>
        <span style={{ fontSize: '12px', color: '#a855f7', fontWeight: 600 }}>
          {rps || 42} RPS
        </span>
      </div>

      <div style={{ height: '140px', display: 'flex', alignItems: 'flex-end', gap: '12px', padding: '10px 0', borderBottom: '1px solid #30363d' }}>
        {points.map((val, idx) => (
          <div key={idx} style={{ flex: 1, display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '6px' }}>
            <div style={{
              width: '100%',
              height: `${(val / 50) * 100}%`,
              backgroundColor: '#a855f7',
              borderRadius: '4px 4px 0 0',
              opacity: 0.85
            }} />
            <span style={{ fontSize: '10px', color: '#8b949e' }}>t-{8 - idx}m</span>
          </div>
        ))}
      </div>
    </div>
  );
}
