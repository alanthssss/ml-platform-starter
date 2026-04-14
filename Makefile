.PHONY: bootstrap deploy train test destroy smoke-test fmt lint build

CLUSTER_NAME ?= ml-platform
KUBECONFIG   ?= $(HOME)/.kube/config
NAMESPACE    ?= ml-platform

## ── Local cluster ────────────────────────────────────────────────────────────

bootstrap: ## Create local k3d cluster and install Argo CD
	k3d cluster create $(CLUSTER_NAME) --wait
	kubectl create namespace argocd --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
	kubectl wait --for=condition=available --timeout=120s deployment/argocd-server -n argocd
	@echo "Argo CD is ready. Run 'make deploy' to apply the app-of-apps."

destroy: ## Delete local k3d cluster
	k3d cluster delete $(CLUSTER_NAME)

## ── GitOps deploy ────────────────────────────────────────────────────────────

deploy: ## Apply Argo CD app-of-apps
	kubectl apply -f deploy/argocd/app-of-apps.yaml

## ── ML pipeline ──────────────────────────────────────────────────────────────

train: ## Train + export + register Iris model
	cd ml/iris && pip install -q -r requirements.txt
	cd ml/iris && python train.py
	cd ml/iris && python export.py
	cd ml/iris && python register.py

## ── Go service ───────────────────────────────────────────────────────────────

build: ## Build api-gateway-go binary
	$(MAKE) -C services/api-gateway-go build

test: ## Run Go unit tests
	$(MAKE) -C services/api-gateway-go test

lint: ## Lint Go code
	$(MAKE) -C services/api-gateway-go lint

fmt: ## Format Go code
	$(MAKE) -C services/api-gateway-go fmt

## ── Smoke tests ──────────────────────────────────────────────────────────────

smoke-test: ## Run smoke tests against local deployment
	@echo "Checking api-gateway-go /healthz..."
	kubectl run smoke --image=curlimages/curl --restart=Never --rm -i \
	  -- curl -sf http://api-gateway.$(NAMESPACE).svc.cluster.local:8080/healthz
	@echo "Smoke tests passed."

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
