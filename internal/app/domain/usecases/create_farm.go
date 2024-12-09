package usecases

import (
	"context"

	"github.com/arthurgavazza/farm-api-challenge/internal/app/domain"
	"github.com/google/uuid"
)

type CreateFarmUseCase interface {
	Execute(ctx context.Context, farm domain.Farm) (*domain.Farm, error)
}
type CreateFarm struct {
	repository domain.FarmRepository
}

func (uc *CreateFarm) Execute(ctx context.Context, farm domain.Farm) (*domain.Farm, error) {
	farmID := uuid.New()
	farm.ID = farmID

	for i := range farm.CropProductions {
		farm.CropProductions[i].ID = uuid.New()
		farm.CropProductions[i].FarmID = farmID
	}
	return uc.repository.CreateFarm(ctx, &farm)
}

func NewCreateFarmUseCase(repo domain.FarmRepository) *CreateFarm {
	return &CreateFarm{
		repository: repo,
	}
}
