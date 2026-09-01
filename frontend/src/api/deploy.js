import apiClient from './client';

// Git Repositories
export const createGitRepository = async (data) => {
  const response = await apiClient.post('/deploy/repositories', data);
  return response.data;
};

export const getGitRepository = async (id) => {
  const response = await apiClient.get(`/deploy/repositories/${id}`);
  return response.data;
};

export const listGitRepositories = async () => {
  const response = await apiClient.get('/deploy/repositories');
  return response.data.repositories || [];
};

export const updateGitRepository = async (id, data) => {
  const response = await apiClient.patch(`/deploy/repositories/${id}`, data);
  return response.data;
};

export const deleteGitRepository = async (id) => {
  const response = await apiClient.delete(`/deploy/repositories/${id}`);
  return response.data;
};

// Pipelines
export const createPipeline = async (data) => {
  const response = await apiClient.post('/deploy/pipelines', data);
  return response.data;
};

export const getPipeline = async (id) => {
  const response = await apiClient.get(`/deploy/pipelines/${id}`);
  return response.data;
};

export const listPipelines = async (environment = '') => {
  const params = environment ? { environment } : {};
  const response = await apiClient.get('/deploy/pipelines', { params });
  return response.data.pipelines || [];
};

export const updatePipeline = async (id, data) => {
  const response = await apiClient.patch(`/deploy/pipelines/${id}`, data);
  return response.data;
};

export const deletePipeline = async (id) => {
  const response = await apiClient.delete(`/deploy/pipelines/${id}`);
  return response.data;
};

// Builds
export const createBuild = async (data) => {
  const response = await apiClient.post('/deploy/builds', data);
  return response.data;
};

export const getBuild = async (id) => {
  const response = await apiClient.get(`/deploy/builds/${id}`);
  return response.data;
};

export const listBuilds = async (pipelineId = '', limit = 20) => {
  const params = { limit };
  if (pipelineId) params.pipeline_id = pipelineId;
  const response = await apiClient.get('/deploy/builds', { params });
  return response.data.builds || [];
};

export const updateBuildStatus = async (id, status, errorMessage = '') => {
  const response = await apiClient.patch(`/deploy/builds/${id}/status`, {
    status,
    error_message: errorMessage,
  });
  return response.data;
};

// Deployments
export const createDeployment = async (data) => {
  const response = await apiClient.post('/deploy/deployments', data);
  return response.data;
};

export const getDeployment = async (id) => {
  const response = await apiClient.get(`/deploy/deployments/${id}`);
  return response.data;
};

export const listDeployments = async (environment = '', limit = 20) => {
  const params = { limit };
  if (environment) params.environment = environment;
  const response = await apiClient.get('/deploy/deployments', { params });
  return response.data.deployments || [];
};

export const rollbackDeployment = async (id, reason = '') => {
  const response = await apiClient.post(`/deploy/deployments/${id}/rollback`, {
    reason,
  });
  return response.data;
};

// Webhooks
export const handleWebhook = async (provider, payload, repositoryId) => {
  const response = await apiClient.post(
    `/deploy/webhooks/${provider}?repository_id=${repositoryId}`,
    payload,
    {
      headers: {
        'Content-Type': 'application/json',
      },
    }
  );
  return response.data;
};

// Statistics
export const getDeployStats = async () => {
  const response = await apiClient.get('/deploy/stats');
  return response.data;
};