import React, { useState, useEffect, useCallback } from 'react';
import {
  Sliders, BookOpen, History as HistoryIcon, Clock, Calendar, RefreshCw, PlayCircle
} from 'lucide-react';
import { apiGet } from '../../api/client.js';

import Dashboard from './Dashboard.jsx';
import Rules from './Rules.jsx';
import Runbooks from './Runbooks.jsx';
import History from './History.jsx';
import Approvals from './Approvals.jsx';
import Scheduler from './Scheduler.jsx';

export default function AutomationSuite() {
  const [activeTab, setActiveTab] = useState('overview');

  const [rules, setRules] = useState(null);
  const [runbooks, setRunbooks] = useState(null);
  const [history, setHistory] = useState(null);
  const [approvals, setApprovals] = useState(null);

  const [loading, setLoading] = useState(false);

  const loadData = useCallback(async () => {
    setLoading(true);
    try {
      const [rRes, rbRes, hRes, aRes] = await Promise.allSettled([
        apiGet('/automation/rules'),
        apiGet('/automation/runbooks'),
        apiGet('/automation/history'),
        apiGet('/automation/approvals'),
      ]);

      if (rRes.status === 'fulfilled') setRules(rRes.value);
      if (rbRes.status === 'fulfilled') setRunbooks(rbRes.value);
      if (hRes.status === 'fulfilled') setHistory(hRes.value);
      if (aRes.status === 'fulfilled') setApprovals(aRes.value);
    } catch (err) {
      console.error('Failed loading automation suite data', err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const tabs = [
    { id: 'overview', label: 'Automation Dashboard', icon: PlayCircle },
    { id: 'rules', label: 'Auto-Remediation Rules', icon: Sliders },
    { id: 'runbooks', label: 'Operational Runbooks', icon: BookOpen },
    { id: 'history', label: 'Execution History', icon: HistoryIcon },
    { id: 'approvals', label: 'Pending Approvals', icon: Clock },
    { id: 'scheduler', label: 'Scheduled Cron Jobs', icon: Calendar },
  ];

  return (
    <div style={{ minHeight: '100vh', backgroundColor: '#0d1117', color: '#c9d1d9', padding: '24px 32px' }}>
      {/* Top Header */}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '28px' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '14px' }}>
          <Sliders size={32} color="#58a6ff" />
          <div>
            <h1 style={{ fontSize: '24px', fontWeight: 700, color: '#f0f6fc', margin: 0 }}>
              AIOps Automation &amp; Auto-Remediation Control Plane
            </h1>
            <span style={{ fontSize: '13px', color: '#8b949e' }}>
              Rule-based auto-remediation, operational runbooks, SSH/WinRM/K8s/Docker execution, and approval workflows
            </span>
          </div>
        </div>

        <button
          onClick={loadData}
          disabled={loading}
          style={{
            backgroundColor: '#21262d',
            color: '#c9d1d9',
            border: '1px solid #30363d',
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
          <RefreshCw size={16} className={loading ? 'animate-spin' : ''} />
          {loading ? 'Refreshing...' : 'Refresh Status'}
        </button>
      </div>

      {/* Navigation Tabs */}
      <div style={{ display: 'flex', gap: '8px', borderBottom: '1px solid #30363d', marginBottom: '24px', overflowX: 'auto' }}>
        {tabs.map((tab) => {
          const Icon = tab.icon;
          const isActive = activeTab === tab.id;
          return (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              style={{
                backgroundColor: 'transparent',
                border: 'none',
                borderBottom: isActive ? '2px solid #58a6ff' : '2px solid transparent',
                color: isActive ? '#58a6ff' : '#8b949e',
                padding: '10px 16px',
                cursor: 'pointer',
                fontWeight: 600,
                fontSize: '14px',
                display: 'flex',
                alignItems: 'center',
                gap: '8px',
                whiteSpace: 'nowrap',
                transition: 'all 0.2s'
              }}
            >
              <Icon size={16} />
              {tab.label}
            </button>
          );
        })}
      </div>

      {/* Render Active Tab */}
      {activeTab === 'overview' && <Dashboard rules={rules} history={history} approvals={approvals} />}
      {activeTab === 'rules' && <Rules rules={rules} onRefresh={loadData} />}
      {activeTab === 'runbooks' && <Runbooks runbooks={runbooks} onRefresh={loadData} />}
      {activeTab === 'history' && <History history={history} />}
      {activeTab === 'approvals' && <Approvals approvals={approvals} onRefresh={loadData} />}
      {activeTab === 'scheduler' && <Scheduler />}
    </div>
  );
}
