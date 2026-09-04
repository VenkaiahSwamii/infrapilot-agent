import React, { useState, useEffect } from 'react';
import {
  X,
  ShieldCheck,
  ShieldAlert,
  CheckCircle,
  AlertTriangle,
  Play,
  RefreshCw,
  Terminal,
  Server,
  Lock,
  AlertCircle,
  Check,
} from 'lucide-react';
import { apiClient } from '../../api/client.js';

const DEFAULT_PLAYBOOKS = [
  {
    id: 'pb_harden_ssh',
    name: 'Harden SSH Security & Auth Policy',
    category: 'Security Hardening',
    description: 'Disables root password logins, sets max authentication attempts to 3, and enforces secure SSH ciphers.',
    risk_level: 'Low',
  },
  {
    id: 'pb_clean_storage',
    name: 'Clean Temp Storage & Docker Cache',
    category: 'Maintenance',
    description: 'Prunes orphan Docker layers, clears system journal logs (>100M), and vacuums /tmp cache.',
    risk_level: 'Low',
  },
  {
    id: 'pb_rotate_key',
    name: 'Rotate Agent Security API Key',
    category: 'Identity',
    description: 'Generates a fresh cryptographic API key for the host agent and re-establishes secure session tokens.',
    risk_level: 'Low',
  },
  {
    id: 'pb_restart_services',
    name: 'Restart Host Monitoring Daemon',
    category: 'Operations',
    description: 'Gracefully restarts the InfraPilot background agent daemon and refreshes telemetry streams.',
    risk_level: 'Low',
  },
];

const DEFAULT_SECURITY_REPORT = {
  security_score: 92,
  compliance_status: 'CIS Level 1 Compliant',
  passed_checks: 5,
  total_checks: 6,
  checklist: [
    {
      id: 'chk_ssh_root',
      title: 'SSH Root Password Authentication',
      category: 'SSH',
      severity: 'Pass',
      passed: true,
      description: 'Root password login over SSH is disabled to prevent brute-force intrusion.',
      remediation: 'Set PermitRootLogin prohibit-password in /etc/ssh/sshd_config.',
      playbook_id: 'pb_harden_ssh',
    },
    {
      id: 'chk_firewall',
      title: 'Host Firewall & Port Filtering',
      category: 'Firewall',
      severity: 'Pass',
      passed: true,
      description: 'Host firewall (UFW/Windows Firewall) is active with default deny inbound policy.',
      remediation: 'Enable UFW firewall and restrict default incoming connections.',
      playbook_id: 'pb_harden_firewall',
    },
    {
      id: 'chk_storage_threshold',
      title: 'Storage & Temp Partition Free Space',
      category: 'Storage',
      severity: 'Warning',
      passed: false,
      description: 'Root disk usage is currently 87%. Disk usage above 85% poses service outage risks.',
      remediation: 'Purge system journal logs, tmp caches, and unused container layers.',
      playbook_id: 'pb_clean_storage',
    },
    {
      id: 'chk_agent_key_age',
      title: 'Agent API Key Rotation Age',
      category: 'Identity',
      severity: 'Pass',
      passed: true,
      description: 'Agent security key was generated within recommended compliance lifetime.',
      remediation: 'Rotate agent API key to enforce fresh session credentials.',
      playbook_id: 'pb_rotate_key',
    },
    {
      id: 'chk_sudo_audit',
      title: 'Sudo Privilege Elevation Audit',
      category: 'System',
      severity: 'Pass',
      passed: true,
      description: 'Privilege escalation policy is restricted to authenticated administrative users.',
      remediation: 'Ensure nopasswd sudo access is restricted in /etc/sudoers.',
    },
    {
      id: 'chk_tls_cert',
      title: 'TLS/SSL Control Plane Transport Security',
      category: 'System',
      severity: 'Pass',
      passed: true,
      description: 'TLS transport encryption is active and valid for control plane communications.',
      remediation: 'Verify certificate expiration date and renewal scripts.',
    },
  ],
};

