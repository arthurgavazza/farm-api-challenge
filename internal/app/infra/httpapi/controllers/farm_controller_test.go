package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/arthurgavazza/farm-api-challenge/internal/app/domain"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/dto"
	"github.com/arthurgavazza/farm-api-challenge/internal/app/models"
	shared "github.com/arthurgavazza/farm-api-challenge/internal/app/shared/errors"
	logger "github.com/arthurgavazza/farm-api-challenge/internal/app/shared/logger"
	"github.com/arthurgavazza/farm-api-challenge/testutils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockCreateFarmUseCase struct {
	mock.Mock
}

func (m *MockCreateFarmUseCase) Execute(ctx context.Context, farm domain.Farm) (*domain.Farm, error) {
	args := m.Called(ctx, farm)
	return args.Get(0).(*domain.Farm), args.Error(1)
}

type MockListFarmsUseCase struct {
	mock.Mock
}

func (m *MockListFarmsUseCase) Execute(ctx context.Context, searchParameters *domain.FarmSearchParameters) (*models.PaginatedResponse[*domain.Farm], error) {
	args := m.Called(ctx, searchParameters)
	return args.Get(0).(*models.PaginatedResponse[*domain.Farm]), args.Error(1)
}

type MockDeleteFarmUseCase struct {
	mock.Mock
}

func (m *MockDeleteFarmUseCase) Execute(ctx context.Context, farmId string) error {
	args := m.Called(ctx, farmId)
	return args.Error(0)
}

type FarmControllerTestSuite struct {
	suite.Suite
	logger *logger.Logger
}

func (cs *FarmControllerTestSuite) SetupSuite() {
	cs.logger = logger.NewLogger()
}

func (cs *FarmControllerTestSuite) TestFarmControllerCreateFarm() {
	tests := []struct {
		name               string
		inputDTO           dto.CreateFarmDTO
		expectedStatusCode int
		mockResponse       *domain.Farm
		mockError          error
		mockRequired       bool
	}{
		{
			name: "Successful Farm Creation",
			inputDTO: dto.CreateFarmDTO{
				Name:            "Test Farm",
				LandArea:        100.5,
				UnitMeasure:     "hectares",
				Address:         "123 Farm Lane",
				CropProductions: []dto.CropProductionDTO{},
			},
			expectedStatusCode: fiber.StatusCreated,
			mockResponse: &domain.Farm{
				ID:          uuid.New(),
				Name:        "Test Farm",
				LandArea:    100.5,
				UnitMeasure: "hectares",
				Address:     "123 Farm Lane",
			},
			mockError:    nil,
			mockRequired: true,
		},
		{
			name: "Bad Request - Invalid Crop Type",
			inputDTO: dto.CreateFarmDTO{
				Name:        "Test Farm",
				LandArea:    100.5,
				UnitMeasure: "hectares",
				Address:     "123 Farm Lane",
				CropProductions: []dto.CropProductionDTO{
					{
						CropType:    "InvalidType",
						IsIrrigated: true,
						IsInsured:   true,
					},
				},
			},
			expectedStatusCode: fiber.StatusBadRequest,
			mockResponse:       nil,
			mockError:          nil,
			mockRequired:       false,
		},
		{
			name: "Internal Server Error - Mock Use Case Error",
			inputDTO: dto.CreateFarmDTO{
				Name:            "Test Farm",
				LandArea:        100.5,
				UnitMeasure:     "hectares",
				Address:         "123 Farm Lane",
				CropProductions: []dto.CropProductionDTO{},
			},
			expectedStatusCode: fiber.StatusInternalServerError,
			mockResponse:       nil,
			mockError:          assert.AnError,
			mockRequired:       true,
		},
	}
	for _, tt := range tests {
		cs.Run(tt.name, func() {
			var mockUseCase *MockCreateFarmUseCase
			if tt.mockRequired {
				mockUseCase = new(MockCreateFarmUseCase)
				mockUseCase.On("Execute", mock.Anything, mock.AnythingOfType("domain.Farm")).
					Return(tt.mockResponse, tt.mockError)
			}

			controller := NewFarmController(mockUseCase, nil, nil, cs.logger)
			app := fiber.New(fiber.Config{
				AppName:       "farm-api-test by @arthurgavazza",
				CaseSensitive: true,
			})
			app.Post("/farms", controller.CreateFarm)

			payload, err := json.Marshal(tt.inputDTO)
			assert.NoError(cs.T(), err)

			req, err := http.NewRequest("POST", "/farms", bytes.NewReader(payload))
			assert.NoError(cs.T(), err)
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			assert.NoError(cs.T(), err)

			assert.Equal(cs.T(), tt.expectedStatusCode, resp.StatusCode)

			if tt.expectedStatusCode == fiber.StatusCreated {
				var responseFarm domain.Farm
				err = json.NewDecoder(resp.Body).Decode(&responseFarm)
				assert.NoError(cs.T(), err)
				assert.Equal(cs.T(), tt.mockResponse.ID, responseFarm.ID)
			} else if tt.expectedStatusCode == fiber.StatusBadRequest || tt.expectedStatusCode == fiber.StatusInternalServerError {
				var response map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(cs.T(), err)
				assert.NotNil(cs.T(), response["error"])
			}
			if tt.mockRequired {
				mockUseCase.AssertExpectations(cs.T())
			}
		})
	}
}

