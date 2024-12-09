package usecases

import (
	"context"

	"github.com/arthurgavazza/farm-api-challenge/internal/app/domain"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/models"
)

type ListFarmsUseCase interface {
	Execute(ctx context.Context, searchParameters *domain.FarmSearchParameters) (*models.PaginatedResponse[*domain.Farm], error)
}
type ListFarms struct {
	repository domain.FarmRepository
}

func (uc *ListFarms) Execute(ctx context.Context, searchParameters *domain.FarmSearchParameters) (*models.PaginatedResponse[*domain.Farm], error) {
	return uc.repository.ListFarms(ctx, searchParameters)
}

func NewListFarmsUseCase(repo domain.FarmRepository) *ListFarms {
	return &ListFarms{
		repository: repo,
	}
}
