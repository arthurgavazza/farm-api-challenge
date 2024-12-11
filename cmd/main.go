package main

import (
	"github.com/arthurgavazza/farm-api-challenge/internal/app/domain/usecases"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/config"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/database"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/httpapi"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/httpapi/controllers"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/httpapi/routers"
	shared "github.com/arthurgavazza/farm-api-challenge/internal/app/shared/logger"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
	"go.uber.org/fx"
)

// @title           Swagger Farms API
// @version         1.0
// @description     This is a farms API
// @termsOfService  http://swagger.io/terms/
// @contact.name    API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io
// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:8080
// @BasePath  /
// @externalDocs.description  OpenAPI
// @externalDocs.url  https://swagger.io/specification/         https://swagger.io/resources/open-api/
func main() {
	uuid.EnableRandPool()

	err := godotenv.Load(".env")
	if err != nil {
		log.Warn("Coudn't load .env file")
	}
	app := fx.New(
		shared.Module,
		config.Module,
		controllers.Module,
		usecases.Module,
		httpapi.Module,
		routers.Module,
		database.Module,
		fx.Invoke(func(*fasthttp.Server) {}),
		fx.NopLogger,
	)

	app.Run()
}
