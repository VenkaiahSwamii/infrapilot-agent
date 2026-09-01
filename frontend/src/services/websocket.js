import { getWebSocketUrl } from '../websocket/liveEvents.js';

let socket = null;

export function connectWebSocket(onMessage) {
  const wsUrl = getWebSocketUrl();
  socket = new WebSocket(wsUrl);

  socket.onopen = () => {
    console.log("✅ WebSocket Connected");
  };

  socket.onmessage = (event) => {
    const data = JSON.parse(event.data);
    onMessage(data);
  };

  socket.onerror = (err) => {
    console.error(err);
  };

  socket.onclose = () => {
    console.log("WebSocket Closed");
  };

  return socket;
}

export function disconnectWebSocket() {
  if (socket) socket.close();
}