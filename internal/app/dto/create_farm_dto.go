package dto

import (
	shared "github.com/arthurgavazza/farm-api-challenge/internal/app/shared/validation"
)

type CropProductionDTO struct {
	CropType    string `json:"crop_type" validate:"required,oneof=RICE CORN COFFEE SOYBEANS"`
	IsIrrigated bool   `json:"is_irrigated"`
	IsInsured   bool   `json:"is_insured"`
}

type CreateFarmDTO struct {
	Name            string              `json:"name" validate:"required"`
	LandArea        float64             `json:"land_area" validate:"required,gt=0"`
	UnitMeasure     string              `json:"unit_measure" validate:"required"`
	Address         string              `json:"address" validate:"required"`
	CropProductions []CropProductionDTO `json:"crop_productions" validate:"dive"`
}

func (dto *CreateFarmDTO) Validate() []shared.ErrorResponse {
	return shared.ValidateStruct(dto)
}
