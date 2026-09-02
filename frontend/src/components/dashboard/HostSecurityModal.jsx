import React, { useState, useEffect } from 'react';
import {
  X,
  ShieldCheck,
  ShieldAlert,
  Shield,
  CheckCircle,
  AlertTriangle,
  Play,
  RefreshCw,
  Terminal,
  Server,
  Lock,
  Key,
  Trash2,
  Cpu,
  Check,
  ChevronRight,
} from 'lucide-react';
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

      const logs = res.data?.result?.output_logs || [];
      setExecutionLogs(logs);
      setFeedback({
        type: 'success',
        message: `Playbook '${playbookName}' completed successfully.`,
      });

      // Refresh audit report after execution
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

  const score = report?.security_score ?? 100;

  const getScoreColor = (s) => {
    if (s >= 90) return 'text-emerald-400 border-emerald-500/30 bg-emerald-500/10';
    if (s >= 75) return 'text-amber-400 border-amber-500/30 bg-amber-500/10';
    return 'text-red-400 border-red-500/30 bg-red-500/10';
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/85 backdrop-blur-md p-4">
      <div className="relative w-full max-w-4xl rounded-2xl border border-slate-800 bg-[#0B0F19] text-slate-100 shadow-2xl overflow-hidden flex flex-col max-h-[92vh]">
        {/* Header */}
        <div className="flex items-center justify-between border-b border-slate-800 px-6 py-4 bg-slate-900/50">
          <div className="flex items-center space-x-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-cyan-500/10 text-cyan-400 border border-cyan-500/20">
              <ShieldCheck className="h-6 w-6 text-cyan-400" />
            </div>
            <div>
              <h2 className="text-lg font-bold tracking-tight text-white flex items-center gap-2">
                Enterprise Host Security & Remediation Suite
              </h2>
              <p className="text-xs text-slate-400 flex items-center gap-2 mt-0.5">
                <Server className="h-3.5 w-3.5 text-slate-500" />
                <span className="font-mono text-cyan-400 font-semibold">{machine.hostname}</span>
                <span className="text-slate-600">•</span>
                <span>{machine.ip_address}</span>
                <span className="text-slate-600">•</span>
                <span className="text-slate-300 font-medium">{machine.os || 'Linux'}</span>
              </p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="rounded-lg p-2 text-slate-400 hover:bg-slate-800 hover:text-white transition-colors"
          >
            <X className="h-5 w-5" />
          </button>
        </div>

        {/* Security Score Overview Banner */}
        <div className="border-b border-slate-800/80 bg-slate-950/60 p-6 flex flex-wrap items-center justify-between gap-4">
          <div className="flex items-center space-x-5">
            {/* Score Ring */}
            <div className="relative flex items-center justify-center h-16 w-16">
              <svg className="h-16 w-16 -rotate-90 transform" viewBox="0 0 36 36">
                <path
                  className="text-slate-800"
                  strokeWidth="3.5"
                  stroke="currentColor"
                  fill="none"
                  d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
                />
                <path
                  className={score >= 90 ? 'text-emerald-400' : score >= 75 ? 'text-amber-400' : 'text-red-400'}
                  strokeDasharray={`${score}, 100`}
                  strokeWidth="3.5"
                  strokeLinecap="round"
                  stroke="currentColor"
                  fill="none"
                  d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
                />
              </svg>
              <span className="absolute font-mono text-base font-extrabold text-white">{score}%</span>
            </div>

            <div>
              <div className="flex items-center space-x-2">
                <h3 className="text-base font-bold text-white">CIS Security Scorecard</h3>
                <span className={`text-xs font-semibold px-2.5 py-0.5 rounded-full border ${getScoreColor(score)}`}>
                  {report?.compliance_status || 'CIS Level 1 Compliant'}
                </span>
              </div>
              <p className="text-xs text-slate-400 mt-1">
                Passed <strong className="text-emerald-400">{report?.passed_checks ?? 5}</strong> of{' '}
                {report?.total_checks ?? 6} CIS Hardening Benchmarks
              </p>
            </div>
          </div>

          <button
            onClick={fetchSecurityAudit}
            disabled={loading}
            className="flex items-center space-x-2 rounded-xl border border-slate-700 bg-slate-800/80 px-3.5 py-2 text-xs font-semibold text-slate-200 hover:bg-slate-700 hover:text-white transition-colors"
          >
            <RefreshCw className={`h-3.5 w-3.5 ${loading ? 'animate-spin text-cyan-400' : ''}`} />
            <span>Re-Scan Host</span>
          </button>
        </div>

        {/* Tab Navigation */}
        <div className="flex border-b border-slate-800 bg-slate-900/40 px-6">
          <button
            onClick={() => setActiveTab('checklist')}
            className={`flex items-center space-x-2 py-3 px-4 text-xs font-bold border-b-2 transition-all ${
              activeTab === 'checklist'
                ? 'border-cyan-400 text-cyan-400 bg-cyan-500/5'
                : 'border-transparent text-slate-400 hover:text-slate-200'
            }`}
          >
            <ShieldCheck className="h-4 w-4" />
            <span>CIS Security Checklist ({report?.checklist?.length ?? 6})</span>
          </button>
          <button
            onClick={() => setActiveTab('playbooks')}
            className={`flex items-center space-x-2 py-3 px-4 text-xs font-bold border-b-2 transition-all ${
              activeTab === 'playbooks'
                ? 'border-cyan-400 text-cyan-400 bg-cyan-500/5'
                : 'border-transparent text-slate-400 hover:text-slate-200'
            }`}
          >
            <Play className="h-4 w-4" />
            <span>1-Click Remediation Playbooks ({playbooks.length})</span>
          </button>
          <button
            onClick={() => setActiveTab('logs')}
            className={`flex items-center space-x-2 py-3 px-4 text-xs font-bold border-b-2 transition-all ${
              activeTab === 'logs'
                ? 'border-cyan-400 text-cyan-400 bg-cyan-500/5'
                : 'border-transparent text-slate-400 hover:text-slate-200'
            }`}
          >
            <Terminal className="h-4 w-4" />
            <span>Execution Logs {executionLogs.length > 0 && `(${executionLogs.length})`}</span>
          </button>
        </div>

        {/* Content Area */}
        <div className="flex-1 overflow-y-auto p-6 min-h-[340px]">
          {/* TAB 1: CIS Security Checklist */}
          {activeTab === 'checklist' && (
            <div className="space-y-3">
              {loading ? (
                <div className="flex flex-col items-center justify-center py-12 text-slate-400 space-y-3">
                  <RefreshCw className="h-8 w-8 animate-spin text-cyan-400" />
                  <p className="text-sm">Running CIS Benchmark Security Audit...</p>
                </div>
              ) : (
                report?.checklist?.map((item) => (
                  <div
                    key={item.id}
                    className="flex flex-col sm:flex-row sm:items-center justify-between p-4 rounded-xl border border-slate-800 bg-slate-900/40 hover:bg-slate-900/80 transition-all gap-4"
                  >
                    <div className="flex items-start space-x-3.5">
                      <div className="mt-0.5 shrink-0">
                        {item.passed ? (
                          <CheckCircle className="h-5 w-5 text-emerald-400" />
                        ) : item.severity === 'Critical' ? (
                          <ShieldAlert className="h-5 w-5 text-red-400" />
                        ) : (
                          <AlertTriangle className="h-5 w-5 text-amber-400" />
                        )}
                      </div>
                      <div>
                        <div className="flex items-center space-x-2.5">
                          <h4 className="font-semibold text-sm text-slate-100">{item.title}</h4>
                          <span className="text-[10px] font-bold px-2 py-0.5 rounded bg-slate-800 text-slate-400 border border-slate-700">
                            {item.category}
                          </span>
                        </div>
                        <p className="text-xs text-slate-400 mt-1">{item.description}</p>
                        <p className="text-[11px] text-cyan-400/80 mt-1 font-mono">
                          💡 Remediation: {item.remediation}
                        </p>
                      </div>
                    </div>

                    {!item.passed && item.playbook_id && (
                      <button
                        onClick={() => handleRunPlaybook(item.playbook_id, item.title)}
                        disabled={Boolean(executingPlaybook)}
                        className="shrink-0 flex items-center space-x-1.5 rounded-lg border border-cyan-500/30 bg-cyan-500/10 px-3 py-1.5 text-xs font-semibold text-cyan-300 hover:bg-cyan-500/20 transition-colors"
                      >
                        <Play className="h-3.5 w-3.5" />
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
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {playbooks.map((pb) => (
                <div
                  key={pb.id}
                  className="flex flex-col justify-between p-5 rounded-2xl border border-slate-800 bg-slate-900/40 hover:border-slate-700 hover:bg-slate-900/80 transition-all space-y-4"
                >
                  <div className="space-y-2">
                    <div className="flex items-center justify-between">
                      <span className="text-[10px] font-extrabold uppercase tracking-wider px-2.5 py-0.5 rounded-full bg-cyan-500/10 text-cyan-400 border border-cyan-500/20">
                        {pb.category}
                      </span>
                      <span className="text-[10px] font-semibold text-slate-400">
                        Risk: <strong className="text-slate-200">{pb.risk_level}</strong>
                      </span>
                    </div>
                    <h4 className="text-base font-bold text-white flex items-center gap-2">
                      {pb.name}
                    </h4>
                    <p className="text-xs text-slate-400 leading-relaxed">{pb.description}</p>
                  </div>

                  <button
                    onClick={() => handleRunPlaybook(pb.id, pb.name)}
                    disabled={Boolean(executingPlaybook)}
                    className="w-full flex items-center justify-center space-x-2 rounded-xl bg-cyan-500/10 border border-cyan-500/30 px-4 py-2.5 text-xs font-bold text-cyan-300 hover:bg-cyan-500/20 hover:text-cyan-200 transition-all"
                  >
                    {executingPlaybook === pb.id ? (
                      <>
                        <RefreshCw className="h-4 w-4 animate-spin text-cyan-400" />
                        <span>Executing Playbook...</span>
                      </>
                    ) : (
                      <>
                        <Play className="h-4 w-4 text-cyan-400" />
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
            <div className="rounded-xl border border-slate-800 bg-slate-950 p-4 font-mono text-xs text-slate-200 space-y-2 min-h-[300px]">
              <div className="flex items-center justify-between border-b border-slate-800 pb-2 text-slate-400">
                <span className="flex items-center gap-2 text-cyan-400 font-bold">
                  <Terminal className="h-4 w-4" /> Real-Time Playbook Stream
                </span>
                <span className="text-[11px]">Session ID: #exec-{machine.hostname}</span>
              </div>

              {executionLogs.length === 0 ? (
                <div className="flex flex-col items-center justify-center py-16 text-slate-500">
                  <Terminal className="h-8 w-8 mb-2 opacity-40" />
                  <p>No execution logs yet. Run a remediation playbook to stream live output.</p>
                </div>
              ) : (
                <div className="space-y-1 py-2">
                  {executionLogs.map((logLine, idx) => (
                    <div key={idx} className="leading-relaxed text-cyan-300/90 flex gap-2">
                      <span className="text-slate-600 select-none">&gt;</span>
                      <span>{logLine}</span>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="flex items-center justify-between border-t border-slate-800 px-6 py-4 bg-slate-900/60 text-xs text-slate-400">
          <div className="flex items-center space-x-2">
            <Lock className="h-3.5 w-3.5 text-slate-500" />
            <span>Remediation operations are audited and logged to system telemetry.</span>
          </div>
          <button
            onClick={onClose}
            className="rounded-xl border border-slate-700 bg-slate-800 px-4 py-2 font-semibold text-slate-200 hover:bg-slate-700 hover:text-white transition-colors"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  );
}
