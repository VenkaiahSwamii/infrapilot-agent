# InfraPilot Enterprise Troubleshooting Guide

This guide covers common issues and their solutions.

## Table of Contents

- [Quick Diagnostics](#quick-diagnostics)
- [Installation Issues](#installation-issues)
- [Backend Problems](#backend-problems)
- [Frontend Issues](#frontend-issues)
- [Agent Issues](#agent-issues)
- [Database Issues](#database-issues)
- [Performance Problems](#performance-problems)
- [Network Issues](#network-issues)

---

## Quick Diagnostics

### Health Check Script

```bash
#!/bin/bash
echo "=== InfraPilot Health Check ==="

# Check backend
echo -n "Backend: "
curl -s http://localhost:8080/healthz | jq -r '.status'

# Check database
echo -n "PostgreSQL: "
docker exec infrapilot-postgres pg_isready -U postgres

# Check Redis
echo -n "Redis: "
docker exec infrapilot-redis redis-cli ping

# Check pods (Kubernetes)
echo "=== Kubernetes Pods ==="
kubectl get pods -n infrapilot

# Check logs
echo "=== Recent Backend Logs ==="
docker logs --tail 50 infrapilot-backend
```

---

## Installation Issues

### Docker Compose fails to start

**Symptom:** `docker compose up -d` hangs or fails

**Solutions:**

1. **Check Docker version**
   ```bash
   docker --version  # Should be 24.0+
   docker compose version  # Should be 2.20+
   ```

2. **Check available resources**
   ```bash
   docker system info
   # Ensure at least 4GB RAM available
   ```

3. **Check port conflicts**
   ```bash
   netstat -tulpn | grep -E '8080|5432|6379'
   # Stop conflicting services
   ```

4. **Clean Docker cache**
   ```bash
   docker system prune -a
   ```

### Port already in use

**Symptom:** `Bind for 0.0.0.0:8080 failed: port is already allocated`

**Solution:**
```bash
# Find process using port
sudo lsof -i :8080
sudo kill -9 <PID>

# Or modify docker-compose.yml to use different ports
ports:
  - "8081:8080"  # Change host port
```

---

## Backend Problems

### Backend won't start

**Symptom:** Backend container exits immediately

**Diagnose:**
```bash
docker logs infrapilot-backend
```

**Common causes:**

1. **Database not ready**
   ```
   Error: connection refused
   ```
   **Solution:** Wait for PostgreSQL health check or increase `depends_on` timeout.

2. **Invalid environment variables**
   ```
   Error: JWT_SECRET is required
   ```
   **Solution:** Check `backend/.env` file exists and has required values.

3. **Missing migrations**
   ```
   Error: relation "users" does not exist
   ```
   **Solution:** Apply migrations:
   ```bash
   docker exec infrapilot-postgres psql -U postgres -f /docker-entrypoint-initdb.d/*.sql
   ```

### High memory usage

**Symptom:** Backend using > 2GB RAM

**Solutions:**

1. **Check for memory leaks**
   ```bash
   docker stats infrapilot-backend
   ```

2. **Enable pprof profiling**
   ```bash
   curl http://localhost:8080/debug/pprof/heap
   ```

3. **Reduce worker count**
   ```env
   # backend/.env
   WORKER_PROCESSES=2
   ```

### Slow API responses

**Symptom:** API latency > 500ms

**Diagnose:**
```bash
# Check database query performance
docker exec infrapilot-postgres psql -U postgres -c "SELECT * FROM pg_stat_activity;"

# Check Redis hit rate
docker exec infrapilot-redis redis-cli info stats | grep keyspace_hits
```

**Solutions:**

1. **Add database indexes**
   ```sql
   CREATE INDEX idx_machines_status ON machines(status);
   CREATE INDEX idx_metrics_timestamp ON metrics(timestamp);
   ```

2. **Increase Redis memory**
   ```yaml
   # docker-compose.yml
   redis:
     command: redis-server --maxmemory 2gb --maxmemory-policy allkeys-lru
   ```

---

## Frontend Issues

### Blank page / White screen

**Symptom:** Frontend loads but shows nothing

**Diagnose:**
```bash
# Check browser console for errors
# Check nginx logs
docker logs infrapilot-frontend
```

**Solutions:**

1. **Check API URL configuration**
   ```bash
   # Verify nginx.conf
   proxy_pass http://backend:8080/api;
   ```

2. **Rebuild frontend**
   ```bash
   docker compose build frontend
   docker compose up -d frontend
   ```

3. **Clear browser cache**
   - Open DevTools (F12)
   - Right-click refresh
   - "Empty Cache and Hard Reload"

### Static assets not loading

**Symptom:** CSS/JS files return 404

**Solution:**
```bash
# Verify build output exists
docker exec infrapilot-frontend ls -la /usr/share/nginx/html

# Rebuild if missing
docker compose build --no-cache frontend
```

---

## Agent Issues

### Agent won't connect

**Symptom:** Agent shows disconnected in dashboard

**Diagnose:**
```bash
# Check agent logs
docker logs infrapilot-agent

# Test WebSocket connection
curl -i -N -H "Connection: Upgrade" -H "Upgrade: websocket" \
  -H "Sec-WebSocket-Key: test" -H "Sec-WebSocket-Version: 13" \
  http://localhost:8080/ws
```

**Solutions:**

1. **Check backend URL**
   ```yaml
   # agent/config.yaml
   backend:
     url: "http://backend:8080"
   ```

2. **Verify firewall allows WebSocket**
   ```bash
   sudo ufw allow 8080/tcp
   ```

3. **Check agent authentication**
   ```bash
   # View agent logs for auth errors
   docker logs infrapilot-agent | grep -i auth
   ```

### Agent high CPU usage

**Symptom:** Agent consuming > 50% CPU

**Solutions:**

1. **Reduce collection frequency**
   ```yaml
   # agent/config.yaml
   collection:
     interval: 30s  # Increase from default 10s
   ```

2. **Disable unused plugins**
   ```yaml
   plugins:
     kubernetes:
       enabled: false  # Disable if not needed
   ```

---

## Database Issues

### Database connection refused

**Symptom:** `connection refused` or `timeout`

**Solutions:**

1. **Check PostgreSQL is running**
   ```bash
   docker ps | grep postgres
   docker logs infrapilot-postgres
   ```

2. **Verify credentials**
   ```bash
   docker exec -it infrapilot-postgres psql -U postgres
   ```

3. **Check max connections**
   ```sql
   SHOW max_connections;
   SELECT count(*) FROM pg_stat_activity;
   ```

### Database full

**Symptom:** `No space left on device`

**Solutions:**

1. **Clean old data**
   ```sql
   -- Delete old metrics
   DELETE FROM metrics WHERE timestamp < NOW() - INTERVAL '30 days';
   VACUUM FULL metrics;
   ```

2. **Increase disk space**
   ```bash
   # Extend volume
   docker volume ls
   docker volume inspect infrapilot_postgres_data
   ```

### Slow queries

**Symptom:** Queries taking > 1 second

**Diagnose:**
```sql
-- Find slow queries
SELECT query, mean_exec_time, calls
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;
```

**Solutions:**
```sql
-- Add indexes
CREATE INDEX CONCURRENTLY idx_metrics_machine_id ON metrics(machine_id);
CREATE INDEX CONCURRENTLY idx_alerts_created_at ON alerts(created_at);

-- Update statistics
ANALYZE;
```

---

## Performance Problems

### High CPU usage

**Symptom:** Server CPU > 90%

**Diagnose:**
```bash
# Check which process is using CPU
top -p $(pgrep -f infrapilot)

# Check database queries
docker exec infrapilot-postgres psql -U postgres -c \
  "SELECT pid, usename, query, state FROM pg_stat_activity WHERE state != 'idle';"
```

**Solutions:**

1. **Scale backend horizontally**
   ```bash
   docker compose up -d --scale backend=3
   ```

2. **Optimize queries**
   - Add indexes (see Database section)
   - Use connection pooling (PgBouncer)

3. **Enable caching**
   ```env
   REDIS_CACHE_TTL=300
   ```

### Memory leaks

**Symptom:** Memory usage grows continuously

**Diagnose:**
```bash
# Monitor memory over time
watch -n 1 'docker stats infrapilot-backend --no-stream'

# Generate heap profile
curl http://localhost:8080/debug/pprof/heap > heap.prof
go tool pprof -http=:8081 heap.prof
```

---

## Network Issues

### WebSocket disconnections

**Symptom:** WebSocket connections drop frequently

**Diagnose:**
```bash
# Check WebSocket logs
docker logs infrapilot-backend | grep -i websocket

# Test with wscat
npm install -g wscat
wscat -c ws://localhost:8080/ws
```

**Solutions:**

1. **Increase timeouts**
   ```yaml
   # nginx.conf
   proxy_read_timeout 3600s;
   proxy_send_timeout 3600s;
   ```

2. **Check load balancer limits**
   ```yaml
   # For AWS ALB
   idle_timeout.timeout_seconds: 3600
   ```

### SSL/TLS errors

**Symptom:** `SSL: CERTIFICATE_VERIFY_FAILED`

**Solutions:**

1. **Verify certificate**
   ```bash
   openssl s_client -connect infrapilot.io:443
   ```

2. **Check certificate expiry**
   ```bash
   echo | openssl s_client -connect infrapilot.io:443 2>/dev/null | \
     openssl x509 -noout -dates
   ```

---

## Common Error Messages

### "JWT token expired"
**Solution:** Refresh token or login again

### "Rate limit exceeded"
**Solution:** Wait 60 seconds or increase rate limits in config

### "Database connection pool exhausted"
**Solution:** Increase max connections or use connection pooling

### "Queue is full"
**Solution:** Increase queue size or scale workers

### "Agent certificate invalid"
**Solution:** Re-enroll agent with new certificate

---

## Getting Help

### Collect Debug Information

```bash
#!/bin/bash
# debug.sh

echo "=== InfraPilot Debug Info ===" > debug.txt
date >> debug.txt

echo -e "\n=== Docker Compose Config ===" >> debug.txt
docker compose config >> debug.txt

echo -e "\n=== Running Containers ===" >> debug.txt
docker ps >> debug.txt

echo -e "\n=== Backend Logs ===" >> debug.txt
docker logs --tail 100 infrapilot-backend >> debug.txt 2>&1

echo -e "\n=== PostgreSQL Logs ===" >> debug.txt
docker logs --tail 100 infrapilot-postgres >> debug.txt 2>&1

echo -e "\n=== Redis Info ===" >> debug.txt
docker exec infrapilot-redis redis-cli INFO all >> debug.txt 2>&1

echo "Debug info saved to debug.txt"
```

### Support Channels

- **Documentation:** https://docs.infrapilot.io
- **GitHub Issues:** https://github.com/yourusername/infrapilot-enterprise/issues
- **Discord:** https://discord.gg/infrapilot
- **Email:** support@infrapilot.io

### Reporting Issues

Include:
1. InfraPilot version
2. Deployment method (Docker/Kubernetes)
3. OS and version
4. Full error message
5. Steps to reproduce
6. Debug log (from script above)