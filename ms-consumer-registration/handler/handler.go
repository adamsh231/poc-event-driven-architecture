package handler

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"ms-consumer-registration/domain"
	"ms-consumer-registration/helpers"
)

type Handler struct {
	Config domain.Config
}

func NewHandler(config domain.Config) Handler {
	return Handler{Config: config}
}

func (handler Handler) RegisterHandler() {
	kafkaHelper := helpers.NewKafkaHelper(handler.Config.KafkaConsumer)

	// transaction topic
	kafkaHelper.AddHandler("transaction", func(message *kafka.Message) {
		fmt.Println("Do action here")
	})

}
