import React, { useState } from 'react';
import { Users as UsersIcon, UserPlus, Shield, CheckCircle, XCircle } from 'lucide-react';

export default function Users({ users, onInvite }) {
  const [showModal, setShowModal] = useState(false);
  const [email, setEmail] = useState('');
  const [role, setRole] = useState('operator');

  const handleSubmit = (e) => {
    e.preventDefault();
    onInvite(email, role);
    setEmail('');
    setShowModal(false);
  };

  const getRoleBadgeColor = (r) => {
    switch (r?.toLowerCase()) {
      case 'owner': return { bg: 'rgba(168, 85, 247, 0.15)', color: '#a855f7' };
      case 'admin': return { bg: 'rgba(88, 166, 255, 0.15)', color: '#58a6ff' };
      case 'devops': return { bg: 'rgba(255, 166, 87, 0.15)', color: '#ffa657' };
      case 'operator': return { bg: 'rgba(63, 185, 80, 0.15)', color: '#3fb950' };
      default: return { bg: 'rgba(139, 148, 158, 0.15)', color: '#8b949e' };
    }
  };

  return (
    <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '12px', padding: '24px' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
        <div>
          <h3 style={{ margin: 0, color: '#f0f6fc', fontSize: '18px', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '8px' }}>
            <UsersIcon size={20} color="#58a6ff" />
            Team Member Management
          </h3>
          <span style={{ fontSize: '13px', color: '#8b949e' }}>Manage access roles and invitations for this organization</span>
        </div>
        <button
          onClick={() => setShowModal(true)}
          style={{
            backgroundColor: '#1f6feb',
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
          <UserPlus size={16} />
          Invite Member
        </button>
      </div>

      {users.length === 0 ? (
        <div style={{ color: '#8b949e', fontStyle: 'italic', padding: '16px 0' }}>No members found in this organization.</div>
      ) : (
        <div style={{ overflowX: 'auto' }}>
          <table style={{ width: '100%', borderCollapse: 'collapse', textAlign: 'left', fontSize: '14px' }}>
            <thead>
              <tr style={{ borderBottom: '1px solid #30363d', color: '#8b949e' }}>
                <th style={{ padding: '12px' }}>USER</th>
                <th style={{ padding: '12px' }}>EMAIL</th>
                <th style={{ padding: '12px' }}>ORGANIZATION ROLE</th>
                <th style={{ padding: '12px' }}>JOINED</th>
              </tr>
            </thead>
            <tbody>
              {users.map((u) => {
                const badge = getRoleBadgeColor(u.role);
                return (
                  <tr key={u.id} style={{ borderBottom: '1px solid #21262d', color: '#c9d1d9' }}>
                    <td style={{ padding: '12px', fontWeight: 600, color: '#f0f6fc' }}>{u.name || u.username || 'User'}</td>
                    <td style={{ padding: '12px' }}>{u.email}</td>
                    <td style={{ padding: '12px' }}>
                      <span style={{
                        padding: '4px 10px',
                        borderRadius: '9999px',
                        fontSize: '12px',
                        fontWeight: 600,
                        backgroundColor: badge.bg,
                        color: badge.color,
                        textTransform: 'capitalize'
                      }}>
                        {u.role}
                      </span>
                    </td>
                    <td style={{ padding: '12px', color: '#8b949e', fontSize: '13px' }}>
                      {new Date(u.created_at).toLocaleDateString()}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}

      {showModal && (
        <div style={{
          position: 'fixed', top: 0, left: 0, right: 0, bottom: 0,
          backgroundColor: 'rgba(0,0,0,0.7)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 1000
        }}>
          <div style={{ backgroundColor: '#161b22', border: '1px solid #30363d', borderRadius: '12px', width: '400px', padding: '24px' }}>
            <h3 style={{ margin: '0 0 16px 0', color: '#f0f6fc' }}>Invite Organization Member</h3>
            <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
              <div>
                <label style={{ fontSize: '13px', color: '#8b949e', display: 'block', marginBottom: '6px' }}>Email Address</label>
                <input
                  type="email"
                  required
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="colleague@company.com"
                  style={{ width: '100%', padding: '8px 12px', background: '#0d1117', border: '1px solid #30363d', color: '#fff', borderRadius: '6px' }}
                />
              </div>
              <div>
                <label style={{ fontSize: '13px', color: '#8b949e', display: 'block', marginBottom: '6px' }}>Role</label>
                <select
                  value={role}
                  onChange={(e) => setRole(e.target.value)}
                  style={{ width: '100%', padding: '8px 12px', background: '#0d1117', border: '1px solid #30363d', color: '#fff', borderRadius: '6px' }}
                >
                  <option value="owner">Organization Owner</option>
                  <option value="admin">Organization Admin</option>
                  <option value="devops">DevOps Engineer</option>
                  <option value="operator">Operator</option>
                  <option value="viewer">Viewer</option>
                </select>
              </div>
              <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '10px', marginTop: '12px' }}>
                <button type="button" onClick={() => setShowModal(false)} style={{ background: '#21262d', color: '#c9d1d9', border: '1px solid #30363d', padding: '8px 14px', borderRadius: '6px', cursor: 'pointer' }}>Cancel</button>
                <button type="submit" style={{ background: '#2ea043', color: '#fff', border: 'none', padding: '8px 14px', borderRadius: '6px', cursor: 'pointer', fontWeight: 600 }}>Send Invitation</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
