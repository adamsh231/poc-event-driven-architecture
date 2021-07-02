package domain

import "github.com/confluentinc/confluent-kafka-go/kafka"

type Config struct {
	KafkaConsumer *kafka.Consumer
}

func LoadConfig() (config Config, err error){

	// kafka
	configuration := kafka.ConfigMap{
		"bootstrap.servers":    "localhost:9092",
		"group.id":             "group1",
		"go.events.channel.enable": true,
		"enable.partition.eof": true,
	}
	consumer, err := kafka.NewConsumer(&configuration)
	if err != nil {
		return config, err
	}
	config.KafkaConsumer = consumer


	return config, err
}

