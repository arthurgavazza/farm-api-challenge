package domain

import (
	"time"

	"github.com/google/uuid"
)

type Farm struct {
	ID              uuid.UUID        `json:"id"`
	Name            string           `json:"name"`
	LandArea        float64          `json:"land_area"`
	UnitMeasure     string           `json:"unit_measure"`
	Address         string           `json:"address"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	DeletedAt       *time.Time       `json:"deleted_at,omitempty"`
	CropProductions []CropProduction `json:"crop_productions"`
}

type FarmSearchParameters struct {
	CropType        *string  `json:"crop_type"`
	MinimumLandArea *float64 `json:"minimum_land_area"`
	MaximumLandArea *float64 `json:"maximum_land_area"`
	Page            int      `json:"page"`
	PerPage         int      `json:"per_page"`
}

func NewFarm(
	name string,
	landArea float64,
	unitMeasure string,
	address string,
	productions []CropProduction,
) (*Farm, error) {
	farm := &Farm{
		ID:              uuid.New(),
		Name:            name,
		LandArea:        landArea,
		UnitMeasure:     unitMeasure,
		Address:         address,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		CropProductions: productions,
	}

	return farm, nil
}
