package service

import (
	"context"
	"fmt"

	"github.com/alanthssss/ml-platform-starter/services/api-gateway-go/internal/metrics"
	"github.com/alanthssss/ml-platform-starter/services/api-gateway-go/internal/triton"
	"go.uber.org/zap"
)

var irisLabels = []string{"setosa", "versicolor", "virginica"}

// PredictRequest is the incoming JSON body for /predict.
type PredictRequest struct {
	Inputs []float32 `json:"inputs"`
}

// PredictResponse is the JSON response from /predict.
type PredictResponse struct {
	Model      string `json:"model"`
	Prediction int    `json:"prediction"`
	Label      string `json:"label"`
}

// PredictService defines the prediction interface.
type PredictService interface {
	Predict(ctx context.Context, req PredictRequest) (PredictResponse, error)
}

type predictService struct {
	triton triton.Client
	log    *zap.Logger
}

// NewPredictService creates a new predictService.
func NewPredictService(tc triton.Client, log *zap.Logger) PredictService {
	return &predictService{triton: tc, log: log}
}

func (s *predictService) Predict(ctx context.Context, req PredictRequest) (PredictResponse, error) {
	if len(req.Inputs) != 4 {
		return PredictResponse{}, fmt.Errorf("expected 4 inputs, got %d", len(req.Inputs))
	}

	timer := metrics.NewTritonTimer()
	defer timer.ObserveDuration()

	classIdx, err := s.triton.Infer(ctx, "iris_onnx", req.Inputs)
	if err != nil {
		return PredictResponse{}, fmt.Errorf("triton infer: %w", err)
	}

	if classIdx < 0 || classIdx >= len(irisLabels) {
		return PredictResponse{}, fmt.Errorf("invalid class index %d", classIdx)
	}

	return PredictResponse{
		Model:      "iris_onnx",
		Prediction: classIdx,
		Label:      irisLabels[classIdx],
	}, nil
}
