package domain

import "time"

type Farm struct {
    ID          uint                `json:"id"`
    Name        string              `json:"name"`
    LandArea    float64             `json:"land_area"`
    UnitMeasure string              `json:"unit_measure"`
    Address     string              `json:"address"`
    CreatedAt   time.Time           `json:"created_at"`
    UpdatedAt   time.Time           `json:"updated_at"`
    DeletedAt   *time.Time          `json:"deleted_at,omitempty"` // nullable field, use *time.Time
    Productions []CropProduction    `json:"productions"`
}


type CropProduction struct {
    ID          uint        `json:"id"`
    FarmID      uint        `json:"farm_id"`
    CropType    string      `json:"crop_type"`
    IsIrrigated bool        `json:"is_irrigated"`
    IsInsured   bool        `json:"is_insured"`
}

type FarmSearchParameters struct {
	CropType string     `json:"crop_type"`
	LandArea float64    `json:"land_area"`
	Page    int         `json:"page"`
	PerPage int         `json:"per_page"`
}