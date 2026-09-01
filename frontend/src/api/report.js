const TOKEN_KEY = 'token';

export const reportApi = {
  async request(endpoint, options = {}) {
    const token = localStorage.getItem(TOKEN_KEY);
    const headers = {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...options.headers,
    };

    const res = await fetch(`/api/v1${endpoint}`, {
      ...options,
      headers,
    });

    if (!res.ok) {
      let msg = `HTTP ${res.status}`;
      try {
        const err = await res.json();
        msg = err.error || msg;
      } catch (_) {}
      throw new Error(msg);
    }

    if (res.status === 204) {
      return null;
    }

    const contentType = res.headers.get('content-type') || '';
    if (contentType.includes('application/json')) {
      return res.json();
    }
    return res.blob();
  },

  list() {
    return this.request('/reports', { method: 'GET' });
  },

  generate(payload) {
    return this.request('/reports/generate', {
      method: 'POST',
      body: JSON.stringify(payload),
    });
  },

  download(id) {
    return this.request(`/reports/${id}/download`, { method: 'GET' });
  },

  delete(id) {
    return this.request(`/reports/${id}`, { method: 'DELETE' });
  },
};
