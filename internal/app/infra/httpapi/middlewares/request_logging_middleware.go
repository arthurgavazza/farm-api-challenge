package middlewares

import (
	"time"

	logger "github.com/arthurgavazza/farm-api-challenge/internal/app/shared/logger"
	"github.com/gofiber/fiber/v2"
)

func RequestLogger(log *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		log.Info(c.Context(), "Incoming request", map[string]interface{}{
			"method":  c.Method(),
			"path":    c.Path(),
			"query":   c.Context().QueryArgs().String(),
			"headers": c.GetReqHeaders(),
		})

		err := c.Next()

		duration := time.Since(start)
		statusCode := c.Response().StatusCode()
		if statusCode < 200 || statusCode >= 300 {
			log.Info(c.Context(), "Response sent", map[string]interface{}{
				"status":   statusCode,
				"duration": duration.String(),
				"body":     string(c.Response().Body()),
			})
		} else {
			log.Info(c.Context(), "Response sent", map[string]interface{}{
				"status":   statusCode,
				"duration": duration.String(),
			})
		}

		return err
	}
}
