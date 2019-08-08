package model

import (
	"math/rand"
	"time"
)

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

// DeviceMetric represents data sent to the api from devices in the field
type DeviceMetric struct {
	DeviceKey string     `json:"deviceKey"`
	Value     float64    `json:"value"`
	Name      MetricName `json:"name"`
	Timestamp time.Time  `json:"timestamp"`
}

// GenerateDeviceData generates a new DeviceMetric
func GenerateDeviceData(client *DeviceDataClient) (*DeviceMetric, error) {
	time := time.Now()
	mn := getMetricName()

	devKey, err := client.GetDeviceKey()
	if err != nil {
		return nil, err
	}

	return &DeviceMetric{
		devKey,
		rand.Float64() * 150,
		mn,
		time,
	}, nil
}

func getMetricName() MetricName {
	idx := rand.Intn(len(metricNames) - 1)
	return metricNames[idx]
}
