package mappers

import (
	"github.com/arthurgavazza/farm-api-challenge/internal/app/domain"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/database/entities"
)

func ToGormFarm(domainFarm *domain.Farm) *entities.Farm {
	return &entities.Farm{
		ID:          domainFarm.ID,
		Name:        domainFarm.Name,
		LandArea:    domainFarm.LandArea,
		UnitMeasure: domainFarm.UnitMeasure,
		Address:     domainFarm.Address,
		Productions: ToGormCropProductions(domainFarm.Productions),
	}
}

func ToGormCropProductions(domainCrops []domain.CropProduction) []entities.CropProduction {
	var crops []entities.CropProduction
	for _, crop := range domainCrops {
		crops = append(crops, entities.CropProduction{
			CropType:    crop.CropType,
			IsIrrigated: crop.IsIrrigated,
			IsInsured:   crop.IsInsured,
		})
	}
	return crops
}
