import React, { useState } from 'react';
import { Play, BookOpen, Terminal, CheckCircle2 } from 'lucide-react';
import { apiPost } from '../../api/client.js';

export default function Runbooks({ runbooks, onRefresh }) {
  const [executing, setExecuting] = useState(null);
  const [resultModal, setResultModal] = useState(null);

  const list = runbooks?.runbooks || [
    { id: '1', name: 'Restart Nginx Web Server', category: 'service', command: 'sudo systemctl restart nginx' },
    { id: '2', name: 'Restart Apache HTTPD', category: 'service', command: 'sudo systemctl restart apache2' },
    { id: '3', name: 'Restart Docker Daemon', category: 'docker', command: 'sudo systemctl restart docker' },
    { id: '4', name: 'Restart Kubernetes Deployment', category: 'kubernetes', command: 'kubectl rollout restart deployment/{{target}}' },
    { id: '5', name: 'Clear Cache & Buffers', category: 'maintenance', command: 'sudo sync && echo 3 | sudo tee /proc/sys/vm/drop_caches' },
    { id: '6', name: 'Rotate Logs', category: 'maintenance', command: 'sudo logrotate -f /etc/logrotate.conf' },
    { id: '7', name: 'Cleanup Disk', category: 'maintenance', command: 'sudo apt-get clean && sudo journalctl --vacuum-time=3d' },
    { id: '8', name: 'Restart Database', category: 'database', command: 'sudo systemctl restart postgresql' },
  ];

  const handleExecute = async (rb) => {
    setExecuting(rb.id);
    try {
      const res = await apiPost('/automation/run', {
        target: 'api-prod-node-01',
        action_type: rb.command,
      });
      setResultModal({ name: rb.name, output: res.output || 'Command executed successfully.' });
      if (onRefresh) onRefresh();
    } catch (err) {
      setResultModal({ name: rb.name, output: `Execution failed: ${err.message}` });
    } finally {
      setExecuting(null);
    }
  };

  return (
    <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '12px', padding: '24px' }}>
      <div style={{ marginBottom: '20px' }}>
        <h3 style={{ margin: 0, color: '#a855f7', fontSize: '18px', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '8px' }}>
          <BookOpen size={20} color="#a855f7" />
          Reusable Operational Runbooks Library
        </h3>
        <span style={{ fontSize: '13px', color: '#8b949e' }}>Pre-approved operational procedures executable on-demand via 1-Click action</span>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))', gap: '20px' }}>
        {list.map((rb) => (
          <div key={rb.id} style={{
            background: '#0d1117',
            border: '1px solid #30363d',
            borderRadius: '10px',
            padding: '20px',
            display: 'flex',
            flexDirection: 'column',
            justifyContent: 'space-between'
          }}>
            <div>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '8px' }}>
                <span style={{
                  padding: '3px 8px',
                  borderRadius: '4px',
                  fontSize: '11px',
                  fontWeight: 700,
                  textTransform: 'uppercase',
                  backgroundColor: 'rgba(168, 85, 247, 0.15)',
                  color: '#a855f7',
                }}>
                  {rb.category}
                </span>
              </div>

              <h4 style={{ margin: '4px 0 8px 0', color: '#f0f6fc', fontSize: '15px' }}>{rb.name}</h4>
              <div style={{ background: '#161b22', padding: '8px 10px', borderRadius: '6px', fontSize: '12px', color: '#58a6ff', fontFamily: 'monospace', marginBottom: '16px' }}>
                {rb.command}
              </div>
            </div>

            <button
              onClick={() => handleExecute(rb)}
              disabled={executing === rb.id}
              style={{
                backgroundColor: '#1f6feb',
                color: '#fff',
                border: 'none',
                padding: '8px 14px',
                borderRadius: '6px',
                cursor: 'pointer',
                fontWeight: 600,
                fontSize: '13px',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                gap: '6px'
              }}
            >
              <Play size={14} />
              {executing === rb.id ? 'Executing...' : 'Run Runbook'}
            </button>
          </div>
        ))}
      </div>

      {resultModal && (
        <div style={{
          position: 'fixed', top: 0, left: 0, right: 0, bottom: 0,
          backgroundColor: 'rgba(0,0,0,0.7)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 1000
        }}>
          <div style={{ backgroundColor: '#161b22', border: '1px solid #30363d', borderRadius: '12px', width: '500px', padding: '24px' }}>
            <h3 style={{ margin: '0 0 12px 0', color: '#f0f6fc' }}>Execution Result: {resultModal.name}</h3>
            <pre style={{ background: '#0d1117', border: '1px solid #30363d', color: '#3fb950', padding: '12px', borderRadius: '6px', fontSize: '12px', overflowX: 'auto' }}>
              {resultModal.output}
            </pre>
            <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: '16px' }}>
              <button onClick={() => setResultModal(null)} style={{ background: '#21262d', color: '#c9d1d9', border: '1px solid #30363d', padding: '8px 16px', borderRadius: '6px', cursor: 'pointer', fontWeight: 600 }}>Close</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
