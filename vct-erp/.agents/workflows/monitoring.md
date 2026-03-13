---
description: Setup monitoring, logging, and alerting for VCT Platform
---

# /monitoring - Setup Monitoring & Alerting

## Architecture

```
Application → Prometheus (metrics) → Grafana (dashboards)
Application → Zerolog (logs) → Loki/ELK (log aggregation)
Grafana → Alertmanager → Slack/Email/PagerDuty
Neon Dashboard → Compute hours, storage, query performance
Supabase Dashboard → Auth, Realtime, Storage, API usage
```

## Steps

### Step 1: Add health check endpoints to backend
```go
// GET /health - basic liveness
// GET /ready - readiness (DB, Redis connected)
// GET /metrics - Prometheus metrics
```

### Step 2: Docker Compose monitoring stack
```yaml
# deploy/docker/docker-compose.monitoring.yml
services:
  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3001:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana_data:/var/lib/grafana

  loki:
    image: grafana/loki:latest
    ports:
      - "3100:3100"

  alertmanager:
    image: prom/alertmanager:latest
    ports:
      - "9093:9093"
```

### Step 3: Prometheus configuration
```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'vct-api'
    scrape_interval: 15s
    static_configs:
      - targets: ['api:8080']
    metrics_path: /metrics

  - job_name: 'postgres'
    static_configs:
      - targets: ['postgres-exporter:9187']

  - job_name: 'redis'
    static_configs:
      - targets: ['redis-exporter:9121']
```

### Step 4: Key metrics to track

| Category | Metric | Alert Threshold |
|----------|--------|----------------|
| **HTTP** | Request rate (req/s) | No alert |
| **HTTP** | Error rate (5xx) | > 1% for 5 min |
| **HTTP** | Response time P95 | > 500ms for 5 min |
| **HTTP** | Response time P99 | > 1s for 5 min |
| **DB** | Active connections | > 80% pool |
| **DB** | Query duration P95 | > 200ms |
| **Neon** | Compute hours | > budget limit |
| **Neon** | Storage | > 80% of plan |
| **Supabase** | API requests | > rate limit |
| **Supabase** | Auth requests | > 1000/hour |
| **Redis** | Memory usage | > 80% |
| **Redis** | Hit rate | < 90% |
| **System** | CPU usage | > 80% for 10 min |
| **System** | Memory usage | > 85% |
| **System** | Disk usage | > 90% |
| **App** | Goroutine count | > 10,000 |

### Step 5: Alert rules
```yaml
# alert_rules.yml
groups:
  - name: vct-platform
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.01
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected (> 1%)"

      - alert: SlowResponses
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 0.5
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "P95 response time > 500ms"

      - alert: DatabaseConnectionsHigh
        expr: pg_stat_activity_count > 160  # 80% of 200
        for: 2m
        labels:
          severity: warning
```

### Step 6: Grafana dashboards
Create dashboards for:
1. **API Overview**: Request rate, error rate, latency histogram
2. **Database (Neon)**: Compute hours, connections, query duration, cache hit ratio
3. **Infrastructure**: CPU, memory, disk, network
4. **Business Metrics**: Active users, registrations, tournament activity
5. **Supabase**: Auth sessions, Realtime connections, Storage usage

### Step 7: Start monitoring stack
```bash
docker compose -f deploy/docker/docker-compose.monitoring.yml up -d
echo "Prometheus:     http://localhost:9090"
echo "Grafana:        http://localhost:3001 (admin/admin)"
echo "Alertmanager:   http://localhost:9093"
echo "Neon Dashboard: https://console.neon.tech"
echo "Supabase:       https://supabase.com/dashboard"
```
