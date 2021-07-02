package router

import (
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

		return ctx.JSON(input)
	})
}
