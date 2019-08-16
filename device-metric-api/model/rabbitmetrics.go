package model

import (
	"encoding/json"
	"globalhealthlimited/core"

	"github.com/streadway/amqp"
)

// RabbitDeviceMetricStore reperesents a rabbit mq api to adding metrics
type RabbitDeviceMetricStore struct {
	*core.RabbitClient
}

// NewRabbitDeviceMetricStore returns a new instance of RabbitDeviceMetricStore
func NewRabbitDeviceMetricStore() (*RabbitDeviceMetricStore, error) {
	client, err := core.NewRabbitClient()
	if err != nil {
		return nil, err
	}

	return &RabbitDeviceMetricStore{client}, nil
}

// SaveDeviceMetric implements the DeviceMetricStore and pushes metric data to Rabbit MQ
func (store *RabbitDeviceMetricStore) SaveDeviceMetric(metric core.DeviceMetric) error {
	q, err := store.CreateQueue("ghl-device-metrics")
	if err != nil {
		return err
	}

	json, err := json.Marshal(metric)
	if err != nil {
		return err
	}

	err = store.Publish("",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(json),
		})
	if err != nil {
		return err
	}

	return nil
}
