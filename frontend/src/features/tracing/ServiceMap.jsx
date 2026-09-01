import React from 'react';
import { Network, Server, Database, Cpu, Layers, Box } from 'lucide-react';

export default function ServiceMap({ topology }) {
  const nodes = topology?.nodes || [
    { id: 'react-ui', label: 'React Frontend UI', type: 'frontend' },
    { id: 'api-gateway', label: 'Go API Gateway', type: 'gateway' },
    { id: 'auth-service', label: 'Auth Service', type: 'service' },
    { id: 'postgres-db', label: 'PostgreSQL Primary DB', type: 'database' },
    { id: 'redis-cache', label: 'Redis Cache & Queue', type: 'cache' },
    { id: 'ai-engine', label: 'AI RAG Engine', type: 'ai' },
    { id: 'docker-host', label: 'Docker Daemon', type: 'container' },
    { id: 'k8s-cluster', label: 'Kubernetes API', type: 'cluster' },
  ];

  const getNodeIcon = (type) => {
    switch (type) {
      case 'frontend': return <Layers size={20} color="#58a6ff" />;
      case 'gateway': return <Server size={20} color="#3fb950" />;
      case 'database': return <Database size={20} color="#ffa657" />;
      case 'cache': return <Box size={20} color="#f78166" />;
      case 'ai': return <Cpu size={20} color="#a855f7" />;
      default: return <Network size={20} color="#79c0ff" />;
    }
  };

  return (
    <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '12px', padding: '24px' }}>
      <h3 style={{ color: '#f0f6fc', margin: '0 0 16px 0', fontSize: '18px', display: 'flex', alignItems: 'center', gap: '8px' }}>
        <Network size={22} color="#58a6ff" /> Live Service Dependency Map Topology
      </h3>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(220px, 1fr))', gap: '16px' }}>
        {nodes.map((node) => (
          <div key={node.id} style={{
            background: '#0d1117',
            border: '1px solid #30363d',
            borderRadius: '10px',
            padding: '16px',
            display: 'flex',
            alignItems: 'center',
            gap: '12px',
            boxShadow: '0 2px 8px rgba(0,0,0,0.2)'
          }}>
            <div style={{ padding: '10px', background: '#161b22', borderRadius: '8px' }}>
              {getNodeIcon(node.type)}
            </div>
            <div>
              <div style={{ fontSize: '14px', fontWeight: 700, color: '#f0f6fc' }}>{node.label}</div>
              <div style={{ fontSize: '11px', color: '#8b949e', textTransform: 'uppercase', marginTop: '2px' }}>{node.type}</div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
