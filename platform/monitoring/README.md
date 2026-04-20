# Monitoring Stack

Prometheus and Grafana deployment for ml-platform-starter.

## Components

- **Prometheus** - scrapes api-gateway-go and Triton metrics
- **Grafana** - dashboards provisioned from ConfigMaps

## Applying

```bash
kubectl create namespace monitoring
kubectl apply -f platform/monitoring/
```

## Accessing

```bash
# Prometheus
kubectl port-forward svc/prometheus -n monitoring 9090:9090
# → http://localhost:9090

# Grafana
kubectl port-forward svc/grafana -n monitoring 3000:3000
# → http://localhost:3000 (admin/admin)
```

## Dashboard

The `ml-platform` dashboard includes:
- API request rate
- API request latency (p99)
- Triton request latency (p99)
- Prediction count

## Metrics scraped

| Job | Target | Path |
|-----|--------|------|
| api-gateway | api-gateway.ml-platform:8080 | /metrics |
| triton | triton.ml-platform:8002 | /metrics |
