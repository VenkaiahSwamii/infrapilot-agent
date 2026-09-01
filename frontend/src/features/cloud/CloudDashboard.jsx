import React from 'react';
import { Cloud, Server, Cpu, DollarSign, AlertTriangle, ShieldCheck, Layers } from 'lucide-react';

export default function CloudDashboard({ accounts, summary }) {
  const accountList = accounts?.accounts || [];
  const hybrid = summary || {
    total_servers: 428,
    total_vms: 376,
    docker_containers: 3200,
    k8s_clusters: 15,
  };

  const awsCount = accountList.filter(a => a.provider === 'aws').length || 3;
  const azCount = accountList.filter(a => a.provider === 'azure').length || 2;
  const gcpCount = accountList.filter(a => a.provider === 'gcp').length || 1;

  const totalAccounts = accountList.length || 6;

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
      {/* Summary KPI Cards Grid */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(220px, 1fr))', gap: '20px' }}>
        <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '10px', padding: '20px' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', color: '#8b949e', fontSize: '12px', fontWeight: 600 }}>
            <span>CONNECTED CLOUD ACCOUNTS</span>
            <Cloud size={18} color="#58a6ff" />
          </div>
          <div style={{ fontSize: '28px', fontWeight: 800, color: '#f0f6fc', margin: '10px 0 4px 0' }}>
            {totalAccounts}
          </div>
          <div style={{ fontSize: '12px', color: '#8b949e' }}>
            AWS: <strong style={{ color: '#ffa657' }}>{awsCount}</strong> | Azure: <strong style={{ color: '#58a6ff' }}>{azCount}</strong> | GCP: <strong style={{ color: '#ea4335' }}>{gcpCount}</strong>
          </div>
        </div>

        <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '10px', padding: '20px' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', color: '#8b949e', fontSize: '12px', fontWeight: 600 }}>
            <span>TOTAL VIRTUAL MACHINES</span>
            <Server size={18} color="#a855f7" />
          </div>
          <div style={{ fontSize: '28px', fontWeight: 800, color: '#f0f6fc', margin: '10px 0 4px 0' }}>
            {hybrid.total_vms}
          </div>
          <div style={{ fontSize: '12px', color: '#3fb950' }}>{hybrid.total_servers} Total Hybrid Hosts</div>
        </div>

        <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '10px', padding: '20px' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', color: '#8b949e', fontSize: '12px', fontWeight: 600 }}>
            <span>CLOUD CLUSTERS</span>
            <Cpu size={18} color="#ffa657" />
          </div>
          <div style={{ fontSize: '28px', fontWeight: 800, color: '#f0f6fc', margin: '10px 0 4px 0' }}>
            {hybrid.k8s_clusters}
          </div>
          <div style={{ fontSize: '12px', color: '#8b949e' }}>Managed EKS, AKS, & GKE</div>
        </div>

        <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '10px', padding: '20px' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', color: '#8b949e', fontSize: '12px', fontWeight: 600 }}>
            <span>ESTIMATED MONTHLY COST</span>
            <DollarSign size={18} color="#3fb950" />
          </div>
          <div style={{ fontSize: '28px', fontWeight: 800, color: '#3fb950', margin: '10px 0 4px 0' }}>
            $7,430
          </div>
          <div style={{ fontSize: '12px', color: '#8b949e' }}>Cross-cloud monthly spend</div>
        </div>
      </div>

      {/* Account Status Grid */}
      <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '12px', padding: '24px' }}>
        <h3 style={{ margin: '0 0 16px 0', color: '#f0f6fc', fontSize: '18px', fontWeight: 600 }}>
          Multi-Cloud Account Overview
        </h3>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))', gap: '16px' }}>
          {accountList.map((acc) => (
            <div key={acc.id} style={{ background: '#0d1117', border: '1px solid #30363d', borderRadius: '8px', padding: '16px' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '8px' }}>
                <span style={{
                  padding: '3px 8px',
                  borderRadius: '4px',
                  fontSize: '11px',
                  fontWeight: 700,
                  textTransform: 'uppercase',
                  backgroundColor: acc.provider === 'aws' ? 'rgba(255, 166, 87, 0.15)' : acc.provider === 'azure' ? 'rgba(88, 166, 255, 0.15)' : 'rgba(234, 67, 53, 0.15)',
                  color: acc.provider === 'aws' ? '#ffa657' : acc.provider === 'azure' ? '#58a6ff' : '#ea4335',
                }}>
                  {acc.provider}
                </span>
                <span style={{ fontSize: '12px', color: '#3fb950', fontWeight: 600 }}>● {acc.status}</span>
              </div>

              <div style={{ fontWeight: 700, color: '#f0f6fc', fontSize: '15px', marginBottom: '4px' }}>{acc.account_name}</div>
              <div style={{ fontSize: '12px', color: '#8b949e', marginBottom: '12px' }}>ID: <code>{acc.account_id}</code></div>

              <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '12px', color: '#c9d1d9', borderTop: '1px solid #21262d', paddingTop: '8px' }}>
                <span>VMs: <strong>{acc.total_vms}</strong></span>
                <span>Clusters: <strong>{acc.total_clusters}</strong></span>
                <span>Monthly: <strong style={{ color: '#3fb950' }}>${acc.monthly_cost}</strong></span>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
