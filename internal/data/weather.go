package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/caricciy/go-weather/internal/entity"
	"io"
	"net/http"
	url2 "net/url"
)

type weatherDTO struct {
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
	Current struct {
		TempC float64 `json:"temp_c"`
		TempF float64 `json:"temp_f"`
	} `json:"current"`
}

type WeatherApiRepository struct {
	apiKey         string
	targetEndpoint string
}

// NewWeatherApiStore creates a new instance of WeatherApiRepository
func NewWeatherApiStore(apiKey string) *WeatherApiRepository {
	return &WeatherApiRepository{
		apiKey:         apiKey,
		targetEndpoint: "https://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no",
	}
}

func (w *WeatherApiRepository) GetWeatherInfo(ctx context.Context, cep *entity.CEP) (*entity.WeatherInfo, error) {
	escapedLocation := url2.QueryEscape(cep.Localidade)
	url := fmt.Sprintf(w.targetEndpoint, w.apiKey, escapedLocation)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch Weather: received status code %d", resp.StatusCode)
	}

	var weatherData weatherDTO
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Simulated response
	return &entity.WeatherInfo{
		Celcius:    weatherData.Current.TempC,
		Fahrenheit: weatherData.Current.TempF,
	}, nil
}
