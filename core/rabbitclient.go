package core

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

// NewRabbitClient returns a new instance of RabbitClient
func NewRabbitClient() (*RabbitClient, error) {
	conn, err := amqp.Dial(buildQueueConnString())
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitClient{conn, ch}, nil
}

// RabbitClient reperesents a rabbit mq client
type RabbitClient struct {
	*amqp.Connection
	*amqp.Channel
}

// Close closes all underlying rabbit mq resources
func (client *RabbitClient) Close() {
	client.Connection.Close()
	client.Channel.Close()
}

// CreateQueue creates a rabbit mq queue with the specified name if it doesn't exist
func (client *RabbitClient) CreateQueue(queueName string) (*amqp.Queue, error) {
	q, err := client.QueueDeclare(
		queueName,
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
