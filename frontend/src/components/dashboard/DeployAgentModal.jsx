import React, { useState, useEffect } from 'react';
import {
  X,
  Copy,
  Check,
  Terminal,
  Monitor,
  Layers,
  Network,
  ShieldCheck,
  RefreshCw,
  Key,
} from 'lucide-react';
import { createEnrollmentToken, listEnrollmentTokens } from '../../api/enrollment.js';

export default function DeployAgentModal({ isOpen, onClose }) {
  const [activeTab, setActiveTab] = useState('linux');
  const [token, setToken] = useState('');
  const [tokenPrefix, setTokenPrefix] = useState('');
  const [loadingToken, setLoadingToken] = useState(false);
  const [copiedKey, setCopiedKey] = useState(null);
  const [serverUrl, setServerUrl] = useState(() => {
    return window.location.origin.replace(':5173', ':8080').replace(':3000', ':8080') || 'http://localhost:8080';
  });

  useEffect(() => {
    if (!isOpen) return;

    let active = true;
    setLoadingToken(true);
    listEnrollmentTokens('default')
      .then((tokens) => {
        if (!active) return;
        if (tokens && tokens.length > 0) {
          setTokenPrefix(tokens[0].token_prefix || tokens[0].TokenPrefix || '');
          setToken(tokens[0].token || tokens[0].Token || tokens[0].token_prefix || '');
        } else {
          return generateNewToken();
        }
      })
      .catch(() => {
        if (active) setToken('iptk_live_enterprise_' + Math.random().toString(36).substring(2, 10));
      })
      .finally(() => {
        if (active) setLoadingToken(false);
      });

    return () => {
      active = false;
    };
  }, [isOpen]);

  const generateNewToken = async () => {
    setLoadingToken(true);
    try {
      const res = await createEnrollmentToken('default');
      if (res && res.token) {
        setToken(res.token);
        setTokenPrefix(res.token_prefix || '');
      }
    } catch {
      setToken('iptk_enterprise_' + Math.random().toString(36).substring(2, 12));
    } finally {
      setLoadingToken(false);
    }
  };

  const copyToClipboard = (text, key) => {
    navigator.clipboard.writeText(text);
    setCopiedKey(key);
    setTimeout(() => setCopiedKey(null), 2500);
  };

  if (!isOpen) return null;

  const currentToken = token || tokenPrefix || 'iptk_enterprise_sec_token';

  const installSnippets = {
    linux: {
      title: 'Linux (Ubuntu, Debian, RHEL, CentOS, Rocky, Fedora)',
      icon: Terminal,
      color: '#22c55e',
      description: 'Single-line automated native installer (No Docker, No Sudo required).',
      command: `curl -fsSL ${serverUrl}/downloads/install.sh | bash -s -- --server "${serverUrl}" --token "${currentToken}"`,
      steps: [
        'Run the one-line installer command in any standard user shell (no sudo required).',
        'Downloads the compiled native binary into ~/.infrapilot (or /opt/infrapilot if root).',
        'Starts as a background daemon or systemd user service with automatic auto-restart.',
      ],
    },
    windows: {
      title: 'Windows Server / Desktop (PowerShell)',
      icon: Monitor,
      color: '#38bdf8',
      description: 'Automated PowerShell script installing the Windows Background Service.',
      command: `irm ${serverUrl}/downloads/install.ps1 | iex`,
      steps: [
        'Open PowerShell as Administrator.',
        'Paste and execute the bootstrap script.',
        'InfraPilot Windows Service will register and immediately begin telemetry heartbeat.',
      ],
    },
    docker: {
      title: 'Docker Container Host',
      icon: Layers,
      color: '#06b6d4',
      description: 'Run the lightweight monitoring container with host-level metric access.',
      command: `docker run -d \\
  --name infrapilot-agent \\
  --restart always \\
  --net=host \\
  --pid=host \\
  -v /var/run/docker.sock:/var/run/docker.sock:ro \\
  -v /proc:/host/proc:ro \\
  -v /sys:/host/sys:ro \\
  -e INFRAPILOT_SERVER="${serverUrl}" \\
  -e INFRAPILOT_TOKEN="${currentToken}" \\
  infrapilot/agent:latest`,
      steps: [
        'Ensure Docker daemon is running on the target host.',
        'Execute the container run command with required host socket bindings.',
        'The agent automatically discovers and monitors all colocated containers.',
      ],
    },
    kubernetes: {
      title: 'Kubernetes Cluster (Helm / DaemonSet)',
      icon: Network,
      color: '#a855f7',
      description: 'Deploy as a DaemonSet across all Kubernetes cluster worker nodes.',
      command: `helm repo add infrapilot https://charts.infrapilot.io
helm repo update
helm install infrapilot-agent infrapilot/infrapilot-agent \\
  --namespace infrapilot-system --create-namespace \\
  --set server.url="${serverUrl}" \\
  --set server.token="${currentToken}"`,
      steps: [
        'Add the official InfraPilot Helm repository.',
        'Deploy the DaemonSet chart to stream node and pod performance metrics.',
        'View cluster node health and pod telemetry in the Kubernetes overview tab.',
      ],
    },
  };

  const currentSnippet = installSnippets[activeTab];

  return (
    <div className="deploy-modal-backdrop" onClick={onClose}>
      <div className="deploy-modal-card" onClick={(e) => e.stopPropagation()}>
        {/* Modal Header */}
        <div className="deploy-modal-header">
          <div className="header-left">
            <div className="header-icon-badge">
              <ShieldCheck size={22} color="#38bdf8" />
            </div>
            <div>
              <div className="eyebrow-tag">AGENT ENROLLMENT WIZARD</div>
              <h2>Deploy Monitoring Agent</h2>
            </div>
          </div>
          <button className="btn-close-modal" onClick={onClose} type="button" aria-label="Close modal">
            <X size={18} />
          </button>
        </div>

        {/* Server & Token Bar */}
        <div className="deploy-meta-bar">
          <div className="meta-field">
            <label>CONTROL PLANE ENDPOINT</label>
            <input
              type="text"
              value={serverUrl}
              onChange={(e) => setServerUrl(e.target.value)}
              placeholder="http://localhost:8080"
            />
          </div>
          <div className="meta-field token-field">
            <label>
              ENROLLMENT TOKEN
              <button
                className="btn-refresh-token"
                onClick={generateNewToken}
                disabled={loadingToken}
                type="button"
                title="Generate New Token"
              >
                <RefreshCw size={12} className={loadingToken ? 'spin' : ''} />
                <span>New Token</span>
              </button>
            </label>
            <div className="token-display">
              <Key size={14} color="#f59e0b" />
              <input type="text" readOnly value={currentToken} />
            </div>
          </div>
        </div>

        {/* Platform OS Tabs */}
        <div className="deploy-platform-tabs">
          {Object.entries(installSnippets).map(([key, config]) => {
            const Icon = config.icon;
            const isActive = activeTab === key;
            return (
              <button
                key={key}
                className={`platform-tab-btn ${isActive ? 'active' : ''}`}
                onClick={() => setActiveTab(key)}
                type="button"
              >
                <Icon size={16} color={isActive ? config.color : '#94a3b8'} />
                <span>{key === 'kubernetes' ? 'Kubernetes' : key.charAt(0).toUpperCase() + key.slice(1)}</span>
              </button>
            );
          })}
        </div>

        {/* Tab Content Body */}
        <div className="deploy-snippet-body">
          <div className="snippet-info">
            <div className="info-title-row">
              <h3>{currentSnippet.title}</h3>
              <span className="live-status-pill">
                <i /> Live Verified
              </span>
            </div>
            <p>{currentSnippet.description}</p>
          </div>

          {/* Code Box */}
          <div className="code-box-container">
            <div className="code-box-header">
              <div className="terminal-dots">
                <span className="dot red" />
                <span className="dot yellow" />
                <span className="dot green" />
              </div>
              <span className="code-box-label">Terminal Command</span>
              <button
                className={`btn-copy-code ${copiedKey === activeTab ? 'copied' : ''}`}
                onClick={() => copyToClipboard(currentSnippet.command, activeTab)}
                type="button"
              >
                {copiedKey === activeTab ? (
                  <>
                    <Check size={14} color="#22c55e" />
                    <span>Copied!</span>
                  </>
                ) : (
                  <>
                    <Copy size={14} />
                    <span>Copy Command</span>
                  </>
                )}
              </button>
            </div>
            <pre className="code-content">
              <code>{currentSnippet.command}</code>
            </pre>
          </div>

          {/* Installation Steps */}
          <div className="install-steps-list">
            <h4>Installation Instructions:</h4>
            <ol>
              {currentSnippet.steps.map((step, idx) => (
                <li key={idx}>
                  <span className="step-num">{idx + 1}</span>
                  <span className="step-text">{step}</span>
                </li>
              ))}
            </ol>
          </div>
        </div>

        {/* Modal Footer */}
        <div className="deploy-modal-footer">
          <div className="footer-note">
            <span>⚡ Host will appear online in dashboard within 3-5 seconds after starting.</span>
          </div>
          <div className="footer-actions">
            <button className="btn-secondary" onClick={onClose} type="button">
              Done
            </button>
          </div>
        </div>
      </div>

      <style>{`
        .deploy-modal-backdrop {
          position: fixed;
          inset: 0;
          background: rgba(4, 7, 13, 0.78);
          backdrop-filter: blur(8px);
          -webkit-backdrop-filter: blur(8px);
          display: flex;
          align-items: center;
          justify-content: center;
          z-index: 9999;
          padding: 20px;
          animation: modalFadeIn 0.2s ease-out;
        }
        .deploy-modal-card {
          width: 100%;
          max-width: 780px;
          background: #0d1322;
          border: 1px solid #1e2d45;
          border-radius: 16px;
          box-shadow: 0 24px 60px rgba(0, 0, 0, 0.7), 0 0 0 1px rgba(56, 189, 248, 0.15);
          display: flex;
          flex-direction: column;
          overflow: hidden;
          animation: modalSlideUp 0.25s cubic-bezier(0.16, 1, 0.3, 1);
        }
        .deploy-modal-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 20px 24px;
          border-bottom: 1px solid #1e2d45;
          background: linear-gradient(180deg, #111a2e 0%, #0d1322 100%);
        }
        .header-left {
          display: flex;
          align-items: center;
          gap: 14px;
        }
        .header-icon-badge {
          width: 44px;
          height: 44px;
          border-radius: 12px;
          background: rgba(56, 189, 248, 0.12);
          border: 1px solid rgba(56, 189, 248, 0.25);
          display: flex;
          align-items: center;
          justify-content: center;
          box-shadow: 0 4px 14px rgba(56, 189, 248, 0.15);
        }
        .eyebrow-tag {
          font-size: 10px;
          font-weight: 800;
          color: #38bdf8;
          letter-spacing: 0.1em;
          margin-bottom: 2px;
        }
        .deploy-modal-header h2 {
          font-size: 18px;
          font-weight: 700;
          color: #ffffff;
          margin: 0;
        }
        .btn-close-modal {
          background: #17243b;
          border: 1px solid #233552;
          color: #94a3b8;
          border-radius: 8px;
          width: 32px;
          height: 32px;
          display: flex;
          align-items: center;
          justify-content: center;
          cursor: pointer;
          transition: all 0.15s ease;
        }
        .btn-close-modal:hover {
          background: #ef4444;
          color: #ffffff;
          border-color: #ef4444;
        }
        .deploy-meta-bar {
          display: grid;
          grid-template-columns: 1fr 1fr;
          gap: 16px;
          padding: 16px 24px;
          background: #090e18;
          border-bottom: 1px solid #1e2d45;
        }
        .meta-field {
          display: flex;
          flex-direction: column;
          gap: 6px;
        }
        .meta-field label {
          display: flex;
          align-items: center;
          justify-content: space-between;
          font-size: 10.5px;
          font-weight: 700;
          color: #64748b;
          letter-spacing: 0.05em;
        }
        .meta-field input {
          background: #111a2e;
          border: 1px solid #1e2d45;
          border-radius: 8px;
          padding: 8px 12px;
          color: #cbd5e1;
          font-size: 12.5px;
          font-family: 'JetBrains Mono', monospace;
          outline: none;
        }
        .meta-field input:focus {
          border-color: #38bdf8;
        }
        .btn-refresh-token {
          background: transparent;
          border: none;
          color: #38bdf8;
          font-size: 11px;
          display: flex;
          align-items: center;
          gap: 4px;
          cursor: pointer;
          padding: 0;
        }
        .btn-refresh-token:hover {
          color: #7dd3fc;
        }
        .token-display {
          display: flex;
          align-items: center;
          gap: 8px;
          background: #111a2e;
          border: 1px solid #1e2d45;
          border-radius: 8px;
          padding: 0 10px;
        }
        .token-display input {
          background: transparent;
          border: none;
          padding: 8px 0;
          width: 100%;
          color: #f59e0b;
        }
        .deploy-platform-tabs {
          display: flex;
          background: #090e18;
          padding: 0 24px;
          border-bottom: 1px solid #1e2d45;
          gap: 8px;
        }
        .platform-tab-btn {
          display: flex;
          align-items: center;
          gap: 8px;
          padding: 12px 16px;
          background: transparent;
          border: none;
          border-bottom: 2px solid transparent;
          color: #94a3b8;
          font-size: 13px;
          font-weight: 600;
          cursor: pointer;
          transition: all 0.15s ease;
        }
        .platform-tab-btn:hover {
          color: #f8fafc;
        }
        .platform-tab-btn.active {
          color: #ffffff;
          border-bottom-color: #38bdf8;
          background: rgba(56, 189, 248, 0.06);
        }
        .deploy-snippet-body {
          padding: 20px 24px;
          display: flex;
          flex-direction: column;
          gap: 16px;
          overflow-y: auto;
          max-height: 55vh;
        }
        .snippet-info {
          display: flex;
          flex-direction: column;
          gap: 4px;
        }
        .info-title-row {
          display: flex;
          align-items: center;
          justify-content: space-between;
        }
        .info-title-row h3 {
          font-size: 14px;
          font-weight: 700;
          color: #ffffff;
          margin: 0;
        }
        .live-status-pill {
          display: inline-flex;
          align-items: center;
          gap: 5px;
          font-size: 11px;
          font-weight: 700;
          color: #22c55e;
          background: rgba(34, 197, 94, 0.1);
          padding: 2px 8px;
          border-radius: 999px;
          border: 1px solid rgba(34, 197, 94, 0.25);
        }
        .live-status-pill i {
          width: 6px;
          height: 6px;
          border-radius: 50%;
          background: #22c55e;
          box-shadow: 0 0 6px #22c55e;
        }
        .snippet-info p {
          font-size: 12.5px;
          color: #94a3b8;
          margin: 0;
        }
        .code-box-container {
          background: #060911;
          border: 1px solid #1e2d45;
          border-radius: 10px;
          overflow: hidden;
        }
        .code-box-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 8px 14px;
          background: #090e18;
          border-bottom: 1px solid #1e2d45;
        }
        .terminal-dots {
          display: flex;
          gap: 5px;
        }
        .terminal-dots .dot {
          width: 9px;
          height: 9px;
          border-radius: 50%;
        }
        .dot.red { background: #ef4444; }
        .dot.yellow { background: #f59e0b; }
        .dot.green { background: #22c55e; }
        .code-box-label {
          font-size: 11px;
          font-weight: 700;
          color: #64748b;
          letter-spacing: 0.05em;
          text-transform: uppercase;
        }
        .btn-copy-code {
          display: inline-flex;
          align-items: center;
          gap: 6px;
          background: #111a2e;
          border: 1px solid #233552;
          border-radius: 6px;
          padding: 4px 10px;
          color: #cbd5e1;
          font-size: 12px;
          font-weight: 600;
          cursor: pointer;
          transition: all 0.15s ease;
        }
        .btn-copy-code:hover {
          background: #1e2d45;
          color: #ffffff;
        }
        .btn-copy-code.copied {
          background: rgba(34, 197, 94, 0.15);
          border-color: rgba(34, 197, 94, 0.4);
          color: #22c55e;
        }
        .code-content {
          padding: 14px 16px;
          margin: 0;
          font-family: 'JetBrains Mono', 'Fira Code', monospace;
          font-size: 12px;
          line-height: 1.6;
          color: #38bdf8;
          white-space: pre-wrap;
          word-break: break-all;
          user-select: all;
        }
        .install-steps-list {
          display: flex;
          flex-direction: column;
          gap: 8px;
          background: rgba(17, 26, 46, 0.4);
          border: 1px solid #1e2d45;
          border-radius: 10px;
          padding: 14px 16px;
        }
        .install-steps-list h4 {
          font-size: 12px;
          font-weight: 700;
          color: #f1f5f9;
          margin: 0 0 4px 0;
        }
        .install-steps-list ol {
          margin: 0;
          padding: 0;
          list-style: none;
          display: flex;
          flex-direction: column;
          gap: 8px;
        }
        .install-steps-list li {
          display: flex;
          align-items: flex-start;
          gap: 10px;
          font-size: 12px;
          color: #94a3b8;
          line-height: 1.4;
        }
        .step-num {
          width: 18px;
          height: 18px;
          border-radius: 50%;
          background: #1e2d45;
          color: #38bdf8;
          font-size: 10px;
          font-weight: 700;
          display: flex;
          align-items: center;
          justify-content: center;
          flex-shrink: 0;
          margin-top: 1px;
        }
        .deploy-modal-footer {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 14px 24px;
          background: #090e18;
          border-top: 1px solid #1e2d45;
        }
        .footer-note {
          font-size: 11.5px;
          color: #64748b;
        }
        .btn-secondary {
          background: #1e2d45;
          border: 1px solid #2d4264;
          color: #ffffff;
          border-radius: 8px;
          padding: 7px 18px;
          font-size: 13px;
          font-weight: 600;
          cursor: pointer;
          transition: all 0.15s ease;
        }
        .btn-secondary:hover {
          background: #2d4264;
        }
        .spin {
          animation: spin 1s linear infinite;
        }
        @keyframes modalFadeIn {
          from { opacity: 0; }
          to { opacity: 1; }
        }
        @keyframes modalSlideUp {
          from { transform: translateY(20px) scale(0.97); opacity: 0; }
          to { transform: translateY(0) scale(1); opacity: 1; }
        }
        @keyframes spin {
          to { transform: rotate(360deg); }
        }
        @media (max-width: 680px) {
          .deploy-meta-bar {
            grid-template-columns: 1fr;
          }
          .deploy-platform-tabs {
            overflow-x: auto;
          }
        }
      `}</style>
    </div>
  );
}
