package domain

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Config struct {
	App *fiber.App
	KafkaProducer *kafka.Producer
	Validator *validator.Validate
}

func LoadConfig() (config Config, err error){

	// kafka
	configuration := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	}
	producer, err := kafka.NewProducer(configuration)
	if err != nil {
		return config, err
	}
	config.KafkaProducer = producer

	// validator
	config.Validator = validator.New()

	// Fiber
	config.App = fiber.New()

	return config, err
}