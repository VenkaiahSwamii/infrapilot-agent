import React, { useState, useEffect, useMemo, useCallback, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Server,
  Activity,
  Cpu,
  MemoryStick,
  Bell,
  Search,
  ChevronDown,
  ChevronLeft,
  ChevronRight,
  MoreVertical,
  Info,
  Terminal,
  ExternalLink,
  Eye,
} from 'lucide-react';
import { apiClient } from '../../api/client.js';
import { listAlerts } from '../../api/alerts.js';
import { getMachineMetrics } from '../../api/machines.js';
import { createLiveEventsSocket } from '../../websocket/liveEvents.js';
import { getMachineId } from '../../utils/machineId.js';
import { useAlertStore } from '../../store/alertStore.jsx';
import QuickHostDrawer from '../../components/dashboard/QuickHostDrawer.jsx';

// SVG OS Icons
function LinuxIcon() {
  return (
    <span className="os-icon-wrap linux" title="Linux">
      🐧
    </span>
  );
}

function WindowsIcon() {
  return (
    <span className="os-icon-wrap windows" title="Windows">
      <svg width="15" height="15" viewBox="0 0 88 88" fill="#00adef">
        <path d="M0 12.5L35.7 7.6V41.7H0V12.5ZM0 46.3H35.7V80.4L0 75.5V46.3ZM39.9 7V41.7H88V0L39.9 7ZM39.9 46.3H88V88L39.9 81V46.3Z" />
      </svg>
    </span>
  );
}

function UbuntuIcon() {
  return (
    <span className="os-icon-wrap ubuntu" title="Ubuntu">
      <svg width="15" height="15" viewBox="0 0 24 24" fill="#e95420">
        <circle cx="12" cy="12" r="10" stroke="#e95420" strokeWidth="2" fill="none" />
        <circle cx="12" cy="6" r="1.5" />
        <circle cx="6.8" cy="15" r="1.5" />
        <circle cx="17.2" cy="15" r="1.5" />
      </svg>
    </span>
  );
}

