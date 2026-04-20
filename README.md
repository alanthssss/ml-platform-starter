# ml-platform-starter

A minimal but complete MLOps / AI service delivery monorepo for learning and as a reusable framework starter.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                         ml-platform-starter                         │
│                                                                     │
│  ┌──────────┐   ┌──────────────┐   ┌──────────────────────────┐   │
│  │  Airflow  │──▶│   MLflow     │──▶│  Triton Inference Server  │   │
│  │  (DAG)   │   │  (registry)  │   │  (model serving)          │   │
│  └──────────┘   └──────────────┘   └────────────┬─────────────┘   │
│                                                   │                 │
│  ┌──────────────────────────────────────────────▼─────────────┐   │
│  │              api-gateway-go  (/predict, /healthz, /metrics) │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
│  ┌──────────┐   ┌──────────────┐                                   │
│  │Prometheus│──▶│   Grafana    │                                   │
│  │(scrape)  │   │(dashboards)  │                                   │
│  └──────────┘   └──────────────┘                                   │
│                                                                     │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │                  Argo CD  (GitOps)                            │  │
│  └──────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────┘
```

## Stack

| Component | Purpose |
|-----------|---------|
| Terraform | Cloud infrastructure provisioning |
| Ansible | VM / node configuration |
| Argo CD | GitOps continuous delivery |
| Prometheus | Metrics collection |
| Grafana | Dashboards and alerting |
| Python | ML training, export, and registration |
| Go | API gateway service |
| MLflow | Experiment tracking and model registry |
| Airflow | Workflow orchestration |
| Triton | High-performance model inference |

## Quick Start

```bash
# 1. Prerequisites: Docker, kubectl, k3d (or kind), Go >=1.21, Python >=3.10

# 2. Bootstrap local cluster
make bootstrap

# 3. Train and export the Iris example model
make train

# 4. Deploy everything via Argo CD
make deploy

# 5. Run smoke tests
make smoke-test

# 6. Tear down
make destroy
```

## Repository Layout

```
ml-platform-starter/
├── docs/                  # Architecture and quickstart docs
├── infra/                 # Terraform modules (cloud infra)
├── ansible/               # Ansible playbooks (node config)
├── platform/              # Platform services (Triton, monitoring, Airflow)
│   ├── triton/            # Triton Inference Server K8s manifests
│   ├── monitoring/        # Prometheus + Grafana
│   └── airflow/           # Airflow deployment
├── services/              # Application services
│   └── api-gateway-go/    # Go API gateway
├── ml/                    # ML code
│   └── iris/              # Iris classification example
├── deploy/                # Argo CD / GitOps
│   └── argocd/            # App-of-apps and child applications
├── examples/              # End-to-end usage examples
├── templates/             # Reusable scaffolding templates
├── Makefile
├── .gitignore
├── .editorconfig
└── README.md
```

## Roadmap

- [x] Repository structure and docs
- [x] Go API gateway with /predict, /healthz, /readyz, /metrics
- [x] Python Iris classifier: train, export to ONNX, register in MLflow
- [x] Triton Inference Server K8s manifests
- [x] Argo CD app-of-apps
- [x] Prometheus + Grafana monitoring stack
- [x] Airflow DAG for end-to-end ML pipeline
- [ ] Terraform module for cloud cluster (GKE/EKS)
- [ ] Ansible playbook for bare-metal node setup
- [ ] Multi-model support in Triton
- [ ] Auth / JWT middleware in API gateway
- [ ] Canary rollout via Argo Rollouts

## License

MIT