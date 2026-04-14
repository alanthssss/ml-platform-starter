# Architecture

## Overview

ml-platform-starter is a local-first, GitOps-driven ML serving platform.

```
Client
  │
  ▼
api-gateway-go  (Go, port 8080)
  │  /predict
  ▼
Triton Inference Server  (HTTP 8000 / gRPC 8001 / metrics 8002)
  │  model: iris_onnx
  ▼
Model Repository  (PVC or local volume)
```

## Components

### api-gateway-go
- Language: Go
- Endpoints: `/healthz`, `/readyz`, `/metrics`, `/predict`
- Calls Triton via HTTP (REST)
- Exposes Prometheus metrics for request rate, latency, and prediction count

### Triton Inference Server
- NVIDIA Triton (CPU-only mode for local dev)
- Serves ONNX models from a model repository volume
- Metrics on port 8002 (Prometheus-compatible)

### MLflow
- Tracks experiment runs, parameters, metrics, and artifact locations
- Used by `ml/iris/register.py`

### Airflow
- Orchestrates: train → export → register → copy to Triton repo
- DAG: `ml_pipeline_iris`

### Prometheus + Grafana
- Prometheus scrapes api-gateway-go and Triton
- Grafana shows request rate, latency, and prediction count

### Argo CD
- App-of-apps pattern
- Manages: triton, api-gateway, prometheus, grafana
- Auto-sync in local dev

## Data Flow

```
Airflow DAG
  ├─ train.py        → ml/iris/artifacts/model.pkl
  ├─ export.py       → ml/iris/artifacts/iris_onnx/1/model.onnx
  ├─ register.py     → MLflow (experiment + artifact)
  └─ copy to Triton repo (PVC)

Client POST /predict
  → api-gateway-go
  → Triton HTTP /v2/models/iris_onnx/infer
  → JSON response
```
