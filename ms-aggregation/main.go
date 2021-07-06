package main

import (
	"log"
	"ms-aggregation/domain"
	"ms-aggregation/router"
)

func main() {
	config, err := domain.LoadConfig()
	if err != nil {
		log.Fatal("Error while loading configuration :", err.Error())
	}
	defer config.KafkaProducer.Close()

	routing := router.NewRouter(config.App, config)
	routing.RegisterRouter()
	log.Fatal(config.App.Listen(":3000"))

}
