package database

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/domain"
	"github.com/arthurgavazza/farm-api-challenge/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type FarmRepositoryTestSuite struct {
	suite.Suite
	conn *sql.DB
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	repo *FarmRepository
	farm *domain.Farm
}

func (rs *FarmRepositoryTestSuite) SetupSuite() {
	var (
		err error
	)

	rs.conn, rs.mock, err = sqlmock.New()
	assert.NoError(rs.T(), err)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 rs.conn,
		PreferSimpleProtocol: true,
	})

	rs.DB, err = gorm.Open(dialector, &gorm.Config{})
	assert.NoError(rs.T(), err)

	rs.repo = NewFarmRepository(rs.DB)
	assert.IsType(rs.T(), &FarmRepository{}, rs.repo)
	farmId := uuid.New()
	coffeeCrop, err := domain.NewCropProduction(uuid.New(), farmId, domain.CropTypeCoffee, true, true)
	if err != nil {
		rs.T().Error(err)
	}
	riceCrop, err := domain.NewCropProduction(uuid.New(), farmId, domain.CropTypeRice, true, true)
	if err != nil {
		rs.T().Error(err)
	}
	rs.farm = &domain.Farm{
		ID:          farmId,
		Name:        "Test Farm",
		LandArea:    100,
		UnitMeasure: "acre",
		Address:     "Test Address",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Productions: []domain.CropProduction{
			*coffeeCrop,
			*riceCrop,
		},
	}
}

func (rs *FarmRepositoryTestSuite) AfterTest(_, _ string) {
	assert.NoError(rs.T(), rs.mock.ExpectationsWereMet())
}

func (rs *FarmRepositoryTestSuite) TestCreateFarm() {
	rs.mock.ExpectBegin()
	rs.mock.ExpectExec(
		regexp.QuoteMeta(`INSERT INTO "farms" ("id","name","land_area","unit_measure","address","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`)).
		WithArgs(
			rs.farm.ID,
			rs.farm.Name,
			rs.farm.LandArea,
			rs.farm.UnitMeasure,
			rs.farm.Address,
			testutils.AnyTime{},
			testutils.AnyTime{},
			nil,
		).WillReturnResult(sqlmock.NewResult(1, 1))
	rs.mock.ExpectExec(`INSERT INTO "crop_productions"`).WillReturnResult(sqlmock.NewResult(2, 2))
	rs.mock.ExpectCommit()

	farm, err := rs.repo.CreateFarm(context.Background(), rs.farm)
	assert.NoError(rs.T(), err)
	assert.NotNil(rs.T(), farm.ID)
	assert.Equal(rs.T(), rs.farm.ID, farm.ID)
	assert.Equal(rs.T(), rs.farm.Name, farm.Name)
	assert.Equal(rs.T(), rs.farm.LandArea, farm.LandArea)
	assert.Equal(rs.T(), rs.farm.UnitMeasure, farm.UnitMeasure)
	assert.Equal(rs.T(), rs.farm.Address, farm.Address)
	for i, expectedCropProduction := range rs.farm.Productions {
		assert.Equal(rs.T(), expectedCropProduction.CropType, farm.Productions[i].CropType)
		assert.Equal(rs.T(), expectedCropProduction.IsIrrigated, farm.Productions[i].IsIrrigated)
		assert.Equal(rs.T(), expectedCropProduction.IsInsured, farm.Productions[i].IsInsured)
	}

}

func TestSuite(t *testing.T) {
	suite.Run(t, new(FarmRepositoryTestSuite))
}
