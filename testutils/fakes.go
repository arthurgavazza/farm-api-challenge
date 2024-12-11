package testutils

import (
	"time"

	"math/rand"

	"github.com/arthurgavazza/farm-api-challenge/internal/app/domain"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

func GenerateFakeFarm(cropType *string, landArea *float64) *domain.Farm {
	var area float64
	if landArea != nil {
		area = *landArea
	} else {
		randomValue, err := faker.RandomInt(100, 800)
		if err != nil {
			area = 102.5
		}
		area = float64(randomValue[0])
	}
	farm := &domain.Farm{
		ID:          uuid.New(),
		Name:        faker.Name(),
		LandArea:    area,
		UnitMeasure: "hectares",
		Address:     "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	cropTypeValue := cropType
	if cropTypeValue == nil {
		generatedCropType := generateRandomCropType()
		cropTypeValue = &generatedCropType
	}
	cropProduction := domain.CropProduction{
		ID:          uuid.New(),
		FarmID:      farm.ID,
		CropType:    *cropTypeValue,
		IsIrrigated: true,
		IsInsured:   true,
	}

	farm.CropProductions = append(farm.CropProductions, cropProduction)

	return farm
}

func generateRandomCropType() string {
	cropTypes := []domain.CropType{domain.CropTypeCoffee, domain.CropTypeRice, domain.CropTypeSoybean, domain.CropTypeCorn}
	return cropTypes[rand.Intn(len(cropTypes))].String()
}

func GenerateFarms(total int, cropType *string, landArea *float64) []*domain.Farm {
	var farms []*domain.Farm
	for i := 0; i < total; i++ {
		farm := GenerateFakeFarm(cropType, landArea)
		farms = append(farms, farm)
	}
	return farms
}
