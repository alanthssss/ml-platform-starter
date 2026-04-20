package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alanthssss/ml-platform-starter/services/api-gateway-go/internal/handler"
	"github.com/alanthssss/ml-platform-starter/services/api-gateway-go/internal/metrics"
	"github.com/alanthssss/ml-platform-starter/services/api-gateway-go/internal/service"
	"github.com/alanthssss/ml-platform-starter/services/api-gateway-go/internal/triton"
	"go.uber.org/zap"
)

func main() {
	log, _ := zap.NewProduction()
	defer log.Sync()

	cfg := loadConfig()
	metrics.Register()

	tritonClient := triton.NewHTTPClient(cfg.TritonURL, cfg.TritonTimeout)
	svc := service.NewPredictService(tritonClient, log)
	h := handler.New(svc, log)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", h.Healthz)
	mux.HandleFunc("/readyz", h.Readyz)
	mux.Handle("/metrics", metrics.Handler())
	mux.HandleFunc("/predict", h.Predict)

	srv := &http.Server{
		Addr:         cfg.ListenAddr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("starting server", zap.String("addr", cfg.ListenAddr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("shutdown error", zap.Error(err))
	}
	log.Info("server stopped")
}

type config struct {
	ListenAddr    string
	TritonURL     string
	TritonTimeout time.Duration
}

func loadConfig() config {
	addr := getEnv("LISTEN_ADDR", ":8080")
	tritonURL := getEnv("TRITON_URL", "http://triton:8000")
	tritonTimeout := 10 * time.Second
	return config{
		ListenAddr:    addr,
		TritonURL:     tritonURL,
		TritonTimeout: tritonTimeout,
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
