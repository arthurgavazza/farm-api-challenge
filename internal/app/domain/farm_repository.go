package domain

import (
	"context"

	"github.com/arthurgavazza/farm-api-challenge/internal/app/models"
)

type FarmRepository interface {
	CreateFarm(ctx context.Context, farm *Farm) (*Farm, error)
	ListFarms(ctx context.Context, searchParameters *FarmSearchParameters) (*models.PaginatedResponse[*Farm], error)
	// DeleteFarm(farmId string) error
}
