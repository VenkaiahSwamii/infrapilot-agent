import React from 'react';
import { Server, Cpu, Database, Zap, HardDrive, Globe } from 'lucide-react';

export default function AWS({ resources }) {
  const list = resources?.resources || [
    { name: 'prod-web-ec2-01', resource_type: 'ec2', region: 'us-east-1', state: 'running' },
    { name: 'prod-api-ec2-02', resource_type: 'ec2', region: 'us-east-1', state: 'running' },
    { name: 'prod-eks-cluster', resource_type: 'eks', region: 'us-east-1', state: 'ACTIVE' },
    { name: 'prod-rds-postgres', resource_type: 'rds', region: 'us-west-2', state: 'available' },
    { name: 'image-processor', resource_type: 'lambda', region: 'us-east-1', state: 'active' },
    { name: 's3-prod-assets-bucket-acme', resource_type: 's3', region: 'us-east-1', state: 'active' },
  ];

  const counts = {
    regions: 3,
    ec2: 24,
    eks: 2,
    rds: 4,
    lambda: 126,
    s3: 18,
  };

  return (
    <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '12px', padding: '24px' }}>
      <div style={{ marginBottom: '20px' }}>
        <h3 style={{ margin: 0, color: '#ffa657', fontSize: '18px', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '8px' }}>
          <CloudIcon />
          Amazon Web Services (AWS) Observability
        </h3>
        <span style={{ fontSize: '13px', color: '#8b949e' }}>EC2, EBS, EKS, RDS, Lambda, S3, and CloudWatch metrics</span>
      </div>

      {/* AWS Stat Bar */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(150px, 1fr))', gap: '16px', marginBottom: '24px' }}>
        <div style={{ background: '#0d1117', border: '1px solid #30363d', borderRadius: '8px', padding: '16px' }}>
          <div style={{ color: '#8b949e', fontSize: '11px', fontWeight: 700 }}>REGIONS</div>
          <div style={{ fontSize: '24px', fontWeight: 800, color: '#f0f6fc', marginTop: '4px' }}>{counts.regions}</div>
        </div>
        <div style={{ background: '#0d1117', border: '1px solid #30363d', borderRadius: '8px', padding: '16px' }}>
          <div style={{ color: '#8b949e', fontSize: '11px', fontWeight: 700 }}>EC2 INSTANCES</div>
          <div style={{ fontSize: '24px', fontWeight: 800, color: '#58a6ff', marginTop: '4px' }}>{counts.ec2}</div>
        </div>
        <div style={{ background: '#0d1117', border: '1px solid #30363d', borderRadius: '8px', padding: '16px' }}>
          <div style={{ color: '#8b949e', fontSize: '11px', fontWeight: 700 }}>EKS CLUSTERS</div>
          <div style={{ fontSize: '24px', fontWeight: 800, color: '#a855f7', marginTop: '4px' }}>{counts.eks}</div>
        </div>
        <div style={{ background: '#0d1117', border: '1px solid #30363d', borderRadius: '8px', padding: '16px' }}>
          <div style={{ color: '#8b949e', fontSize: '11px', fontWeight: 700 }}>RDS DATABASES</div>
          <div style={{ fontSize: '24px', fontWeight: 800, color: '#3fb950', marginTop: '4px' }}>{counts.rds}</div>
        </div>
        <div style={{ background: '#0d1117', border: '1px solid #30363d', borderRadius: '8px', padding: '16px' }}>
          <div style={{ color: '#8b949e', fontSize: '11px', fontWeight: 700 }}>LAMBDA FUNCTIONS</div>
          <div style={{ fontSize: '24px', fontWeight: 800, color: '#ffa657', marginTop: '4px' }}>{counts.lambda}</div>
        </div>
      </div>

      {/* Resource Table */}
      <div style={{ overflowX: 'auto' }}>
        <table style={{ width: '100%', borderCollapse: 'collapse', textAlign: 'left', fontSize: '13px' }}>
          <thead>
            <tr style={{ borderBottom: '1px solid #30363d', color: '#8b949e' }}>
              <th style={{ padding: '12px' }}>RESOURCE NAME</th>
              <th style={{ padding: '12px' }}>SERVICE TYPE</th>
              <th style={{ padding: '12px' }}>REGION</th>
              <th style={{ padding: '12px' }}>STATUS</th>
            </tr>
          </thead>
          <tbody>
            {list.map((r, idx) => (
              <tr key={idx} style={{ borderBottom: '1px solid #21262d', color: '#c9d1d9' }}>
                <td style={{ padding: '12px', fontWeight: 600, color: '#f0f6fc' }}>{r.name}</td>
                <td style={{ padding: '12px', textTransform: 'uppercase' }}>
                  <span style={{ padding: '2px 6px', borderRadius: '4px', background: '#21262d', color: '#ffa657', fontSize: '11px', fontWeight: 700 }}>
                    {r.resource_type}
                  </span>
                </td>
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

function CloudIcon() {
  return (
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M17.5 19H9a7 7 0 1 1 6.71-9h1.79a4.5 4.5 0 1 1 0 9Z"/>
    </svg>
  );
}