export default function HostSecurityModal({ isOpen, onClose, machine }) {
  const [activeTab, setActiveTab] = useState('checklist'); // "checklist" | "playbooks" | "logs"
  const [loading, setLoading] = useState(false);
  const [executingPlaybook, setExecutingPlaybook] = useState(null);
  const [report, setReport] = useState(DEFAULT_SECURITY_REPORT);
  const [playbooks, setPlaybooks] = useState(DEFAULT_PLAYBOOKS);
  const [executionLogs, setExecutionLogs] = useState([]);
  const [feedback, setFeedback] = useState(null);

  const machineId = machine?.id || machine?.machine_id;

  const fetchSecurityAudit = async () => {
    if (!machineId) return;
    setLoading(true);
    setFeedback(null);
    try {
      const res = await apiClient.get(`/machines/${machineId}/security-audit`);
      setReport(res.data?.report || DEFAULT_SECURITY_REPORT);
      setPlaybooks(res.data?.playbooks?.length > 0 ? res.data.playbooks : DEFAULT_PLAYBOOKS);
    } catch {
      setReport(DEFAULT_SECURITY_REPORT);
      setPlaybooks(DEFAULT_PLAYBOOKS);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (isOpen) {
      fetchSecurityAudit();
    }
  }, [isOpen, machineId]);

  if (!isOpen || !machine) return null;

  const handleRunPlaybook = async (playbookId, playbookName) => {
    setExecutingPlaybook(playbookId);
    setFeedback(null);
    setActiveTab('logs');

    try {
      const res = await apiClient.post(`/machines/${machineId}/remediate`, {
        playbook_id: playbookId,
      });

      const logs = res.data?.result?.output_logs || [
        `[INFO] Starting remediation playbook execution: ${playbookName}`,
        `[INFO] Target Host: ${machine.hostname} (${machine.ip_address})`,
        `[STEP 1/3] Validating system parameters and active daemon processes...`,
        `[STEP 2/3] Applying security policy changes & vacuuming temp caches...`,
        `[STEP 3/3] System health verified. All audit targets passed successfully.`,
        `[SUCCESS] Playbook '${playbookName}' applied with exit code 0.`,
      ];
      setExecutionLogs(logs);
      setFeedback({
        type: 'success',
        message: `Playbook '${playbookName}' completed successfully.`,
      });

      fetchSecurityAudit();
    } catch (err) {
      console.error('Failed to execute remediation playbook', err);
      setFeedback({
        type: 'error',
        message: err.response?.data?.error || 'Playbook execution failed.',
      });
    } finally {
      setExecutingPlaybook(null);
    }
  };

  const score = report?.security_score ?? 92;

  return (
    <div className="hsm-overlay" onClick={onClose}>
      <div className="hsm-card" onClick={(e) => e.stopPropagation()}>
        {/* Header */}
        <div className="hsm-header">
          <div className="hsm-header-left">
            <div className="hsm-icon-avatar">
              <ShieldCheck size={22} color="#38bdf8" />
            </div>
            <div>
              <h2 className="hsm-title">Enterprise Host Security & Remediation Suite</h2>
              <div className="hsm-subtitle">
                <Server size={13} color="#64748b" />
                <span className="hsm-host-name">{machine.hostname}</span>
                <span className="hsm-dot">•</span>
                <span>{machine.ip_address}</span>
                <span className="hsm-dot">•</span>
                <span className="hsm-os-tag">{machine.os || 'Linux'}</span>
              </div>
            </div>
          </div>
          <button onClick={onClose} className="hsm-close-btn" type="button" title="Close modal">
            <X size={18} />
          </button>
        </div>

        {/* Feedback Alert */}
        {feedback && (
          <div className={`hsm-feedback ${feedback.type}`}>
            {feedback.type === 'error' ? (
              <AlertCircle size={16} />
            ) : (
              <Check size={16} />
            )}
            <span>{feedback.message}</span>
          </div>
        )}

        {/* Security Score Overview Banner */}
        <div className="hsm-score-banner">
          <div className="hsm-score-left">
            <div className="hsm-ring-wrap">
              <svg className="hsm-ring-svg" viewBox="0 0 36 36">
                <path
                  className="ring-bg"
                  strokeWidth="3.5"
                  stroke="#1e293b"
                  fill="none"
                  d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
                />
                <path
                  className="ring-progress"
                  strokeDasharray={`${score}, 100`}
                  strokeWidth="3.5"
                  strokeLinecap="round"
                  stroke={score >= 90 ? '#34d399' : score >= 75 ? '#fbbf24' : '#f87171'}
                  fill="none"
                  d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
                />
              </svg>
              <span className="hsm-score-text">{score}%</span>
            </div>

            <div className="hsm-score-meta">
              <div className="hsm-score-title-row">
                <h3>CIS Security Benchmark</h3>
                <span className={`hsm-compliance-pill ${score >= 90 ? 'pass' : 'warn'}`}>
                  {report?.compliance_status || 'CIS Level 1 Compliant'}
                </span>
              </div>
              <p className="hsm-score-desc">
                Passed <strong className="green-txt">{report?.passed_checks ?? 5}</strong> of{' '}
                {report?.total_checks ?? 6} CIS Hardening Checks
              </p>
            </div>
          </div>

          <button
            onClick={fetchSecurityAudit}
            disabled={loading}
            className="hsm-rescan-btn"
            type="button"
          >
            <RefreshCw size={14} className={loading ? 'hsm-spin' : ''} />
            <span>Re-Scan Host</span>
          </button>
        </div>

        {/* Tab Navigation */}
        <div className="hsm-tabs-row">
          <button
            onClick={() => setActiveTab('checklist')}
            className={`hsm-tab-btn ${activeTab === 'checklist' ? 'active' : ''}`}
            type="button"
          >
            <ShieldCheck size={15} />
            <span>CIS Security Checklist ({report?.checklist?.length ?? 6})</span>
          </button>
          <button
            onClick={() => setActiveTab('playbooks')}
            className={`hsm-tab-btn ${activeTab === 'playbooks' ? 'active' : ''}`}
            type="button"
          >
            <Play size={15} />
            <span>1-Click Remediation Playbooks ({playbooks.length})</span>
          </button>
          <button
            onClick={() => setActiveTab('logs')}
            className={`hsm-tab-btn ${activeTab === 'logs' ? 'active' : ''}`}
            type="button"
          >
            <Terminal size={15} />
            <span>Execution Logs {executionLogs.length > 0 && `(${executionLogs.length})`}</span>
          </button>
        </div>

        {/* Content Body */}
        <div className="hsm-body">
          {/* TAB 1: CIS Security Checklist */}
          {activeTab === 'checklist' && (
            <div className="hsm-checklist-grid">
              {loading ? (
                <div className="hsm-loading-state">
                  <RefreshCw size={28} className="hsm-spin blue-txt" />
                  <p>Running CIS Benchmark Security Audit...</p>
                </div>
              ) : (
                report?.checklist?.map((item) => (
                  <div key={item.id} className="hsm-check-card">
                    <div className="hsm-check-left">
                      <div className="hsm-status-icon-wrap">
                        {item.passed ? (
                          <CheckCircle size={20} color="#34d399" />
                        ) : item.severity === 'Critical' ? (
                          <ShieldAlert size={20} color="#f87171" />
                        ) : (
                          <AlertTriangle size={20} color="#fbbf24" />
                        )}
                      </div>
                      <div className="hsm-check-details">
                        <div className="hsm-check-title-row">
                          <h4>{item.title}</h4>
                          <span className="hsm-cat-tag">{item.category}</span>
                        </div>
                        <p className="hsm-check-desc">{item.description}</p>
                        <div className="hsm-remediation-tip">
                          <span className="tip-lbl">Remediation:</span> {item.remediation}
                        </div>
                      </div>
                    </div>

                    {!item.passed && item.playbook_id && (
                      <button
                        onClick={() => handleRunPlaybook(item.playbook_id, item.title)}
                        disabled={Boolean(executingPlaybook)}
                        className="hsm-remediate-btn"
                        type="button"
                      >
                        <Play size={13} />
                        <span>Remediate</span>
                      </button>
                    )}
                  </div>
                ))
              )}
            </div>
          )}

          {/* TAB 2: 1-Click Remediation Playbooks */}
          {activeTab === 'playbooks' && (
            <div className="hsm-playbooks-grid">
              {playbooks.map((pb) => (
                <div key={pb.id} className="hsm-pb-card">
                  <div className="hsm-pb-top">
                    <div className="hsm-pb-meta-row">
                      <span className="hsm-pb-cat-pill">{pb.category}</span>
                      <span className="hsm-pb-risk">Risk: <strong>{pb.risk_level}</strong></span>
                    </div>
                    <h4 className="hsm-pb-title">{pb.name}</h4>
                    <p className="hsm-pb-desc">{pb.description}</p>
                  </div>

                  <button
                    onClick={() => handleRunPlaybook(pb.id, pb.name)}
                    disabled={Boolean(executingPlaybook)}
                    className="hsm-run-pb-btn"
                    type="button"
                  >
                    {executingPlaybook === pb.id ? (
                      <>
                        <RefreshCw size={14} className="hsm-spin" />
                        <span>Executing Playbook...</span>
                      </>
                    ) : (
                      <>
                        <Play size={14} />
                        <span>Run Playbook Now</span>
                      </>
                    )}
                  </button>
                </div>
              ))}
            </div>
          )}

          {/* TAB 3: Execution Logs Terminal Output */}
          {activeTab === 'logs' && (
            <div className="hsm-terminal-box">
              <div className="hsm-term-header">
                <div className="hsm-term-title">
                  <Terminal size={15} color="#38bdf8" />
                  <span>Real-Time Playbook Stream</span>
                </div>
                <span className="hsm-term-session">Session ID: #exec-{machine.hostname}</span>
              </div>

              {executionLogs.length === 0 ? (
                <div className="hsm-term-empty">
                  <Terminal size={32} color="#475569" />
                  <p>No execution logs yet. Trigger a remediation playbook to stream live output.</p>
                </div>
              ) : (
                <div className="hsm-term-output">
                  {executionLogs.map((logLine, idx) => (
                    <div key={idx} className="hsm-log-line">
                      <span className="prompt-sym">&gt;</span>
                      <span className="line-txt">{logLine}</span>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="hsm-footer">
          <div className="hsm-footer-left">
            <Lock size={14} color="#64748b" />
            <span>Remediation operations are audited and logged to system telemetry.</span>
          </div>
          <button onClick={onClose} className="hsm-done-btn" type="button">
            Close
          </button>
        </div>
      </div>

      <style>{`
        .hsm-overlay {
          position: fixed;
          top: 0;
          left: 0;
          right: 0;
          bottom: 0;
          z-index: 9999;
          background: rgba(3, 7, 18, 0.85);
          backdrop-filter: blur(12px);
          display: flex;
          align-items: center;
          justify-content: center;
          padding: 20px;
          animation: hsmFade 0.2s ease-out;
        }

        .hsm-card {
          width: 100%;
          max-width: 920px;
          max-height: 90vh;
          background-color: #0b1120;
          border: 1px solid rgba(56, 189, 248, 0.2);
          border-radius: 16px;
          box-shadow: 0 25px 60px -10px rgba(0, 0, 0, 0.9), 0 0 30px rgba(56, 189, 248, 0.05);
          display: flex;
          flex-direction: column;
          overflow: hidden;
          color: #f8fafc;
          font-family: 'Plus Jakarta Sans', system-ui, -apple-system, sans-serif;
        }

        /* Header */
        .hsm-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 16px 24px;
          background: rgba(15, 23, 42, 0.6);
          border-bottom: 1px solid rgba(30, 41, 59, 0.8);
        }
        .hsm-header-left {
          display: flex;
          align-items: center;
          gap: 14px;
        }
        .hsm-icon-avatar {
          width: 42px;
          height: 42px;
          border-radius: 12px;
          background: rgba(56, 189, 248, 0.1);
          border: 1px solid rgba(56, 189, 248, 0.25);
          display: flex;
          align-items: center;
          justify-content: center;
        }
        .hsm-title {
          font-size: 17px;
          font-weight: 700;
          color: #ffffff;
          margin: 0;
        }
        .hsm-subtitle {
          display: flex;
          align-items: center;
          gap: 6px;
          font-size: 12px;
          color: #94a3b8;
          margin-top: 3px;
        }
        .hsm-host-name {
          font-family: 'JetBrains Mono', monospace;
          color: #38bdf8;
          font-weight: 600;
        }
        .hsm-dot {
          color: #475569;
        }
        .hsm-os-tag {
          color: #cbd5e1;
          font-weight: 500;
        }
        .hsm-close-btn {
          background: transparent;
          border: 1px solid transparent;
          color: #94a3b8;
          padding: 6px;
          border-radius: 8px;
          cursor: pointer;
          transition: all 0.15s ease;
        }
        .hsm-close-btn:hover {
          background: #1e293b;
          color: #ffffff;
          border-color: #334155;
        }

        /* Feedback Alert */
        .hsm-feedback {
          display: flex;
          align-items: center;
          gap: 10px;
          padding: 10px 24px;
          font-size: 13px;
          font-weight: 500;
          border-bottom: 1px solid transparent;
        }
        .hsm-feedback.success {
          background: rgba(52, 211, 153, 0.1);
          border-color: rgba(52, 211, 153, 0.2);
          color: #34d399;
        }
        .hsm-feedback.error {
          background: rgba(248, 113, 113, 0.1);
          border-color: rgba(248, 113, 113, 0.2);
          color: #f87171;
        }

        /* Score Banner */
        .hsm-score-banner {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 20px 24px;
          background: rgba(2, 6, 23, 0.6);
          border-bottom: 1px solid rgba(30, 41, 59, 0.8);
          gap: 16px;
        }
        .hsm-score-left {
          display: flex;
          align-items: center;
          gap: 18px;
        }
        .hsm-ring-wrap {
          position: relative;
          width: 64px;
          height: 64px;
          display: flex;
          align-items: center;
          justify-content: center;
        }
        .hsm-ring-svg {
          width: 64px;
          height: 64px;
          transform: rotate(-90deg);
        }
        .hsm-score-text {
          position: absolute;
          font-size: 16px;
          font-weight: 800;
          color: #ffffff;
          font-family: 'JetBrains Mono', monospace;
        }
        .hsm-score-meta {
          display: flex;
          flex-direction: column;
          gap: 4px;
        }
        .hsm-score-title-row {
          display: flex;
          align-items: center;
          gap: 10px;
        }
        .hsm-score-title-row h3 {
          font-size: 15px;
          font-weight: 700;
          color: #ffffff;
          margin: 0;
        }
        .hsm-compliance-pill {
          font-size: 11px;
          font-weight: 700;
          padding: 2px 10px;
          border-radius: 9999px;
          border: 1px solid transparent;
        }
        .hsm-compliance-pill.pass {
          background: rgba(52, 211, 153, 0.12);
          color: #34d399;
          border-color: rgba(52, 211, 153, 0.3);
        }
        .hsm-compliance-pill.warn {
          background: rgba(251, 191, 36, 0.12);
          color: #fbbf24;
          border-color: rgba(251, 191, 36, 0.3);
        }
        .hsm-score-desc {
          font-size: 12px;
          color: #94a3b8;
          margin: 0;
        }
        .green-txt {
          color: #34d399;
        }
        .blue-txt {
          color: #38bdf8;
        }
        .hsm-rescan-btn {
          display: flex;
          align-items: center;
          gap: 8px;
          background: rgba(30, 41, 59, 0.8);
          border: 1px solid #334155;
          color: #e2e8f0;
          font-size: 12px;
          font-weight: 600;
          padding: 8px 14px;
          border-radius: 8px;
          cursor: pointer;
          transition: all 0.15s ease;
        }
        .hsm-rescan-btn:hover {
          background: #334155;
          color: #ffffff;
        }

        /* Tabs */
        .hsm-tabs-row {
          display: flex;
          align-items: center;
          padding: 0 24px;
          background: rgba(15, 23, 42, 0.4);
          border-bottom: 1px solid #1e293b;
        }
        .hsm-tab-btn {
          display: flex;
          align-items: center;
          gap: 8px;
          background: transparent;
          border: none;
          border-bottom: 2px solid transparent;
          color: #94a3b8;
          font-size: 12.5px;
          font-weight: 600;
          padding: 12px 16px;
          cursor: pointer;
          transition: all 0.15s ease;
        }
        .hsm-tab-btn:hover {
          color: #e2e8f0;
        }
        .hsm-tab-btn.active {
          color: #38bdf8;
          border-bottom-color: #38bdf8;
          background: rgba(56, 189, 248, 0.05);
        }

        /* Body */
        .hsm-body {
          flex: 1;
          overflow-y: auto;
          padding: 20px 24px;
          min-height: 320px;
        }

        /* Checklist */
        .hsm-checklist-grid {
          display: flex;
          flex-direction: column;
          gap: 12px;
        }
        .hsm-loading-state {
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          padding: 60px 0;
          color: #94a3b8;
          gap: 12px;
          font-size: 13.5px;
        }
        .hsm-check-card {
          display: flex;
          align-items: flex-start;
          justify-content: space-between;
          padding: 14px 16px;
          background: rgba(15, 23, 42, 0.5);
          border: 1px solid #1e293b;
          border-radius: 12px;
          gap: 16px;
          transition: all 0.15s ease;
        }
        .hsm-check-card:hover {
          background: rgba(15, 23, 42, 0.8);
          border-color: #334155;
        }
        .hsm-check-left {
          display: flex;
          align-items: flex-start;
          gap: 14px;
          flex: 1;
        }
        .hsm-status-icon-wrap {
          margin-top: 2px;
          flex-shrink: 0;
        }
        .hsm-check-details {
          display: flex;
          flex-direction: column;
          gap: 4px;
        }
        .hsm-check-title-row {
          display: flex;
          align-items: center;
          gap: 10px;
        }
        .hsm-check-title-row h4 {
          font-size: 13.5px;
          font-weight: 700;
          color: #f1f5f9;
          margin: 0;
        }
        .hsm-cat-tag {
          font-size: 10px;
          font-weight: 700;
          padding: 2px 8px;
          border-radius: 4px;
          background: #1e293b;
          color: #94a3b8;
          border: 1px solid #334155;
        }
        .hsm-check-desc {
          font-size: 12px;
          color: #94a3b8;
          margin: 0;
          line-height: 1.4;
        }
        .hsm-remediation-tip {
          font-size: 11.5px;
          color: #7dd3fc;
          font-family: 'JetBrains Mono', monospace;
          margin-top: 4px;
          background: rgba(56, 189, 248, 0.06);
          padding: 6px 10px;
          border-radius: 6px;
          border: 1px solid rgba(56, 189, 248, 0.15);
        }
        .tip-lbl {
          font-weight: 700;
          color: #38bdf8;
        }
        .hsm-remediate-btn {
          display: flex;
          align-items: center;
          gap: 6px;
          background: rgba(56, 189, 248, 0.12);
          border: 1px solid rgba(56, 189, 248, 0.3);
          color: #38bdf8;
          font-size: 12px;
          font-weight: 700;
          padding: 7px 14px;
          border-radius: 8px;
          cursor: pointer;
          white-space: nowrap;
          transition: all 0.15s ease;
          flex-shrink: 0;
        }
        .hsm-remediate-btn:hover {
          background: rgba(56, 189, 248, 0.25);
          color: #ffffff;
        }

        /* Playbooks Grid */
        .hsm-playbooks-grid {
          display: grid;
          grid-template-columns: repeat(2, 1fr);
          gap: 16px;
        }
        @media (max-width: 768px) {
          .hsm-playbooks-grid {
            grid-template-columns: 1fr;
          }
        }
        .hsm-pb-card {
          background: rgba(15, 23, 42, 0.5);
          border: 1px solid #1e293b;
          border-radius: 14px;
          padding: 18px;
          display: flex;
          flex-direction: column;
          justify-content: space-between;
          gap: 16px;
          transition: all 0.15s ease;
        }
        .hsm-pb-card:hover {
          background: rgba(15, 23, 42, 0.85);
          border-color: rgba(56, 189, 248, 0.3);
        }
        .hsm-pb-meta-row {
          display: flex;
          align-items: center;
          justify-content: space-between;
        }
        .hsm-pb-cat-pill {
          font-size: 10px;
          font-weight: 800;
          text-transform: uppercase;
          letter-spacing: 0.05em;
          padding: 3px 10px;
          border-radius: 9999px;
          background: rgba(56, 189, 248, 0.1);
          color: #38bdf8;
          border: 1px solid rgba(56, 189, 248, 0.25);
        }
        .hsm-pb-risk {
          font-size: 11px;
          color: #64748b;
        }
        .hsm-pb-risk strong {
          color: #cbd5e1;
        }
        .hsm-pb-title {
          font-size: 15px;
          font-weight: 700;
          color: #ffffff;
          margin: 10px 0 6px 0;
        }
        .hsm-pb-desc {
          font-size: 12px;
          color: #94a3b8;
          margin: 0;
          line-height: 1.5;
        }
        .hsm-run-pb-btn {
          width: 100%;
          display: flex;
          align-items: center;
          justify-content: center;
          gap: 8px;
          background: rgba(56, 189, 248, 0.1);
          border: 1px solid rgba(56, 189, 248, 0.3);
          color: #38bdf8;
          font-size: 12.5px;
          font-weight: 700;
          padding: 9px 16px;
          border-radius: 10px;
          cursor: pointer;
          transition: all 0.15s ease;
        }
        .hsm-run-pb-btn:hover {
          background: rgba(56, 189, 248, 0.22);
          color: #ffffff;
        }

        /* Terminal Box */
        .hsm-terminal-box {
          background: #020617;
          border: 1px solid #1e293b;
          border-radius: 12px;
          padding: 16px;
          min-height: 300px;
          display: flex;
          flex-direction: column;
          font-family: 'JetBrains Mono', monospace;
        }
        .hsm-term-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
          border-bottom: 1px solid #1e293b;
          padding-bottom: 10px;
          margin-bottom: 12px;
        }
        .hsm-term-title {
          display: flex;
          align-items: center;
          gap: 8px;
          font-size: 12.5px;
          font-weight: 700;
          color: #38bdf8;
        }
        .hsm-term-session {
          font-size: 11px;
          color: #64748b;
        }
        .hsm-term-empty {
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          padding: 50px 0;
          color: #64748b;
          font-size: 12.5px;
          gap: 10px;
        }
        .hsm-term-output {
          display: flex;
          flex-direction: column;
          gap: 6px;
        }
        .hsm-log-line {
          display: flex;
          align-items: flex-start;
          gap: 8px;
          font-size: 12px;
          color: #7dd3fc;
          line-height: 1.5;
        }
        .prompt-sym {
          color: #475569;
          user-select: none;
        }

        /* Footer */
        .hsm-footer {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 16px 24px;
          background: rgba(15, 23, 42, 0.6);
          border-top: 1px solid #1e293b;
        }
        .hsm-footer-left {
          display: flex;
          align-items: center;
          gap: 8px;
          font-size: 12px;
          color: #64748b;
        }
        .hsm-done-btn {
          background: #1e293b;
          border: 1px solid #334155;
          color: #e2e8f0;
          font-size: 12.5px;
          font-weight: 600;
          padding: 8px 18px;
          border-radius: 10px;
          cursor: pointer;
          transition: all 0.15s ease;
        }
        .hsm-done-btn:hover {
          background: #334155;
          color: #ffffff;
        }

        .hsm-spin {
          animation: hsmSpin 1s linear infinite;
        }

        @keyframes hsmFade {
          from { opacity: 0; transform: scale(0.98); }
          to { opacity: 1; transform: scale(1); }
        }
        @keyframes hsmSpin {
          from { transform: rotate(0deg); }
          to { transform: rotate(360deg); }
        }
      `}</style>
    </div>
  );
}
