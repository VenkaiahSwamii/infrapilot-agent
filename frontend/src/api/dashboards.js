import { apiClient } from './client';

export const dashboardsApi = {
  // Get all dashboards
  getDashboards: async (params = {}) => {
    const response = await apiClient.get('/dashboards', { params });
    return response.data;
  },

  // Get default dashboard
  getDefaultDashboard: async () => {
    const response = await apiClient.get('/dashboards/default');
    return response.data;
  },

  // Get dashboard by ID
  getDashboard: async (id) => {
    const response = await apiClient.get(`/dashboards/${id}`);
    return response.data;
  },

  // Create dashboard
  createDashboard: async (data) => {
    const response = await apiClient.post('/dashboards', data);
    return response.data;
  },

  // Update dashboard
  updateDashboard: async (id, data) => {
    const response = await apiClient.put(`/dashboards/${id}`, data);
    return response.data;
  },

  // Delete dashboard
  deleteDashboard: async (id) => {
    const response = await apiClient.delete(`/dashboards/${id}`);
    return response.data;
  },

  // Set default dashboard
  setDefaultDashboard: async (id) => {
    const response = await apiClient.post(`/dashboards/${id}/set-default`);
    return response.data;
  },

  // Duplicate dashboard
  duplicateDashboard: async (id, name) => {
    const response = await apiClient.post(`/dashboards/${id}/duplicate`, { name });
    return response.data;
  },

  // Share dashboard
  shareDashboard: async (id, data) => {
    const response = await apiClient.post(`/dashboards/${id}/share`, data);
    return response.data;
  },

  // Get widgets
  getWidgets: async (dashboardId) => {
    const response = await apiClient.get(`/dashboards/${dashboardId}/widgets`);
    return response.data;
  },

  // Create widget
  createWidget: async (dashboardId, data) => {
    const response = await apiClient.post(`/dashboards/${dashboardId}/widgets`, data);
    return response.data;
  },

  // Update widget
  updateWidget: async (id, data) => {
    const response = await apiClient.patch(`/widgets/${id}`, data);
    return response.data;
  },

  // Delete widget
  deleteWidget: async (id) => {
    const response = await apiClient.delete(`/widgets/${id}`);
    return response.data;
  },

  // Get templates
  getTemplates: async (category = '') => {
    const response = await apiClient.get('/dashboard-templates', {
      params: { category }
    });
    return response.data;
  },

  // Get template by ID
  getTemplate: async (id) => {
    const response = await apiClient.get(`/dashboard-templates/${id}`);
    return response.data;
  },

  // Create dashboard from template
  createFromTemplate: async (templateId, name) => {
    const template = await dashboardsApi.getTemplate(templateId);
    const dashboard = await dashboardsApi.createDashboard({
      name,
      template_id: templateId,
    });

    // Add widgets from template
    if (template.data.widgets) {
      for (const widget of template.data.widgets) {
        await dashboardsApi.createWidget(dashboard.data.id, widget);
      }
    }

    return dashboard;
  },
};

export default dashboardsApi;