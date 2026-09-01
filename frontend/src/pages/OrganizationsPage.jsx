import React from 'react';
import OrganizationManagement from '../features/organization/OrganizationManagement.jsx';

export default function OrganizationsPage() {
  return (
    <div className="organizations-page-wrapper">
      <OrganizationManagement />
      <style>{`
        /* Remove double body scrollbars/background conflicts */
        .organizations-page-wrapper > div {
          min-height: auto !important;
          background-color: transparent !important;
          padding: 0 !important;
        }
      `}</style>
    </div>
  );
}
