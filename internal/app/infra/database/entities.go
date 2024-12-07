package database

import "time"

type Farm struct {
    ID          uint            `gorm:"primaryKey" json:"id"`
    Name        string          `gorm:"size:255;not null" json:"name"`
    LandArea    float64         `gorm:"not null" json:"land_area"`
    UnitMeasure string          `gorm:"size:50;not null" json:"unit_measure"`
    Address     string          `gorm:"size:255;not null" json:"address"`
    CreatedAt   time.Time       `json:"created_at"`
    UpdatedAt   time.Time       `json:"updated_at"`
    DeletedAt   *time.Time      `json:"deleted_at,omitempty"` // nullable field, use *time.Time
    Productions []CropProduction `gorm:"foreignKey:FarmID;constraint:OnDelete:CASCADE;" json:"productions"`
}

type CropProduction struct {
    ID         uint   `gorm:"primaryKey" json:"id"`
    FarmID     uint   `gorm:"not null" json:"farm_id"`
    CropType   string `gorm:"size:50;not null" json:"crop_type"` // Enum-like values: RICE, BEANS, etc.
    IsIrrigated bool  `gorm:"not null" json:"is_irrigated"`
    IsInsured   bool  `gorm:"not null" json:"is_insured"`
}
