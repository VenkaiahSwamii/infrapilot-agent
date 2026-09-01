import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 50 },  // Ramp-up to 50 concurrent users
    { duration: '1m', target: 100 },  // Peak load: 100 concurrent users / 1,000 agents
    { duration: '30s', target: 0 },   // Ramp-down
  ],
  thresholds: {
    http_req_duration: ['p(95)<100'], // 95% of API requests must complete in <100ms
  },
};

export default function () {
  const res = http.get('http://localhost:8080/api/v1/platform/health');
  check(res, {
    'status is 200': (r) => r.status === 200,
    'health status is HEALTHY': (r) => JSON.parse(r.body).status === 'HEALTHY',
  });
  sleep(1);
}
