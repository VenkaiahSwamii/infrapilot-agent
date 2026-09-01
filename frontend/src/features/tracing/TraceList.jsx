import React from 'react';
import { Clock, AlertTriangle, Eye, ArrowRight } from 'lucide-react';

export default function TraceList({ traces, onSelectTrace }) {
  const list = traces || [];

  const getLatencyColor = (duration, hasError) => {
    if (hasError) return '#f78166';
    if (duration > 500) return '#ffa657';
    if (duration > 200) return '#e3b341';
    return '#3fb950';
  };

  return (
    <div style={{ overflowX: 'auto' }}>
      <table style={{ width: '100%', borderCollapse: 'collapse', textAlign: 'left', fontSize: '13px' }}>
        <thead>
          <tr style={{ borderBottom: '1px solid #30363d', color: '#8b949e' }}>
            <th style={{ padding: '12px' }}>TRACE ID</th>
            <th style={{ padding: '12px' }}>OPERATION / ROUTE</th>
            <th style={{ padding: '12px' }}>SERVICE</th>
            <th style={{ padding: '12px' }}>STATUS</th>
            <th style={{ padding: '12px' }}>LATENCY</th>
            <th style={{ padding: '12px' }}>TIMESTAMP</th>
            <th style={{ padding: '12px', textAlign: 'right' }}>ACTION</th>
          </tr>
        </thead>
        <tbody>
          {list.length === 0 ? (
            <tr>
              <td colSpan="7" style={{ textAlign: 'center', padding: '24px', color: '#8b949e', fontStyle: 'italic' }}>
                No distributed traces recorded.
              </td>
            </tr>
          ) : (
            list.map((t) => (
              <tr key={t.id} style={{ borderBottom: '1px solid #21262d', color: '#c9d1d9' }}>
                <td style={{ padding: '12px', fontFamily: 'monospace', color: '#58a6ff', fontSize: '12px' }}>
                  {t.trace_id?.substring(0, 18)}...
                </td>
                <td style={{ padding: '12px', fontWeight: 600, color: '#f0f6fc' }}>
                  {t.name}
                </td>
                <td style={{ padding: '12px', color: '#a855f7', fontWeight: 600 }}>
                  {t.service_name || 'api-server'}
                </td>
                <td style={{ padding: '12px' }}>
                  <span style={{
                    padding: '2px 8px',
                    borderRadius: '4px',
                    fontSize: '11px',
                    fontWeight: 700,
                    backgroundColor: t.status_code >= 400 ? 'rgba(247, 129, 102, 0.15)' : 'rgba(63, 185, 80, 0.15)',
                    color: t.status_code >= 400 ? '#f78166' : '#3fb950',
                  }}>
                    {t.status_code}
                  </span>
                </td>
                <td style={{ padding: '12px', fontWeight: 700, color: getLatencyColor(t.duration_ms, t.has_error) }}>
                  {t.duration_ms} ms {t.is_slow && <span style={{ fontSize: '10px', background: 'rgba(255,166,87,0.2)', padding: '2px 4px', borderRadius: '3px', marginLeft: '4px' }}>SLOW</span>}
                </td>
                <td style={{ padding: '12px', color: '#8b949e', fontSize: '12px' }}>
                  {new Date(t.timestamp).toLocaleString()}
                </td>
                <td style={{ padding: '12px', textAlign: 'right' }}>
                  <button
                    onClick={() => onSelectTrace(t)}
                    style={{ background: '#21262d', color: '#58a6ff', border: '1px solid #30363d', padding: '4px 10px', borderRadius: '4px', cursor: 'pointer', fontSize: '12px', display: 'inline-flex', alignItems: 'center', gap: '4px' }}
                  >
                    <Eye size={14} /> Timeline
                  </button>
                </td>
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  );
}
