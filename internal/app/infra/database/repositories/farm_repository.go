package repositories

import (
	"context"

	"github.com/arthurgavazza/farm-api-challenge/internal/app/domain"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/database/entities"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/database/mappers"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/models"
	"gorm.io/gorm"
)

type FarmRepository struct {
	db *gorm.DB
}

func NewFarmRepository(db *gorm.DB) *FarmRepository {
	return &FarmRepository{
		db: db,
	}
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
	farm.CreatedAt = ormFarm.CreatedAt
	farm.UpdatedAt = ormFarm.UpdatedAt
	return farm, nil
}

func (f *FarmRepository) ListFarms(ctx context.Context, searchParameters *domain.FarmSearchParameters) (*models.PaginatedResponse[*domain.Farm], error) {
	var farms []entities.Farm
	var totalCount int64

	query := f.db.Model(&entities.Farm{}).Preload("Productions")

	if searchParameters.CropType != nil {
		query = query.Joins("JOIN crop_productions ON crop_productions.farm_id = farms.id").
			Where("crop_productions.crop_type ILIKE ?", "%"+*searchParameters.CropType+"%")
	}
	if searchParameters.LandArea != nil {
		query = query.Where("land_area = ?", *searchParameters.LandArea)
	}

	// Count total items for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply pagination
	offset := (searchParameters.Page - 1) * searchParameters.PerPage
	if searchParameters.Page < 1 {
		offset = 0
	}
	if searchParameters.PerPage < 1 {
		searchParameters.PerPage = 10 // Default items per page
	}

	if err := query.Offset(offset).Limit(searchParameters.PerPage).Find(&farms).Error; err != nil {
		return nil, err
	}

	var domainFarms []*domain.Farm
	for _, farm := range farms {
		domainFarms = append(domainFarms, mappers.ToDomainFarm(&farm))
	}

	response := &models.PaginatedResponse[*domain.Farm]{
		Items:       domainFarms,
		TotalCount:  totalCount,
		CurrentPage: searchParameters.Page,
		PerPage:     searchParameters.PerPage,
	}

	return response, nil
}
