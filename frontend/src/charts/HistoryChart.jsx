import React, { useMemo } from 'react';

const COLORS = {
  cpu: '#06b6d4',
  memory: '#22c55e',
  disk: '#eab308',
  network: '#a78bfa',
};

export default function HistoryChart({ samples = [], range = '1h' }) {
  const chartData = useMemo(() => {
    if (!Array.isArray(samples) || samples.length === 0) return [];
    return [...samples].sort((a, b) => {
      const ta = new Date(a.timestamp || a.created_at || 0).getTime();
      const tb = new Date(b.timestamp || b.created_at || 0).getTime();
      return ta - tb;
    });
  }, [samples]);

  if (chartData.length === 0) {
    return (
      <div
        style={{
          padding: '60px',
          textAlign: 'center',
          color: '#64748b',
          backgroundColor: '#0d1220',
          border: '1px solid #1f2e44',
          borderRadius: '12px',
        }}
      >
        No historical telemetry samples available.
      </div>
    );
  }

  const width = 800;
  const height = 300;
  const paddingX = 50;
  const paddingY = 30;

  const count = chartData.length;
  const getCoordinates = (metricKey) => {
    return chartData.map((sample, idx) => {
      const x = paddingX + (idx / Math.max(1, count - 1)) * (width - 2 * paddingX);
      let val = 0;
      if (metricKey === 'cpu') val = sample.cpu_usage ?? sample.cpu ?? 0;
      else if (metricKey === 'memory') val = sample.memory_usage ?? sample.memory ?? 0;
      else if (metricKey === 'disk') val = sample.disk_usage ?? sample.disk ?? 0;
      else if (metricKey === 'network') val = Number(sample.upload_mbps || 0) + Number(sample.download_mbps || 0);

      val = Math.max(0, Math.min(val, 100));
      const y = height - paddingY - (val / 100) * (height - 2 * paddingY);
      return { x, y, value: val };
    });
  };

  const cpuCoords = getCoordinates('cpu');
  const memCoords = getCoordinates('memory');
  const diskCoords = getCoordinates('disk');

  const getPath = (coords) =>
    coords.reduce(
      (acc, point, idx) => (idx === 0 ? `M ${point.x} ${point.y}` : `${acc} L ${point.x} ${point.y}`),
      ''
    );

  return (
    <div
      style={{
        backgroundColor: '#0d1220',
        border: '1px solid #1f2e44',
        borderRadius: '12px',
        padding: '20px',
        display: 'flex',
        flexDirection: 'column',
        gap: '16px',
      }}
    >
      <div style={{ display: 'flex', gap: '16px', alignItems: 'center', fontSize: '12px' }}>
        <span style={{ color: COLORS.cpu, fontWeight: 600 }}>● CPU %</span>
        <span style={{ color: COLORS.memory, fontWeight: 600 }}>● Memory %</span>
        <span style={{ color: COLORS.disk, fontWeight: 600 }}>● Disk %</span>
      </div>

      <div style={{ width: '100%', overflowX: 'auto' }}>
        <svg viewBox={`0 0 ${width} ${height}`} style={{ width: '100%', height: 'auto' }}>
          <line x1={paddingX} y1={paddingY} x2={width - paddingX} y2={paddingY} stroke="#1f2e44" strokeDasharray="4 4" />
          <line x1={paddingX} y1={height / 2} x2={width - paddingX} y2={height / 2} stroke="#1f2e44" strokeDasharray="4 4" />
          <line x1={paddingX} y1={height - paddingY} x2={width - paddingX} y2={height - paddingY} stroke="#1f2e44" />

          <path d={getPath(cpuCoords)} fill="none" stroke={COLORS.cpu} strokeWidth="2.5" />
          <path d={getPath(memCoords)} fill="none" stroke={COLORS.memory} strokeWidth="2.5" />
          <path d={getPath(diskCoords)} fill="none" stroke={COLORS.disk} strokeWidth="2.5" />
        </svg>
      </div>
    </div>
  );
}