package repositories

import (
	"context"
	"time"

	"github.com/arthurgavazza/farm-api-challenge/internal/app/domain"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/database/entities"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/database/mappers"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/models"
	shared "github.com/arthurgavazza/farm-api-challenge/internal/app/shared/errors"
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
	var farmIDs []string
	var rawResults []farmWithCropProduction
	var totalCount int64

	baseQuery := f.db.Model(&entities.Farm{}).Joins("JOIN crop_productions ON crop_productions.farm_id = farms.id")

	if searchParameters.CropType != nil {
		baseQuery = baseQuery.Where("crop_productions.crop_type = ?", *searchParameters.CropType)
	}

	if searchParameters.MinimumLandArea != nil && searchParameters.MaximumLandArea != nil {
		baseQuery = baseQuery.Where("farms.land_area BETWEEN ? AND ?", *searchParameters.MinimumLandArea, *searchParameters.MaximumLandArea)
	} else if searchParameters.MinimumLandArea != nil {
		baseQuery = baseQuery.Where("farms.land_area >= ?", *searchParameters.MinimumLandArea)
	} else if searchParameters.MaximumLandArea != nil {
		baseQuery = baseQuery.Where("farms.land_area <= ?", *searchParameters.MaximumLandArea)
	}

	// Count distinct farms
	if err := baseQuery.Distinct("farms.id").Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Calculate pagination
	offset := (searchParameters.Page - 1) * searchParameters.PerPage
	if searchParameters.Page < 1 {
		offset = 0
	}
	if searchParameters.PerPage < 1 {
		searchParameters.PerPage = 10
	}

	if err := baseQuery.
		Select("farms.id").
		Group("farms.id").
		Offset(offset).
		Limit(searchParameters.PerPage).
		Pluck("farms.id", &farmIDs).Error; err != nil {
		return nil, err
	}

	// Fetch full details for the selected farm IDs
	if err := f.db.Model(&entities.Farm{}).
		Joins("JOIN crop_productions ON crop_productions.farm_id = farms.id").
		Where("farms.id IN ?", farmIDs).
		Select(`farms.id AS farm_id, farms.name, farms.land_area, farms.unit_measure, farms.address, farms.created_at, farms.updated_at, farms.deleted_at,
                crop_productions.id AS crop_production_id, crop_productions.farm_id AS crop_production_farm_id, crop_productions.crop_type, crop_productions.is_irrigated, crop_productions.is_insured`).
		Find(&rawResults).Error; err != nil {
		return nil, err
	}

	// Parse results and create the response
	domainFarms := f.parseRawFarmResults(rawResults)
	response := &models.PaginatedResponse[*domain.Farm]{
		Items:       domainFarms,
		TotalCount:  totalCount,
		CurrentPage: searchParameters.Page,
		PerPage:     searchParameters.PerPage,
	}

	return response, nil
}

func (f *FarmRepository) DeleteFarm(ctx context.Context, farmId string) error {
	tx := f.db.Delete(&entities.Farm{}, "id = ?", farmId)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return &shared.NotFoundError{
			Resource: "Farm",
			ID:       farmId,
		}
	}

	return nil
}
