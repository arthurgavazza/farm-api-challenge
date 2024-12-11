package database

import (
	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/config"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/database/repositories"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var Module = fx.Options(
	fx.Provide(func(config *config.Config) *gorm.DB {
		connectionString := generateConnectionString(config)
		return NewPostgresDatabase(connectionString)
	}),
	repositories.Module,
)
