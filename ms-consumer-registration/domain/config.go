package domain

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-redis/redis/v8"
)

type Config struct {
	KafkaConfig *kafka.ConfigMap
	RedisClient   *redis.Client
}

func LoadConfig() (config Config, err error) {

	// kafka config
	config.KafkaConfig = &kafka.ConfigMap{
		"bootstrap.servers":        "10.130.105.161:9092",
		"group.id":                 "consumer-registration",
		"go.events.channel.enable": true,
		"enable.partition.eof":     true,
	}

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