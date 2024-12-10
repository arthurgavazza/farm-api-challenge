package controllers

import (
	"fmt"
	"strconv"
	"strings"

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
	if errs := dto.Validate(); len(errs) > 0 && errs[0].Error {
		errMsgs := make([]string, 0)
		for _, err := range errs {
			errMsgs = append(errMsgs, fmt.Sprintf(
				"[%s]: '%v' | Needs to implement '%s'",
				err.FailedField,
				err.Value,
				err.Tag,
			))
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": strings.Join(errMsgs, " and ")})
	}
	var productions []domain.CropProduction
	for _, production := range dto.CropProductions {
		domainCropProduction := domain.CropProduction{
			CropType:    production.CropType,
			IsInsured:   production.IsInsured,
			IsIrrigated: production.IsIrrigated,
		}
		productions = append(productions, domainCropProduction)
	}
	farm, err := fc.createFarmUsecase.Execute(c.Context(), domain.Farm{
		Name:            dto.Name,
		LandArea:        dto.LandArea,
		UnitMeasure:     dto.UnitMeasure,
		Address:         dto.Address,
		CropProductions: productions,
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
		Page:    c.QueryInt("page", 1),
		PerPage: c.QueryInt("per_page", 10),
	}
	if cropType, exists := queries["crop_type"]; exists {
		searchParameters.CropType = &cropType
	}

	if minLandAreaStr, exists := queries["minimum_land_area"]; exists {
		landArea, err := strconv.ParseFloat(minLandAreaStr, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": `Query parameter "minimum_land_area" must be a valid floating-point number`,
			})
		}
		searchParameters.MinimumLandArea = &landArea
	}
	if maximumLandAreaStr, exists := queries["maximum_land_area"]; exists {
		landArea, err := strconv.ParseFloat(maximumLandAreaStr, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": `Query parameter "maximum_land_area" must be a valid floating-point number`,
			})
		}
		searchParameters.MaximumLandArea = &landArea
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
