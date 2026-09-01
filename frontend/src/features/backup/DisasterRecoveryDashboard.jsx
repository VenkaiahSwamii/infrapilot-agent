// Sprint 9.6: Disaster Recovery Dashboard
import React, { useState, useEffect } from 'react';
import { backupAPI } from '../../api/backup';

export default function DisasterRecoveryDashboard() {
  const [stats, setStats] = useState(null);
  const [backups, setBackups] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      const [statsRes, backupsRes] = await Promise.all([
        backupAPI.getStats(),
        backupAPI.getBackups(),
      ]);
      setStats(statsRes);
      setBackups(backupsRes.backups || []);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleTriggerBackup = async (type) => {
    try {
      await backupAPI.triggerBackup(type);
      alert('Backup triggered successfully');
      loadData();
    } catch (err) {
      alert('Backup failed: ' + err.message);
    }
  };

  const handleRestore = async (backupId) => {
    const component = prompt('Enter component to restore (postgresql, qdrant, kubernetes, config, grafana, full):');
    if (!component) return;
    if (!confirm('WARNING: This will overwrite existing data. Type "RESTORE" to confirm.')) return;
    
    try {
      await backupAPI.restoreBackup(backupId, component);
      alert('Restore completed successfully');
      loadData();
    } catch (err) {
      alert('Restore failed: ' + err.message);
    }
  };

  const handleDelete = async (id) => {
    if (!confirm('Are you sure you want to delete this backup?')) return;
    try {
      await backupAPI.deleteBackup(id);
      loadData();
    } catch (err) {
      alert('Delete failed: ' + err.message);
    }
  };

  if (loading) return <div className="p-6">Loading Disaster Recovery Dashboard...</div>;
  if (error) return <div className="p-6 text-red-600">Error: {error}</div>;

  const formatBytes = (bytes) => {
    if (!bytes) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i];
  };

  const formatDate = (dateStr) => {
    if (!dateStr) return 'Never';
    return new Date(dateStr).toLocaleString();
  };

  return (
    <div className="p-6 space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold text-gray-900">Disaster Recovery Dashboard</h1>
        <div className="space-x-2">
          <button
            onClick={() => handleTriggerBackup('full')}
            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            Trigger Full Backup
          </button>
          <button
            onClick={() => handleTriggerBackup('postgresql')}
            className="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700"
          >
            Backup PostgreSQL
          </button>
        </div>
      </div>

      {/* Stats Grid */}
      {stats && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <div className="bg-white p-6 rounded-lg shadow">
            <h3 className="text-sm font-medium text-gray-500">Last Backup</h3>
            <p className="text-2xl font-semibold text-gray-900">
              {formatDate(stats.last_backup?.created_at)}
            </p>
          </div>
          <div className="bg-white p-6 rounded-lg shadow">
            <h3 className="text-sm font-medium text-gray-500">Total Backup Size</h3>
            <p className="text-2xl font-semibold text-gray-900">
              {formatBytes(stats.backup_size)}
            </p>
          </div>
          <div className="bg-white p-6 rounded-lg shadow">
            <h3 className="text-sm font-medium text-gray-500">Backup Status</h3>
            <p className="text-2xl font-semibold text-green-600">{stats.backup_status}</p>
          </div>
          <div className="bg-white p-6 rounded-lg shadow">
            <h3 className="text-sm font-medium text-gray-500">RTO / RPO</h3>
            <p className="text-2xl font-semibold text-gray-900">
              {stats.rto_seconds}s / {stats.rpo_seconds}s
            </p>
          </div>
        </div>
      )}

      {/* Backup List */}
      <div className="bg-white rounded-lg shadow">
        <div className="px-6 py-4 border-b border-gray-200">
          <h2 className="text-xl font-semibold text-gray-900">Backups</h2>
        </div>
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Type</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Size</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Created</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Encrypted</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {backups.map((backup) => (
                <tr key={backup.id}>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                    {backup.type}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm">
                    <span className={`px-2 py-1 rounded-full text-xs font-medium ${
                      backup.status === 'completed' ? 'bg-green-100 text-green-800' :
                      backup.status === 'verified' ? 'bg-blue-100 text-blue-800' :
                      backup.status === 'failed' ? 'bg-red-100 text-red-800' :
                      'bg-yellow-100 text-yellow-800'
                    }`}>
                      {backup.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {formatBytes(backup.size)}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {formatDate(backup.created_at)}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {backup.encrypted ? 'Yes' : 'No'}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm space-x-2">
                    <button
                      onClick={() => handleRestore(backup.id)}
                      className="text-blue-600 hover:text-blue-800"
                    >
                      Restore
                    </button>
                    <button
                      onClick={() => backupAPI.downloadBackup(backup.id).then(blob => {
                        const url = window.URL.createObjectURL(blob);
                        const a = document.createElement('a');
                        a.href = url;
                        a.download = backup.filename || 'backup';
                        a.click();
                      })}
                      className="text-green-600 hover:text-green-800"
                    >
                      Download
                    </button>
                    <button
                      onClick={() => handleDelete(backup.id)}
                      className="text-red-600 hover:text-red-800"
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}