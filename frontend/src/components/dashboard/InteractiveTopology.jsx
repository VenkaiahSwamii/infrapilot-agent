import React, { useState } from 'react';
import {
  ShieldCheck,
  Server,
  Monitor,
  Layers,
  Network,
  Cloud,
  HardDrive,
  Database,
} from 'lucide-react';

const ICON_MAP = {
  linux: Server,
  windows: Monitor,
  docker: Layers,
  kubernetes: Network,
  cloud: Cloud,
  storage: HardDrive,
  database: Database,
};

const COLOR_MAP = {
  linux: '#22c55e',
  windows: '#38bdf8',
  docker: '#06b6d4',
  kubernetes: '#a855f7',
  cloud: '#3b82f6',
  storage: '#eab308',
  database: '#ef4444',
};

export default function InteractiveTopology({ categories = [], selectedCategory = 'all', onSelectCategory }) {
  const [hoveredNode, setHoveredNode] = useState(null);

  const activeCategories = categories.filter((c) => c.id !== 'all').slice(0, 6);
  const displayNodes = activeCategories.length > 0 ? activeCategories : [
    { id: 'linux', label: 'Linux Servers', count: 0, online: 0 },
    { id: 'windows', label: 'Windows Fleet', count: 0, online: 0 },
    { id: 'docker', label: 'Docker Hosts', count: 0, online: 0 },
    { id: 'kubernetes', label: 'Kubernetes Pods', count: 0, online: 0 },
    { id: 'cloud', label: 'Cloud VPCs', count: 0, online: 0 },
    { id: 'storage', label: 'Storage Volumes', count: 0, online: 0 },
  ];

  return (
    <section className="topology-interactive-card">
      <div className="topology-card-header">
        <div>
          <span className="eyebrow-label">INFRASTRUCTURE ARCHITECTURE</span>
          <h3>Cluster Topology Map</h3>
        </div>
        <div className="topology-legend">
          <span className="legend-dot pulse" />
          <span>Real-time Agent Mesh</span>
        </div>
      </div>

      <div className="topology-canvas-container">
        {/* Animated Radial Background Grid */}
        <div className="topology-radial-bg" />

        {/* Center Control Plane Hub */}
        <div className="topology-center-hub">
          <div className="hub-pulse-ring" />
          <div className="hub-core">
            <ShieldCheck size={26} color="#38bdf8" />
            <span>CONTROL PLANE</span>
            <small>Active Telemetry</small>
          </div>
        </div>

        {/* Orbiting Infrastructure Nodes */}
        <div className="topology-nodes-orbit">
          {displayNodes.map((cat, idx) => {
            const Icon = ICON_MAP[cat.id] || Server;
            const color = COLOR_MAP[cat.id] || '#38bdf8';
            const isSelected = selectedCategory === cat.id;
            const isHovered = hoveredNode === cat.id;

            return (
              <div
                key={cat.id}
                className={`orbit-node-item orbit-node-${idx} ${isSelected ? 'selected' : ''} ${isHovered ? 'hovered' : ''}`}
                onClick={() => onSelectCategory(isSelected ? 'all' : cat.id)}
                onMouseEnter={() => setHoveredNode(cat.id)}
                onMouseLeave={() => setHoveredNode(null)}
                style={{ '--node-color': color }}
              >
                {/* SVG Conduit Beam */}
                <div className="conduit-line" />

                <div className="node-badge-card">
                  <div className="node-icon-box" style={{ backgroundColor: `${color}18`, borderColor: `${color}40` }}>
                    <Icon size={16} color={color} />
                  </div>
                  <div className="node-info">
                    <strong>{cat.label}</strong>
                    <span>{cat.count || 0} Instances ({cat.online || 0} active)</span>
                  </div>
                  <span className="node-health-dot" style={{ backgroundColor: cat.count === 0 ? '#64748b' : color }} />
                </div>
              </div>
            );
          })}
        </div>
      </div>

      <div className="topology-footer-bar">
        <span>Click any cluster node to filter fleet inventory.</span>
        <button
          className="btn-reset-filter"
          onClick={() => onSelectCategory('all')}
          style={{ opacity: selectedCategory === 'all' ? 0.4 : 1 }}
          type="button"
        >
          Reset Filter ({selectedCategory})
        </button>
      </div>

      <style>{`
        .topology-interactive-card {
          background: #0d1322;
          border: 1px solid #1e2d45;
          border-radius: 12px;
          padding: 16px;
          display: flex;
          flex-direction: column;
          box-shadow: 0 4px 20px rgba(0, 0, 0, 0.35);
        }
        .topology-card-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
          margin-bottom: 12px;
        }
        .eyebrow-label {
          font-size: 10px;
          font-weight: 800;
          color: #38bdf8;
          letter-spacing: 0.1em;
          display: block;
        }
        .topology-card-header h3 {
          font-size: 14px;
          font-weight: 700;
          color: #ffffff;
          margin: 2px 0 0 0;
        }
        .topology-legend {
          display: flex;
          align-items: center;
          gap: 6px;
          font-size: 11px;
          color: #94a3b8;
        }
        .legend-dot {
          width: 6px;
          height: 6px;
          border-radius: 50%;
          background: #22c55e;
        }
        .legend-dot.pulse {
          box-shadow: 0 0 8px #22c55e;
          animation: pulseGreen 2s infinite;
        }
        .topology-canvas-container {
          position: relative;
          height: 250px;
          background: #070b14;
          border: 1px solid #162033;
          border-radius: 10px;
          overflow: hidden;
          display: flex;
          align-items: center;
          justify-content: center;
        }
        .topology-radial-bg {
          position: absolute;
          inset: 0;
          background: radial-gradient(circle at center, rgba(56, 189, 248, 0.08) 0%, transparent 70%);
          pointer-events: none;
        }
        .topology-center-hub {
          position: relative;
          z-index: 10;
          display: flex;
          align-items: center;
          justify-content: center;
        }
        .hub-pulse-ring {
          position: absolute;
          width: 100px;
          height: 100px;
          border-radius: 50%;
          border: 1px solid rgba(56, 189, 248, 0.25);
          animation: hubExpand 3s infinite ease-out;
        }
        .hub-core {
          width: 78px;
          height: 78px;
          border-radius: 50%;
          background: #0d1322;
          border: 2px solid #38bdf8;
          box-shadow: 0 0 20px rgba(56, 189, 248, 0.35);
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          text-align: center;
          padding: 4px;
        }
        .hub-core span {
          font-size: 8px;
          font-weight: 800;
          color: #ffffff;
          letter-spacing: 0.05em;
          margin-top: 2px;
        }
        .hub-core small {
          font-size: 7px;
          color: #38bdf8;
        }
        .topology-nodes-orbit {
          position: absolute;
          inset: 0;
          display: grid;
          grid-template-columns: 1fr 1fr;
          grid-template-rows: 1fr 1fr 1fr;
          padding: 10px 14px;
          gap: 8px;
          pointer-events: none;
        }
        .orbit-node-item {
          pointer-events: auto;
          display: flex;
          align-items: center;
          cursor: pointer;
          transition: all 0.2s ease;
        }
        .orbit-node-0, .orbit-node-2, .orbit-node-4 {
          justify-content: flex-start;
        }
        .orbit-node-1, .orbit-node-3, .orbit-node-5 {
          justify-content: flex-end;
        }
        .node-badge-card {
          display: flex;
          align-items: center;
          gap: 8px;
          background: #0d1322;
          border: 1px solid #1e2d45;
          border-radius: 8px;
          padding: 5px 10px;
          transition: all 0.15s ease;
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.4);
        }
        .orbit-node-item:hover .node-badge-card {
          border-color: var(--node-color);
          transform: translateY(-2px);
          background: #111a2e;
        }
        .orbit-node-item.selected .node-badge-card {
          border-color: var(--node-color);
          box-shadow: 0 0 14px var(--node-color);
          background: #142036;
        }
        .node-icon-box {
          width: 26px;
          height: 26px;
          border-radius: 6px;
          border: 1px solid transparent;
          display: flex;
          align-items: center;
          justify-content: center;
        }
        .node-info {
          display: flex;
          flex-direction: column;
          line-height: 1.1;
        }
        .node-info strong {
          font-size: 11px;
          color: #f1f5f9;
          white-space: nowrap;
        }
        .node-info span {
          font-size: 9.5px;
          color: #94a3b8;
          white-space: nowrap;
        }
        .node-health-dot {
          width: 6px;
          height: 6px;
          border-radius: 50%;
          flex-shrink: 0;
          margin-left: 2px;
        }
        .topology-footer-bar {
          display: flex;
          align-items: center;
          justify-content: space-between;
          margin-top: 10px;
          font-size: 11px;
          color: #64748b;
        }
        .btn-reset-filter {
          background: #111a2e;
          border: 1px solid #1e2d45;
          color: #38bdf8;
          border-radius: 6px;
          padding: 3px 8px;
          font-size: 10.5px;
          font-weight: 600;
          cursor: pointer;
          transition: all 0.15s ease;
        }
        .btn-reset-filter:hover {
          background: #1e2d45;
          color: #ffffff;
        }
        @keyframes pulseGreen {
          0%, 100% { opacity: 1; transform: scale(1); }
          50% { opacity: 0.6; transform: scale(1.15); }
        }
        @keyframes hubExpand {
          0% { transform: scale(0.8); opacity: 0.8; }
          100% { transform: scale(1.6); opacity: 0; }
        }
      `}</style>
    </section>
  );
}
