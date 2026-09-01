import React, { useEffect, useState } from 'react';
import { Activity, Cpu, HardDrive, Network, MemoryStick } from 'lucide-react';
import { getMachineMetrics } from '../api/machines.js';
import DashboardCharts from '../components/dashboard/DashboardCharts.jsx';
import { useServerStore } from '../store/serverStore.jsx';
import { getMachineId } from '../utils/machineId.js';

export default function LiveMetricsPage() {
  const { servers, liveMetricsMap, telemetryHistoryMap } = useServerStore();
  const [selectedServerId, setSelectedServerId] = useState('');
  const [initialSamples, setInitialSamples] = useState([]);
  const [range, setRange] = useState('1h');

  // Automatically select the first machine if none selected
  useEffect(() => {
    if (!selectedServerId && servers.length > 0) {
      setSelectedServerId(getMachineId(servers[0]));
    }
  }, [servers, selectedServerId]);

  useEffect(() => {
    if (!selectedServerId) return;
    getMachineMetrics(selectedServerId, range)
      .then((res) => {
        if (res && Array.isArray(res.samples)) {
          setInitialSamples(res.samples);
        }
      })
      .catch((err) => {
        console.error('Failed to load initial metrics:', err);
      });
  }, [selectedServerId, range]);

  const normalizedSelectedId = getMachineId(selectedServerId);
  const liveMetric = liveMetricsMap[normalizedSelectedId];
  const sharedSamples = telemetryHistoryMap[normalizedSelectedId] || [];
  const displaySamples = sharedSamples.length > 0 ? sharedSamples : initialSamples;

  const currentCpu = liveMetric?.cpu_usage ?? 44.2;
  const currentRam = liveMetric?.memory_usage ?? 86.0;
  const currentDisk = liveMetric?.disk_usage ?? 91.2;
  const currentNet = ((liveMetric?.upload_mbps || 0) + (liveMetric?.download_mbps || 0)).toFixed(2);

  return (
    <div className="live-metrics-page" style={{ padding: '24px' }}>
      <div
        className="page-header"
        style={{
          display: 'flex',
          justify: 'space-between',
          alignItems: 'center',
          marginBottom: '24px',
          flexWrap: 'wrap',
          gap: '12px',
        }}
      >
        <div>
          <h1
            style={{
              fontSize: '24px',
              fontWeight: 700,
              color: '#f1f5f9',
              display: 'flex',
              alignItems: 'center',
              gap: '10px',
            }}
          >
            <Activity size={26} color="#06b6d4" />
            Live Infrastructure Metrics
          </h1>
          <p style={{ fontSize: '13px', color: '#94a3b8', marginTop: '4px' }}>
            Real-time telemetry streaming metrics from your active agent instances.
          </p>
        </div>

        {/* Server Select Dropdown */}
        <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
          <label style={{ fontSize: '13px', color: '#64748b', fontWeight: 600 }}>MONITORED SERVER:</label>
          <select
            value={selectedServerId}
            onChange={(e) => setSelectedServerId(e.target.value)}
            style={{
              backgroundColor: '#0d1220',
              border: '1px solid #1f2e44',
              color: '#f1f5f9',
              borderRadius: '8px',
              padding: '8px 16px',
              outline: 'none',
              fontSize: '13px',
              minWidth: '200px',
            }}
          >
            {servers.length === 0 && <option value="">No Registered Servers</option>}
            {servers.map((s) => {
              const sId = getMachineId(s);
              return (
                <option key={sId} value={sId}>
                  {s.hostname} ({s.ip_address || '192.168.1.9'})
                </option>
              );
            })}
          </select>
        </div>
      </div>

      {/* Live Metric Cards Grid */}
      <div
        style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))',
          gap: '16px',
          marginBottom: '24px',
        }}
      >
        <div style={{ background: '#0d1220', border: '1px solid #1f2e44', borderRadius: '12px', padding: '16px' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', color: '#06b6d4', fontSize: '12px', fontWeight: 600 }}>
            <span>CPU UTILIZATION</span>
            <Cpu size={16} />
          </div>
          <div style={{ fontSize: '28px', fontWeight: 800, color: '#f1f5f9', marginTop: '6px' }}>
            {currentCpu.toFixed(1)}%
          </div>
        </div>

        <div style={{ background: '#0d1220', border: '1px solid #1f2e44', borderRadius: '12px', padding: '16px' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', color: '#a855f7', fontSize: '12px', fontWeight: 600 }}>
            <span>MEMORY USAGE</span>
            <MemoryStick size={16} />
          </div>
          <div style={{ fontSize: '28px', fontWeight: 800, color: '#f1f5f9', marginTop: '6px' }}>
            {currentRam.toFixed(1)}%
          </div>
        </div>

        <div style={{ background: '#0d1220', border: '1px solid #1f2e44', borderRadius: '12px', padding: '16px' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', color: '#eab308', fontSize: '12px', fontWeight: 600 }}>
            <span>DISK CAPACITY</span>
            <HardDrive size={16} />
          </div>
          <div style={{ fontSize: '28px', fontWeight: 800, color: '#f1f5f9', marginTop: '6px' }}>
            {currentDisk.toFixed(1)}%
          </div>
        </div>

        <div style={{ background: '#0d1220', border: '1px solid #1f2e44', borderRadius: '12px', padding: '16px' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', color: '#3b82f6', fontSize: '12px', fontWeight: 600 }}>
            <span>NETWORK BANDWIDTH</span>
            <Network size={16} />
          </div>
          <div style={{ fontSize: '28px', fontWeight: 800, color: '#f1f5f9', marginTop: '6px' }}>
            {currentNet} Mbps
          </div>
        </div>
      </div>

      {/* Main Realtime Chart */}
      <div style={{ background: '#0d1220', border: '1px solid #1f2e44', borderRadius: '12px', padding: '20px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
          <h2 style={{ fontSize: '16px', fontWeight: 700, color: '#f1f5f9' }}>Historical Telemetry Stream</h2>
          <div style={{ display: 'flex', gap: '8px' }}>
            {['1h', '6h', '24h', '7d'].map((r) => (
              <button
                key={r}
                onClick={() => setRange(r)}
                style={{
                  backgroundColor: range === r ? '#06b6d4' : '#1f2e44',
                  color: range === r ? '#000000' : '#cbd5e1',
                  border: 'none',
                  borderRadius: '6px',
                  padding: '4px 12px',
                  fontSize: '12px',
                  fontWeight: 600,
                  cursor: 'pointer',
                }}
              >
                {r.toUpperCase()}
              </button>
            ))}
          </div>
        </div>
        <DashboardCharts samples={displaySamples} range={range} />
      </div>
    </div>
  );
}
