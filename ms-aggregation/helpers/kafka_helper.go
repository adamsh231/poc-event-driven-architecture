package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaHelper struct {
	KafkaProducer *kafka.Producer
}

func NewKafkaHelper(kafkaProducer *kafka.Producer) KafkaHelper {
	return KafkaHelper{KafkaProducer: kafkaProducer}
}

func (helper KafkaHelper) Publish(topic string, partition int32, message interface{}) (err error) {
	messageInBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	if partition == -1 {
		partition = kafka.PartitionAny
	}

	messageConfiguration := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: partition,
		},
		Value: messageInBytes,
	}

	if err = helper.KafkaProducer.Produce(messageConfiguration, nil); err != nil {
		return err
	}

	if r := helper.KafkaProducer.Flush(10000); r > 0 {
		fmt.Printf("⚠️ Failed to flush all messages after 10 seconds. %d message(s) remain\n", r)
	} else {
		fmt.Println("✨ All messages flushed from the queue")
	}

	return err
}
