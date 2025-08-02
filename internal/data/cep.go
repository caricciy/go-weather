package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/caricciy/go-weather/internal/entity"
	"io"
	"net/http"
)

// cepDTO is a data transfer object (DTO) for the CEP entity.
// It represents the data retrieved from a data source, such as a database or an external API.
type cepDTO struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade,omitempty"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type ViaCEPStore struct {
	targetEndpoint string
}

// NewViaCEPStore creates a new instance of ViaCEPStore
func NewViaCEPStore() *ViaCEPStore {
	return &ViaCEPStore{
		targetEndpoint: "https://viacep.com.br/ws/%s/json",
	}
}

// GetCEP retrieves information for a given CEP
func (s *ViaCEPStore) GetCEP(ctx context.Context, cep string) (*entity.CEP, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf(s.targetEndpoint, cep), nil)
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
		return nil, fmt.Errorf("failed to fetch CEP: received status code %d", resp.StatusCode)
	}

	var cepData cepDTO
	if err := json.NewDecoder(resp.Body).Decode(&cepData); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Simulated response
	return &entity.CEP{
		Localidade: cepData.Localidade,
	}, nil
}
