import React, { useState, useEffect } from 'react';
import { Building2, ChevronDown, Check, Plus } from 'lucide-react';
import { getOrganizations } from '../../api/organization.js';

export default function OrganizationSwitcher({ onSelectOrg }) {
  const [organizations, setOrganizations] = useState([]);
  const [selectedOrg, setSelectedOrg] = useState(null);
  const [isOpen, setIsOpen] = useState(false);

  useEffect(() => {
    loadOrganizations();
  }, []);

  const loadOrganizations = async () => {
    try {
      const res = await getOrganizations();
      const orgs = res.organizations || [];
      setOrganizations(orgs);

      const savedOrgId = localStorage.getItem('activeOrgId');
      const active = orgs.find(o => o.id === savedOrgId) || orgs[0] || { id: 'default', name: 'Default Organization' };
      setSelectedOrg(active);
      if (onSelectOrg) onSelectOrg(active);
    } catch (err) {
      console.error('Failed to load organizations for switcher', err);
      const fallback = { id: 'default', name: 'Default Organization' };
      setSelectedOrg(fallback);
    }
  };

  const handleSelect = (org) => {
    setSelectedOrg(org);
    localStorage.setItem('activeOrgId', org.id);
    setIsOpen(false);
    if (onSelectOrg) onSelectOrg(org);
    window.location.reload();
  };

  return (
    <div style={{ position: 'relative', display: 'inline-block' }}>
      <button
        onClick={() => setIsOpen(!isOpen)}
        style={{
          display: 'flex',
          alignItems: 'center',
          gap: '8px',
          background: '#161b22',
          border: '1px solid #30363d',
          borderRadius: '6px',
          padding: '6px 12px',
          color: '#f0f6fc',
          fontSize: '13px',
          fontWeight: 600,
          cursor: 'pointer',
        }}
      >
        <Building2 size={16} color="#58a6ff" />
        <span>{selectedOrg?.name || 'Select Organization'}</span>
        <ChevronDown size={14} color="#8b949e" />
      </button>

      {isOpen && (
        <div style={{
          position: 'absolute',
          top: '100%',
          left: 0,
          marginTop: '6px',
          width: '220px',
          background: '#161b22',
          border: '1px solid #30363d',
          borderRadius: '8px',
          boxShadow: '0 8px 24px rgba(0,0,0,0.4)',
          zIndex: 1000,
          padding: '6px 0',
        }}>
          <div style={{ padding: '6px 12px', fontSize: '11px', fontWeight: 700, color: '#8b949e', textTransform: 'uppercase' }}>
            Tenants & Organizations
          </div>
          {organizations.map((org) => (
            <div
              key={org.id}
              onClick={() => handleSelect(org)}
              style={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                padding: '8px 12px',
                fontSize: '13px',
                color: '#c9d1d9',
                cursor: 'pointer',
                backgroundColor: selectedOrg?.id === org.id ? '#21262d' : 'transparent',
              }}
            >
              <span style={{ fontWeight: selectedOrg?.id === org.id ? 600 : 400 }}>{org.name}</span>
              {selectedOrg?.id === org.id && <Check size={14} color="#3fb950" />}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
