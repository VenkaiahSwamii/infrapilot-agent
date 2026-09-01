import React from 'react';
import { Calendar, TrendingUp, AlertTriangle } from 'lucide-react';

export default function Predictions({ predictions, loading }) {
  if (loading) {
    return <div style={{ color: '#94a3b8', padding: '20px' }}>Generating predictions...</div>;
  }

  if (!predictions || predictions.length === 0) {
    return (
      <div style={{
        padding: '24px',
        backgroundColor: '#0f172a',
        borderRadius: '12px',
        border: '1px solid #1e293b',
        textAlign: 'center',
        color: '#94a3b8'
      }}>
        No resource exhaustion predictions registered.
      </div>
    );
  }

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
      {predictions.map((p) => {
        const progressPct = (p.current_value / p.expected_value) * 100;
        const color = progressPct >= 85 ? '#ef4444' : progressPct >= 70 ? '#f59e0b' : '#3b82f6';

        return (
          <div key={p.id || p.hostname + p.metric_name} style={{
            backgroundColor: '#0f172a',
            border: '1px solid #1e293b',
            borderRadius: '12px',
            padding: '20px',
            boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)',
            display: 'flex',
            flexDirection: 'column',
            gap: '16px'
          }}>
            {/* Header */}
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
              <div>
                <span style={{
                  fontSize: '11px',
                  fontWeight: '600',
                  color: '#94a3b8',
                  textTransform: 'uppercase',
                  letterSpacing: '0.05em'
                }}>
                  {p.hostname}
                </span>
                <h3 style={{ margin: '4px 0 0 0', fontSize: '18px', color: '#f8fafc', fontWeight: '600' }}>
                  {p.metric_name} Saturation
                </h3>
              </div>
              <div style={{
                display: 'flex',
                alignItems: 'center',
                gap: '8px',
                padding: '4px 12px',
                borderRadius: '6px',
                backgroundColor: 'rgba(239, 68, 68, 0.1)',
                border: '1px solid rgba(239, 68, 68, 0.2)'
              }}>
                <AlertTriangle size={16} color="#ef4444" />
                <span style={{ fontSize: '12px', fontWeight: '600', color: '#ef4444' }}>
                  Action Needed
                </span>
              </div>
            </div>

            {/* Metrics Bar */}
            <div>
              <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '13px', marginBottom: '8px' }}>
                <span style={{ color: '#94a3b8' }}>
                  Current: <strong style={{ color: '#cbd5e1' }}>{p.current_value}%</strong>
                </span>
                <span style={{ color: '#ef4444' }}>
                  Critical limit: <strong style={{ color: '#ef4444' }}>{p.expected_value}%</strong>
                </span>
              </div>
              <div style={{ height: '10px', width: '100%', backgroundColor: '#1e293b', borderRadius: '5px', overflow: 'hidden' }}>
                <div style={{ height: '100%', width: `${progressPct}%`, backgroundColor: color, borderRadius: '5px' }} />
              </div>
            </div>

            {/* Expected Exhaustion */}
            <div style={{
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'space-between',
              paddingTop: '12px',
              borderTop: '1px solid #1e293b'
            }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '8px', color: '#94a3b8', fontSize: '13px' }}>
                <Calendar size={16} />
                <span>Expected saturation in:</span>
                <strong style={{ color: '#f59e0b' }}>{p.timeframe_days} days</strong>
              </div>
              <div style={{ display: 'flex', alignItems: 'center', gap: '6px', color: '#3b82f6', fontSize: '13px' }}>
                <TrendingUp size={16} />
                <span>Linear Trend Projection</span>
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
}
