import { apiGet } from './client.js';

export function getEnterpriseOverview() {
  return apiGet('/overview');
}
