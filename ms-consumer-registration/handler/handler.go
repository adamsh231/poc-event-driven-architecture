package handler

import (
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
	kafkaHelper := helpers.NewKafkaHelper(handler.Config.KafkaConsumer)
	redisHelper := helpers.NewRedisHelper(handler.Config.RedisClient)

	// transaction topic
	kafkaHelper.AddHandler("transaction", func(message *kafka.Message) {
		redisHelper.QueueingOutlet(message)
		time.Sleep(3 * time.Second)
	})

}
