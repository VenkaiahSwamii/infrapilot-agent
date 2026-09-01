import React, { useEffect, useState } from 'react';
import { History, ShieldAlert, Clock, RefreshCw, Search } from 'lucide-react';
import { apiClient } from '../api/client.js';
import { useDashboardStore } from '../store/dashboardStore.jsx';

export default function AuditLogsPage() {
  const [logs, setLogs] = useState([]);
  const [loading, setLoading] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const { addToast } = useDashboardStore();

  const fetchAuditLogs = async () => {
    setLoading(true);
    try {
      const res = await apiClient.get('/audit-logs');
      // The backend returns an array of audit logs
      if (Array.isArray(res.data)) {
        setLogs(res.data);
      } else if (res.data && Array.isArray(res.data.audit_logs)) {
        setLogs(res.data.audit_logs);
      }
    } catch (err) {
      console.error('Failed to fetch global audit logs:', err);
      addToast('warning', 'Audit Logs Error', 'Unable to retrieve audit logs trail.');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchAuditLogs();
  }, []);

  const filteredLogs = React.useMemo(() => {
    return logs.filter(l => {
      if (!searchQuery) return true;
      const q = searchQuery.toLowerCase();
      return (
        String(l.username || '').toLowerCase().includes(q) ||
        String(l.action || '').toLowerCase().includes(q) ||
        String(l.resource || '').toLowerCase().includes(q) ||
        String(l.result || '').toLowerCase().includes(q)
      );
    });
  }, [logs, searchQuery]);

  return (
    <div className="audit-logs-page">
      <div className="page-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px' }}>
        <div>
          <h1 style={{ fontSize: '24px', fontWeight: 700, color: '#f1f5f9', display: 'flex', alignItems: 'center', gap: '10px' }}>
            <History size={26} color="#06b6d4" />
            Global Audit Trails
          </h1>
          <p style={{ fontSize: '13px', color: '#94a3b8', marginTop: '4px' }}>
            Permanent, tamper-evident cryptographic log of all administrative actions in the platform.
          </p>
        </div>
        <button
          className="refresh-btn"
          onClick={fetchAuditLogs}
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
          {loading ? 'Refreshing...' : 'Refresh Logs'}
        </button>
      </div>

      {/* Toolbar */}
      <div className="audit-toolbar" style={{ display: 'flex', gap: '16px', marginBottom: '20px' }}>
        <div className="search-box" style={{ display: 'flex', alignItems: 'center', gap: '8px', backgroundColor: '#0d1220', border: '1px solid #1f2e44', borderRadius: '8px', padding: '8px 12px', flex: 1 }}>
          <Search size={16} color="#64748b" />
          <input
            type="text"
            placeholder="Search audit trail by user, action, target..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            style={{ background: 'none', border: 'none', color: '#f1f5f9', outline: 'none', width: '100%', fontSize: '13px' }}
          />
        </div>
      </div>

      {/* Logs Table */}
      <div className="table-wrapper" style={{ backgroundColor: '#0d1220', border: '1px solid #1f2e44', borderRadius: '12px', overflow: 'hidden' }}>
        {filteredLogs.length === 0 ? (
          <div className="empty-state" style={{ textAlign: 'center', padding: '60px 24px', color: '#64748b' }}>
            <ShieldAlert size={40} style={{ opacity: 0.3, marginBottom: '12px' }} />
            <strong>No audit events found</strong>
          </div>
        ) : (
          <table style={{ width: '100%', borderCollapse: 'collapse', textAlign: 'left', fontSize: '13px' }}>
            <thead>
              <tr style={{ backgroundColor: '#080c14', borderBottom: '1px solid #1f2e44', color: '#64748b' }}>
                <th style={{ padding: '14px 16px', fontWeight: 600 }}>Timestamp</th>
                <th style={{ padding: '14px 16px', fontWeight: 600 }}>User Actor</th>
                <th style={{ padding: '14px 16px', fontWeight: 600 }}>Event Action</th>
                <th style={{ padding: '14px 16px', fontWeight: 600 }}>Target Resource</th>
                <th style={{ padding: '14px 16px', fontWeight: 600 }}>Status Result</th>
              </tr>
            </thead>
            <tbody>
              {filteredLogs.map((l, index) => {
                const isFail = String(l.result || '').toLowerCase().includes('fail') || String(l.result || '').toLowerCase().includes('block');
                return (
                  <tr key={l.id || index} style={{ borderBottom: '1px solid rgba(31, 46, 68, 0.5)', color: '#cbd5e1' }}>
                    <td style={{ padding: '14px 16px', color: '#64748b' }}>
                      <span style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
                        <Clock size={12} />
                        {new Date(l.created_at || l.CreatedAt || Date.now()).toLocaleString()}
                      </span>
                    </td>
                    <td style={{ padding: '14px 16px', fontWeight: 600, color: '#f1f5f9' }}>{l.username || 'System'}</td>
                    <td style={{ padding: '14px 16px' }}>{l.action}</td>
                    <td style={{ padding: '14px 16px', fontFamily: 'monospace', color: '#64748b' }}>{l.resource || '--'}</td>
                    <td style={{ padding: '14px 16px' }}>
                      <span
                        style={{
                          display: 'inline-flex',
                          padding: '3px 8px',
                          borderRadius: '4px',
                          fontSize: '11px',
                          fontWeight: 600,
                          backgroundColor: isFail ? 'rgba(239, 68, 68, 0.1)' : 'rgba(34, 197, 94, 0.1)',
                          color: isFail ? '#ef4444' : '#22c55e',
                        }}
                      >
                        {l.result || 'Success'}
                      </span>
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
