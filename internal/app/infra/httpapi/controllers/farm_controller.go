package controllers

import (
	"github.com/arthurgavazza/farm-api-challenge/internal/app/domain"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/domain/usecases"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FarmController struct {
	createFarmUsecase *usecases.CreateFarm
}

func (fc *FarmController) CreateFarm(c *fiber.Ctx) error {
	var dto dto.CreateFarmDTO

	if err := c.BodyParser(&dto); err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "bad request"})
	}
	validationErrors := dto.Validate() 
	if validationErrors != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": validationErrors})
	}
	var productions []domain.CropProduction
	for _, production := range dto.Productions {
		domainCropProduction := domain.CropProduction{
			ID:          uuid.New(),
			CropType:    production.CropType,
			IsInsured:   production.IsInsured,
			IsIrrigated: production.IsIrrigated,
		}
		productions = append(productions, domainCropProduction)
	}
	farm, err := fc.createFarmUsecase.Execute(c.Context(),domain.Farm{
		Name:        dto.Name,
		LandArea:    dto.LandArea,
		UnitMeasure: dto.UnitMeasure,
		Address:     dto.Address,
		Productions: productions,
	})
	if err != nil {
		return c.
			Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "internal server error"})
	}
	c.Set("Location", "/farms/"+farm.ID.String())
	return c.Status(fiber.StatusCreated).JSON(farm)
}

func NewFarmController(
	createFarmUsecase *usecases.CreateFarm,
) *FarmController {
	return &FarmController{
		createFarmUsecase: createFarmUsecase,
	}
}
