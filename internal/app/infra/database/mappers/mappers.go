package mappers

import (
	"github.com/arthurgavazza/farm-api-challenge/internal/app/domain"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/database"
)

func ToGormFarm(domainFarm *domain.Farm) *database.Farm {
	return &database.Farm{
		ID:          domainFarm.ID,
		Name:        domainFarm.Name,
		LandArea:    domainFarm.LandArea,
		UnitMeasure: domainFarm.UnitMeasure,
		Address:     domainFarm.Address,
		Productions: ToGormCropProductions(domainFarm.Productions),
	}
}

func ToGormCropProductions(domainCrops []domain.CropProduction) []database.CropProduction {
	var crops []database.CropProduction
	for _, crop := range domainCrops {
		crops = append(crops, database.CropProduction{
			CropType:    crop.CropType,
			IsIrrigated: crop.IsIrrigated,
			IsInsured:   crop.IsInsured,
		})
	}
	return crops
}
