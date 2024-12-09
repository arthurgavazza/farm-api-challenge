package repositories

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
		CropProductions: []domain.CropProduction{
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
	for i, expectedCropProduction := range rs.farm.CropProductions {
		assert.Equal(rs.T(), expectedCropProduction.CropType, farm.CropProductions[i].CropType)
		assert.Equal(rs.T(), expectedCropProduction.IsIrrigated, farm.CropProductions[i].IsIrrigated)
		assert.Equal(rs.T(), expectedCropProduction.IsInsured, farm.CropProductions[i].IsInsured)
	}

}

func (rs *FarmRepositoryTestSuite) TestListFarmsEmptyResults() {
	rs.mock.ExpectQuery(`SELECT count(.*) FROM "farms"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	rs.mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WillReturnRows(sqlmock.NewRows([]string{
			"farm_id", "name", "land_area", "unit_measure", "address", "created_at", "updated_at", "deleted_at",
			"crop_production_id", "crop_production_farm_id", "crop_type", "is_irrigated", "is_insured",
		}))

	searchParams := &domain.FarmSearchParameters{
		Page:    1,
		PerPage: 10,
	}
	response, err := rs.repo.ListFarms(context.Background(), searchParams)

	assert.NoError(rs.T(), err)
	assert.NotNil(rs.T(), response)
	assert.Equal(rs.T(), 0, len(response.Items))
	assert.Equal(rs.T(), int64(0), response.TotalCount)
	assert.Equal(rs.T(), 1, response.CurrentPage)
	assert.Equal(rs.T(), 10, response.PerPage)
}

func (rs *FarmRepositoryTestSuite) TestListFarmsWithMultipleCropProductions() {
	rs.mock.ExpectQuery(`SELECT count(.*) FROM "farms"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	rows := sqlmock.NewRows([]string{
		"farm_id", "name", "land_area", "unit_measure", "address", "created_at", "updated_at", "deleted_at",
		"crop_production_id", "crop_production_farm_id", "crop_type", "is_irrigated", "is_insured",
	}).AddRow(
		rs.farm.ID, rs.farm.Name, rs.farm.LandArea, rs.farm.UnitMeasure, rs.farm.Address, rs.farm.CreatedAt, rs.farm.UpdatedAt, nil,
		rs.farm.CropProductions[0].ID, rs.farm.ID, rs.farm.CropProductions[0].CropType, rs.farm.CropProductions[0].IsIrrigated, rs.farm.CropProductions[0].IsInsured,
	).AddRow(
		rs.farm.ID, rs.farm.Name, rs.farm.LandArea, rs.farm.UnitMeasure, rs.farm.Address, rs.farm.CreatedAt, rs.farm.UpdatedAt, nil,
		rs.farm.CropProductions[1].ID, rs.farm.ID, rs.farm.CropProductions[1].CropType, rs.farm.CropProductions[1].IsIrrigated, rs.farm.CropProductions[1].IsInsured,
	)
	rs.mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WillReturnRows(rows)

	searchParams := &domain.FarmSearchParameters{
		Page:    1,
		PerPage: 10,
	}
	response, err := rs.repo.ListFarms(context.Background(), searchParams)

	assert.NoError(rs.T(), err)
	assert.NotNil(rs.T(), response)
	assert.Equal(rs.T(), 1, len(response.Items))
	assert.Equal(rs.T(), rs.farm.ID, response.Items[0].ID)
	assert.Equal(rs.T(), 2, len(response.Items[0].CropProductions))
}

func (rs *FarmRepositoryTestSuite) TestListFarmsWithCropTypeFilter() {

	rs.mock.ExpectQuery(`SELECT count(.*) FROM "farms"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	rows := sqlmock.NewRows([]string{
		"farm_id", "name", "land_area", "unit_measure", "address", "created_at", "updated_at", "deleted_at",
		"crop_production_id", "crop_production_farm_id", "crop_type", "is_irrigated", "is_insured",
	}).AddRow(
		rs.farm.ID, rs.farm.Name, rs.farm.LandArea, rs.farm.UnitMeasure, rs.farm.Address, rs.farm.CreatedAt, rs.farm.UpdatedAt, nil,
		rs.farm.CropProductions[0].ID, rs.farm.ID, rs.farm.CropProductions[0].CropType, rs.farm.CropProductions[0].IsIrrigated, rs.farm.CropProductions[0].IsInsured,
	)

	perPage := 10
	rs.mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WithArgs(domain.CropTypeCoffee, perPage).
		WillReturnRows(rows)

	searchParams := &domain.FarmSearchParameters{
		Page:     1,
		PerPage:  perPage,
		CropType: testutils.PointerTo(domain.CropTypeCoffee.String()),
	}
	response, err := rs.repo.ListFarms(context.Background(), searchParams)

	// Assertions
	assert.NoError(rs.T(), err)
	assert.NotNil(rs.T(), response)
	assert.Equal(rs.T(), 1, len(response.Items)) // One farm
	assert.Equal(rs.T(), rs.farm.ID, response.Items[0].ID)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(FarmRepositoryTestSuite))
}
