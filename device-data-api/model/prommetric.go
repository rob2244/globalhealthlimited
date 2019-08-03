package model

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// NewPromDeviceMetricStore returns a new instance of NewPromDeviceMetricStore
func NewPromDeviceMetricStore() *PromDeviceMetricStore {
	return &PromDeviceMetricStore{Metrics: make(map[string]prometheus.Histogram)}
}

// PromDeviceMetricStore reperesents a Prometheus api to adding metrics
type PromDeviceMetricStore struct {
	Metrics map[string]prometheus.Histogram
}

// SaveDeviceMetric implements the DeviceMetricStore and creates a new metric
func (prom *PromDeviceMetricStore) SaveDeviceMetric(metric DeviceMetric) error {
	metricKey := getMetricKey(metric)
	met, ok := prom.Metrics[metricKey]

	if !ok {
		met = createMetric(metric)
		prom.Metrics[metricKey] = met
	}

	met.Observe(metric.Value)
	return nil
}

func getMetricKey(metric DeviceMetric) string {
	return fmt.Sprintf("%s-%s", metric.Name, metric.DeviceKey)
}

func createMetric(metric DeviceMetric) prometheus.Histogram {
	return promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "ghl",
		Subsystem: metric.DeviceKey,
		Name:      metric.Name,
		Help: fmt.Sprintf("Metric Name: %s, Device Identifier: %s",
			metric.Name, metric.DeviceKey),
		ConstLabels: prometheus.Labels{"device_id": metric.DeviceKey, "metric_name": metric.Name},
	})
}
