package domain

import "context"

type FarmRepository interface {
	CreateFarm(ctx context.Context,farm *Farm) (*Farm, error)
	// ListFarms(searchParameters FarmSearchParameters) (*models.PaginatedResponse[Farm], error)
	// DeleteFarm(farmId string) error
}
