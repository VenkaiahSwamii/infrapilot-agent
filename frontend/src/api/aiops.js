import { apiGet, apiPost } from "./client.js";

export const getAIOpsDashboard = () => apiGet("/aiops/dashboard");
export const getAnomalies = (machineId) => apiGet(`/aiops/anomalies${machineId ? `?machine_id=${machineId}` : ""}`);
export const getPredictions = () => apiGet("/aiops/predictions");
export const getForecasts = () => apiGet("/aiops/forecasts");
export const getRootCause = (incidentId) => apiGet(`/aiops/root-cause/${incidentId}`);
export const getRecommendations = () => apiGet("/aiops/recommendations");
export const executeRemediation = (recommendationId, mode) => apiPost("/aiops/remediations/execute", { recommendation_id: recommendationId, mode });
export const getHealthScore = () => apiGet("/aiops/health-score");
export const queryAIChat = (query) => apiPost("/aiops/chat", { query });
export const getCostOptimization = () => apiGet("/aiops/cost-optimization");
export const getSLAMetrics = () => apiGet("/aiops/sla");
export const getSecurityAnalytics = () => apiGet("/aiops/security");
