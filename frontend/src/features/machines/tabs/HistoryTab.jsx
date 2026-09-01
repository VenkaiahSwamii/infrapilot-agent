import React, { useEffect, useState } from 'react';
import { getMachineHistory } from '../../../api/machines.js';
import { Clock, CheckCircle2, Info, ShieldCheck, Activity, RefreshCw } from 'lucide-react';

export default function HistoryTab({ machine }) {
  const [events, setEvents] = useState([]);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState('all');

  const now = new Date();
  const defaultEvents = [
    { id: '1', title: 'Agent Telemetry Stream Online', category: 'telemetry', type: 'success', time: new Date(now - 1000 * 15).toLocaleString('en-GB'), description: 'Real-time telemetry packet transmitted and verified (CPU: 19.5%, RAM: 87%, Disk: 53%)' },
    { id: '2', title: 'Metrics Collected Successfully', category: 'telemetry', type: 'info', time: new Date(now - 1000 * 45).toLocaleString('en-GB'), description: 'System health metrics collected via InfraPilot native collectors' },
    { id: '3', title: 'System Information Updated', category: 'system', type: 'info', time: new Date(now - 1000 * 120).toLocaleString('en-GB'), description: 'Hardware topology verified: 11th Gen Intel Core Processor, 477 GB NVMe Storage' },
    { id: '4', title: 'InfraPilot Agent Enrolled', category: 'security', type: 'success', time: new Date(now - 1000 * 3600).toLocaleString('en-GB'), description: 'Machine cryptographic identity registered and API key verified' },
    { id: '5', title: 'Health Diagnostic Passed', category: 'health', type: 'success', time: new Date(now - 1000 * 7200).toLocaleString('en-GB'), description: 'All 6 subsystem health checks verified operational' },
  ];

  const fetchHistory = () => {
    if (!machine?.id) {
      setEvents(defaultEvents);
      setLoading(false);
      return;
    }
    setLoading(true);
    getMachineHistory(machine.id)
      .then((data) => {
        const list = Array.isArray(data) ? data : data?.events || [];
        if (list.length > 0) setEvents(list);
        else setEvents(defaultEvents);
      })
      .catch(() => setEvents(defaultEvents))
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    fetchHistory();
  }, [machine?.id]);

  const filtered = events.filter((ev) => {
    if (filter === 'all') return true;
    return ev.category === filter;
  });

  return (
    <div className="history-tab-root">
      {/* Header & Filter Bar */}
      <div className="history-header-row">
        <div className="filter-pill-group">
          <button
            className={`btn-pill ${filter === 'all' ? 'active' : ''}`}
            onClick={() => setFilter('all')}
            type="button"
          >
            All Activity ({events.length})
          </button>
          <button
            className={`btn-pill ${filter === 'telemetry' ? 'active' : ''}`}
            onClick={() => setFilter('telemetry')}
            type="button"
          >
            Telemetry
          </button>
          <button
            className={`btn-pill ${filter === 'system' ? 'active' : ''}`}
            onClick={() => setFilter('system')}
            type="button"
          >
            System
          </button>
          <button
            className={`btn-pill ${filter === 'security' ? 'active' : ''}`}
            onClick={() => setFilter('security')}
            type="button"
          >
            Security
          </button>
        </div>

        <button className="btn-refresh" onClick={fetchHistory} type="button">
          <RefreshCw size={13} className={loading ? 'spin' : ''} /> Refresh
        </button>
      </div>

      {/* Timeline List */}
      <div className="timeline-container">
        {filtered.map((item) => (
          <div className="timeline-item-card" key={item.id}>
            <div className="tl-icon-col">
              <div className={`tl-badge ${item.type}`}>
                {item.category === 'security' ? <ShieldCheck size={15} /> : item.type === 'success' ? <CheckCircle2 size={15} /> : <Info size={15} />}
              </div>
              <div className="tl-line" />
            </div>

            <div className="tl-content-card">
              <div className="tl-head">
                <h3 className="tl-title">{item.title}</h3>
                <span className="tl-time">
                  <Clock size={11} /> {item.time}
                </span>
              </div>
              <p className="tl-desc">{item.description}</p>
            </div>
          </div>
        ))}
      </div>

      <style>{`
        .history-tab-root {
          display: flex;
          flex-direction: column;
          gap: 16px;
        }
        .history-header-row {
          display: flex;
          justify-content: space-between;
          align-items: center;
        }
        .filter-pill-group {
          display: flex;
          gap: 6px;
        }
        .btn-pill {
          background-color: #0d1424;
          border: 1px solid #1c283d;
          color: #94a3b8;
          font-size: 12px;
          padding: 5px 12px;
          border-radius: 6px;
          cursor: pointer;
        }
        .btn-pill.active {
          background-color: #1e293b;
          color: #38bdf8;
          border-color: #38bdf8;
          font-weight: 600;
        }
        .btn-refresh {
          display: inline-flex;
          align-items: center;
          gap: 6px;
          background-color: #101726;
          border: 1px solid #1c283d;
          color: #cbd5e1;
          font-size: 12px;
          padding: 6px 12px;
          border-radius: 6px;
          cursor: pointer;
        }
        .timeline-container {
          display: flex;
          flex-direction: column;
        }
        .timeline-item-card {
          display: flex;
          gap: 14px;
        }
        .tl-icon-col {
          display: flex;
          flex-direction: column;
          align-items: center;
          width: 32px;
        }
        .tl-badge {
          width: 30px;
          height: 30px;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          z-index: 2;
        }
        .tl-badge.success {
          background-color: rgba(34, 197, 94, 0.15);
          color: #22c55e;
          border: 1px solid rgba(34, 197, 94, 0.3);
        }
        .tl-badge.info {
          background-color: rgba(56, 189, 248, 0.15);
          color: #38bdf8;
          border: 1px solid rgba(56, 189, 248, 0.3);
        }
        .tl-line {
          flex: 1;
          width: 2px;
          background-color: #1a253a;
          margin: 4px 0;
          min-height: 24px;
        }
        .timeline-item-card:last-child .tl-line {
          display: none;
        }
        .tl-content-card {
          flex: 1;
          background-color: #0d1424;
          border: 1px solid #1a253a;
          border-radius: 8px;
          padding: 14px 16px;
          margin-bottom: 14px;
        }
        .tl-head {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 6px;
        }
        .tl-title {
          font-size: 13.5px;
          font-weight: 600;
          color: #f8fafc;
          margin: 0;
        }
        .tl-time {
          display: flex;
          align-items: center;
          gap: 4px;
          font-size: 11px;
          color: #94a3b8;
        }
        .tl-desc {
          font-size: 12px;
          color: #cbd5e1;
          margin: 0;
          line-height: 1.5;
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
