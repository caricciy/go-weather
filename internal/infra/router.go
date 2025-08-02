package infra

import (
	"github.com/caricciy/go-weather/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

func NewAppRouter() *chi.Mux {
	router := chi.NewRouter()

	// Add middleware
	router.Use(middleware.Logger, middleware.Recoverer, middleware.RequestID)

	// Set up health check route
	router.Get("/health", health)

	return router
}

func health(w http.ResponseWriter, r *http.Request) {
	status := struct {
		Status string    `json:"status"`
		Time   time.Time `json:"time"`
	}{Status: "UP", Time: time.Now()}

	util.SendJSON(w, status, http.StatusOK)
}
