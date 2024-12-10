package routers

import (
	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/httpapi/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type FarmRouter struct {
	controller *controllers.FarmController
}

func (f *FarmRouter) Load(r *fiber.App) {
	log.Info("Loading farm routes")
	r.Post("/farms", f.controller.CreateFarm)
	r.Get("/farms", f.controller.ListFarms)
	r.Delete("/farms/:id", f.controller.DeleteFarm)
}

func NewFarmRouter(
	controller *controllers.FarmController,
) *FarmRouter {
	return &FarmRouter{
		controller: controller,
	}
}
