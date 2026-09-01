import React, { useState, useEffect, useCallback } from 'react';
import {
  Cloud, DollarSign, ShieldAlert, Layers, Plus, RefreshCw, Server
} from 'lucide-react';
import { apiGet, apiPost } from '../../api/client.js';

import CloudDashboard from './CloudDashboard.jsx';
import AWS from './AWS.jsx';
import Azure from './Azure.jsx';
import GCP from './GCP.jsx';
import CostAnalysis from './CostAnalysis.jsx';
import Security from './Security.jsx';
import Resources from './Resources.jsx';

export default function CloudSuite() {
  const [activeTab, setActiveTab] = useState('overview');

  const [accounts, setAccounts] = useState(null);
  const [resources, setResources] = useState(null);
  const [costs, setCosts] = useState(null);
  const [security, setSecurity] = useState(null);
  const [alerts, setAlerts] = useState(null);

  const [loading, setLoading] = useState(false);
  const [showConnectModal, setShowConnectModal] = useState(false);
  const [providerType, setProviderType] = useState('aws');
  const [connectForm, setConnectForm] = useState({ account_id: '', account_name: '', access_key: '', secret_key: '', role_arn: '', tenant_id: '', client_id: '', client_secret: '', service_account_json: '' });
  const [msg, setMsg] = useState(null);

  const loadData = useCallback(async () => {
    setLoading(true);
    try {
      const [accRes, resRes, costRes, secRes, alRes] = await Promise.allSettled([
        apiGet('/cloud/accounts'),
        apiGet('/cloud/resources'),
        apiGet('/cloud/costs'),
        apiGet('/cloud/security'),
        apiGet('/cloud/alerts'),
      ]);

      if (accRes.status === 'fulfilled') setAccounts(accRes.value);
      if (resRes.status === 'fulfilled') setResources(resRes.value);
      if (costRes.status === 'fulfilled') setCosts(costRes.value);
      if (secRes.status === 'fulfilled') setSecurity(secRes.value);
      if (alRes.status === 'fulfilled') setAlerts(alRes.value);
    } catch (err) {
      console.error('Failed loading cloud suite data', err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const handleConnectAccount = async (e) => {
    e.preventDefault();
    try {
      let endpoint = '/cloud/aws/connect';
      let payload = { account_id: connectForm.account_id, account_name: connectForm.account_name, access_key: connectForm.access_key, secret_key: connectForm.secret_key, role_arn: connectForm.role_arn };

      if (providerType === 'azure') {
        endpoint = '/cloud/azure/connect';
        payload = { subscription_id: connectForm.account_id, account_name: connectForm.account_name, tenant_id: connectForm.tenant_id, client_id: connectForm.client_id, client_secret: connectForm.client_secret };
      } else if (providerType === 'gcp') {
        endpoint = '/cloud/gcp/connect';
        payload = { project_id: connectForm.account_id, account_name: connectForm.account_name, service_account_json: connectForm.service_account_json };
      }

      await apiPost(endpoint, payload);
      setMsg({ type: 'success', text: `${providerType.toUpperCase()} account connected successfully!` });
      setShowConnectModal(false);
      setConnectForm({ account_id: '', account_name: '', access_key: '', secret_key: '', role_arn: '', tenant_id: '', client_id: '', client_secret: '', service_account_json: '' });
      loadData();
    } catch (err) {
      setMsg({ type: 'error', text: `Connection failed: ${err.message}` });
    }
  };

  const tabs = [
    { id: 'overview', label: 'Multi-Cloud Overview', icon: Cloud },
    { id: 'aws', label: 'Amazon Web Services', icon: Cloud },
    { id: 'azure', label: 'Microsoft Azure', icon: Cloud },
    { id: 'gcp', label: 'Google Cloud (GCP)', icon: Cloud },
    { id: 'costs', label: 'Cost Analysis & FinOps', icon: DollarSign },
    { id: 'security', label: 'Cloud Security (CSPM)', icon: ShieldAlert },
    { id: 'resources', label: 'Hybrid Cloud Inventory', icon: Layers },
  ];

  return (
    <div style={{ minHeight: '100vh', backgroundColor: '#0d1117', color: '#c9d1d9', padding: '24px 32px' }}>
      {/* Top Header */}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '28px' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '14px' }}>
          <Cloud size={32} color="#58a6ff" />
          <div>
            <h1 style={{ fontSize: '24px', fontWeight: 700, color: '#f0f6fc', margin: 0 }}>
              Hybrid Multi-Cloud Observability Suite
            </h1>
            <span style={{ fontSize: '13px', color: '#8b949e' }}>
              AWS, Azure, and GCP Cloud Monitoring alongside On-Premises Physical & Virtual Infrastructure
            </span>
          </div>
        </div>

        <div style={{ display: 'flex', gap: '12px' }}>
          <button
            onClick={() => setShowConnectModal(true)}
            style={{
              backgroundColor: '#238636',
              color: '#ffffff',
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
            <Plus size={16} />
            Connect Cloud Account
          </button>

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
            {loading ? 'Syncing...' : 'Sync Cloud Accounts'}
          </button>
        </div>
      </div>

      {msg && (
        <div style={{
          backgroundColor: msg.type === 'success' ? '#1f6feb22' : '#f7816622',
          border: `1px solid ${msg.type === 'success' ? '#1f6feb' : '#f78166'}`,
          borderRadius: '6px',
          padding: '12px 16px',
          marginBottom: '20px',
          color: msg.type === 'success' ? '#58a6ff' : '#f78166',
          fontSize: '14px'
        }}>
          {msg.text}
        </div>
      )}

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
      {activeTab === 'overview' && <CloudDashboard accounts={accounts} summary={resources} />}
      {activeTab === 'aws' && <AWS resources={resources} />}
      {activeTab === 'azure' && <Azure resources={resources} />}
      {activeTab === 'gcp' && <GCP resources={resources} />}
      {activeTab === 'costs' && <CostAnalysis costs={costs} />}
      {activeTab === 'security' && <Security findings={security} />}
      {activeTab === 'resources' && <Resources summary={resources} />}

      {/* Connect Cloud Account Modal */}
      {showConnectModal && (
        <div style={{
          position: 'fixed', top: 0, left: 0, right: 0, bottom: 0,
          backgroundColor: 'rgba(0,0,0,0.7)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 1000
        }}>
          <div style={{ backgroundColor: '#161b22', border: '1px solid #30363d', borderRadius: '12px', width: '460px', padding: '24px' }}>
            <h3 style={{ margin: '0 0 16px 0', color: '#f0f6fc' }}>Connect Enterprise Cloud Account</h3>
            <form onSubmit={handleConnectAccount} style={{ display: 'flex', flexDirection: 'column', gap: '14px' }}>
              <div>
                <label style={{ fontSize: '12px', color: '#8b949e', display: 'block', marginBottom: '4px' }}>Cloud Provider</label>
                <select
                  value={providerType}
                  onChange={(e) => setProviderType(e.target.value)}
                  style={{ width: '100%', padding: '8px 12px', background: '#0d1117', border: '1px solid #30363d', color: '#fff', borderRadius: '6px' }}
                >
                  <option value="aws">Amazon Web Services (AWS)</option>
                  <option value="azure">Microsoft Azure</option>
                  <option value="gcp">Google Cloud Platform (GCP)</option>
                </select>
              </div>

              <div>
                <label style={{ fontSize: '12px', color: '#8b949e', display: 'block', marginBottom: '4px' }}>Account Name</label>
                <input
                  type="text"
                  required
                  value={connectForm.account_name}
                  onChange={(e) => setConnectForm({ ...connectForm, account_name: e.target.value })}
                  placeholder="Production AWS Account"
                  style={{ width: '100%', padding: '8px 12px', background: '#0d1117', border: '1px solid #30363d', color: '#fff', borderRadius: '6px' }}
                />
              </div>

              <div>
                <label style={{ fontSize: '12px', color: '#8b949e', display: 'block', marginBottom: '4px' }}>
                  {providerType === 'aws' ? 'AWS Account ID' : providerType === 'azure' ? 'Azure Subscription ID' : 'GCP Project ID'}
                </label>
                <input
                  type="text"
                  required
                  value={connectForm.account_id}
                  onChange={(e) => setConnectForm({ ...connectForm, account_id: e.target.value })}
                  placeholder={providerType === 'aws' ? '123456789012' : providerType === 'azure' ? 'sub-az-01' : 'gcp-prod-infra'}
                  style={{ width: '100%', padding: '8px 12px', background: '#0d1117', border: '1px solid #30363d', color: '#fff', borderRadius: '6px' }}
                />
              </div>

              {providerType === 'aws' && (
                <div>
                  <label style={{ fontSize: '12px', color: '#8b949e', display: 'block', marginBottom: '4px' }}>IAM Role ARN (AssumeRole)</label>
                  <input
                    type="text"
                    value={connectForm.role_arn}
                    onChange={(e) => setConnectForm({ ...connectForm, role_arn: e.target.value })}
                    placeholder="arn:aws:iam::123456789012:role/InfraPilotRole"
                    style={{ width: '100%', padding: '8px 12px', background: '#0d1117', border: '1px solid #30363d', color: '#fff', borderRadius: '6px' }}
                  />
                </div>
              )}

              <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '10px', marginTop: '12px' }}>
                <button type="button" onClick={() => setShowConnectModal(false)} style={{ background: '#21262d', color: '#c9d1d9', border: '1px solid #30363d', padding: '8px 14px', borderRadius: '6px', cursor: 'pointer' }}>Cancel</button>
                <button type="submit" style={{ background: '#238636', color: '#fff', border: 'none', padding: '8px 14px', borderRadius: '6px', cursor: 'pointer', fontWeight: 600 }}>Connect Account</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
