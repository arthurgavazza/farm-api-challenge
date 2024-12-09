package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/arthurgavazza/farm-api-challenge/internal/app/domain"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
)

type MockFarmRepository struct {
	mock.Mock
}

func (m *MockFarmRepository) ListFarms(ctx context.Context, searchParameters *domain.FarmSearchParameters) (*models.PaginatedResponse[*domain.Farm], error) {
	panic("unimplemented")
}

func (m *MockFarmRepository) CreateFarm(ctx context.Context, farm *domain.Farm) (*domain.Farm, error) {
	args := m.Called(ctx, farm)
	return args.Get(0).(*domain.Farm), args.Error(1)
}

func TestCreateFarmSuccess(t *testing.T) {
	mockRepo := new(MockFarmRepository)
	useCase := NewCreateFarmUseCase(mockRepo)

	ctx := context.Background()
	farm := domain.Farm{
		Name:        "Test Farm",
		LandArea:    100.5,
		UnitMeasure: "acres",
		Address:     "123 Farm Lane",
		Productions: []domain.CropProduction{
			{CropType: "RICE"},
		},
	}

	expectedFarm := farm
	expectedFarm.ID = uuid.New()
	expectedFarm.Productions[0].ID = uuid.New()
	expectedFarm.Productions[0].FarmID = expectedFarm.ID

	mockRepo.On("CreateFarm", ctx, mock.MatchedBy(func(f *domain.Farm) bool {
		return f.Name == farm.Name && f.LandArea == farm.LandArea && len(f.Productions) == 1
	})).Return(&expectedFarm, nil)

	result, err := useCase.Execute(ctx, farm)

	assert.NoError(t, err)
	assert.Equal(t, expectedFarm.ID, result.ID)
	assert.Equal(t, expectedFarm.Productions[0].ID, result.Productions[0].ID)
	mockRepo.AssertExpectations(t)
}

func TestCreateFarmRepositoryError(t *testing.T) {
	mockRepo := new(MockFarmRepository)
	useCase := NewCreateFarmUseCase(mockRepo)

	ctx := context.Background()
	farm := domain.Farm{
		Name:        "Test Farm",
		LandArea:    100.5,
		UnitMeasure: "acres",
		Address:     "123 Farm Lane",
		Productions: []domain.CropProduction{
			{CropType: "RICE"},
		},
	}

	mockRepo.On("CreateFarm", ctx, mock.Anything).Return((*domain.Farm)(nil), errors.New("database error"))

	result, err := useCase.Execute(ctx, farm)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "database error")
	mockRepo.AssertExpectations(t)
}
