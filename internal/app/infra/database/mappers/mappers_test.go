package mappers

import (
	"testing"

	"github.com/arthurgavazza/farm-api-challenge/internal/app/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestToGormFarm(t *testing.T) {
	domainFarm := &domain.Farm{
		ID:          uuid.New(),
		Name:        "Test Farm",
		LandArea:    100.5,
		UnitMeasure: "hectares",
		Address:     "123 Farm Lane",
		Productions: []domain.CropProduction{
			{
				ID:          uuid.New(),
				CropType:    domain.CropTypeCoffee.String(),
				IsIrrigated: true,
				IsInsured:   false,
				FarmID:      uuid.New(),
			},
			{
				ID:          uuid.New(),
				CropType:    domain.CropTypeRice.String(),
				IsIrrigated: false,
				IsInsured:   true,
				FarmID:      uuid.New(),
			},
		},
	}

	result := ToGormFarm(domainFarm)

	assert.NotNil(t, result)
	assert.Equal(t, domainFarm.ID, result.ID)
	assert.Equal(t, domainFarm.Name, result.Name)
	assert.Equal(t, domainFarm.LandArea, result.LandArea)
	assert.Equal(t, domainFarm.UnitMeasure, result.UnitMeasure)
	assert.Equal(t, domainFarm.Address, result.Address)
	assert.Len(t, result.Productions, len(domainFarm.Productions))

	for i, crop := range domainFarm.Productions {
		assert.Equal(t, crop.CropType, result.Productions[i].CropType)
		assert.Equal(t, crop.IsIrrigated, result.Productions[i].IsIrrigated)
		assert.Equal(t, crop.IsInsured, result.Productions[i].IsInsured)
		assert.Equal(t, crop.ID, result.Productions[i].ID)
		assert.Equal(t, crop.FarmID, result.Productions[i].FarmID)
	}
}

func TestToGormCropProductions(t *testing.T) {
	domainCrops := []domain.CropProduction{
		{
			ID:          uuid.New(),
			CropType:    domain.CropTypeCoffee.String(),
			IsIrrigated: true,
			IsInsured:   false,
			FarmID:      uuid.New(),
		},
		{
			ID:          uuid.New(),
			CropType:    domain.CropTypeRice.String(),
			IsIrrigated: false,
			IsInsured:   true,
			FarmID:      uuid.New(),
		},
	}

	result := ToGormCropProductions(domainCrops)

	assert.Len(t, result, len(domainCrops))
	for i, crop := range domainCrops {
		assert.Equal(t, crop.CropType, result[i].CropType)
		assert.Equal(t, crop.IsIrrigated, result[i].IsIrrigated)
		assert.Equal(t, crop.IsInsured, result[i].IsInsured)
		assert.Equal(t, crop.ID, result[i].ID)
		assert.Equal(t, crop.FarmID, result[i].FarmID)
	}
}
