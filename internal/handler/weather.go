package handler

import (
	"context"
	"errors"
	"github.com/caricciy/go-weather/internal/usecase"
	"github.com/caricciy/go-weather/internal/util"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

// WeatherResponse represents the structure of the weather response
type getWeatherByCEPResponse struct {
	Celcius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

type errorResponse struct {
	Message string `json:"message"`
}

type WeatherHandler struct {
	cepUseCases *usecase.WeatherUseCases
}

func NewWeatherHandler(cepUseCases *usecase.WeatherUseCases) *WeatherHandler {
	return &WeatherHandler{
		cepUseCases: cepUseCases,
	}
}

// HandleGetWeatherByCEP handles the request to get CEP information
func (h *WeatherHandler) HandleGetWeatherByCEP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	paramCEP := chi.URLParam(r, "cep")

	weather, err := h.cepUseCases.GetWeatherByCEP(ctx, paramCEP)

	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidCEP):
			util.SendJSON(w, errorResponse{"invalid zipcode"}, http.StatusUnprocessableEntity)
		case errors.Is(err, usecase.ErrCEPNotFound):
			util.SendJSON(w, errorResponse{"can not find zipcode"}, http.StatusNotFound)
		default:
			util.SendJSON(w, errorResponse{"An unexpected error occurred"}, http.StatusInternalServerError)
		}
		return
	}
	response := getWeatherByCEPResponse{
		Celcius:    weather.Celcius,
		Fahrenheit: weather.Fahrenheit,
		Kelvin:     weather.Kelvin,
	}

	util.SendJSON(w, response, http.StatusOK)

}
