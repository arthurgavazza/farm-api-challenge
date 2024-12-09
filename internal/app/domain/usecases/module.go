package usecases

import (
	"go.uber.org/fx"
)

var Module = fx.Provide(
	NewCreateFarmUseCase,
	NewListFarmsUseCase,
	fx.Annotate(
		NewCreateFarmUseCase,
		fx.As(new(CreateFarmUseCase)),
	),
	fx.Annotate(
		NewListFarmsUseCase,
		fx.As(new(ListFarmsUseCase)),
	),
)
