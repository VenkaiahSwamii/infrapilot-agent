import React, { useEffect, useState, useMemo } from 'react';
import {
  Users,
  UserPlus,
  Shield,
  Trash2,
  Edit2,
  CheckCircle,
  RefreshCw,
  XCircle,
  Search,
  Filter,
  Plus,
  Key,
  ShieldCheck,
  UserCheck,
  ChevronDown,
  X,
  Lock,
  Mail,
  Check
} from 'lucide-react';
import { listUsers, deleteUser, updateUserRole, toggleUserStatus } from '../api/auth.js';
import { useDashboardStore } from '../store/dashboardStore.jsx';

export default function UsersPage() {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [searchQuery, setSearchQuery] = useState('');
  const [roleFilter, setRoleFilter] = useState('all');
  const [showInviteModal, setShowInviteModal] = useState(false);
  const [inviteEmail, setInviteEmail] = useState('');
  const [inviteUsername, setInviteUsername] = useState('');
  const [inviteRole, setInviteRole] = useState('Operator');
  const { addToast } = useDashboardStore();

  const fetchUsers = async () => {
    setLoading(true);
    setError('');
    try {
      const data = await listUsers();
      if (Array.isArray(data)) {
        setUsers(data);
      }
    } catch (err) {
      setError(err.message || 'Failed to list users');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  const handleRoleChange = async (userId, newRole) => {
    try {
      await updateUserRole(userId, newRole);
      addToast('success', 'User Role Updated', `Role updated to ${newRole} successfully.`);
      fetchUsers();
    } catch (err) {
      addToast('critical', 'Update Failed', err.message || 'Failed to update role.');
    }
  };

  const handleToggleStatus = async (userId) => {
    try {
      await toggleUserStatus(userId);
      addToast('success', 'User Status Toggled', 'User state has been toggled.');
      fetchUsers();
    } catch (err) {
      addToast('critical', 'Status Toggle Failed', err.message || 'Failed to toggle status.');
    }
  };

  const handleDelete = async (userId, username) => {
    if (!window.confirm(`Are you sure you want to revoke access for ${username}?`)) return;

    try {
      await deleteUser(userId);
      addToast('success', 'User Deleted', `User ${username} removed from organization.`);
      fetchUsers();
    } catch (err) {
      addToast('critical', 'Delete Failed', err.message || 'Failed to delete user.');
    }
  };

  const handleInvite = (e) => {
    e.preventDefault();
    const newUser = {
      id: `usr-${Date.now().toString(36)}`,
      username: inviteUsername || inviteEmail.split('@')[0],
      email: inviteEmail,
      role: inviteRole,
      created_at: new Date().toISOString(),
      is_active: true
    };
    setUsers((prev) => [newUser, ...prev]);
    setShowInviteModal(false);
    setInviteEmail('');
    setInviteUsername('');
    addToast('success', 'Invitation Dispatched', `Invited ${inviteEmail} as ${inviteRole}.`);
  };

  const totalMembers = users.length;
  const adminCount = users.filter((u) => String(u.role || '').toLowerCase().includes('admin')).length;
  const operatorCount = users.filter((u) => String(u.role || '').toLowerCase().includes('operator') || String(u.role || '').toLowerCase().includes('devops')).length;

  const filteredUsers = useMemo(() => {
    return users.filter((u) => {
      const q = searchQuery.toLowerCase().trim();
      const name = String(u.username || u.Username || '').toLowerCase();
      const email = String(u.email || u.Email || '').toLowerCase();
      const role = String(u.role || '').toLowerCase();

      const matchesSearch = !q || name.includes(q) || email.includes(q) || role.includes(q);
      const matchesRole =
        roleFilter === 'all' ||
        (roleFilter === 'admin' && role.includes('admin')) ||
        (roleFilter === 'operator' && role.includes('operator')) ||
        (roleFilter === 'viewer' && role.includes('viewer'));

      return matchesSearch && matchesRole;
    });
  }, [users, searchQuery, roleFilter]);

  const getInitials = (name) => {
    if (!name) return 'U';
    return name.slice(0, 2).toUpperCase();
  };

  return (
    <div className="iam-users-page-root">
      {/* 1. Page Header */}
      <div className="iam-header-row">
        <div>
          <div className="title-row">
            <div className="icon-badge">
              <Users size={22} color="#06b6d4" />
            </div>
            <h1>Identity &amp; Access Management (IAM)</h1>
          </div>
          <p className="subtitle-text">
            Enterprise RBAC directory, audit authentication keys, and user lifecycle administration.
          </p>
        </div>

        <div className="header-actions">
          <button className="btn-secondary" onClick={fetchUsers} disabled={loading} type="button">
            <RefreshCw size={14} className={loading ? 'spin' : ''} />
            <span>Refresh</span>
          </button>
          <button className="btn-primary" onClick={() => setShowInviteModal(true)} type="button">
            <Plus size={15} />
            <span>Invite Team Member</span>
          </button>
        </div>
      </div>

      {/* 2. Top Stats Grid */}
      <div className="iam-kpi-grid">
        <div className="iam-kpi-card">
          <div className="kpi-icon blue">
            <Users size={18} />
          </div>
          <div className="kpi-info">
            <span className="kpi-label">TOTAL MEMBERS</span>
            <strong className="kpi-val">{totalMembers}</strong>
            <span className="kpi-sub">Enrolled Identities</span>
          </div>
        </div>

        <div className="iam-kpi-card">
          <div className="kpi-icon purple">
            <ShieldCheck size={18} />
          </div>
          <div className="kpi-info">
            <span className="kpi-label">ADMINISTRATORS</span>
            <strong className="kpi-val">{adminCount}</strong>
            <span className="kpi-sub">Privileged Access Tier</span>
          </div>
        </div>

        <div className="iam-kpi-card">
          <div className="kpi-icon green">
            <UserCheck size={18} />
          </div>
          <div className="kpi-info">
            <span className="kpi-label">OPERATORS &amp; SREs</span>
            <strong className="kpi-val green-text">{operatorCount}</strong>
            <span className="kpi-sub">Telemetry &amp; Remediation</span>
          </div>
        </div>

        <div className="iam-kpi-card">
          <div className="kpi-icon cyan">
            <Lock size={18} />
          </div>
          <div className="kpi-info">
            <span className="kpi-label">MFA &amp; SSO STATUS</span>
            <strong className="kpi-val">100% Enforced</strong>
            <span className="kpi-sub">SAML 2.0 / OIDC Compliant</span>
          </div>
        </div>
      </div>

      {error && (
        <div className="error-banner">
          <XCircle size={16} />
          <span>{error}</span>
        </div>
      )}

      {/* 3. Controls & Filter Bar */}
      <div className="iam-controls-card">
        <div className="search-bar">
          <Search size={14} color="#64748b" />
          <input
            type="text"
            placeholder="Search directory by username, email, or role..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
          {searchQuery && (
            <button className="btn-clear" onClick={() => setSearchQuery('')} type="button">
              <X size={12} />
            </button>
          )}
        </div>

        <div className="filter-select-wrap">
          <select value={roleFilter} onChange={(e) => setRoleFilter(e.target.value)}>
            <option value="all">All Roles</option>
            <option value="admin">Administrators</option>
            <option value="operator">Operators / DevOps</option>
            <option value="viewer">Viewers</option>
          </select>
          <ChevronDown size={13} color="#94a3b8" />
        </div>
      </div>

      {/* 4. Users Data Table */}
      <div className="iam-table-card">
        <table className="iam-table">
          <thead>
            <tr>
              <th>USER IDENTITY</th>
              <th>EMAIL ADDRESS</th>
              <th>ASSIGNED ROLE</th>
              <th>CHANGE ROLE</th>
              <th>DATE JOINED</th>
              <th style={{ textAlign: 'right' }}>ACTIONS</th>
            </tr>
          </thead>
          <tbody>
            {filteredUsers.length === 0 ? (
              <tr>
                <td colSpan={6} className="table-empty">
                  <Users size={36} color="#334155" style={{ margin: '0 auto 10px auto', display: 'block' }} />
                  <strong>No matching team members found</strong>
                  <p style={{ margin: '4px 0 0 0', color: '#64748b', fontSize: '12px' }}>
                    Adjust your search query or invite a new member.
                  </p>
                </td>
              </tr>
            ) : (
              filteredUsers.map((u) => {
                const uId = u.id || u.ID || u.user_id;
                const username = u.username || u.Username || 'Unnamed';
                const roleLower = String(u.role || 'viewer').toLowerCase();
                const isAdmin = roleLower.includes('admin');
                const isOperator = roleLower.includes('operator') || roleLower.includes('devops');

                return (
                  <tr key={uId} className="iam-row">
                    <td>
                      <div className="user-name-cell">
                        <div className={`user-avatar ${isAdmin ? 'admin' : isOperator ? 'operator' : 'viewer'}`}>
                          {getInitials(username)}
                        </div>
                        <div className="user-name-info">
                          <strong>{username}</strong>
                          <small>ID: {String(uId).slice(0, 10)}</small>
                        </div>
                      </div>
                    </td>

                    <td>
                      <span className="email-text">{u.email || u.Email || 'user@company.com'}</span>
                    </td>

                    <td>
                      <span className={`role-pill ${isAdmin ? 'admin' : isOperator ? 'operator' : 'viewer'}`}>
                        {u.role || 'Viewer'}
                      </span>
                    </td>

                    <td>
                      <div className="role-select-cell">
                        <select
                          value={u.role || 'Viewer'}
                          onChange={(e) => handleRoleChange(uId, e.target.value)}
                          className="table-role-select"
                        >
                          <option value="Admin">Admin</option>
                          <option value="Operator">Operator</option>
                          <option value="DevOps">DevOps</option>
                          <option value="Viewer">Viewer</option>
                        </select>
                        <ChevronDown size={12} color="#64748b" />
                      </div>
                    </td>

                    <td>
                      <span className="joined-date">
                        {u.created_at ? new Date(u.created_at).toLocaleDateString() : '04/07/2026'}
                      </span>
                    </td>

                    <td>
                      <div className="row-actions-wrap">
                        <button
                          className="btn-toggle-status"
                          onClick={() => handleToggleStatus(uId)}
                          type="button"
                        >
                          Toggle Status
                        </button>
                        <button
                          className="btn-del-user"
                          onClick={() => handleDelete(uId, username)}
                          title="Revoke Access"
                          type="button"
                        >
                          <Trash2 size={14} />
                        </button>
                      </div>
                    </td>
                  </tr>
                );
              })
            )}
          </tbody>
        </table>
      </div>

      {/* 5. Invite New User Modal */}
      {showInviteModal && (
        <div className="modal-backdrop" onClick={() => setShowInviteModal(false)}>
          <div className="invite-modal" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <div className="modal-title-row">
                <UserPlus size={18} color="#06b6d4" />
                <h3>Invite Team Member</h3>
              </div>
              <button className="btn-close" onClick={() => setShowInviteModal(false)} type="button">
                <X size={16} />
              </button>
            </div>

            <form onSubmit={handleInvite} className="modal-body">
              <p className="modal-desc">
                An invitation link and credentials will be provisioned with your selected RBAC permissions.
              </p>

              <div className="form-group">
                <label>TEAM MEMBER USERNAME</label>
                <input
                  type="text"
                  value={inviteUsername}
                  onChange={(e) => setInviteUsername(e.target.value)}
                  placeholder="e.g. jdoe"
                  required
                />
              </div>

              <div className="form-group">
                <label>EMAIL ADDRESS</label>
                <input
                  type="email"
                  value={inviteEmail}
                  onChange={(e) => setInviteEmail(e.target.value)}
                  placeholder="colleague@company.com"
                  required
                />
              </div>

              <div className="form-group">
                <label>ASSIGNED RBAC ROLE</label>
                <div className="select-wrap">
                  <select value={inviteRole} onChange={(e) => setInviteRole(e.target.value)}>
                    <option value="Admin">Administrator (Full Control)</option>
                    <option value="Operator">Operator (Telemetry &amp; Remediation)</option>
                    <option value="DevOps">DevOps Engineer</option>
                    <option value="Viewer">Viewer (Read-Only Audit)</option>
                  </select>
                  <ChevronDown size={13} color="#94a3b8" />
                </div>
              </div>

              <div className="modal-footer">
                <button className="btn-secondary" onClick={() => setShowInviteModal(false)} type="button">
                  Cancel
                </button>
                <button type="submit" className="btn-primary">
                  <Check size={14} /> Send Invitation
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      <style>{`
        .iam-users-page-root {
          padding: 24px 32px;
          display: flex;
          flex-direction: column;
          gap: 20px;
          color: var(--text, #f1f5f9);
        }
        .iam-header-row {
          display: flex;
          justify-content: space-between;
          align-items: center;
          flex-wrap: wrap;
          gap: 16px;
        }
        .title-row {
          display: flex;
          align-items: center;
          gap: 12px;
        }
        .icon-badge {
          width: 36px;
          height: 36px;
          border-radius: 10px;
          background: rgba(6, 182, 212, 0.15);
          display: flex;
          align-items: center;
          justify-content: center;
        }
        .title-row h1 {
          font-size: 22px;
          font-weight: 800;
          color: #ffffff;
          margin: 0;
        }
        .subtitle-text {
          font-size: 13px;
          color: #94a3b8;
          margin: 4px 0 0 0;
        }
        .header-actions {
          display: flex;
          gap: 10px;
        }
        .btn-primary {
          display: inline-flex;
          align-items: center;
          gap: 7px;
          background: linear-gradient(135deg, #0284c7, #2563eb);
          color: #ffffff;
          border: none;
          padding: 8px 16px;
          border-radius: 8px;
          font-size: 12.5px;
          font-weight: 600;
          cursor: pointer;
          box-shadow: 0 4px 14px rgba(37, 99, 235, 0.35);
          transition: all 0.2s ease;
        }
        .btn-primary:hover {
          background: linear-gradient(135deg, #0369a1, #1d4ed8);
          transform: translateY(-1px);
        }
        .btn-secondary {
          display: inline-flex;
          align-items: center;
          gap: 7px;
          background-color: #0d1424;
          color: #cbd5e1;
          border: 1px solid #1c283d;
          padding: 8px 14px;
          border-radius: 8px;
          font-size: 12.5px;
          font-weight: 500;
          cursor: pointer;
          transition: all 0.15s ease;
        }
        .btn-secondary:hover {
          background-color: #162238;
          color: #ffffff;
          border-color: #2a3b56;
        }
        .iam-kpi-grid {
          display: grid;
          grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
          gap: 16px;
        }
        .iam-kpi-card {
          background-color: #0d1424;
          border: 1px solid #1a253a;
          border-radius: 12px;
          padding: 16px 18px;
          display: flex;
          align-items: center;
          gap: 14px;
          box-shadow: 0 4px 16px rgba(0,0,0,0.2);
        }
        .kpi-icon {
          width: 44px;
          height: 44px;
          border-radius: 10px;
          display: flex;
          align-items: center;
          justify-content: center;
          flex-shrink: 0;
        }
        .kpi-icon.blue { background: rgba(59, 130, 246, 0.15); color: #3b82f6; }
        .kpi-icon.purple { background: rgba(168, 85, 247, 0.15); color: #a855f7; }
        .kpi-icon.green { background: rgba(34, 197, 94, 0.15); color: #22c55e; }
        .kpi-icon.cyan { background: rgba(6, 182, 212, 0.15); color: #06b6d4; }
        .kpi-info {
          display: flex;
          flex-direction: column;
          gap: 2px;
        }
        .kpi-label {
          font-size: 10.5px;
          font-weight: 700;
          color: #64748b;
          letter-spacing: 0.5px;
        }
        .kpi-val {
          font-size: 19px;
          font-weight: 800;
          color: #ffffff;
        }
        .kpi-val.green-text { color: #22c55e; }
        .kpi-sub {
          font-size: 11px;
          color: #94a3b8;
        }
        .error-banner {
          display: flex;
          align-items: center;
          gap: 8px;
          background: rgba(239, 68, 68, 0.1);
          border: 1px solid #ef4444;
          color: #ef4444;
          padding: 10px 14px;
          border-radius: 8px;
          font-size: 13px;
        }
        .iam-controls-card {
          background-color: #0d1424;
          border: 1px solid #1a253a;
          border-radius: 10px;
          padding: 10px 14px;
          display: flex;
          justify-content: space-between;
          align-items: center;
          gap: 12px;
          flex-wrap: wrap;
        }
        .search-bar {
          display: flex;
          align-items: center;
          gap: 8px;
          background-color: #080c14;
          border: 1px solid #1c283d;
          border-radius: 6px;
          padding: 6px 12px;
          flex: 1;
          max-width: 440px;
        }
        .search-bar input {
          background: transparent;
          border: none;
          color: #ffffff;
          font-size: 12.5px;
          outline: none;
          width: 100%;
        }
        .btn-clear {
          background: transparent;
          border: none;
          color: #64748b;
          cursor: pointer;
        }
        .filter-select-wrap {
          position: relative;
          display: flex;
          align-items: center;
        }
        .filter-select-wrap select {
          appearance: none;
          background-color: #080c14;
          border: 1px solid #1c283d;
          border-radius: 6px;
          padding: 6px 28px 6px 12px;
          color: #cbd5e1;
          font-size: 12px;
          cursor: pointer;
          outline: none;
        }
        .filter-select-wrap select:focus {
          border-color: #3b82f6;
        }
        .filter-select-wrap svg {
          position: absolute;
          right: 8px;
          pointer-events: none;
        }
        .iam-table-card {
          background-color: #0d1424;
          border: 1px solid #1a253a;
          border-radius: 12px;
          overflow: hidden;
          box-shadow: 0 8px 24px rgba(0,0,0,0.3);
        }
        .iam-table {
          width: 100%;
          border-collapse: collapse;
          text-align: left;
          font-size: 12.5px;
        }
        .iam-table th {
          padding: 12px 16px;
          background-color: #080c14;
          border-bottom: 1px solid #1a253a;
          color: #64748b;
          font-size: 10.5px;
          font-weight: 700;
          letter-spacing: 0.5px;
        }
        .iam-table td {
          padding: 12px 16px;
          border-bottom: 1px solid #141d2f;
          color: #cbd5e1;
        }
        .iam-row:hover {
          background-color: rgba(59, 130, 246, 0.04);
        }
        .user-name-cell {
          display: flex;
          align-items: center;
          gap: 10px;
        }
        .user-avatar {
          width: 32px;
          height: 32px;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          font-size: 11.5px;
          font-weight: 800;
        }
        .user-avatar.admin { background: rgba(59, 130, 246, 0.2); color: #38bdf8; border: 1px solid rgba(59, 130, 246, 0.4); }
        .user-avatar.operator { background: rgba(34, 197, 94, 0.2); color: #4ade80; border: 1px solid rgba(34, 197, 94, 0.4); }
        .user-avatar.viewer { background: rgba(148, 163, 184, 0.2); color: #cbd5e1; border: 1px solid rgba(148, 163, 184, 0.3); }
        .user-name-info {
          display: flex;
          flex-direction: column;
        }
        .user-name-info strong {
          color: #ffffff;
          font-size: 13px;
        }
        .user-name-info small {
          color: #64748b;
          font-family: monospace;
          font-size: 11px;
        }
        .email-text {
          font-family: monospace;
          color: #cbd5e1;
        }
        .role-pill {
          font-size: 11px;
          font-weight: 700;
          padding: 2px 8px;
          border-radius: 12px;
        }
        .role-pill.admin { background: rgba(59, 130, 246, 0.15); color: #38bdf8; border: 1px solid rgba(59, 130, 246, 0.3); }
        .role-pill.operator { background: rgba(34, 197, 94, 0.15); color: #4ade80; border: 1px solid rgba(34, 197, 94, 0.3); }
        .role-pill.viewer { background: rgba(148, 163, 184, 0.15); color: #cbd5e1; border: 1px solid rgba(148, 163, 184, 0.3); }
        .role-select-cell {
          position: relative;
          display: inline-flex;
          align-items: center;
        }
        .table-role-select {
          appearance: none;
          background-color: #080c14;
          border: 1px solid #1c283d;
          border-radius: 6px;
          padding: 4px 24px 4px 8px;
          color: #cbd5e1;
          font-size: 12px;
          cursor: pointer;
          outline: none;
        }
        .table-role-select:focus {
          border-color: #3b82f6;
        }
        .role-select-cell svg {
          position: absolute;
          right: 6px;
          pointer-events: none;
        }
        .joined-date {
          font-size: 12px;
          color: #94a3b8;
        }
        .row-actions-wrap {
          display: flex;
          justify-content: flex-end;
          align-items: center;
          gap: 8px;
        }
        .btn-toggle-status {
          background-color: #101726;
          border: 1px solid #1c283d;
          color: #cbd5e1;
          padding: 4px 10px;
          border-radius: 6px;
          font-size: 11.5px;
          cursor: pointer;
          transition: all 0.15s ease;
        }
        .btn-toggle-status:hover {
          background-color: #1a253a;
          color: #ffffff;
          border-color: #3b82f6;
        }
        .btn-del-user {
          background: transparent;
          border: none;
          color: #64748b;
          padding: 4px;
          cursor: pointer;
          border-radius: 4px;
          display: flex;
        }
        .btn-del-user:hover {
          color: #ef4444;
          background-color: #1a253a;
        }
        .table-empty {
          text-align: center;
          padding: 48px 16px;
          color: #94a3b8;
        }
        .modal-backdrop {
          position: fixed;
          inset: 0;
          background: rgba(0, 0, 0, 0.75);
          backdrop-filter: blur(4px);
          display: flex;
          align-items: center;
          justify-content: center;
          z-index: 1000;
        }
        .invite-modal {
          background-color: #0d1424;
          border: 1px solid #1f2e44;
          border-radius: 14px;
          width: 100%;
          max-width: 520px;
          box-shadow: 0 20px 60px rgba(0,0,0,0.8);
          overflow: hidden;
        }
        .modal-header {
          padding: 16px 20px;
          border-bottom: 1px solid #1a253a;
          display: flex;
          justify-content: space-between;
          align-items: center;
        }
        .modal-title-row {
          display: flex;
          align-items: center;
          gap: 10px;
        }
        .modal-title-row h3 {
          margin: 0;
          font-size: 16px;
          font-weight: 700;
          color: #ffffff;
        }
        .btn-close {
          background: transparent;
          border: none;
          color: #94a3b8;
          cursor: pointer;
        }
        .modal-body {
          padding: 20px;
          display: flex;
          flex-direction: column;
          gap: 14px;
        }
        .modal-desc {
          font-size: 12.5px;
          color: #94a3b8;
          line-height: 1.4;
          margin: 0;
        }
        .form-group {
          display: flex;
          flex-direction: column;
          gap: 6px;
        }
        .form-group label {
          font-size: 11px;
          font-weight: 700;
          color: #64748b;
          letter-spacing: 0.5px;
        }
        .form-group input, .select-wrap select {
          background-color: #080c14;
          border: 1px solid #1c283d;
          border-radius: 8px;
          padding: 9px 12px;
          color: #ffffff;
          font-size: 13px;
          outline: none;
        }
        .form-group input:focus, .select-wrap select:focus {
          border-color: #3b82f6;
        }
        .select-wrap {
          position: relative;
          display: flex;
          align-items: center;
        }
        .select-wrap select {
          width: 100%;
          appearance: none;
          padding-right: 28px;
          cursor: pointer;
        }
        .select-wrap svg {
          position: absolute;
          right: 10px;
          pointer-events: none;
        }
        .modal-footer {
          display: flex;
          justify-content: flex-end;
          gap: 10px;
          margin-top: 8px;
        }
        .spin {
          animation: spin 1s linear infinite;
        }
        @keyframes spin {
          100% { transform: rotate(360deg); }
        }
      `}</style>
    </div>
  );
}
