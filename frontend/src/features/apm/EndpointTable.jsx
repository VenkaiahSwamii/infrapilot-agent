import React from 'react';
import { AlertCircle, Flame } from 'lucide-react';

export default function EndpointTable({ endpoints, slowEndpoints, topAPIs }) {
  const allEndpoints = endpoints || [];
  const slowList = slowEndpoints || [];

  return (
    <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '20px' }}>
      {/* Top Endpoints */}
      <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '12px', padding: '20px' }}>
        <h4 style={{ margin: '0 0 16px 0', color: '#f0f6fc', display: 'flex', alignItems: 'center', gap: '8px', fontSize: '15px' }}>
          <Flame size={18} color="#58a6ff" /> Top APIs &amp; Traffic Volume (%)
        </h4>
        <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
          {allEndpoints.slice(0, 5).map((ep, idx) => (
            <div key={idx} style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '12px' }}>
                <span style={{ fontWeight: 600, color: '#f0f6fc' }}>
                  <span style={{ color: ep.Method === 'GET' ? '#58a6ff' : '#a855f7', marginRight: '6px' }}>{ep.Method}</span>
                  {ep.Endpoint}
                </span>
                <span style={{ color: '#8b949e', fontWeight: 600 }}>{ep.Percentage}% ({ep.AvgLatencyMs} ms)</span>
              </div>
              <div style={{ background: '#0d1117', height: '6px', borderRadius: '3px', overflow: 'hidden' }}>
                <div style={{ width: `${ep.Percentage * 3}%`, height: '100%', backgroundColor: '#58a6ff' }} />
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Slow Endpoints Flagged */}
      <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '12px', padding: '20px' }}>
        <h4 style={{ margin: '0 0 16px 0', color: '#f0f6fc', display: 'flex', alignItems: 'center', gap: '8px', fontSize: '15px' }}>
          <AlertCircle size={18} color="#ffa657" /> Flagged Slow Endpoints (&gt;500ms)
        </h4>
        <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
          {slowList.map((ep, idx) => (
            <div key={idx} style={{ background: '#0d1117', border: '1px solid rgba(255,166,87,0.3)', borderRadius: '8px', padding: '12px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <div>
                <span style={{ fontWeight: 700, color: '#ffa657', fontSize: '13px', display: 'block' }}>
                  {ep.Method} {ep.Endpoint}
                </span>
                <span style={{ fontSize: '11px', color: '#8b949e' }}>Requests: {ep.RequestCount}</span>
              </div>
              <span style={{ background: 'rgba(255,166,87,0.2)', color: '#ffa657', padding: '4px 10px', borderRadius: '4px', fontSize: '12px', fontWeight: 700 }}>
                {ep.AvgLatencyMs} ms
              </span>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
