package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"globalhealthlimited/device-data-generator/model"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var numWorkers int
var reqPerWorkerPerMin int

func init() {
	const (
		nwDesc    = "Number of worker go routines to spawn"
		rpwpmDesc = "Number of minutes between worker Http post requests"
	)

	flag.IntVar(&numWorkers, "n", 10, nwDesc)
	flag.IntVar(&reqPerWorkerPerMin, "r", 1, rpwpmDesc)
	flag.Parse()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	metrics := make(chan *model.DeviceMetric, 100)
	defer close(metrics)

	for i := 0; i < numWorkers; i++ {
		go uploadDeviceData(metrics)
	}

	dbClient, err := model.NewDeviceDataClient()
	defer dbClient.Close()
	if err != nil {
		log.Fatal(err)
	}

	for {
		dd, err := model.GenerateDeviceData(dbClient)
		if err != nil {
			log.Fatal(err)
		}
		metrics <- dd
	}
}

func uploadDeviceData(metrics <-chan *model.DeviceMetric) {
	url := getURL()

	for metric := range metrics {
		payload, err := json.Marshal(metric)

		if err != nil {
			log.Fatalln(err)
		}

		_, err = http.Post(url, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			log.Fatalln(err)
		}

		time.Sleep(time.Duration(reqPerWorkerPerMin) * time.Minute)
	}
}

func getURL() string {
	if env := os.Getenv("ENVIRONMENT"); env == "PRODUCTION" {
		return "http://device-data-api-service/api/device"
	}

	return "http://localhost:8080/api/device"
}
