import React from 'react';
import { AlertTriangle, Brain, CheckCircle } from 'lucide-react';

export default function ErrorRateCard({ errorRate, recommendations }) {
  const recs = recommendations || [
    {
      issue: 'API latency increased by 43% on /api/v1/reports',
      possible_causes: ['Slow PostgreSQL unindexed query', 'High CPU usage during PDF compilation'],
      recommendation: 'Optimize SQL index on historical_snapshots table and introduce Redis caching layer for /api/v1/reports.',
      target_endpoint: 'GET /api/v1/reports',
    },
  ];

  return (
    <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '12px', padding: '20px' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
        <h4 style={{ margin: 0, color: '#f0f6fc', display: 'flex', alignItems: 'center', gap: '8px', fontSize: '15px' }}>
          <Brain size={18} color="#a855f7" /> AI-Assisted APM Insights &amp; Root Cause Analysis
        </h4>
        <span style={{
          backgroundColor: errorRate > 1.0 ? 'rgba(247, 129, 102, 0.15)' : 'rgba(63, 185, 80, 0.15)',
          color: errorRate > 1.0 ? '#f78166' : '#3fb950',
          border: `1px solid ${errorRate > 1.0 ? '#f78166' : '#3fb950'}`,
          padding: '4px 10px',
          borderRadius: '6px',
          fontSize: '12px',
          fontWeight: 700
        }}>
          Error Rate: {errorRate || 0.8}%
        </span>
      </div>

      <div style={{ display: 'flex', flexDirection: 'column', gap: '14px' }}>
        {recs.map((rec, idx) => (
          <div key={idx} style={{ background: '#0d1117', border: '1px solid #30363d', borderRadius: '8px', padding: '16px' }}>
            <div style={{ color: '#f78166', fontWeight: 700, fontSize: '14px', marginBottom: '6px', display: 'flex', alignItems: 'center', gap: '6px' }}>
              <AlertTriangle size={16} /> {rec.issue}
            </div>
            <div style={{ fontSize: '12px', color: '#8b949e', marginBottom: '8px' }}>
              <strong>Possible Causes:</strong> {rec.possible_causes?.join(' • ')}
            </div>
            <div style={{ fontSize: '13px', color: '#3fb950', background: 'rgba(63,185,80,0.1)', padding: '10px', borderRadius: '6px', borderLeft: '3px solid #3fb950' }}>
              <strong>AI Recommendation:</strong> {rec.recommendation}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
