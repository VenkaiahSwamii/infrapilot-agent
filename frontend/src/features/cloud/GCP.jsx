import React from 'react';
import { Server, Cpu, Database, Cloud } from 'lucide-react';

export default function GCP({ resources }) {
  const list = resources?.resources || [
    { name: 'gcp-instance-worker-01', resource_type: 'compute', project_id: 'gcp-prod-infra', zone: 'us-central1-a', state: 'RUNNING' },
    { name: 'gcp-gke-prod', resource_type: 'gke', project_id: 'gcp-prod-infra', zone: 'us-central1', state: 'RUNNING' },
    { name: 'cloudsql-replica-01', resource_type: 'sql', project_id: 'gcp-prod-infra', zone: 'us-central1-b', state: 'RUNNABLE' },
  ];

  return (
    <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '12px', padding: '24px' }}>
      <div style={{ marginBottom: '20px' }}>
        <h3 style={{ margin: 0, color: '#ea4335', fontSize: '18px', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '8px' }}>
          <Cloud size={20} color="#ea4335" />
          Google Cloud Platform (GCP) Observability
        </h3>
        <span style={{ fontSize: '13px', color: '#8b949e' }}>Compute Engine, GKE, Cloud SQL, Cloud Storage, Cloud Functions</span>
      </div>

      <div style={{ overflowX: 'auto' }}>
        <table style={{ width: '100%', borderCollapse: 'collapse', textAlign: 'left', fontSize: '13px' }}>
          <thead>
            <tr style={{ borderBottom: '1px solid #30363d', color: '#8b949e' }}>
              <th style={{ padding: '12px' }}>RESOURCE NAME</th>
              <th style={{ padding: '12px' }}>TYPE</th>
              <th style={{ padding: '12px' }}>PROJECT ID</th>
              <th style={{ padding: '12px' }}>ZONE</th>
              <th style={{ padding: '12px' }}>STATUS</th>
            </tr>
          </thead>
          <tbody>
            {list.map((r, idx) => (
              <tr key={idx} style={{ borderBottom: '1px solid #21262d', color: '#c9d1d9' }}>
                <td style={{ padding: '12px', fontWeight: 600, color: '#f0f6fc' }}>{r.name}</td>
                <td style={{ padding: '12px', textTransform: 'uppercase' }}>
                  <span style={{ padding: '2px 6px', borderRadius: '4px', background: '#21262d', color: '#ea4335', fontSize: '11px', fontWeight: 700 }}>
                    {r.resource_type}
                  </span>
                </td>
                <td style={{ padding: '12px', color: '#8b949e' }}>{r.project_id}</td>
                <td style={{ padding: '12px', color: '#58a6ff' }}>{r.zone}</td>
                <td style={{ padding: '12px', color: '#3fb950', fontWeight: 600 }}>{r.state}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
