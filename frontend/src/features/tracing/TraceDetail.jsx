import React from 'react';
import { GitCommit, Clock, ArrowRight, X } from 'lucide-react';

export default function TraceDetail({ trace, onClose }) {
  if (!trace) return null;

  const totalDuration = trace.duration_ms || 1;
  const spans = trace.spans || [
    { span_id: '1', name: 'HTTP Request Handler', service: 'api-gateway', duration_ms: trace.duration_ms, start_offset: 0 },
    { span_id: '2', name: 'Verify JWT & Org Quota', service: 'auth-service', duration_ms: 12, start_offset: 0 },
    { span_id: '3', name: 'SELECT * FROM metrics', service: 'postgres-db', duration_ms: 78, start_offset: 14 },
    { span_id: '4', name: 'AI Incident Analysis Inference', service: 'ai-engine', duration_ms: 340, start_offset: 93 },
  ];

  const getServiceColor = (service) => {
    switch (service) {
      case 'api-gateway': return '#58a6ff';
      case 'auth-service': return '#3fb950';
      case 'postgres-db': return '#ffa657';
      case 'ai-engine': return '#a855f7';
      case 'redis': return '#f78166';
      default: return '#79c0ff';
    }
  };

  return (
    <div style={{
      position: 'fixed', top: 0, left: 0, right: 0, bottom: 0,
      backgroundColor: 'rgba(0,0,0,0.75)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 1000
    }}>
      <div style={{ backgroundColor: '#161b22', border: '1px solid #30363d', borderRadius: '12px', width: '740px', padding: '24px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
          <div>
            <h3 style={{ margin: 0, color: '#f0f6fc', display: 'flex', alignItems: 'center', gap: '8px', fontSize: '18px' }}>
              <GitCommit size={20} color="#58a6ff" /> Trace Execution Breakdown: {trace.name}
            </h3>
            <div style={{ fontSize: '12px', color: '#8b949e', marginTop: '4px' }}>Trace ID: <code>{trace.trace_id}</code></div>
          </div>
          <button onClick={onClose} style={{ background: 'transparent', border: 'none', color: '#8b949e', cursor: 'pointer' }}>
            <X size={20} />
          </button>
        </div>

        {/* Trace Summary Bar */}
        <div style={{ background: '#0d1117', border: '1px solid #30363d', borderRadius: '8px', padding: '16px', marginBottom: '20px', display: 'flex', justifyContent: 'space-between' }}>
          <div><span style={{ color: '#8b949e', fontSize: '12px' }}>TOTAL LATENCY:</span> <strong style={{ color: trace.duration_ms > 500 ? '#ffa657' : '#3fb950', fontSize: '16px', marginLeft: '6px' }}>{trace.duration_ms} ms</strong></div>
          <div><span style={{ color: '#8b949e', fontSize: '12px' }}>TOTAL SPANS:</span> <strong style={{ color: '#f0f6fc', fontSize: '16px', marginLeft: '6px' }}>{spans.length}</strong></div>
          <div><span style={{ color: '#8b949e', fontSize: '12px' }}>STATUS:</span> <strong style={{ color: trace.has_error ? '#f78166' : '#3fb950', fontSize: '16px', marginLeft: '6px' }}>{trace.status_code || 200}</strong></div>
        </div>

        {/* Gantt Timeline Chart Breakdown */}
        <h4 style={{ color: '#f0f6fc', margin: '0 0 12px 0', fontSize: '14px' }}>Timeline Execution Graph (Gantt View)</h4>
        <div style={{ background: '#0d1117', border: '1px solid #30363d', borderRadius: '8px', padding: '16px', display: 'flex', flexDirection: 'column', gap: '14px', maxHeight: '320px', overflowY: 'auto' }}>
          {spans.map((span, idx) => {
            const startPct = Math.min(90, ((span.start_offset || (idx * 20)) / totalDuration) * 100);
            const widthPct = Math.max(5, (span.duration_ms / totalDuration) * 100);
            const color = getServiceColor(span.service);

            return (
              <div key={idx} style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '12px' }}>
                  <span style={{ fontWeight: 600, color: '#f0f6fc' }}>
                    <span style={{ color, fontWeight: 700, marginRight: '6px' }}>[{span.service.toUpperCase()}]</span>
                    {span.name}
                  </span>
                  <span style={{ fontWeight: 700, color }}>{span.duration_ms} ms</span>
                </div>

                {/* Timeline Bar Track */}
                <div style={{ background: '#161b22', height: '10px', borderRadius: '5px', width: '100%', position: 'relative' }}>
                  <div style={{
                    position: 'absolute',
                    left: `${startPct}%`,
                    width: `${widthPct}%`,
                    height: '100%',
                    backgroundColor: color,
                    borderRadius: '5px',
                    boxShadow: `0 0 8px ${color}66`
                  }} />
                </div>
              </div>
            );
          })}
        </div>

        <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: '20px' }}>
          <button onClick={onClose} style={{ background: '#238636', color: '#fff', border: 'none', padding: '8px 16px', borderRadius: '6px', cursor: 'pointer', fontWeight: 600, fontSize: '13px' }}>Close</button>
        </div>
      </div>
    </div>
  );
}
