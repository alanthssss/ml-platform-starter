# Quick Start

## Prerequisites

| Tool | Minimum version |
|------|----------------|
| Docker | 24+ |
| kubectl | 1.28+ |
| k3d | 5.6+ |
| Go | 1.21+ |
| Python | 3.10+ |
| make | any |

Install k3d:
```bash
curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash
```

## 1. Bootstrap the cluster

```bash
make bootstrap
```

This creates a k3d cluster named `ml-platform` and installs Argo CD.

## 2. Train and export the model

```bash
make train
```

This runs the Iris training pipeline and places the ONNX model in
`ml/iris/artifacts/iris_onnx/1/model.onnx`.

## 3. Deploy platform via Argo CD

```bash
make deploy
```

This applies the Argo CD app-of-apps which syncs Triton, api-gateway, Prometheus,
and Grafana.

## 4. Smoke test

```bash
make smoke-test
```

## 5. Manual prediction

Port-forward the api-gateway:
```bash
kubectl port-forward svc/api-gateway -n ml-platform 8080:8080
```

Then:
```bash
curl -X POST http://localhost:8080/predict \
  -H 'Content-Type: application/json' \
  -d '{"inputs": [5.1, 3.5, 1.4, 0.2]}'
```

Expected response:
```json
{"model":"iris_onnx","prediction":0,"label":"setosa"}
```

## 6. Access Grafana

```bash
kubectl port-forward svc/grafana -n monitoring 3000:3000
```

Open http://localhost:3000 (admin/admin).

## 7. Tear down

```bash
make destroy
```
