import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 10 },    // Warm-up 10 users
    { duration: '1m',  target: 100 },   // Ramp-up 100 users
    { duration: '2m',  target: 1000 },  // Scale 1,000 users
    { duration: '2m',  target: 10000 }, // Peak Stress 10,000 users
    { duration: '1m',  target: 0 },     // Ramp-down
  ],
  thresholds: {
    http_req_duration: ['p(95)<100'], // 95% of requests must complete < 100ms
    http_req_failed: ['rate<0.01'],    // < 1% error rate
  },
};

const BASE_URL = __ENV.TARGET_URL || 'http://localhost:8080/api/v1';
const JWT_TOKEN = __ENV.TOKEN || 'Bearer sample_jwt_token';

export default function () {
  const params = {
    headers: {
      'Authorization': JWT_TOKEN,
      'Content-Type': 'application/json',
      'X-Organization-ID': 'default',
    },
  };

  // 1. Health check & Readiness
  const resReadiness = http.get(`${BASE_URL}/hardening/readiness`, params);
  check(resReadiness, {
    'readiness status is 200': (r) => r.status === 200,
    'latency < 100ms': (r) => r.timings.duration < 100,
  });

  // 2. Telemetry metrics query
  const resMetrics = http.get(`${BASE_URL}/aiops/dashboard`, params);
  check(resMetrics, {
    'aiops dashboard status is 200': (r) => r.status === 200,
  });

  sleep(1);
}
