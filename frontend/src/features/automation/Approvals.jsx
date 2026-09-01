import React, { useState } from 'react';
import { Clock, Check, X, ShieldAlert } from 'lucide-react';
import { apiPost } from '../../api/client.js';

export default function Approvals({ approvals, onRefresh }) {
  const [processingId, setProcessingId] = useState(null);
  const list = approvals?.approvals || [];

  const handleDecision = async (reqId, approve) => {
    setProcessingId(reqId);
    try {
      await apiPost('/automation/approve', {
        request_id: reqId,
        approve,
        reason: approve ? 'Admin approved execution' : 'Rejected by security policy',
      });
      if (onRefresh) onRefresh();
    } catch (err) {
      console.error('Failed to process approval', err);
    } finally {
      setProcessingId(null);
    }
  };

  return (
    <div style={{ background: '#161b22', border: '1px solid #30363d', borderRadius: '12px', padding: '24px' }}>
      <div style={{ marginBottom: '20px' }}>
        <h3 style={{ margin: 0, color: '#ffa657', fontSize: '18px', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '8px' }}>
          <Clock size={20} color="#ffa657" />
          Pending Approval Requests Queue
        </h3>
        <span style={{ fontSize: '13px', color: '#8b949e' }}>Human-in-the-loop approval workflow for sensitive and destructive operational actions</span>
      </div>

      {list.length === 0 ? (
        <div style={{ color: '#8b949e', fontStyle: 'italic', padding: '16px 0' }}>No pending approval requests in queue.</div>
      ) : (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
          {list.map((req) => (
            <div key={req.id} style={{
              background: '#0d1117',
              border: '1px solid #ffa657',
              borderRadius: '10px',
              padding: '20px',
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center'
            }}>
              <div>
                <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                  <h4 style={{ margin: 0, color: '#f0f6fc', fontSize: '16px' }}>{req.rule_name}</h4>
                  <span style={{ padding: '2px 8px', borderRadius: '4px', background: 'rgba(255, 166, 87, 0.15)', color: '#ffa657', fontSize: '11px', fontWeight: 700 }}>
                    {req.status}
                  </span>
                </div>
                <div style={{ fontSize: '13px', color: '#c9d1d9', marginTop: '6px' }}>
                  Action: <strong style={{ color: '#58a6ff' }}>{req.action}</strong>
                </div>
                <div style={{ fontSize: '12px', color: '#8b949e', marginTop: '4px' }}>
                  Target: <code>{req.target_resource}</code> | Requested by: {req.requested_by}
                </div>
              </div>

              {req.status === 'PENDING' ? (
                <div style={{ display: 'flex', gap: '10px' }}>
                  <button
                    onClick={() => handleDecision(req.id, false)}
                    disabled={processingId === req.id}
                    style={{
                      backgroundColor: 'rgba(247, 129, 102, 0.15)',
                      color: '#f78166',
                      border: '1px solid #f78166',
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
                    <X size={16} /> Reject
                  </button>

                  <button
                    onClick={() => handleDecision(req.id, true)}
                    disabled={processingId === req.id}
                    style={{
                      backgroundColor: '#2ea043',
                      color: '#fff',
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
                    <Check size={16} /> Approve &amp; Execute
                  </button>
                </div>
              ) : (
                <span style={{ color: '#8b949e', fontSize: '13px', fontStyle: 'italic' }}>
                  Decided by {req.approver || 'Admin'}
                </span>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
