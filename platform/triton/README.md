# Triton Inference Server

Kubernetes manifests for deploying NVIDIA Triton Inference Server in the `ml-platform` namespace.

## Files

| File | Description |
|------|-------------|
| namespace.yaml | Creates the ml-platform namespace |
| pvc.yaml | 1 Gi PVC for the model repository |
| deployment.yaml | Triton server deployment with probes and resources |
| service.yaml | Service exposing HTTP (8000), gRPC (8001), metrics (8002) |
| model-repo/iris_onnx/config.pbtxt | Triton model config for the Iris ONNX model |

## Ports

| Port | Protocol | Purpose |
|------|----------|---------|
| 8000 | HTTP | KServe v2 inference API |
| 8001 | gRPC | gRPC inference API |
| 8002 | HTTP | Prometheus metrics |

## Model Repository Layout

```
/models/
└── iris_onnx/
    ├── config.pbtxt
    └── 1/
        └── model.onnx
```

## Applying

```bash
kubectl apply -f platform/triton/
```

## Copying the model

```bash
kubectl cp ml/iris/artifacts/iris_onnx/ ml-platform/$(kubectl get pod -n ml-platform -l app=triton -o jsonpath='{.items[0].metadata.name}'):/models/iris_onnx
```
