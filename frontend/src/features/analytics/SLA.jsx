import React from 'react';
import { Activity, Clock, ShieldCheck, CheckCircle } from 'lucide-react';

export default function SLA({ sla }) {
  const data = sla || {
    availability_pct: 99.95,
    uptime_minutes: 43182,
    downtime_minutes: 18,
    mttr_minutes: 12,
    mtbf_days: 14,
    incidents_count: 3,
  };

  return (
    <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '12px', padding: '24px' }}>
      <div style={{ marginBottom: '20px' }}>
        <h3 style={{ margin: 0, color: '#f0f6fc', fontSize: '18px', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '8px' }}>
          <Activity size={20} color="#3fb950" />
          Service Level Agreement (SLA) Monitoring
        </h3>
        <span style={{ fontSize: '13px', color: '#8b949e' }}>Real-time uptime tracking, downtime minutes, MTTR, and MTBF metrics</span>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(220px, 1fr))', gap: '20px' }}>
        <div style={{ background: '#0d1117', border: '1px solid #30363d', borderRadius: '10px', padding: '20px' }}>
          <div style={{ color: '#8b949e', fontSize: '12px', fontWeight: 600 }}>AVAILABILITY %</div>
          <div style={{ fontSize: '32px', fontWeight: 800, color: '#3fb950', margin: '10px 0 4px 0' }}>
            {data.availability_pct}%
          </div>
          <div style={{ fontSize: '12px', color: '#3fb950', display: 'flex', alignItems: 'center', gap: '4px' }}>
            <CheckCircle size={14} /> Meets 99.9% Target
          </div>
        </div>

        <div style={{ background: '#0d1117', border: '1px solid #30363d', borderRadius: '10px', padding: '20px' }}>
          <div style={{ color: '#8b949e', fontSize: '12px', fontWeight: 600 }}>DOWNTIME</div>
          <div style={{ fontSize: '32px', fontWeight: 800, color: '#ffa657', margin: '10px 0 4px 0' }}>
            {data.downtime_minutes} <span style={{ fontSize: '14px', color: '#8b949e' }}>Minutes</span>
          </div>
          <div style={{ fontSize: '12px', color: '#8b949e' }}>Total outage duration</div>
        </div>

        <div style={{ background: '#0d1117', border: '1px solid #30363d', borderRadius: '10px', padding: '20px' }}>
          <div style={{ color: '#8b949e', fontSize: '12px', fontWeight: 600 }}>MTTR (MEAN TIME TO REPAIR)</div>
          <div style={{ fontSize: '32px', fontWeight: 800, color: '#58a6ff', margin: '10px 0 4px 0' }}>
            {data.mttr_minutes} <span style={{ fontSize: '14px', color: '#8b949e' }}>Minutes</span>
          </div>
          <div style={{ fontSize: '12px', color: '#3fb950' }}>Avg resolution speed</div>
        </div>

        <div style={{ background: '#0d1117', border: '1px solid #30363d', borderRadius: '10px', padding: '20px' }}>
          <div style={{ color: '#8b949e', fontSize: '12px', fontWeight: 600 }}>MTBF (MEAN TIME BETWEEN FAILURES)</div>
          <div style={{ fontSize: '32px', fontWeight: 800, color: '#a855f7', margin: '10px 0 4px 0' }}>
            {data.mtbf_days} <span style={{ fontSize: '14px', color: '#8b949e' }}>Days</span>
          </div>
          <div style={{ fontSize: '12px', color: '#8b949e' }}>Avg operational stability</div>
        </div>
      </div>
    </div>
  );
}
