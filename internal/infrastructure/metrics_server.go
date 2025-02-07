package infrastructure

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	ShutdownTimeoutSec = 3
)

type MetricsServer struct {
	Port   int
	Server *http.Server
}

func NewMetricsServer(port int) *MetricsServer {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	return &MetricsServer{
		Port: port,
		Server: &http.Server{
			Handler: mux,
			Addr:    fmt.Sprintf(":%d", port),
		},
	}
}

func (m *MetricsServer) Start(ctx context.Context) {

	go func() {
		if err := m.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return
		}
	}()

	<-ctx.Done()
	err := m.Stop()
	if err != nil {
		return
	}

	return
}

func (m *MetricsServer) Stop() error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), ShutdownTimeoutSec*time.Second)
	defer cancel()

	return m.Server.Shutdown(shutdownCtx)
}
