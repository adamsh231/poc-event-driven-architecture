package main

import (
	"log"
	"ms-consumer-registration/domain"
	"ms-consumer-registration/handler"
)

func main() {
	config, err := domain.LoadConfig()
	if err != nil {
		log.Fatal("Error while loading configuration:", err.Error())
	}
	defer config.KafkaConsumer.Close()

	handler.NewHandler(config).RegisterHandler()
}

