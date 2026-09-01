import React, { useState, useEffect, useCallback } from 'react';
import { Activity, RefreshCw, Zap, Clock, AlertTriangle, CheckCircle, ShieldAlert, Cpu } from 'lucide-react';
import { apiGet } from '../../api/client.js';
import LatencyChart from './LatencyChart.jsx';
import ThroughputChart from './ThroughputChart.jsx';
import EndpointTable from './EndpointTable.jsx';
import ErrorRateCard from './ErrorRateCard.jsx';
import ServiceHealth from './ServiceHealth.jsx';

export default function Dashboard() {
  const [apmData, setApmData] = useState(null);
  const [loading, setLoading] = useState(false);

  const loadAPMData = useCallback(async () => {
    setLoading(true);
    try {
      const data = await apiGet('/apm');
      setApmData(data);
    } catch (err) {
      console.error('Failed to load APM telemetry', err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadAPMData();
    const interval = setInterval(loadAPMData, 10000); // Live poll every 10s
    return () => clearInterval(interval);
  }, [loadAPMData]);

  const overview = apmData?.overview || {
    api_health: 99.98,
    average_latency_ms: 58.4,
    requests_per_sec: 42.0,
    error_rate: 0.8,
    slow_api_count: 3,
  };

  return (
    <div style={{ minHeight: '100vh', backgroundColor: '#0d1117', color: '#c9d1d9', padding: '24px 32px' }}>
      {/* Header */}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '14px' }}>
          <Activity size={32} color="#3fb950" />
          <div>
            <h1 style={{ fontSize: '24px', fontWeight: 700, color: '#f0f6fc', margin: 0 }}>
              Application Performance Monitoring (APM)
            </h1>
            <span style={{ fontSize: '13px', color: '#8b949e' }}>
              Real-time API response times, RPS throughput, error rates, slow endpoint detection, &amp; AI insights
            </span>
          </div>
        </div>

        <button
          onClick={loadAPMData}
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
          {loading ? 'Fetching APM Telemetry...' : 'Refresh APM'}
        </button>
      </div>

      {/* KPI Cards Bar */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(180px, 1fr))', gap: '16px', marginBottom: '24px' }}>
        <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '10px', padding: '16px' }}>
          <div style={{ color: '#8b949e', fontSize: '11px', fontWeight: 700 }}>API HEALTH SCORE</div>
          <div style={{ fontSize: '24px', fontWeight: 800, color: '#3fb950', marginTop: '4px' }}>{overview.api_health}%</div>
        </div>
        <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '10px', padding: '16px' }}>
          <div style={{ color: '#8b949e', fontSize: '11px', fontWeight: 700 }}>AVERAGE LATENCY</div>
          <div style={{ fontSize: '24px', fontWeight: 800, color: '#58a6ff', marginTop: '4px' }}>{overview.average_latency_ms} ms</div>
        </div>
        <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '10px', padding: '16px' }}>
          <div style={{ color: '#8b949e', fontSize: '11px', fontWeight: 700 }}>REQUESTS / SEC (RPS)</div>
          <div style={{ fontSize: '24px', fontWeight: 800, color: '#a855f7', marginTop: '4px' }}>{overview.requests_per_sec}</div>
        </div>
        <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '10px', padding: '16px' }}>
          <div style={{ color: '#8b949e', fontSize: '11px', fontWeight: 700 }}>ERROR RATE</div>
          <div style={{ fontSize: '24px', fontWeight: 800, color: overview.error_rate > 1.0 ? '#f78166' : '#3fb950', marginTop: '4px' }}>{overview.error_rate}%</div>
        </div>
        <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '10px', padding: '16px' }}>
          <div style={{ color: '#8b949e', fontSize: '11px', fontWeight: 700 }}>FLAGGED SLOW APIS</div>
          <div style={{ fontSize: '24px', fontWeight: 800, color: '#ffa657', marginTop: '4px' }}>{overview.slow_api_count}</div>
        </div>
      </div>

      {/* Latency and Throughput Charts Grid */}
      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '20px', marginBottom: '24px' }}>
        <LatencyChart avgLatency={overview.average_latency_ms} />
        <ThroughputChart rps={overview.requests_per_sec} />
      </div>

      {/* Endpoint Table & Slow APIs */}
      <div style={{ marginBottom: '24px' }}>
        <EndpointTable
          endpoints={apmData?.endpoints}
          slowEndpoints={apmData?.slow_endpoints}
          topAPIs={apmData?.top_apis}
        />
      </div>

      {/* Microservice Health & AI Insights Grid */}
      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '20px' }}>
        <ServiceHealth services={apmData?.services} />
        <ErrorRateCard errorRate={overview.error_rate} recommendations={apmData?.recommendations} />
      </div>
    </div>
  );
}
