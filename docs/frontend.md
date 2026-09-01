# InfraPilot Enterprise Frontend Guide

## Overview

The InfraPilot Frontend is a modern React-based single-page application (SPA) that provides a real-time dashboard for monitoring infrastructure. It communicates with the backend via REST APIs and WebSockets.

## Technology Stack

- **Framework:** React 18
- **Build Tool:** Vite
- **Language:** JavaScript (ES6+)
- **Styling:** CSS3 with CSS Variables
- **State Management:** Context API + Hooks
- **Routing:** React Router
- **Charts:** ApexCharts / Chart.js
- **Terminal:** xterm.js
- **Icons:** React Icons / Heroicons
- **HTTP Client:** Axios
- **WebSocket:** Native WebSocket API

## Architecture

### Directory Structure

```
frontend/
├── public/
│   ├── index.html
│   └── favicon.ico
├── src/
│   ├── app/
│   │   └── App.jsx              # Main app component
│   │   └── Routes.jsx           # Route definitions
│   ├── features/
│   │   ├── overview/
│   │   │   └── EnterpriseDashboard.jsx
│   │   ├── machines/
│   │   │   ├── MachineList.jsx
│   │   │   ├── MachineDetails.jsx
│   │   │   └── MachineMetrics.jsx
│   │   ├── terminal/
│   │   │   └── Terminal.jsx
│   │   ├── files/
│   │   │   └── FileManager.jsx
│   │   ├── logs/
│   │   │   └── LogViewer.jsx
│   │   ├── alerts/
│   │   │   ├── AlertList.jsx
│   │   │   └── AlertForm.jsx
│   │   ├── reports/
│   │   │   ├── ReportsPage.jsx
│   │   │   └── ReportGenerator.jsx
│   │   └── settings/
│   │       └── Settings.jsx
│   ├── components/
│   │   ├── Layout.jsx
│   │   ├── Sidebar.jsx
│   │   ├── Header.jsx
│   │   ├── MetricCard.jsx
│   │   └── ChartWidget.jsx
│   ├── services/
│   │   ├── api.js               # Axios instance
│   │   ├── websocket.js         # WebSocket manager
│   │   ├── auth.js              # Authentication
│   │   └── metrics.js           # Metrics helpers
│   ├── hooks/
│   │   ├── useAuth.js
│   │   ├── useWebSocket.js
│   │   └── useMetrics.js
│   ├── utils/
│   │   ├── formatters.js
│   │   └── validators.js
│   └── styles/
│       └── main.css
├── Dockerfile
├── nginx.conf
└── package.json
```

### Component Hierarchy

```
App
├── AuthProvider
├── Router
│   ├── Layout
│   │   ├── Sidebar
│   │   ├── Header
│   │   └── Content
│   │       ├── Dashboard
│   │       ├── Machines
│   │       ├── Terminal
│   │       ├── Files
│   │       ├── Logs
│   │       ├── Alerts
│   │       ├── Reports
│   │       └── Settings
```

## Key Components

### Dashboard

**File:** `src/features/overview/EnterpriseDashboard.jsx`

Real-time overview of all monitored infrastructure:

```javascript
import { useState, useEffect } from 'react';
import axios from 'axios';
import MetricCard from '../../components/MetricCard';
import ChartWidget from '../../components/ChartWidget';

export default function EnterpriseDashboard() {
  const [machines, setMachines] = useState([]);
  const [metrics, setMetrics] = useState({});
  const [alerts, setAlerts] = useState([]);

  useEffect(() => {
    // Fetch initial data
    loadDashboardData();

    // Subscribe to real-time updates
    const ws = new WebSocket(`${API_URL}/ws`);
    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.type === 'metric') {
        setMetrics(prev => ({ ...prev, ...data.payload }));
      }
    };

    return () => ws.close();
  }, []);

  return (
    <div className="dashboard">
      <div className="metrics-grid">
        <MetricCard title="Active Machines" value={machines.length} />
        <MetricCard title="Active Alerts" value={alerts.length} alert />
        <MetricCard title="Avg CPU" value={metrics.avg_cpu} unit="%" />
        <MetricCard title="Avg Memory" value={metrics.avg_memory} unit="%" />
      </div>

      <div className="charts-section">
        <ChartWidget
          title="CPU Usage"
          data={metrics.cpu_history}
          type="area"
        />
        <ChartWidget
          title="Memory Usage"
          data={metrics.memory_history}
          type="area"
        />
      </div>
    </div>
  );
}
```

