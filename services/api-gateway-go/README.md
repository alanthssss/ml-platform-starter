# api-gateway-go

Go HTTP API gateway that proxies prediction requests to Triton Inference Server.

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | /healthz | Liveness probe |
| GET | /readyz | Readiness probe |
| GET | /metrics | Prometheus metrics |
| POST | /predict | Forward prediction to Triton |

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| LISTEN_ADDR | :8080 | HTTP listen address |
| TRITON_URL | http://triton:8000 | Triton base URL |

## Example

```bash
curl -X POST http://localhost:8080/predict \
  -H 'Content-Type: application/json' \
  -d '{"inputs": [5.1, 3.5, 1.4, 0.2]}'
```

Response:
```json
{"model":"iris_onnx","prediction":0,"label":"setosa"}
```

## Running locally

```bash
# Start a mock Triton or point at a real one
export TRITON_URL=http://localhost:8000
make build
./bin/api-gateway
```

## Running tests

```bash
make test
```
