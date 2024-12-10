package usecases

import (
	"context"

	"github.com/arthurgavazza/farm-api-challenge/internal/app/domain"
)

type DeleteFarmUseCase interface {
	Execute(ctx context.Context, farmId string) error
}
type DeleteFarm struct {
	repository domain.FarmRepository
}

func (uc *DeleteFarm) Execute(ctx context.Context, farmId string) error {
	return uc.repository.DeleteFarm(ctx, farmId)
}

func NewDeleteFarmUseCase(repo domain.FarmRepository) *DeleteFarm {
	return &DeleteFarm{
		repository: repo,
	}
}
