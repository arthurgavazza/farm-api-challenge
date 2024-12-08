package repositories

import (
	"github.com/arthurgavazza/farm-api-challenge/internal/app/domain"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(
			NewFarmRepository,
			fx.As(new(domain.FarmRepository)),
		),
	),
)
