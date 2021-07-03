package domain

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-redis/redis/v8"
)

type Config struct {
	KafkaConsumer *kafka.Consumer
	RedisClient *redis.Client
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

	// redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	ping := redisClient.Ping(context.Background())
	if ping.Err() != nil {
		return config, ping.Err()
	}
	config.RedisClient = redisClient

	return config, err
}

