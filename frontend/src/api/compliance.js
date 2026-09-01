import { apiGet, apiPost } from "./client.js";

export const getSecurityDashboard = () => apiGet("/compliance/dashboard");
export const getFrameworks = () => apiGet("/compliance/frameworks");
export const getAuditLogs = () => apiGet("/compliance/audit-logs");
export const setupMFA = () => apiPost("/compliance/mfa/setup", {});
export const verifyMFA = (secret, code) => apiPost("/compliance/mfa/verify", { secret, code });
export const getPolicies = () => apiGet("/compliance/policies");
export const getCertificates = () => apiGet("/compliance/certificates");
export const getSecrets = () => apiGet("/compliance/secrets");
export const exportComplianceReport = (framework, format) => apiPost("/compliance/reports/export", { framework, format });
