import { apiGet, apiPost, apiPatch, apiDelete } from "./client.js";

export const getOrganizations = () => apiGet("/organizations");
export const createOrganization = (payload) => apiPost("/organizations", payload);
export const getOrganization = (orgId) => apiGet(`/organizations/${orgId}`);
export const updateOrganization = (orgId, payload) => apiPatch(`/organizations/${orgId}`, payload);
export const deleteOrganization = (orgId) => apiDelete(`/organizations/${orgId}`);

export const getOrganizationUsers = (orgId) => apiGet(`/organizations/${orgId}/users`);
export const getOrganizationSettings = (orgId) => apiGet(`/organizations/${orgId}/settings`);
export const updateOrganizationSettings = (orgId, payload) => apiPatch(`/organizations/${orgId}/settings`, payload);

export const getOrganizationQuotas = (orgId) => apiGet(`/organizations/${orgId}/quotas`);
export const updateOrganizationQuotas = (orgId, payload) => apiPatch(`/organizations/${orgId}/quotas`, payload);

export const getOrganizationBilling = (orgId) => apiGet(`/organizations/${orgId}/billing`);
export const createInvitation = (orgId, payload) => apiPost(`/organizations/${orgId}/invitations`, payload);

export const getSuperAdminMetrics = () => apiGet("/admin/super-dashboard");
