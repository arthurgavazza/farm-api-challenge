package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Farm struct {
	ID          uuid.UUID        `gorm:"primaryKey"`
	Name        string           `gorm:"size:255;not null"`
	LandArea    float64          `gorm:"not null"`
	UnitMeasure string           `gorm:"size:50;not null"`
	Address     string           `gorm:"size:255;not null"`
	Productions []CropProduction `gorm:"foreignKey:FarmID;constraint:OnDelete:CASCADE;"`
	CreatedAt   time.Time        `gorm:"not null"`
	UpdatedAt   time.Time        `gorm:"not null"`
	DeletedAt   gorm.DeletedAt   `gorm:"index"`
}
