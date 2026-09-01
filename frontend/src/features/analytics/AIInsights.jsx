import React, { useCallback, useEffect, useState } from 'react';
import { Activity, AlertTriangle, Lightbulb, RefreshCw, ShieldCheck } from 'lucide-react';
import { apiClient } from '../../api/client.js';

const INSIGHT_CATEGORIES = [
  { key: 'performance', label: 'Performance', tone: 'cyan' },
  { key: 'capacity', label: 'Capacity', tone: 'amber' },
  { key: 'security', label: 'Security', tone: 'red' },
  { key: 'reliability', label: 'Reliability', tone: 'green' },
  { key: 'cost', label: 'Cost Optimization', tone: 'violet' },
];

export default function AIInsightsPage() {
  const [insights, setInsights] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const fetchInsights = useCallback(async () => {
    setLoading(true);
    setError('');
    try {
      const data = await apiClient.get('/analytics/insights').catch(() => ({ data: [] }));
      const list = Array.isArray(data.data) ? data.data : [];
      setInsights(list);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchInsights();
    const timer = setInterval(fetchInsights, 20000);
    return () => clearInterval(timer);
  }, [fetchInsights]);

  const insightsByCategory = React.useMemo(() => {
    const map = {};
    for (const category of INSIGHT_CATEGORIES) {
      map[category.key] = [];
    }
    for (const insight of insights) {
      const key = String(insight.category || insight.type || 'performance').toLowerCase();
      if (!map[key]) map[key] = [];
      map[key].push(insight);
    }
    return map;
  }, [insights]);

  return (
    <div className="insights-page">
      <div className="insights-header">
        <div className="insights-title-section">
          <Lightbulb size={28} />
          <h1>AI Insights</h1>
          <span className="insights-count">{insights.length}</span>
        </div>
        <button className="refresh-btn" onClick={fetchInsights} disabled={loading}>
          <RefreshCw size={16} className={loading ? 'spin' : ''} />
          Refresh
        </button>
      </div>

      {error && <div className="insights-error">{error}</div>}

      {loading ? (
        <div className="insights-loading">
          <div className="spinner" />
          <p>Analyzing telemetry...</p>
        </div>
      ) : (
        <div className="insights-layout">
          {INSIGHT_CATEGORIES.map((category) => (
            <section key={category.key} className={`insight-group insight-${category.tone}`}>
              <div className="insight-group-header">
                <div className="insight-group-title">
                  <Activity size={18} />
                  <h2>{category.label}</h2>
                </div>
                <span className="insight-badge">{insightsByCategory[category.key]?.length || 0}</span>
              </div>
              <div className="insight-list">
                {(insightsByCategory[category.key] || []).map((insight, index) => (
                  <article key={insight.id || `${category.key}-${index}`} className="insight-card">
                    <div className="insight-icon">
                      {insight.severity === 'critical' || insight.severity === 'high' ? (
                        <AlertTriangle size={20} />
                      ) : (
                        <ShieldCheck size={20} />
                      )}
                    </div>
                    <div className="insight-body">
                      <div className="insight-title">{insight.title || insight.message || 'Insight'}</div>
                      <p className="insight-message">{insight.description || insight.message || ''}</p>
                      <div className="insight-meta">
                        <span>{insight.hostname || insight.machine_name || 'Platform'}</span>
                        <span>{insight.created_at ? new Date(insight.created_at).toLocaleString() : ''}</span>
                      </div>
                    </div>
                  </article>
                ))}
                {!insightsByCategory[category.key]?.length && (
                  <div className="insight-empty">
                    <p>No insights for this category.</p>
                  </div>
                )}
              </div>
            </section>
          ))}
        </div>
      )}

      <style>{`
        .insights-page {
          padding: 24px;
          max-width: 1200px;
          margin: 0 auto;
          font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
        }

        .insights-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 24px;
        }

        .insights-title-section {
          display: flex;
          align-items: center;
          gap: 12px;
          color: #f1f5f9;
        }

        .insights-title-section h1 {
          margin: 0;
          font-size: 24px;
          font-weight: 700;
        }

        .insights-count {
          background: #334155;
          color: #94a3b8;
          padding: 2px 10px;
          border-radius: 12px;
          font-size: 14px;
        }

        .refresh-btn {
          display: flex;
          align-items: center;
          gap: 6px;
          padding: 8px 16px;
          background: #334155;
          color: #e2e8f0;
          border: 1px solid #475569;
          border-radius: 8px;
          cursor: pointer;
          font-size: 14px;
        }

        .refresh-btn:hover { background: #475569; }

        .insights-error {
          background: #7f1d1d;
          color: #fca5a5;
          padding: 12px 16px;
          border-radius: 8px;
          margin-bottom: 16px;
        }

        .insights-loading {
          text-align: center;
          padding: 60px 20px;
          color: #94a3b8;
        }

        .spinner {
          width: 32px;
          height: 32px;
          border: 3px solid #334155;
          border-top-color: #a78bfa;
          border-radius: 50%;
          animation: spin 0.8s linear infinite;
          margin: 0 auto 12px;
        }

        @keyframes spin { to { transform: rotate(360deg); } }
        .spin { animation: spin 0.8s linear infinite; }

        .insights-layout {
          display: grid;
          grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
          gap: 16px;
        }

        .insight-group {
          background: #0f172a;
          border: 1px solid #334155;
          border-radius: 12px;
          padding: 14px;
        }

        .insight-group-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 10px;
        }

        .insight-group-title {
          display: flex;
          align-items: center;
          gap: 8px;
          color: #f1f5f9;
        }

        .insight-group-title h2 {
          margin: 0;
          font-size: 16px;
          font-weight: 700;
        }

        .insight-badge {
          background: #1e293b;
          color: #e2e8f0;
          padding: 2px 10px;
          border-radius: 12px;
          font-size: 13px;
        }

        .insight-list {
          display: flex;
          flex-direction: column;
          gap: 8px;
        }

        .insight-card {
          display: flex;
          gap: 10px;
          background: #1e293b;
          border: 1px solid #334155;
          border-radius: 10px;
          padding: 10px;
        }

        .insight-icon {
          color: #a78bfa;
        }

        .insight-body {
          flex: 1;
          min-width: 0;
        }

        .insight-title {
          color: #f1f5f9;
          font-size: 14px;
          font-weight: 600;
          margin-bottom: 4px;
        }

        .insight-message {
          color: #94a3b8;
          font-size: 13px;
          margin: 0 0 6px;
        }

        .insight-meta {
          display: flex;
          justify-content: space-between;
          gap: 10px;
          font-size: 12px;
          color: #64748b;
        }

        .insight-empty {
          text-align: center;
          color: #64748b;
          padding: 10px;
        }
      `}</style>
    </div>
  );
}