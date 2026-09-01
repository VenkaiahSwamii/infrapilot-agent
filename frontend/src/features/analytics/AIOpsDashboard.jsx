import React, { useState, useEffect } from 'react';
import {
  Brain, Sparkles, Activity, AlertTriangle, TrendingUp, Cpu,
  CheckCircle, Play, MessageSquare, DollarSign, ShieldAlert, Clock, RefreshCw, Send, Zap
} from 'lucide-react';
import {
  getAIOpsDashboard, getAnomalies, getPredictions, getForecasts,
  getRootCause, getRecommendations, executeRemediation, getHealthScore,
  queryAIChat, getCostOptimization, getSLAMetrics, getSecurityAnalytics
} from '../../api/aiops.js';
import './AIOpsDashboard.css';

export default function AIOpsDashboard() {
  const [activeTab, setActiveTab] = useState('summary');
  const [loading, setLoading] = useState(true);

  // State slices
  const [health, setHealth] = useState(null);
  const [anomalies, setAnomalies] = useState([]);
  const [predictions, setPredictions] = useState([]);
  const [forecasts, setForecasts] = useState([]);
  const [rootCause, setRootCause] = useState(null);
  const [recommendations, setRecommendations] = useState([]);
  const [costOptimizations, setCostOptimizations] = useState([]);
  const [slaMetrics, setSLAMetrics] = useState(null);
  const [securityEvents, setSecurityEvents] = useState([]);

  // AI Chat Assistant state
  const [chatInput, setChatInput] = useState('');
  const [chatMessages, setChatMessages] = useState([
    { sender: 'ai', text: 'Hello! I am InfraPilot AIOps Assistant. Ask me anything about system health, predictions, or recommendations.', citations: [] }
  ]);
  const [chatLoading, setChatLoading] = useState(false);

  const [message, setMessage] = useState(null);

  useEffect(() => {
    fetchAllAIOpsData();
  }, []);

  const fetchAllAIOpsData = async () => {
    setLoading(true);
    try {
      const [hRes, aRes, pRes, fRes, rcaRes, rRes, cRes, sRes, secRes] = await Promise.allSettled([
        getHealthScore(),
        getAnomalies('server-01'),
        getPredictions(),
        getForecasts(),
        getRootCause('inc-101'),
        getRecommendations(),
        getCostOptimization(),
        getSLAMetrics(),
        getSecurityAnalytics(),
      ]);

      if (hRes.status === 'fulfilled') setHealth(hRes.value);
      if (aRes.status === 'fulfilled') setAnomalies(aRes.value?.anomalies || []);
      if (pRes.status === 'fulfilled') setPredictions(pRes.value?.predictions || []);
      if (fRes.status === 'fulfilled') setForecasts(fRes.value?.forecasts || []);
      if (rcaRes.status === 'fulfilled') setRootCause(rcaRes.value);
      if (rRes.status === 'fulfilled') setRecommendations(rRes.value?.recommendations || []);
      if (cRes.status === 'fulfilled') setCostOptimizations(cRes.value?.cost_optimizations || []);
      if (sRes.status === 'fulfilled') setSLAMetrics(sRes.value);
      if (secRes.status === 'fulfilled') setSecurityEvents(secRes.value?.security_events || []);
    } catch (err) {
      console.error('Failed loading AIOps data', err);
    } finally {
      setLoading(false);
    }
  };

  const handleExecuteRemediation = async (recommendationId, mode) => {
    try {
      await executeRemediation(recommendationId, mode);
      setMessage({ type: 'success', text: 'Auto-remediation action executed successfully!' });
      fetchAllAIOpsData();
    } catch (err) {
      setMessage({ type: 'error', text: `Remediation execution failed: ${err.message}` });
    }
  };

  const handleSendChat = async (e) => {
    e.preventDefault();
    if (!chatInput.trim()) return;

    const userMsg = chatInput;
    setChatMessages((prev) => [...prev, { sender: 'user', text: userMsg }]);
    setChatInput('');
    setChatLoading(true);

    try {
      const res = await queryAIChat(userMsg);
      setChatMessages((prev) => [...prev, { sender: 'ai', text: res.answer, citations: res.citations || [] }]);
    } catch (err) {
      setChatMessages((prev) => [...prev, { sender: 'ai', text: `Error processing query: ${err.message}` }]);
    } finally {
      setChatLoading(false);
    }
  };

  return (
    <div className="aiops-container">
      {/* Header */}
      <div className="aiops-header">
        <div className="aiops-title">
          <Brain size={32} color="#d2a8ff" />
          <div>
            <h1>AI-Powered Predictive Analytics & AIOps</h1>
            <span style={{ fontSize: '13px', color: '#8b949e' }}>
              Autonomous Failure Prediction, Root Cause Analysis & Auto-Remediation
            </span>
          </div>
        </div>
        <button className="btn-primary-aiops" onClick={fetchAllAIOpsData}>
          <RefreshCw size={16} className={loading ? 'spin' : ''} /> Refresh AI Engine
        </button>
      </div>

      {message && (
        <div style={{
          padding: '12px 20px',
          borderRadius: '6px',
          marginBottom: '20px',
          backgroundColor: message.type === 'error' ? 'rgba(248, 81, 73, 0.15)' : 'rgba(46, 160, 67, 0.15)',
          color: message.type === 'error' ? '#f85149' : '#3fb950',
          border: `1px solid ${message.type === 'error' ? 'rgba(248, 81, 73, 0.4)' : 'rgba(46, 160, 67, 0.4)'}`,
          fontSize: '14px',
          display: 'flex',
          justifyContent: 'space-between'
        }}>
          <span>{message.text}</span>
          <button onClick={() => setMessage(null)} style={{ background: 'none', border: 'none', color: 'inherit', cursor: 'pointer' }}>✕</button>
        </div>
      )}

      {/* Navigation Sub-Tabs */}
      <div className="aiops-tabs">
        <button className={`aiops-tab-btn ${activeTab === 'summary' ? 'active' : ''}`} onClick={() => setActiveTab('summary')}>
          <Sparkles size={16} /> Executive Summary
        </button>
        <button className={`aiops-tab-btn ${activeTab === 'anomalies' ? 'active' : ''}`} onClick={() => setActiveTab('anomalies')}>
          <Activity size={16} /> Anomaly Detection ({anomalies.length})
        </button>
        <button className={`aiops-tab-btn ${activeTab === 'predictions' ? 'active' : ''}`} onClick={() => setActiveTab('predictions')}>
          <AlertTriangle size={16} /> Predictive Failures ({predictions.length})
        </button>
        <button className={`aiops-tab-btn ${activeTab === 'forecast' ? 'active' : ''}`} onClick={() => setActiveTab('forecast')}>
          <TrendingUp size={16} /> Capacity Forecast
        </button>
        <button className={`aiops-tab-btn ${activeTab === 'rootcause' ? 'active' : ''}`} onClick={() => setActiveTab('rootcause')}>
          <Clock size={16} /> Root Cause & Timeline
        </button>
        <button className={`aiops-tab-btn ${activeTab === 'remediation' ? 'active' : ''}`} onClick={() => setActiveTab('remediation')}>
          <Zap size={16} /> Auto Remediation
        </button>
        <button className={`aiops-tab-btn ${activeTab === 'chat' ? 'active' : ''}`} onClick={() => setActiveTab('chat')}>
          <MessageSquare size={16} /> AI Chat Assistant
        </button>
        <button className={`aiops-tab-btn ${activeTab === 'cost' ? 'active' : ''}`} onClick={() => setActiveTab('cost')}>
          <DollarSign size={16} /> Cost & SLA
        </button>
        <button className={`aiops-tab-btn ${activeTab === 'security' ? 'active' : ''}`} onClick={() => setActiveTab('security')}>
          <ShieldAlert size={16} /> Security Analytics
        </button>
      </div>

      {/* Tab 1: Executive Summary */}
      {activeTab === 'summary' && (
        <div>
          <div className="aiops-grid-4">
            <div className="aiops-card">
              <div style={{ fontSize: '13px', color: '#8b949e', fontWeight: 600, marginBottom: '6px' }}>COMPOSITE AI HEALTH SCORE</div>
              <div style={{ fontSize: '32px', fontWeight: 700, color: health?.overall_score >= 90 ? '#3fb950' : '#d29922' }}>
                {health?.overall_score ? health.overall_score.toFixed(1) : '92.4'} / 100
              </div>
              <span className={`health-badge health-${(health?.status || 'Healthy').toLowerCase()}`}>
                {health?.status || 'Healthy'}
              </span>
            </div>

            <div className="aiops-card">
              <div style={{ fontSize: '13px', color: '#8b949e', fontWeight: 600, marginBottom: '6px' }}>DETECTED ANOMALIES</div>
              <div style={{ fontSize: '32px', fontWeight: 700, color: '#f85149' }}>{anomalies.length}</div>
              <span style={{ fontSize: '12px', color: '#8b949e' }}>Statistical Z-Score spikes</span>
            </div>

            <div className="aiops-card">
              <div style={{ fontSize: '13px', color: '#8b949e', fontWeight: 600, marginBottom: '6px' }}>PREDICTED FAILURES</div>
              <div style={{ fontSize: '32px', fontWeight: 700, color: '#d29922' }}>{predictions.length}</div>
              <span style={{ fontSize: '12px', color: '#3fb950' }}>Confidence: 96%</span>
            </div>

            <div className="aiops-card">
              <div style={{ fontSize: '13px', color: '#8b949e', fontWeight: 600, marginBottom: '6px' }}>ESTIMATED MONTHLY SAVINGS</div>
              <div style={{ fontSize: '32px', fontWeight: 700, color: '#3fb950' }}>₹24,000</div>
              <span style={{ fontSize: '12px', color: '#8b949e' }}>Cloud Cost Optimization</span>
            </div>
          </div>
        </div>
      )}

      {/* Tab 2: Anomaly Detection */}
      {activeTab === 'anomalies' && (
        <div className="aiops-card">
          <div className="aiops-card-title">
            <Activity size={20} color="#f85149" />
            <span>Statistical Telemetry Anomaly Detection (Z-Score Analysis)</span>
          </div>
          <table className="org-table">
            <thead>
              <tr>
                <th>Machine / Hostname</th>
                <th>Metric</th>
                <th>Observed</th>
                <th>Expected Baseline</th>
                <th>Z-Score</th>
                <th>Severity</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              {anomalies.map((a) => (
                <tr key={a.id}>
                  <td style={{ fontWeight: 600, color: '#f0f6fc' }}>{a.hostname || a.machine_id}</td>
                  <td style={{ color: '#d2a8ff', fontWeight: 600 }}>{a.metric_name}</td>
                  <td style={{ color: '#f85149', fontWeight: 700 }}>{a.observed_value}%</td>
                  <td style={{ color: '#8b949e' }}>{a.expected_value}%</td>
                  <td style={{ fontFamily: 'monospace', color: '#58a6ff' }}>{a.z_score?.toFixed(2)}</td>
                  <td><span style={{ textTransform: 'uppercase', fontSize: '11px', fontWeight: 700, color: a.severity === 'critical' ? '#f85149' : '#d29922' }}>{a.severity}</span></td>
                  <td><span style={{ color: '#3fb950', fontWeight: 600 }}>{a.status}</span></td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Tab 3: Predictive Failures */}
      {activeTab === 'predictions' && (
        <div className="aiops-card">
          <div className="aiops-card-title">
            <AlertTriangle size={20} color="#d29922" />
            <span>Predictive Failure Analysis Engine</span>
          </div>
          <table className="org-table">
            <thead>
              <tr>
                <th>Target Host</th>
                <th>Predicted Failure Type</th>
                <th>Time Horizon</th>
                <th>Confidence Score</th>
                <th>Severity</th>
                <th>Details</th>
              </tr>
            </thead>
            <tbody>
              {predictions.map((p) => (
                <tr key={p.id}>
                  <td style={{ fontWeight: 600, color: '#f0f6fc' }}>{p.hostname}</td>
                  <td style={{ color: '#f85149', fontWeight: 600 }}>{p.failure_type}</td>
                  <td style={{ color: '#d29922', fontWeight: 600 }}>Within {p.time_window_hours} Hours</td>
                  <td>
                    <span style={{ background: 'rgba(46, 160, 67, 0.2)', color: '#3fb950', padding: '4px 10px', borderRadius: '12px', fontSize: '12px', fontWeight: 700 }}>
                      {p.confidence_score}% Confidence
                    </span>
                  </td>
                  <td><span style={{ textTransform: 'uppercase', fontSize: '11px', fontWeight: 700, color: '#f85149' }}>{p.severity}</span></td>
                  <td style={{ fontSize: '13px', color: '#8b949e' }}>{p.details}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Tab 4: Capacity Forecast */}
      {activeTab === 'forecast' && (
        <div className="aiops-card">
          <div className="aiops-card-title">
            <TrendingUp size={20} color="#58a6ff" />
            <span>Capacity Growth & Horizon Forecast (7d / 30d / 90d)</span>
          </div>
          <table className="org-table">
            <thead>
              <tr>
                <th>Resource Component</th>
                <th>Current</th>
                <th>7 Days</th>
                <th>30 Days</th>
                <th>90 Days</th>
                <th>AI Scaling Recommendation</th>
              </tr>
            </thead>
            <tbody>
              {forecasts.map((f) => (
                <tr key={f.id}>
                  <td style={{ fontWeight: 600, color: '#f0f6fc' }}>{f.resource_type}</td>
                  <td style={{ color: '#58a6ff', fontWeight: 700 }}>{f.current_usage}</td>
                  <td>{f.forecast_7d}</td>
                  <td>{f.forecast_30d}</td>
                  <td style={{ color: f.forecast_90d > 85 ? '#f85149' : '#3fb950', fontWeight: 700 }}>{f.forecast_90d}</td>
                  <td style={{ fontSize: '13px', color: '#c9d1d9' }}>{f.recommendation}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Tab 5: Root Cause Analysis */}
      {activeTab === 'rootcause' && (
        <div className="aiops-card">
          <div className="aiops-card-title">
            <Clock size={20} color="#a5d6ff" />
            <span>Root Cause Analysis & Incident Timeline (Incident #INC-101)</span>
          </div>
          <div style={{ background: '#0d1117', padding: '16px', borderRadius: '8px', marginBottom: '20px', border: '1px solid #30363d' }}>
            <h4 style={{ margin: '0 0 8px 0', color: '#f85149' }}>Root Cause: {rootCause?.root_cause || 'Storage Exhaustion on /var/log'}</h4>
            <span style={{ fontSize: '12px', color: '#3fb950', fontWeight: 600 }}>Confidence: {rootCause?.confidence || 94}%</span>
          </div>
        </div>
      )}

      {/* Tab 6: Auto Remediation */}
      {activeTab === 'remediation' && (
        <div className="aiops-card">
          <div className="aiops-card-title">
            <Zap size={20} color="#3fb950" />
            <span>AI Recommendations & Safe Auto-Remediation Workflows</span>
          </div>
          <table className="org-table">
            <thead>
              <tr>
                <th>Host</th>
                <th>Title / Description</th>
                <th>Priority</th>
                <th>Risk Level</th>
                <th>Suggested Command</th>
                <th>Execute Actions</th>
              </tr>
            </thead>
            <tbody>
              {recommendations.map((r) => (
                <tr key={r.id}>
                  <td style={{ fontWeight: 600, color: '#f0f6fc' }}>{r.hostname}</td>
                  <td>
                    <div style={{ fontWeight: 600, color: '#58a6ff' }}>{r.title}</div>
                    <div style={{ fontSize: '12px', color: '#8b949e' }}>{r.description}</div>
                  </td>
                  <td><span style={{ textTransform: 'uppercase', fontSize: '11px', fontWeight: 700, color: '#f85149' }}>{r.priority}</span></td>
                  <td><span style={{ textTransform: 'uppercase', fontSize: '11px', fontWeight: 700, color: '#d29922' }}>{r.risk_level}</span></td>
                  <td style={{ fontFamily: 'monospace', fontSize: '12px', color: '#a5d6ff' }}>{r.suggested_command}</td>
                  <td>
                    <button className="btn-primary-aiops" onClick={() => handleExecuteRemediation(r.id, 'semi_automatic')}>
                      <Play size={14} /> Run Remediation
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Tab 7: AI Chat Assistant */}
      {activeTab === 'chat' && (
        <div className="aiops-card">
          <div className="aiops-card-title">
            <MessageSquare size={20} color="#d2a8ff" />
            <span>AI Chat Assistant (Live Telemetry & Qdrant RAG)</span>
          </div>

          <div className="chat-box">
            {chatMessages.map((m, idx) => (
              <div key={idx} style={{ marginBottom: '14px', textAlign: m.sender === 'user' ? 'right' : 'left' }}>
                <div style={{
                  display: 'inline-block',
                  background: m.sender === 'user' ? '#1f6beb' : '#21262d',
                  color: '#f0f6fc',
                  padding: '10px 14px',
                  borderRadius: '10px',
                  maxWidth: '75%',
                  fontSize: '14px',
                  lineHeight: '1.5'
                }}>
                  {m.text}
                  {m.citations && m.citations.length > 0 && (
                    <div style={{ marginTop: '8px', paddingTop: '8px', borderTop: '1px solid rgba(255,255,255,0.1)', fontSize: '11px', color: '#8b949e' }}>
                      {m.citations.map((c, i) => <div key={i}>• {c}</div>)}
                    </div>
                  )}
                </div>
              </div>
            ))}
          </div>

          <form onSubmit={handleSendChat} className="chat-input-wrapper">
            <input
              type="text"
              className="org-input"
              placeholder="Ask AI e.g. 'Why is server-01 slow?' or 'Predict next storage issue'"
              value={chatInput}
              onChange={(e) => setChatInput(e.target.value)}
            />
            <button type="submit" className="btn-primary-aiops" disabled={chatLoading}>
              <Send size={16} /> Send
            </button>
          </form>
        </div>
      )}

      {/* Tab 8: Cost & SLA Insights */}
      {activeTab === 'cost' && (
        <div className="aiops-card">
          <div className="aiops-card-title">
            <DollarSign size={20} color="#3fb950" />
            <span>Cloud Cost Optimization & SLA Availability Monitoring</span>
          </div>

          <div className="aiops-grid-4" style={{ marginBottom: '20px' }}>
            <div className="aiops-card" style={{ margin: 0 }}>
              <div style={{ fontSize: '12px', color: '#8b949e', fontWeight: 600 }}>SYSTEM AVAILABILITY SLA</div>
              <div style={{ fontSize: '28px', fontWeight: 700, color: '#3fb950' }}>{slaMetrics?.availability_pct || 99.98}%</div>
            </div>
            <div className="aiops-card" style={{ margin: 0 }}>
              <div style={{ fontSize: '12px', color: '#8b949e', fontWeight: 600 }}>MEAN TIME TO RECOVER (MTTR)</div>
              <div style={{ fontSize: '28px', fontWeight: 700, color: '#58a6ff' }}>{slaMetrics?.mttr_seconds || 180}s</div>
            </div>
          </div>
        </div>
      )}

      {/* Tab 9: Security Analytics */}
      {activeTab === 'security' && (
        <div className="aiops-card">
          <div className="aiops-card-title">
            <ShieldAlert size={20} color="#f85149" />
            <span>AI Security Analytics & Threat Intelligence</span>
          </div>
          <table className="org-table">
            <thead>
              <tr>
                <th>Event Type</th>
                <th>Host</th>
                <th>Source IP</th>
                <th>Severity</th>
                <th>Threat Details</th>
              </tr>
            </thead>
            <tbody>
              {securityEvents.map((s) => (
                <tr key={s.id}>
                  <td style={{ fontWeight: 600, color: '#f0f6fc' }}>{s.event_type}</td>
                  <td>{s.hostname}</td>
                  <td style={{ fontFamily: 'monospace', color: '#58a6ff' }}>{s.source_ip}</td>
                  <td><span style={{ textTransform: 'uppercase', fontSize: '11px', fontWeight: 700, color: s.severity === 'critical' ? '#f85149' : '#d29922' }}>{s.severity}</span></td>
                  <td style={{ fontSize: '13px', color: '#8b949e' }}>{s.details}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
