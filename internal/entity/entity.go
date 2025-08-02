package entity

import "context"

type CEPRepository interface {
	GetCEP(ctx context.Context, cep string) (*CEP, error)
}

type WeatherRepository interface {
	GetWeatherInfo(ctx context.Context, cep *CEP) (*WeatherInfo, error)
}

type CEP struct {
	Localidade string
}

type WeatherInfo struct {
	Celcius    float64
	Fahrenheit float64
	Kelvin     float64
}
