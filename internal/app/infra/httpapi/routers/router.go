package routers

import (
	_ "github.com/arthurgavazza/farm-api-challenge/docs"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type Router interface {
	Load()
}

func MakeRouter(
	farmRouter *FarmRouter,
	config *config.Config,
) *fiber.App {
	cfg := fiber.Config{
		AppName:       "farm-api by @arthurgavazza",
		CaseSensitive: true,
	}

	r := fiber.New(cfg)
	r.Get("/swagger/*", swagger.HandlerDefault) // default
	r.Get("/healthcheck", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "healthy",
			"appName": "farm-api",
		})
	})

	farmRouter.Load(r)

	return r
}
