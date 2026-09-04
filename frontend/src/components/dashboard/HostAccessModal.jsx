import React, { useState, useEffect } from 'react';
import {
  X,
  ShieldCheck,
  UserX,
  Search,
  Check,
  AlertCircle,
  Lock,
  RefreshCw,
  Server,
} from 'lucide-react';
import { apiClient } from '../../api/client.js';

const DEFAULT_USER_PERMISSIONS = [
  {
    id: 'p-1',
    user_id: 'u-1',
    username: 'admin',
    email: 'admin@infrapilot.io',
    user_role: 'SuperAdmin',
    permission_level: 'Full Control',
    granted_by_name: 'System Role (SuperAdmin)',
  },
  {
    id: 'p-2',
    user_id: 'u-2',
    username: 'devops_lead',
    email: 'devops@infrapilot.io',
    user_role: 'DevOps',
    permission_level: 'Operator',
    granted_by_name: 'admin',
  },
  {
    id: 'p-3',
    user_id: 'u-3',
    username: 'secops_auditor',
    email: 'secops@infrapilot.io',
    user_role: 'Viewer',
    permission_level: 'Viewer',
    granted_by_name: 'admin',
  },
  {
    id: 'p-4',
    user_id: 'u-4',
    username: 'operator_user',
    email: 'operator@infrapilot.io',
    user_role: 'Operator',
    permission_level: 'None',
    granted_by_name: 'admin',
  },
];

