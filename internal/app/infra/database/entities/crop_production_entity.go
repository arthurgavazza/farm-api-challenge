package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CropProduction struct {
	ID          uuid.UUID      `gorm:"primaryKey;default:gen_random_uuid()"`
	FarmID      uuid.UUID      `gorm:"not null"`
	CropType    string         `gorm:"size:50;not null"` // Enum-like values: RICE, BEANS, etc.
	IsIrrigated bool           `gorm:"not null"`
	IsInsured   bool           `gorm:"not null"`
	CreatedAt   time.Time      `gorm:"not null"`
	UpdatedAt   time.Time      `gorm:"not null"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
