import React from 'react';
import { DollarSign, PieChart, TrendingUp } from 'lucide-react';

export default function CostAnalysis({ costs }) {
  const data = costs || {
    total_monthly: 7430,
    aws_monthly: 842,
    azure_monthly: 511,
    gcp_monthly: 218,
    breakdown: [
      { provider: 'aws', category: 'compute', monthly_spend: 420 },
      { provider: 'aws', category: 'database', monthly_spend: 210 },
      { provider: 'aws', category: 'storage', monthly_spend: 112 },
      { provider: 'aws', category: 'network', monthly_spend: 100 },
      { provider: 'azure', category: 'compute', monthly_spend: 280 },
      { provider: 'azure', category: 'database', monthly_spend: 140 },
      { provider: 'azure', category: 'storage', monthly_spend: 91 },
      { provider: 'gcp', category: 'compute', monthly_spend: 120 },
      { provider: 'gcp', category: 'database', monthly_spend: 68 },
      { provider: 'gcp', category: 'storage', monthly_spend: 30 },
    ],
  };

  const providers = [
    { name: 'Amazon Web Services', code: 'aws', amount: data.aws_monthly || 842, color: '#ffa657' },
    { name: 'Microsoft Azure', code: 'azure', amount: data.azure_monthly || 511, color: '#58a6ff' },
    { name: 'Google Cloud Platform', code: 'gcp', amount: data.gcp_monthly || 218, color: '#ea4335' },
  ];

  return (
    <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '12px', padding: '24px' }}>
      <div style={{ marginBottom: '20px' }}>
        <h3 style={{ margin: 0, color: '#3fb950', fontSize: '18px', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '8px' }}>
          <DollarSign size={20} color="#3fb950" />
          Cloud Cost Management & FinOps Dashboard
        </h3>
        <span style={{ fontSize: '13px', color: '#8b949e' }}>Multi-cloud spend tracking for Compute, Storage, Network, and Databases</span>
      </div>

      {/* Provider Cost Cards */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(220px, 1fr))', gap: '20px', marginBottom: '24px' }}>
        {providers.map((p) => (
          <div key={p.code} style={{ background: '#0d1117', border: '1px solid #30363d', borderRadius: '10px', padding: '20px' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', color: '#8b949e', fontSize: '11px', fontWeight: 700 }}>
              <span>{p.name.toUpperCase()}</span>
              <DollarSign size={16} color={p.color} />
            </div>
            <div style={{ fontSize: '28px', fontWeight: 800, color: p.color, margin: '10px 0 4px 0' }}>
              ${p.amount}
            </div>
            <div style={{ fontSize: '12px', color: '#8b949e' }}>Estimated monthly spend</div>
          </div>
        ))}
      </div>

      {/* Resource Category Breakdown */}
      <h4 style={{ color: '#f0f6fc', margin: '0 0 12px 0', fontSize: '15px' }}>Category Spend Breakdown</h4>
      <div style={{ overflowX: 'auto' }}>
        <table style={{ width: '100%', borderCollapse: 'collapse', textAlign: 'left', fontSize: '13px' }}>
          <thead>
            <tr style={{ borderBottom: '1px solid #30363d', color: '#8b949e' }}>
              <th style={{ padding: '12px' }}>PROVIDER</th>
              <th style={{ padding: '12px' }}>RESOURCE CATEGORY</th>
              <th style={{ padding: '12px' }}>MONTHLY SPEND</th>
            </tr>
          </thead>
          <tbody>
            {(data.breakdown || []).map((b, idx) => (
              <tr key={idx} style={{ borderBottom: '1px solid #21262d', color: '#c9d1d9' }}>
                <td style={{ padding: '12px', textTransform: 'uppercase', fontWeight: 700, color: b.provider === 'aws' ? '#ffa657' : b.provider === 'azure' ? '#58a6ff' : '#ea4335' }}>
                  {b.provider}
                </td>
                <td style={{ padding: '12px', textTransform: 'capitalize' }}>{b.category}</td>
                <td style={{ padding: '12px', fontWeight: 700, color: '#3fb950' }}>${b.monthly_spend}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
