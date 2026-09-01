import React from 'react';
import { ShieldAlert, Clock, User, FileText } from 'lucide-react';

export default function AuditLogs({ auditLogs }) {
  const logs = auditLogs?.audit_logs || [];

  return (
    <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '12px', padding: '24px' }}>
      <div style={{ marginBottom: '20px' }}>
        <h3 style={{ margin: 0, color: '#f0f6fc', fontSize: '18px', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '8px' }}>
          <ShieldAlert size={20} color="#3fb950" />
          Organization Audit Log Trail
        </h3>
        <span style={{ fontSize: '13px', color: '#8b949e' }}>Chronological activity tracking for administrative events inside this organization</span>
      </div>

      {logs.length === 0 ? (
        <div style={{ color: '#8b949e', fontStyle: 'italic', padding: '16px 0' }}>No tenant audit logs recorded yet.</div>
      ) : (
        <div style={{ overflowX: 'auto' }}>
          <table style={{ width: '100%', borderCollapse: 'collapse', textAlign: 'left', fontSize: '13px' }}>
            <thead>
              <tr style={{ borderBottom: '1px solid #30363d', color: '#8b949e' }}>
                <th style={{ padding: '12px' }}>TIMESTAMP</th>
                <th style={{ padding: '12px' }}>USER</th>
                <th style={{ padding: '12px' }}>ACTION</th>
                <th style={{ padding: '12px' }}>RESOURCE</th>
                <th style={{ padding: '12px' }}>RESULT</th>
              </tr>
            </thead>
            <tbody>
              {logs.map((l) => (
                <tr key={l.id} style={{ borderBottom: '1px solid #21262d', color: '#c9d1d9' }}>
                  <td style={{ padding: '12px', color: '#8b949e', fontSize: '12px' }}>
                    {new Date(l.created_at).toLocaleString()}
                  </td>
                  <td style={{ padding: '12px', fontWeight: 600, color: '#f0f6fc' }}>{l.username}</td>
                  <td style={{ padding: '12px' }}>{l.action}</td>
                  <td style={{ padding: '12px', color: '#58a6ff' }}>{l.resource || '-'}</td>
                  <td style={{ padding: '12px' }}>
                    <span style={{
                      padding: '3px 8px',
                      borderRadius: '4px',
                      fontSize: '11px',
                      fontWeight: 600,
                      backgroundColor: l.result?.startsWith('Failure') ? 'rgba(247, 129, 102, 0.15)' : 'rgba(63, 185, 80, 0.15)',
                      color: l.result?.startsWith('Failure') ? '#f78166' : '#3fb950',
                    }}>
                      {l.result || 'Success'}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
