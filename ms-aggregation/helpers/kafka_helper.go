package helpers

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaHelper struct {
	KafkaProducer *kafka.Producer
}

func NewKafkaHelper(kafkaProducer *kafka.Producer) KafkaHelper {
	return KafkaHelper{KafkaProducer: kafkaProducer}
}


func (helper KafkaHelper) Publish(topic string, message interface{}) (err error) {
	messageInBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	messageConfiguration := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: messageInBytes,
	}

	if err = helper.KafkaProducer.Produce(messageConfiguration, nil); err != nil {
		return err
	}

	// TODO: flush for acknowledge

	return err
}
