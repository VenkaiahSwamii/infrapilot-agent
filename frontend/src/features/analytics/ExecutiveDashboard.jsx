import React from 'react';
import { Building2, Server, Layers, Cpu, AlertTriangle, ShieldCheck, Activity, Award } from 'lucide-react';

export default function ExecutiveDashboard({ overview }) {
  const data = overview || {
    total_servers: 480,
    docker_containers: 3200,
    k8s_clusters: 28,
    critical_incidents: 4,
    availability_pct: 99.97,
  };

  const widgets = [
    { label: 'INFRASTRUCTURE HEALTH', value: '96%', color: '#3fb950', icon: Activity, desc: 'Global system health' },
    { label: 'ORGANIZATIONS', value: '12', color: '#58a6ff', icon: Building2, desc: 'Enterprise tenants' },
    { label: 'SERVERS MONITORED', value: data.total_servers || '480', color: '#58a6ff', icon: Server, desc: 'Active agent hosts' },
    { label: 'CONTAINERS', value: (data.docker_containers || 3200).toLocaleString(), color: '#a855f7', icon: Layers, desc: 'Docker container instances' },
    { label: 'KUBERNETES CLUSTERS', value: data.k8s_clusters || '28', color: '#ffa657', icon: Cpu, desc: 'Production K8s clusters' },
    { label: 'CRITICAL ALERTS', value: data.critical_incidents || '4', color: '#f78166', icon: AlertTriangle, desc: 'Requires immediate action' },
    { label: 'SLA COMPLIANCE', value: `${data.availability_pct || 99.97}%`, color: '#3fb950', icon: ShieldCheck, desc: 'Monthly uptime rating' },
  ];

  return (
    <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '12px', padding: '24px' }}>
      <div style={{ marginBottom: '20px' }}>
        <h3 style={{ margin: 0, color: '#f0f6fc', fontSize: '18px', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '8px' }}>
          <Award size={20} color="#58a6ff" />
          Executive KPI Dashboard
        </h3>
        <span style={{ fontSize: '13px', color: '#8b949e' }}>High-level strategic insights, organization scale & operational compliance</span>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))', gap: '20px' }}>
        {widgets.map((w, idx) => {
          const Icon = w.icon;
          return (
            <div key={idx} style={{
              background: '#0d1117',
              border: '1px solid #30363d',
              borderRadius: '10px',
              padding: '20px',
              display: 'flex',
              flexDirection: 'column',
              justifyContent: 'space-between'
            }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', color: '#8b949e', fontSize: '11px', fontWeight: 700 }}>
                <span>{w.label}</span>
                <Icon size={18} color={w.color} />
              </div>

              <div style={{ fontSize: '32px', fontWeight: 800, color: '#f0f6fc', margin: '14px 0 4px 0' }}>
                {w.value}
              </div>

              <div style={{ fontSize: '12px', color: '#8b949e' }}>{w.desc}</div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
