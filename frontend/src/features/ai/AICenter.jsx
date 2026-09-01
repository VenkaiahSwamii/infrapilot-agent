import React, { useEffect, useState } from 'react';
import { Activity, AlertTriangle, Lightbulb, MessageCircle, RefreshCw, ShieldCheck, TrendingUp } from 'lucide-react';
import { apiClient } from '../../api/client.js';

const STATUS_TONE = {
  OPEN: 'red',
  RESOLVED: 'green',
};

export default function AICenter() {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [incidents, setIncidents] = useState([]);
  const [recommendations, setRecommendations] = useState([]);
  const [predictions, setPredictions] = useState([]);
  const [chatHistory, setChatHistory] = useState([]);
  const [chatInput, setChatInput] = useState('');
  const [chatSending, setChatSending] = useState(false);

  const fetchData = async () => {
    setLoading(true);
    setError('');
    try {
      const [incRes, recRes, predRes] = await Promise.all([
        apiClient.get('/ai/incidents').catch(() => ({ data: [] })),
        apiClient.get('/ai/recommendations').catch(() => ({ data: [] })),
        apiClient.get('/ai/predictions').catch(() => ({ data: [] })),
      ]);

      setIncidents(Array.isArray(incRes.data) ? incRes.data : []);
      setRecommendations(Array.isArray(recRes.data) ? recRes.data : []);
      setPredictions(Array.isArray(predRes.data?.predictions) ? predRes.data.predictions : []);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
    const timer = setInterval(fetchData, 30000);
    return () => clearInterval(timer);
  }, []);

  const sendChat = async () => {
    const question = chatInput.trim();
    if (!question) return;
    setChatSending(true);
    setChatHistory((prev) => [...prev, { role: 'user', text: question }]);
    try {
      const res = await apiClient.post('/ai/chat', { question });
      const answer = res.data?.answer || 'No response';
      setChatHistory((prev) => [...prev, { role: 'assistant', text: answer }]);
    } catch {
      setChatHistory((prev) => [...prev, { role: 'assistant', text: 'I could not process that request right now.' }]);
    } finally {
      setChatSending(false);
      setChatInput('');
    }
  };

  return (
    <div className="ai-center">
      <div className="ai-center-header">
        <div className="ai-center-title">
          <Activity size={26} />
          <h1>AI Center</h1>
        </div>
        <button className="refresh-btn" onClick={fetchData} disabled={loading}>
          <RefreshCw size={16} className={loading ? 'spin' : ''} />
          Refresh
        </button>
      </div>

      {error && <div className="ai-error">{error}</div>}

      {loading ? (
        <div className="ai-loading">Analyzing infrastructure telemetry...</div>
      ) : (
        <div className="ai-layout">
          <section className="ai-panel">
            <div className="ai-section-header">
              <AlertTriangle size={18} />
              <h2>Active Incidents</h2>
              <span className="badge">{incidents.length}</span>
            </div>
            <div className="ai-list">
              {incidents.map((inc) => (
                <article key={inc.id} className="ai-card incident-card">
                  <div className="ai-card-main">
                    <div className="ai-card-title">{inc.title || inc.incident_type || 'Incident'}</div>
                    <p className="ai-card-desc">{inc.description || inc.root_cause || ''}</p>
                    <div className="ai-chip-row">
                      <span className={`chip ${STATUS_TONE[inc.status] || 'neutral'}`}>{inc.status || 'OPEN'}</span>
                      <span className="chip neutral">Severity: {inc.severity}</span>
                      <span className="chip neutral">Confidence: {Math.round((inc.confidence || 0) * 100)}%</span>
                    </div>
                  </div>
                </article>
              ))}
              {!incidents.length && <div className="ai-empty">No active incidents.</div>}
            </div>
          </section>

          <section className="ai-panel">
            <div className="ai-section-header">
              <Lightbulb size={18} />
              <h2>Recommendations</h2>
              <span className="badge">{recommendations.length}</span>
            </div>
            <div className="ai-list">
              {recommendations.map((item) => (
                <article key={item.id} className="ai-card recommendation-card">
                  <div className="ai-card-main">
                    <div className="ai-card-title">{item.title}</div>
                    <p className="ai-card-desc">{item.reason || item.priority || ''}</p>
                    <div className="ai-chip-row">
                      <span className="chip neutral">{item.hostname || 'Platform'}</span>
                      <span className="chip neutral">Priority: {item.priority}</span>
                    </div>
                  </div>
                </article>
              ))}
              {!recommendations.length && <div className="ai-empty">No recommendations.</div>}
            </div>
          </section>

          <section className="ai-panel">
            <div className="ai-section-header">
              <TrendingUp size={18} />
              <h2>Predictions</h2>
              <span className="badge">{predictions.length}</span>
            </div>
            <div className="ai-list">
              {predictions.map((item) => (
                <article key={`${item.machine_id}-${item.metric}`} className="ai-card prediction-card">
                  <div className="ai-card-main">
                    <div className="ai-card-title">{item.metric} usage ({item.machine_id})</div>
                    <p className="ai-card-desc">{item.message || ''}</p>
                    <div className="ai-chip-row">
                      <span className="chip neutral">Current: {item.current_value}%</span>
                      <span className="chip amber">{item.alert_severity || 'Predicted: ' + item.predicted_value + '%'}</span>
                      <span className="chip neutral">{item.days_remaining} days</span>
                    </div>
                  </div>
                </article>
              ))}
              {!predictions.length && <div className="ai-empty">No predictions.</div>}
            </div>
          </section>

          <section className="ai-panel chat-panel">
            <div className="ai-section-header">
              <MessageCircle size={18} />
              <h2>Infrastructure Chat</h2>
            </div>
            <div className="chat-window">
              {chatHistory.map((m, idx) => (
                <div key={idx} className={`chat-msg ${m.role === 'user' ? 'user' : 'assistant'}`}>
                  <div className="chat-bubble">{m.text}</div>
                </div>
              ))}
              {chatSending && <div className="chat-msg assistant"><div className="chat-bubble">Thinking...</div></div>}
            </div>
            <form
              className="chat-form"
              onSubmit={(e) => {
                e.preventDefault();
                sendChat();
              }}
            >
              <input
                value={chatInput}
                onChange={(e) => setChatInput(e.target.value)}
                placeholder="Ask InfraPilot AI: slow server, outages, Docker restarts..."
              />
              <button type="submit" disabled={!chatInput.trim() || chatSending}>
                Send
              </button>
            </form>
          </section>
        </div>
      )}

      <style>{`
        .ai-center {
          padding: 24px;
          max-width: 1200px;
          margin: 0 auto;
          font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
          color: #e2e8f0;
        }
        .ai-center-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 24px;
        }
        .ai-center-title {
          display: flex;
          align-items: center;
          gap: 12px;
        }
        .ai-center-title h1 {
          margin: 0;
          font-size: 22px;
          font-weight: 700;
        }
        .refresh-btn {
          display: inline-flex;
          align-items: center;
          gap: 8px;
          padding: 8px 14px;
          background: #334155;
          color: #e2e8f0;
          border: 1px solid #475569;
          border-radius: 8px;
          cursor: pointer;
        }
        .refresh-btn:hover { background: #475569; }
        .ai-error {
          background: #7f1d1d;
          color: #fca5a5;
          padding: 12px 16px;
          border-radius: 8px;
          margin-bottom: 16px;
        }
        .ai-loading, .ai-empty {
          text-align: center;
          color: #94a3b8;
          padding: 40px 16px;
        }
        .ai-layout {
          display: grid;
          grid-template-columns: repeat(2, 1fr);
          gap: 18px;
        }
        .ai-panel {
          background: #0f172a;
          border: 1px solid #334155;
          border-radius: 12px;
          padding: 14px;
          grid-column: span 1;
        }
        .chat-panel {
          grid-column: 1 / -1;
        }
        .ai-section-header {
          display: flex;
          align-items: center;
          gap: 8px;
          color: #f1f5f9;
          margin-bottom: 10px;
        }
        .ai-section-header h2 {
          margin: 0;
          font-size: 16px;
          font-weight: 700;
        }
        .badge {
          margin-left: auto;
          background: #1e293b;
          color: #e2e8f0;
          padding: 2px 10px;
          border-radius: 12px;
          font-size: 12px;
        }
        .ai-list {
          display: flex;
          flex-direction: column;
          gap: 10px;
        }
        .ai-card {
          background: #1e293b;
          border: 1px solid #334155;
          border-radius: 10px;
          padding: 12px;
        }
        .ai-card-title {
          color: #f1f5f9;
          font-weight: 600;
          margin-bottom: 4px;
        }
        .ai-card-desc {
          color: #94a3b8;
          margin: 0 0 8px;
          font-size: 13px;
        }
        .ai-chip-row {
          display: flex;
          gap: 8px;
          flex-wrap: wrap;
        }
        .chip {
          font-size: 12px;
          padding: 2px 10px;
          border-radius: 999px;
          background: #0f172a;
          color: #e2e8f0;
          border: 1px solid #334155;
        }
        .chip.red { color: #fca5a5; }
        .chip.green { color: #86efac; }
        .chip.amber { color: #fcd34d; }
        .chip.violet { color: #d8b4fe; }
        .chip.neutral { color: #cbd5e1; }
        .chat-window {
          display: flex;
          flex-direction: column;
          gap: 8px;
          min-height: 180px;
          max-height: 340px;
          overflow: auto;
          padding: 10px;
          background: #0b1220;
          border-radius: 10px;
          border: 1px solid #1f2b3d;
          margin-bottom: 10px;
        }
        .chat-msg { display: flex; }
        .chat-msg.user { justify-content: flex-end; }
        .chat-msg.assistant { justify-content: flex-start; }
        .chat-bubble {
          max-width: 340px;
          padding: 8px 12px;
          border-radius: 12px;
          font-size: 14px;
          line-height: 1.4;
          white-space: pre-wrap;
        }
        .chat-msg.user .chat-bubble {
          background: #2563eb;
          color: #fff;
          border-bottom-right-radius: 4px;
        }
        .chat-msg.assistant .chat-bubble {
          background: #334155;
          color: #e2e8f0;
          border-bottom-left-radius: 4px;
        }
        .chat-form {
          display: flex;
          gap: 8px;
        }
        .chat-form input {
          flex: 1;
          padding: 10px 12px;
          background: #0b1220;
          color: #e2e8f0;
          border: 1px solid #334155;
          border-radius: 10px;
        }
        .chat-form button {
          padding: 10px 14px;
          background: #2563eb;
          color: white;
          border: none;
          border-radius: 10px;
        }
        .spin { animation: spin 0.8s linear infinite; }
        @keyframes spin { to { transform: rotate(360deg); } }
        @media (max-width: 900px) {
          .ai-layout { grid-template-columns: 1fr; }
          .chat-panel { grid-column: 1; }
        }
      `}</style>
    </div>
  );
}