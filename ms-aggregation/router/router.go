package router

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"ms-aggregation/domain"
	"ms-aggregation/domain/requests"
	"ms-aggregation/helpers"
)

type Router struct {
	App    *fiber.App
	Config domain.Config
}

func NewRouter(app *fiber.App, config domain.Config) Router{
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
		if err := kafkaHelper.Publish(input.Topic, input.Message); err != nil{
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

		fake := gofakeit.New(0)
		message := map[string]interface{}{
			"name": fake.Name(),
			"gender": fake.Gender(),
			"pet": fake.Animal(),
			"friends": []string{
				fake.Name(),
				fake.Name(),
				fake.Name(),
			},
		}

		kafkaHelper := helpers.NewKafkaHelper(router.Config.KafkaProducer)
		if err := kafkaHelper.Publish(input.Topic, message); err != nil{
			return ctx.SendString(err.Error())
		}

		fmt.Println(message)

		return ctx.JSON(message)
	})
}
