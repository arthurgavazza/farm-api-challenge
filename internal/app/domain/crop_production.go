package domain

import (
	"errors"

	"github.com/google/uuid"
)

type CropProduction struct {
	ID          uuid.UUID `json:"id"`
	FarmID      uuid.UUID `json:"farm_id"`
	CropType    string    `json:"crop_type"`
	IsIrrigated bool      `json:"is_irrigated"`
	IsInsured   bool      `json:"is_insured"`
}

type CropType string

const (
	CropTypeRice    CropType = "RICE"
	CropTypeCorn    CropType = "CORN"
	CropTypeSoybean CropType = "SOYBEANS"
	CropTypeCoffee  CropType = "COFFEE"
)

func (c CropType) IsValid() bool {
	switch c {
	case CropTypeRice, CropTypeCorn, CropTypeSoybean, CropTypeCoffee:
		return true
	default:
		return false
	}
}

func (c CropType) String() string {
	return string(c)
}

var (
	ErrInvalidCropType = errors.New("invalid crop type")
	ErrInvalidFarmID   = errors.New("invalid farm ID")
)

func NewCropProduction(
	id uuid.UUID,
	farmId uuid.UUID,
	cropType CropType,
	isIrrigated bool,
	isInsured bool,
) (*CropProduction, error) {
	if farmId == uuid.Nil {
		return nil, ErrInvalidFarmID
	}
	if !cropType.IsValid() {
		return nil, ErrInvalidCropType
	}

	return &CropProduction{
		ID:          id,
		FarmID:      farmId,
		CropType:    cropType.String(),
		IsIrrigated: isIrrigated,
		IsInsured:   isInsured,
	}, nil
}
