package infra

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// NewHttpServer initializes a new HTTP server with configurations
func NewHttpServer(handler *chi.Mux) *http.Server {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return server
}

// WaitForShutdown listens for OS signals and gracefully shuts down the server
func WaitForShutdown(server *http.Server) {
	// Create a channel to listen for OS signals
	shudown := make(chan os.Signal, 1)
	signal.Notify(shudown, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Wait for the shutdown signal
	<-shudown

	// Initiate graceful shutdown
	log.Println("Received shudown signal, shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		panic("Server forced to shutdown" + err.Error())
	}

	log.Println("Server gracefully stopped")
}
