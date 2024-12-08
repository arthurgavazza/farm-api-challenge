package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Farm struct {
	ID          uuid.UUID        `gorm:"primaryKey" json:"id"`
	Name        string           `gorm:"size:255;not null" json:"name"`
	LandArea    float64          `gorm:"not null" json:"land_area"`
	UnitMeasure string           `gorm:"size:50;not null" json:"unit_measure"`
	Address     string           `gorm:"size:255;not null" json:"address"`
	Productions []CropProduction `gorm:"foreignKey:FarmID;constraint:OnDelete:CASCADE;" json:"productions"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	DeletedAt   gorm.DeletedAt   `gorm:"index" json:"deleted_at,omitempty"`
}