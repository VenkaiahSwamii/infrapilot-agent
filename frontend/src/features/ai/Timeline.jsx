import React from 'react';
import { Clock, AlertTriangle, Info, AlertOctagon, CheckCircle } from 'lucide-react';

export default function Timeline({ timelines }) {
  if (!timelines || timelines.length === 0) {
    return (
      <div className="timeline-empty" style={{ padding: '15px', color: '#94a3b8', fontStyle: 'italic' }}>
        No timeline events recorded.
      </div>
    );
  }

  const getIcon = (type) => {
    switch (type ? type.toLowerCase() : '') {
      case 'critical':
      case 'error':
        return <AlertOctagon size={16} color="#ef4444" />;
      case 'warning':
        return <AlertTriangle size={16} color="#f59e0b" />;
      case 'normal':
      case 'info':
        return <Info size={16} color="#3b82f6" />;
      case 'resolved':
      case 'success':
        return <CheckCircle size={16} color="#10b981" />;
      default:
        return <Clock size={16} color="#94a3b8" />;
    }
  };

  return (
    <div className="ai-timeline-container" style={{
      display: 'flex',
      flexDirection: 'column',
      gap: '16px',
      padding: '16px 8px',
      position: 'relative'
    }}>
      {/* Central line */}
      <div style={{
        position: 'absolute',
        left: '23px',
        top: '24px',
        bottom: '24px',
        width: '2px',
        backgroundColor: '#334155'
      }} />

      {timelines.map((item, index) => {
        const timeStr = item.Timestamp ? new Date(item.Timestamp).toLocaleTimeString() : '';
        return (
          <div key={item.id || index} style={{
            display: 'flex',
            gap: '16px',
            alignItems: 'flex-start',
            position: 'relative',
            zIndex: 1
          }}>
            {/* Dot & Icon */}
            <div style={{
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              width: '32px',
              height: '32px',
              borderRadius: '50%',
              backgroundColor: '#1e293b',
              border: '2px solid #334155',
              flexShrink: 0
            }}>
              {getIcon(item.Type || item.type)}
            </div>

            {/* Content card */}
            <div style={{
              flexGrow: 1,
              backgroundColor: '#0f172a',
              border: '1px solid #1e293b',
              borderRadius: '8px',
              padding: '12px 16px',
              boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)'
            }}>
              <div style={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                marginBottom: '4px'
              }}>
                <span style={{
                  fontSize: '11px',
                  fontWeight: '600',
                  textTransform: 'uppercase',
                  color: (item.Type || item.type) === 'Critical' ? '#ef4444' : (item.Type || item.type) === 'Warning' ? '#f59e0b' : '#3b82f6',
                  backgroundColor: (item.Type || item.type) === 'Critical' ? 'rgba(239, 68, 68, 0.1)' : (item.Type || item.type) === 'Warning' ? 'rgba(245, 158, 11, 0.1)' : 'rgba(59, 130, 246, 0.1)',
                  padding: '2px 8px',
                  borderRadius: '4px'
                }}>
                  {item.Type || item.type || 'INFO'}
                </span>
                <span style={{ fontSize: '11px', color: '#64748b' }}>
                  {timeStr}
                </span>
              </div>
              <p style={{
                fontSize: '13px',
                color: '#e2e8f0',
                margin: 0,
                lineHeight: '1.5'
              }}>
                {item.Message || item.message || 'No description provided'}
              </p>
            </div>
          </div>
        );
      })}
    </div>
  );
}
