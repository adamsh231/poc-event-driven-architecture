package handler

import (
	"fmt"
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
	topics := []string{
		"transaction",
	}
	fmt.Println("Listening on topic:", topics)

	kafkaHelper := helpers.NewKafkaHelper(handler.Config.KafkaConsumer)
	for _, topic := range topics {
		kafkaHelper.AddHandler(topic)
	}

}
