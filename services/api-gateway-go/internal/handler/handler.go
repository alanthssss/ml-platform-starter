package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/alanthssss/ml-platform-starter/services/api-gateway-go/internal/metrics"
	"github.com/alanthssss/ml-platform-starter/services/api-gateway-go/internal/service"
	"go.uber.org/zap"
)

// Handler holds HTTP handlers for the API gateway.
type Handler struct {
	svc service.PredictService
	log *zap.Logger
}

// New creates a new Handler.
func New(svc service.PredictService, log *zap.Logger) *Handler {
	return &Handler{svc: svc, log: log}
}

// Healthz returns 200 OK for liveness checks.
func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// Readyz returns 200 OK when the service is ready.
func (h *Handler) Readyz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// Predict accepts a JSON body, calls Triton, and returns the result.
func (h *Handler) Predict(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	start := time.Now()
	metrics.RequestsTotal.WithLabelValues("predict").Inc()

	var req service.PredictRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Warn("bad request", zap.Error(err))
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	resp, err := h.svc.Predict(r.Context(), req)
	if err != nil {
		h.log.Error("prediction failed", zap.Error(err))
		http.Error(w, "prediction error", http.StatusInternalServerError)
		return
	}

	metrics.RequestLatency.WithLabelValues("predict").Observe(time.Since(start).Seconds())
	metrics.PredictionsTotal.Inc()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
