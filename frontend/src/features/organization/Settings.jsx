import React, { useState, useEffect } from 'react';
import {
  Palette,
  Save,
  Clock,
  Shield,
  Sliders,
  Globe,
  CheckCircle2,
  Building,
  Image as ImageIcon,
  Sparkles,
  ExternalLink,
  RotateCcw,
  Check
} from 'lucide-react';

export default function Settings({ settings, onSave }) {
  const [form, setForm] = useState({
    company_name: 'InfraPilot Enterprise',
    logo_url: '',
    timezone: 'UTC',
    primary_color: '#3b82f6',
    sidebar_color: '#0d1424',
    custom_domain: 'monitoring.company.internal',
    license_tier: 'enterprise',
  });
  const [savedMsg, setSavedMsg] = useState(false);

  useEffect(() => {
    if (settings) {
      setForm({
        company_name: settings.company_name || 'InfraPilot Enterprise',
        logo_url: settings.logo_url || '',
        timezone: settings.timezone || 'UTC',
        primary_color: settings.primary_color || '#3b82f6',
        sidebar_color: settings.sidebar_color || '#0d1424',
        custom_domain: settings.custom_domain || 'monitoring.company.internal',
        license_tier: settings.license_tier || 'enterprise',
      });
    }
  }, [settings]);

  const colorPresets = [
    { name: 'Cyber Cyan', primary: '#06b6d4', sidebar: '#08131e' },
    { name: 'Emerald Pro', primary: '#10b981', sidebar: '#081c15' },
    { name: 'Royal Purple', primary: '#a855f7', sidebar: '#140c24' },
    { name: 'Electric Blue', primary: '#3b82f6', sidebar: '#0d1424' },
    { name: 'Sunset Amber', primary: '#f59e0b', sidebar: '#1c1308' },
  ];

  const applyColorPreset = (preset) => {
    setForm((prev) => ({
      ...prev,
      primary_color: preset.primary,
      sidebar_color: preset.sidebar,
    }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    if (onSave) onSave(form);
    setSavedMsg(true);
    setTimeout(() => setSavedMsg(false), 3000);
  };

  return (
    <div className="whitelabel-settings-root">
      {/* Header Row */}
      <div className="settings-header-box">
        <div className="title-row">
          <div className="icon-badge">
            <Palette size={20} color="#38bdf8" />
          </div>
          <div>
            <h3>White-Label & Organization Branding</h3>
            <p>Customize enterprise portal identity, custom CNAME domains, visual themes, and timezones.</p>
          </div>
        </div>
      </div>

      {savedMsg && (
        <div className="save-alert-success">
          <CheckCircle2 size={16} />
          <span>Organization and branding configurations saved successfully!</span>
        </div>
      )}

      {/* Main Form & Live Preview Grid */}
      <div className="settings-content-grid">
        {/* Left: Settings Form */}
        <form onSubmit={handleSubmit} className="settings-form-card">
          <div className="form-section-title">
            <Building size={15} color="#38bdf8" />
            <span>ORGANIZATION IDENTITY & DOMAIN</span>
          </div>

          <div className="form-group">
            <label>COMPANY DISPLAY NAME</label>
            <input
              type="text"
              value={form.company_name}
              onChange={(e) => setForm({ ...form, company_name: e.target.value })}
              placeholder="e.g. Acme Cloud Systems Inc."
              required
            />
          </div>

          <div className="form-group">
            <label>CUSTOM LOGO URL (PNG / SVG)</label>
            <input
              type="url"
              value={form.logo_url}
              onChange={(e) => setForm({ ...form, logo_url: e.target.value })}
              placeholder="https://assets.company.com/brand/logo.svg"
            />
          </div>

          <div className="form-row-2">
            <div className="form-group">
              <label>CUSTOM DOMAIN (CNAME)</label>
              <input
                type="text"
                value={form.custom_domain}
                onChange={(e) => setForm({ ...form, custom_domain: e.target.value })}
                placeholder="monitoring.company.com"
              />
            </div>

            <div className="form-group">
              <label>DEFAULT TIMEZONE</label>
              <select
                value={form.timezone}
                onChange={(e) => setForm({ ...form, timezone: e.target.value })}
              >
                <option value="UTC">UTC (Coordinated Universal Time)</option>
                <option value="America/New_York">America/New_York (EST / EDT)</option>
                <option value="Europe/London">Europe/London (GMT / BST)</option>
                <option value="Asia/Kolkata">Asia/Kolkata (IST)</option>
                <option value="Asia/Tokyo">Asia/Tokyo (JST)</option>
                <option value="Australia/Sydney">Australia/Sydney (AEST)</option>
              </select>
            </div>
          </div>

          <div className="form-section-title" style={{ marginTop: '12px' }}>
            <Palette size={15} color="#a855f7" />
            <span>THEME & PALETTE CUSTOMIZATION</span>
          </div>

          {/* Color Preset Swatches */}
          <div className="presets-row">
            <span className="presets-label">Quick Themes:</span>
            <div className="presets-list">
              {colorPresets.map((p) => (
                <button
                  key={p.name}
                  type="button"
                  className="preset-swatch-btn"
                  onClick={() => applyColorPreset(p)}
                >
                  <span className="swatch-dot" style={{ backgroundColor: p.primary }} />
                  <span>{p.name}</span>
                </button>
              ))}
            </div>
          </div>

          <div className="form-row-2">
            <div className="form-group">
              <label>PRIMARY ACCENT COLOR</label>
              <div className="color-picker-wrap">
                <input
                  type="color"
                  value={form.primary_color}
                  onChange={(e) => setForm({ ...form, primary_color: e.target.value })}
                  className="color-input"
                />
                <span className="color-hex">{form.primary_color}</span>
              </div>
            </div>

            <div className="form-group">
              <label>SIDEBAR BASE COLOR</label>
              <div className="color-picker-wrap">
                <input
                  type="color"
                  value={form.sidebar_color}
                  onChange={(e) => setForm({ ...form, sidebar_color: e.target.value })}
                  className="color-input"
                />
                <span className="color-hex">{form.sidebar_color}</span>
              </div>
            </div>
          </div>

          <div className="form-footer">
            <button type="submit" className="btn-save-settings">
              <Save size={15} /> Save Organization Settings
            </button>
          </div>
        </form>

        {/* Right: Live Interactive Branding Preview */}
        <div className="preview-container-card">
          <div className="preview-header-label">
            <Sparkles size={14} color="#38bdf8" />
            <span>LIVE WHITE-LABEL PORTAL PREVIEW</span>
          </div>

          <div className="mockup-browser-frame">
            <div className="mockup-browser-bar">
              <div className="mockup-dots">
                <span className="dot red" />
                <span className="dot yellow" />
                <span className="dot green" />
              </div>
              <div className="mockup-url-bar">
                <Globe size={11} color="#64748b" />
                <span>https://{form.custom_domain || 'monitoring.company.internal'}</span>
              </div>
            </div>

            <div className="mockup-portal-body">
              {/* Mockup Topbar */}
              <div className="mockup-topbar" style={{ backgroundColor: form.sidebar_color }}>
                <div className="mockup-brand">
                  <div
                    className="mockup-logo-box"
                    style={{ backgroundColor: form.primary_color }}
                  >
                    <Building size={14} color="#ffffff" />
                  </div>
                  <strong>{form.company_name || 'InfraPilot Enterprise'}</strong>
                </div>

                <div className="mockup-nav-pills">
                  <span
                    className="mockup-pill active"
                    style={{ color: form.primary_color, borderColor: form.primary_color }}
                  >
                    Infrastructure
                  </span>
                  <span className="mockup-pill">APM</span>
                  <span className="mockup-pill">Alerts</span>
                </div>
              </div>

              {/* Mockup Dashboard Content */}
              <div className="mockup-content-sample">
                <div className="mockup-metric-card">
                  <span className="m-label">ACTIVE FLEET</span>
                  <strong style={{ color: form.primary_color }}>3 Connected</strong>
                  <div
                    className="m-line"
                    style={{ backgroundColor: form.primary_color }}
                  />
                </div>

                <div className="mockup-metric-card">
                  <span className="m-label">HEALTH SCORE</span>
                  <strong style={{ color: '#22c55e' }}>99.9% Uptime</strong>
                  <div className="m-line" style={{ backgroundColor: '#22c55e' }} />
                </div>
              </div>

              <div className="mockup-footer-note">
                <CheckCircle2 size={13} color="#22c55e" />
                <span>Custom SSL &amp; Enterprise RBAC Active</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <style>{`
        .whitelabel-settings-root {
          display: flex;
          flex-direction: column;
          gap: 20px;
        }
        .settings-header-box {
          background-color: #0d1424;
          border: 1px solid #1a253a;
          border-radius: 12px;
          padding: 18px 24px;
        }
        .title-row {
          display: flex;
          align-items: center;
          gap: 12px;
        }
        .icon-badge {
          width: 36px;
          height: 36px;
          border-radius: 10px;
          background: rgba(56, 189, 248, 0.15);
          display: flex;
          align-items: center;
          justify-content: center;
        }
        .title-row h3 {
          margin: 0;
          font-size: 17px;
          font-weight: 700;
          color: #ffffff;
        }
        .title-row p {
          margin: 3px 0 0 0;
          font-size: 12.5px;
          color: #94a3b8;
        }
        .save-alert-success {
          display: flex;
          align-items: center;
          gap: 8px;
          background-color: rgba(34, 197, 94, 0.12);
          border: 1px solid rgba(34, 197, 94, 0.3);
          color: #4ade80;
          padding: 10px 16px;
          border-radius: 8px;
          font-size: 13px;
          font-weight: 600;
        }
        .settings-content-grid {
          display: grid;
          grid-template-columns: 1.2fr 1fr;
          gap: 20px;
        }
        @media (max-width: 1024px) {
          .settings-content-grid {
            grid-template-columns: 1fr;
          }
        }
        .settings-form-card, .preview-container-card {
          background-color: #0d1424;
          border: 1px solid #1a253a;
          border-radius: 12px;
          padding: 22px 24px;
          display: flex;
          flex-direction: column;
          gap: 16px;
          box-shadow: 0 4px 20px rgba(0,0,0,0.25);
        }
        .form-section-title {
          display: flex;
          align-items: center;
          gap: 6px;
          font-size: 11px;
          font-weight: 700;
          color: #64748b;
          letter-spacing: 0.5px;
        }
        .form-group {
          display: flex;
          flex-direction: column;
          gap: 6px;
        }
        .form-group label {
          font-size: 11px;
          font-weight: 700;
          color: #94a3b8;
          letter-spacing: 0.5px;
        }
        .form-group input, .form-group select {
          background-color: #080c14;
          border: 1px solid #1c283d;
          border-radius: 8px;
          padding: 9px 12px;
          color: #ffffff;
          font-size: 13px;
          outline: none;
        }
        .form-group input:focus, .form-group select:focus {
          border-color: #3b82f6;
        }
        .form-row-2 {
          display: grid;
          grid-template-columns: 1fr 1fr;
          gap: 14px;
        }
        .presets-row {
          display: flex;
          align-items: center;
          gap: 10px;
          flex-wrap: wrap;
        }
        .presets-label {
          font-size: 11.5px;
          color: #64748b;
          font-weight: 600;
        }
        .presets-list {
          display: flex;
          gap: 6px;
          flex-wrap: wrap;
        }
        .preset-swatch-btn {
          display: inline-flex;
          align-items: center;
          gap: 6px;
          background-color: #080c14;
          border: 1px solid #1c283d;
          color: #cbd5e1;
          padding: 4px 10px;
          border-radius: 6px;
          font-size: 11.5px;
          cursor: pointer;
          transition: all 0.15s ease;
        }
        .preset-swatch-btn:hover {
          border-color: #3b82f6;
          color: #ffffff;
        }
        .swatch-dot {
          width: 8px;
          height: 8px;
          border-radius: 50%;
        }
        .color-picker-wrap {
          display: flex;
          align-items: center;
          gap: 10px;
          background-color: #080c14;
          border: 1px solid #1c283d;
          border-radius: 8px;
          padding: 6px 12px;
        }
        .color-input {
          width: 32px;
          height: 28px;
          border: none;
          background: transparent;
          cursor: pointer;
          border-radius: 4px;
        }
        .color-hex {
          font-family: monospace;
          font-size: 12px;
          color: #94a3b8;
        }
        .form-footer {
          display: flex;
          justify-content: flex-end;
          margin-top: 10px;
        }
        .btn-save-settings {
          display: inline-flex;
          align-items: center;
          gap: 7px;
          background: linear-gradient(135deg, #10b981, #059669);
          color: #ffffff;
          border: none;
          padding: 10px 20px;
          border-radius: 8px;
          font-size: 13px;
          font-weight: 700;
          cursor: pointer;
          box-shadow: 0 4px 14px rgba(16, 185, 129, 0.35);
          transition: all 0.2s ease;
        }
        .btn-save-settings:hover {
          background: linear-gradient(135deg, #059669, #047857);
          transform: translateY(-1px);
        }
        .preview-header-label {
          display: flex;
          align-items: center;
          gap: 6px;
          font-size: 11px;
          font-weight: 700;
          color: #38bdf8;
          letter-spacing: 0.5px;
        }
        .mockup-browser-frame {
          background-color: #080c14;
          border: 1px solid #1c283d;
          border-radius: 10px;
          overflow: hidden;
          box-shadow: 0 10px 30px rgba(0,0,0,0.5);
        }
        .mockup-browser-bar {
          background-color: #0f172a;
          border-bottom: 1px solid #1e293b;
          padding: 8px 12px;
          display: flex;
          align-items: center;
          gap: 12px;
        }
        .mockup-dots {
          display: flex;
          gap: 5px;
        }
        .mockup-dots .dot {
          width: 8px;
          height: 8px;
          border-radius: 50%;
        }
        .mockup-dots .dot.red { background: #ef4444; }
        .mockup-dots .dot.yellow { background: #f59e0b; }
        .mockup-dots .dot.green { background: #22c55e; }
        .mockup-url-bar {
          background-color: #080c14;
          border-radius: 4px;
          padding: 2px 10px;
          font-size: 10.5px;
          color: #94a3b8;
          display: flex;
          align-items: center;
          gap: 6px;
          flex: 1;
        }
        .mockup-portal-body {
          padding: 16px;
          display: flex;
          flex-direction: column;
          gap: 12px;
        }
        .mockup-topbar {
          border-radius: 8px;
          padding: 10px 14px;
          display: flex;
          justify-content: space-between;
          align-items: center;
          border: 1px solid rgba(255,255,255,0.06);
        }
        .mockup-brand {
          display: flex;
          align-items: center;
          gap: 8px;
        }
        .mockup-logo-box {
          width: 24px;
          height: 24px;
          border-radius: 6px;
          display: flex;
          align-items: center;
          justify-content: center;
        }
        .mockup-brand strong {
          font-size: 12.5px;
          color: #ffffff;
        }
        .mockup-nav-pills {
          display: flex;
          gap: 6px;
        }
        .mockup-pill {
          font-size: 10px;
          color: #94a3b8;
          padding: 2px 6px;
          border-radius: 4px;
          border: 1px solid transparent;
        }
        .mockup-pill.active {
          font-weight: 700;
          background: rgba(255,255,255,0.05);
        }
        .mockup-content-sample {
          display: grid;
          grid-template-columns: 1fr 1fr;
          gap: 10px;
        }
        .mockup-metric-card {
          background-color: #0c1220;
          border: 1px solid #1a253a;
          border-radius: 8px;
          padding: 10px;
          display: flex;
          flex-direction: column;
          gap: 4px;
        }
        .m-label {
          font-size: 9.5px;
          color: #64748b;
          font-weight: 700;
        }
        .mockup-metric-card strong {
          font-size: 14px;
        }
        .m-line {
          width: 100%;
          height: 3px;
          border-radius: 2px;
          margin-top: 4px;
        }
        .mockup-footer-note {
          display: flex;
          align-items: center;
          gap: 6px;
          font-size: 11px;
          color: #94a3b8;
          margin-top: 4px;
        }
      `}</style>
    </div>
  );
}
