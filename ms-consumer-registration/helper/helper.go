package helper

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

type KafkaHelper struct {
	KafkaConsumer *kafka.Consumer
}

func NewKafkaHelper(kafkaConsumer *kafka.Consumer) KafkaHelper {
	return KafkaHelper{KafkaConsumer: kafkaConsumer}
}

func (helper KafkaHelper) AddHandler(topic string) {

	doneChan := make(chan bool)

	if err := helper.KafkaConsumer.Subscribe(topic, nil); err != nil {
		log.Fatal("Failed to add handler:", err.Error())
	}

	// Handle the events that we get
	go func() {
		for {
			select {
			case event := <-helper.KafkaConsumer.Events():
				switch event.(type) {

				case *kafka.Message:
					// It's a message
					km := event.(*kafka.Message)
					fmt.Printf("✅ Message '%v' received from topic '%v' (partition %d at offset %d)\n",
						string(km.Value),
						string(*km.TopicPartition.Topic),
						km.TopicPartition.Partition,
						km.TopicPartition.Offset)

				case kafka.PartitionEOF:
					pe := event.(kafka.PartitionEOF)
					fmt.Printf("🌆 Got to the end of partition %v on topic %v at offset %v\n",
						pe.Partition,
						string(*pe.Topic),
						pe.Offset)

				case kafka.OffsetsCommitted:
					continue

				case kafka.Error:
					em := event.(kafka.Error)
					fmt.Printf("☠️ Uh oh, caught an error:\n\t%v\n", em)

				default:
					fmt.Printf("Got an event that's not a Message, Error, or PartitionEOF 👻\n\t%v\n", event)

				}

			}
		}
	}()

	<-doneChan // close(doneChan)

}
