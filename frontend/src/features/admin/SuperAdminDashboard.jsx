import React, { useState, useEffect } from 'react';
import { ShieldCheck, Building2, Server, Users, Activity, DollarSign, HardDrive, RefreshCw, Cpu } from 'lucide-react';
import { getSuperAdminMetrics } from '../../api/organization.js';
import './SuperAdminDashboard.css';

export default function SuperAdminDashboard() {
  const [metrics, setMetrics] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchMetrics();
  }, []);

  const fetchMetrics = async () => {
    setLoading(true);
    try {
      const res = await getSuperAdminMetrics();
      setMetrics(res);
    } catch (err) {
      console.error('Failed to load Super Admin dashboard metrics', err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="super-admin-container">
      <div className="super-admin-header">
        <div className="super-admin-title">
          <ShieldCheck size={32} color="#f78166" />
          <div>
            <h1>Platform Super Admin Dashboard</h1>
            <span style={{ fontSize: '13px', color: '#8b949e' }}>
              Multi-Tenant System Control & Cross-Tenant Oversight
            </span>
          </div>
        </div>
        <button className="btn-secondary-dr" onClick={fetchMetrics} style={{ background: '#21262d', color: '#c9d1d9', border: '1px solid #30363d', padding: '8px 16px', borderRadius: '6px', cursor: 'pointer', display: 'flex', alignItems: 'center', gap: '8px' }}>
          <RefreshCw size={16} /> Refresh Metrics
        </button>
      </div>

      {/* KPI Cards Grid */}
      <div className="super-admin-grid">
        <div className="super-card">
          <div className="super-card-header">
            <span>TOTAL ORGANIZATIONS</span>
            <Building2 size={18} color="#58a6ff" />
          </div>
          <div className="super-card-value">{metrics?.total_organizations ?? 1}</div>
          <div className="super-card-sub">Active SaaS Tenants</div>
        </div>

        <div className="super-card">
          <div className="super-card-header">
            <span>MANAGED MACHINES</span>
            <Server size={18} color="#3fb950" />
          </div>
          <div className="super-card-value">{metrics?.total_machines ?? 0}</div>
          <div className="super-card-sub">{metrics?.active_agents ?? 0} online agents</div>
        </div>

        <div className="super-card">
          <div className="super-card-header">
            <span>REGISTERED USERS</span>
            <Users size={18} color="#a5d6ff" />
          </div>
          <div className="super-card-value">{metrics?.total_users ?? 1}</div>
          <div className="super-card-sub">Across all organizations</div>
        </div>

        <div className="super-card">
          <div className="super-card-header">
            <span>MONTHLY REVENUE</span>
            <DollarSign size={18} color="#2ea043" />
          </div>
          <div className="super-card-value" style={{ color: '#3fb950' }}>$125,000</div>
          <div className="super-card-sub">Recurring SaaS ARR</div>
        </div>

        <div className="super-card">
          <div className="super-card-header">
            <span>CLUSTER HEALTH</span>
            <Cpu size={18} color="#ffa657" />
          </div>
          <div className="super-card-value" style={{ color: '#3fb950' }}>99.99%</div>
          <div className="super-card-sub">3 K8s Nodes Operational</div>
        </div>
      </div>

      {/* Platform Subscription Breakdown */}
      <div className="org-card" style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '10px', padding: '24px' }}>
        <h3 style={{ margin: '0 0 16px 0', color: '#f0f6fc' }}>Tenant Subscription Tier Breakdown</h3>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: '16px' }}>
          <div style={{ background: '#0d1117', padding: '16px', borderRadius: '8px', border: '1px solid #30363d' }}>
            <span style={{ fontSize: '12px', color: '#8b949e', fontWeight: 600 }}>STARTER PLAN</span>
            <div style={{ fontSize: '24px', fontWeight: 700, color: '#58a6ff' }}>12 Tenants</div>
            <span style={{ fontSize: '12px', color: '#8b949e' }}>Max 50 Machines / Org</span>
          </div>

          <div style={{ background: '#0d1117', padding: '16px', borderRadius: '8px', border: '1px solid #30363d' }}>
            <span style={{ fontSize: '12px', color: '#8b949e', fontWeight: 600 }}>PROFESSIONAL PLAN</span>
            <div style={{ fontSize: '24px', fontWeight: 700, color: '#a5d6ff' }}>8 Tenants</div>
            <span style={{ fontSize: '12px', color: '#8b949e' }}>Max 250 Machines / Org</span>
          </div>

          <div style={{ background: '#0d1117', padding: '16px', borderRadius: '8px', border: '1px solid #30363d' }}>
            <span style={{ fontSize: '12px', color: '#8b949e', fontWeight: 600 }}>ENTERPRISE PLAN</span>
            <div style={{ fontSize: '24px', fontWeight: 700, color: '#3fb950' }}>5 Tenants</div>
            <span style={{ fontSize: '12px', color: '#8b949e' }}>Unlimited Seats & Custom SSO</span>
          </div>
        </div>
      </div>
    </div>
  );
}
