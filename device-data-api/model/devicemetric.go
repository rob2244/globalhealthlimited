package model

import "time"

// DeviceMetricStore represents presistance of device metrics to a data store
type DeviceMetricStore interface {
	SaveDeviceMetric(metric DeviceMetric) error
}

// DeviceMetric represents data sent to the api from devices in the field
type DeviceMetric struct {
	DeviceKey string    `json:"deviceKey"`
	Value     float64   `json:"value"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}
