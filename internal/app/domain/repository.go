package domain

import "github.com/arthurgavazza/farm-api-challenge/internal/app/models"


type FarmRepository interface {
	CreateFarm(farm *Farm) (*Farm, error)
	ListFarms(searchParameters FarmSearchParameters) (*models.PaginatedResponse[Farm], error)
	DeleteFarm(farmId string) error
}