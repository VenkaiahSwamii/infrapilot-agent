import React, { useState, useMemo } from 'react';
import { Server, Search, RefreshCw } from 'lucide-react';
import { useServerStore } from '../store/serverStore.jsx';
import { useNavigate } from 'react-router-dom';
import { getMachineId } from '../utils/machineId.js';

export default function InfrastructurePage() {
  const { servers, loading, fetchServers, liveMetricsMap } = useServerStore();
  const [filterType, setFilterType] = useState('all');
  const [searchQuery, setSearchQuery] = useState('');
  const navigate = useNavigate();

  const stats = useMemo(() => {
    const total = servers.length;
    const online = servers.filter((s) => String(s.status || '').toUpperCase() === 'ONLINE').length;
    const offline = total - online;
    const warning = servers.filter((s) => {
      const live = liveMetricsMap[getMachineId(s)] || {};
      const cpu = Number(live.cpu_usage ?? s.cpu_usage ?? 0);
      const ram = Number(live.memory_usage ?? s.memory_usage ?? 0);
      return cpu > 80 || ram > 85;
    }).length;

    return { total, online, offline, warning };
  }, [servers, liveMetricsMap]);

  const filteredServers = useMemo(() => {
    return servers.filter((s) => {
      const sId = getMachineId(s);
      const matchesSearch =
        String(s.hostname || '').toLowerCase().includes(searchQuery.toLowerCase()) ||
        String(s.ip_address || '').includes(searchQuery) ||
        String(s.os || '').toLowerCase().includes(searchQuery.toLowerCase()) ||
        sId.includes(searchQuery.toLowerCase());

      const live = liveMetricsMap[sId] || {};
      const statusStr = String(s.status || '').toUpperCase();
      const cpu = Number(live.cpu_usage ?? s.cpu_usage ?? 0);
      const ram = Number(live.memory_usage ?? s.memory_usage ?? 0);
      const hasWarning = cpu > 80 || ram > 85;

      if (filterType === 'online') return matchesSearch && statusStr === 'ONLINE' && !hasWarning;
      if (filterType === 'warning') return matchesSearch && statusStr === 'ONLINE' && hasWarning;
      if (filterType === 'offline') return matchesSearch && statusStr !== 'ONLINE';
      return matchesSearch;
    });
  }, [servers, searchQuery, filterType, liveMetricsMap]);

  return (
    <div className="infrastructure-page" style={{ padding: '24px' }}>
      <div
        className="page-header"
        style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px' }}
      >
        <div>
          <h1 style={{ fontSize: '24px', fontWeight: 700, color: '#f1f5f9', display: 'flex', alignItems: 'center', gap: '10px' }}>
            <Server size={26} color="#06b6d4" />
            Infrastructure Inventory
          </h1>
          <p style={{ fontSize: '13px', color: '#94a3b8', marginTop: '4px' }}>
            Operational health and configurations of all enrolled server endpoints.
          </p>
        </div>
        <button
          className="refresh-btn"
          onClick={fetchServers}
          disabled={loading}
          style={{
            display: 'flex',
            alignItems: 'center',
            gap: '6px',
            padding: '8px 16px',
            backgroundColor: '#1f2e44',
            color: '#f1f5f9',
            border: '1px solid #2e3f5a',
            borderRadius: '8px',
            cursor: 'pointer',
            fontSize: '13px',
            fontWeight: 500,
          }}
        >
          <RefreshCw size={14} className={loading ? 'spin' : ''} />
          {loading ? 'Refreshing...' : 'Refresh'}
        </button>
      </div>

      {/* KPI Cards */}
      <div
        className="infra-kpis"
        style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(180px, 1fr))', gap: '16px', marginBottom: '24px' }}
      >
        <button
          className={`kpi-card ${filterType === 'all' ? 'active' : ''}`}
          onClick={() => setFilterType('all')}
          style={{
            background: '#0d1220',
            border: filterType === 'all' ? '1px solid #06b6d4' : '1px solid #1f2e44',
            borderRadius: '12px',
            padding: '16px',
            color: '#f1f5f9',
            cursor: 'pointer',
            textAlign: 'left',
          }}
        >
          <span style={{ fontSize: '11px', fontWeight: 700, color: '#64748b', textTransform: 'uppercase' }}>TOTAL SERVERS</span>
          <strong style={{ display: 'block', fontSize: '24px', fontWeight: 800, marginTop: '4px' }}>{stats.total}</strong>
        </button>
        <button
          className={`kpi-card ${filterType === 'online' ? 'active' : ''}`}
          onClick={() => setFilterType('online')}
          style={{
            background: '#0d1220',
            border: filterType === 'online' ? '1px solid #22c55e' : '1px solid #1f2e44',
            borderRadius: '12px',
            padding: '16px',
            color: '#f1f5f9',
            cursor: 'pointer',
            textAlign: 'left',
          }}
        >
          <span style={{ fontSize: '11px', fontWeight: 700, color: '#22c55e', textTransform: 'uppercase' }}>ONLINE</span>
          <strong style={{ display: 'block', fontSize: '24px', fontWeight: 800, marginTop: '4px' }}>{stats.online}</strong>
        </button>
        <button
          className={`kpi-card ${filterType === 'warning' ? 'active' : ''}`}
          onClick={() => setFilterType('warning')}
          style={{
            background: '#0d1220',
            border: filterType === 'warning' ? '1px solid #f59e0b' : '1px solid #1f2e44',
            borderRadius: '12px',
            padding: '16px',
            color: '#f1f5f9',
            cursor: 'pointer',
            textAlign: 'left',
          }}
        >
          <span style={{ fontSize: '11px', fontWeight: 700, color: '#f59e0b', textTransform: 'uppercase' }}>WARNING</span>
          <strong style={{ display: 'block', fontSize: '24px', fontWeight: 800, marginTop: '4px' }}>{stats.warning}</strong>
        </button>
        <button
          className={`kpi-card ${filterType === 'offline' ? 'active' : ''}`}
          onClick={() => setFilterType('offline')}
          style={{
            background: '#0d1220',
            border: filterType === 'offline' ? '1px solid #ef4444' : '1px solid #1f2e44',
            borderRadius: '12px',
            padding: '16px',
            color: '#f1f5f9',
            cursor: 'pointer',
            textAlign: 'left',
          }}
        >
          <span style={{ fontSize: '11px', fontWeight: 700, color: '#ef4444', textTransform: 'uppercase' }}>OFFLINE</span>
          <strong style={{ display: 'block', fontSize: '24px', fontWeight: 800, marginTop: '4px' }}>{stats.offline}</strong>
        </button>
      </div>

      {/* Toolbar */}
      <div className="infra-toolbar" style={{ display: 'flex', gap: '16px', marginBottom: '20px' }}>
        <div
          className="search-box"
          style={{
            display: 'flex',
            alignItems: 'center',
            gap: '8px',
            backgroundColor: '#0d1220',
            border: '1px solid #1f2e44',
            borderRadius: '8px',
            padding: '8px 12px',
            flex: 1,
          }}
        >
          <Search size={16} color="#64748b" />
          <input
            type="text"
            placeholder="Search servers by hostname, IP, OS..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            style={{ background: 'none', border: 'none', color: '#f1f5f9', outline: 'none', width: '100%', fontSize: '13px' }}
          />
        </div>
      </div>

      {/* Servers Table */}
      <div
        className="table-wrapper"
        style={{ backgroundColor: '#0d1220', border: '1px solid #1f2e44', borderRadius: '12px', overflow: 'hidden' }}
      >
        {filteredServers.length === 0 ? (
          <div className="empty-state" style={{ textAlign: 'center', padding: '60px 24px', color: '#64748b' }}>
            <Server size={40} style={{ opacity: 0.3, marginBottom: '12px' }} />
            <strong>No matching servers found</strong>
            <p style={{ fontSize: '12px', marginTop: '4px' }}>Try relaxing your search query or filters.</p>
          </div>
        ) : (
          <table style={{ width: '100%', borderCollapse: 'collapse', textAlign: 'left', fontSize: '13px' }}>
            <thead>
              <tr style={{ backgroundColor: '#080c14', borderBottom: '1px solid #1f2e44', color: '#64748b' }}>
                <th style={{ padding: '14px 16px', fontWeight: 600 }}>Server</th>
                <th style={{ padding: '14px 16px', fontWeight: 600 }}>IP Address</th>
                <th style={{ padding: '14px 16px', fontWeight: 600 }}>OS / Kernel</th>
                <th style={{ padding: '14px 16px', fontWeight: 600 }}>Status</th>
                <th style={{ padding: '14px 16px', fontWeight: 600 }}>CPU Usage</th>
                <th style={{ padding: '14px 16px', fontWeight: 600 }}>Memory Usage</th>
                <th style={{ padding: '14px 16px', fontWeight: 600 }}>Last Seen</th>
                <th style={{ padding: '14px 16px', fontWeight: 600 }}>Actions</th>
              </tr>
            </thead>
            <tbody>
              {filteredServers.map((s) => {
                const sId = getMachineId(s);
                const live = liveMetricsMap[sId] || {};
                const statusStr = String(s.status || '').toUpperCase();
                const isOnline = statusStr === 'ONLINE';

                const cpuVal = Number(live.cpu_usage ?? s.cpu_usage ?? 0);
                const ramVal = Number(live.memory_usage ?? s.memory_usage ?? 0);
                const isWarn = isOnline && (cpuVal > 80 || ramVal > 85);

                let badgeBg = 'rgba(239, 68, 68, 0.1)';
                let badgeColor = '#ef4444';
                let dotColor = '#ef4444';
                let statusLabel = 'Offline';

                if (isOnline) {
                  if (isWarn) {
                    badgeBg = 'rgba(245, 158, 11, 0.1)';
                    badgeColor = '#f59e0b';
                    dotColor = '#f59e0b';
                    statusLabel = 'Warning';
                  } else {
                    badgeBg = 'rgba(34, 197, 94, 0.1)';
                    badgeColor = '#22c55e';
                    dotColor = '#22c55e';
                    statusLabel = 'Online';
                  }
                }

                return (
                  <tr key={sId} style={{ borderBottom: '1px solid rgba(31, 46, 68, 0.5)', color: '#cbd5e1' }}>
                    <td style={{ padding: '14px 16px', fontWeight: 600, color: '#f1f5f9' }}>{s.hostname}</td>
                    <td style={{ padding: '14px 16px' }}>{s.ip_address || s.IPAddress || '--'}</td>
                    <td style={{ padding: '14px 16px', color: '#64748b' }}>
                      {s.os || 'Windows'} / {s.kernel || '--'}
                    </td>
                    <td style={{ padding: '14px 16px' }}>
                      <span
                        style={{
                          display: 'inline-flex',
                          alignItems: 'center',
                          gap: '6px',
                          padding: '3px 8px',
                          borderRadius: '12px',
                          fontSize: '11px',
                          fontWeight: 600,
                          backgroundColor: badgeBg,
                          color: badgeColor,
                        }}
                      >
                        <span style={{ width: 6, height: 6, borderRadius: '50%', backgroundColor: dotColor }} />
                        {statusLabel}
                      </span>
                    </td>
                    <td style={{ padding: '14px 16px' }}>{isOnline ? `${cpuVal.toFixed(1)}%` : '--'}</td>
                    <td style={{ padding: '14px 16px' }}>{isOnline ? `${ramVal.toFixed(1)}%` : '--'}</td>
                    <td style={{ padding: '14px 16px', color: '#64748b' }}>
                      {live.created_at ? new Date(live.created_at).toLocaleTimeString() : s.last_seen || 'never'}
                    </td>
                    <td style={{ padding: '14px 16px' }}>
                      <button
                        onClick={() => navigate(`/machines/${sId}`)}
                        style={{
                          background: 'none',
                          border: '1px solid #1f2e44',
                          color: '#06b6d4',
                          padding: '4px 10px',
                          borderRadius: '6px',
                          cursor: 'pointer',
                          fontSize: '12px',
                        }}
                      >
                        Details
                      </button>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}
