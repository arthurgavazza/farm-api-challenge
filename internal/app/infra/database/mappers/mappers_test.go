package mappers

import (
	"testing"

	"github.com/arthurgavazza/farm-api-challenge/internal/app/domain"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/database/entities"
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
		CropProductions: []domain.CropProduction{
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
	assert.Len(t, result.CropProductions, len(domainFarm.CropProductions))

	for i, crop := range domainFarm.CropProductions {
		assert.Equal(t, crop.CropType, result.CropProductions[i].CropType)
		assert.Equal(t, crop.IsIrrigated, result.CropProductions[i].IsIrrigated)
		assert.Equal(t, crop.IsInsured, result.CropProductions[i].IsInsured)
		assert.Equal(t, crop.ID, result.CropProductions[i].ID)
		assert.Equal(t, crop.FarmID, result.CropProductions[i].FarmID)
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

func TestToDomainFarm(t *testing.T) {
	gormFarm := &entities.Farm{
		ID:          uuid.New(),
		Name:        "Test Farm",
		LandArea:    100.5,
		UnitMeasure: "hectares",
		Address:     "123 Farm Lane",
		CropProductions: []entities.CropProduction{
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

	result := ToDomainFarm(gormFarm)

	assert.NotNil(t, result)
	assert.Equal(t, gormFarm.ID, result.ID)
	assert.Equal(t, gormFarm.Name, result.Name)
	assert.Equal(t, gormFarm.LandArea, result.LandArea)
	assert.Equal(t, gormFarm.UnitMeasure, result.UnitMeasure)
	assert.Equal(t, gormFarm.Address, result.Address)
	assert.Len(t, result.CropProductions, len(gormFarm.CropProductions))

	for i, crop := range gormFarm.CropProductions {
		assert.Equal(t, crop.CropType, result.CropProductions[i].CropType)
		assert.Equal(t, crop.IsIrrigated, result.CropProductions[i].IsIrrigated)
		assert.Equal(t, crop.IsInsured, result.CropProductions[i].IsInsured)
		assert.Equal(t, crop.ID, result.CropProductions[i].ID)
		assert.Equal(t, crop.FarmID, result.CropProductions[i].FarmID)
	}
}

func TestToDomainCropProductions(t *testing.T) {
	gormCrops := []entities.CropProduction{
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

	result := ToDomainCropProductions(gormCrops)

	assert.Len(t, result, len(gormCrops))
	for i, crop := range gormCrops {
		assert.Equal(t, crop.CropType, result[i].CropType)
		assert.Equal(t, crop.IsIrrigated, result[i].IsIrrigated)
		assert.Equal(t, crop.IsInsured, result[i].IsInsured)
		assert.Equal(t, crop.ID, result[i].ID)
		assert.Equal(t, crop.FarmID, result[i].FarmID)
	}
}
