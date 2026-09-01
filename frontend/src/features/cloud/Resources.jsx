import React from 'react';
import { Server, Layers, Cpu, Cloud, Globe } from 'lucide-react';

export default function Resources({ summary }) {
  const hybrid = summary?.hybrid_summary || {
    total_servers: 428,
    total_vms: 376,
    docker_containers: 3200,
    k8s_clusters: 15,
    aws_resource_count: 156,
    azure_resource_count: 84,
    gcp_resource_count: 48,
  };

  const categories = [
    { title: 'Physical Servers (On-Premises)', count: 52, icon: Server, color: '#58a6ff' },
    { title: 'Cloud Virtual Machines', count: hybrid.total_vms, icon: Cloud, color: '#ffa657' },
    { title: 'Docker Containers', count: hybrid.docker_containers, icon: Layers, color: '#a855f7' },
    { title: 'Kubernetes Clusters (EKS/AKS/GKE/On-Prem)', count: hybrid.k8s_clusters, icon: Cpu, color: '#3fb950' },
    { title: 'AWS Cloud Resources', count: hybrid.aws_resource_count, icon: Globe, color: '#ffa657' },
    { title: 'Azure Cloud Resources', count: hybrid.azure_resource_count, icon: Globe, color: '#58a6ff' },
    { title: 'GCP Cloud Resources', count: hybrid.gcp_resource_count, icon: Globe, color: '#ea4335' },
  ];

  return (
    <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '12px', padding: '24px' }}>
      <div style={{ marginBottom: '20px' }}>
        <h3 style={{ margin: 0, color: '#f0f6fc', fontSize: '18px', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '8px' }}>
          <Layers size={20} color="#58a6ff" />
          Unified Hybrid Cloud Infrastructure Inventory
        </h3>
        <span style={{ fontSize: '13px', color: '#8b949e' }}>Single-pane-of-glass discovery for physical servers, VMs, containers, K8s, and multi-cloud services</span>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(260px, 1fr))', gap: '20px' }}>
        {categories.map((c, idx) => {
          const Icon = c.icon;
          return (
            <div key={idx} style={{ background: '#0d1117', border: '1px solid #30363d', borderRadius: '10px', padding: '20px' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', color: '#8b949e', fontSize: '12px', fontWeight: 600 }}>
                <span>{c.title.toUpperCase()}</span>
                <Icon size={18} color={c.color} />
              </div>
              <div style={{ fontSize: '32px', fontWeight: 800, color: '#f0f6fc', margin: '12px 0 4px 0' }}>
                {c.count}
              </div>
              <div style={{ fontSize: '12px', color: '#3fb950' }}>● Active & Monitored</div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
