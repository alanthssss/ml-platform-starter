# Argo CD App-of-Apps

GitOps deployment using Argo CD app-of-apps pattern.

## Structure

```
deploy/argocd/
├── app-of-apps.yaml   # Root Application pointing to apps/
├── project.yaml       # AppProject for access control
├── apps/
│   ├── triton.yaml    # Triton Inference Server
│   ├── api-gateway.yaml
│   ├── prometheus.yaml
│   └── grafana.yaml
└── README.md
```

## Bootstrap Argo CD

```bash
# Install Argo CD
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Create project and app-of-apps
kubectl apply -f deploy/argocd/project.yaml
kubectl apply -f deploy/argocd/app-of-apps.yaml

# Get initial admin password
kubectl -n argocd get secret argocd-initial-admin-secret \
  -o jsonpath="{.data.password}" | base64 -d

# Port-forward UI
kubectl port-forward svc/argocd-server -n argocd 8443:443
```

Open https://localhost:8443 (admin / <password from above>).

## Sync Policy

All apps are configured for **automated sync** with prune and self-heal enabled.
This means any push to the repo will automatically be applied to the cluster.
