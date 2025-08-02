package data

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

// createMockServer creates a mock HTTP server that simulates the ViaCEP API response.
// (server port will be random)
func createMockServer(mockResponse any) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a delay if no mock response is provided
		if mockResponse == nil || reflect.ValueOf(mockResponse).IsZero() {
			time.Sleep(2 * time.Second)
			w.WriteHeader(http.StatusOK)
			return
		}

		// Handle specific mock response for a known path
		if r.URL.Path == "/ws/12345678/json" {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(mockResponse)
			return
		}

		// Default to 404 for unknown paths
		w.WriteHeader(http.StatusNotFound)
	}))
}

func TestGetCEP(t *testing.T) {
	// Mock data for a valid CEP response
	mockCEP := cepDTO{
		Localidade: "São Paulo",
	}

	mockServer := createMockServer(mockCEP)
	defer mockServer.Close()

	// Replace the target endpoint with the mock server URL
	store := &ViaCEPStore{targetEndpoint: mockServer.URL + "/ws/%s/json"}

	t.Run("Valid CEP", func(t *testing.T) {
		cep, err := store.GetCEP(context.Background(), "12345678")
		assert.NoError(t, err)
		assert.NotNil(t, cep)
		assert.Equal(t, "São Paulo", cep.Localidade)
	})

	t.Run("Invalid CEP", func(t *testing.T) {
		cep, err := store.GetCEP(context.Background(), "00000000")
		assert.Error(t, err)
		assert.Nil(t, cep)
	})

	t.Run("Valid CEP with Invalid Context", func(t *testing.T) {
		cep, err := store.GetCEP(nil, "12345678")
		assert.Error(t, err)
		assert.Nil(t, cep)
	})
}

func TestGetCEP_Timeout(t *testing.T) {
	// Create a mock HTTP server that delays its response (server port will be random)
	mockServer := createMockServer(0) // No response data, just a delay

	defer mockServer.Close()

	// Replace the target endpoint with the mock server URL
	store := &ViaCEPStore{targetEndpoint: mockServer.URL + "/ws/%s/json"}

	// Set a timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Call GetCEP with the timeout context
	cep, err := store.GetCEP(ctx, "12345678")

	// Assert that a timeout error occurred
	assert.Error(t, err)
	assert.Nil(t, cep)
}

func TestGetCEP_InvalidResponse(t *testing.T) {
	mockServer := createMockServer(map[string]string{"invalid_field": "This is not a valid response"})
	defer mockServer.Close()

	store := &ViaCEPStore{targetEndpoint: mockServer.URL + "/ws/%s/json"}

	cep, err := store.GetCEP(context.Background(), "12345678")
	assert.NoError(t, err)
	assert.NotNil(t, cep)
	// Expecting Localidade to be empty due to invalid response
	assert.Empty(t, cep.Localidade)
}
