package main

import (
	"globalhealthlimited/device-data-api/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

func main() {
	r := gin.Default()
	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)

	promStore := model.NewPromDeviceMetricStore()
	rabbitStore, err := model.NewRabbitDeviceMetricStore()

	handleError(err)
	defer rabbitStore.Close()

	r.POST("/api/device", postDevice([]model.DeviceMetricStore{promStore, rabbitStore}))
	r.Run(":8080")
}

func postDevice(stores []model.DeviceMetricStore) func(c *gin.Context) {
	return func(c *gin.Context) {
		var json model.DeviceMetric

		if err := c.BindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Printf("Recieved metric for %s from device %s", json.Name, json.DeviceKey)
		for _, store := range stores {
			store.SaveDeviceMetric(json)
		}

		c.Status(http.StatusOK)
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
