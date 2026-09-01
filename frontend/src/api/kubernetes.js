import { apiGet, apiPost } from './client.js';

export function getKubernetesOverview(id) {
  return apiGet(`/kubernetes/overview/${id}`);
}

export function getMachineKubernetesData(id) {
  return apiGet(`/machines/${id}/kubernetes`);
}

export function getKubernetesClusters() {
  return apiGet('/kubernetes/clusters');
}

export function getKubernetesNodes(clusterId) {
  if (clusterId && clusterId !== 'all') {
    return apiGet(`/kubernetes/nodes/${clusterId}`);
  }
  return apiGet('/kubernetes/nodes');
}

export function getKubernetesPods(clusterId) {
  if (clusterId && clusterId !== 'all') {
    return apiGet(`/kubernetes/pods/${clusterId}`);
  }
  return apiGet('/kubernetes/pods');
}

export function getKubernetesDeployments(clusterId) {
  if (clusterId && clusterId !== 'all') {
    return apiGet(`/kubernetes/deployments/${clusterId}`);
  }
  return apiGet('/kubernetes/deployments');
}

export function getKubernetesServices(clusterId) {
  if (clusterId && clusterId !== 'all') {
    return apiGet(`/kubernetes/services/${clusterId}`);
  }
  return apiGet('/kubernetes/services');
}

export function getKubernetesStatefulSets(clusterId) {
  return apiGet(`/kubernetes/statefulsets/${clusterId || 'default'}`);
}

export function getKubernetesDaemonSets(clusterId) {
  return apiGet(`/kubernetes/daemonsets/${clusterId || 'default'}`);
}

export function getKubernetesNamespaces(clusterId) {
  return apiGet(`/kubernetes/namespaces/${clusterId || 'default'}`);
}

export function getKubernetesStorage(clusterId) {
  return apiGet(`/kubernetes/storage/${clusterId || 'default'}`);
}

export function getKubernetesEvents(clusterId) {
  return apiGet(`/kubernetes/events/${clusterId || 'default'}`);
}

export function getPodLogs(podName, namespace = 'default', tail = 100) {
  return apiGet(`/kubernetes/pods/logs?pod=${encodeURIComponent(podName)}&namespace=${encodeURIComponent(namespace)}&tail=${tail}`);
}

export function getResourceYAML(kind = 'deployment', name = 'api-gateway', namespace = 'default') {
  return apiGet(`/kubernetes/yaml?kind=${encodeURIComponent(kind)}&name=${encodeURIComponent(name)}&namespace=${encodeURIComponent(namespace)}`);
}
