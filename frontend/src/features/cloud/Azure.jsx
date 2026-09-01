import React from 'react';
import { Server, Cpu, Database, Layers, Cloud } from 'lucide-react';

export default function Azure({ resources }) {
  const list = resources?.resources || [
    { name: 'az-vm-web-01', resource_type: 'vm', resource_group: 'rg-prod', region: 'eastus', state: 'running' },
    { name: 'az-aks-cluster', resource_type: 'aks', resource_group: 'rg-prod', region: 'eastus', state: 'Succeeded' },
    { name: 'az-sql-analytics', resource_type: 'sql', resource_group: 'rg-prod', region: 'westeurope', state: 'Online' },
  ];

  return (
    <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '12px', padding: '24px' }}>
      <div style={{ marginBottom: '20px' }}>
        <h3 style={{ margin: 0, color: '#58a6ff', fontSize: '18px', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '8px' }}>
          <Cloud size={20} color="#58a6ff" />
          Microsoft Azure Observability
        </h3>
        <span style={{ fontSize: '13px', color: '#8b949e' }}>Azure Virtual Machines, AKS, Azure SQL, Storage Accounts, App Services</span>
      </div>

      <div style={{ overflowX: 'auto' }}>
        <table style={{ width: '100%', borderCollapse: 'collapse', textAlign: 'left', fontSize: '13px' }}>
          <thead>
            <tr style={{ borderBottom: '1px solid #30363d', color: '#8b949e' }}>
              <th style={{ padding: '12px' }}>RESOURCE NAME</th>
              <th style={{ padding: '12px' }}>TYPE</th>
              <th style={{ padding: '12px' }}>RESOURCE GROUP</th>
              <th style={{ padding: '12px' }}>REGION</th>
              <th style={{ padding: '12px' }}>STATUS</th>
            </tr>
          </thead>
          <tbody>
            {list.map((r, idx) => (
              <tr key={idx} style={{ borderBottom: '1px solid #21262d', color: '#c9d1d9' }}>
                <td style={{ padding: '12px', fontWeight: 600, color: '#f0f6fc' }}>{r.name}</td>
                <td style={{ padding: '12px', textTransform: 'uppercase' }}>
                  <span style={{ padding: '2px 6px', borderRadius: '4px', background: '#21262d', color: '#58a6ff', fontSize: '11px', fontWeight: 700 }}>
                    {r.resource_type}
                  </span>
                </td>
                <td style={{ padding: '12px', color: '#8b949e' }}>{r.resource_group}</td>
                <td style={{ padding: '12px', color: '#58a6ff' }}>{r.region}</td>
                <td style={{ padding: '12px', color: '#3fb950', fontWeight: 600 }}>{r.state}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
