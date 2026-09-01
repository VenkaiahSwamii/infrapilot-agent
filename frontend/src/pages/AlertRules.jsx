import React, { useCallback, useEffect, useState } from 'react';
import {
  Plus,
  Edit2,
  Trash2,
  ToggleLeft,
  ToggleRight,
  AlertTriangle,
  RefreshCw,
  X,
  CheckCircle,
} from 'lucide-react';
import { apiClient } from '../api/client.js';

const METRIC_OPTIONS = [
  { value: 'cpu_usage', label: 'CPU Usage' },
  { value: 'memory_percent', label: 'Memory %' },
  { value: 'disk_percent', label: 'Disk %' },
  { value: 'latency_ms', label: 'Latency (ms)' },
  { value: 'cpu_temperature', label: 'CPU Temperature' },
  { value: 'packet_loss', label: 'Packet Loss' },
  { value: 'disk_read_bps', label: 'Disk Read (Bps)' },
  { value: 'disk_write_bps', label: 'Disk Write (Bps)' },
  { value: 'upload_mbps', label: 'Upload (Mbps)' },
  { value: 'download_mbps', label: 'Download (Mbps)' },
];

const OPERATOR_OPTIONS = [
  { value: '>', label: '>' },
  { value: '>=', label: '>=' },
  { value: '<', label: '<' },
  { value: '<=', label: '<=' },
  { value: '==', label: '=' },
  { value: '!=', label: '!=' },
];

const SEVERITY_OPTIONS = [
  { value: 'critical', label: 'Critical' },
  { value: 'warning', label: 'Warning' },
  { value: 'info', label: 'Info' },
];

const SEVERITY_COLORS = {
  critical: '#ef4444',
  warning: '#f97316',
  info: '#06b6d4',
};

const EMPTY_FORM = {
  name: '',
  metric: 'cpu_usage',
  operator: '>',
  threshold: '',
  severity: 'warning',
  enabled: true,
};