func (cs *FarmControllerTestSuite) TestFarmControllerListFarms() {
	tests := []struct {
		name               string
		expectedStatusCode int
		mockResponse       *models.PaginatedResponse[*domain.Farm]
		mockError          error
		mockRequired       bool
		queryString        string
	}{
		{
			name:               "Successful farms retrieval with valid query string",
			expectedStatusCode: fiber.StatusOK,
			mockResponse: &models.PaginatedResponse[*domain.Farm]{
				TotalCount:  5,
				PerPage:     10,
				CurrentPage: 1,
				Items:       testutils.GenerateFarms(5, testutils.PointerTo(domain.CropTypeCoffee.String()), nil),
			},
			mockError:    nil,
			mockRequired: true,
			queryString:  fmt.Sprintf("?crop_type=%s", domain.CropTypeCoffee.String()),
		},
		{
			name:               "Successful farms retrieval with empty query string",
			expectedStatusCode: fiber.StatusOK,
			mockResponse: &models.PaginatedResponse[*domain.Farm]{
				TotalCount:  5,
				PerPage:     10,
				CurrentPage: 1,
				Items:       testutils.GenerateFarms(5, testutils.PointerTo(domain.CropTypeCoffee.String()), nil),
			},
			mockError:    nil,
			mockRequired: true,
			queryString:  "",
		},
		{
			name:               "Invalid query parameters request",
			expectedStatusCode: fiber.StatusBadRequest,
			mockResponse:       nil,
			mockError:          nil,
			mockRequired:       false,
			queryString:        "?maximum_land_area=test",
		},
		{
			name:               "Unknown exception in use case layer",
			expectedStatusCode: fiber.StatusInternalServerError,
			mockResponse:       nil,
			mockError:          errors.New("Unknown error"),
			mockRequired:       true,
			queryString:        fmt.Sprintf("?crop_type=%s", domain.CropTypeCoffee.String()),
		},
	}

	for _, tt := range tests {
		cs.Run(tt.name, func() {
			var mockUseCase *MockListFarmsUseCase
			if tt.mockRequired {
				mockUseCase = new(MockListFarmsUseCase)
				mockUseCase.On("Execute", mock.Anything, mock.AnythingOfType("*domain.FarmSearchParameters")).
					Return(tt.mockResponse, tt.mockError)
			}

			controller := NewFarmController(nil, mockUseCase, nil, cs.logger)
			app := fiber.New(fiber.Config{
				AppName:       "farm-api-test by @arthurgavazza",
				CaseSensitive: true,
			})
			app.Get("/farms", controller.ListFarms)
			path := "/farms"
			route := fmt.Sprintf("%s%s", path, tt.queryString)
			req, err := http.NewRequest("GET", route, nil)
			assert.NoError(cs.T(), err)
			resp, err := app.Test(req)
			assert.NoError(cs.T(), err)

			assert.Equal(cs.T(), tt.expectedStatusCode, resp.StatusCode)

			if tt.expectedStatusCode == fiber.StatusOK {
				var response models.PaginatedResponse[*domain.Farm]
				err = json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(cs.T(), err)
				for i, item := range response.Items {
					assert.Equal(cs.T(), tt.mockResponse.Items[i].ID.String(), item.ID.String())
				}
				assert.Equal(cs.T(), tt.mockResponse.PerPage, response.PerPage)
				assert.Equal(cs.T(), tt.mockResponse.TotalCount, response.TotalCount)
				assert.Equal(cs.T(), tt.mockResponse.CurrentPage, response.CurrentPage)
			} else if tt.expectedStatusCode == fiber.StatusBadRequest || tt.expectedStatusCode == fiber.StatusInternalServerError {
				var response map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(cs.T(), err)
				assert.NotNil(cs.T(), response["error"])
			}
			if tt.mockRequired {
				mockUseCase.AssertExpectations(cs.T())
			}
		})
	}
}

func (cs *FarmControllerTestSuite) TestFarmControllerDeleteFarm() {
	farmId := uuid.New().String()
	notFoundErr := &shared.NotFoundError{
		Resource: "Farm",
		ID:       farmId,
	}
	tests := []struct {
		name               string
		expectedStatusCode int
		mockError          error
		mockRequired       bool
		farmId             string
	}{
		{
			name:               "Successful farm deletion",
			expectedStatusCode: fiber.StatusNoContent,
			mockError:          nil,
			mockRequired:       true,
			farmId:             farmId,
		},
		{
			name:               "Invalid request -  farm not found",
			expectedStatusCode: fiber.StatusNotFound,
			mockError:          notFoundErr,
			mockRequired:       true,
			farmId:             farmId,
		},
	}

	for _, tt := range tests {
		cs.Run(tt.name, func() {
			var mockUseCase *MockDeleteFarmUseCase
			if tt.mockRequired {
				mockUseCase = new(MockDeleteFarmUseCase)
				mockUseCase.On("Execute", mock.Anything, mock.AnythingOfType("string")).
					Return(tt.mockError)
			}

			controller := NewFarmController(nil, nil, mockUseCase, cs.logger)
			app := fiber.New(fiber.Config{
				AppName:       "farm-api-test by @arthurgavazza",
				CaseSensitive: true,
			})
			app.Delete("/farms/:id", controller.DeleteFarm)
			route := fmt.Sprintf("/farms/%s", tt.farmId)
			req, err := http.NewRequest("DELETE", route, nil)
			assert.NoError(cs.T(), err)
			resp, err := app.Test(req)
			assert.NoError(cs.T(), err)

			assert.Equal(cs.T(), tt.expectedStatusCode, resp.StatusCode)
			if tt.expectedStatusCode == fiber.StatusBadRequest || tt.expectedStatusCode == fiber.StatusInternalServerError {
				var response map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(cs.T(), err)
				assert.NotNil(cs.T(), response["error"])
			}
			if tt.mockRequired {
				mockUseCase.AssertExpectations(cs.T())
			}
		})
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(FarmControllerTestSuite))
}
