import React from 'react';
import { ShieldCheck } from 'lucide-react';

export default function AlertsPanel({ alerts }) {
  return (
    <section className="alerts-panel">
      <div className="panel-head compact">
        <div>
          <span className="eyebrow">Alerts</span>
          <h2>Active Incidents</h2>
        </div>
        <span>{alerts.length}</span>
      </div>
      <div className="alert-list">
        {alerts.length === 0 ? (
          <div className="quiet-state">
            <ShieldCheck size={28} />
            <strong>All clear</strong>
            <p>No active alerts are reported by the backend.</p>
          </div>
        ) : (
          alerts.map((alert, index) => (
            <article className={`alert-item alert-${alert.severity || 'warning'}`} key={alert.id || index}>
              <div>
                <strong>{alert.title || alert.severity || 'warning'}</strong>
                <span>{alert.created_at ? new Date(alert.created_at).toLocaleTimeString() : 'now'}</span>
              </div>
              <p>{alert.message}</p>
            </article>
          ))
        )}
      </div>
    </section>
  );
}
