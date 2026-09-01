import React, { useState, useEffect, useCallback } from 'react';
import { GitCommit, Search, RefreshCw, Layers, AlertTriangle, Network, ShieldCheck } from 'lucide-react';
import { apiGet } from '../../api/client.js';
import TraceList from './TraceList.jsx';
import TraceDetail from './TraceDetail.jsx';
import ServiceMap from './ServiceMap.jsx';

export default function TracingSuite() {
  const [traces, setTraces] = useState([]);
  const [topology, setTopology] = useState(null);
  const [loading, setLoading] = useState(false);
  const [activeTab, setActiveTab] = useState('traces'); // 'traces' | 'servicemap'

  const [searchQuery, setSearchQuery] = useState('');
  const [onlySlow, setOnlySlow] = useState(false);
  const [selectedTrace, setSelectedTrace] = useState(null);

  const loadTracingData = useCallback(async () => {
    setLoading(true);
    try {
      let url = `/traces?only_slow=${onlySlow}`;
      if (searchQuery) url += `&query=${encodeURIComponent(searchQuery)}`;

      const traceRes = await apiGet(url);
      setTraces(traceRes.traces || []);

      const topoRes = await apiGet('/traces/services');
      setTopology(topoRes);
    } catch (err) {
      console.error('Failed to load distributed traces', err);
    } finally {
      setLoading(false);
    }
  }, [onlySlow, searchQuery]);

  useEffect(() => {
    loadTracingData();
  }, [loadTracingData]);

  const handleSearchSubmit = (e) => {
    e.preventDefault();
    loadTracingData();
  };

  const slowCount = traces.filter(t => t.is_slow).length;
  const errorCount = traces.filter(t => t.has_error).length;

  return (
    <div style={{ minHeight: '100vh', backgroundColor: '#0d1117', color: '#c9d1d9', padding: '24px 32px' }}>
      {/* Header */}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '14px' }}>
          <GitCommit size={32} color="#a855f7" />
          <div>
            <h1 style={{ fontSize: '24px', fontWeight: 700, color: '#f0f6fc', margin: 0 }}>
              Distributed Tracing &amp; Service Dependency Map
            </h1>
            <span style={{ fontSize: '13px', color: '#8b949e' }}>
              End-to-end request tracing with Trace IDs, Gantt timeline breakdown, and slow request detection (&gt;500ms)
            </span>
          </div>
        </div>

        <button
          onClick={loadTracingData}
          disabled={loading}
          style={{
            backgroundColor: '#21262d',
            color: '#c9d1d9',
            border: '1px solid #30363d',
            padding: '8px 16px',
            borderRadius: '6px',
            cursor: 'pointer',
            fontWeight: 600,
            fontSize: '13px',
            display: 'flex',
            alignItems: 'center',
            gap: '6px'
          }}
        >
          <RefreshCw size={16} className={loading ? 'animate-spin' : ''} />
          {loading ? 'Fetching Traces...' : 'Refresh Traces'}
        </button>
      </div>

      {/* Navigation Tabs */}
      <div style={{ display: 'flex', gap: '12px', marginBottom: '24px', borderBottom: '1px solid #30363d', paddingBottom: '12px' }}>
        <button
          onClick={() => setActiveTab('traces')}
          style={{
            background: activeTab === 'traces' ? '#21262d' : 'transparent',
            color: activeTab === 'traces' ? '#58a6ff' : '#8b949e',
            border: activeTab === 'traces' ? '1px solid #30363d' : 'none',
            padding: '8px 16px',
            borderRadius: '6px',
            cursor: 'pointer',
            fontWeight: 600,
            fontSize: '13px',
            display: 'flex',
            alignItems: 'center',
            gap: '6px'
          }}
        >
          <GitCommit size={16} /> Distributed Traces List
        </button>
        <button
          onClick={() => setActiveTab('servicemap')}
          style={{
            background: activeTab === 'servicemap' ? '#21262d' : 'transparent',
            color: activeTab === 'servicemap' ? '#58a6ff' : '#8b949e',
            border: activeTab === 'servicemap' ? '1px solid #30363d' : 'none',
            padding: '8px 16px',
            borderRadius: '6px',
            cursor: 'pointer',
            fontWeight: 600,
            fontSize: '13px',
            display: 'flex',
            alignItems: 'center',
            gap: '6px'
          }}
        >
          <Network size={16} /> Service Dependency Map
        </button>
      </div>

      {/* KPI Cards */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(180px, 1fr))', gap: '16px', marginBottom: '24px' }}>
        <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '10px', padding: '16px' }}>
          <div style={{ color: '#8b949e', fontSize: '11px', fontWeight: 700 }}>RECORDED TRACES</div>
          <div style={{ fontSize: '24px', fontWeight: 800, color: '#f0f6fc', marginTop: '4px' }}>{traces.length}</div>
        </div>
        <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '10px', padding: '16px' }}>
          <div style={{ color: '#8b949e', fontSize: '11px', fontWeight: 700 }}>SLOW REQUESTS (&gt;500ms)</div>
          <div style={{ fontSize: '24px', fontWeight: 800, color: '#ffa657', marginTop: '4px' }}>{slowCount}</div>
        </div>
        <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '10px', padding: '16px' }}>
          <div style={{ color: '#8b949e', fontSize: '11px', fontWeight: 700 }}>FAILED REQUESTS</div>
          <div style={{ fontSize: '24px', fontWeight: 800, color: '#f78166', marginTop: '4px' }}>{errorCount}</div>
        </div>
        <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '10px', padding: '16px' }}>
          <div style={{ color: '#8b949e', fontSize: '11px', fontWeight: 700 }}>AVG LATENCY</div>
          <div style={{ fontSize: '24px', fontWeight: 800, color: '#3fb950', marginTop: '4px' }}>
            {traces.length > 0 ? Math.round(traces.reduce((acc, t) => acc + t.duration_ms, 0) / traces.length) : 0} ms
          </div>
        </div>
      </div>

      {activeTab === 'traces' ? (
        <>
          {/* Search Controls */}
          <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '12px', padding: '20px', marginBottom: '24px' }}>
            <form onSubmit={handleSearchSubmit} style={{ display: 'grid', gridTemplateColumns: '1fr 180px 100px', gap: '14px', alignItems: 'center' }}>
              <div style={{ position: 'relative' }}>
                <Search size={16} color="#8b949e" style={{ position: 'absolute', left: '12px', top: '10px' }} />
                <input
                  type="text"
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  placeholder="Filter traces by route or trace ID (e.g. /aiops, /reports)..."
                  style={{
                    width: '100%',
                    padding: '8px 12px 8px 36px',
                    background: '#0d1117',
                    border: '1px solid #30363d',
                    color: '#fff',
                    borderRadius: '6px',
                    fontSize: '13px'
                  }}
                />
              </div>

              <div style={{ display: 'flex', alignItems: 'center', gap: '8px', fontSize: '13px', color: '#c9d1d9' }}>
                <input
                  type="checkbox"
                  id="onlySlow"
                  checked={onlySlow}
                  onChange={(e) => setOnlySlow(e.target.checked)}
                />
                <label htmlFor="onlySlow" style={{ cursor: 'pointer' }}>Only Slow (&gt;500ms)</label>
              </div>

              <button
                type="submit"
                style={{
                  backgroundColor: '#1f6feb',
                  color: '#fff',
                  border: 'none',
                  padding: '8px 16px',
                  borderRadius: '6px',
                  cursor: 'pointer',
                  fontWeight: 600,
                  fontSize: '13px'
                }}
              >
                Search
              </button>
            </form>
          </div>

          {/* Trace List Table */}
          <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '12px', padding: '24px' }}>
            <TraceList traces={traces} onSelectTrace={setSelectedTrace} />
          </div>
        </>
      ) : (
        <ServiceMap topology={topology} />
      )}

      {/* Trace Timeline Modal */}
      {selectedTrace && <TraceDetail trace={selectedTrace} onClose={() => setSelectedTrace(null)} />}
    </div>
  );
}
