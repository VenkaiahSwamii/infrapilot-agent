import React from 'react';
import { Gauge, AlertCircle, CheckCircle2 } from 'lucide-react';

export default function Quotas({ quotas, quotaReport, onSave }) {
  const report = quotaReport?.quota_report || {
    machines: { current_usage: 12, max_limit: 100, percentage: 12.0, warning_trigger: false },
    users: { current_usage: 8, max_limit: 25, percentage: 32.0, warning_trigger: false },
    storage_gb: { current_usage: 120, max_limit: 500, percentage: 24.0, warning_trigger: false },
    reports: { current_usage: 14, max_limit: 100, percentage: 14.0, warning_trigger: false },
  };

  const getBarColor = (pct) => {
    if (pct >= 90) return '#f78166';
    if (pct >= 80) return '#ffa657';
    return '#3fb950';
  };

  return (
    <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '12px', padding: '24px' }}>
      <div style={{ marginBottom: '20px' }}>
        <h3 style={{ margin: 0, color: '#f0f6fc', fontSize: '18px', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '8px' }}>
          <Gauge size={20} color="#ffa657" />
          Tenant Resource Quotas & Usage Warnings
        </h3>
        <span style={{ fontSize: '13px', color: '#8b949e' }}>Real-time consumption tracking against plan allocation limits</span>
      </div>

      <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
        {Object.entries(report).map(([key, item]) => {
          const barColor = getBarColor(item.percentage);
          return (
            <div key={key} style={{ background: '#0d1117', border: '1px solid #30363d', borderRadius: '8px', padding: '16px' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '8px' }}>
                <span style={{ color: '#f0f6fc', fontWeight: 600, textTransform: 'uppercase', fontSize: '13px' }}>
                  {key.replace('_', ' ')}
                </span>
                <span style={{ fontSize: '13px', color: '#8b949e' }}>
                  <strong style={{ color: '#f0f6fc' }}>{item.current_usage}</strong> / {item.max_limit === 0 ? 'Unlimited' : item.max_limit} ({item.percentage.toFixed(1)}%)
                </span>
              </div>

              {/* Progress Bar */}
              <div style={{ width: '100%', height: '8px', background: '#21262d', borderRadius: '4px', overflow: 'hidden' }}>
                <div style={{ width: `${Math.min(item.percentage, 100)}%`, height: '100%', background: barColor, transition: 'width 0.3s' }} />
              </div>

              {item.warning_trigger && (
                <div style={{ display: 'flex', alignItems: 'center', gap: '6px', color: '#ffa657', fontSize: '12px', marginTop: '8px' }}>
                  <AlertCircle size={14} />
                  <span>Warning: Resource consumption is at {item.percentage.toFixed(1)}% of allocated plan quota.</span>
                </div>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
}
