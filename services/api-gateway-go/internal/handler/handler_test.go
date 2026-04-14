package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alanthssss/ml-platform-starter/services/api-gateway-go/internal/handler"
	"github.com/alanthssss/ml-platform-starter/services/api-gateway-go/internal/service"
	"go.uber.org/zap"
)

// mockService is a test double for service.PredictService.
type mockService struct {
	resp service.PredictResponse
	err  error
}

func (m *mockService) Predict(_ context.Context, _ service.PredictRequest) (service.PredictResponse, error) {
	return m.resp, m.err
}

func TestHealthz(t *testing.T) {
	h := handler.New(&mockService{}, zap.NewNop())
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	h.Healthz(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestReadyz(t *testing.T) {
	h := handler.New(&mockService{}, zap.NewNop())
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rr := httptest.NewRecorder()
	h.Readyz(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestPredict_OK(t *testing.T) {
	mock := &mockService{
		resp: service.PredictResponse{Model: "iris_onnx", Prediction: 0, Label: "setosa"},
	}
	h := handler.New(mock, zap.NewNop())

	body, _ := json.Marshal(service.PredictRequest{Inputs: []float32{5.1, 3.5, 1.4, 0.2}})
	req := httptest.NewRequest(http.MethodPost, "/predict", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	h.Predict(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp service.PredictResponse
	json.NewDecoder(rr.Body).Decode(&resp)
	if resp.Label != "setosa" {
		t.Errorf("expected setosa, got %s", resp.Label)
	}
}

func TestPredict_BadMethod(t *testing.T) {
	h := handler.New(&mockService{}, zap.NewNop())
	req := httptest.NewRequest(http.MethodGet, "/predict", nil)
	rr := httptest.NewRecorder()
	h.Predict(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rr.Code)
	}
}
