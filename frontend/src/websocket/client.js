import { WEBSOCKET_EVENTS } from './events.js';

class EnterpriseWebSocketClient {
  constructor() {
    this.socket = null;
    this.reconnectCount = 0;
    this.eventListeners = new Map();
    this.onStatusChangeCallback = null;
    this.onMetricsUpdateCallback = null;
    this.reconnectTimeout = null;
    this.pingInterval = null;
    this.msgCount = 0;
    this.rateInterval = null;
    this.onStatsChangeCallback = null;
  }

  connect(url) {
    if (this.socket) {
      if (this.socket.readyState === WebSocket.OPEN || this.socket.readyState === WebSocket.CONNECTING) {
        console.log("WebSocket connection already open or connecting. Skipping new connection.");
        return;
      }
      this.socket.onclose = null; // Prevent the old socket from triggering reconnect
      this.socket.close();
    }

    this.socket = new WebSocket(url);

    this.socket.onopen = () => {
      this.updateStatus('Connected');
      this.subscribe('global');
      this.startPingLoop();
      this.startRateCalculation();
    };

    this.socket.onmessage = (event) => {
      this.msgCount++;

      const payloads = event.data.split("\n");

      for (const raw of payloads) {
        if (!raw.trim()) continue;

        try {
          const message = JSON.parse(raw);
          this.dispatchEvent(message);
        } catch (err) {
          console.error("Failed to parse WebSocket message chunk:", raw, err);
        }
      }
    };

    this.socket.onclose = () => {
      this.updateStatus('Disconnected');
      this.stopPingLoop();
      this.stopRateCalculation();
      this.reconnect(url);
    };

    this.socket.onerror = (err) => {
      console.error("WS ERROR");
      console.log("readyState =", this.socket.readyState);
      console.log("url =", this.socket.url);
      console.log(err);
    };
  }

  reconnect(url) {
    if (this.reconnectTimeout) return;
    this.reconnectCount++;
    this.triggerStatsChange({ reconnects: this.reconnectCount });
    this.reconnectTimeout = setTimeout(() => {
      this.reconnectTimeout = null;
      this.connect(url);
    }, 5000);
  }

  subscribe(room) {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify({ event: 'subscribe', room }));
    }
  }

  unsubscribe(room) {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify({ event: 'unsubscribe', room }));
    }
  }

  startPingLoop() {
    this.pingInterval = setInterval(() => {
      if (this.socket && this.socket.readyState === WebSocket.OPEN) {
        this.socket.send(JSON.stringify({ event: 'ping', timestamp: Date.now() }));
      }
    }, 5000);
  }

  stopPingLoop() {
    if (this.pingInterval) {
      clearInterval(this.pingInterval);
      this.pingInterval = null;
    }
  }

  startRateCalculation() {
    this.rateInterval = setInterval(() => {
      this.triggerStatsChange({ eventsPerSec: this.msgCount });
      this.msgCount = 0;
    }, 1000);
  }

  stopRateCalculation() {
    if (this.rateInterval) {
      clearInterval(this.rateInterval);
      this.rateInterval = null;
    }
  }

  addEventListener(event, callback) {
    if (!this.eventListeners.has(event)) {
      this.eventListeners.set(event, new Set());
    }
    this.eventListeners.get(event).add(callback);
  }

  removeEventListener(event, callback) {
    if (this.eventListeners.has(event)) {
      this.eventListeners.get(event).delete(callback);
    }
  }

  dispatchEvent(message) {
    // Standard event routing
    if (message.event) {
      if (message.event === 'pong') {
        const rtt = Date.now() - message.timestamp;
        this.triggerStatsChange({ latency: rtt });
        return;
      }

      const listeners = this.eventListeners.get(message.event);
      if (listeners) {
        listeners.forEach((cb) => cb(message));
      }

      if (this.onMetricsUpdateCallback) {
        this.onMetricsUpdateCallback(message);
      }
    } else if (message.type) {
      if (this.onMetricsUpdateCallback) {
        this.onMetricsUpdateCallback(message);
      }
    }
  }

  updateStatus(status) {
    if (this.onStatusChangeCallback) {
      this.onStatusChangeCallback(status);
    }
  }

  triggerStatsChange(stats) {
    if (this.onStatsChangeCallback) {
      this.onStatsChangeCallback(stats);
    }
  }

  disconnect() {
    if (this.socket) {
      this.socket.onclose = null; // Prevent reconnection trigger
      this.socket.close();
      this.socket = null;
    }
    this.stopPingLoop();
    this.stopRateCalculation();
    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout);
      this.reconnectTimeout = null;
    }
  }
}

export const wsClientInstance = new EnterpriseWebSocketClient();
