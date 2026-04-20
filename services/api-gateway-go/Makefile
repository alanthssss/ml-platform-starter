.PHONY: build test lint fmt docker-build docker-push

IMAGE ?= api-gateway-go:local
BINARY = bin/api-gateway

build:
	go build -o $(BINARY) ./cmd/server

test:
	go test ./... -v -race

lint:
	golangci-lint run ./...

fmt:
	gofmt -s -w .

docker-build:
	docker build -t $(IMAGE) .

docker-push:
	docker push $(IMAGE)
