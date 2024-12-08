package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Farm struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	LandArea    float64          `json:"land_area"`
	UnitMeasure string           `json:"unit_measure"`
	Address     string           `json:"address"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	DeletedAt   *time.Time       `json:"deleted_at,omitempty"` // nullable field, use *time.Time
	Productions []CropProduction `json:"productions"`
}

type FarmSearchParameters struct {
	CropType string  `json:"crop_type"`
	LandArea float64 `json:"land_area"`
	Page     int     `json:"page"`
	PerPage  int     `json:"per_page"`
}

func NewFarm(
	name string,
	landArea float64,
	unitMeasure string,
	address string,
	productions []CropProduction,
) (*Farm, error) {
	if name == "" {
		return nil, errors.New("farm name cannot be empty")
	}
	if landArea <= 0 {
		return nil, errors.New("land area must be greater than zero")
	}
	if unitMeasure == "" {
		return nil, errors.New("unit measure cannot be empty")
	}
	if address == "" {
		return nil, errors.New("address cannot be empty")
	}

	farm := &Farm{
		ID:          uuid.New(),
		Name:        name,
		LandArea:    landArea,
		UnitMeasure: unitMeasure,
		Address:     address,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Productions: productions,
	}

	return farm, nil
}
