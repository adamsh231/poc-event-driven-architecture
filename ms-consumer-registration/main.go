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

	doneChan := make(chan bool)
	handler.NewHandler(config).RegisterHandler()
	fmt.Println("Listening on broker...")
	<-doneChan
}

