package router

import (
	"fmt"
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

		kafkaHelper := helpers.NewKafkaHelper(router.Config.KafkaProducer)
		if err := kafkaHelper.Publish(input.Topic, input.Message); err != nil {
			return ctx.SendString(err.Error())
		}

		fmt.Println(input)

		return ctx.JSON(input)
	})

	router.App.Post("/publish/generator", func(ctx *fiber.Ctx) error {

		input := new(requests.PublishRandomRequest)
		if err := ctx.BodyParser(input); err != nil {
			return ctx.SendString(err.Error())
		}
		if err := router.Config.Validator.Struct(input); err != nil {
			return ctx.SendString(err.Error())
		}

		go router.sendMessage("transaction")
		go router.sendMessage("order")

		return ctx.JSON("success")
	})
}

func (router Router) sendMessage(topic string) {
	for i := 0; i < 10000; i++ {
		message := map[string]interface{}{
			"id_"+topic:   i,
		}

		kafkaHelper := helpers.NewKafkaHelper(router.Config.KafkaProducer)

		if err := kafkaHelper.Publish(topic, message); err != nil{
			fmt.Println(err.Error())
		}

		fmt.Println(topic, message)
	}
}
