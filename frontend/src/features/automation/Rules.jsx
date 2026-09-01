import React, { useState } from 'react';
import { Sliders, Plus, CheckCircle, ShieldAlert, ToggleLeft, ToggleRight } from 'lucide-react';
import { apiPost } from '../../api/client.js';

export default function Rules({ rules, onRefresh }) {
  const [showModal, setShowModal] = useState(false);
  const [name, setName] = useState('');
  const [triggerEventType, setTriggerEventType] = useState('high_cpu');
  const [actionType, setActionType] = useState('restart_service');
  const [targetResource, setTargetResource] = useState('');
  const [requiresApproval, setRequiresApproval] = useState(false);

  const ruleList = rules?.rules || [];

  const handleCreateRule = async (e) => {
    e.preventDefault();
    try {
      await apiPost('/automation/rules', {
        name,
        trigger_event_type: triggerEventType,
        action_type: actionType,
        target_resource: targetResource || 'nginx.service',
        requires_approval: requiresApproval,
        enabled: true,
      });
      setShowModal(false);
      setName('');
      setTargetResource('');
      if (onRefresh) onRefresh();
    } catch (err) {
      console.error('Failed to create rule', err);
    }
  };

  return (
    <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '12px', padding: '24px' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
        <div>
          <h3 style={{ margin: 0, color: '#f0f6fc', fontSize: '18px', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '8px' }}>
            <Sliders size={20} color="#58a6ff" />
            Auto-Remediation Rules Engine
          </h3>
          <span style={{ fontSize: '13px', color: '#8b949e' }}>IF Alert Condition THEN Execute Remediation Action</span>
        </div>

        <button
          onClick={() => setShowModal(true)}
          style={{
            backgroundColor: '#238636',
            color: '#fff',
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
          <Plus size={16} />
          Create Rule
        </button>
      </div>

      <div style={{ display: 'flex', flexDirection: 'column', gap: '14px' }}>
        {ruleList.map((r) => (
          <div key={r.id} style={{ background: '#0d1117', border: '1px solid #30363d', borderRadius: '8px', padding: '16px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <div>
              <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                <h4 style={{ margin: 0, color: '#f0f6fc', fontSize: '15px' }}>{r.name}</h4>
                {r.requires_approval ? (
                  <span style={{ padding: '2px 8px', borderRadius: '4px', background: 'rgba(255, 166, 87, 0.15)', color: '#ffa657', fontSize: '11px', fontWeight: 700 }}>
                    APPROVAL REQUIRED
                  </span>
                ) : (
                  <span style={{ padding: '2px 8px', borderRadius: '4px', background: 'rgba(63, 185, 80, 0.15)', color: '#3fb950', fontSize: '11px', fontWeight: 700 }}>
                    AUTO EXECUTE
                  </span>
                )}
              </div>
              <div style={{ fontSize: '13px', color: '#8b949e', marginTop: '4px' }}>{r.description || `Trigger: ${r.trigger_event_type} -> Action: ${r.action_type}`}</div>
              <div style={{ fontSize: '12px', color: '#58a6ff', marginTop: '4px' }}>Target: <code>{r.target_resource}</code></div>
            </div>

            <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
              <span style={{ color: r.enabled ? '#3fb950' : '#8b949e', fontWeight: 600, fontSize: '13px' }}>
                {r.enabled ? 'ACTIVE' : 'DISABLED'}
              </span>
            </div>
          </div>
        ))}
      </div>

      {showModal && (
        <div style={{
          position: 'fixed', top: 0, left: 0, right: 0, bottom: 0,
          backgroundColor: 'rgba(0,0,0,0.7)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 1000
        }}>
          <div style={{ backgroundColor: '#161b22', border: '1px solid #30363d', borderRadius: '12px', width: '440px', padding: '24px' }}>
            <h3 style={{ margin: '0 0 16px 0', color: '#f0f6fc' }}>Create Auto-Remediation Rule</h3>
            <form onSubmit={handleCreateRule} style={{ display: 'flex', flexDirection: 'column', gap: '14px' }}>
              <div>
                <label style={{ fontSize: '12px', color: '#8b949e', display: 'block', marginBottom: '4px' }}>Rule Name</label>
                <input
                  type="text"
                  required
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="Auto-Restart Nginx on High CPU"
                  style={{ width: '100%', padding: '8px 12px', background: '#0d1117', border: '1px solid #30363d', color: '#fff', borderRadius: '6px' }}
                />
              </div>

              <div>
                <label style={{ fontSize: '12px', color: '#8b949e', display: 'block', marginBottom: '4px' }}>Trigger Event</label>
                <select
                  value={triggerEventType}
                  onChange={(e) => setTriggerEventType(e.target.value)}
                  style={{ width: '100%', padding: '8px 12px', background: '#0d1117', border: '1px solid #30363d', color: '#fff', borderRadius: '6px' }}
                >
                  <option value="high_cpu">IF High CPU (&gt;95% for 10m)</option>
                  <option value="low_disk">IF Low Disk Space (&gt;90%)</option>
                  <option value="pod_crash">IF K8s Pod CrashLoopBackOff</option>
                  <option value="container_down">IF Docker Container Stopped</option>
                </select>
              </div>

              <div>
                <label style={{ fontSize: '12px', color: '#8b949e', display: 'block', marginBottom: '4px' }}>Action Type</label>
                <select
                  value={actionType}
                  onChange={(e) => setActionType(e.target.value)}
                  style={{ width: '100%', padding: '8px 12px', background: '#0d1117', border: '1px solid #30363d', color: '#fff', borderRadius: '6px' }}
                >
                  <option value="restart_service">Restart System Service (SSH)</option>
                  <option value="clean_disk">Clean Temp Files &amp; Vacuum Logs</option>
                  <option value="restart_pod">Restart Kubernetes Deployment</option>
                  <option value="restart_container">Restart Docker Container</option>
                </select>
              </div>

              <div>
                <label style={{ fontSize: '12px', color: '#8b949e', display: 'block', marginBottom: '4px' }}>Target Resource</label>
                <input
                  type="text"
                  value={targetResource}
                  onChange={(e) => setTargetResource(e.target.value)}
                  placeholder="nginx.service / redis-prod / api-deployment"
                  style={{ width: '100%', padding: '8px 12px', background: '#0d1117', border: '1px solid #30363d', color: '#fff', borderRadius: '6px' }}
                />
              </div>

              <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                <input
                  type="checkbox"
                  id="reqApp"
                  checked={requiresApproval}
                  onChange={(e) => setRequiresApproval(e.target.checked)}
                />
                <label htmlFor="reqApp" style={{ fontSize: '13px', color: '#c9d1d9', cursor: 'pointer' }}>
                  Require Admin Approval before execution
                </label>
              </div>

              <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '10px', marginTop: '12px' }}>
                <button type="button" onClick={() => setShowModal(false)} style={{ background: '#21262d', color: '#c9d1d9', border: '1px solid #30363d', padding: '8px 14px', borderRadius: '6px', cursor: 'pointer' }}>Cancel</button>
                <button type="submit" style={{ background: '#238636', color: '#fff', border: 'none', padding: '8px 14px', borderRadius: '6px', cursor: 'pointer', fontWeight: 600 }}>Create Rule</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
