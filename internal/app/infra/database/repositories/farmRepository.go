package database

import (
	"context"

	"github.com/arthurgavazza/farm-api-challenge/internal/app/domain"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/database/mappers"
	"gorm.io/gorm"
)

type FarmRepository struct {
	db  *gorm.DB
}

func (f *FarmRepository) CreateFarm(ctx context.Context, farm *domain.Farm) (*domain.Farm, error) {
	ormFarm := mappers.ToGormFarm(farm)
    err := f.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(&ormFarm).Error; err != nil {
            return err
        }
        return nil
    })
    if err != nil {
        return nil, err
    }
	farm.ID = ormFarm.ID
	farm.CreatedAt = ormFarm.CreatedAt
	farm.UpdatedAt = ormFarm.UpdatedAt
    return farm, nil
}