export default function AlertRulesPage() {
  const [rules, setRules] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showModal, setShowModal] = useState(false);
  const [editingRule, setEditingRule] = useState(null);
  const [form, setForm] = useState(EMPTY_FORM);
  const [saving, setSaving] = useState(false);

  const fetchRules = useCallback(async () => {
    setLoading(true);
    setError('');
    try {
      const data = await apiClient.get('/alert-rules');
      setRules(Array.isArray(data.data) ? data.data : []);
    } catch (err) {
      setError(err.response?.data?.error || err.message);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchRules();
  }, [fetchRules]);

  function openCreate() {
    setEditingRule(null);
    setForm(EMPTY_FORM);
    setShowModal(true);
  }

  function openEdit(rule) {
    setEditingRule(rule);
    setForm({
      name: rule.name || '',
      metric: rule.metric || 'cpu_usage',
      operator: rule.operator || '>',
      threshold: rule.threshold ?? '',
      severity: rule.severity || 'warning',
      enabled: rule.is_enabled ?? true,
    });
    setShowModal(true);
  }

  function closeModal() {
    setShowModal(false);
    setEditingRule(null);
    setForm(EMPTY_FORM);
  }

  async function handleSubmit(e) {
    e.preventDefault();
    setSaving(true);
    setError('');

    const payload = {
      name: form.name,
      metric: form.metric,
      operator: form.operator,
      threshold: Number(form.threshold),
      severity: form.severity,
      enabled: form.enabled,
    };

    try {
      if (editingRule) {
        await apiClient.patch(`/alert-rules/${editingRule.id}`, payload);
      } else {
        await apiClient.post('/alert-rules', payload);
      }
      closeModal();
      fetchRules();
    } catch (err) {
      setError(err.response?.data?.error || err.message);
    } finally {
      setSaving(false);
    }
  }

  async function handleDelete(rule) {
    if (!window.confirm(`Delete alert rule "${rule.name}"?`)) return;
    try {
      await apiClient.delete(`/alert-rules/${rule.id}`);
      fetchRules();
    } catch (err) {
      setError(err.response?.data?.error || err.message);
    }
  }

  async function handleToggle(rule) {
    try {
      await apiClient.patch(`/alert-rules/${rule.id}/toggle`, {
        enabled: !rule.is_enabled,
      });
      fetchRules();
    } catch (err) {
      setError(err.response?.data?.error || err.message);
    }
  }

  return (
    <div className="alert-rules-page">
      <div className="alert-rules-header">
        <div className="alert-rules-title-section">
          <AlertTriangle size={28} />
          <h1>Alert Rules</h1>
          <span className="rules-count">{rules.length}</span>
        </div>
        <div className="alert-rules-actions">
          <button className="refresh-btn" onClick={fetchRules} disabled={loading}>
            <RefreshCw size={16} className={loading ? 'spin' : ''} />
            Refresh
          </button>
          <button className="primary-btn" onClick={openCreate}>
            <Plus size={16} />
            New Rule
          </button>
        </div>
      </div>

      {error && <div className="alert-rules-error">{error}</div>}

      {loading ? (
        <div className="alert-rules-loading">
          <div className="spinner" />
          <p>Loading alert rules...</p>
        </div>
      ) : rules.length === 0 ? (
        <div className="alert-rules-empty">
          <CheckCircle size={48} />
          <h3>No alert rules</h3>
          <p>Create rules to define when alerts should be generated.</p>
        </div>
      ) : (
        <div className="alert-rules-table-wrapper">
          <table className="alert-rules-table">
            <thead>
              <tr>
                <th>Name</th>
                <th>Metric</th>
                <th>Operator</th>
                <th>Threshold</th>
                <th>Severity</th>
                <th>Status</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {rules.map((rule) => (
                <tr key={rule.id} className={`rule-${rule.severity}`}>
                  <td className="rule-name">{rule.name}</td>
                  <td>{rule.metric}</td>
                  <td>
                    <span className="operator-badge">{rule.operator}</span>
                  </td>
                  <td>{rule.threshold}</td>
                  <td>
                    <span
                      className="severity-badge"
                      style={{
                        color: SEVERITY_COLORS[rule.severity] || '#94a3b8',
                        borderColor: SEVERITY_COLORS[rule.severity] || '#94a3b8',
                      }}
                    >
                      {rule.severity}
                    </span>
                  </td>
                  <td>
                    <span className={`status-badge ${rule.is_enabled ? 'enabled' : 'disabled'}`}>
                      {rule.is_enabled ? 'Enabled' : 'Disabled'}
                    </span>
                  </td>
                  <td>
                    <div className="rule-actions">
                      <button
                        className="icon-btn"
                        title={rule.is_enabled ? 'Disable' : 'Enable'}
                        onClick={() => handleToggle(rule)}
                      >
                        {rule.is_enabled ? <ToggleRight size={18} /> : <ToggleLeft size={18} />}
                      </button>
                      <button className="icon-btn" title="Edit" onClick={() => openEdit(rule)}>
                        <Edit2 size={16} />
                      </button>
                      <button className="icon-btn danger" title="Delete" onClick={() => handleDelete(rule)}>
                        <Trash2 size={16} />
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {showModal && (
        <div className="modal-overlay">
          <div className="modal">
            <div className="modal-header">
              <h2>{editingRule ? 'Edit Alert Rule' : 'New Alert Rule'}</h2>
              <button className="icon-btn" onClick={closeModal}>
                <X size={18} />
              </button>
            </div>
            <form onSubmit={handleSubmit}>
              <div className="form-group">
                <label>Name</label>
                <input
                  type="text"
                  value={form.name}
                  onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
                  placeholder="e.g. High CPU"
                  required
                />
              </div>
              <div className="form-row">
                <div className="form-group">
                  <label>Metric</label>
                  <select
                    value={form.metric}
                    onChange={(e) => setForm((f) => ({ ...f, metric: e.target.value }))}
                  >
                    {METRIC_OPTIONS.map((m) => (
                      <option key={m.value} value={m.value}>
                        {m.label}
                      </option>
                    ))}
                  </select>
                </div>
                <div className="form-group">
                  <label>Operator</label>
                  <select
                    value={form.operator}
                    onChange={(e) => setForm((f) => ({ ...f, operator: e.target.value }))}
                  >
                    {OPERATOR_OPTIONS.map((o) => (
                      <option key={o.value} value={o.value}>
                        {o.label}
                      </option>
                    ))}
                  </select>
                </div>
                <div className="form-group">
                  <label>Threshold</label>
                  <input
                    type="number"
                    step="any"
                    value={form.threshold}
                    onChange={(e) => setForm((f) => ({ ...f, threshold: e.target.value }))}
                    required
                  />
                </div>
              </div>
              <div className="form-group">
                <label>Severity</label>
                <select
                  value={form.severity}
                  onChange={(e) => setForm((f) => ({ ...f, severity: e.target.value }))}
                >
                  {SEVERITY_OPTIONS.map((s) => (
                    <option key={s.value} value={s.value}>
                      {s.label}
                    </option>
                  ))}
                </select>
              </div>
              <div className="form-actions">
                <button type="button" className="secondary-btn" onClick={closeModal} disabled={saving}>
                  Cancel
                </button>
                <button type="submit" className="primary-btn" disabled={saving}>
                  {saving ? 'Saving...' : editingRule ? 'Update Rule' : 'Create Rule'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      <style>{`
        .alert-rules-page {
          padding: 24px;
          max-width: 1200px;
          margin: 0 auto;
          font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
        }

        .alert-rules-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 24px;
        }

        .alert-rules-title-section {
          display: flex;
          align-items: center;
          gap: 12px;
          color: #f1f5f9;
        }

        .alert-rules-title-section h1 {
          margin: 0;
          font-size: 24px;
          font-weight: 700;
        }

        .rules-count {
          background: #334155;
          color: #94a3b8;
          padding: 2px 10px;
          border-radius: 12px;
          font-size: 14px;
        }

        .alert-rules-actions {
          display: flex;
          gap: 10px;
        }

        .refresh-btn {
          display: flex;
          align-items: center;
          gap: 6px;
          padding: 8px 16px;
          background: #334155;
          color: #e2e8f0;
          border: 1px solid #475569;
          border-radius: 8px;
          cursor: pointer;
          font-size: 14px;
        }

        .refresh-btn:hover { background: #475569; }

        .primary-btn {
          display: inline-flex;
          align-items: center;
          gap: 6px;
          padding: 8px 16px;
          background: #2563eb;
          color: #fff;
          border: none;
          border-radius: 8px;
          cursor: pointer;
          font-size: 14px;
        }

        .primary-btn:hover { background: #1d4ed8; }

        .secondary-btn {
          display: inline-flex;
          align-items: center;
          gap: 6px;
          padding: 8px 16px;
          background: #334155;
          color: #e2e8f0;
          border: 1px solid #475569;
          border-radius: 8px;
          cursor: pointer;
          font-size: 14px;
        }

        .secondary-btn:hover { background: #475569; }

        .alert-rules-error {
          background: #7f1d1d;
          color: #fecaca;
          padding: 10px 14px;
          border-radius: 8px;
          margin-bottom: 16px;
        }

        .alert-rules-loading,
        .alert-rules-empty {
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          padding: 60px 20px;
          color: #94a3b8;
          gap: 12px;
        }

        .alert-rules-table-wrapper {
          background: #0f172a;
          border: 1px solid #334155;
          border-radius: 12px;
          overflow: hidden;
        }

        .alert-rules-table {
          width: 100%;
          border-collapse: collapse;
          color: #e2e8f0;
        }

        .alert-rules-table th,
        .alert-rules-table td {
          padding: 12px 16px;
          text-align: left;
          border-bottom: 1px solid #1e293b;
        }

        .alert-rules-table th {
          background: #111827;
          color: #94a3b8;
          font-weight: 600;
          font-size: 12px;
          text-transform: uppercase;
          letter-spacing: 0.05em;
        }

        .alert-rules-table tr:hover {
          background: #111827;
        }

        .rule-name {
          font-weight: 600;
          color: #f8fafc;
        }

        .operator-badge {
          display: inline-block;
          padding: 2px 8px;
          background: #1e293b;
          border: 1px solid #334155;
          border-radius: 6px;
          font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
          color: #e2e8f0;
        }

        .severity-badge {
          display: inline-block;
          padding: 2px 8px;
          border-radius: 6px;
          background: #111827;
          border: 1px solid;
          font-weight: 600;
          text-transform: uppercase;
          font-size: 12px;
        }

        .status-badge {
          display: inline-flex;
          padding: 2px 8px;
          border-radius: 999px;
          font-size: 12px;
          font-weight: 600;
          text-transform: uppercase;
        }

        .status-badge.enabled {
          background: #14532d;
          color: #4ade80;
        }

        .status-badge.disabled {
          background: #3f3f46;
          color: #a1a1aa;
        }

        .rule-actions {
          display: flex;
          gap: 6px;
        }

        .icon-btn {
          display: inline-flex;
          align-items: center;
          justify-content: center;
          gap: 6px;
          padding: 6px 10px;
          background: transparent;
          color: #cbd5e1;
          border: 1px solid #334155;
          border-radius: 6px;
          cursor: pointer;
        }

        .icon-btn:hover {
          background: #1e293b;
          color: #f8fafc;
        }

        .icon-btn.danger {
          color: #f87171;
          border-color: #7f1d1d;
        }

        .icon-btn.danger:hover {
          background: #7f1d1d;
          color: #fecaca;
        }

        .modal-overlay {
          position: fixed;
          inset: 0;
          background: rgba(0, 0, 0, 0.6);
          display: flex;
          align-items: center;
          justify-content: center;
          z-index: 1000;
        }

        .modal {
          background: #0f172a;
          border: 1px solid #334155;
          border-radius: 12px;
          width: 100%;
          max-width: 520px;
          padding: 20px;
          color: #e2e8f0;
        }

        .modal-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 16px;
        }

        .modal-header h2 {
          margin: 0;
          font-size: 18px;
        }

        .form-group {
          display: flex;
          flex-direction: column;
          gap: 6px;
          margin-bottom: 14px;
        }

        .form-row {
          display: grid;
          grid-template-columns: 1fr 1fr 1fr;
          gap: 12px;
        }

        .form-group label {
          font-size: 12px;
          color: #94a3b8;
          text-transform: uppercase;
          font-weight: 600;
        }

        .form-group input,
        .form-group select {
          background: #020617;
          color: #f8fafc;
          border: 1px solid #334155;
          border-radius: 8px;
          padding: 8px 10px;
          font-size: 14px;
        }

        .form-actions {
          display: flex;
          justify-content: flex-end;
          gap: 10px;
          margin-top: 6px;
        }
      `}</style>
    </div>
  );
}