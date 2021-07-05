package handler

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"ms-consumer-registration/domain"
	"ms-consumer-registration/helpers"
	"time"
)

type Handler struct {
	Config domain.Config
}

func NewHandler(config domain.Config) Handler {
	return Handler{Config: config}
}

func (handler Handler) RegisterHandler() {

	kafkaHelper := helpers.NewKafkaHelper(handler.Config.KafkaConfig)
	//redisHelper := helpers.NewRedisHelper(handler.Config.RedisClient)

	kafkaHelper.AddHandlerToggle("transaction", func(message *kafka.Message) {
		//redisHelper.QueueingOutlet(message)
		time.Sleep(500 * time.Millisecond)
		fmt.Println("transaction_ok")
	})

	kafkaHelper.AddHandlerToggle("order", func(message *kafka.Message) {
		//redisHelper.QueueingOutlet(message)
		time.Sleep(500 * time.Millisecond)
		fmt.Println("order_ok")
	})

}
