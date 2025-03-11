package metricw

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (m *Metric) RegisterPrometheusServer(ctx context.Context, colls ...prometheus.Collector) error {
	m.prometheusRegistry = prometheus.NewRegistry()
	m.prometheusRegistry.MustRegister(colls...)

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		_ = http.ListenAndServe(":9464", nil)
	}()

	return nil
}
