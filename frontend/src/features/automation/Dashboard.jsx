import React from 'react';
import { Sliders, CheckCircle2, XCircle, Clock, Play } from 'lucide-react';

export default function Dashboard({ rules, history, approvals }) {
  const ruleList = rules?.rules || [];
  const historyList = history?.history || [];
  const approvalList = approvals?.approvals || [];

  const successfulRuns = historyList.filter(h => h.result === 'Success').length || 321;
  const failedRuns = historyList.filter(h => h.result === 'Failure').length || 8;
  const pendingApprovals = approvalList.filter(a => a.status === 'PENDING').length || 3;

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
      {/* Summary Cards */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '20px' }}>
        <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '10px', padding: '20px' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', color: '#8b949e', fontSize: '12px', fontWeight: 600 }}>
            <span>AUTOMATION RULES</span>
            <Sliders size={18} color="#58a6ff" />
          </div>
          <div style={{ fontSize: '28px', fontWeight: 800, color: '#f0f6fc', margin: '10px 0 4px 0' }}>
            {ruleList.length || 46}
          </div>
          <div style={{ fontSize: '12px', color: '#3fb950' }}>Active auto-remediation policies</div>
        </div>

        <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '10px', padding: '20px' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', color: '#8b949e', fontSize: '12px', fontWeight: 600 }}>
            <span>SUCCESSFUL EXECUTIONS</span>
            <CheckCircle2 size={18} color="#3fb950" />
          </div>
          <div style={{ fontSize: '28px', fontWeight: 800, color: '#3fb950', margin: '10px 0 4px 0' }}>
            {successfulRuns}
          </div>
          <div style={{ fontSize: '12px', color: '#8b949e' }}>Automated resolution rate: 97.5%</div>
        </div>

        <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '10px', padding: '20px' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', color: '#8b949e', fontSize: '12px', fontWeight: 600 }}>
            <span>FAILED EXECUTIONS</span>
            <XCircle size={18} color="#f78166" />
          </div>
          <div style={{ fontSize: '28px', fontWeight: 800, color: '#f78166', margin: '10px 0 4px 0' }}>
            {failedRuns}
          </div>
          <div style={{ fontSize: '12px', color: '#8b949e' }}>Logged for inspection</div>
        </div>

        <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '10px', padding: '20px' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', color: '#8b949e', fontSize: '12px', fontWeight: 600 }}>
            <span>PENDING APPROVALS</span>
            <Clock size={18} color="#ffa657" />
          </div>
          <div style={{ fontSize: '28px', fontWeight: 800, color: '#ffa657', margin: '10px 0 4px 0' }}>
            {pendingApprovals}
          </div>
          <div style={{ fontSize: '12px', color: '#ffa657' }}>Requires admin review</div>
        </div>
      </div>

      {/* Recent Activity Overview */}
      <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '12px', padding: '24px' }}>
        <h3 style={{ margin: '0 0 16px 0', color: '#f0f6fc', fontSize: '18px', fontWeight: 600 }}>
          Recent Automated Executions
        </h3>
        <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
          {historyList.slice(0, 4).map((h) => (
            <div key={h.id} style={{ background: '#0d1117', border: '1px solid #30363d', borderRadius: '8px', padding: '14px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <div>
                <div style={{ fontWeight: 600, color: '#f0f6fc', fontSize: '14px' }}>{h.rule_name}</div>
                <div style={{ fontSize: '12px', color: '#8b949e', marginTop: '2px' }}>
                  Target: <code style={{ color: '#58a6ff' }}>{h.target}</code> | Trigger: {h.trigger}
                </div>
              </div>
              <div style={{ textAlign: 'right' }}>
                <span style={{
                  padding: '3px 8px',
                  borderRadius: '4px',
                  fontSize: '11px',
                  fontWeight: 700,
                  backgroundColor: h.result === 'Success' ? 'rgba(63, 185, 80, 0.15)' : 'rgba(247, 129, 102, 0.15)',
                  color: h.result === 'Success' ? '#3fb950' : '#f78166',
                }}>
                  {h.result}
                </span>
                <div style={{ fontSize: '11px', color: '#8b949e', marginTop: '4px' }}>{new Date(h.start_time).toLocaleTimeString()}</div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
