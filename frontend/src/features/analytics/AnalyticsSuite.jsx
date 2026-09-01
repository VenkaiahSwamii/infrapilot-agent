import React, { useState, useEffect, useCallback } from 'react';
import {
  BarChart3, TrendingUp, Gauge, Activity, Award, FileText, ShieldAlert, Filter, RefreshCw
} from 'lucide-react';
import { apiGet } from '../../api/client.js';

import Overview from './Overview.jsx';
import Trends from './Trends.jsx';
import Capacity from './Capacity.jsx';
import SLA from './SLA.jsx';
import ExecutiveDashboard from './ExecutiveDashboard.jsx';
import Reports from './Reports.jsx';
import IncidentAnalytics from './IncidentAnalytics.jsx';

export default function AnalyticsSuite() {
  const [activeTab, setActiveTab] = useState('overview');
  const [trendRange, setTrendRange] = useState('1d');

  const [overviewData, setOverviewData] = useState(null);
  const [trendsData, setTrendsData] = useState(null);
  const [capacityData, setCapacityData] = useState(null);
  const [slaData, setSLAData] = useState(null);
  const [reportsData, setReportsData] = useState(null);
  const [incidentsData, setIncidentsData] = useState(null);

  const [loading, setLoading] = useState(false);

  const loadData = useCallback(async () => {
    setLoading(true);
    try {
      const [oRes, tRes, cRes, sRes, rRes, iRes] = await Promise.allSettled([
        apiGet('/analytics/overview'),
        apiGet(`/analytics/trends?range=${trendRange}`),
        apiGet('/analytics/capacity'),
        apiGet('/analytics/sla'),
        apiGet('/reports'),
        apiGet('/analytics/incidents'),
      ]);

      if (oRes.status === 'fulfilled') setOverviewData(oRes.value);
      if (tRes.status === 'fulfilled') setTrendsData(tRes.value?.trends || []);
      if (cRes.status === 'fulfilled') setCapacityData(cRes.value);
      if (sRes.status === 'fulfilled') setSLAData(sRes.value);
      if (rRes.status === 'fulfilled') setReportsData(rRes.value);
      if (iRes.status === 'fulfilled') setIncidentsData(iRes.value);
    } catch (err) {
      console.error('Failed loading analytics suite data', err);
    } finally {
      setLoading(false);
    }
  }, [trendRange]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const tabs = [
    { id: 'overview', label: 'Infrastructure Analytics', icon: BarChart3 },
    { id: 'executive', label: 'Executive Dashboard', icon: Award },
    { id: 'trends', label: 'Historical Trends', icon: TrendingUp },
    { id: 'capacity', label: 'Capacity Planning', icon: Gauge },
    { id: 'sla', label: 'SLA Monitoring', icon: Activity },
    { id: 'reports', label: 'Reports & Export', icon: FileText },
    { id: 'incidents', label: 'Incident Analytics', icon: ShieldAlert },
  ];

  return (
    <div style={{ minHeight: '100vh', backgroundColor: '#0d1117', color: '#c9d1d9', padding: '24px 32px' }}>
      {/* Top Header */}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '28px' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '14px' }}>
          <BarChart3 size={32} color="#58a6ff" />
          <div>
            <h1 style={{ fontSize: '24px', fontWeight: 700, color: '#f0f6fc', margin: 0 }}>
              Enterprise Analytics, SLA & Executive Insights
            </h1>
            <span style={{ fontSize: '13px', color: '#8b949e' }}>
              Business-ready KPIs, trend analysis, predictive capacity forecasting, and custom report generation
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
          {loading ? 'Refreshing...' : 'Refresh Insights'}
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

      {/* Render Active Feature Tab */}
      {activeTab === 'overview' && <Overview data={overviewData} />}
      {activeTab === 'executive' && <ExecutiveDashboard overview={overviewData} />}
      {activeTab === 'trends' && (
        <Trends
          trends={trendsData}
          range={trendRange}
          onRangeChange={(newRange) => setTrendRange(newRange)}
        />
      )}
      {activeTab === 'capacity' && <Capacity predictions={capacityData} />}
      {activeTab === 'sla' && <SLA sla={slaData} />}
      {activeTab === 'reports' && <Reports reports={reportsData} onRefresh={loadData} />}
      {activeTab === 'incidents' && <IncidentAnalytics incidents={incidentsData} />}
    </div>
  );
}
