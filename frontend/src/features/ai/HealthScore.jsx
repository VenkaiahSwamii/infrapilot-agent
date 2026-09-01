import React from 'react';
import { Heart, Activity, AlertTriangle, ShieldCheck } from 'lucide-react';

export default function HealthScore({ scores, loading }) {
  if (loading) {
    return <div style={{ color: '#94a3b8', padding: '20px' }}>Analyzing health status...</div>;
  }

  if (!scores || scores.length === 0) {
    return (
      <div style={{
        padding: '24px',
        backgroundColor: '#0f172a',
        borderRadius: '12px',
        border: '1px solid #1e293b',
        textAlign: 'center',
        color: '#94a3b8'
      }}>
        No health metrics collected yet.
      </div>
    );
  }

  const getScoreColor = (score) => {
    if (score >= 90) return '#10b981'; // Green
    if (score >= 70) return '#f59e0b'; // Amber
    return '#ef4444'; // Red
  };

  const getRiskIcon = (risk) => {
    switch (risk ? risk.toLowerCase() : '') {
      case 'low':
        return <ShieldCheck size={20} color="#10b981" />;
      case 'medium':
        return <Activity size={20} color="#f59e0b" />;
      case 'high':
        return <AlertTriangle size={20} color="#ef4444" />;
      default:
        return <Heart size={20} color="#94a3b8" />;
    }
  };

  return (
    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))', gap: '20px' }}>
      {scores.map((hs) => {
        const color = getScoreColor(hs.score);
        const factorsList = hs.factors ? hs.factors.split(', ') : [];

        return (
          <div key={hs.id || hs.hostname} style={{
            backgroundColor: '#0f172a',
            border: '1px solid #1e293b',
            borderRadius: '12px',
            padding: '24px',
            boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)',
            transition: 'transform 0.2s ease, border-color 0.2s ease',
            position: 'relative',
            overflow: 'hidden'
          }}>
            {/* Top Border Indicator */}
            <div style={{
              position: 'absolute',
              top: 0,
              left: 0,
              right: 0,
              height: '4px',
              backgroundColor: color
            }} />

            {/* Header */}
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
              <div>
                <h3 style={{ margin: 0, fontSize: '18px', color: '#f8fafc', fontWeight: '600' }}>
                  {hs.hostname}
                </h3>
                <span style={{ fontSize: '12px', color: '#64748b' }}>
                  ID: {hs.server_id ? hs.server_id.slice(0, 8) : 'N/A'}
                </span>
              </div>
              <div style={{
                display: 'flex',
                alignItems: 'center',
                gap: '6px',
                padding: '4px 10px',
                borderRadius: '20px',
                backgroundColor: 'rgba(30, 41, 59, 0.5)',
                border: '1px solid #334155'
              }}>
                {getRiskIcon(hs.risk)}
                <span style={{ fontSize: '12px', fontWeight: '600', color: '#e2e8f0' }}>
                  {hs.risk} Risk
                </span>
              </div>
            </div>

            {/* Score Ring / Bar */}
            <div style={{ display: 'flex', alignItems: 'center', gap: '24px', marginBottom: '24px' }}>
              <div style={{
                width: '72px',
                height: '72px',
                borderRadius: '50%',
                border: `6px solid ${color}`,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                backgroundColor: 'rgba(30, 41, 59, 0.2)',
                flexShrink: 0
              }}>
                <span style={{ fontSize: '20px', fontWeight: '700', color: '#f8fafc' }}>
                  {hs.score}
                </span>
              </div>
              <div style={{ flexGrow: 1 }}>
                <div style={{ height: '8px', width: '100%', backgroundColor: '#1e293b', borderRadius: '4px', overflow: 'hidden', marginBottom: '8px' }}>
                  <div style={{ height: '100%', width: `${hs.score}%`, backgroundColor: color, borderRadius: '4px' }} />
                </div>
                <span style={{ fontSize: '12px', color: '#94a3b8' }}>
                  Overall stability index based on real-time factors.
                </span>
              </div>
            </div>

            {/* Stability Factors */}
            <div>
              <h4 style={{ margin: '0 0 12px 0', fontSize: '13px', color: '#94a3b8', fontWeight: '600', textTransform: 'uppercase', letterSpacing: '0.05em' }}>
                Health Factors
              </h4>
              <ul style={{ margin: 0, paddingLeft: '18px', display: 'flex', flexDirection: 'column', gap: '8px' }}>
                {factorsList.map((factor, i) => (
                  <li key={i} style={{ fontSize: '13px', color: '#cbd5e1', lineHeight: '1.4' }}>
                    {factor}
                  </li>
                ))}
              </ul>
            </div>
          </div>
        );
      })}
    </div>
  );
}
