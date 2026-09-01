import React, { useState, useMemo, useRef, useEffect } from 'react';
import { Cpu, HardDrive, Network, Layers, ChevronDown } from 'lucide-react';

export default function DashboardCharts({ samples = [], range = '1h', onRangeChange }) {
  const [activeMetrics, setActiveMetrics] = useState({
    cpu: true,
    memory: true,
    disk: false,
    network: false,
  });

  const svgRef = useRef(null);
  const [hoverIndex, setHoverIndex] = useState(null);
  const [tooltipPos, setTooltipPos] = useState({ x: 0, y: 0 });

  // 1. Process samples
  const chartData = useMemo(() => {
    if (!Array.isArray(samples) || samples.length === 0) return [];
    
    // Sort chronologically
    return [...samples].sort((a, b) => {
      const ta = new Date(a.timestamp || a.created_at || 0).getTime();
      const tb = new Date(b.timestamp || b.created_at || 0).getTime();
      return ta - tb;
    });
  }, [samples]);

  // 2. Math parameters
  const width = 800;
  const height = 280;
  const paddingX = 60;
  const paddingY = 40;

  // Max network value in samples to scale network line
  const maxNetwork = useMemo(() => {
    if (chartData.length === 0) return 10;
    const maxVal = Math.max(
      ...chartData.map(s => Number(s.upload_mbps || 0) + Number(s.download_mbps || 0)),
      10 // baseline
    );
    return Math.ceil(maxVal * 1.1); // add 10% headroom
  }, [chartData]);

  // Transform coordinates helper
  const getCoordinates = (metricKey) => {
    if (chartData.length === 0) return [];
    
    const count = chartData.length;
    return chartData.map((sample, idx) => {
      // X coordinate spaced evenly
      const x = paddingX + (idx / Math.max(1, count - 1)) * (width - 2 * paddingX);
      
      // Get raw value
      let val = 0;
      let max = 100;
      if (metricKey === 'cpu') val = sample.cpu_usage ?? sample.cpu ?? 0;
      else if (metricKey === 'memory') val = sample.memory_usage ?? sample.memory ?? 0;
      else if (metricKey === 'disk') val = sample.disk_usage ?? sample.disk ?? 0;
      else if (metricKey === 'network') {
        val = Number(sample.upload_mbps || 0) + Number(sample.download_mbps || 0);
        max = maxNetwork;
      }
      
      // Clip value
      val = Math.max(0, Math.min(val, max));
      
      // Y coordinate (SVG starts at top left, so we invert)
      const y = height - paddingY - (val / max) * (height - 2 * paddingY);
      
      return { x, y, value: val };
    });
  };

  const lines = useMemo(() => {
    return {
      cpu: getCoordinates('cpu'),
      memory: getCoordinates('memory'),
      disk: getCoordinates('disk'),
      network: getCoordinates('network'),
    };
  }, [chartData, maxNetwork]);

  // Create SVG path string
  const getPathString = (coords) => {
    if (coords.length === 0) return '';
    return coords.reduce((acc, point, idx) => {
      return idx === 0 ? `M ${point.x} ${point.y}` : `${acc} L ${point.x} ${point.y}`;
    }, '');
  };

  // Create SVG area fill path string (closes path to bottom axis)
  const getAreaString = (coords) => {
    if (coords.length === 0) return '';
    const linePath = getPathString(coords);
    const lastPoint = coords[coords.length - 1];
    const firstPoint = coords[0];
    const yBaseline = height - paddingY;
    return `${linePath} L ${lastPoint.x} ${yBaseline} L ${firstPoint.x} ${yBaseline} Z`;
  };

  // Handle SVG Mousemove for Interactive Tooltip
  const handleMouseMove = (e) => {
    if (!svgRef.current || chartData.length === 0) return;
    
    const rect = svgRef.current.getBoundingClientRect();
    const clientX = e.clientX - rect.left;
    const clientY = e.clientY - rect.top;

    // Scale back to viewBox coordinates (0 to 800 width)
    const svgX = (clientX / rect.width) * width;
    
    // Find closest data point index based on x position
    const step = (width - 2 * paddingX) / Math.max(1, chartData.length - 1);
    let index = Math.round((svgX - paddingX) / step);
    index = Math.max(0, Math.min(index, chartData.length - 1));
    
    setHoverIndex(index);
    
    // Position tooltip
    const actualX = paddingX + index * step;
    setTooltipPos({
      x: (actualX / width) * rect.width,
      y: (clientY / rect.height) * rect.height - 80,
    });
  };

  const handleMouseLeave = () => {
    setHoverIndex(null);
  };

  // Generate X axis tick labels (5 labels)
  const xTicks = useMemo(() => {
    if (chartData.length === 0) return [];
    
    const count = chartData.length;
    const indexes = [0, Math.floor(count * 0.25), Math.floor(count * 0.5), Math.floor(count * 0.75), count - 1].filter(
      (val, idx, self) => self.indexOf(val) === idx && val >= 0 && val < count
    );

    return indexes.map(idx => {
      const sample = chartData[idx];
      const timeStr = sample.timestamp || sample.created_at || '';
      if (!timeStr) return '';
      
      try {
        const d = new Date(timeStr);
        if (range === '24h') {
          return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
        }
        return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' });
      } catch {
        return 'now';
      }
    });
  }, [chartData, range]);

  const toggleMetric = (key) => {
    setActiveMetrics(prev => ({ ...prev, [key]: !prev[key] }));
  };

  return (
    <section className="chart-panel">
      <div className="chart-head">
        <div className="chart-title-area">
          <span className="eyebrow">WebSocket Stream</span>
          <h2>Resource History</h2>
        </div>

        <div className="chart-controls">
          {/* Timeframe Selectors */}
          <div className="range-selector">
            {['5m', '15m', '1h', '24h'].map(r => (
              <button
                key={r}
                className={`range-btn ${range === r ? 'active' : ''}`}
                onClick={() => onRangeChange && onRangeChange(r)}
                type="button"
              >
                {r === '5m' ? '5 Min' : r === '15m' ? '15 Min' : r === '1h' ? '1 Hour' : '24 Hours'}
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* KPI Toggles */}
      <div className="metric-selectors">
        <button
          className={`selector-btn cpu ${activeMetrics.cpu ? 'active' : ''}`}
          onClick={() => toggleMetric('cpu')}
          type="button"
        >
          <span className="bullet-cpu" />
          <Cpu size={14} />
          <span>CPU Usage</span>
        </button>
        <button
          className={`selector-btn memory ${activeMetrics.memory ? 'active' : ''}`}
          onClick={() => toggleMetric('memory')}
          type="button"
        >
          <span className="bullet-mem" />
          <Layers size={14} />
          <span>Memory Usage</span>
        </button>
        <button
          className={`selector-btn disk ${activeMetrics.disk ? 'active' : ''}`}
          onClick={() => toggleMetric('disk')}
          type="button"
        >
          <span className="bullet-disk" />
          <HardDrive size={14} />
          <span>Disk Usage</span>
        </button>
        <button
          className={`selector-btn network ${activeMetrics.network ? 'active' : ''}`}
          onClick={() => toggleMetric('network')}
          type="button"
        >
          <span className="bullet-net" />
          <Network size={14} />
          <span>Network Speed</span>
        </button>
      </div>

      {/* Main SVG Graph */}
      <div className="svg-container" style={{ position: 'relative' }}>
        {chartData.length === 0 ? (
          <div className="chart-placeholder">
            <Activity className="pulse-slow" size={32} />
            <strong>Awaiting Telemetry Ingestion...</strong>
            <p>Historical agent metrics will draw here in real-time.</p>
          </div>
        ) : (
          <svg
            ref={svgRef}
            viewBox={`0 0 ${width} ${height}`}
            className="chart-svg"
            onMouseMove={handleMouseMove}
            onMouseLeave={handleMouseLeave}
          >
            {/* Definitions for gradients */}
            <defs>
              <linearGradient id="cpu-grad" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stopColor="#06b6d4" stopOpacity="0.25" />
                <stop offset="100%" stopColor="#06b6d4" stopOpacity="0.0" />
              </linearGradient>
              <linearGradient id="mem-grad" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stopColor="#22c55e" stopOpacity="0.25" />
                <stop offset="100%" stopColor="#22c55e" stopOpacity="0.0" />
              </linearGradient>
              <linearGradient id="disk-grad" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stopColor="#f59e0b" stopOpacity="0.25" />
                <stop offset="100%" stopColor="#f59e0b" stopOpacity="0.0" />
              </linearGradient>
              <linearGradient id="net-grad" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stopColor="#a78bfa" stopOpacity="0.25" />
                <stop offset="100%" stopColor="#a78bfa" stopOpacity="0.0" />
              </linearGradient>
            </defs>

            {/* Horizontal Grid lines */}
            {[0, 0.25, 0.5, 0.75, 1].map((ratio, index) => {
              const y = paddingY + ratio * (height - 2 * paddingY);
              const val = Math.round(100 - ratio * 100);
              return (
                <g key={index}>
                  <line
                    x1={paddingX}
                    y1={y}
                    x2={width - paddingX}
                    y2={y}
                    stroke="#1f2e44"
                    strokeWidth="1"
                    strokeDasharray="4 4"
                  />
                  {/* Left percentages Y-axis */}
                  <text
                    x={paddingX - 10}
                    y={y + 4}
                    fill="#64748b"
                    fontSize="11"
                    textAnchor="end"
                  >
                    {val}%
                  </text>
                  {/* Right Mbps Y-axis */}
                  {activeMetrics.network && (
                    <text
                      x={width - paddingX + 10}
                      y={y + 4}
                      fill="#a78bfa"
                      fontSize="11"
                      textAnchor="start"
                    >
                      {Math.round((maxNetwork - ratio * maxNetwork))} Mb/s
                    </text>
                  )}
                </g>
              );
            })}

            {/* Vertical grid lines & labels */}
            {xTicks.map((tick, idx) => {
              if (!tick) return null;
              const x = paddingX + (idx / Math.max(1, xTicks.length - 1)) * (width - 2 * paddingX);
              return (
                <g key={idx}>
                  <line
                    x1={x}
                    y1={paddingY}
                    x2={x}
                    y2={height - paddingY}
                    stroke="#1f2e44"
                    strokeWidth="1"
                    opacity="0.5"
                  />
                  <text
                    x={x}
                    y={height - paddingY + 20}
                    fill="#64748b"
                    fontSize="11"
                    textAnchor="middle"
                  >
                    {tick}
                  </text>
                </g>
              );
            })}

            {/* Draw Area Fills & Lines for active metrics */}
            {activeMetrics.disk && lines.disk.length > 0 && (
              <>
                <path d={getAreaString(lines.disk)} fill="url(#disk-grad)" />
                <path d={getPathString(lines.disk)} fill="none" stroke="#f59e0b" strokeWidth="2.5" strokeLinecap="round" />
              </>
            )}

            {activeMetrics.network && lines.network.length > 0 && (
              <>
                <path d={getAreaString(lines.network)} fill="url(#net-grad)" />
                <path d={getPathString(lines.network)} fill="none" stroke="#a78bfa" strokeWidth="2.5" strokeLinecap="round" />
              </>
            )}

            {activeMetrics.memory && lines.memory.length > 0 && (
              <>
                <path d={getAreaString(lines.memory)} fill="url(#mem-grad)" />
                <path d={getPathString(lines.memory)} fill="none" stroke="#22c55e" strokeWidth="2.5" strokeLinecap="round" />
              </>
            )}

            {activeMetrics.cpu && lines.cpu.length > 0 && (
              <>
                <path d={getAreaString(lines.cpu)} fill="url(#cpu-grad)" />
                <path d={getPathString(lines.cpu)} fill="none" stroke="#06b6d4" strokeWidth="2.5" strokeLinecap="round" />
              </>
            )}

            {/* Hover indicator line & points */}
            {hoverIndex !== null && (
              <>
                {/* Vertical marker line */}
                <line
                  x1={paddingX + (hoverIndex / Math.max(1, chartData.length - 1)) * (width - 2 * paddingX)}
                  y1={paddingY}
                  x2={paddingX + (hoverIndex / Math.max(1, chartData.length - 1)) * (width - 2 * paddingX)}
                  y2={height - paddingY}
                  stroke="#38bdf8"
                  strokeWidth="1.5"
                />

                {/* Bullets on data points */}
                {activeMetrics.cpu && lines.cpu[hoverIndex] && (
                  <circle cx={lines.cpu[hoverIndex].x} cy={lines.cpu[hoverIndex].y} r="5" fill="#06b6d4" stroke="#080c14" strokeWidth="2" />
                )}
                {activeMetrics.memory && lines.memory[hoverIndex] && (
                  <circle cx={lines.memory[hoverIndex].x} cy={lines.memory[hoverIndex].y} r="5" fill="#22c55e" stroke="#080c14" strokeWidth="2" />
                )}
                {activeMetrics.disk && lines.disk[hoverIndex] && (
                  <circle cx={lines.disk[hoverIndex].x} cy={lines.disk[hoverIndex].y} r="5" fill="#f59e0b" stroke="#080c14" strokeWidth="2" />
                )}
                {activeMetrics.network && lines.network[hoverIndex] && (
                  <circle cx={lines.network[hoverIndex].x} cy={lines.network[hoverIndex].y} r="5" fill="#a78bfa" stroke="#080c14" strokeWidth="2" />
                )}
              </>
            )}
          </svg>
        )}

        {/* Live Hover Tooltip Panel */}
        {hoverIndex !== null && chartData[hoverIndex] && (
          <div
            className="chart-tooltip"
            style={{
              position: 'absolute',
              left: `${tooltipPos.x + 15}px`,
              top: `${tooltipPos.y}px`,
              pointerEvents: 'none',
            }}
          >
            <div className="tooltip-time">
              {new Date(chartData[hoverIndex].timestamp || chartData[hoverIndex].created_at).toLocaleTimeString()}
            </div>
            <div className="tooltip-rows">
              {activeMetrics.cpu && (
                <div className="tooltip-row cpu">
                  <span className="dot" />
                  <span>CPU:</span>
                  <strong>{(chartData[hoverIndex].cpu_usage ?? chartData[hoverIndex].cpu ?? 0).toFixed(1)}%</strong>
                </div>
              )}
              {activeMetrics.memory && (
                <div className="tooltip-row mem">
                  <span className="dot" />
                  <span>RAM:</span>
                  <strong>{(chartData[hoverIndex].memory_usage ?? chartData[hoverIndex].memory ?? 0).toFixed(1)}%</strong>
                </div>
              )}
              {activeMetrics.disk && (
                <div className="tooltip-row disk">
                  <span className="dot" />
                  <span>Disk:</span>
                  <strong>{(chartData[hoverIndex].disk_usage ?? chartData[hoverIndex].disk ?? 0).toFixed(1)}%</strong>
                </div>
              )}
              {activeMetrics.network && (
                <div className="tooltip-row net">
                  <span className="dot" />
                  <span>Net Speed:</span>
                  <strong>
                    {(
                      Number(chartData[hoverIndex].upload_mbps || 0) +
                      Number(chartData[hoverIndex].download_mbps || 0)
                    ).toFixed(2)} Mb/s
                  </strong>
                </div>
              )}
            </div>
          </div>
        )}
      </div>

      <style>{`
        .chart-panel {
          background-color: #0d1220;
          border: 1px solid #1f2e44;
          border-radius: 12px;
          padding: 24px;
          display: flex;
          flex-direction: column;
          gap: 16px;
        }
        .chart-head {
          display: flex;
          justify-content: space-between;
          align-items: center;
          flex-wrap: wrap;
          gap: 12px;
        }
        .chart-title-area h2 {
          font-size: 18px;
          font-weight: 700;
          color: #f1f5f9;
        }
        .chart-controls {
          display: flex;
          align-items: center;
          gap: 12px;
        }
        .range-selector {
          display: flex;
          background-color: #080c14;
          border: 1px solid #1f2e44;
          border-radius: 8px;
          padding: 3px;
        }
        .range-btn {
          background: none;
          border: none;
          color: #64748b;
          font-size: 12px;
          font-weight: 600;
          padding: 6px 12px;
          border-radius: 6px;
          transition: all 0.2s;
          cursor: pointer;
        }
        .range-btn:hover {
          color: #cbd5e1;
        }
        .range-btn.active {
          background-color: #1f2e44;
          color: #06b6d4;
        }
        .metric-selectors {
          display: flex;
          flex-wrap: wrap;
          gap: 12px;
          margin-bottom: 8px;
        }
        .selector-btn {
          display: flex;
          align-items: center;
          gap: 8px;
          background-color: #080c14;
          border: 1px solid #1f2e44;
          color: #94a3b8;
          padding: 8px 14px;
          border-radius: 8px;
          font-size: 13px;
          font-weight: 500;
          transition: all 0.2s;
          cursor: pointer;
        }
        .selector-btn:hover {
          color: #e2e8f0;
          border-color: #2e3f5a;
        }
        .selector-btn.active {
          color: #f1f5f9;
        }
        .selector-btn.cpu.active { border-color: #06b6d4; background-color: rgba(6, 182, 212, 0.05); }
        .selector-btn.memory.active { border-color: #22c55e; background-color: rgba(34, 197, 150, 0.05); }
        .selector-btn.disk.active { border-color: #f59e0b; background-color: rgba(245, 158, 11, 0.05); }
        .selector-btn.network.active { border-color: #a78bfa; background-color: rgba(167, 139, 250, 0.05); }
        
        .bullet-cpu, .bullet-mem, .bullet-disk, .bullet-net {
          width: 8px;
          height: 8px;
          border-radius: 50%;
          display: inline-block;
        }
        .bullet-cpu { background-color: #06b6d4; }
        .bullet-mem { background-color: #22c55e; }
        .bullet-disk { background-color: #f59e0b; }
        .bullet-net { background-color: #a78bfa; }

        .svg-container {
          background-color: #080c14;
          border: 1px solid rgba(31, 46, 68, 0.5);
          border-radius: 10px;
          padding: 16px 8px;
          min-height: 250px;
          display: flex;
          align-items: center;
          justify-content: center;
        }
        .chart-svg {
          width: 100%;
          height: 100%;
          overflow: visible;
        }
        .chart-placeholder {
          display: flex;
          flex-direction: column;
          align-items: center;
          gap: 8px;
          color: #64748b;
          text-align: center;
        }
        .chart-placeholder strong {
          color: #94a3b8;
          font-size: 14px;
        }
        .chart-placeholder p {
          font-size: 12px;
          max-width: 320px;
        }
        .pulse-slow {
          animation: pulse-op 2s infinite ease-in-out;
        }
        @keyframes pulse-op {
          0% { opacity: 0.3; }
          50% { opacity: 0.8; }
          100% { opacity: 0.3; }
        }

        /* Tooltip Panel */
        .chart-tooltip {
          background-color: #0d1220;
          border: 1px solid #1f2e44;
          border-radius: 8px;
          padding: 10px 14px;
          box-shadow: 0 10px 25px rgba(0,0,0,0.5);
          min-width: 140px;
          z-index: 10;
        }
        .tooltip-time {
          font-size: 11px;
          font-weight: 700;
          color: #64748b;
          margin-bottom: 6px;
          border-bottom: 1px solid #1f2e44;
          padding-bottom: 4px;
        }
        .tooltip-rows {
          display: flex;
          flex-direction: column;
          gap: 4px;
        }
        .tooltip-row {
          display: flex;
          align-items: center;
          gap: 6px;
          font-size: 12px;
          color: #cbd5e1;
        }
        .tooltip-row .dot {
          width: 6px;
          height: 6px;
          border-radius: 50%;
        }
        .tooltip-row strong {
          margin-left: auto;
          color: #f1f5f9;
        }
        .tooltip-row.cpu .dot { background-color: #06b6d4; }
        .tooltip-row.mem .dot { background-color: #22c55e; }
        .tooltip-row.disk .dot { background-color: #f59e0b; }
        .tooltip-row.net .dot { background-color: #a78bfa; }
      `}</style>
    </section>
  );
}
