package usecase

import (
	"context"
	"errors"
	"github.com/caricciy/go-weather/internal/entity"
	"github.com/caricciy/go-weather/internal/util"
)

var (
	ErrInvalidCEP           = errors.New("invalid cep")
	ErrCEPNotFound          = errors.New("cep not found")
	ErrCouldNotFetchCEP     = errors.New("could not fetch cep information")
	ErrCouldNotFetchWeather = errors.New("could not fetch weather information")
	ErrWeatherNotFound      = errors.New("weather information not found")
)

type WeatherUseCases struct {
	cepRepository     entity.CEPRepository
	weatherRepository entity.WeatherRepository
}

// NewWeatherUseCases creates a new instance of WeatherUseCases
func NewWeatherUseCases(cepRepository entity.CEPRepository, weatherRepository entity.WeatherRepository) *WeatherUseCases {
	return &WeatherUseCases{
		cepRepository:     cepRepository,
		weatherRepository: weatherRepository,
	}
}

// GetWeatherByCEP retrieves weather information based on the provided CEP (postal code).
func (s *WeatherUseCases) GetWeatherByCEP(ctx context.Context, cep string) (*entity.WeatherInfo, error) {

	if !util.CheckCEPIsValid(cep) {
		return nil, ErrInvalidCEP
	}

	// Get CEP information
	c, err := s.cepRepository.GetCEP(ctx, cep)

	if err != nil {
		return nil, ErrCouldNotFetchCEP
	}

	if c.Localidade == "" {
		return nil, ErrCEPNotFound
	}

	// Get WeatherInfo based on the CEP information
	stepWeatherInfo, err := s.weatherRepository.GetWeatherInfo(ctx, c)

	if err != nil {
		return nil, ErrCouldNotFetchWeather
	}

	if stepWeatherInfo.Fahrenheit == 0 && stepWeatherInfo.Celcius == 0 {
		return nil, ErrWeatherNotFound
	}

	kelvin := stepWeatherInfo.Celcius + 273.15

	return &entity.WeatherInfo{
		Fahrenheit: stepWeatherInfo.Fahrenheit,
		Celcius:    stepWeatherInfo.Celcius,
		Kelvin:     kelvin,
	}, nil
}
