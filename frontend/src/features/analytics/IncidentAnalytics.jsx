import React from 'react';
import { AlertCircle, Clock, Server, ShieldAlert } from 'lucide-react';

export default function IncidentAnalytics({ incidents }) {
  const list = incidents?.incidents || [
    {
      id: '1',
      server_name: 'api-prod-node-01',
      root_cause: 'OOM Memory Pressure',
      severity: 'CRITICAL',
      resolution_time: 720,
      occurred_at: new Date(Date.now() - 48 * 3600 * 1000).toISOString(),
    },
    {
      id: '2',
      server_name: 'db-primary-postgres',
      root_cause: 'Disk I/O Bottleneck',
      severity: 'HIGH',
      resolution_time: 480,
      occurred_at: new Date(Date.now() - 120 * 3600 * 1000).toISOString(),
    },
  ];

  return (
    <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '12px', padding: '24px' }}>
      <div style={{ marginBottom: '20px' }}>
        <h3 style={{ margin: 0, color: '#f0f6fc', fontSize: '18px', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '8px' }}>
          <ShieldAlert size={20} color="#f78166" />
          Incident Analytics & Root Cause Distribution
        </h3>
        <span style={{ fontSize: '13px', color: '#8b949e' }}>Historical incident frequency, resolution MTTR, and most affected infrastructure</span>
      </div>

      <div style={{ overflowX: 'auto' }}>
        <table style={{ width: '100%', borderCollapse: 'collapse', textAlign: 'left', fontSize: '13px' }}>
          <thead>
            <tr style={{ borderBottom: '1px solid #30363d', color: '#8b949e' }}>
              <th style={{ padding: '12px' }}>AFFECTED SERVER</th>
              <th style={{ padding: '12px' }}>ROOT CAUSE</th>
              <th style={{ padding: '12px' }}>SEVERITY</th>
              <th style={{ padding: '12px' }}>RESOLUTION TIME</th>
              <th style={{ padding: '12px' }}>OCCURRED AT</th>
            </tr>
          </thead>
          <tbody>
            {list.map((item) => (
              <tr key={item.id} style={{ borderBottom: '1px solid #21262d', color: '#c9d1d9' }}>
                <td style={{ padding: '12px', fontWeight: 600, color: '#f0f6fc' }}>{item.server_name}</td>
                <td style={{ padding: '12px', color: '#58a6ff' }}>{item.root_cause}</td>
                <td style={{ padding: '12px' }}>
                  <span style={{
                    padding: '3px 8px',
                    borderRadius: '4px',
                    fontSize: '11px',
                    fontWeight: 700,
                    backgroundColor: item.severity === 'CRITICAL' ? 'rgba(247, 129, 102, 0.15)' : 'rgba(255, 166, 87, 0.15)',
                    color: item.severity === 'CRITICAL' ? '#f78166' : '#ffa657',
                  }}>
                    {item.severity}
                  </span>
                </td>
                <td style={{ padding: '12px' }}>{Math.round(item.resolution_time / 60)} Minutes</td>
                <td style={{ padding: '12px', color: '#8b949e', fontSize: '12px' }}>{new Date(item.occurred_at).toLocaleString()}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
