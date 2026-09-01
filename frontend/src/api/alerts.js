import { apiGet, apiPost, apiDelete, apiPatch } from './client.js';

export function listAlerts(params = {}) {
  const query = new URLSearchParams();
  if (params.status) query.append('status', params.status);
  if (params.severity) query.append('severity', params.severity);
  if (params.category) query.append('category', params.category);
  if (params.machine_id) query.append('machine_id', params.machine_id);
  if (params.search) query.append('search', params.search);
  if (params.limit) query.append('limit', params.limit);
  if (params.offset) query.append('offset', params.offset);

  const qs = query.toString();
  return apiGet(`/alerts${qs ? `?${qs}` : ''}`);
}

export function getAlertStats() {
  return apiGet('/alerts/stats');
}

export function acknowledgeAlert(alertId) {
  return apiPost(`/alerts/${alertId}/ack`);
}

export function resolveAlert(alertId, resolutionNote = 'Resolved manually by operator') {
  return apiPost(`/alerts/${alertId}/resolve`, { resolution_note: resolutionNote });
}

export function silenceAlert(alertId, durationMinutes = 60, reason = 'Scheduled maintenance') {
  return apiPost(`/alerts/${alertId}/silence`, { duration_minutes: durationMinutes, reason });
}

export function bulkAcknowledgeAlerts(payload) {
  return apiPost('/alerts/bulk/ack', payload);
}

export function bulkResolveAlerts(payload) {
  return apiPost('/alerts/bulk/resolve', payload);
}

export function purgeResolvedAlerts() {
  return apiDelete('/alerts/purge-resolved');
}

export function cleanupDuplicateAlerts() {
  return apiPost('/alerts/cleanup-duplicates');
}

export function aiAnalyzeAlert(alertId) {
  return apiPost(`/alerts/${alertId}/ai-analyze`);
}

export function remediateAlert(alertId) {
  return apiPost(`/alerts/${alertId}/remediate`);
}

export function listAlertRules(orgId = '') {
  if (orgId) {
    return apiGet(`/orgs/${orgId}/alert-rules`);
  }
  return apiGet('/orgs/default/alert-rules').catch(() => []);
}

export function createAlertRule(orgId = 'default', data) {
  return apiPost(`/orgs/${orgId}/alert-rules`, data);
}

export function updateAlertRule(orgId = 'default', ruleId, data) {
  return apiPatch(`/orgs/${orgId}/alert-rules/${ruleId}`, data);
}

export function deleteAlertRule(orgId = 'default', ruleId) {
  return apiDelete(`/orgs/${orgId}/alert-rules/${ruleId}`);
}
