package main

import (
	"fmt"
	"globalhealthlimited/data-ingester/data"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func main() {

	conn, err := amqp.Dial(buildQueueConnString())
	handleError(err)

	defer conn.Close()
	defer data.Close()

	ch, err := conn.Channel()
	handleError(err)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"ghl-device-url",
		false,
		false,
		false,
		false,
		nil,
	)
	handleError(err)

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	handleError(err)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			url := string(d.Body)
			log.Printf("Download of url %s starting", url)

			filepath, err := data.DownloadFile(url)
			if err != nil {
				d.Reject(true)
				handleError(err)
			}

			log.Printf("Download Complete, saving file %s to database...", filepath)
			defer os.Remove(filepath)

			if err = data.UnzipAndSave(filepath); err != nil {
				d.Reject(true)
				handleError(err)
			}

			log.Printf("Successfully saved file %s to database", filepath)
			d.Ack(false)
		}
	}()

	<-forever
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

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
