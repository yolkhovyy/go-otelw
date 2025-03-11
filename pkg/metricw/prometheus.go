package metricw

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (m *Metric) RegisterPrometheusCollectors(ctx context.Context, colls ...prometheus.Collector) error {
	m.registry = prometheus.NewRegistry()
	m.registry.MustRegister(colls...)

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		_ = http.ListenAndServe(":9464", nil)
	}()

	return nil
}
