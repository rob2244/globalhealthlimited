package main

import "time"

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

// DeviceMetric represents data sent to the api from devices in the field
type DeviceMetric struct {
	DeviceKey string    `json:"deviceKey"`
	Value     float64   `json:"value"`
	Name      string    `json:"name"`
	Unit      string    `json:"unit"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {

}

func generateDeviceData() {
	time := time.Now()
}
