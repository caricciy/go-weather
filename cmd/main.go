package main

import (
	"errors"
	"github.com/caricciy/go-weather/internal/infra"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables directly")
	}

	// Configure logger
	logOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	l := slog.New(slog.NewJSONHandler(os.Stdout, logOpts))
	l = l.With(slog.String("app_name", "go-weather"))
	slog.SetDefault(l)
}

func main() {
	router := infra.NewAppRouter()

	// Initialize handlers
	weatherHandler := infra.NewWeatherHandler()

	// Define routes
	router.Get("/weather/{cep}", weatherHandler.HandleGetWeatherByCEP)

	server := infra.NewHttpServer(router)

	// Start the server in a goroutine so it doesn't block
	go func() {
		log.Printf("Starting server on port %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
		}
	}()

	infra.WaitForShutdown(server)
}
