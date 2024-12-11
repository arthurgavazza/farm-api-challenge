package integration_tests

import (
	"context"
	"log"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/arthurgavazza/farm-api-challenge/internal/app/domain"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/config"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/database"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/infra/database/repositories"
	"github.com/arthurgavazza/farm-api-challenge/testutils"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	logger "github.com/arthurgavazza/farm-api-challenge/internal/app/shared/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type IntegrationTestsSuite struct {
	suite.Suite
	postgresContainer *postgres.PostgresContainer
	repo              *repositories.FarmRepository
}

func (is *IntegrationTestsSuite) SetupSuite() {
	uuid.EnableRandPool()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalf("unable to get caller info")
	}
	projectRoot := filepath.Join(filepath.Dir(filename), "../")
	log.Println(projectRoot)
	envFilePath := filepath.Join(projectRoot, ".env.test")
	err := godotenv.Load(envFilePath)
	if err != nil {
		log.Fatalf("unable to load env file")
	}
	ctx := context.Background()
	configuration := config.NewConfig()
	postgresContainer, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase(configuration.Database.Name),
		postgres.WithUsername(configuration.Database.User),
		postgres.WithPassword(configuration.Database.Password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}
	connectionString, err := postgresContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("Failed to get connection string %s", err)
	}

	is.postgresContainer = postgresContainer
	db := database.NewPostgresDatabase(connectionString)
	repo := repositories.NewFarmRepository(db, logger.NewLogger())
	is.repo = repo

}

func (is *IntegrationTestsSuite) TestCreateFarm() {
	ctx := context.Background()
	farm := testutils.GenerateFakeFarm(nil, nil)
	defer is.repo.DeleteFarm(ctx, farm.ID.String())
	createdFarm, err := is.repo.CreateFarm(ctx, farm)

	assert.NoError(is.T(), err)
	assert.NotNil(is.T(), createdFarm.ID)
	assert.Equal(is.T(), farm.Name, createdFarm.Name)
	assert.Equal(is.T(), farm.LandArea, createdFarm.LandArea)
}

func (is *IntegrationTestsSuite) TestListFarms() {
	ctx := context.Background()
	var farms []*domain.Farm
	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
		for _, farm := range farms {
			is.repo.DeleteFarm(ctx, farm.ID.String())
		}
	}()

	cofeeCropType := domain.CropTypeCoffee.String()
	riceCropType := domain.CropTypeRice.String()
	cornCropType := domain.CropTypeCorn.String()
	landArea := float64(100)
	farmsWithCoffeeProductions := testutils.GenerateFarms(2, &cofeeCropType, &landArea)
	farmsWithRiceProductions := testutils.GenerateFarms(6, &riceCropType, &landArea)
	biggerLandArea := 15 * landArea
	farmsWithBiggerLandAreas := testutils.GenerateFarms(3, &cornCropType, &biggerLandArea)
	farms = append(farms, farmsWithCoffeeProductions...)
	farms = append(farms, farmsWithRiceProductions...)
	farms = append(farms, farmsWithBiggerLandAreas...)

	for _, farm := range farms {
		wg.Add(1)
		go func(farm *domain.Farm) {
			defer wg.Done()
			_, err := is.repo.CreateFarm(ctx, farm)
			if err != nil {
				is.T().Errorf("failed to insert farm: %v", err)
			}
		}(farm)
	}
	wg.Wait()
	searchParams := &domain.FarmSearchParameters{
		Page:     1,
		PerPage:  3,
		CropType: &cofeeCropType,
	}
	coffeeCrops, err := is.repo.ListFarms(ctx, searchParams)
	assert.NoError(is.T(), err)
	// since there are two farms with coffee production and the provided PerPage argument = 3, the returned list should have 2 items
	assert.NotNil(is.T(), coffeeCrops)
	assert.Equal(is.T(), int64(len(farmsWithCoffeeProductions)), coffeeCrops.TotalCount)
	assert.Equal(is.T(), len(farmsWithCoffeeProductions), len(coffeeCrops.Items))

	// since there are 6 farms with rice production and the provided PerPage argument = 3, the returned list should have 3 items
	searchParams.CropType = &riceCropType
	riceCrops, err := is.repo.ListFarms(ctx, searchParams)
	assert.NoError(is.T(), err)
	assert.NotNil(is.T(), riceCrops)
	assert.Equal(is.T(), searchParams.PerPage, len(riceCrops.Items))
	assert.Equal(is.T(), int64(len(farmsWithRiceProductions)), riceCrops.TotalCount)

	// since there are 3 farms with land area between 1000 and 2000 the returned result should reflect this number
	searchParams.CropType = nil
	searchParams.MinimumLandArea = testutils.PointerTo(float64(1000))
	searchParams.MaximumLandArea = testutils.PointerTo(float64(2000))

	biggerLandAreaFarms, err := is.repo.ListFarms(ctx, searchParams)
	assert.NoError(is.T(), err)
	assert.NotNil(is.T(), biggerLandAreaFarms)
	assert.Equal(is.T(), len(farmsWithBiggerLandAreas), len(biggerLandAreaFarms.Items))

}

func (is *IntegrationTestsSuite) TestDeleteFarm() {
	ctx := context.Background()
	farm := testutils.GenerateFakeFarm(nil, nil)
	defer is.repo.DeleteFarm(ctx, farm.ID.String())
	createdFarm, err := is.repo.CreateFarm(ctx, farm)
	require.NoError(is.T(), err)

	err = is.repo.DeleteFarm(ctx, createdFarm.ID.String())
	assert.NoError(is.T(), err)

	err = is.repo.DeleteFarm(ctx, createdFarm.ID.String())
	assert.Error(is.T(), err)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestsSuite))
}

func (is *IntegrationTestsSuite) TearDownSuite() {
	if err := testcontainers.TerminateContainer(is.postgresContainer); err != nil {
		log.Printf("failed to terminate container: %s", err)
	}
}
