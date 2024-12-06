package routers

import (
	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/config"
	"github.com/gofiber/fiber/v2"
)

type Router interface {
	Load()
}

func MakeRouter(
	config *config.Config,
) *fiber.App {
	cfg := fiber.Config{
		AppName:       "farm-api by @arthurgavazza",
		CaseSensitive: true,
	}

	r := fiber.New(cfg)

	r.Get("/healthcheck", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":   "healthy",
			"appName":  "farm-api",
		})
	})

	return r
}