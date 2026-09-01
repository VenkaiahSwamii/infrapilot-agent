import { apiClient, apiGet, apiPost, apiPatch, apiDelete } from './client.js';

// Incident CRUD operations
export const incidentApi = {
  // List incidents with filtering
  getIncidents: (filters = {}) => {
    const params = new URLSearchParams();
    if (filters.machine_id) params.append('machine_id', filters.machine_id);
    if (filters.status) params.append('status', filters.status);
    if (filters.severity) params.append('severity', filters.severity);
    
    return apiGet(`/incidents?${params.toString()}`);
  },

  // Get single incident
  getIncident: (id) => apiGet(`/incidents/${id}`),

  // Create incident manually
  createIncident: (data) => apiPost('/incidents', data),

  // Update incident
  updateIncident: (id, data) => apiPatch(`/incidents/${id}`, data),

  // Delete incident
  deleteIncident: (id) => apiDelete(`/incidents/${id}`),

  // Assign incident to user
  assignIncident: (id, userId) => apiPost(`/incidents/${id}/assign`, { user_id: userId }),

  // Add comment
  addComment: (id, comment) => apiPost(`/incidents/${id}/comment`, { comment }),

  // Resolve incident
  resolveIncident: (id, resolutionNote) => apiPost(`/incidents/${id}/resolve`, { resolution_note: resolutionNote }),

  // Close incident
  closeIncident: (id) => apiPost(`/incidents/${id}/close`),

  // Get incident timeline
  getTimeline: (id) => apiGet(`/incidents/${id}/timeline`),

  // Get incident alerts
  getAlerts: (id) => apiGet(`/incidents/${id}/alerts`),

  // Get incident analytics
  getAnalytics: () => apiGet('/incidents/analytics'),

  // Get AI summary
  getAISummary: (id) => apiGet(`/incidents/${id}/summary`),
};

export default incidentApi;