import React from 'react';
import { Sparkles, Terminal, CheckCircle } from 'lucide-react';

export default function Recommendations({ recommendations, loading }) {
  if (loading) {
    return <div style={{ color: '#94a3b8', padding: '20px' }}>Generating recommendations...</div>;
  }

  if (!recommendations || recommendations.length === 0) {
    return (
      <div style={{
        padding: '24px',
        backgroundColor: '#0f172a',
        borderRadius: '12px',
        border: '1px solid #1e293b',
        textAlign: 'center',
        color: '#94a3b8'
      }}>
        No active recommendations.
      </div>
    );
  }

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
      {recommendations.map((r) => {
        const priorityColor = r.priority === 'High' || r.priority === 'Critical' ? '#ef4444' : '#f59e0b';
        return (
          <div key={r.id} style={{
            backgroundColor: '#0f172a',
            border: '1px solid #1e293b',
            borderRadius: '12px',
            padding: '20px',
            boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)',
            position: 'relative'
          }}>
            {/* Priority Indicator */}
            <span style={{
              position: 'absolute',
              top: '20px',
              right: '20px',
              fontSize: '11px',
              fontWeight: '600',
              textTransform: 'uppercase',
              color: priorityColor,
              backgroundColor: priorityColor + '1a',
              border: `1px solid ${priorityColor}33`,
              padding: '2px 8px',
              borderRadius: '4px'
            }}>
              {r.priority}
            </span>

            {/* Title & Host */}
            <div style={{ marginBottom: '12px' }}>
              <span style={{ fontSize: '12px', color: '#64748b', fontWeight: '500' }}>
                Recommendation for: <strong style={{ color: '#94a3b8' }}>{r.hostname}</strong>
              </span>
              <h3 style={{ margin: '4px 0 0 0', fontSize: '16px', color: '#f8fafc', fontWeight: '600' }}>
                {r.title}
              </h3>
            </div>

            {/* Explanation / Reason */}
            <p style={{
              fontSize: '14px',
              color: '#cbd5e1',
              margin: '0 0 16px 0',
              lineHeight: '1.5'
            }}>
              {r.reason}
            </p>

            {/* Evidence / Reason Box */}
            <div style={{
              backgroundColor: 'rgba(30, 41, 59, 0.4)',
              border: '1px solid #1e293b',
              borderRadius: '8px',
              padding: '12px 16px',
              display: 'flex',
              alignItems: 'center',
              gap: '12px'
            }}>
              <Sparkles size={16} color="#3b82f6" />
              <span style={{ fontSize: '13px', color: '#94a3b8' }}>
                Evidence: Metrics values exceeded threshold indicators.
              </span>
            </div>
          </div>
        );
      })}
    </div>
  );
}
