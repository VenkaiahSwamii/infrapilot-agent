import React from 'react';
import { ShieldAlert, AlertTriangle, Lock, ShieldCheck } from 'lucide-react';

export default function Security({ findings }) {
  const list = findings?.security_findings || [
    {
      provider: 'aws',
      finding_type: 'public_s3',
      severity: 'CRITICAL',
      title: 'Public S3 Bucket Detected',
      description: "Bucket 's3-acme-public-logs' has global read ACL permissions enabled.",
      remediation: 'Enable S3 Block Public Access at bucket policy level.',
    },
    {
      provider: 'aws',
      finding_type: 'open_security_group',
      severity: 'HIGH',
      title: 'Open SSH Port (22) to Internet',
      description: 'Security group ingress rule allows 0.0.0.0/0 on port 22.',
      remediation: 'Restrict ingress to corporate VPN IP cidr.',
    },
    {
      provider: 'azure',
      finding_type: 'unencrypted_disk',
      severity: 'MEDIUM',
      title: 'Unencrypted Managed Disk',
      description: 'Azure VM OS disk does not have customer-managed encryption (CMEK) enabled.',
      remediation: 'Enable Disk Encryption Set using Azure Key Vault.',
    },
  ];

  return (
    <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '12px', padding: '24px' }}>
      <div style={{ marginBottom: '20px' }}>
        <h3 style={{ margin: 0, color: '#f78166', fontSize: '18px', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '8px' }}>
          <ShieldAlert size={20} color="#f78166" />
          Cloud Security Posture Management (CSPM)
        </h3>
        <span style={{ fontSize: '13px', color: '#8b949e' }}>Detect public storage buckets, open security groups, unencrypted disks, and weak IAM policies</span>
      </div>

      <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
        {list.map((item, idx) => (
          <div key={idx} style={{
            background: '#0d1117',
            border: `1px solid ${item.severity === 'CRITICAL' ? '#f78166' : item.severity === 'HIGH' ? '#ffa657' : '#30363d'}`,
            borderRadius: '10px',
            padding: '20px'
          }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '8px' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                <span style={{
                  padding: '3px 8px',
                  borderRadius: '4px',
                  fontSize: '11px',
                  fontWeight: 700,
                  textTransform: 'uppercase',
                  backgroundColor: item.provider === 'aws' ? 'rgba(255, 166, 87, 0.15)' : 'rgba(88, 166, 255, 0.15)',
                  color: item.provider === 'aws' ? '#ffa657' : '#58a6ff',
                }}>
                  {item.provider}
                </span>
                <h4 style={{ margin: 0, color: '#f0f6fc', fontSize: '15px' }}>{item.title}</h4>
              </div>

              <span style={{
                padding: '4px 10px',
                borderRadius: '9999px',
                fontSize: '11px',
                fontWeight: 700,
                backgroundColor: item.severity === 'CRITICAL' ? 'rgba(247, 129, 102, 0.15)' : 'rgba(255, 166, 87, 0.15)',
                color: item.severity === 'CRITICAL' ? '#f78166' : '#ffa657',
              }}>
                {item.severity}
              </span>
            </div>

            <p style={{ color: '#c9d1d9', fontSize: '13px', margin: '8px 0 12px 0' }}>{item.description}</p>

            <div style={{ background: '#161b22', padding: '10px 14px', borderRadius: '6px', fontSize: '12px', color: '#58a6ff' }}>
              <strong>RECOMMENDED REMEDIATION:</strong> {item.remediation}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
