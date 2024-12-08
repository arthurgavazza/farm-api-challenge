package dto

import (
	"github.com/go-playground/validator/v10"
)


type CropProductionDTO struct {
	CropType    string `json:"crop_type" validate:"required,oneof=RICE CORN COFFEE SOYBEANS"`
	IsIrrigated bool   `json:"is_irrigated"`
	IsInsured   bool   `json:"is_insured"`
}

type CreateFarmDTO struct {
	Name        string              `json:"name" validate:"required"`
	LandArea    float64             `json:"land_area" validate:"required,gt=0"`
	UnitMeasure string              `json:"unit_measure" validate:"required"`
	Address     string              `json:"address" validate:"required"`
	Productions []CropProductionDTO `json:"productions" validate:"dive"`
}



func (dto *CreateFarmDTO) Validate() *string{
	validate := validator.New()
	err := validate.Struct(dto)
   
	if err != nil {
		errorMessage := err.(validator.ValidationErrors).Error()
		return &errorMessage
	
	} 
	return nil
}