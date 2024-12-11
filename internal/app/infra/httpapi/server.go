package httpapi

import (
	"context"
	"fmt"

	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/config"
	shared "github.com/arthurgavazza/farm-api-challenge/internal/app/shared/logger"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"go.uber.org/fx"
)

func NewServer(
	lifecycle fx.Lifecycle,
	router *fiber.App,
	config *config.Config,
	logger *shared.Logger,
	_ *gorm.DB,
) *fasthttp.Server {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Info(ctx, "Starting the server...")
				addr := fmt.Sprintf(":%s", config.Server.Port)
				if err := router.Listen(addr); err != nil {
					logger.Fatal(ctx, "Error starting the server: %s\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info(ctx, "Stopping the server...")
			logger.Close()
			return router.ShutdownWithContext(ctx)
		},
	})

	return router.Server()
}
