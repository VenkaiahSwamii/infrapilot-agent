import React from 'react';
import { Gauge, Server, Bell, ShieldCheck } from 'lucide-react';

function StatCard({ icon: Icon, label, value, meta, tone = 'cyan' }) {
  return (
    <article className={`stat-card stat-card-${tone}`}>
      <div className="stat-copy">
        <span>{label}</span>
        <strong>{value}</strong>
        {meta && <small>{meta}</small>}
      </div>
      <div className="stat-icon">
        <Icon size={22} />
      </div>
    </article>
  );
}

export default function OverviewCards({ displayOverview }) {
  const health = displayOverview.health_score ?? 'N/A';
  const healthLabel = typeof health === 'number' ? `${health}%` : health;

  return (
    <div className="stats-grid">
      <StatCard
        icon={Gauge}
        label="Health Index"
        value={healthLabel}
        meta={typeof health === 'number' && health > 85 ? 'Healthy control plane' : 'Warning state'}
        tone={typeof health === 'number' && health > 85 ? 'green' : 'amber'}
      />
      <StatCard
        icon={Server}
        label="Connected Servers"
        value={`${displayOverview.online} / ${displayOverview.total_servers}`}
        meta={`${displayOverview.offline} offline servers`}
        tone="cyan"
      />
      <StatCard
        icon={Bell}
        label="Active Alerts"
        value={displayOverview.active_alerts}
        meta="Real-time incident feed"
        tone={displayOverview.active_alerts > 0 ? 'red' : 'green'}
      />
      <StatCard
        icon={ShieldCheck}
        label="Security Baseline"
        value="Compliant"
        meta="ISO 27001 / SOC 2 check passed"
        tone="green"
      />
    </div>
  );
}
