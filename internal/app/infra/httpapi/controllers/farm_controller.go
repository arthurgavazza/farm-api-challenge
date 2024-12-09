package controllers

import (
	"strconv"

	"github.com/arthurgavazza/farm-api-challenge/internal/app/domain"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/domain/usecases"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/dto"
	"github.com/gofiber/fiber/v2"
)

type FarmController struct {
	createFarmUsecase usecases.CreateFarmUseCase
	listFarmsUseCase  usecases.ListFarmsUseCase
}

func (fc *FarmController) CreateFarm(c *fiber.Ctx) error {
	var dto dto.CreateFarmDTO
	if err := c.BodyParser(&dto); err != nil {
		return err
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
			CropType:    production.CropType,
			IsInsured:   production.IsInsured,
			IsIrrigated: production.IsIrrigated,
		}
		productions = append(productions, domainCropProduction)
	}
	farm, err := fc.createFarmUsecase.Execute(c.Context(), domain.Farm{
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

func (fc *FarmController) ListFarms(c *fiber.Ctx) error {
	queries := c.Queries()
	searchParameters := &domain.FarmSearchParameters{
		Page:    0,
		PerPage: 10,
	}
	if cropType, exists := queries["crop_type"]; exists {
		searchParameters.CropType = &cropType
	}
	if landAreaStr, exists := queries["land_area"]; exists {
		landArea, err := strconv.ParseFloat(landAreaStr, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": `Query parameter "land_area" must be a valid floating-point number`,
			})
		}
		searchParameters.LandArea = &landArea
	}
	if pageStr, exists := queries["page"]; exists {
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": `Query parameter "page" must be a valid integer`,
			})
		}
		searchParameters.Page = page
	}
	result, err := fc.listFarmsUseCase.Execute(c.Context(), searchParameters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(result)
}

func NewFarmController(
	createFarmUsecase usecases.CreateFarmUseCase,
	listFarmsUsecase usecases.ListFarmsUseCase,
) *FarmController {
	return &FarmController{
		createFarmUsecase: createFarmUsecase,
		listFarmsUseCase:  listFarmsUsecase,
	}
}
