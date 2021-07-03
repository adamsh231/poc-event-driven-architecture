package main

import (
	"fmt"
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

	fmt.Println("Listening on broker...")
	handler.NewHandler(config).RegisterHandler()
}

