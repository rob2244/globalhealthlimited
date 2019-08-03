package model

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

// NewRabbitDeviceMetricStore returns a new instance of RabbitDeviceMetricStore
func NewRabbitDeviceMetricStore() (*RabbitDeviceMetricStore, error) {
	conn, err := amqp.Dial(buildQueueConnString())
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitDeviceMetricStore{conn, ch}, nil
}

// RabbitDeviceMetricStore reperesents a rabbit mq api to adding metrics
type RabbitDeviceMetricStore struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// Close closes all underlying rabbit mq resources
func (store *RabbitDeviceMetricStore) Close() {
	store.channel.Close()
	store.conn.Close()
}

// SaveDeviceMetric implements the DeviceMetricStore and pushes metric data to Rabbit MQ
func (store *RabbitDeviceMetricStore) SaveDeviceMetric(metric DeviceMetric) error {
	q, err := store.createQueue()
	if err != nil {
		return err
	}

	json, err := json.Marshal(metric)
	if err != nil {
		return err
	}

	err = store.channel.Publish("",
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

func (store *RabbitDeviceMetricStore) createQueue() (*amqp.Queue, error) {
	q, err := store.channel.QueueDeclare(
		"ghl-device-metrics",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	return &q, nil
}

func buildQueueConnString() string {
	host := "localhost"
	password := "guest"
	username := "guest"

	env := os.Getenv("ENVIRONMENT")

	if env == "PROD" {
		host = "rabbitmq"
		password = os.Getenv("RABBITMQ_PASSWORD")
		username = "user"
	}

	return fmt.Sprintf("amqp://%s:%s@%s:5672", username, password, host)
}
