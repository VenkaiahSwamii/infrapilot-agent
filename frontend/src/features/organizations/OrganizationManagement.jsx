import React, { useState, useEffect } from 'react';
import {
  Building2, Users, Palette, Gauge, CreditCard, Lock, Plus, Save, UserPlus, CheckCircle, ShieldCheck
} from 'lucide-react';
import {
  getOrganizations, createOrganization, getOrganizationSettings,
  updateOrganizationSettings, getOrganizationQuotas, updateOrganizationQuotas,
  getOrganizationBilling, getOrganizationUsers, createInvitation
} from '../../api/organization.js';
import OrganizationSwitcher from '../../components/layout/OrganizationSwitcher.jsx';
import './OrganizationManagement.css';

export default function OrganizationManagement() {
  const [activeTab, setActiveTab] = useState('overview');
  const [organizations, setOrganizations] = useState([]);
  const [currentOrg, setCurrentOrg] = useState(null);
  const [users, setUsers] = useState([]);
  const [settings, setSettings] = useState({
    company_name: '',
    logo_url: '',
    primary_color: '#58a6ff',
    sidebar_color: '#0d1117',
    custom_domain: '',
    license_tier: 'enterprise',
  });
  const [quotas, setQuotas] = useState({
    max_machines: 50,
    max_users: 20,
    max_storage_gb: 500,
    retention_days: 90,
  });
  const [billing, setBilling] = useState({
    plan: 'enterprise',
    status: 'active',
    payment_provider: 'stripe',
  });

  const [message, setMessage] = useState(null);
  const [showCreateOrgModal, setShowCreateOrgModal] = useState(false);
  const [showInviteModal, setShowInviteModal] = useState(false);

  const [newOrgForm, setNewOrgForm] = useState({ name: '', slug: '' });
  const [inviteForm, setInviteForm] = useState({ email: '', role: 'operator' });

  useEffect(() => {
    loadInitialData();
  }, []);

  const loadInitialData = async () => {
    try {
      const orgsRes = await getOrganizations();
      const orgList = orgsRes.organizations || [];
      setOrganizations(orgList);

      const activeId = localStorage.getItem('activeOrgId') || orgList[0]?.id || 'default';
      const activeOrg = orgList.find(o => o.id === activeId) || { id: activeId, name: 'Default Organization' };
      setCurrentOrg(activeOrg);

      if (activeOrg.id && activeOrg.id !== 'default') {
        loadOrgDetails(activeOrg.id);
      }
    } catch (err) {
      console.error('Failed to load multi-tenant data', err);
    }
  };

  const loadOrgDetails = async (orgId) => {
    try {
      const [settRes, quotaRes, billRes, userRes] = await Promise.allSettled([
        getOrganizationSettings(orgId),
        getOrganizationQuotas(orgId),
        getOrganizationBilling(orgId),
        getOrganizationUsers(orgId),
      ]);

      if (settRes.status === 'fulfilled') setSettings(settRes.value);
      if (quotaRes.status === 'fulfilled') setQuotas(quotaRes.value);
      if (billRes.status === 'fulfilled') setBilling(billRes.value);
      if (userRes.status === 'fulfilled') setUsers(userRes.value?.members || []);
    } catch (err) {
      console.error('Failed loading org details', err);
    }
  };

  const handleCreateOrg = async (e) => {
    e.preventDefault();
    try {
      const res = await createOrganization(newOrgForm);
      setMessage({ type: 'success', text: `Organization '${res.name}' created successfully!` });
      setShowCreateOrgModal(false);
      setNewOrgForm({ name: '', slug: '' });
      loadInitialData();
    } catch (err) {
      setMessage({ type: 'error', text: `Failed to create organization: ${err.message}` });
    }
  };

  const handleSaveSettings = async (e) => {
    e.preventDefault();
    if (!currentOrg?.id) return;
    try {
      const res = await updateOrganizationSettings(currentOrg.id, settings);
      setSettings(res);
      setMessage({ type: 'success', text: 'White-label branding & settings saved successfully!' });
    } catch (err) {
      setMessage({ type: 'error', text: `Failed to save settings: ${err.message}` });
    }
  };

  const handleSaveQuotas = async (e) => {
    e.preventDefault();
    if (!currentOrg?.id) return;
    try {
      const res = await updateOrganizationQuotas(currentOrg.id, quotas);
      setQuotas(res);
      setMessage({ type: 'success', text: 'Organization resource quotas updated!' });
    } catch (err) {
      setMessage({ type: 'error', text: `Failed to update quotas: ${err.message}` });
    }
  };

  const handleSendInvite = async (e) => {
    e.preventDefault();
    if (!currentOrg?.id) return;
    try {
      await createInvitation(currentOrg.id, inviteForm);
      setMessage({ type: 'success', text: `Invitation sent to ${inviteForm.email}` });
      setShowInviteModal(false);
      setInviteForm({ email: '', role: 'operator' });
    } catch (err) {
      setMessage({ type: 'error', text: `Failed to send invitation: ${err.message}` });
    }
  };

  return (
    <div className="org-container">
      {/* Header */}
      <div className="org-header">
        <div className="org-header-title">
          <Building2 size={32} color="#58a6ff" />
          <div>
            <h1>Multi-Tenant Enterprise Management</h1>
            <span style={{ fontSize: '13px', color: '#8b949e' }}>
              Centralized SaaS Administration & Tenant Isolation
            </span>
          </div>
        </div>
        <div style={{ display: 'flex', gap: '12px', alignItems: 'center' }}>
          <OrganizationSwitcher onSelectOrg={(org) => {
            setCurrentOrg(org);
            if (org.id) loadOrgDetails(org.id);
          }} />
          <button className="btn-primary-org" onClick={() => setShowCreateOrgModal(true)}>
            <Plus size={16} /> New Tenant
          </button>
        </div>
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

      {/* Navigation Tabs */}
      <div className="org-tabs">
        <button className={`org-tab-btn ${activeTab === 'overview' ? 'active' : ''}`} onClick={() => setActiveTab('overview')}>
          <Building2 size={16} /> Organizations
        </button>
        <button className={`org-tab-btn ${activeTab === 'branding' ? 'active' : ''}`} onClick={() => setActiveTab('branding')}>
          <Palette size={16} /> White-Label Branding
        </button>
        <button className={`org-tab-btn ${activeTab === 'users' ? 'active' : ''}`} onClick={() => setActiveTab('users')}>
          <Users size={16} /> Users & RBAC
        </button>
        <button className={`org-tab-btn ${activeTab === 'quotas' ? 'active' : ''}`} onClick={() => setActiveTab('quotas')}>
          <Gauge size={16} /> Quotas & Limits
        </button>
        <button className={`org-tab-btn ${activeTab === 'billing' ? 'active' : ''}`} onClick={() => setActiveTab('billing')}>
          <CreditCard size={16} /> Billing Readiness
        </button>
        <button className={`org-tab-btn ${activeTab === 'sso' ? 'active' : ''}`} onClick={() => setActiveTab('sso')}>
          <Lock size={16} /> Enterprise SSO
        </button>
      </div>

      {/* Tab 1: Overview & Organizations List */}
      {activeTab === 'overview' && (
        <div className="org-card">
          <div className="org-card-title">
            <Building2 size={20} color="#58a6ff" />
            <span>Active Tenant Organizations ({organizations.length})</span>
          </div>
          <table className="org-table">
            <thead>
              <tr>
                <th>Organization Name</th>
                <th>Slug</th>
                <th>Status</th>
                <th>License Tier</th>
                <th>Created At</th>
              </tr>
            </thead>
            <tbody>
              {organizations.map((o) => (
                <tr key={o.id}>
                  <td style={{ fontWeight: 600, color: '#f0f6fc' }}>{o.name}</td>
                  <td style={{ fontFamily: 'monospace' }}>{o.slug}</td>
                  <td>
                    <span style={{ color: o.status === 'active' ? '#3fb950' : '#8b949e', fontWeight: 600 }}>
                      ● {o.status || 'active'}
                    </span>
                  </td>
                  <td><span style={{ textTransform: 'uppercase', fontSize: '12px', fontWeight: 600, color: '#58a6ff' }}>Enterprise</span></td>
                  <td>{new Date(o.created_at || Date.now()).toLocaleDateString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Tab 2: White-Label Branding */}
      {activeTab === 'branding' && (
        <div className="org-card">
          <div className="org-card-title">
            <Palette size={20} color="#58a6ff" />
            <span>Tenant White-Label Customization - {currentOrg?.name}</span>
          </div>
          <form onSubmit={handleSaveSettings}>
            <div className="org-grid-2">
              <div className="org-input-group">
                <label>Company / Display Name</label>
                <input
                  type="text"
                  className="org-input"
                  value={settings.company_name || ''}
                  onChange={(e) => setSettings({ ...settings, company_name: e.target.value })}
                />
              </div>

              <div className="org-input-group">
                <label>Custom Domain</label>
                <input
                  type="text"
                  className="org-input"
                  placeholder="infra.yourcompany.com"
                  value={settings.custom_domain || ''}
                  onChange={(e) => setSettings({ ...settings, custom_domain: e.target.value })}
                />
              </div>

              <div className="org-input-group">
                <label>Primary Brand Color</label>
                <div style={{ display: 'flex', alignItems: 'center' }}>
                  <input
                    type="text"
                    className="org-input"
                    value={settings.primary_color || '#58a6ff'}
                    onChange={(e) => setSettings({ ...settings, primary_color: e.target.value })}
                  />
                  <span className="color-preview-box" style={{ backgroundColor: settings.primary_color }} />
                </div>
              </div>

              <div className="org-input-group">
                <label>Sidebar Accent Color</label>
                <div style={{ display: 'flex', alignItems: 'center' }}>
                  <input
                    type="text"
                    className="org-input"
                    value={settings.sidebar_color || '#0d1117'}
                    onChange={(e) => setSettings({ ...settings, sidebar_color: e.target.value })}
                  />
                  <span className="color-preview-box" style={{ backgroundColor: settings.sidebar_color }} />
                </div>
              </div>
            </div>

            <button type="submit" className="btn-primary-org" style={{ marginTop: '16px' }}>
              <Save size={16} /> Save Branding & Settings
            </button>
          </form>
        </div>
      )}

      {/* Tab 3: Users & Team RBAC */}
      {activeTab === 'users' && (
        <div className="org-card">
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '18px' }}>
            <div className="org-card-title" style={{ margin: 0 }}>
              <Users size={20} color="#58a6ff" />
              <span>Team Members & Roles - {currentOrg?.name}</span>
            </div>
            <button className="btn-primary-org" onClick={() => setShowInviteModal(true)}>
              <UserPlus size={16} /> Invite User
            </button>
          </div>

          <table className="org-table">
            <thead>
              <tr>
                <th>User</th>
                <th>Email</th>
                <th>Role</th>
              </tr>
            </thead>
            <tbody>
              {users.map((u) => (
                <tr key={u.id}>
                  <td style={{ fontWeight: 600, color: '#f0f6fc' }}>{u.name || 'Team Member'}</td>
                  <td>{u.email}</td>
                  <td>
                    <span style={{ textTransform: 'uppercase', fontSize: '12px', fontWeight: 700, color: '#a5d6ff' }}>
                      {u.role}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Tab 4: Quotas & Resource Limits */}
      {activeTab === 'quotas' && (
        <div className="org-card">
          <div className="org-card-title">
            <Gauge size={20} color="#58a6ff" />
            <span>Resource Consumption Quotas - {currentOrg?.name}</span>
          </div>
          <form onSubmit={handleSaveQuotas}>
            <div className="org-grid-2">
              <div className="org-input-group">
                <label>Max Machines (0 = Unlimited)</label>
                <input
                  type="number"
                  className="org-input"
                  value={quotas.max_machines}
                  onChange={(e) => setQuotas({ ...quotas, max_machines: parseInt(e.target.value) || 0 })}
                />
              </div>

              <div className="org-input-group">
                <label>Max User Seats</label>
                <input
                  type="number"
                  className="org-input"
                  value={quotas.max_users}
                  onChange={(e) => setQuotas({ ...quotas, max_users: parseInt(e.target.value) || 0 })}
                />
              </div>

              <div className="org-input-group">
                <label>Max Metric Storage (GB)</label>
                <input
                  type="number"
                  className="org-input"
                  value={quotas.max_storage_gb}
                  onChange={(e) => setQuotas({ ...quotas, max_storage_gb: parseInt(e.target.value) || 0 })}
                />
              </div>

              <div className="org-input-group">
                <label>Data Retention (Days)</label>
                <input
                  type="number"
                  className="org-input"
                  value={quotas.retention_days}
                  onChange={(e) => setQuotas({ ...quotas, retention_days: parseInt(e.target.value) || 0 })}
                />
              </div>
            </div>

            <button type="submit" className="btn-primary-org">
              <Save size={16} /> Save Quota Limits
            </button>
          </form>
        </div>
      )}

      {/* Tab 5: Billing Readiness */}
      {activeTab === 'billing' && (
        <div className="org-card">
          <div className="org-card-title">
            <CreditCard size={20} color="#58a6ff" />
            <span>SaaS Billing Readiness & Payment Gateways</span>
          </div>
          <div className="org-grid-2">
            <div style={{ background: '#0d1117', padding: '16px', borderRadius: '8px', border: '1px solid #30363d' }}>
              <h4 style={{ margin: '0 0 8px 0', color: '#f0f6fc' }}>Active Plan Tier</h4>
              <p style={{ fontSize: '20px', fontWeight: 700, color: '#3fb950', textTransform: 'uppercase', margin: 0 }}>
                {billing.plan || 'Enterprise'}
              </p>
              <span style={{ fontSize: '12px', color: '#8b949e' }}>Status: {billing.status}</span>
            </div>

            <div style={{ background: '#0d1117', padding: '16px', borderRadius: '8px', border: '1px solid #30363d' }}>
              <h4 style={{ margin: '0 0 8px 0', color: '#f0f6fc' }}>Supported Payment Providers</h4>
              <p style={{ fontSize: '13px', color: '#c9d1d9', margin: 0 }}>
                Stripe • Razorpay • PayPal • AWS Marketplace • Azure Marketplace
              </p>
            </div>
          </div>
        </div>
      )}

      {/* Tab 6: Enterprise SSO */}
      {activeTab === 'sso' && (
        <div className="org-card">
          <div className="org-card-title">
            <Lock size={20} color="#58a6ff" />
            <span>Enterprise Single Sign-On (SSO) Integration</span>
          </div>
          <div className="org-grid-2">
            <div className="org-input-group">
              <label>SSO Provider</label>
              <select className="org-select">
                <option value="saml">SAML 2.0</option>
                <option value="oidc">OpenID Connect (OIDC)</option>
                <option value="azure">Azure AD / Entra ID</option>
                <option value="google">Google Workspace</option>
                <option value="okta">Okta</option>
                <option value="ldap">Active Directory / LDAP</option>
              </select>
            </div>

            <div className="org-input-group">
              <label>Identity Provider (IdP) Metadata URL</label>
              <input type="text" className="org-input" placeholder="https://idp.yourcompany.com/saml/metadata" />
            </div>
          </div>
          <button className="btn-primary-org">
            <ShieldCheck size={16} /> Enable SSO Provider
          </button>
        </div>
      )}

      {/* Create Org Modal */}
      {showCreateOrgModal && (
        <div style={{ position: 'fixed', top: 0, left: 0, right: 0, bottom: 0, background: 'rgba(0,0,0,0.8)', display: 'flex', justifyContent: 'center', alignItems: 'center', zIndex: 1000 }}>
          <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '10px', padding: '24px', width: '450px' }}>
            <h3 style={{ margin: '0 0 16px 0', color: '#f0f6fc' }}>Create New Tenant Organization</h3>
            <form onSubmit={handleCreateOrg}>
              <div className="org-input-group">
                <label>Organization Name</label>
                <input
                  type="text"
                  className="org-input"
                  placeholder="Acme Corp"
                  value={newOrgForm.name}
                  onChange={(e) => setNewOrgForm({ ...newOrgForm, name: e.target.value, slug: e.target.value.toLowerCase().replace(/[^a-z0-9]/g, '-') })}
                  required
                />
              </div>

              <div className="org-input-group">
                <label>Slug (URL Key)</label>
                <input
                  type="text"
                  className="org-input"
                  placeholder="acme-corp"
                  value={newOrgForm.slug}
                  onChange={(e) => setNewOrgForm({ ...newOrgForm, slug: e.target.value })}
                  required
                />
              </div>

              <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '10px', marginTop: '20px' }}>
                <button type="button" className="org-tab-btn" onClick={() => setShowCreateOrgModal(false)}>Cancel</button>
                <button type="submit" className="btn-primary-org">Create Tenant</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Invite User Modal */}
      {showInviteModal && (
        <div style={{ position: 'fixed', top: 0, left: 0, right: 0, bottom: 0, background: 'rgba(0,0,0,0.8)', display: 'flex', justifyContent: 'center', alignItems: 'center', zIndex: 1000 }}>
          <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '10px', padding: '24px', width: '450px' }}>
            <h3 style={{ margin: '0 0 16px 0', color: '#f0f6fc' }}>Invite User to {currentOrg?.name}</h3>
            <form onSubmit={handleSendInvite}>
              <div className="org-input-group">
                <label>Email Address</label>
                <input
                  type="email"
                  className="org-input"
                  placeholder="engineer@company.com"
                  value={inviteForm.email}
                  onChange={(e) => setInviteForm({ ...inviteForm, email: e.target.value })}
                  required
                />
              </div>

              <div className="org-input-group">
                <label>Tenant Role</label>
                <select
                  className="org-select"
                  value={inviteForm.role}
                  onChange={(e) => setInviteForm({ ...inviteForm, role: e.target.value })}
                >
                  <option value="owner">Owner</option>
                  <option value="admin">Admin</option>
                  <option value="operator">Operator</option>
                  <option value="viewer">Viewer</option>
                  <option value="auditor">Auditor</option>
                  <option value="billing_admin">Billing Admin</option>
                </select>
              </div>

              <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '10px', marginTop: '20px' }}>
                <button type="button" className="org-tab-btn" onClick={() => setShowInviteModal(false)}>Cancel</button>
                <button type="submit" className="btn-primary-org">Send Invitation</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
