package core

import (
	"math/rand"
	"time"
)

// DeviceMetricStore represents presistance of device metrics to a data store
type DeviceMetricStore interface {
	SaveDeviceMetric(metric DeviceMetric) error
}

// DeviceMetric represents data sent to the api from devices in the field
type DeviceMetric struct {
	DeviceKey string     `json:"deviceKey"`
	Value     float64    `json:"value"`
	Name      MetricName `json:"name"`
	Unit      string     `json:"unit"`
	Timestamp time.Time  `json:"timestamp"`
	Latitude  float32    `json:"latitude"`
	Longitude float32    `json:"longitude"`
}

// MetricName reperesents a string constant for the names of metrics
type MetricName string

// Constants representing the available measurement types in the system
const (
	HeartRate       MetricName = "heart_rate"
	BloodSugar      MetricName = "blood_sugar"
	LungCapacity    MetricName = "lung_capacity"
	NumberOfToes    MetricName = "number_of_toes"
	NumberOfFingers MetricName = "number_of_fingers"
	Weight          MetricName = "weight"
	Height          MetricName = "height"
	SpinalFluid     MetricName = "spinal_fluid"
	Blood           MetricName = "blood"
)

var metricNames = []MetricName{
	HeartRate, BloodSugar, LungCapacity,
	NumberOfToes, NumberOfFingers, Weight,
	Height, SpinalFluid, Blood,
}

var metricUnits = map[MetricName]string{
	HeartRate: "BPM", BloodSugar: "mmol/L", LungCapacity: "L",
	NumberOfToes: "base10", NumberOfFingers: "base10", Weight: "LBS",
	Height: "FT", SpinalFluid: "L", Blood: "L",
}

// GetMetricName gets a random Device metric name
func GetMetricName() (MetricName, string) {
	idx := rand.Intn(len(metricNames) - 1)
	name := metricNames[idx]
	return name, metricUnits[name]
}
