package main

import (
	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/config"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/database"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/httpapi"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/httpapi/routers"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
	"go.uber.org/fx"
)

func main() {
	
	uuid.EnableRandPool()
	
	err := godotenv.Load(".env")
	if err != nil {
		log.Warn("Coudn't load .env file")
	}
	log.Warn("Coudn't load .env file")
	app := fx.New(
		config.Module,
		httpapi.Module,
		routers.Module,
		database.Module,
		fx.Invoke(func(*fasthttp.Server) {}),
		fx.NopLogger,
	)

	app.Run()
}