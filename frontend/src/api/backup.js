// Sprint 9.6: Disaster Recovery & Backup Automation API

const API_BASE = import.meta.env.VITE_API_BASE_URL || '/api/v1';

// Get auth headers
const getHeaders = () => {
  const token = localStorage.getItem('token');
  return {
    'Content-Type': 'application/json',
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
  };
};

export const backupAPI = {
  // List all backups
  getBackups: async () => {
    const res = await fetch(`${API_BASE}/backups`, { headers: getHeaders() });
    if (!res.ok) throw new Error('Failed to fetch backups');
    return res.json();
  },

  // Trigger a new backup
  triggerBackup: async (type = 'full', retention = 'daily') => {
    const res = await fetch(`${API_BASE}/backups/trigger`, {
      method: 'POST',
      headers: getHeaders(),
      body: JSON.stringify({ type, retention, compress: true, encrypt: true }),
    });
    if (!res.ok) throw new Error('Failed to trigger backup');
    return res.json();
  },

  // Download backup
  downloadBackup: async (id) => {
    const res = await fetch(`${API_BASE}/backups/${id}/download`, { headers: getHeaders() });
    if (!res.ok) throw new Error('Failed to download backup');
    return res.blob();
  },

  // Delete backup
  deleteBackup: async (id) => {
    const res = await fetch(`${API_BASE}/backups/${id}`, {
      method: 'DELETE',
      headers: getHeaders(),
    });
    if (!res.ok) throw new Error('Failed to delete backup');
    return res.json();
  },

  // Restore backup
  restoreBackup: async (backupId, component) => {
    const res = await fetch(`${API_BASE}/backups/restore`, {
      method: 'POST',
      headers: getHeaders(),
      body: JSON.stringify({ backup_id: backupId, component, confirm_phrase: 'RESTORE' }),
    });
    if (!res.ok) throw new Error('Failed to restore backup');
    return res.json();
  },

  // Get stats
  getStats: async () => {
    const res = await fetch(`${API_BASE}/backups/stats`, { headers: getHeaders() });
    if (!res.ok) throw new Error('Failed to fetch stats');
    return res.json();
  },

  // Run DR test
  runDRTest: async () => {
    const res = await fetch(`${API_BASE}/backups/dr-test`, {
      method: 'POST',
      headers: getHeaders(),
    });
    if (!res.ok) throw new Error('Failed to run DR test');
    return res.json();
  },

  // Get DR test history
  getDRTestHistory: async () => {
    const res = await fetch(`${API_BASE}/backups/dr-test/history`, { headers: getHeaders() });
    if (!res.ok) throw new Error('Failed to fetch DR test history');
    return res.json();
  },

  // Rotate encryption key
  rotateKey: async () => {
    const res = await fetch(`${API_BASE}/backups/keys/rotate`, {
      method: 'POST',
      headers: getHeaders(),
    });
    if (!res.ok) throw new Error('Failed to rotate key');
    return res.json();
  },

  // Get key info
  getKeyInfo: async () => {
    const res = await fetch(`${API_BASE}/backups/keys/info`, { headers: getHeaders() });
    if (!res.ok) throw new Error('Failed to fetch key info');
    return res.json();
  },

  // Enforce retention
  enforceRetention: async () => {
    const res = await fetch(`${API_BASE}/backups/retention/enforce`, {
      method: 'POST',
      headers: getHeaders(),
    });
    if (!res.ok) throw new Error('Failed to enforce retention');
    return res.json();
  },
};