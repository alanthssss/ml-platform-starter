package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// RequestsTotal counts incoming HTTP requests by endpoint.
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_gateway_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"endpoint"},
	)

	// RequestLatency tracks HTTP request latency by endpoint.
	RequestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_gateway_request_duration_seconds",
			Help:    "HTTP request latency in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint"},
	)

	// PredictionsTotal counts successful predictions.
	PredictionsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "api_gateway_predictions_total",
			Help: "Total number of successful predictions.",
		},
	)

	// TritonLatency tracks downstream Triton request latency.
	TritonLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "api_gateway_triton_request_duration_seconds",
			Help:    "Downstream Triton request latency in seconds.",
			Buckets: prometheus.DefBuckets,
		},
	)
)

// Register registers all metrics with the default Prometheus registry.
func Register() {
	prometheus.MustRegister(RequestsTotal, RequestLatency, PredictionsTotal, TritonLatency)
}

// Handler returns the Prometheus HTTP handler.
func Handler() http.Handler {
	return promhttp.Handler()
}

// tritonTimer is a helper for observing Triton request duration.
type tritonTimer struct {
	start time.Time
}

// NewTritonTimer starts a new timer for Triton request latency.
func NewTritonTimer() *tritonTimer {
	return &tritonTimer{start: time.Now()}
}

// ObserveDuration records the elapsed time.
func (t *tritonTimer) ObserveDuration() {
	TritonLatency.Observe(time.Since(t.start).Seconds())
}