export default function HostAccessModal({ isOpen, onClose, machine }) {
  const [loading, setLoading] = useState(false);
  const [savingUserId, setSavingUserId] = useState(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [permissions, setPermissions] = useState(DEFAULT_USER_PERMISSIONS);
  const [feedback, setFeedback] = useState(null);

  const machineId = machine?.id || machine?.machine_id;

  useEffect(() => {
    if (!isOpen) return;

    let active = true;
    setLoading(true);
    setFeedback(null);

    if (!machineId || String(machineId).startsWith('m-')) {
      setPermissions(DEFAULT_USER_PERMISSIONS);
      setLoading(false);
      return;
    }

    apiClient
      .get(`/machines/${machineId}/access`)
      .then((res) => {
        if (!active) return;
        const fetched = res.data?.permissions || [];
        setPermissions(fetched.length > 0 ? fetched : DEFAULT_USER_PERMISSIONS);
      })
      .catch((err) => {
        if (!active) return;
        console.error('Failed to load host permissions', err);
        setPermissions(DEFAULT_USER_PERMISSIONS);
      })
      .finally(() => {
        if (active) setLoading(false);
      });

    return () => {
      active = false;
    };
  }, [isOpen, machineId]);

  if (!isOpen || !machine) return null;

  const filteredUsers = permissions.filter((p) => {
    const q = searchQuery.toLowerCase().trim();
    if (!q) return true;
    return (
      p.username?.toLowerCase().includes(q) ||
      p.email?.toLowerCase().includes(q) ||
      p.user_role?.toLowerCase().includes(q)
    );
  });

  const handleUpdatePermission = async (targetUserId, newLevel) => {
    setSavingUserId(targetUserId);
    setFeedback(null);

    try {
      await apiClient.post(`/machines/${machineId}/access`, {
        user_id: targetUserId,
        permission_level: newLevel,
      });

      setPermissions((prev) =>
        prev.map((item) =>
          item.user_id === targetUserId ? { ...item, permission_level: newLevel } : item
        )
      );

      setFeedback({
        type: 'success',
        message: 'Host access updated successfully.',
      });
      setTimeout(() => setFeedback(null), 3000);
    } catch (err) {
      console.error('Failed to update host access', err);
      setFeedback({
        type: 'error',
        message: err.response?.data?.error || 'Failed to update user access level.',
      });
    } finally {
      setSavingUserId(null);
    }
  };

  const getLevelPillClass = (level) => {
    switch (level) {
      case 'Full Control':
        return 'full-control';
      case 'Operator':
        return 'operator';
      case 'Viewer':
        return 'viewer';
      default:
        return 'none';
    }
  };

  return (
    <div className="ham-overlay" onClick={onClose}>
      <div className="ham-card" onClick={(e) => e.stopPropagation()}>
        {/* Header */}
        <div className="ham-header">
          <div className="ham-header-left">
            <div className="ham-icon-avatar">
              <ShieldCheck size={22} color="#38bdf8" />
            </div>
            <div>
              <h2 className="ham-title">Host Access Control & Delegation</h2>
              <div className="ham-subtitle">
                <Server size={13} color="#64748b" />
                <span className="ham-host-name">{machine.hostname}</span>
                <span className="ham-dot">•</span>
                <span>{machine.ip_address}</span>
              </div>
            </div>
          </div>
          <button onClick={onClose} className="ham-close-btn" type="button" title="Close modal">
            <X size={18} />
          </button>
        </div>

        {/* Feedback Alert */}
        {feedback && (
          <div className={`ham-feedback ${feedback.type}`}>
            {feedback.type === 'error' ? (
              <AlertCircle size={16} />
            ) : (
              <Check size={16} />
            )}
            <span>{feedback.message}</span>
          </div>
        )}

        {/* Search Bar */}
        <div className="ham-search-row">
          <div className="ham-input-wrap">
            <Search size={15} color="#64748b" className="ham-search-icon" />
            <input
              type="text"
              placeholder="Search user by name, email, or role..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="ham-search-input"
            />
          </div>
          <span className="ham-user-count">
            {filteredUsers.length} User{filteredUsers.length !== 1 ? 's' : ''}
          </span>
        </div>

        {/* User Permission List */}
        <div className="ham-body">
          {loading ? (
            <div className="ham-loading-state">
              <RefreshCw size={28} className="ham-spin blue-txt" />
              <p>Loading user permissions...</p>
            </div>
          ) : filteredUsers.length === 0 ? (
            <div className="ham-empty-state">
              <UserX size={32} color="#475569" />
              <p>No users found matching your search query.</p>
            </div>
          ) : (
            <div className="ham-users-list">
              {filteredUsers.map((user) => {
                const isAdmin =
                  user.user_role === 'Admin' ||
                  user.user_role === 'SuperAdmin' ||
                  user.user_role === 'OrgAdmin';

                return (
                  <div key={user.user_id} className="ham-user-card">
                    {/* Left: User Avatar & Info */}
                    <div className="ham-user-left">
                      <div className="ham-avatar-circle">
                        {user.username?.charAt(0).toUpperCase() || 'U'}
                      </div>
                      <div className="ham-user-info">
                        <div className="ham-user-name-row">
                          <span className="ham-username">{user.username}</span>
                          <span className={`ham-role-badge ${isAdmin ? 'admin' : 'standard'}`}>
                            {user.user_role || 'User'}
                          </span>
                        </div>
                        <span className="ham-user-email">{user.email}</span>
                      </div>
                    </div>

                    {/* Right: Permission Dropdown */}
                    <div className="ham-user-right">
                      {savingUserId === user.user_id ? (
                        <RefreshCw size={16} className="ham-spin blue-txt" />
                      ) : (
                        <select
                          value={user.permission_level || 'None'}
                          onChange={(e) => handleUpdatePermission(user.user_id, e.target.value)}
                          disabled={isAdmin}
                          className={`ham-select-level ${getLevelPillClass(user.permission_level)} ${
                            isAdmin ? 'disabled' : ''
                          }`}
                        >
                          <option value="Full Control" className="opt-green">
                            ⚡ Full Control
                          </option>
                          <option value="Operator" className="opt-sky">
                            🛠️ Operator
                          </option>
                          <option value="Viewer" className="opt-purple">
                            👁️ Viewer
                          </option>
                          <option value="None" className="opt-muted">
                            🚫 No Access
                          </option>
                        </select>
                      )}
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="ham-footer">
          <div className="ham-footer-left">
            <Lock size={14} color="#64748b" />
            <span>Permissions take effect immediately for active sessions.</span>
          </div>
          <button onClick={onClose} className="ham-done-btn" type="button">
            Done
          </button>
        </div>
      </div>

      <style>{`
        .ham-overlay {
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
          animation: hamFade 0.2s ease-out;
        }

        .ham-card {
          width: 100%;
          max-width: 680px;
          max-height: 88vh;
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
        .ham-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 18px 24px;
          background: rgba(15, 23, 42, 0.6);
          border-bottom: 1px solid rgba(30, 41, 59, 0.8);
        }
        .ham-header-left {
          display: flex;
          align-items: center;
          gap: 14px;
        }
        .ham-icon-avatar {
          width: 42px;
          height: 42px;
          border-radius: 12px;
          background: rgba(56, 189, 248, 0.1);
          border: 1px solid rgba(56, 189, 248, 0.25);
          display: flex;
          align-items: center;
          justify-content: center;
        }
        .ham-title {
          font-size: 16.5px;
          font-weight: 700;
          color: #ffffff;
          margin: 0;
        }
        .ham-subtitle {
          display: flex;
          align-items: center;
          gap: 6px;
          font-size: 12px;
          color: #94a3b8;
          margin-top: 3px;
        }
        .ham-host-name {
          font-family: 'JetBrains Mono', monospace;
          color: #38bdf8;
          font-weight: 600;
        }
        .ham-dot {
          color: #475569;
        }
        .ham-close-btn {
          background: transparent;
          border: 1px solid transparent;
          color: #94a3b8;
          padding: 6px;
          border-radius: 8px;
          cursor: pointer;
          transition: all 0.15s ease;
        }
        .ham-close-btn:hover {
          background: #1e293b;
          color: #ffffff;
          border-color: #334155;
        }

        /* Feedback Alert */
        .ham-feedback {
          display: flex;
          align-items: center;
          gap: 10px;
          padding: 10px 24px;
          font-size: 13px;
          font-weight: 500;
          border-bottom: 1px solid transparent;
        }
        .ham-feedback.success {
          background: rgba(52, 211, 153, 0.1);
          border-color: rgba(52, 211, 153, 0.2);
          color: #34d399;
        }
        .ham-feedback.error {
          background: rgba(248, 113, 113, 0.1);
          border-color: rgba(248, 113, 113, 0.2);
          color: #f87171;
        }

        /* Search Row */
        .ham-search-row {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 14px 24px;
          background: rgba(2, 6, 23, 0.4);
          border-bottom: 1px solid #1e293b;
          gap: 16px;
        }
        .ham-input-wrap {
          position: relative;
          flex: 1;
        }
        .ham-search-icon {
          position: absolute;
          left: 12px;
          top: 50%;
          transform: translateY(-50%);
        }
        .ham-search-input {
          width: 100%;
          background: #020617;
          border: 1px solid #1e293b;
          border-radius: 10px;
          padding: 8px 14px 8px 36px;
          color: #f8fafc;
          font-size: 13px;
          outline: none;
          transition: all 0.15s ease;
        }
        .ham-search-input:focus {
          border-color: #38bdf8;
          box-shadow: 0 0 0 2px rgba(56, 189, 248, 0.2);
        }
        .ham-user-count {
          font-size: 12px;
          font-weight: 600;
          color: #64748b;
          white-space: nowrap;
        }

        /* Body */
        .ham-body {
          flex: 1;
          overflow-y: auto;
          padding: 16px 24px;
          min-height: 280px;
        }
        .ham-loading-state,
        .ham-empty-state {
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          padding: 50px 0;
          color: #64748b;
          font-size: 13px;
          gap: 10px;
        }
        .ham-users-list {
          display: flex;
          flex-direction: column;
          gap: 10px;
        }
        .ham-user-card {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 12px 16px;
          background: rgba(15, 23, 42, 0.5);
          border: 1px solid #1e293b;
          border-radius: 12px;
          transition: all 0.15s ease;
        }
        .ham-user-card:hover {
          background: rgba(15, 23, 42, 0.85);
          border-color: #334155;
        }
        .ham-user-left {
          display: flex;
          align-items: center;
          gap: 12px;
        }
        .ham-avatar-circle {
          width: 36px;
          height: 36px;
          border-radius: 50%;
          background: #1e293b;
          border: 1px solid #334155;
          color: #38bdf8;
          font-size: 14px;
          font-weight: 700;
          display: flex;
          align-items: center;
          justify-content: center;
        }
        .ham-user-info {
          display: flex;
          flex-direction: column;
          gap: 2px;
        }
        .ham-user-name-row {
          display: flex;
          align-items: center;
          gap: 8px;
        }
        .ham-username {
          font-size: 13.5px;
          font-weight: 700;
          color: #f1f5f9;
        }
        .ham-role-badge {
          font-size: 10px;
          font-weight: 700;
          padding: 2px 8px;
          border-radius: 9999px;
          border: 1px solid transparent;
        }
        .ham-role-badge.admin {
          background: rgba(168, 85, 247, 0.12);
          color: #c084fc;
          border-color: rgba(168, 85, 247, 0.3);
        }
        .ham-role-badge.standard {
          background: #1e293b;
          color: #94a3b8;
          border-color: #334155;
        }
        .ham-user-email {
          font-size: 11.5px;
          color: #64748b;
        }

        /* Dropdown */
        .ham-select-level {
          background: #020617;
          border: 1px solid #1e293b;
          border-radius: 8px;
          padding: 6px 12px;
          font-size: 12px;
          font-weight: 700;
          color: #f1f5f9;
          outline: none;
          cursor: pointer;
          transition: all 0.15s ease;
        }
        .ham-select-level.full-control {
          color: #34d399;
          border-color: rgba(52, 211, 153, 0.3);
          background: rgba(52, 211, 153, 0.08);
        }
        .ham-select-level.operator {
          color: #38bdf8;
          border-color: rgba(56, 189, 248, 0.3);
          background: rgba(56, 189, 248, 0.08);
        }
        .ham-select-level.viewer {
          color: #c084fc;
          border-color: rgba(168, 85, 247, 0.3);
          background: rgba(168, 85, 247, 0.08);
        }
        .ham-select-level.none {
          color: #64748b;
          border-color: #334155;
        }
        .ham-select-level.disabled {
          opacity: 0.8;
          cursor: not-allowed;
        }

        .opt-green { background: #0b1120; color: #34d399; }
        .opt-sky { background: #0b1120; color: #38bdf8; }
        .opt-purple { background: #0b1120; color: #c084fc; }
        .opt-muted { background: #0b1120; color: #64748b; }

        /* Footer */
        .ham-footer {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 16px 24px;
          background: rgba(15, 23, 42, 0.6);
          border-top: 1px solid #1e293b;
        }
        .ham-footer-left {
          display: flex;
          align-items: center;
          gap: 8px;
          font-size: 12px;
          color: #64748b;
        }
        .ham-done-btn {
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
        .ham-done-btn:hover {
          background: #334155;
          color: #ffffff;
        }

        .ham-spin {
          animation: hamSpin 1s linear infinite;
        }
        .blue-txt {
          color: #38bdf8;
        }

        @keyframes hamFade {
          from { opacity: 0; transform: scale(0.98); }
          to { opacity: 1; transform: scale(1); }
        }
        @keyframes hamSpin {
          from { transform: rotate(0deg); }
          to { transform: rotate(360deg); }
        }
      `}</style>
    </div>
  );
}
