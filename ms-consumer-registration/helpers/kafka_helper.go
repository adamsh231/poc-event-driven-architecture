package helpers

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

type KafkaHelper struct {
	KafkaConfig *kafka.ConfigMap
}

func NewKafkaHelper(kafkaConfig *kafka.ConfigMap) KafkaHelper {
	return KafkaHelper{KafkaConfig: kafkaConfig}
}

func (helper KafkaHelper) AddHandlerToggle(topic string, action func(message *kafka.Message)) {
	helper.AddHandlerBasic(topic, action)
	//helper.AddHandlerWithWorker(topic, action)
}

func (helper KafkaHelper) AddHandlerBasic(topic string, action func(message *kafka.Message)) {

	consumer, err := kafka.NewConsumer(helper.KafkaConfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	if err = consumer.Subscribe(topic, nil); err != nil {
		log.Fatal(err.Error())
	}

	go func() {
		defer consumer.Close()
		for {
			select {
			case event := <-consumer.Events():
				switch event.(type) {

				case *kafka.Message:
					km := event.(*kafka.Message)
					fmt.Printf("✅ Message '%v' received from topic '%v' (partition %d at offset %d)\n",
						string(km.Value),
						string(*km.TopicPartition.Topic),
						km.TopicPartition.Partition,
						km.TopicPartition.Offset)

					// action
					action(km)

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
}

type ActionStruct struct {
	Action  func(message *kafka.Message)
	Message *kafka.Message
}

func (helper KafkaHelper) AddHandlerWithWorker(topic string, action func(message *kafka.Message)) {

	consumer, err := kafka.NewConsumer(helper.KafkaConfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	if err = consumer.Subscribe(topic, nil); err != nil {
		log.Fatal(err.Error())
	}

	bufferedChannel := make(chan ActionStruct, 5) // buffer

	go func() {
		defer consumer.Close()
		for {
			select {
			case event := <-consumer.Events():
				switch event.(type) {

				case *kafka.Message:
					km := event.(*kafka.Message)
					fmt.Printf("✅ Message '%v' received from topic '%v' (partition %d at offset %d)\n",
						string(km.Value),
						string(*km.TopicPartition.Topic),
						km.TopicPartition.Partition,
						km.TopicPartition.Offset)

					// action
					bufferedChannel <- ActionStruct{
						Action:  action,
						Message: km,
					}

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

	go func() {
		for {
			if len(bufferedChannel) == 5 {
				for val := range bufferedChannel {
					go val.Action(val.Message)
					fmt.Println("buffer length:", len(bufferedChannel))
				}
			}
		}
	}()
}
