import React from 'react';
import { ShieldCheck, Check, X } from 'lucide-react';

export default function Roles() {
  const roles = [
    { name: 'Organization Owner', slug: 'owner', desc: 'Full ownership and billing authority over tenant organization.' },
    { name: 'Organization Admin', slug: 'admin', desc: 'Full access to manage users, settings, and infrastructure.' },
    { name: 'DevOps Engineer', slug: 'devops', desc: 'Can configure servers, deployment pipelines, AI, and K8s.' },
    { name: 'Operator', slug: 'operator', desc: 'Can restart services, manage docker, and handle incidents.' },
    { name: 'Viewer', slug: 'viewer', desc: 'Read-only access to tenant dashboards, logs, and alerts.' },
  ];

  const permissions = [
    { key: 'dashboard', label: 'View Dashboard', owner: true, admin: true, devops: true, operator: true, viewer: true },
    { key: 'servers', label: 'Manage Servers & Hosts', owner: true, admin: true, devops: true, operator: true, viewer: false },
    { key: 'deploy', label: 'Deploy Agents & Pipelines', owner: true, admin: true, devops: true, operator: false, viewer: false },
    { key: 'logs', label: 'View & Stream Logs', owner: true, admin: true, devops: true, operator: true, viewer: true },
    { key: 'alerts', label: 'Manage & Resolve Alerts', owner: true, admin: true, devops: true, operator: true, viewer: false },
    { key: 'users', label: 'User & Role Management', owner: true, admin: true, devops: false, operator: false, viewer: false },
    { key: 'ai', label: 'Configure AI & Incident Ops', owner: true, admin: true, devops: true, operator: false, viewer: false },
    { key: 'k8s', label: 'Configure Kubernetes Clusters', owner: true, admin: true, devops: true, operator: false, viewer: false },
  ];

  return (
    <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '12px', padding: '24px' }}>
      <div style={{ marginBottom: '20px' }}>
        <h3 style={{ margin: 0, color: '#f0f6fc', fontSize: '18px', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '8px' }}>
          <ShieldCheck size={20} color="#a855f7" />
          Organization Permission Matrix
        </h3>
        <span style={{ fontSize: '13px', color: '#8b949e' }}>Role-Based Access Control (RBAC) permission scoping per tenant</span>
      </div>

      <div style={{ overflowX: 'auto' }}>
        <table style={{ width: '100%', borderCollapse: 'collapse', textAlign: 'left', fontSize: '13px' }}>
          <thead>
            <tr style={{ borderBottom: '1px solid #30363d', color: '#8b949e' }}>
              <th style={{ padding: '12px' }}>PERMISSION</th>
              <th style={{ padding: '12px', textAlign: 'center' }}>OWNER</th>
              <th style={{ padding: '12px', textAlign: 'center' }}>ADMIN</th>
              <th style={{ padding: '12px', textAlign: 'center' }}>DEVOPS</th>
              <th style={{ padding: '12px', textAlign: 'center' }}>OPERATOR</th>
              <th style={{ padding: '12px', textAlign: 'center' }}>VIEWER</th>
            </tr>
          </thead>
          <tbody>
            {permissions.map((p) => (
              <tr key={p.key} style={{ borderBottom: '1px solid #21262d', color: '#c9d1d9' }}>
                <td style={{ padding: '12px', fontWeight: 600, color: '#f0f6fc' }}>{p.label}</td>
                <td style={{ padding: '12px', textAlign: 'center' }}>{p.owner ? <Check size={18} color="#3fb950" /> : <X size={18} color="#f78166" />}</td>
                <td style={{ padding: '12px', textAlign: 'center' }}>{p.admin ? <Check size={18} color="#3fb950" /> : <X size={18} color="#f78166" />}</td>
                <td style={{ padding: '12px', textAlign: 'center' }}>{p.devops ? <Check size={18} color="#3fb950" /> : <X size={18} color="#f78166" />}</td>
                <td style={{ padding: '12px', textAlign: 'center' }}>{p.operator ? <Check size={18} color="#3fb950" /> : <X size={18} color="#f78166" />}</td>
                <td style={{ padding: '12px', textAlign: 'center' }}>{p.viewer ? <Check size={18} color="#3fb950" /> : <X size={18} color="#f78166" />}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
