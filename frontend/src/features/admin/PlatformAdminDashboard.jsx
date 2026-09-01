import React, { useState, useEffect, useCallback } from 'react';
import {
  ShieldAlert, Activity, RefreshCw, Server, Database, HardDrive,
  Cpu, Calendar, Play, Download, Trash2, CheckCircle
} from 'lucide-react';
import { apiGet, apiPost } from '../../api/client.js';
import './SuperAdminDashboard.css'; // Reuse core styles

export default function PlatformAdminDashboard() {
  const [status, setStatus] = useState(null);
  const [backups, setBackups] = useState([]);
  const [loading, setLoading] = useState(true);
  const [actionMessage, setActionMessage] = useState('');

  const fetchStatusAndBackups = useCallback(async () => {
    setLoading(true);
    try {
      const [statusRes, backupRes] = await Promise.allSettled([
        apiGet('/admin/platform/status'),
        apiGet('/admin/backup/list')
      ]);

      if (statusRes.status === 'fulfilled') {
        setStatus(statusRes.value);
      }
      if (backupRes.status === 'fulfilled') {
        setBackups(backupRes.value?.backups || []);
      }
    } catch (err) {
      console.error('Failed to fetch platform metrics:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchStatusAndBackups();
    const timer = setInterval(fetchStatusAndBackups, 10000);
    return () => clearInterval(timer);
  }, [fetchStatusAndBackups]);

  const handleRunBackup = async () => {
    setActionMessage('Running database backup...');
    try {
      const res = await apiPost('/admin/backup/run');
      setActionMessage(`Backup completed successfully: ${res.file_name}`);
      fetchStatusAndBackups();
    } catch (err) {
      setActionMessage(`Backup failed: ${err.message}`);
    }
  };

  const handleRestoreBackup = async (fileName) => {
    if (!window.confirm(`Are you sure you want to restore database from ${fileName}? This will overwrite current table data.`)) {
      return;
    }
    setActionMessage('Restoring database...');
    try {
      await apiPost('/admin/backup/restore', { file_name: fileName });
      setActionMessage('Database restored successfully.');
      fetchStatusAndBackups();
    } catch (err) {
      setActionMessage(`Restore failed: ${err.message}`);
    }
  };

  const getStatusColor = (val) => {
    return val === 'Healthy' ? '#3fb950' : '#f78166';
  };

  return (
    <div className="super-admin-container" style={{ minHeight: '100vh', backgroundColor: '#0d1117', color: '#c9d1d9', padding: '24px 32px' }}>
      {/* Header */}
      <div className="super-admin-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '28px' }}>
        <div className="super-admin-title" style={{ display: 'flex', alignItems: 'center', gap: '14px' }}>
          <ShieldAlert size={32} color="#58a6ff" />
          <div>
            <h1 style={{ fontSize: '26px', fontWeight: 700, color: '#f0f6fc', margin: 0 }}>
              Platform System Administration
            </h1>
            <span style={{ fontSize: '13px', color: '#8b949e' }}>
              Real-Time Platform Infrastructure Control & Backup Governance
            </span>
          </div>
        </div>
        <button
          className="btn-secondary-dr"
          onClick={fetchStatusAndBackups}
          disabled={loading}
          style={{
            background: '#21262d',
            color: '#c9d1d9',
            border: '1px solid #30363d',
            padding: '8px 16px',
            borderRadius: '6px',
            cursor: 'pointer',
            display: 'flex',
            alignItems: 'center',
            gap: '8px'
          }}
        >
          <RefreshCw size={16} className={loading ? 'spin' : ''} />
          Sync Diagnostics
        </button>
      </div>

      {actionMessage && (
        <div style={{
          backgroundColor: '#161b22',
          border: '1px solid #30363d',
          borderRadius: '8px',
          padding: '12px 16px',
          marginBottom: '24px',
          color: '#58a6ff',
          fontSize: '14px',
          fontWeight: 500
        }}>
          {actionMessage}
        </div>
      )}

      {/* Component Status Grid */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))', gap: '20px', marginBottom: '32px' }}>
        <div className="super-card" style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '10px', padding: '20px', position: 'relative' }}>
          <div className="super-card-header" style={{ display: 'flex', justifyContent: 'space-between', color: '#8b949e', fontSize: '13px' }}>
            <span>BACKEND STATUS</span>
            <Activity size={18} color="#58a6ff" />
          </div>
          <div className="super-card-value" style={{ color: getStatusColor(status?.backend_status), fontSize: '26px', fontWeight: 700 }}>
            {status?.backend_status || 'Checking...'}
          </div>
          <div className="super-card-sub" style={{ fontSize: '12px', color: '#8b949e' }}>Stateless API Nodes operational</div>
        </div>

        <div className="super-card" style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '10px', padding: '20px', position: 'relative' }}>
          <div className="super-card-header" style={{ display: 'flex', justifyContent: 'space-between', color: '#8b949e', fontSize: '13px' }}>
            <span>DATABASE CLUSTER</span>
            <Database size={18} color="#3fb950" />
          </div>
          <div className="super-card-value" style={{ color: getStatusColor(status?.database_status), fontSize: '26px', fontWeight: 700 }}>
            {status?.database_status || 'Checking...'}
          </div>
          <div className="super-card-sub" style={{ fontSize: '12px', color: '#8b949e' }}>PostgreSQL server connection</div>
        </div>

        <div className="super-card" style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '10px', padding: '20px', position: 'relative' }}>
          <div className="super-card-header" style={{ display: 'flex', justifyContent: 'space-between', color: '#8b949e', fontSize: '13px' }}>
            <span>REDIS CACHE / PUB-SUB</span>
            <HardDrive size={18} color="#ffa657" />
          </div>
          <div className="super-card-value" style={{ color: getStatusColor(status?.redis_status), fontSize: '26px', fontWeight: 700 }}>
            {status?.redis_status || 'Checking...'}
          </div>
          <div className="super-card-sub" style={{ fontSize: '12px', color: '#8b949e' }}>Session cache and event bus</div>
        </div>
      </div>

      {/* Stats Counter Grid */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '20px', marginBottom: '32px' }}>
        <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '8px', padding: '16px', textAlign: 'center' }}>
          <span style={{ fontSize: '12px', color: '#8b949e', fontWeight: 600 }}>CONNECTED AGENTS</span>
          <div style={{ fontSize: '32px', fontWeight: 700, color: '#f0f6fc', margin: '8px 0' }}>
            {status?.connected_agents ?? 0}
          </div>
          <span style={{ fontSize: '12px', color: '#8b949e' }}>Streaming live metrics</span>
        </div>

        <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '8px', padding: '16px', textAlign: 'center' }}>
          <span style={{ fontSize: '12px', color: '#8b949e', fontWeight: 600 }}>CONNECTED DASHBOARDS</span>
          <div style={{ fontSize: '32px', fontWeight: 700, color: '#f0f6fc', margin: '8px 0' }}>
            {status?.connected_dashboards ?? 0}
          </div>
          <span style={{ fontSize: '12px', color: '#8b949e' }}>Active WS browser links</span>
        </div>

        <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '8px', padding: '16px', textAlign: 'center' }}>
          <span style={{ fontSize: '12px', color: '#8b949e', fontWeight: 600 }}>THROUGHPUT RATE</span>
          <div style={{ fontSize: '32px', fontWeight: 700, color: '#f0f6fc', margin: '8px 0' }}>
            {status?.api_requests_per_sec ?? 0} r/s
          </div>
          <span style={{ fontSize: '12px', color: '#8b949e' }}>API Server Requests</span>
        </div>

        <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '8px', padding: '16px', textAlign: 'center' }}>
          <span style={{ fontSize: '12px', color: '#8b949e', fontWeight: 600 }}>WEBSOCKET CLIENTS</span>
          <div style={{ fontSize: '32px', fontWeight: 700, color: '#f0f6fc', margin: '8px 0' }}>
            {status?.websocket_clients ?? 0}
          </div>
          <span style={{ fontSize: '12px', color: '#8b949e' }}>Hub connections</span>
        </div>
      </div>

      {/* Backup and Restore Governance */}
      <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '10px', padding: '24px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
          <h3 style={{ margin: 0, color: '#f0f6fc', fontSize: '18px', fontWeight: 600 }}>
            Database Backup & Restore Governance
          </h3>
          <button
            onClick={handleRunBackup}
            style={{
              backgroundColor: '#2ea043',
              color: '#ffffff',
              border: 'none',
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
            <Download size={16} />
            Trigger Manual Backup
          </button>
        </div>

        {backups.length === 0 ? (
          <div style={{ color: '#8b949e', fontSize: '14px', fontStyle: 'italic', padding: '12px 0' }}>
            No backup snapshots found in directory.
          </div>
        ) : (
          <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
            {backups.map((b) => (
              <div key={b.file_name} style={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                background: '#0d1117',
                padding: '12px 18px',
                borderRadius: '8px',
                border: '1px solid #30363d'
              }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
                  <Calendar size={18} color="#58a6ff" />
                  <div>
                    <span style={{ fontWeight: 600, color: '#c9d1d9', fontSize: '14px' }}>{b.file_name}</span>
                    <div style={{ fontSize: '12px', color: '#8b949e', marginTop: '2px' }}>
                      Size: {Math.round(b.size_bytes / 1024)} KB | Mod time: {new Date(b.created_at).toLocaleString()}
                    </div>
                  </div>
                </div>
                <button
                  onClick={() => handleRestoreBackup(b.file_name)}
                  style={{
                    backgroundColor: '#1f6feb',
                    color: '#ffffff',
                    border: 'none',
                    padding: '6px 12px',
                    borderRadius: '4px',
                    cursor: 'pointer',
                    fontSize: '13px',
                    fontWeight: 600,
                    display: 'flex',
                    alignItems: 'center',
                    gap: '4px'
                  }}
                >
                  <Play size={12} />
                  Restore
                </button>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
