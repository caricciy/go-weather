package usecase

import (
	"context"
	"errors"
	"github.com/caricciy/go-weather/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// MockCEPRepository is a mock implementation of the CEPRepository interface.
type MockCEPRepository struct {
	mock.Mock
}

func (m *MockCEPRepository) GetCEP(ctx context.Context, cep string) (*entity.CEP, error) {
	args := m.Called(ctx, cep)
	return args.Get(0).(*entity.CEP), args.Error(1)
}

// MockWeatherRepository is a mock implementation of the WeatherRepository interface.
type MockWeatherRepository struct {
	mock.Mock
}

func (m *MockWeatherRepository) GetWeatherInfo(ctx context.Context, cep *entity.CEP) (*entity.WeatherInfo, error) {
	args := m.Called(ctx, cep)
	return args.Get(0).(*entity.WeatherInfo), args.Error(1)
}

type testRow struct {
	name                          string
	cep                           string
	mockCEP                       *entity.CEP
	mockCEPError                  error
	mockWeather                   *entity.WeatherInfo
	mockWeatherErr                error
	expectedError                 error
	expectedResult                *entity.WeatherInfo
	mockWeatherRepoShouldBeCalled bool
	mockCEPRepoShouldBeCalled     bool
}

func TestGetWeatherByCEP(t *testing.T) {
	mockCEPRepo := new(MockCEPRepository)
	mockWeatherRepo := new(MockWeatherRepository)
	useCases := NewWeatherUseCases(mockCEPRepo, mockWeatherRepo)

	testTable := []testRow{
		{
			name:          "Valid CEP and Weather Info",
			cep:           "12345678",
			mockCEP:       &entity.CEP{Localidade: "São Paulo"},
			mockWeather:   &entity.WeatherInfo{Fahrenheit: 77.0, Celcius: 25.0},
			expectedError: nil,
			expectedResult: &entity.WeatherInfo{
				Fahrenheit: 77.0,
				Celcius:    25.0,
				Kelvin:     298.15,
			},
			mockWeatherRepoShouldBeCalled: true,
			mockCEPRepoShouldBeCalled:     true,
		},
		{
			name:          "Invalid CEP",
			cep:           "invalid",
			expectedError: ErrInvalidCEP,
		},
		{
			name:                      "CEP Not Found",
			cep:                       "00000000",
			mockCEP:                   &entity.CEP{Localidade: ""},
			expectedError:             ErrCEPNotFound,
			mockCEPRepoShouldBeCalled: true,
		},
		{
			name:                          "Weather Info Not Found",
			cep:                           "12345678",
			mockCEP:                       &entity.CEP{Localidade: "São Paulo"},
			mockWeather:                   &entity.WeatherInfo{Fahrenheit: 0, Celcius: 0},
			expectedError:                 ErrWeatherNotFound,
			mockWeatherRepoShouldBeCalled: true,
			mockCEPRepoShouldBeCalled:     true,
		},
		{
			name:                      "Error Fetching CEP",
			cep:                       "12345678",
			mockCEPError:              errors.New("error fetching CEP"),
			expectedError:             ErrCouldNotFetchCEP,
			mockCEPRepoShouldBeCalled: true,
		},
		{
			name:                          "Error Fetching Weather Info",
			cep:                           "12345678",
			mockCEP:                       &entity.CEP{Localidade: "São Paulo"},
			mockWeatherErr:                errors.New("error fetching weather"),
			expectedError:                 ErrCouldNotFetchWeather,
			mockWeatherRepoShouldBeCalled: true,
			mockCEPRepoShouldBeCalled:     true,
		},
	}

	for _, tr := range testTable {
		t.Run(tr.name, func(t *testing.T) {
			// Mock CEP repository behavior
			if tr.mockCEPRepoShouldBeCalled {
				mockCEPRepo.On("GetCEP", mock.Anything, tr.cep).Return(tr.mockCEP, tr.mockCEPError).Once()
			}

			// Mock Weather repository behavior
			if tr.mockCEP != nil && tr.mockCEP.Localidade != "" && tr.mockWeatherRepoShouldBeCalled {
				mockWeatherRepo.On("GetWeatherInfo", mock.Anything, tr.mockCEP).Return(tr.mockWeather, tr.mockWeatherErr).Once()
			}

			// Call the method
			result, err := useCases.GetWeatherByCEP(context.Background(), tr.cep)

			// Assert results
			assert.Equal(t, tr.expectedError, err)
			assert.Equal(t, tr.expectedResult, result)

			// Assert expectations
			if tr.mockCEPRepoShouldBeCalled {
				mockCEPRepo.AssertExpectations(t)
			}

			if tr.mockWeatherRepoShouldBeCalled {
				mockWeatherRepo.AssertExpectations(t)
			}
		})
	}
}
