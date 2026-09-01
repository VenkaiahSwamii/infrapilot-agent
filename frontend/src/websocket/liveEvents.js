const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';

export const WS_URL = import.meta.env.VITE_WS_URL || `${wsProtocol}//${window.location.hostname}:8080/ws`;

export function createLiveEventsSocket() {
  return new WebSocket(getWebSocketUrl());
}

export function getWebSocketUrl() {
  const token = localStorage.getItem('token');
  if (token) {
    // Add token query parameter for WS auth if required
    return `${WS_URL}?token=${encodeURIComponent(token)}`;
  }
  return WS_URL;
}

