import React, { useState, useEffect } from 'react';
import {
  X,
  Shield,
  ShieldCheck,
  UserCheck,
  UserX,
  Search,
  Check,
  AlertCircle,
  Key,
  Lock,
  Zap,
  Eye,
  Settings,
  RefreshCw,
  Server,
} from 'lucide-react';
import apiClient from '../../api/client.js';

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

  const getLevelBadgeClass = (level) => {
    switch (level) {
      case 'Full Control':
        return 'bg-emerald-500/10 text-emerald-400 border-emerald-500/30';
      case 'Operator':
        return 'bg-sky-500/10 text-sky-400 border-sky-500/30';
      case 'Viewer':
        return 'bg-purple-500/10 text-purple-400 border-purple-500/30';
      default:
        return 'bg-slate-800 text-slate-400 border-slate-700';
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm p-4">
      <div className="relative w-full max-w-2xl rounded-2xl border border-slate-800 bg-[#0B0F19] text-slate-100 shadow-2xl overflow-hidden flex flex-col max-h-[90vh]">
        {/* Header */}
        <div className="flex items-center justify-between border-b border-slate-800 px-6 py-4 bg-slate-900/50">
          <div className="flex items-center space-x-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-cyan-500/10 text-cyan-400 border border-cyan-500/20">
              <ShieldCheck className="h-5 w-5" />
            </div>
            <div>
              <h2 className="text-lg font-bold tracking-tight text-white flex items-center gap-2">
                Host Access Control & Delegation
              </h2>
              <p className="text-xs text-slate-400 flex items-center gap-2 mt-0.5">
                <Server className="h-3.5 w-3.5 text-slate-500" />
                <span className="font-mono text-cyan-400 font-semibold">{machine.hostname}</span>
                <span className="text-slate-600">•</span>
                <span>{machine.ip_address}</span>
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

        {/* Feedback Alert */}
        {feedback && (
          <div
            className={`px-6 py-3 text-sm flex items-center space-x-2 border-b ${
              feedback.type === 'error'
                ? 'bg-red-500/10 border-red-500/20 text-red-400'
                : 'bg-emerald-500/10 border-emerald-500/20 text-emerald-400'
            }`}
          >
            {feedback.type === 'error' ? (
              <AlertCircle className="h-4 w-4 shrink-0" />
            ) : (
              <Check className="h-4 w-4 shrink-0" />
            )}
            <span>{feedback.message}</span>
          </div>
        )}

        {/* Search Bar */}
        <div className="px-6 pt-4 pb-2 border-b border-slate-800/60 bg-slate-900/30 flex items-center gap-3">
          <div className="relative flex-1">
            <Search className="absolute left-3 top-2.5 h-4 w-4 text-slate-500" />
            <input
              type="text"
              placeholder="Search user by name, email, or role..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full rounded-xl border border-slate-800 bg-slate-950/80 pl-9 pr-4 py-2 text-sm text-slate-200 placeholder-slate-500 focus:border-cyan-500/50 focus:outline-none focus:ring-1 focus:ring-cyan-500/50"
            />
          </div>
          <div className="text-xs text-slate-400 font-medium px-2">
            {filteredUsers.length} User{filteredUsers.length !== 1 ? 's' : ''}
          </div>
        </div>

        {/* User Permission Table */}
        <div className="flex-1 overflow-y-auto p-6 space-y-3 min-h-[300px]">
          {loading ? (
            <div className="flex flex-col items-center justify-center py-12 text-slate-400 space-y-3">
              <RefreshCw className="h-7 w-7 animate-spin text-cyan-400" />
              <p className="text-sm">Loading user permissions...</p>
            </div>
          ) : filteredUsers.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-12 text-slate-500">
              <UserX className="h-8 w-8 mb-2 opacity-50" />
              <p className="text-sm">No users found matching your search.</p>
            </div>
          ) : (
            filteredUsers.map((user) => {
              const isAdmin =
                user.user_role === 'Admin' ||
                user.user_role === 'SuperAdmin' ||
                user.user_role === 'OrgAdmin';

              return (
                <div
                  key={user.user_id}
                  className="flex items-center justify-between p-3.5 rounded-xl border border-slate-800/80 bg-slate-900/40 hover:bg-slate-900/80 transition-all"
                >
                  {/* User info */}
                  <div className="flex items-center space-x-3">
                    <div className="h-9 w-9 rounded-full bg-slate-800 border border-slate-700 flex items-center justify-center font-bold text-slate-300 text-sm">
                      {user.username?.charAt(0).toUpperCase() || 'U'}
                    </div>
                    <div>
                      <div className="flex items-center space-x-2">
                        <span className="font-semibold text-sm text-slate-100">
                          {user.username}
                        </span>
                        <span
                          className={`text-[10px] font-semibold px-2 py-0.5 rounded-full border ${
                            isAdmin
                              ? 'bg-purple-500/10 text-purple-300 border-purple-500/30'
                              : 'bg-slate-800 text-slate-400 border-slate-700'
                          }`}
                        >
                          {user.user_role || 'User'}
                        </span>
                      </div>
                      <span className="text-xs text-slate-500">{user.email}</span>
                    </div>
                  </div>

                  {/* Level dropdown */}
                  <div className="flex items-center space-x-3">
                    {savingUserId === user.user_id ? (
                      <RefreshCw className="h-4 w-4 animate-spin text-cyan-400" />
                    ) : (
                      <select
                        value={user.permission_level || 'None'}
                        onChange={(e) => handleUpdatePermission(user.user_id, e.target.value)}
                        disabled={isAdmin}
                        className={`rounded-lg border px-3 py-1.5 text-xs font-semibold focus:outline-none cursor-pointer transition-colors ${getLevelBadgeClass(
                          user.permission_level
                        )} ${isAdmin ? 'opacity-80 cursor-not-allowed' : ''}`}
                      >
                        <option value="Full Control" className="bg-slate-900 text-emerald-400">
                          ⚡ Full Control
                        </option>
                        <option value="Operator" className="bg-slate-900 text-sky-400">
                          🛠️ Operator
                        </option>
                        <option value="Viewer" className="bg-slate-900 text-purple-400">
                          👁️ Viewer
                        </option>
                        <option value="None" className="bg-slate-900 text-slate-400">
                          🚫 No Access
                        </option>
                      </select>
                    )}
                  </div>
                </div>
              );
            })
          )}
        </div>

        {/* Footer */}
        <div className="flex items-center justify-between border-t border-slate-800 px-6 py-4 bg-slate-900/60 text-xs text-slate-400">
          <div className="flex items-center space-x-2">
            <Lock className="h-3.5 w-3.5 text-slate-500" />
            <span>Permissions take effect immediately for active sessions.</span>
          </div>
          <button
            onClick={onClose}
            className="rounded-xl border border-slate-700 bg-slate-800 px-4 py-2 font-semibold text-slate-200 hover:bg-slate-700 hover:text-white transition-colors"
          >
            Done
          </button>
        </div>
      </div>
    </div>
  );
}
