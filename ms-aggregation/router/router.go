package router

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gofiber/fiber/v2"
	"ms-aggregation/domain"
	"ms-aggregation/domain/requests"
	"ms-aggregation/helpers"
)

type Router struct {
	App    *fiber.App
	Config domain.Config
}

func NewRouter(app *fiber.App, config domain.Config) Router {
	return Router{App: app, Config: config}
}

func (router Router) RegisterRouter() {

	router.App.Post("/publish", func(ctx *fiber.Ctx) error {

		input := new(requests.PublishRequest)
		if err := ctx.BodyParser(input); err != nil {
			return ctx.SendString(err.Error())
		}
		if err := router.Config.Validator.Struct(input); err != nil {
			return ctx.SendString(err.Error())
		}

		for i := 0; i < 10; i++ {
			kafkaHelper := helpers.NewKafkaHelper(router.Config.KafkaProducer)
			if err := kafkaHelper.Publish(input.Topic, input.Partition, fmt.Sprintf("%s - %d",input.Message,i)); err != nil {
				return ctx.SendString(err.Error())
			}
		}

		return ctx.JSON("success")
	})

	go router.watcher()

}

func (router Router) watcher() {
	for {
		select {
		case ev := <-router.Config.KafkaProducer.Events():
			switch ev.(type) {

			case *kafka.Message:
				km := ev.(*kafka.Message)
				if km.TopicPartition.Error != nil {
					fmt.Printf("☠️ Failed to send message '%v' to topic '%v'\n\tErr: %v",
						string(km.Value),
						string(*km.TopicPartition.Topic),
						km.TopicPartition.Error)
				} else {
					fmt.Printf("✅ Message '%v' delivered to topic '%v' (partition %d at offset %d)\n",
						string(km.Value),
						string(*km.TopicPartition.Topic),
						km.TopicPartition.Partition,
						km.TopicPartition.Offset)
				}

			case kafka.Error:
				// It's an error
				em := ev.(kafka.Error)
				fmt.Printf("☠️ Uh oh, caught an error:\n\t%v\n", em)
			default:
				// It's not anything we were expecting
				fmt.Printf("Got an event that's not a Message or Error 👻\n\t%v\n", ev)

			}
		}
	}
}
