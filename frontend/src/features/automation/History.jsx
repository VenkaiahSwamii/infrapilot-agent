import React, { useState } from 'react';
import { History as HistoryIcon, Terminal, CheckCircle, XCircle } from 'lucide-react';

export default function History({ history }) {
  const [selectedItem, setSelectedItem] = useState(null);
  const list = history?.history || [];

  return (
    <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '12px', padding: '24px' }}>
      <div style={{ marginBottom: '20px' }}>
        <h3 style={{ margin: 0, color: '#f0f6fc', fontSize: '18px', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '8px' }}>
          <HistoryIcon size={20} color="#3fb950" />
          Automation Execution History & Terminal Audit
        </h3>
        <span style={{ fontSize: '13px', color: '#8b949e' }}>Complete record of rule triggers, manual runbooks, output logs, and execution status</span>
      </div>

      {list.length === 0 ? (
        <div style={{ color: '#8b949e', fontStyle: 'italic', padding: '16px 0' }}>No execution history recorded yet.</div>
      ) : (
        <div style={{ overflowX: 'auto' }}>
          <table style={{ width: '100%', borderCollapse: 'collapse', textAlign: 'left', fontSize: '13px' }}>
            <thead>
              <tr style={{ borderBottom: '1px solid #30363d', color: '#8b949e' }}>
                <th style={{ padding: '12px' }}>RULE / RUNBOOK</th>
                <th style={{ padding: '12px' }}>TRIGGER</th>
                <th style={{ padding: '12px' }}>EXECUTED BY</th>
                <th style={{ padding: '12px' }}>TARGET HOST</th>
                <th style={{ padding: '12px' }}>RESULT</th>
                <th style={{ padding: '12px' }}>TIME</th>
                <th style={{ padding: '12px', textAlign: 'right' }}>LOGS</th>
              </tr>
            </thead>
            <tbody>
              {list.map((h) => (
                <tr key={h.id} style={{ borderBottom: '1px solid #21262d', color: '#c9d1d9' }}>
                  <td style={{ padding: '12px', fontWeight: 600, color: '#f0f6fc' }}>{h.rule_name}</td>
                  <td style={{ padding: '12px', color: '#58a6ff' }}>{h.trigger}</td>
                  <td style={{ padding: '12px' }}>{h.user}</td>
                  <td style={{ padding: '12px', fontFamily: 'monospace' }}>{h.target}</td>
                  <td style={{ padding: '12px' }}>
                    <span style={{
                      padding: '3px 8px',
                      borderRadius: '4px',
                      fontSize: '11px',
                      fontWeight: 700,
                      backgroundColor: h.result === 'Success' ? 'rgba(63, 185, 80, 0.15)' : 'rgba(247, 129, 102, 0.15)',
                      color: h.result === 'Success' ? '#3fb950' : '#f78166',
                    }}>
                      {h.result}
                    </span>
                  </td>
                  <td style={{ padding: '12px', color: '#8b949e', fontSize: '12px' }}>{new Date(h.start_time).toLocaleString()}</td>
                  <td style={{ padding: '12px', textAlign: 'right' }}>
                    <button
                      onClick={() => setSelectedItem(h)}
                      style={{ background: '#21262d', color: '#58a6ff', border: '1px solid #30363d', padding: '4px 10px', borderRadius: '4px', cursor: 'pointer', fontSize: '12px' }}
                    >
                      View Output
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {selectedItem && (
        <div style={{
          position: 'fixed', top: 0, left: 0, right: 0, bottom: 0,
          backgroundColor: 'rgba(0,0,0,0.7)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 1000
        }}>
          <div style={{ backgroundColor: '#161b22', border: '1px solid #30363d', borderRadius: '12px', width: '550px', padding: '24px' }}>
            <h3 style={{ margin: '0 0 12px 0', color: '#f0f6fc', display: 'flex', alignItems: 'center', gap: '8px' }}>
              <Terminal size={18} color="#58a6ff" /> Terminal Output: {selectedItem.rule_name}
            </h3>
            <pre style={{ background: '#0d1117', border: '1px solid #30363d', color: '#3fb950', padding: '14px', borderRadius: '6px', fontSize: '12px', overflowX: 'auto', maxHeight: '300px' }}>
              {selectedItem.output_snippet || 'No output recorded.'}
            </pre>
            <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: '16px' }}>
              <button onClick={() => setSelectedItem(null)} style={{ background: '#21262d', color: '#c9d1d9', border: '1px solid #30363d', padding: '8px 16px', borderRadius: '6px', cursor: 'pointer', fontWeight: 600 }}>Close</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
