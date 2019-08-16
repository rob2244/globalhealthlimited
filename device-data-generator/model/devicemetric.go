package model

import (
	"globalhealthlimited/core"
	"math/rand"
	"time"

	"github.com/icrowley/fake"
)

// GenerateDeviceData generates a new DeviceMetric
func GenerateDeviceData(client *DeviceDataClient) (*core.DeviceMetric, error) {
	time := time.Now()
	mn, unit := core.GetMetricName()

	devKey, err := client.GetDeviceKey()
	if err != nil {
		return nil, err
	}

	return &core.DeviceMetric{
		DeviceKey: devKey,
		Value:     rand.Float64() * 150,
		Name:      mn,
		Unit:      unit,
		Timestamp: time,
		Latitude:  fake.Latitude(),
		Longitude: fake.Longitude(),
	}, nil
}
