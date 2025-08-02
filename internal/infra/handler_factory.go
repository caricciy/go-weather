package infra

import (
	"github.com/caricciy/go-weather/internal/data"
	"github.com/caricciy/go-weather/internal/handler"
	"github.com/caricciy/go-weather/internal/usecase"
	"os"
)

func NewWeatherHandler() *handler.WeatherHandler {
	weatherApiKey := os.Getenv("WEATHER_API_KEY")
	vcs := data.NewViaCEPStore()
	ws := data.NewWeatherApiStore(weatherApiKey)
	uc := usecase.NewWeatherUseCases(vcs, ws)
	return handler.NewWeatherHandler(uc)
}
