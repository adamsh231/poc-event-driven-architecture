package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	topic := "hello"
	configuration := kafka.ConfigMap{
		"bootstrap.servers":    "localhost:9092",
		"group.id":             "group1",
		"enable.partition.eof": true,
	}

	// Variable p holds the new Consumer instance.
	consumer, err := kafka.NewConsumer(&configuration)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		err = consumer.Subscribe(topic, nil)
		if err != nil {
			print("Error:", err)
		} else {
			doTerm := true
			for doTerm {
				event := consumer.Poll(1000)
				if event == nil {
					fmt.Println("...")
					continue
				} else {
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
						// We've finished reading messages on this partition so let's wrap up
						// n.b. this is a BIG assumption that we are only consuming from one partition
						pe := event.(kafka.PartitionEOF)
						fmt.Printf("🌆 Got to the end of partition %v on topic %v at offset %v\n",
							pe.Partition,
							string(*pe.Topic),
							pe.Offset)
						doTerm = true

					case kafka.OffsetsCommitted:
						continue

					case kafka.Error:
						// It's an error
						em := event.(kafka.Error)
						fmt.Printf("☠️ Uh oh, caught an error:\n\t%v\n", em)

					default:
						// It's not anything we were expecting
						fmt.Printf("Got an event that's not a Message, Error, or PartitionEOF 👻\n\t%v\n", event)

					}
				}
			}
		}
	}
	defer consumer.Close()

}
