package service_test

import (
	"context"
	"testing"

	"github.com/alanthssss/ml-platform-starter/services/api-gateway-go/internal/service"
	"go.uber.org/zap"
)

// mockTriton satisfies triton.Client for testing.
type mockTriton struct {
	result int
	err    error
}

func (m *mockTriton) Infer(_ context.Context, _ string, _ []float32) (int, error) {
	return m.result, m.err
}

func TestPredict_Setosa(t *testing.T) {
	svc := service.NewPredictService(&mockTriton{result: 0}, zap.NewNop())
	resp, err := svc.Predict(context.Background(), service.PredictRequest{
		Inputs: []float32{5.1, 3.5, 1.4, 0.2},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Label != "setosa" {
		t.Errorf("expected setosa, got %s", resp.Label)
	}
	if resp.Prediction != 0 {
		t.Errorf("expected 0, got %d", resp.Prediction)
	}
}

func TestPredict_WrongInputs(t *testing.T) {
	svc := service.NewPredictService(&mockTriton{}, zap.NewNop())
	_, err := svc.Predict(context.Background(), service.PredictRequest{
		Inputs: []float32{1.0},
	})
	if err == nil {
		t.Fatal("expected error for wrong number of inputs")
	}
}
