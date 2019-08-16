package main

import (
	"encoding/json"
	"fmt"
	"globalhealthlimited/core"
	"log"
	"strings"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	client, err := core.NewRabbitClient()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	q, err := client.CreateQueue("ghl-device-metrics")

	if err != nil {
		log.Fatal(err)
	}

	msgs, err := client.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	metricsChan := make(chan *core.DeviceMetric, 10)

	go getMetrics(metricsChan, msgs)

	report, err := createReport(metricsChan)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(report)
}

func getMetrics(metricChannel chan<- *core.DeviceMetric, msgs <-chan amqp.Delivery) {
	defer close(metricChannel)

	for i := 0; i < 10; i++ {
		msg := <-msgs
		met, err := processMessage(msg)
		if err != nil {
			msg.Nack(false, true)
			log.Fatal(err)
		}

		metricChannel <- met

		err = msg.Ack(false)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func processMessage(dlv amqp.Delivery) (*core.DeviceMetric, error) {
	met := &core.DeviceMetric{}
	err := json.Unmarshal(dlv.Body, met)

	if err != nil {
		return nil, err
	}

	return met, nil
}

func createReport(metrics <-chan *core.DeviceMetric) (string, error) {
	builder := strings.Builder{}

	str := printHeader()

	_, err := builder.WriteString(str)
	if err != nil {
		return "", err
	}

	for met := range metrics {
		str = fmt.Sprintf("%s|%-6.2f|%-20s|%s\n", met.DeviceKey, met.Value, met.Name, met.Timestamp.Format(time.RFC822))
		_, err = builder.WriteString(str)
		if err != nil {
			return "", err
		}
	}

	return builder.String(), nil
}

func printHeader() string {
	return fmt.Sprintf("%-36s|%-6s|%-20s|%s\n", "Device Key", "Value", "Name", "Timestamp")
}
