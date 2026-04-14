# ---- build stage ----
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /api-gateway ./cmd/server

# ---- runtime stage ----
FROM gcr.io/distroless/static:nonroot
COPY --from=builder /api-gateway /api-gateway
EXPOSE 8080
ENTRYPOINT ["/api-gateway"]