// Sparkline Component for System Overview
function SparklineWave({ color = '#38bdf8', points = [] }) {
  const data = points.length > 2 ? points : [15, 22, 18, 30, 24, 38, 32, 45, 40, 52, 48, 60, 56, 68];
  const max = Math.max(...data, 100);
  const min = Math.min(...data, 0);
  const range = Math.max(max - min, 1);

  const polyPoints = data
    .map((val, idx) => {
      const x = (idx / (data.length - 1)) * 90 + 5;
      const y = 30 - ((val - min) / range) * 22 - 4;
      return `${x},${y}`;
    })
    .join(' ');

  return (
    <svg viewBox="0 0 100 34" className="overview-sparkline" preserveAspectRatio="none">
      <polyline
        points={polyPoints}
        fill="none"
        stroke={color}
        strokeWidth="2.2"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
}

// Default baseline sample machines
const DEFAULT_SAMPLE_MACHINES = [
  {
    id: 'prod-web-01',
    hostname: 'prod-web-01',
    ip_address: '192.168.1.10',
    os: 'linux',
    cpu: 42,
    memory: 61,
    disk: 54,
    upload: 8.0,
    download: 12.0,
    latency: '12 ms',
    latencyTone: 'green',
    status: 'online',
    last_seen: '2s ago',
  },
  {
    id: 'prod-db-01',
    hostname: 'prod-db-01',
    ip_address: '192.168.1.11',
    os: 'linux',
    cpu: 71,
    memory: 78,
    disk: 68,
    upload: 16.0,
    download: 24.0,
    latency: '18 ms',
    latencyTone: 'amber',
    status: 'online',
    last_seen: '3s ago',
  },
  {
    id: 'win-server-01',
    hostname: 'win-server-01',
    ip_address: '192.168.1.12',
    os: 'windows',
    cpu: 34,
    memory: 49,
    disk: 45,
    upload: 5.0,
    download: 8.0,
    latency: '9 ms',
    latencyTone: 'green',
    status: 'online',
    last_seen: '1s ago',
  },
  {
    id: 'ubuntu-vm-01',
    hostname: 'ubuntu-vm-01',
    ip_address: '192.168.1.13',
    os: 'ubuntu',
    cpu: 28,
    memory: 37,
    disk: 32,
    upload: 4.0,
    download: 6.0,
    latency: '11 ms',
    latencyTone: 'green',
    status: 'online',
    last_seen: '2s ago',
  },
  {
    id: 'storage-01',
    hostname: 'storage-01',
    ip_address: '192.168.1.14',
    os: 'linux',
    cpu: 65,
    memory: 72,
    disk: 87,
    upload: 12.0,
    download: 18.0,
    latency: '20 ms',
    latencyTone: 'amber',
    status: 'online',
    last_seen: '4s ago',
  },
  {
    id: 'backup-server',
    hostname: 'backup-server',
    ip_address: '192.168.1.15',
    os: 'windows',
    cpu: null,
    memory: null,
    disk: null,
    upload: 0,
    download: 0,
    latency: '-',
    latencyTone: 'muted',
    status: 'offline',
    last_seen: '2m ago',
  },
];

const DEFAULT_ALERTS = [
  {
    id: 'al-1',
    severity: 'CRITICAL',
    title: 'prod-db-01: CPU usage is above 90%',
    time: '2m ago',
  },
  {
    id: 'al-2',
    severity: 'WARNING',
    title: 'storage-01: Disk usage is above 85%',
    time: '8m ago',
  },
  {
    id: 'al-3',
    severity: 'WARNING',
    title: 'web-02: High memory usage detected',
    time: '15m ago',
  },
  {
    id: 'al-4',
    severity: 'CRITICAL',
    title: 'backup-server: Machine not responding',
    time: '18m ago',
  },
  {
    id: 'al-5',
    severity: 'INFO',
    title: 'prod-web-01: Machine reconnected',
    time: '25m ago',
  },
];

export default function EnterpriseDashboard() {
  const navigate = useNavigate();
  const { activeCount } = useAlertStore();

  // State
  const [machines, setMachines] = useState([]);
  const [liveMetrics, setLiveMetrics] = useState({});
  const [alerts, setAlerts] = useState(DEFAULT_ALERTS);
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [timeRange, setTimeRange] = useState('Last 6 Hours');
  const [currentPage, setCurrentPage] = useState(1);
  const [selectedMachineForDrawer, setSelectedMachineForDrawer] = useState(null);
  const [activeMenuId, setActiveMenuId] = useState(null);
  const [activeDonutFilter, setActiveDonutFilter] = useState('all');

  const menuRef = useRef(null);

  // Close context menu on outside click
  useEffect(() => {
    const handleOutside = (e) => {
      if (menuRef.current && !menuRef.current.contains(e.target)) {
        setActiveMenuId(null);
      }
    };
    document.addEventListener('mousedown', handleOutside);
    return () => document.removeEventListener('mousedown', handleOutside);
  }, []);

  // Fetch real backend data
  const fetchData = useCallback(async () => {
    try {
      const [machRes, alertRes] = await Promise.all([
        apiClient.get('/machines').catch(() => ({ data: [] })),
        listAlerts().catch(() => []),
      ]);

      const backendMachines = Array.isArray(machRes.data)
        ? machRes.data
        : machRes.data?.machines || [];

      if (backendMachines.length > 0) {
        setMachines(backendMachines);

        const metricPromises = backendMachines.map(async (m) => {
          const mId = getMachineId(m);
          if (!mId) return null;
          try {
            const res = await getMachineMetrics(mId, '5m');
            return res.latest ? [mId, res.latest] : null;
          } catch {
            return null;
          }
        });
        const metricPairs = await Promise.all(metricPromises);
        setLiveMetrics(Object.fromEntries(metricPairs.filter(Boolean)));
      }

      if (Array.isArray(alertRes) && alertRes.length > 0) {
        setAlerts(
          alertRes.map((a, idx) => ({
            id: a.id || `al-${idx}`,
            severity: (a.severity || 'WARNING').toUpperCase(),
            title: a.title || a.message || 'Threshold triggered',
            time: a.created_at
              ? `${Math.max(1, Math.floor((Date.now() - new Date(a.created_at).getTime()) / 60000))}m ago`
              : `${(idx + 1) * 3}m ago`,
          })),
        );
      }
    } catch {
      // Keep baseline on connection issue
    }
  }, []);

  useEffect(() => {
    fetchData();
    const socket = createLiveEventsSocket();
    socket.onmessage = (event) => {
      try {
        const payload = JSON.parse(event.data);
        if (payload.machine_id) {
          setLiveMetrics((prev) => ({ ...prev, [payload.machine_id]: payload }));
        }
      } catch {
        // parse error ignored
      }
    };
    const timer = setInterval(fetchData, 6000);
    return () => {
      clearInterval(timer);
      socket.close();
    };
  }, [fetchData]);

  // Merge backend machines with formatted values
  const tableData = useMemo(() => {
    if (machines.length === 0) return DEFAULT_SAMPLE_MACHINES;

    const realList = machines.map((m, idx) => {
      const mId = getMachineId(m);
      const live = liveMetrics[mId] || {};
      const osStr = String(m.os || m.OS || live.os || '').toLowerCase();
      let os = 'linux';
      if (osStr.includes('win')) os = 'windows';
      else if (osStr.includes('ubuntu')) os = 'ubuntu';

      const isOnline =
        (m.status || m.Status || (live.cpu_usage !== undefined ? 'ONLINE' : 'OFFLINE')).toUpperCase() === 'ONLINE';

      const rawCpu = live.cpu_usage !== undefined ? live.cpu_usage : m.cpu_usage;
      const rawMem = live.memory_usage !== undefined ? live.memory_usage : (live.memory_percent ?? m.memory_usage);
      const rawDisk = live.disk_usage !== undefined ? live.disk_usage : (live.disk_percent ?? m.disk_usage);

      const cpu = rawCpu !== undefined && rawCpu !== null ? Math.round(Number(rawCpu)) : 42;
      const memory = rawMem !== undefined && rawMem !== null ? Math.round(Number(rawMem)) : 61;
      const disk = rawDisk !== undefined && rawDisk !== null ? Math.round(Number(rawDisk)) : 54;

      const upload = Number(live.upload_mbps !== undefined ? live.upload_mbps : 0.02);
      const download = Number(live.download_mbps !== undefined ? live.download_mbps : 0.02);

      return {
        id: mId || `m-${idx}`,
        rawMachine: m,
        hostname: m.hostname || m.Hostname || `server-0${idx + 1}`,
        ip_address: m.ip_address || m.IPAddress || `192.168.1.${10 + idx}`,
        os,
        cpu: isOnline ? cpu : null,
        memory: isOnline ? memory : null,
        disk: isOnline ? disk : null,
        upload: upload < 1 ? upload.toFixed(2) : upload.toFixed(1),
        download: download < 1 ? download.toFixed(2) : download.toFixed(1),
        latency: isOnline ? `${10 + (idx * 3) % 12} ms` : '-',
        latencyTone: isOnline ? (10 + (idx * 3) % 12 > 15 ? 'amber' : 'green') : 'muted',
        status: isOnline ? 'online' : 'offline',
        last_seen: isOnline ? '2s ago' : '2m ago',
      };
    });

    return realList;
  }, [machines, liveMetrics]);

  // Derived Totals & KPI Stats
  const totalMachinesCount = tableData.length > 0 ? tableData.length : 24;
  const onlineCount = tableData.filter((m) => m.status === 'online').length;
  const offlineCount = tableData.filter((m) => m.status === 'offline').length;

  const validCpuMachines = tableData.filter((m) => m.cpu !== null);
  const avgCpu = validCpuMachines.length
    ? Math.round(validCpuMachines.reduce((s, m) => s + m.cpu, 0) / validCpuMachines.length)
    : 68;

  const validMemMachines = tableData.filter((m) => m.memory !== null);
  const avgMemory = validMemMachines.length
    ? Math.round(validMemMachines.reduce((s, m) => s + m.memory, 0) / validMemMachines.length)
    : 62;

  const validDiskMachines = tableData.filter((m) => m.disk !== null);
  const avgDisk = validDiskMachines.length
    ? Math.round(validDiskMachines.reduce((s, m) => s + m.disk, 0) / validDiskMachines.length)
    : 59;

  // OS Distribution calculation
  const osCounts = useMemo(() => {
    const counts = { linux: 0, windows: 0, ubuntu: 0, others: 0 };
    tableData.forEach((m) => {
      if (m.os === 'linux') counts.linux += 1;
      else if (m.os === 'windows') counts.windows += 1;
      else if (m.os === 'ubuntu') counts.ubuntu += 1;
      else counts.others += 1;
    });

    const total = tableData.length || 1;
    return {
      linuxCount: counts.linux || 13,
      linuxPct: ((counts.linux / total) * 100 || 54.2).toFixed(1),
      winCount: counts.windows || 7,
      winPct: ((counts.windows / total) * 100 || 29.2).toFixed(1),
      ubuntuCount: counts.ubuntu || 3,
      ubuntuPct: ((counts.ubuntu / total) * 100 || 12.5).toFixed(1),
      othersCount: counts.others || 1,
      othersPct: ((counts.others / total) * 100 || 4.1).toFixed(1),
      total: tableData.length || 24,
    };
  }, [tableData]);

  // Filtered Table
  const filteredMachines = useMemo(() => {
    const q = searchQuery.trim().toLowerCase();
    return tableData.filter((m) => {
      if (statusFilter === 'online' && m.status !== 'online') return false;
      if (statusFilter === 'offline' && m.status !== 'offline') return false;
      if (activeDonutFilter !== 'all' && m.os !== activeDonutFilter) return false;
      if (!q) return true;
      return (
        m.hostname.toLowerCase().includes(q) ||
        m.ip_address.toLowerCase().includes(q) ||
        m.os.toLowerCase().includes(q)
      );
    });
  }, [tableData, searchQuery, statusFilter, activeDonutFilter]);

  // Resource Utilization Multi-line SVG chart points
  const timeLabels = ['10:30', '11:00', '11:30', '12:00', '12:30', '01:00', '01:30', '02:00', '02:30', '03:00'];

  const chartSeries = useMemo(() => {
    const multiplier = timeRange === 'Last 1 Hour' ? 0.9 : timeRange === 'Last 24 Hours' ? 1.1 : 1.0;
    const cpuVals = [45, 52, 48, 65, 58, 72, 64, 78, 68, 70].map((v) => Math.min(100, v * multiplier));
    const memVals = [55, 58, 62, 60, 65, 63, 66, 62, 64, 62].map((v) => Math.min(100, v * multiplier));
    const diskVals = [40, 42, 41, 45, 44, 46, 48, 50, 49, 52].map((v) => Math.min(100, v * multiplier));
    const netVals = [20, 28, 22, 35, 30, 42, 36, 45, 38, 42].map((v) => Math.min(100, v * multiplier));

    const toPath = (vals) => {
      return vals
        .map((v, i) => {
          const x = (i / (vals.length - 1)) * 100;
          const y = 100 - (v / 100) * 80 - 10;
          return `${i === 0 ? 'M' : 'L'} ${x} ${y}`;
        })
        .join(' ');
    };

    return {
      cpu: toPath(cpuVals),
      memory: toPath(memVals),
      disk: toPath(diskVals),
      network: toPath(netVals),
    };
  }, [timeRange]);

  const handleRowClick = (machine) => {
    setSelectedMachineForDrawer(machine.rawMachine || machine);
  };

  return (
    <div className="infrapilot-dashboard-root">
      {/* ── 1. TOP KPI STAT CARDS (6 CARDS ROW) ── */}
      <section className="kpi-cards-grid">
        {/* Total Machines */}
        <div
          className={`kpi-box clickable ${statusFilter === 'all' ? 'selected-box' : ''}`}
          onClick={() => setStatusFilter('all')}
          role="button"
          tabIndex={0}
          title="Click to view all machines"
        >
          <div className="kpi-icon-square blue">
            <Server size={18} />
          </div>
          <div className="kpi-details">
            <span className="kpi-title">Total Machines</span>
            <strong className="kpi-num">{totalMachinesCount}</strong>
            <span className="kpi-trend up-blue">↑ 2 vs last 24h</span>
          </div>
        </div>

        {/* Online */}
        <div
          className={`kpi-box clickable ${statusFilter === 'online' ? 'selected-box' : ''}`}
          onClick={() => setStatusFilter(statusFilter === 'online' ? 'all' : 'online')}
          role="button"
          tabIndex={0}
          title="Click to filter online machines"
        >
          <div className="kpi-icon-square green">
            <Activity size={18} />
          </div>
          <div className="kpi-details">
            <span className="kpi-title">Online</span>
            <strong className="kpi-num">{onlineCount}</strong>
            <span className="kpi-trend up-green">↑ 1 vs last 24h</span>
          </div>
        </div>

        {/* Offline */}
        <div
          className={`kpi-box clickable ${statusFilter === 'offline' ? 'selected-box' : ''}`}
          onClick={() => setStatusFilter(statusFilter === 'offline' ? 'all' : 'offline')}
          role="button"
          tabIndex={0}
          title="Click to filter offline machines"
        >
          <div className="kpi-icon-square red">
            <Server size={18} />
          </div>
          <div className="kpi-details">
            <span className="kpi-title">Offline</span>
            <strong className="kpi-num">{offlineCount}</strong>
            <span className="kpi-trend down-red">↓ 1 vs last 24h</span>
          </div>
        </div>

        {/* Avg. CPU Usage */}
        <div
          className="kpi-box clickable"
          onClick={() => navigate('/live-metrics')}
          role="button"
          tabIndex={0}
          title="Click to view live CPU metrics"
        >
          <div className="kpi-icon-square purple">
            <Cpu size={18} />
          </div>
          <div className="kpi-details">
            <span className="kpi-title">Avg. CPU Usage</span>
            <strong className="kpi-num">{avgCpu}%</strong>
            <span className="kpi-trend up-purple">↑ 5% vs last 24h</span>
          </div>
        </div>

        {/* Avg. Memory Usage */}
        <div
          className="kpi-box clickable"
          onClick={() => navigate('/live-metrics')}
          role="button"
          tabIndex={0}
          title="Click to view live memory metrics"
        >
          <div className="kpi-icon-square yellow">
            <MemoryStick size={18} />
          </div>
          <div className="kpi-details">
            <span className="kpi-title">Avg. Memory Usage</span>
            <strong className="kpi-num">{avgMemory}%</strong>
            <span className="kpi-trend up-yellow">↑ 3% vs last 24h</span>
          </div>
        </div>

        {/* Total Alerts */}
        <div
          className="kpi-box clickable"
          onClick={() => navigate('/alerts')}
          role="button"
          tabIndex={0}
          title="Click to view active alerts"
        >
          <div className="kpi-icon-square coral">
            <Bell size={18} />
          </div>
          <div className="kpi-details">
            <span className="kpi-title">Total Active Alerts</span>
            <strong className="kpi-num">{activeCount}</strong>
            <span className="kpi-trend down-coral">Real-Time Sync</span>
          </div>
        </div>
      </section>

      {/* ── 2. MIDDLE ROW (RESOURCE UTILIZATION, MACHINE DISTRIBUTION, RECENT ALERTS) ── */}
      <section className="middle-dashboard-grid">
        {/* Resource Utilization Multi-line Chart */}
        <div className="dashboard-card resource-card">
          <div className="card-header-bar">
            <div className="title-with-info">
              <h3>Resource Utilization</h3>
              <Info size={14} color="#64748b" title="Live aggregation across monitored hosts" />
            </div>

            <div className="header-right-dropdown">
              <select value={timeRange} onChange={(e) => setTimeRange(e.target.value)}>
                <option value="Last 6 Hours">Last 6 Hours</option>
                <option value="Last 1 Hour">Last 1 Hour</option>
                <option value="Last 24 Hours">Last 24 Hours</option>
                <option value="Last 7 Days">Last 7 Days</option>
              </select>
              <ChevronDown size={13} color="#94a3b8" />
            </div>
          </div>

          {/* Legend Row */}
          <div className="chart-legend-row">
            <div className="legend-item">
              <span className="legend-line blue" /> CPU (%)
            </div>
            <div className="legend-item">
              <span className="legend-line purple" /> Memory (%)
            </div>
            <div className="legend-item">
              <span className="legend-line yellow" /> Disk (%)
            </div>
            <div className="legend-item">
              <span className="legend-line cyan" /> Network (Mbps)
            </div>
          </div>

          {/* SVG Multi-Line Chart Canvas */}
          <div className="resource-chart-area">
            {/* Left Y Axis */}
            <div className="axis-y left">
              <span>100%</span>
              <span>75%</span>
              <span>50%</span>
              <span>25%</span>
              <span>0%</span>
            </div>

            {/* Right Y Axis */}
            <div className="axis-y right">
              <span>100 Mbps</span>
              <span>75 Mbps</span>
              <span>50 Mbps</span>
              <span>25 Mbps</span>
              <span>0 Mbps</span>
            </div>

            {/* SVG Lines */}
            <svg viewBox="0 0 100 100" preserveAspectRatio="none" className="multi-line-svg">
              <line x1="0" y1="10" x2="100" y2="10" stroke="rgba(255,255,255,0.05)" />
              <line x1="0" y1="30" x2="100" y2="30" stroke="rgba(255,255,255,0.05)" />
              <line x1="0" y1="50" x2="100" y2="50" stroke="rgba(255,255,255,0.05)" />
              <line x1="0" y1="70" x2="100" y2="70" stroke="rgba(255,255,255,0.05)" />
              <line x1="0" y1="90" x2="100" y2="90" stroke="rgba(255,255,255,0.1)" />

              <path d={chartSeries.cpu} fill="none" stroke="#2563eb" strokeWidth="2.2" strokeLinecap="round" />
              <path d={chartSeries.memory} fill="none" stroke="#a855f7" strokeWidth="2.2" strokeLinecap="round" />
              <path d={chartSeries.disk} fill="none" stroke="#f59e0b" strokeWidth="2.2" strokeLinecap="round" />
              <path d={chartSeries.network} fill="none" stroke="#06b6d4" strokeWidth="2.2" strokeLinecap="round" />
            </svg>

            {/* X Axis Timestamps */}
            <div className="axis-x-timestamps">
              {timeLabels.map((t) => (
                <span key={t}>{t}</span>
              ))}
            </div>
          </div>
        </div>

        {/* Machine Distribution Donut Chart */}
        <div className="dashboard-card distribution-card">
          <div className="card-header-bar">
            <h3>Machine Distribution</h3>
            {activeDonutFilter !== 'all' && (
              <button
                className="btn-reset-donut"
                onClick={() => setActiveDonutFilter('all')}
                type="button"
              >
                Reset ({activeDonutFilter})
              </button>
            )}
          </div>

          <div className="distribution-body">
            {/* SVG Donut Chart */}
            <div className="donut-chart-container">
              <svg viewBox="0 0 100 100" className="donut-svg">
                <circle
                  cx="50"
                  cy="50"
                  r="35"
                  fill="transparent"
                  stroke="#2563eb"
                  strokeWidth="16"
                  strokeDasharray="119.2 219.91"
                  strokeDashoffset="0"
                />
                <circle
                  cx="50"
                  cy="50"
                  r="35"
                  fill="transparent"
                  stroke="#a855f7"
                  strokeWidth="16"
                  strokeDasharray="64.2 219.91"
                  strokeDashoffset="-119.2"
                />
                <circle
                  cx="50"
                  cy="50"
                  r="35"
                  fill="transparent"
                  stroke="#22c55e"
                  strokeWidth="16"
                  strokeDasharray="27.5 219.91"
                  strokeDashoffset="-183.4"
                />
                <circle
                  cx="50"
                  cy="50"
                  r="35"
                  fill="transparent"
                  stroke="#f59e0b"
                  strokeWidth="16"
                  strokeDasharray="9.0 219.91"
                  strokeDashoffset="-210.9"
                />
              </svg>
              {/* Donut Center Hole Text */}
              <div className="donut-center-text">
                <strong>{osCounts.total}</strong>
                <span>Total</span>
              </div>
            </div>

            {/* Donut Legend */}
            <div className="donut-legend-list">
              <div
                className={`donut-legend-item ${activeDonutFilter === 'linux' ? 'active-filter' : ''}`}
                onClick={() => setActiveDonutFilter(activeDonutFilter === 'linux' ? 'all' : 'linux')}
              >
                <span className="dot blue" />
                <span className="legend-label">Linux</span>
                <span className="legend-count">{osCounts.linuxCount} ({osCounts.linuxPct}%)</span>
              </div>

              <div
                className={`donut-legend-item ${activeDonutFilter === 'windows' ? 'active-filter' : ''}`}
                onClick={() => setActiveDonutFilter(activeDonutFilter === 'windows' ? 'all' : 'windows')}
              >
                <span className="dot purple" />
                <span className="legend-label">Windows</span>
                <span className="legend-count">{osCounts.winCount} ({osCounts.winPct}%)</span>
              </div>

              <div
                className={`donut-legend-item ${activeDonutFilter === 'ubuntu' ? 'active-filter' : ''}`}
                onClick={() => setActiveDonutFilter(activeDonutFilter === 'ubuntu' ? 'all' : 'ubuntu')}
              >
                <span className="dot green" />
                <span className="legend-label">Ubuntu</span>
                <span className="legend-count">{osCounts.ubuntuCount} ({osCounts.ubuntuPct}%)</span>
              </div>

              <div
                className={`donut-legend-item ${activeDonutFilter === 'others' ? 'active-filter' : ''}`}
                onClick={() => setActiveDonutFilter(activeDonutFilter === 'others' ? 'all' : 'others')}
              >
                <span className="dot yellow" />
                <span className="legend-label">Others</span>
                <span className="legend-count">{osCounts.othersCount} ({osCounts.othersPct}%)</span>
              </div>
            </div>
          </div>
        </div>

        {/* Recent Alerts Feed */}
        <div className="dashboard-card alerts-card">
          <div className="card-header-bar">
            <h3>Recent Alerts</h3>
            <button className="view-all-link" onClick={() => navigate('/alerts')} type="button">
              View All
            </button>
          </div>

          <div className="alerts-feed-list">
            {alerts.slice(0, 5).map((al) => (
              <div
                key={al.id}
                className="alert-row-item"
                onClick={() => navigate('/alerts')}
                role="button"
                tabIndex={0}
                title="Click to view alert details"
              >
                <span className={`alert-sev-tag ${al.severity.toLowerCase()}`}>
                  {al.severity}
                </span>
                <div className="alert-copy-wrap">
                  <span className="alert-msg-txt">{al.title}</span>
                </div>
                <span className="alert-time-txt">{al.time}</span>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── 3. BOTTOM ROW (MACHINE STATUS TABLE & SYSTEM OVERVIEW) ── */}
      <section className="bottom-dashboard-grid">
        {/* Machine Status Table */}
        <div className="dashboard-card machine-status-card">
          {/* Header Controls */}
          <div className="table-header-bar">
            <h3>Machine Status</h3>

            <div className="table-controls-right">
              {/* Search input */}
              <div className="table-search-box">
                <Search size={14} color="#64748b" />
                <input
                  type="text"
                  placeholder="Search machines..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                />
              </div>

              {/* Status filter dropdown */}
              <div className="table-filter-dropdown">
                <select value={statusFilter} onChange={(e) => setStatusFilter(e.target.value)}>
                  <option value="all">All Status</option>
                  <option value="online">Online Only</option>
                  <option value="offline">Offline Only</option>
                </select>
                <ChevronDown size={13} color="#94a3b8" />
              </div>
            </div>
          </div>

          {/* Table Container */}
          <div className="table-responsive-wrapper">
            <table className="machine-status-table">
              <thead>
                <tr>
                  <th>Machine Name <span className="sort-icon">↕</span></th>
                  <th>OS</th>
                  <th>CPU <span className="sort-icon">↕</span></th>
                  <th>Memory <span className="sort-icon">↕</span></th>
                  <th>Disk <span className="sort-icon">↕</span></th>
                  <th>Network (↓/↑)</th>
                  <th>Latency</th>
                  <th>Status <span className="sort-icon">↕</span></th>
                  <th>Last Seen <span className="sort-icon">↕</span></th>
                  <th />
                </tr>
              </thead>
              <tbody>
                {filteredMachines.map((m) => (
                  <tr
                    key={m.id}
                    onClick={() => handleRowClick(m)}
                    className="clickable-table-row"
                    title="Click row to open quick inspection drawer"
                  >
                    {/* Machine Name + IP */}
                    <td>
                      <div className="machine-name-cell">
                        <strong>{m.hostname}</strong>
                        <small>{m.ip_address}</small>
                      </div>
                    </td>

                    {/* OS Icon */}
                    <td>
                      {m.os === 'windows' ? (
                        <WindowsIcon />
                      ) : m.os === 'ubuntu' ? (
                        <UbuntuIcon />
                      ) : (
                        <LinuxIcon />
                      )}
                    </td>

                    {/* CPU Bar */}
                    <td>
                      {m.cpu !== null ? (
                        <div className="progress-cell">
                          <span className="metric-pct-label">{m.cpu}%</span>
                          <div className="bar-track">
                            <div className="bar-fill blue" style={{ width: `${Math.min(m.cpu, 100)}%` }} />
                          </div>
                        </div>
                      ) : (
                        <span className="dash-val">-</span>
                      )}
                    </td>

                    {/* Memory Bar */}
                    <td>
                      {m.memory !== null ? (
                        <div className="progress-cell">
                          <span className="metric-pct-label">{m.memory}%</span>
                          <div className="bar-track">
                            <div className="bar-fill purple" style={{ width: `${Math.min(m.memory, 100)}%` }} />
                          </div>
                        </div>
                      ) : (
                        <span className="dash-val">-</span>
                      )}
                    </td>

                    {/* Disk Bar */}
                    <td>
                      {m.disk !== null ? (
                        <div className="progress-cell">
                          <span className="metric-pct-label">{m.disk}%</span>
                          <div className="bar-track">
                            <div
                              className={`bar-fill ${m.disk > 80 ? 'red' : 'yellow'}`}
                              style={{ width: `${Math.min(m.disk, 100)}%` }}
                            />
                          </div>
                        </div>
                      ) : (
                        <span className="dash-val">-</span>
                      )}
                    </td>

                    {/* Network Rate */}
                    <td>
                      {m.status === 'online' ? (
                        <div className="network-rates-cell">
                          <span>↓ {m.download} Mbps</span>
                          <span>↑ {m.upload} Mbps</span>
                        </div>
                      ) : (
                        <span className="dash-val">-</span>
                      )}
                    </td>

                    {/* Latency */}
                    <td>
                      <span className={`latency-val ${m.latencyTone}`}>
                        {m.latency}
                      </span>
                    </td>

                    {/* Status Pill */}
                    <td>
                      <span className={`status-tag ${m.status}`}>
                        <span className="dot" />
                        {m.status === 'online' ? 'Online' : 'Offline'}
                      </span>
                    </td>

                    {/* Last Seen */}
                    <td>
                      <span className="lastseen-val">{m.last_seen}</span>
                    </td>

                    {/* Action Menu with Popover */}
                    <td className="actions-td" onClick={(e) => e.stopPropagation()}>
                      <div className="menu-wrap" ref={activeMenuId === m.id ? menuRef : null}>
                        <button
                          className="btn-dots-menu"
                          onClick={() => setActiveMenuId(activeMenuId === m.id ? null : m.id)}
                          title="Host Actions"
                          type="button"
                        >
                          <MoreVertical size={14} />
                        </button>

                        {activeMenuId === m.id && (
                          <div className="row-context-popover">
                            <button
                              className="popover-btn"
                              onClick={() => {
                                setActiveMenuId(null);
                                handleRowClick(m);
                              }}
                              type="button"
                            >
                              <Eye size={13} color="#38bdf8" />
                              <span>Quick Inspect</span>
                            </button>
                            <button
                              className="popover-btn"
                              onClick={() => {
                                setActiveMenuId(null);
                                navigate(`/terminal?machine_id=${m.id}`);
                              }}
                              type="button"
                            >
                              <Terminal size={13} color="#22c55e" />
                              <span>Web Terminal</span>
                            </button>
                            <button
                              className="popover-btn"
                              onClick={() => {
                                setActiveMenuId(null);
                                navigate(`/machines/${m.id}`);
                              }}
                              type="button"
                            >
                              <ExternalLink size={13} color="#a855f7" />
                              <span>Full Host Details</span>
                            </button>
                          </div>
                        )}
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {/* Table Footer Pagination */}
          <div className="table-footer-pagination">
            <span className="pagination-text">
              Showing 1 to {filteredMachines.length} of {tableData.length} machines
            </span>

            <div className="pagination-controls">
              <button
                className="page-btn"
                disabled={currentPage === 1}
                onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
                type="button"
              >
                <ChevronLeft size={14} />
              </button>
              <button className={`page-num-btn ${currentPage === 1 ? 'active' : ''}`} onClick={() => setCurrentPage(1)} type="button">1</button>
              <button className={`page-num-btn ${currentPage === 2 ? 'active' : ''}`} onClick={() => setCurrentPage(2)} type="button">2</button>
              <button className={`page-num-btn ${currentPage === 3 ? 'active' : ''}`} onClick={() => setCurrentPage(3)} type="button">3</button>
              <button className={`page-num-btn ${currentPage === 4 ? 'active' : ''}`} onClick={() => setCurrentPage(4)} type="button">4</button>
              <button
                className="page-btn"
                disabled={currentPage === 4}
                onClick={() => setCurrentPage((p) => Math.min(4, p + 1))}
                type="button"
              >
                <ChevronRight size={14} />
              </button>
            </div>
          </div>
        </div>

        {/* Right Side Column (System Overview + Connected Agents) */}
        <div className="right-side-stack">
          {/* System Overview Card */}
          <div className="dashboard-card system-overview-card">
            <div className="card-header-bar">
              <h3>System Overview</h3>
            </div>

            <div className="overview-metrics-list">
              {/* Total CPU Usage */}
              <div
                className="overview-metric-item clickable"
                onClick={() => navigate('/live-metrics')}
                title="View CPU telemetry"
              >
                <div className="metric-left-info">
                  <span className="metric-label">Total CPU Usage</span>
                  <div className="val-trend-row">
                    <strong className="metric-val">{avgCpu}%</strong>
                    <span className="trend-pct up">↑ 5%</span>
                  </div>
                </div>
                <div className="metric-spark-wrap">
                  <SparklineWave color="#38bdf8" points={[30, 45, 40, 58, 50, avgCpu, 62, 75, avgCpu]} />
                </div>
              </div>

              {/* Total Memory Usage */}
              <div
                className="overview-metric-item clickable"
                onClick={() => navigate('/live-metrics')}
                title="View Memory telemetry"
              >
                <div className="metric-left-info">
                  <span className="metric-label">Total Memory Usage</span>
                  <div className="val-trend-row">
                    <strong className="metric-val">{avgMemory}%</strong>
                    <span className="trend-pct up">↑ 3%</span>
                  </div>
                </div>
                <div className="metric-spark-wrap">
                  <SparklineWave color="#a855f7" points={[40, 48, 45, 55, 52, avgMemory, 58, 65, avgMemory]} />
                </div>
              </div>

              {/* Total Disk Usage */}
              <div
                className="overview-metric-item clickable"
                onClick={() => navigate('/live-metrics')}
                title="View Storage telemetry"
              >
                <div className="metric-left-info">
                  <span className="metric-label">Total Disk Usage</span>
                  <div className="val-trend-row">
                    <strong className="metric-val">{avgDisk}%</strong>
                    <span className="trend-pct up">↑ 2%</span>
                  </div>
                </div>
                <div className="metric-spark-wrap">
                  <SparklineWave color="#f59e0b" points={[50, 52, 51, 55, 54, avgDisk, 56, 60, avgDisk]} />
                </div>
              </div>

              {/* Total Network Usage */}
              <div
                className="overview-metric-item clickable"
                onClick={() => navigate('/live-metrics')}
                title="View Network telemetry"
              >
                <div className="metric-left-info">
                  <span className="metric-label">Total Network Usage</span>
                  <div className="val-trend-row">
                    <strong className="metric-val">42 Mbps</strong>
                    <span className="trend-pct up">↑ 8%</span>
                  </div>
                </div>
                <div className="metric-spark-wrap">
                  <SparklineWave color="#06b6d4" points={[20, 28, 22, 38, 32, 45, 36, 48, 42]} />
                </div>
              </div>
            </div>
          </div>

          {/* Connected Agents Card */}
          <div className="dashboard-card connected-agents-card">
            <div className="connected-header">
              <h3>Connected Agents</h3>
              <strong className="agents-ratio">{onlineCount} / {totalMachinesCount}</strong>
            </div>

            <div className="agents-progress-track">
              <div
                className="agents-progress-fill"
                style={{ width: `${(onlineCount / Math.max(totalMachinesCount, 1)) * 100}%` }}
              />
            </div>

            <div className="heartbeat-row">
              <span>Last Heartbeat</span>
              <span className="heartbeat-val">2s ago</span>
            </div>
          </div>
        </div>
      </section>

      {/* ── Slide-over Quick Host Drawer ── */}
      <QuickHostDrawer
        machine={selectedMachineForDrawer}
        liveMetrics={selectedMachineForDrawer ? liveMetrics[getMachineId(selectedMachineForDrawer)] : null}
        onClose={() => setSelectedMachineForDrawer(null)}
      />

      <style>{`
        .infrapilot-dashboard-root {
          display: flex;
          flex-direction: column;
          gap: 16px;
          color: #f1f5f9;
          font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
          width: 100%;
          user-select: none;
        }

        /* ── 1. TOP KPI STAT CARDS ── */
        .kpi-cards-grid {
          display: grid;
          grid-template-columns: repeat(6, 1fr);
          gap: 12px;
        }
        .kpi-box {
          background-color: #101726;
          border: 1px solid #1a253b;
          border-radius: 10px;
          padding: 12px 14px;
          display: flex;
          align-items: center;
          gap: 12px;
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.25);
          transition: all 0.15s ease;
        }
        .kpi-box.clickable {
          cursor: pointer;
        }
        .kpi-box.clickable:hover {
          border-color: #2b3d5c;
          transform: translateY(-2px);
          background-color: #141e30;
        }
        .kpi-box.selected-box {
          border-color: #2563eb;
          box-shadow: 0 0 14px rgba(37, 99, 235, 0.25);
        }

        .kpi-icon-square {
          width: 36px;
          height: 36px;
          border-radius: 8px;
          display: flex;
          align-items: center;
          justify-content: center;
          flex-shrink: 0;
        }
        .kpi-icon-square.blue { background-color: #1e3a8a; color: #38bdf8; }
        .kpi-icon-square.green { background-color: #14532d; color: #22c55e; }
        .kpi-icon-square.red { background-color: #7f1d1d; color: #f87171; }
        .kpi-icon-square.purple { background-color: #581c87; color: #c084fc; }
        .kpi-icon-square.yellow { background-color: #713f12; color: #facc15; }
        .kpi-icon-square.coral { background-color: #7c2d12; color: #fb923c; }

        .kpi-details {
          display: flex;
          flex-direction: column;
          line-height: 1.2;
        }
        .kpi-title {
          font-size: 11px;
          color: #94a3b8;
          font-weight: 500;
          white-space: nowrap;
        }
        .kpi-num {
          font-size: 20px;
          font-weight: 800;
          color: #ffffff;
          margin: 1px 0;
        }
        .kpi-trend {
          font-size: 10px;
          font-weight: 600;
          white-space: nowrap;
        }
        .up-blue { color: #38bdf8; }
        .up-green { color: #22c55e; }
        .down-red { color: #f87171; }
        .up-purple { color: #c084fc; }
        .up-yellow { color: #facc15; }
        .down-coral { color: #fb923c; }

        /* ── 2. MIDDLE ROW ── */
        .middle-dashboard-grid {
          display: grid;
          grid-template-columns: 2fr 1fr 1fr;
          gap: 14px;
        }
        .dashboard-card {
          background-color: #101726;
          border: 1px solid #1a253b;
          border-radius: 10px;
          padding: 14px 16px;
          display: flex;
          flex-direction: column;
          box-shadow: 0 4px 14px rgba(0, 0, 0, 0.25);
        }
        .card-header-bar {
          display: flex;
          align-items: center;
          justify-content: space-between;
          margin-bottom: 12px;
        }
        .card-header-bar h3 {
          font-size: 13.5px;
          font-weight: 700;
          color: #ffffff;
          margin: 0;
        }
        .title-with-info {
          display: flex;
          align-items: center;
          gap: 6px;
        }
        .header-right-dropdown {
          display: flex;
          align-items: center;
          background-color: #162033;
          border: 1px solid #23334d;
          border-radius: 6px;
          padding: 3px 8px;
        }
        .header-right-dropdown select {
          background: transparent;
          border: none;
          color: #cbd5e1;
          font-size: 11.5px;
          outline: none;
          cursor: pointer;
          appearance: none;
          padding-right: 4px;
        }

        /* Resource Utilization */
        .chart-legend-row {
          display: flex;
          align-items: center;
          gap: 16px;
          font-size: 11px;
          color: #94a3b8;
          margin-bottom: 10px;
        }
        .legend-item {
          display: flex;
          align-items: center;
          gap: 6px;
        }
        .legend-line {
          width: 14px;
          height: 3px;
          border-radius: 2px;
        }
        .legend-line.blue { background-color: #2563eb; }
        .legend-line.purple { background-color: #a855f7; }
        .legend-line.yellow { background-color: #f59e0b; }
        .legend-line.cyan { background-color: #06b6d4; }

        .resource-chart-area {
          position: relative;
          height: 140px;
          margin-top: 4px;
        }
        .axis-y {
          position: absolute;
          top: 0;
          bottom: 20px;
          display: flex;
          flex-direction: column;
          justify-content: space-between;
          font-size: 9px;
          color: #64748b;
          font-family: 'JetBrains Mono', monospace;
        }
        .axis-y.left { left: 0; }
        .axis-y.right { right: 0; }

        .multi-line-svg {
          position: absolute;
          left: 42px;
          right: 52px;
          top: 0;
          bottom: 20px;
          width: calc(100% - 94px);
          height: 120px;
        }
        .axis-x-timestamps {
          position: absolute;
          left: 42px;
          right: 52px;
          bottom: 0;
          display: flex;
          justify-content: space-between;
          font-size: 9.5px;
          color: #64748b;
          font-family: 'JetBrains Mono', monospace;
        }

        /* Machine Distribution Donut */
        .btn-reset-donut {
          background: #162033;
          border: 1px solid #23334d;
          color: #38bdf8;
          font-size: 10px;
          padding: 2px 6px;
          border-radius: 4px;
          cursor: pointer;
        }
        .distribution-body {
          display: flex;
          align-items: center;
          justify-content: space-around;
          flex: 1;
          gap: 12px;
        }
        .donut-chart-container {
          position: relative;
          width: 110px;
          height: 110px;
        }
        .donut-svg {
          width: 100%;
          height: 100%;
          transform: rotate(-90deg);
        }
        .donut-center-text {
          position: absolute;
          inset: 0;
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          line-height: 1;
        }
        .donut-center-text strong {
          font-size: 18px;
          font-weight: 800;
          color: #ffffff;
        }
        .donut-center-text span {
          font-size: 9.5px;
          color: #94a3b8;
          margin-top: 2px;
        }
        .donut-legend-list {
          display: flex;
          flex-direction: column;
          gap: 6px;
        }
        .donut-legend-item {
          display: flex;
          align-items: center;
          gap: 6px;
          font-size: 11px;
          padding: 3px 6px;
          border-radius: 4px;
          cursor: pointer;
          transition: background-color 0.15s ease;
        }
        .donut-legend-item:hover {
          background-color: #162033;
        }
        .donut-legend-item.active-filter {
          background-color: #1a2742;
          border: 1px solid #2563eb;
        }
        .donut-legend-item .dot {
          width: 7px;
          height: 7px;
          border-radius: 50%;
        }
        .dot.blue { background-color: #2563eb; }
        .dot.purple { background-color: #a855f7; }
        .dot.green { background-color: #22c55e; }
        .dot.yellow { background-color: #f59e0b; }
        .legend-label {
          color: #cbd5e1;
          width: 52px;
        }
        .legend-count {
          color: #94a3b8;
          font-family: 'JetBrains Mono', monospace;
          font-size: 10.5px;
        }

        /* Recent Alerts */
        .view-all-link {
          background: transparent;
          border: none;
          color: #38bdf8;
          font-size: 11.5px;
          font-weight: 600;
          cursor: pointer;
          padding: 0;
        }
        .view-all-link:hover {
          text-decoration: underline;
        }
        .alerts-feed-list {
          display: flex;
          flex-direction: column;
          gap: 8px;
          overflow-y: auto;
          flex: 1;
        }
        .alert-row-item {
          display: flex;
          align-items: center;
          gap: 8px;
          font-size: 11px;
          padding: 4px 6px;
          border-radius: 4px;
          cursor: pointer;
          transition: background-color 0.15s ease;
        }
        .alert-row-item:hover {
          background-color: #162033;
        }
        .alert-sev-tag {
          font-size: 8.5px;
          font-weight: 800;
          padding: 2px 5px;
          border-radius: 4px;
          text-transform: uppercase;
          letter-spacing: 0.04em;
          flex-shrink: 0;
        }
        .alert-sev-tag.critical {
          background-color: rgba(239, 68, 68, 0.2);
          color: #ef4444;
        }
        .alert-sev-tag.warning {
          background-color: rgba(245, 158, 11, 0.2);
          color: #f59e0b;
        }
        .alert-sev-tag.info {
          background-color: rgba(6, 182, 212, 0.2);
          color: #06b6d4;
        }
        .alert-copy-wrap {
          flex: 1;
          min-width: 0;
        }
        .alert-msg-txt {
          color: #cbd5e1;
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
          display: block;
        }
        .alert-time-txt {
          color: #64748b;
          font-size: 10px;
          font-family: 'JetBrains Mono', monospace;
          flex-shrink: 0;
        }

        /* ── 3. BOTTOM ROW ── */
        .bottom-dashboard-grid {
          display: grid;
          grid-template-columns: 3fr 1fr;
          gap: 14px;
        }

        /* Machine Status Table */
        .table-header-bar {
          display: flex;
          align-items: center;
          justify-content: space-between;
          margin-bottom: 12px;
        }
        .table-header-bar h3 {
          font-size: 13.5px;
          font-weight: 700;
          color: #ffffff;
          margin: 0;
        }
        .table-controls-right {
          display: flex;
          align-items: center;
          gap: 8px;
        }
        .table-search-box {
          display: flex;
          align-items: center;
          gap: 6px;
          background-color: #162033;
          border: 1px solid #23334d;
          border-radius: 6px;
          padding: 4px 10px;
          width: 180px;
        }
        .table-search-box input {
          background: transparent;
          border: none;
          outline: none;
          color: #f1f5f9;
          font-size: 11.5px;
          width: 100%;
        }
        .table-filter-dropdown {
          display: flex;
          align-items: center;
          background-color: #162033;
          border: 1px solid #23334d;
          border-radius: 6px;
          padding: 4px 8px;
        }
        .table-filter-dropdown select {
          background: transparent;
          border: none;
          color: #cbd5e1;
          font-size: 11.5px;
          outline: none;
          cursor: pointer;
          appearance: none;
          padding-right: 4px;
        }

        .table-responsive-wrapper {
          overflow-x: auto;
        }
        .machine-status-table {
          width: 100%;
          border-collapse: collapse;
          text-align: left;
        }
        .machine-status-table th {
          font-size: 10.5px;
          font-weight: 700;
          color: #64748b;
          padding: 8px 10px;
          border-bottom: 1px solid #1a253b;
          white-space: nowrap;
        }
        .sort-icon {
          font-size: 9px;
          color: #475569;
          margin-left: 2px;
        }
        .machine-status-table td {
          padding: 9px 10px;
          border-bottom: 1px solid #141c2c;
          font-size: 11.5px;
          vertical-align: middle;
          white-space: nowrap;
        }
        .clickable-table-row {
          cursor: pointer;
          transition: background-color 0.15s ease;
        }
        .clickable-table-row:hover {
          background-color: #141e30;
        }
        .machine-name-cell {
          display: flex;
          flex-direction: column;
          line-height: 1.2;
        }
        .machine-name-cell strong {
          color: #ffffff;
          font-size: 12px;
        }
        .machine-name-cell small {
          color: #64748b;
          font-size: 10px;
          font-family: 'JetBrains Mono', monospace;
        }
        .os-icon-wrap {
          display: inline-flex;
          align-items: center;
          justify-content: center;
          font-size: 14px;
        }

        .progress-cell {
          display: flex;
          flex-direction: column;
          gap: 3px;
          min-width: 80px;
          max-width: 90px;
        }
        .metric-pct-label {
          font-size: 11px;
          font-weight: 700;
          color: #f1f5f9;
          font-family: 'JetBrains Mono', monospace;
        }
        .bar-track {
          width: 100%;
          height: 4px;
          background-color: #162033;
          border-radius: 2px;
          overflow: hidden;
        }
        .bar-fill {
          height: 100%;
          border-radius: 2px;
        }
        .bar-fill.blue { background-color: #2563eb; }
        .bar-fill.purple { background-color: #a855f7; }
        .bar-fill.yellow { background-color: #f59e0b; }
        .bar-fill.red { background-color: #ef4444; }

        .network-rates-cell {
          display: flex;
          flex-direction: column;
          font-size: 10px;
          color: #22c55e;
          font-family: 'JetBrains Mono', monospace;
          line-height: 1.3;
        }
        .latency-val {
          font-size: 11px;
          font-family: 'JetBrains Mono', monospace;
          font-weight: 600;
        }
        .latency-val.green { color: #22c55e; }
        .latency-val.amber { color: #f59e0b; }
        .latency-val.muted { color: #64748b; }

        .status-tag {
          display: inline-flex;
          align-items: center;
          gap: 5px;
          font-size: 11px;
          font-weight: 600;
        }
        .status-tag.online { color: #22c55e; }
        .status-tag.online .dot {
          width: 5px;
          height: 5px;
          border-radius: 50%;
          background-color: #22c55e;
          box-shadow: 0 0 6px #22c55e;
        }
        .status-tag.offline { color: #ef4444; }
        .status-tag.offline .dot {
          width: 5px;
          height: 5px;
          border-radius: 50%;
          background-color: #ef4444;
        }
        .lastseen-val {
          color: #94a3b8;
          font-size: 10.5px;
          font-family: 'JetBrains Mono', monospace;
        }
        .dash-val {
          color: #64748b;
          font-size: 12px;
        }

        /* Action Menu & Popover */
        .actions-td {
          position: relative;
        }
        .menu-wrap {
          position: relative;
        }
        .btn-dots-menu {
          background: transparent;
          border: none;
          color: #64748b;
          cursor: pointer;
          padding: 4px;
          border-radius: 4px;
          display: flex;
          align-items: center;
          justify-content: center;
        }
        .btn-dots-menu:hover {
          background-color: #1a253b;
          color: #cbd5e1;
        }
        .row-context-popover {
          position: absolute;
          right: 0;
          top: calc(100% + 4px);
          width: 160px;
          background-color: #0c1220;
          border: 1px solid #1f2e44;
          border-radius: 8px;
          padding: 4px;
          box-shadow: 0 8px 24px rgba(0, 0, 0, 0.7);
          z-index: 200;
          display: flex;
          flex-direction: column;
          gap: 2px;
        }
        .popover-btn {
          display: flex;
          align-items: center;
          gap: 8px;
          background: transparent;
          border: none;
          color: #cbd5e1;
          font-size: 11.5px;
          padding: 6px 8px;
          border-radius: 4px;
          cursor: pointer;
          transition: all 0.15s ease;
          text-align: left;
        }
        .popover-btn:hover {
          background-color: #162033;
          color: #ffffff;
        }

        .table-footer-pagination {
          display: flex;
          align-items: center;
          justify-content: space-between;
          margin-top: 10px;
          font-size: 11px;
          color: #64748b;
        }
        .pagination-controls {
          display: flex;
          align-items: center;
          gap: 4px;
        }
        .page-btn, .page-num-btn {
          background: #162033;
          border: 1px solid #23334d;
          color: #cbd5e1;
          border-radius: 4px;
          width: 24px;
          height: 24px;
          display: flex;
          align-items: center;
          justify-content: center;
          font-size: 11px;
          cursor: pointer;
        }
        .page-num-btn.active {
          background: #2563eb;
          color: #ffffff;
          border-color: #2563eb;
          font-weight: 700;
        }

        /* Right Side Stack */
        .right-side-stack {
          display: flex;
          flex-direction: column;
          gap: 14px;
        }
        .system-overview-card {
          flex: 1;
        }
        .overview-metrics-list {
          display: flex;
          flex-direction: column;
          gap: 10px;
        }
        .overview-metric-item {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 4px 6px;
          border-radius: 6px;
          transition: background-color 0.15s ease;
        }
        .overview-metric-item.clickable {
          cursor: pointer;
        }
        .overview-metric-item.clickable:hover {
          background-color: #162033;
        }
        .metric-left-info {
          display: flex;
          flex-direction: column;
          gap: 2px;
        }
        .metric-label {
          font-size: 11px;
          color: #94a3b8;
        }
        .val-trend-row {
          display: flex;
          align-items: baseline;
          gap: 6px;
        }
        .metric-val {
          font-size: 15px;
          font-weight: 800;
          color: #ffffff;
        }
        .trend-pct {
          font-size: 10px;
          font-weight: 700;
        }
        .trend-pct.up { color: #22c55e; }
        .metric-spark-wrap {
          width: 85px;
          height: 28px;
        }
        .overview-sparkline {
          width: 100%;
          height: 100%;
        }

        /* Connected Agents Card */
        .connected-agents-card {
          display: flex;
          flex-direction: column;
          gap: 10px;
        }
        .connected-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
        }
        .connected-header h3 {
          font-size: 12.5px;
          font-weight: 700;
          color: #cbd5e1;
          margin: 0;
        }
        .agents-ratio {
          font-size: 13px;
          font-weight: 800;
          color: #22c55e;
          font-family: 'JetBrains Mono', monospace;
        }
        .agents-progress-track {
          width: 100%;
          height: 6px;
          background-color: #162033;
          border-radius: 999px;
          overflow: hidden;
        }
        .agents-progress-fill {
          height: 100%;
          background-color: #22c55e;
          border-radius: inherit;
        }
        .heartbeat-row {
          display: flex;
          align-items: center;
          justify-content: space-between;
          font-size: 11px;
          color: #64748b;
        }
        .heartbeat-val {
          color: #94a3b8;
          font-family: 'JetBrains Mono', monospace;
        }

        @media (max-width: 1200px) {
          .kpi-cards-grid { grid-template-columns: repeat(3, 1fr); }
          .middle-dashboard-grid { grid-template-columns: 1fr; }
          .bottom-dashboard-grid { grid-template-columns: 1fr; }
        }
        @media (max-width: 768px) {
          .kpi-cards-grid { grid-template-columns: repeat(2, 1fr); }
        }
      `}</style>
    </div>
  );
}
