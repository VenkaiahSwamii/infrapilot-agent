import { apiGet, apiPost } from './client.js';

export function createEnrollmentToken(organizationId = 'default') {
  return apiPost(`/organizations/${organizationId}/enrollment-tokens`, {
    name: 'Dashboard generated enrollment token',
    max_uses: 1000,
  });
}

export function listEnrollmentTokens(organizationId = 'default') {
  return apiGet(`/organizations/${organizationId}/enrollment-tokens`);
}
