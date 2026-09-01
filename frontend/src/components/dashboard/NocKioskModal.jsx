import React, { useState, useEffect } from 'react';
import {
  X,
  Maximize,
  Minimize,
  ShieldCheck,
  Cpu,
  MemoryStick,
  HardDrive,
  Network,
  Radio,
  Clock,
} from 'lucide-react';

export default function NocKioskModal({ isOpen, onClose, overview, telemetryHistory, assets, alerts }) {
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [currentTime, setCurrentTime] = useState(() => new Date());
  const [activeMetric, setActiveMetric] = useState('cpu');

  useEffect(() => {
    if (!isOpen) return;

    const timer = setInterval(() => {
      setCurrentTime(new Date());
    }, 1000);

    return () => clearInterval(timer);
  }, [isOpen]);

  const toggleFullscreen = () => {
    if (!document.fullscreenElement) {
      document.documentElement.requestFullscreen().catch(() => {});
      setIsFullscreen(true);
    } else {
      document.exitFullscreen().catch(() => {});
      setIsFullscreen(false);
    }
  };

  if (!isOpen) return null;

  const onlineCount = assets.filter((a) => a.status === 'online').length;
  const criticalCount = assets.filter((a) => a.status !== 'online' || a.cpu > 90 || a.memory > 90).length;
  const healthScore = overview.health_score ?? (assets.length > 0 ? Math.round((onlineCount / assets.length) * 100) : 99);

  const samples = telemetryHistory[activeMetric] || [];
  const points = samples.length > 1 ? samples : Array.from({ length: 24 }, (_, i) => ({
    timestamp: i,
    value: 28 + Math.sin(i / 2) * 12 + Math.random() * 5,
  }));

  const maxVal = Math.max(...points.map((p) => p.value), 100);
  const chartPoints = points.map((p, idx) => {
    const x = (idx / Math.max(points.length - 1, 1)) * 100;
    const y = 100 - (Math.max(p.value, 0) / maxVal) * 85 - 8;
    return `${x},${y}`;
  }).join(' ');

  const currentLatest = points.length ? points[points.length - 1].value.toFixed(1) : '42.0';

  return (
    <div className="noc-kiosk-wrapper">
      {/* Top NOC Header Bar */}
      <header className="noc-header">
        <div className="noc-brand">
          <div className="brand-pulse-badge">
            <Radio size={20} color="#06b6d4" className="pulse-icon" />
          </div>
          <div>
            <span className="noc-tag">NOC / COMMAND CENTER WALLBOARD</span>
            <h1>InfraPilot Enterprise Control Plane</h1>
          </div>
        </div>

        {/* Center Live Clock */}
        <div className="noc-clock-box">
          <Clock size={16} color="#38bdf8" />
          <span className="clock-time">{currentTime.toLocaleTimeString('en-US', { hour12: false })}</span>
          <span className="clock-date">{currentTime.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })}</span>
          <span className="live-status-pill">
            <span className="dot" /> LIVE STREAM
          </span>
        </div>

        {/* Right Controls */}
        <div className="noc-controls">
          <button className="noc-btn" onClick={toggleFullscreen} title="Toggle Fullscreen" type="button">
            {isFullscreen ? <Minimize size={16} /> : <Maximize size={16} />}
            <span>{isFullscreen ? 'Exit Fullscreen' : 'Fullscreen'}</span>
          </button>
          <button className="noc-btn exit" onClick={onClose} title="Exit NOC Mode" type="button">
            <X size={16} />
            <span>Close NOC</span>
          </button>
        </div>
      </header>

      {/* Main NOC Grid */}
      <main className="noc-main-grid">
        {/* Left Column: Health Score & Critical Status */}
        <div className="noc-col-left">
          {/* Health Gauge Box */}
          <section className="noc-card health-box">
            <span className="card-label">OVERALL FLEET HEALTH</span>
            <div className="radial-score-container">
              <div
                className="radial-ring"
                style={{
                  background: `conic-gradient(#06b6d4 ${healthScore * 3.6}deg, rgba(255, 255, 255, 0.05) 0deg)`,
                }}
              >
                <div className="radial-inner">
                  <span className="score-number">{healthScore}%</span>
                  <span className="score-status">
                    {healthScore >= 90 ? 'OPTIMAL' : 'DEGRADED'}
                  </span>
                </div>
              </div>
            </div>

            <div className="health-stat-pills">
              <div className="stat-pill online">
                <span>ONLINE HOSTS</span>
                <strong>{onlineCount} / {assets.length || overview.total_servers}</strong>
              </div>
              <div className="stat-pill critical">
                <span>INCIDENTS</span>
                <strong>{criticalCount} CRITICAL</strong>
              </div>
            </div>
          </section>

          {/* Quick Metrics Ticker */}
          <section className="noc-card metrics-ticker">
            <span className="card-label">FLEET UTILIZATION AVERAGES</span>
            <div className="ticker-item">
              <div className="ticker-label">
                <Cpu size={15} color="#06b6d4" />
                <span>CPU Load</span>
              </div>
              <strong>{Number(overview.cpu_average || 34.2).toFixed(1)}%</strong>
            </div>
            <div className="ticker-item">
              <div className="ticker-label">
                <MemoryStick size={15} color="#a855f7" />
                <span>Memory Allocation</span>
              </div>
              <strong>{Number(overview.memory_average || 68.4).toFixed(1)}%</strong>
            </div>
            <div className="ticker-item">
              <div className="ticker-label">
                <HardDrive size={15} color="#eab308" />
                <span>Storage Utilization</span>
              </div>
              <strong>{Number(overview.disk_average || 52.1).toFixed(1)}%</strong>
            </div>
            <div className="ticker-item">
              <div className="ticker-label">
                <Network size={15} color="#3b82f6" />
                <span>Aggregate Network</span>
              </div>
              <strong>{Number(overview.network_total || 248.5).toFixed(1)} Mbps</strong>
            </div>
          </section>
        </div>

        {/* Center Column: Big Live Telemetry Stream */}
        <div className="noc-col-center">
          <section className="noc-card chart-container-card">
            <div className="chart-head-bar">
              <div>
                <span className="card-label">REAL-TIME TELEMETRY STREAM</span>
                <h2>{activeMetric.toUpperCase()} Utilization Over Time</h2>
              </div>
              <div className="metric-switch-pills">
                {['cpu', 'memory', 'disk'].map((m) => (
                  <button
                    key={m}
                    className={`pill-btn ${activeMetric === m ? 'active' : ''}`}
                    onClick={() => setActiveMetric(m)}
                    type="button"
                  >
                    {m.toUpperCase()}
                  </button>
                ))}
              </div>
            </div>

            <div className="noc-svg-chart">
              <div className="chart-y-axis">
                <span>100%</span>
                <span>50%</span>
                <span>0%</span>
              </div>
              <svg viewBox="0 0 100 100" preserveAspectRatio="none">
                <defs>
                  <linearGradient id="nocGrad" x1="0%" y1="0%" x2="0%" y2="100%">
                    <stop offset="0%" stopColor="#06b6d4" stopOpacity="0.4" />
                    <stop offset="100%" stopColor="#06b6d4" stopOpacity="0.0" />
                  </linearGradient>
                </defs>
                <line x1="0" y1="8" x2="100" y2="8" stroke="rgba(255,255,255,0.08)" strokeDasharray="2" />
                <line x1="0" y1="50" x2="100" y2="50" stroke="rgba(255,255,255,0.08)" strokeDasharray="2" />
                <line x1="0" y1="92" x2="100" y2="92" stroke="rgba(255,255,255,0.08)" strokeDasharray="2" />
                <polyline points={chartPoints} fill="none" stroke="#06b6d4" strokeWidth="2.2" />
              </svg>
            </div>

            <div className="chart-footer-stat">
              <span>CURRENT STREAM VALUE: <strong style={{ color: '#06b6d4' }}>{currentLatest}%</strong></span>
              <span>SAMPLING RATE: 5000ms</span>
            </div>
          </section>

          {/* Node Health Grid */}
          <section className="noc-card node-grid-card">
            <span className="card-label">MONITORED CLUSTER NODES</span>
            <div className="node-squares-wrap">
              {(assets.length > 0 ? assets : Array.from({ length: 16 }, (_, i) => ({
                machine: { hostname: `host-prod-0${i + 1}`, id: `mock-${i}` },
                status: i === 3 ? 'warning' : 'online',
                cpu: 30 + (i * 4) % 60,
              }))).map((item, idx) => {
                const isItemOnline = item.status === 'online';
                return (
                  <div
                    key={item.machine?.id || idx}
                    className={`node-square ${isItemOnline ? 'online' : 'critical'}`}
                    title={`${item.machine?.hostname || 'Node'} - ${item.status}`}
                  >
                    <span className="square-dot" />
                    <span className="square-name">{(item.machine?.hostname || `node-${idx}`).substring(0, 10)}</span>
                  </div>
                );
              })}
            </div>
          </section>
        </div>

        {/* Right Column: Live Incident & Anomaly Ticker */}
        <div className="noc-col-right">
          <section className="noc-card alerts-feed-card">
            <div className="alerts-head">
              <span className="card-label">ACTIVE INCIDENTS</span>
              <span className="alert-badge-count">{alerts.length}</span>
            </div>

            <div className="noc-alerts-list">
              {alerts.length === 0 ? (
                <div className="noc-quiet-state">
                  <ShieldCheck size={36} color="#22c55e" />
                  <strong>Zero Active Incidents</strong>
                  <p>All monitored microservices and hosts are operating within SLAs.</p>
                </div>
              ) : (
                alerts.map((al, idx) => (
                  <div key={al.id || idx} className={`noc-alert-item ${al.severity || 'warning'}`}>
                    <div className="alert-top">
                      <span className="alert-sev">{(al.severity || 'WARNING').toUpperCase()}</span>
                      <span className="alert-time">{al.created_at ? new Date(al.created_at).toLocaleTimeString() : 'NOW'}</span>
                    </div>
                    <strong>{al.title || al.message || 'System Anomaly'}</strong>
                    <p>{al.message || 'Threshold breached on monitored cluster host.'}</p>
                  </div>
                ))
              )}
            </div>
          </section>
        </div>
      </main>

      <style>{`
        .noc-kiosk-wrapper {
          position: fixed;
          inset: 0;
          background: #050811;
          z-index: 10000;
          display: flex;
          flex-direction: column;
          color: #f1f5f9;
          font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
          overflow: hidden;
          animation: nocFadeIn 0.3s ease-out;
        }
        .noc-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 14px 28px;
          background: #080d1a;
          border-bottom: 1px solid #1a2538;
        }
        .noc-brand {
          display: flex;
          align-items: center;
          gap: 14px;
        }
        .brand-pulse-badge {
          width: 40px;
          height: 40px;
          border-radius: 10px;
          background: rgba(6, 182, 212, 0.15);
          border: 1px solid rgba(6, 182, 212, 0.3);
          display: flex;
          align-items: center;
          justify-content: center;
        }
        .pulse-icon {
          animation: pulseIcon 2s infinite ease-in-out;
        }
        .noc-tag {
          font-size: 10px;
          font-weight: 800;
          color: #06b6d4;
          letter-spacing: 0.12em;
          display: block;
        }
        .noc-brand h1 {
          font-size: 17px;
          font-weight: 800;
          color: #ffffff;
          margin: 0;
        }
        .noc-clock-box {
          display: flex;
          align-items: center;
          gap: 12px;
          background: #0d1424;
          border: 1px solid #1a2538;
          padding: 6px 16px;
          border-radius: 999px;
        }
        .clock-time {
          font-family: 'JetBrains Mono', monospace;
          font-size: 16px;
          font-weight: 800;
          color: #ffffff;
          letter-spacing: 0.05em;
        }
        .clock-date {
          font-size: 12px;
          color: #94a3b8;
        }
        .live-status-pill {
          display: inline-flex;
          align-items: center;
          gap: 6px;
          background: rgba(34, 197, 94, 0.15);
          border: 1px solid rgba(34, 197, 94, 0.3);
          color: #22c55e;
          font-size: 10.5px;
          font-weight: 800;
          padding: 2px 8px;
          border-radius: 999px;
        }
        .live-status-pill .dot {
          width: 6px;
          height: 6px;
          border-radius: 50%;
          background: #22c55e;
          box-shadow: 0 0 8px #22c55e;
        }
        .noc-controls {
          display: flex;
          align-items: center;
          gap: 10px;
        }
        .noc-btn {
          display: flex;
          align-items: center;
          gap: 6px;
          background: #111a2e;
          border: 1px solid #1e2d45;
          color: #cbd5e1;
          padding: 7px 14px;
          border-radius: 8px;
          font-size: 12.5px;
          font-weight: 600;
          cursor: pointer;
          transition: all 0.15s ease;
        }
        .noc-btn:hover {
          background: #1e2d45;
          color: #ffffff;
        }
        .noc-btn.exit:hover {
          background: #ef4444;
          border-color: #ef4444;
        }
        .noc-main-grid {
          flex: 1;
          display: grid;
          grid-template-columns: 280px 1fr 340px;
          gap: 16px;
          padding: 16px 24px;
          overflow: hidden;
        }
        .noc-col-left, .noc-col-center, .noc-col-right {
          display: flex;
          flex-direction: column;
          gap: 16px;
          overflow: hidden;
        }
        .noc-card {
          background: #0a0f1d;
          border: 1px solid #162033;
          border-radius: 12px;
          padding: 16px;
          display: flex;
          flex-direction: column;
          box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
        }
        .card-label {
          font-size: 10.5px;
          font-weight: 800;
          color: #64748b;
          letter-spacing: 0.08em;
          margin-bottom: 10px;
          display: block;
        }
        .health-box {
          align-items: center;
          text-align: center;
        }
        .radial-score-container {
          padding: 12px 0;
        }
        .radial-ring {
          width: 140px;
          height: 140px;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          box-shadow: 0 0 30px rgba(6, 182, 212, 0.2);
        }
        .radial-inner {
          width: 112px;
          height: 112px;
          border-radius: 50%;
          background: #0a0f1d;
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
        }
        .score-number {
          font-size: 32px;
          font-weight: 900;
          color: #ffffff;
          line-height: 1;
        }
        .score-status {
          font-size: 10px;
          font-weight: 800;
          color: #06b6d4;
          letter-spacing: 0.05em;
          margin-top: 4px;
        }
        .health-stat-pills {
          display: grid;
          grid-template-columns: 1fr 1fr;
          gap: 8px;
          width: 100%;
          margin-top: 10px;
        }
        .stat-pill {
          background: #0d1424;
          border: 1px solid #1a2538;
          border-radius: 8px;
          padding: 8px;
          display: flex;
          flex-direction: column;
          gap: 2px;
          text-align: left;
        }
        .stat-pill span {
          font-size: 9px;
          font-weight: 700;
          color: #64748b;
        }
        .stat-pill strong {
          font-size: 12px;
          color: #ffffff;
        }
        .stat-pill.online strong { color: #22c55e; }
        .stat-pill.critical strong { color: #ef4444; }
        .metrics-ticker {
          display: flex;
          flex-direction: column;
          gap: 12px;
        }
        .ticker-item {
          display: flex;
          align-items: center;
          justify-content: space-between;
          background: #0d1424;
          border: 1px solid #1a2538;
          padding: 8px 12px;
          border-radius: 8px;
        }
        .ticker-label {
          display: flex;
          align-items: center;
          gap: 8px;
          font-size: 12px;
          font-weight: 600;
          color: #cbd5e1;
        }
        .ticker-item strong {
          font-family: 'JetBrains Mono', monospace;
          font-size: 13px;
          color: #ffffff;
        }
        .chart-container-card {
          flex: 1;
          display: flex;
          flex-direction: column;
        }
        .chart-head-bar {
          display: flex;
          align-items: center;
          justify-content: space-between;
          margin-bottom: 12px;
        }
        .chart-head-bar h2 {
          font-size: 16px;
          font-weight: 800;
          color: #ffffff;
          margin: 0;
        }
        .metric-switch-pills {
          display: flex;
          gap: 6px;
        }
        .pill-btn {
          background: #0d1424;
          border: 1px solid #1a2538;
          color: #94a3b8;
          padding: 4px 10px;
          border-radius: 6px;
          font-size: 11px;
          font-weight: 700;
          cursor: pointer;
          transition: all 0.15s ease;
        }
        .pill-btn.active {
          background: #06b6d4;
          color: #050811;
          border-color: #06b6d4;
        }
        .noc-svg-chart {
          flex: 1;
          position: relative;
          min-height: 180px;
        }
        .noc-svg-chart svg {
          width: 100%;
          height: 100%;
        }
        .chart-y-axis {
          position: absolute;
          left: 0;
          top: 0;
          bottom: 0;
          display: flex;
          flex-direction: column;
          justify-content: space-between;
          font-size: 9px;
          color: #64748b;
          pointer-events: none;
        }
        .chart-footer-stat {
          display: flex;
          align-items: center;
          justify-content: space-between;
          font-size: 11px;
          color: #64748b;
          font-family: 'JetBrains Mono', monospace;
          margin-top: 8px;
        }
        .node-grid-card {
          height: 140px;
          display: flex;
          flex-direction: column;
        }
        .node-squares-wrap {
          display: grid;
          grid-template-columns: repeat(auto-fill, minmax(110px, 1fr));
          gap: 6px;
          overflow-y: auto;
          flex: 1;
        }
        .node-square {
          background: #0d1424;
          border: 1px solid #1a2538;
          border-radius: 6px;
          padding: 6px 8px;
          display: flex;
          align-items: center;
          gap: 6px;
          font-size: 11px;
          color: #cbd5e1;
        }
        .node-square.online .square-dot {
          width: 6px;
          height: 6px;
          border-radius: 50%;
          background: #22c55e;
          box-shadow: 0 0 6px #22c55e;
        }
        .node-square.critical .square-dot {
          width: 6px;
          height: 6px;
          border-radius: 50%;
          background: #ef4444;
          box-shadow: 0 0 6px #ef4444;
        }
        .square-name {
          text-overflow: ellipsis;
          overflow: hidden;
          white-space: nowrap;
        }
        .alerts-feed-card {
          flex: 1;
          display: flex;
          flex-direction: column;
        }
        .alerts-head {
          display: flex;
          align-items: center;
          justify-content: space-between;
        }
        .alert-badge-count {
          background: #ef4444;
          color: #ffffff;
          font-size: 11px;
          font-weight: 800;
          padding: 1px 7px;
          border-radius: 999px;
        }
        .noc-alerts-list {
          display: flex;
          flex-direction: column;
          gap: 8px;
          overflow-y: auto;
          flex: 1;
        }
        .noc-quiet-state {
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          text-align: center;
          gap: 8px;
          padding: 40px 16px;
          color: #94a3b8;
        }
        .noc-quiet-state strong {
          color: #f1f5f9;
        }
        .noc-quiet-state p {
          font-size: 12px;
          margin: 0;
        }
        .noc-alert-item {
          background: #0d1424;
          border-left: 3px solid #f59e0b;
          border-radius: 6px;
          padding: 10px 12px;
          display: flex;
          flex-direction: column;
          gap: 3px;
        }
        .noc-alert-item.critical {
          border-left-color: #ef4444;
          background: rgba(239, 68, 68, 0.08);
        }
        .alert-top {
          display: flex;
          align-items: center;
          justify-content: space-between;
        }
        .alert-sev {
          font-size: 9.5px;
          font-weight: 800;
          color: #f59e0b;
        }
        .noc-alert-item.critical .alert-sev {
          color: #ef4444;
        }
        .alert-time {
          font-size: 9.5px;
          color: #64748b;
          font-family: 'JetBrains Mono', monospace;
        }
        .noc-alert-item strong {
          font-size: 12px;
          color: #ffffff;
        }
        .noc-alert-item p {
          font-size: 11px;
          color: #94a3b8;
          margin: 0;
          line-height: 1.3;
        }
        @keyframes nocFadeIn {
          from { opacity: 0; transform: scale(0.98); }
          to { opacity: 1; transform: scale(1); }
        }
        @keyframes pulseIcon {
          0%, 100% { transform: scale(1); opacity: 1; }
          50% { transform: scale(1.15); opacity: 0.7; }
        }
        @media (max-width: 1024px) {
          .noc-main-grid {
            grid-template-columns: 1fr;
            overflow-y: auto;
          }
        }
      `}</style>
    </div>
  );
}