### Terminal Component

**File:** `src/features/terminal/Terminal.jsx`

xterm.js integration for remote shell access:

```javascript
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import { useWebSocket } from '../../hooks/useWebSocket';

export default function Terminal({ machineId, sessionId }) {
  const terminalRef = useRef(null);
  const { send, lastMessage } = useWebSocket(`/ws/terminal/${sessionId}`);

  useEffect(() => {
    const term = new Terminal({
      cursorBlink: true,
      theme: { background: '#1e1e1e' },
      fontSize: 14
    });

    const fitAddon = new FitAddon();
    term.loadAddon(fitAddon);
    term.open(terminalRef.current);
    fitAddon.fit();

    term.onData((data) => {
      send({ type: 'input', data });
    });

    term.onResize(({ cols, rows }) => {
      send({ type: 'resize', cols, rows });
    });

    // Handle incoming output
    if (lastMessage?.type === 'output') {
      term.write(lastMessage.data);
    }

    return () => term.dispose();
  }, []);

  return <div ref={terminalRef} className="terminal-container" />;
}
```

### Real-time WebSocket Hook

**File:** `src/hooks/useWebSocket.js`

```javascript
import { useState, useEffect, useRef } from 'react';

export function useWebSocket(url) {
  const [lastMessage, setLastMessage] = useState(null);
  const ws = useRef(null);

  useEffect(() => {
    const token = localStorage.getItem('token');
    ws.current = new WebSocket(`${API_URL}${url}?token=${token}`);

    ws.current.onmessage = (event) => {
      const data = JSON.parse(event.data);
      setLastMessage(data);
    };

    return () => ws.current?.close();
  }, [url]);

  const send = (data) => {
    if (ws.current?.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify(data));
    }
  };

  return { lastMessage, send };
}
```

## State Management

### Authentication Context

**File:** `src/app/AuthProvider.jsx`

```javascript
import { createContext, useContext, useState } from 'react';

const AuthContext = createContext();

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [token, setToken] = useState(localStorage.getItem('token'));

  const login = async (email, password) => {
    const response = await axios.post('/auth/login', { email, password });
    const { access_token, refresh_token, user } = response.data;
    
    setToken(access_token);
    setUser(user);
    localStorage.setItem('token', access_token);
    localStorage.setItem('refresh_token', refresh_token);
  };

  const logout = () => {
    setToken(null);
    setUser(null);
    localStorage.removeItem('token');
  };

  return (
    <AuthContext.Provider value={{ user, token, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => useContext(AuthContext);
```

## API Integration

### Axios Instance

**File:** `src/services/api.js`

```javascript
import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL;

export const api = axios.create({
  baseURL: API_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Interceptor for auth
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Interceptor for token refresh
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.status === 401) {
      // Try to refresh token
      const refresh_token = localStorage.getItem('refresh_token');
      try {
        const response = await axios.post(`${API_URL}/auth/refresh`, null, {
          headers: { Authorization: `Bearer ${refresh_token}` }
        });
        localStorage.setItem('token', response.data.access_token);
        error.config.headers.Authorization = `Bearer ${response.data.access_token}`;
        return api(error.config);
      } catch (refreshError) {
        localStorage.removeItem('token');
        window.location.href = '/login';
      }
    }
    return Promise.reject(error);
  }
);
```

## UI Components

### Metric Card

