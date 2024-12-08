package domain

import (
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
