import React, { useState } from 'react';
import { Calendar, Plus, Clock, Play } from 'lucide-react';

export default function Scheduler() {
  const scheduledJobs = [
    { name: 'Nightly Database Snapshot Backup', cron: '0 2 * * *', frequency: 'Daily at 02:00 UTC', next_run: 'Tonight at 02:00' },
    { name: 'Weekly Log Rotation & Cleanup', cron: '0 3 * * 0', frequency: 'Every Sunday at 03:00 UTC', next_run: 'Next Sunday' },
    { name: 'Monthly System Cache Maintenance', cron: '0 4 1 * *', frequency: '1st of month at 04:00 UTC', next_run: '1st of next month' },
  ];

  return (
    <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '12px', padding: '24px' }}>
      <div style={{ marginBottom: '20px' }}>
        <h3 style={{ margin: 0, color: '#58a6ff', fontSize: '18px', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '8px' }}>
          <Calendar size={20} color="#58a6ff" />
          Scheduled Cron Automation Runner
        </h3>
        <span style={{ fontSize: '13px', color: '#8b949e' }}>Automate recurring backups, log rotations, maintenance, and compliance audits</span>
      </div>

      <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
        {scheduledJobs.map((j, idx) => (
          <div key={idx} style={{ background: '#0d1117', border: '1px solid #30363d', borderRadius: '10px', padding: '16px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <div>
              <h4 style={{ margin: '0 0 4px 0', color: '#f0f6fc', fontSize: '15px' }}>{j.name}</h4>
              <div style={{ fontSize: '12px', color: '#8b949e' }}>
                Schedule: <code style={{ color: '#58a6ff' }}>{j.cron}</code> ({j.frequency})
              </div>
            </div>

            <div style={{ textAlign: 'right' }}>
              <span style={{ padding: '3px 8px', borderRadius: '4px', background: 'rgba(63, 185, 80, 0.15)', color: '#3fb950', fontSize: '11px', fontWeight: 700 }}>
                SCHEDULED
              </span>
              <div style={{ fontSize: '11px', color: '#8b949e', marginTop: '4px' }}>Next run: {j.next_run}</div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
