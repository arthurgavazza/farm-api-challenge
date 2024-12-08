package database

import (
	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/database/repositories"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewPostgresDatabase,
	),
	repositories.Module,
)
