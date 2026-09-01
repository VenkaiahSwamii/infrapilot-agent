import React from 'react';
import { Server, Activity, ShieldCheck, AlertTriangle } from 'lucide-react';

export default function ServiceHealth({ services }) {
  const list = services || [
    { service_name: 'api-gateway', avg_latency_ms: 42.5, p95_latency_ms: 110.0, request_count: 14200, throughput: 325.0, error_rate: 0.08 },
    { service_name: 'auth-service', avg_latency_ms: 12.1, p95_latency_ms: 28.0, request_count: 18400, throughput: 410.0, error_rate: 0.02 },
    { service_name: 'postgres-db', avg_latency_ms: 78.4, p95_latency_ms: 240.0, request_count: 22100, throughput: 520.0, error_rate: 0.08 },
    { service_name: 'ai-engine', avg_latency_ms: 340.2, p95_latency_ms: 890.0, request_count: 3800, throughput: 85.0, error_rate: 0.39 },
  ];

  return (
    <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '12px', padding: '20px' }}>
      <h4 style={{ margin: '0 0 16px 0', color: '#f0f6fc', display: 'flex', alignItems: 'center', gap: '8px', fontSize: '15px' }}>
        <Server size={18} color="#3fb950" /> Microservice Performance &amp; Health Matrix
      </h4>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '14px' }}>
        {list.map((s, idx) => (
          <div key={idx} style={{ background: '#0d1117', border: '1px solid #30363d', borderRadius: '8px', padding: '14px' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '8px' }}>
              <span style={{ fontWeight: 700, color: '#f0f6fc', fontSize: '14px' }}>{s.service_name}</span>
              <ShieldCheck size={16} color="#3fb950" />
            </div>

            <div style={{ display: 'flex', flexDirection: 'column', gap: '4px', fontSize: '12px', color: '#8b949e' }}>
              <div>Avg Latency: <strong style={{ color: '#58a6ff' }}>{s.avg_latency_ms} ms</strong></div>
              <div>P95 Latency: <strong style={{ color: '#ffa657' }}>{s.p95_latency_ms} ms</strong></div>
              <div>Throughput: <strong style={{ color: '#a855f7' }}>{s.throughput} req/min</strong></div>
              <div>Error Rate: <strong style={{ color: s.error_rate > 1 ? '#f78166' : '#3fb950' }}>{s.error_rate}%</strong></div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
