package repositories

import (
	"context"
	"time"

	"github.com/arthurgavazza/farm-api-challenge/internal/app/domain"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/database/entities"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/database/mappers"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/models"
	"github.com/google/uuid"
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

type farmWithCropProduction struct {
	FarmID      uuid.UUID  `json:"farm_id"`
	Name        string     `json:"name"`
	LandArea    float64    `json:"land_area"`
	UnitMeasure string     `json:"unit_measure"`
	Address     string     `json:"address"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`

	CropProductionID     uuid.UUID `json:"crop_production_id"`
	CropProductionFarmID uuid.UUID `json:"crop_production_farm_id"`
	CropType             string    `json:"crop_type"`
	IsIrrigated          bool      `json:"is_irrigated"`
	IsInsured            bool      `json:"is_insured"`
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

func (f *FarmRepository) parseRawFarmResults(rawResults []farmWithCropProduction) []*domain.Farm {
	farmsMap := make(map[uuid.UUID]*domain.Farm)

	for _, row := range rawResults {
		if farm, exists := farmsMap[row.FarmID]; exists {
			farm.CropProductions = append(farm.CropProductions, domain.CropProduction{
				ID:          row.CropProductionID,
				FarmID:      row.CropProductionFarmID,
				CropType:    row.CropType,
				IsIrrigated: row.IsIrrigated,
				IsInsured:   row.IsInsured,
			})
		} else {
			farmsMap[row.FarmID] = &domain.Farm{
				ID:          row.FarmID,
				Name:        row.Name,
				LandArea:    row.LandArea,
				UnitMeasure: row.UnitMeasure,
				Address:     row.Address,
				CreatedAt:   row.CreatedAt,
				UpdatedAt:   row.UpdatedAt,
				DeletedAt:   row.DeletedAt,
				CropProductions: []domain.CropProduction{
					{
						ID:          row.CropProductionID,
						FarmID:      row.CropProductionFarmID,
						CropType:    row.CropType,
						IsIrrigated: row.IsIrrigated,
						IsInsured:   row.IsInsured,
					},
				},
			}
		}
	}

	var domainFarms []*domain.Farm
	for _, farm := range farmsMap {
		domainFarms = append(domainFarms, farm)
	}
	return domainFarms
}

func (f *FarmRepository) ListFarms(ctx context.Context, searchParameters *domain.FarmSearchParameters) (*models.PaginatedResponse[*domain.Farm], error) {
	var rawResults []farmWithCropProduction
	var totalCount int64

	query := f.db.Model(&entities.Farm{}).Joins("JOIN crop_productions ON crop_productions.farm_id = farms.id")

	if searchParameters.CropType != nil {
		query = query.Where("crop_productions.crop_type = ?", *searchParameters.CropType)
	}

	if searchParameters.MinimumLandArea != nil && searchParameters.MaximumLandArea != nil {
		query = query.Where("farms.land_area BETWEEN ? AND ?", *searchParameters.MinimumLandArea, *searchParameters.MaximumLandArea)
	} else if searchParameters.MinimumLandArea != nil {
		query = query.Where("farms.land_area >= ?", *searchParameters.MinimumLandArea)
	} else if searchParameters.MaximumLandArea != nil {
		query = query.Where("farms.land_area <= ?", *searchParameters.MaximumLandArea)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	offset := (searchParameters.Page - 1) * searchParameters.PerPage
	if searchParameters.Page < 1 {
		offset = 0
	}
	if searchParameters.PerPage < 1 {
		searchParameters.PerPage = 10
	}
	queryFields := `farms.id AS farm_id, farms.name, farms.land_area, farms.unit_measure, farms.address, farms.created_at, farms.updated_at, farms.deleted_at,
        crop_productions.id AS crop_production_id, crop_productions.farm_id AS crop_production_farm_id, crop_productions.crop_type, crop_productions.is_irrigated, crop_productions.is_insured`

	if err := query.Select(queryFields).
		Offset(offset).Limit(searchParameters.PerPage).Find(&rawResults).Error; err != nil {
		return nil, err
	}

	domainFarms := f.parseRawFarmResults(rawResults)
	response := &models.PaginatedResponse[*domain.Farm]{
		Items:       domainFarms,
		TotalCount:  totalCount,
		CurrentPage: searchParameters.Page,
		PerPage:     searchParameters.PerPage,
	}

	return response, nil
}
