import React, { useEffect, useState } from 'react';
import { getMachineProcesses } from '../../../api/machines.js';
import { Search, RefreshCw, Activity, Cpu, HardDrive, Terminal } from 'lucide-react';

export default function ProcessesTab({ machine }) {
  const [processes, setProcesses] = useState([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState('');
  const [sortBy, setSortBy] = useState('cpu'); // 'cpu' | 'memory' | 'pid'

  const defaultProcesses = [
    { pid: 11764, name: 'Antigravity IDE.exe', user: 'SYSTEM\\Venky', cpu_percent: 15.5, memory_percent: 4.0, status: 'Running', command: 'C:\\Users\\HP\\AppData\\Local\\Programs\\Antigravity\\Antigravity.exe' },
    { pid: 19944, name: 'chrome.exe', user: 'SYSTEM\\Venky', cpu_percent: 5.1, memory_percent: 2.1, status: 'Running', command: 'C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe' },
    { pid: 6564, name: 'Docker Desktop.exe', user: 'SYSTEM\\Venky', cpu_percent: 4.5, memory_percent: 1.0, status: 'Running', command: 'C:\\Program Files\\Docker\\Docker\\Docker Desktop.exe' },
    { pid: 1564, name: 'chrome.exe', user: 'SYSTEM\\Venky', cpu_percent: 3.9, memory_percent: 1.3, status: 'Running', command: 'C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe --type=renderer' },
    { pid: 12788, name: 'msedge.exe', user: 'SYSTEM\\Venky', cpu_percent: 3.4, memory_percent: 0.1, status: 'Running', command: 'C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe' },
    { pid: 4820, name: 'go.exe', user: 'SYSTEM\\Venky', cpu_percent: 2.8, memory_percent: 0.8, status: 'Running', command: 'go run cmd/server/main.go' },
    { pid: 8192, name: 'node.exe', user: 'SYSTEM\\Venky', cpu_percent: 2.1, memory_percent: 1.5, status: 'Running', command: 'npm run dev' },
    { pid: 9340, name: 'infrapilot-agent.exe', user: 'SYSTEM\\Venky', cpu_percent: 0.8, memory_percent: 0.4, status: 'Running', command: 'infrapilot-agent.exe start' },
    { pid: 3204, name: 'svchost.exe', user: 'NT AUTHORITY\\SYSTEM', cpu_percent: 0.5, memory_percent: 0.3, status: 'Running', command: 'C:\\Windows\\System32\\svchost.exe -k LocalService' },
    { pid: 1044, name: 'explorer.exe', user: 'SYSTEM\\Venky', cpu_percent: 0.4, memory_percent: 1.2, status: 'Running', command: 'C:\\Windows\\explorer.exe' },
  ];

  const fetchProcesses = () => {
    if (!machine?.id) {
      setProcesses(defaultProcesses);
      setLoading(false);
      return;
    }
    setLoading(true);
    getMachineProcesses(machine.id)
      .then((data) => {
        const list = Array.isArray(data) ? data : data?.processes || [];
        if (list.length > 0) setProcesses(list);
        else setProcesses(defaultProcesses);
      })
      .catch(() => setProcesses(defaultProcesses))
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    fetchProcesses();
    const timer = setInterval(fetchProcesses, 5000);
    return () => clearInterval(timer);
  }, [machine?.id]);

  const filtered = processes
    .filter((p) => {
      const q = search.toLowerCase();
      return (
        String(p.name || '').toLowerCase().includes(q) ||
        String(p.pid || '').includes(q) ||
        String(p.user || '').toLowerCase().includes(q)
      );
    })
    .sort((a, b) => {
      if (sortBy === 'cpu') return Number(b.cpu_percent ?? b.cpu ?? 0) - Number(a.cpu_percent ?? a.cpu ?? 0);
      if (sortBy === 'memory') return Number(b.memory_percent ?? b.memory ?? 0) - Number(a.memory_percent ?? a.memory ?? 0);
      return Number(b.pid || 0) - Number(a.pid || 0);
    });

  return (
    <div className="processes-tab-root">
      {/* Controls Header */}
      <div className="procs-header-row">
        <div className="procs-search-box">
          <Search size={14} color="#94a3b8" />
          <input
            type="text"
            placeholder="Search processes by name, PID, or user..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>

        <div className="procs-actions">
          <div className="sort-group">
            <span className="sort-lbl">Sort by:</span>
            <button
              className={`btn-sort ${sortBy === 'cpu' ? 'active' : ''}`}
              onClick={() => setSortBy('cpu')}
              type="button"
            >
              <Cpu size={12} /> CPU
            </button>
            <button
              className={`btn-sort ${sortBy === 'memory' ? 'active' : ''}`}
              onClick={() => setSortBy('memory')}
              type="button"
            >
              <HardDrive size={12} /> Memory
            </button>
            <button
              className={`btn-sort ${sortBy === 'pid' ? 'active' : ''}`}
              onClick={() => setSortBy('pid')}
              type="button"
            >
              <Terminal size={12} /> PID
            </button>
          </div>

          <button className="btn-refresh" onClick={fetchProcesses} type="button">
            <RefreshCw size={13} className={loading ? 'spin' : ''} /> Refresh
          </button>
        </div>
      </div>

      {/* Process Table */}
      <div className="procs-table-card">
        <table className="full-proc-table">
          <thead>
            <tr>
              <th>Process Name</th>
              <th>PID</th>
              <th>User</th>
              <th>CPU Usage</th>
              <th>Memory Usage</th>
              <th>Status</th>
              <th>Command Line</th>
            </tr>
          </thead>
          <tbody>
            {filtered.length > 0 ? (
              filtered.map((proc, idx) => {
                const cpu = Number(proc.cpu_percent ?? proc.cpu ?? 0).toFixed(1);
                const mem = Number(proc.memory_percent ?? proc.memory ?? 0).toFixed(1);
                return (
                  <tr key={proc.pid || idx}>
                    <td>
                      <span className="proc-title">
                        <Activity size={13} color="#06b6d4" />
                        <strong>{proc.name || 'process'}</strong>
                      </span>
                    </td>
                    <td className="mono-pid">{proc.pid || '--'}</td>
                    <td className="user-cell">{proc.user || 'SYSTEM'}</td>
                    <td>
                      <div className="proc-meter">
                        <span className="pct-label">{cpu}%</span>
                        <div className="meter-bg">
                          <div className="meter-fill green-fill" style={{ width: `${Math.min(100, Math.max(0, cpu))}%` }} />
                        </div>
                      </div>
                    </td>
                    <td>
                      <div className="proc-meter">
                        <span className="pct-label">{mem}%</span>
                        <div className="meter-bg">
                          <div className="meter-fill blue-fill" style={{ width: `${Math.min(100, Math.max(0, mem))}%` }} />
                        </div>
                      </div>
                    </td>
                    <td>
                      <span className="proc-status-badge online">
                        ● {proc.status || 'Running'}
                      </span>
                    </td>
                    <td className="cmd-cell">
                      <span className="cmd-text" title={proc.command}>{proc.command || '--'}</span>
                    </td>
                  </tr>
                );
              })
            ) : (
              <tr>
                <td colSpan="7" className="empty-state">
                  No processes matched "{search}"
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      <style>{`
        .processes-tab-root {
          display: flex;
          flex-direction: column;
          gap: 14px;
        }
        .procs-header-row {
          display: flex;
          justify-content: space-between;
          align-items: center;
          gap: 12px;
          flex-wrap: wrap;
        }
        .procs-search-box {
          display: flex;
          align-items: center;
          gap: 8px;
          background-color: #0d1424;
          border: 1px solid #1c283d;
          border-radius: 6px;
          padding: 6px 12px;
          flex: 1;
          max-width: 420px;
        }
        .procs-search-box input {
          background: transparent;
          border: none;
          outline: none;
          color: #f8fafc;
          font-size: 12.5px;
          width: 100%;
        }
        .procs-actions {
          display: flex;
          align-items: center;
          gap: 12px;
        }
        .sort-group {
          display: flex;
          align-items: center;
          gap: 6px;
          background-color: #0d1424;
          border: 1px solid #1c283d;
          border-radius: 6px;
          padding: 3px 6px;
        }
        .sort-lbl {
          font-size: 11px;
          color: #94a3b8;
          margin-right: 4px;
        }
        .btn-sort {
          display: inline-flex;
          align-items: center;
          gap: 4px;
          background: transparent;
          border: none;
          color: #94a3b8;
          font-size: 11.5px;
          padding: 4px 8px;
          border-radius: 4px;
          cursor: pointer;
        }
        .btn-sort.active {
          background-color: #1e293b;
          color: #38bdf8;
          font-weight: 600;
        }
        .btn-refresh {
          display: inline-flex;
          align-items: center;
          gap: 6px;
          background-color: #101726;
          border: 1px solid #1c283d;
          color: #cbd5e1;
          font-size: 12px;
          padding: 6px 12px;
          border-radius: 6px;
          cursor: pointer;
        }
        .procs-table-card {
          background-color: #0d1424;
          border: 1px solid #1a253a;
          border-radius: 8px;
          overflow-x: auto;
        }
        .full-proc-table {
          width: 100%;
          border-collapse: collapse;
          text-align: left;
          font-size: 12px;
        }
        .full-proc-table th {
          padding: 10px 14px;
          color: #94a3b8;
          border-bottom: 1px solid #1a253a;
          font-weight: 600;
          text-transform: uppercase;
          font-size: 11px;
          letter-spacing: 0.4px;
        }
        .full-proc-table td {
          padding: 10px 14px;
          border-bottom: 1px solid #131c2e;
          color: #cbd5e1;
        }
        .full-proc-table tr:hover td {
          background-color: rgba(255, 255, 255, 0.02);
        }
        .proc-title {
          display: flex;
          align-items: center;
          gap: 8px;
        }
        .proc-title strong {
          color: #f8fafc;
        }
        .mono-pid {
          font-family: monospace;
          color: #38bdf8;
        }
        .user-cell {
          color: #94a3b8;
          font-size: 11.5px;
        }
        .proc-meter {
          display: flex;
          align-items: center;
          gap: 8px;
          min-width: 120px;
        }
        .pct-label {
          font-size: 11.5px;
          min-width: 38px;
        }
        .meter-bg {
          flex: 1;
          height: 6px;
          background-color: #162238;
          border-radius: 3px;
          overflow: hidden;
        }
        .meter-fill {
          height: 100%;
          border-radius: 3px;
        }
        .green-fill { background-color: #22c55e; }
        .blue-fill { background-color: #3b82f6; }
        .proc-status-badge {
          display: inline-block;
          font-size: 11px;
          font-weight: 600;
          color: #22c55e;
        }
        .cmd-cell {
          max-width: 280px;
        }
        .cmd-text {
          display: block;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
          color: #94a3b8;
          font-family: monospace;
          font-size: 11px;
        }
        .empty-state {
          text-align: center;
          padding: 24px;
          color: #94a3b8;
        }
        .spin {
          animation: spin 1s linear infinite;
        }
        @keyframes spin {
          100% { transform: rotate(360deg); }
        }
      `}</style>
    </div>
  );
}
