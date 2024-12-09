package testutils

import (
	"time"

	"math/rand"

	"github.com/arthurgavazza/farm-api-challenge/internal/app/domain"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

func GenerateFakeFarm(cropType *string) *domain.Farm {
	farm := &domain.Farm{
		ID:          uuid.New(),
		Name:        faker.Name(),
		LandArea:    102.5,
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

func GenerateFarms(total int, cropType *string) []*domain.Farm {
	var farms []*domain.Farm
	for i := 0; i < total; i++ {
		farm := GenerateFakeFarm(cropType)
		farms = append(farms, farm)
	}
	return farms
}
