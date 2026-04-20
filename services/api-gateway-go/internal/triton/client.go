package triton

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client is the interface for calling Triton Inference Server.
type Client interface {
	Infer(ctx context.Context, modelName string, inputs []float32) (int, error)
}

// httpClient calls Triton via its HTTP/REST API (KServe v2 protocol).
type httpClient struct {
	baseURL string
	http    *http.Client
}

// NewHTTPClient creates a new Triton HTTP client.
func NewHTTPClient(baseURL string, timeout time.Duration) Client {
	return &httpClient{
		baseURL: baseURL,
		http:    &http.Client{Timeout: timeout},
	}
}

// tritonRequest is the KServe v2 inference request body.
type tritonRequest struct {
	Inputs []tritonInput `json:"inputs"`
}

type tritonInput struct {
	Name     string    `json:"name"`
	Shape    []int     `json:"shape"`
	Datatype string    `json:"datatype"`
	Data     []float32 `json:"data"`
}

// tritonResponse is the KServe v2 inference response body.
type tritonResponse struct {
	Outputs []tritonOutput `json:"outputs"`
}

type tritonOutput struct {
	Name string  `json:"name"`
	Data []int64 `json:"data"`
}

// Infer sends inference request to Triton and returns the predicted class index.
func (c *httpClient) Infer(ctx context.Context, modelName string, inputs []float32) (int, error) {
	reqBody := tritonRequest{
		Inputs: []tritonInput{
			{
				Name:     "float_input",
				Shape:    []int{1, 4},
				Datatype: "FP32",
				Data:     inputs,
			},
		},
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return 0, fmt.Errorf("marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/v2/models/%s/infer", c.baseURL, modelName)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return 0, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return 0, fmt.Errorf("triton request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("triton status %d", resp.StatusCode)
	}

	var tritonResp tritonResponse
	if err := json.NewDecoder(resp.Body).Decode(&tritonResp); err != nil {
		return 0, fmt.Errorf("decode response: %w", err)
	}

	if len(tritonResp.Outputs) == 0 || len(tritonResp.Outputs[0].Data) == 0 {
		return 0, fmt.Errorf("empty triton response")
	}

	return int(tritonResp.Outputs[0].Data[0]), nil
}