```javascript
export default function MetricCard({ title, value, unit, trend, alert }) {
  const trendIcon = trend > 0 ? '↑' : '↓';
  const trendClass = trend > 0 ? 'trend-up' : 'trend-down';

  return (
    <div className={`metric-card ${alert ? 'alert' : ''}`}>
      <h3>{title}</h3>
      <div className="metric-value">
        {value}{unit && <span className="unit">{unit}</span>}
      </div>
      {trend && (
        <div className={`trend ${trendClass}`}>
          {trendIcon} {Math.abs(trend)}%
        </div>
      )}
    </div>
  );
}
```

### Chart Widget

```javascript
export default function ChartWidget({ title, data, type }) {
  const options = {
    chart: {
      type,
      height: 300,
      animations: { enabled: false },
    },
    series: [{ name: title, data }],
    xaxis: { type: 'datetime' },
    stroke: { curve: 'smooth' },
  };

  return (
    <div className="chart-widget">
      <h4>{title}</h4>
      <ReactApexChart options={options} series={options.series} type={type} />
    </div>
  );
}
```

## Build & Deployment

### Development

```bash
# Install dependencies
npm install

# Start dev server
npm run dev

# Lint
npm run lint

# Type check (if using TypeScript)
npm run typecheck
```

### Production Build

```bash
# Build for production
npm run build

# Preview production build
npm run preview
```

Build output goes to `dist/` directory.

### Docker Deployment

See `Dockerfile.frontend` and `nginx.conf`.

## Performance Optimization

### Code Splitting

```javascript
import { lazy, Suspense } from 'react';

const Terminal = lazy(() => import('./features/terminal/Terminal'));
const Reports = lazy(() => import('./features/reports/ReportsPage'));

function App() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <Routes>
        <Route path="/terminal" element={<Terminal />} />
        <Route path="/reports" element={<Reports />} />
      </Routes>
    </Suspense>
  );
}
```

### Memoization

```javascript
import { memo, useMemo } from 'react';

const MachineList = memo(({ machines, onSelect }) => {
  const sortedMachines = useMemo(() => {
    return machines.sort((a, b) => a.name.localeCompare(b.name));
  }, [machines]);

  return (
    <ul>
      {sortedMachines.map(machine => (
        <li key={machine.id} onClick={() => onSelect(machine)}>
          {machine.name}
        </li>
      ))}
    </ul>
  );
});
```

### Virtual Scrolling

For long lists:

```javascript
import { FixedSizeList } from 'react-window';

const MachineList = ({ machines }) => (
  <FixedSizeList
    height={600}
    itemCount={machines.length}
    itemSize={50}
    width="100%"
  >
    {({ index, style }) => (
      <div style={style}>
        {machines[index].name}
      </div>
    )}
  </FixedSizeList>
);
```

## Testing

### Unit Tests

```javascript
import { render, screen, fireEvent } from '@testing-library/react';
import MetricCard from '../components/MetricCard';

test('renders metric card', () => {
  render(<MetricCard title="CPU" value={45.2} unit="%" />);
  expect(screen.getByText('CPU')).toBeInTheDocument();
  expect(screen.getByText('45.2%')).toBeInTheDocument();
});
```

### E2E Tests

```javascript
import { test, expect } from '@playwright/test';

test('login flow', async ({ page }) => {
  await page.goto('/login');
  await page.fill('[name=email]', 'admin@infrapilot.io');
  await page.fill('[name=password]', 'admin123');
  await page.click('button[type=submit]');
  await expect(page).toHaveURL('/dashboard');
});
```

## Browser Support

- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

## Accessibility

- WCAG 2.1 Level AA compliance
- Keyboard navigation support
- Screen reader compatibility
- High contrast mode support
- Focus indicators

## Troubleshooting

### WebSocket connection fails

```javascript
// Check configuration
console.log('WS URL:', import.meta.env.VITE_WS_URL);

// Verify token
const token = localStorage.getItem('token');
```

### Charts not rendering

```bash
# Check ApexCharts is installed
npm list apexcharts

# Reinstall if needed
npm install apexcharts react-apexcharts
```

### Build fails

```bash
# Clear cache and reinstall
rm -rf node_modules package-lock.json
npm install
npm run build