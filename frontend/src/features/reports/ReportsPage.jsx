import { useState, useEffect } from 'react';
import { reportsApi } from '../../api/reports';

export function ReportsPage() {
  const [reports, setReports] = useState([]);
  const [loading, setLoading] = useState(false);
  const [generating, setGenerating] = useState(false);

  const [form, setForm] = useState({
    name: '',
    type: 'infrastructure_summary',
    format: 'pdf',
    date_from: '',
    date_to: '',
  });

  useEffect(() => {
    fetchReports();
  }, []);

  const fetchReports = async () => {
    setLoading(true);
    try {
      const data = await reportsApi.list();
      setReports(data);
    } catch (err) {
      console.error('Failed to fetch reports', err);
    } finally {
      setLoading(false);
    }
  };

  const handleGenerate = async (e) => {
    e.preventDefault();
    setGenerating(true);
    try {
      const payload = {
        name: form.name || `${form.type} report`,
        type: form.type,
        format: form.format,
        date_from: form.date_from || null,
        date_to: form.date_to || null,
      };
      await reportsApi.generate(payload);
      await fetchReports();
      setForm({ name: '', type: 'infrastructure_summary', format: 'pdf', date_from: '', date_to: '' });
    } catch (err) {
      console.error('Failed to generate report', err);
    } finally {
      setGenerating(false);
    }
  };

  const handleDownload = async (id) => {
    try {
      const blob = await reportsApi.download(id);
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `report_${id}`;
      document.body.appendChild(a);
      a.click();
      a.remove();
      window.URL.revokeObjectURL(url);
    } catch (err) {
      console.error('Download failed', err);
    }
  };

  const handleDelete = async (id) => {
    try {
      await reportsApi.delete(id);
      await fetchReports();
    } catch (err) {
      console.error('Delete failed', err);
    }
  };

  return (
    <div className="reports-page">
      <h1>Reports</h1>

      <section className="card generate-report">
        <h2>Generate Report</h2>
        <form onSubmit={handleGenerate} className="form-grid">
          <label>
            Name
            <input
              value={form.name}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
              placeholder="My report"
            />
          </label>

          <label>
            Type
            <select
              value={form.type}
              onChange={(e) => setForm({ ...form, type: e.target.value })}
            >
              <option value="infrastructure_summary">Infrastructure Summary</option>
              <option value="machine_health">Machine Health</option>
              <option value="cpu_usage">CPU Usage</option>
              <option value="memory_usage">Memory Usage</option>
              <option value="disk_usage">Disk Usage</option>
              <option value="docker_status">Docker Status</option>
              <option value="kubernetes_status">Kubernetes Status</option>
              <option value="alerts">Alerts</option>
              <option value="security_events">Security Events</option>
              <option value="inventory">Inventory</option>
            </select>
          </label>

          <label>
            Format
            <select
              value={form.format}
              onChange={(e) => setForm({ ...form, format: e.target.value })}
            >
              <option value="pdf">PDF</option>
              <option value="excel">Excel</option>
              <option value="csv">CSV</option>
            </select>
          </label>

          <label>
            Date From
            <input
              type="date"
              value={form.date_from}
              onChange={(e) => setForm({ ...form, date_from: e.target.value })}
            />
          </label>

          <label>
            Date To
            <input
              type="date"
              value={form.date_to}
              onChange={(e) => setForm({ ...form, date_to: e.target.value })}
            />
          </label>

          <button type="submit" disabled={generating} className="btn primary">
            {generating ? 'Generating...' : 'Generate'}
          </button>
        </form>
      </section>

      <section className="card">
        <h2>Report History</h2>
        {loading && <p>Loading...</p>}
        {!loading && reports.length === 0 && <p>No reports generated yet.</p>}
        {!loading && reports.length > 0 && (
          <table className="table">
            <thead>
              <tr>
                <th>Name</th>
                <th>Type</th>
                <th>Format</th>
                <th>Status</th>
                <th>Created</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {reports.map((r) => (
                <tr key={r.id}>
                  <td>{r.name}</td>
                  <td>{r.type}</td>
                  <td>{r.format.toUpperCase()}</td>
                  <td>{r.status}</td>
                  <td>{new Date(r.created_at).toLocaleString()}</td>
                  <td>
                    {r.status === 'completed' && (
                      <button className="btn small" onClick={() => handleDownload(r.id)}>Download</button>
                    )}
                    <button className="btn small danger" onClick={() => handleDelete(r.id)}>Delete</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </section>
    </div>
  );
}