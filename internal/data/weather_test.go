package data

import (
	"context"
	"encoding/json"
	"github.com/caricciy/go-weather/internal/entity"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

// createMockWeatherServer creates a mock HTTP server that simulates the Weather API response.
func createMockWeatherServer(mockResponse any) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a delay if no mock response is provided
		if mockResponse == nil || reflect.ValueOf(mockResponse).IsZero() {
			time.Sleep(2 * time.Second)
			w.WriteHeader(http.StatusOK)
			return
		}

		// Handle specific mock response for a known path
		if r.URL.Path == "/v1/current.json" {
			// Check if the query parameters match the expected format
			if r.URL.Query().Get("q") == "Unknown" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(mockResponse)
			return
		}

		// Default to 404 for unknown paths
		w.WriteHeader(http.StatusNotFound)
	}))
}

func TestGetWeatherInfo(t *testing.T) {
	// Mock data for a valid weather response
	mockWeather := weatherDTO{
		Location: struct {
			Name string `json:"name"`
		}{Name: "São Paulo"},
		Current: struct {
			TempC float64 `json:"temp_c"`
			TempF float64 `json:"temp_f"`
		}{TempC: 25.0, TempF: 77.0},
	}

	mockServer := createMockWeatherServer(mockWeather)
	defer mockServer.Close()

	// Replace the target endpoint with the mock server URL
	store := &WeatherApiRepository{
		apiKey:         "test-api-key",
		targetEndpoint: mockServer.URL + "/v1/current.json?key=%s&q=%s&aqi=no",
	}

	t.Run("Valid Weather Info", func(t *testing.T) {
		cep := &entity.CEP{Localidade: "São Paulo"}
		weather, err := store.GetWeatherInfo(context.Background(), cep)
		assert.NoError(t, err)
		assert.NotNil(t, weather)
		assert.Equal(t, 25.0, weather.Celcius)
		assert.Equal(t, 77.0, weather.Fahrenheit)
	})

	t.Run("Invalid Weather Info", func(t *testing.T) {
		cep := &entity.CEP{Localidade: "Unknown"}
		weather, err := store.GetWeatherInfo(context.Background(), cep)
		assert.Error(t, err)
		assert.Nil(t, weather)
	})
}

func TestGetWeatherInfo_Timeout(t *testing.T) {
	// Create a mock HTTP server that delays its response
	mockServer := createMockWeatherServer(nil) // No response data, just a delay
	defer mockServer.Close()

	// Replace the target endpoint with the mock server URL
	store := &WeatherApiRepository{
		apiKey:         "test-api-key",
		targetEndpoint: mockServer.URL + "/v1/current.json?key=%s&q=%s&aqi=no",
	}

	// Set a timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	cep := &entity.CEP{Localidade: "São Paulo"}
	weather, err := store.GetWeatherInfo(ctx, cep)

	// Assert that a timeout error occurred
	assert.Error(t, err)
	assert.Nil(t, weather)
}

func TestGetWeatherInfo_InvalidResponse(t *testing.T) {
	mockServer := createMockWeatherServer(map[string]string{"invalid_field": "This is not a valid response"})
	defer mockServer.Close()

	store := &WeatherApiRepository{
		apiKey:         "test-api-key",
		targetEndpoint: mockServer.URL + "/v1/current.json?key=%s&q=%s&aqi=no",
	}

	cep := &entity.CEP{Localidade: "São Paulo"}
	weather, err := store.GetWeatherInfo(context.Background(), cep)
	assert.NoError(t, err)
	assert.NotNil(t, weather)

	assert.Empty(t, weather.Fahrenheit)
	assert.Empty(t, weather.Celcius)
	assert.Empty(t, weather.Kelvin)
}
