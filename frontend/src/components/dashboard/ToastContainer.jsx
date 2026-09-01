import React from 'react';
import { useDashboardStore } from '../../store/dashboardStore.jsx';
import { X, AlertCircle, CheckCircle, Info, AlertTriangle } from 'lucide-react';

export default function ToastContainer() {
  const { toasts, removeToast } = useDashboardStore();

  return (
    <div className="toast-container">
      {toasts.map((toast) => {
        let Icon = Info;
        if (toast.type === 'critical') Icon = AlertCircle;
        if (toast.type === 'warning') Icon = AlertTriangle;
        if (toast.type === 'success') Icon = CheckCircle;

        return (
          <div key={toast.id} className={`toast-card toast-${toast.type}`}>
            <Icon className="toast-icon" size={18} />
            <div className="toast-content">
              <span className="toast-title">{toast.title}</span>
              <span className="toast-message">{toast.message}</span>
            </div>
            <button className="toast-close" onClick={() => removeToast(toast.id)} type="button">
              <X size={14} />
            </button>
          </div>
        );
      })}
    </div>
  );
}
