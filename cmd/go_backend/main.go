package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aayuskarki/go_backend/internal/config"
	"github.com/aayuskarki/go_backend/internal/http/handlers/student"
)

func main() {
	// Load configuration
	cfg := config.MustLoad()

	// Set up router
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New())

	// Initialize HTTP server with loaded configuration
	server := http.Server{
		Addr:    cfg.HTTPServer.Address,
		Handler: router,
	}

	fmt.Println("Server started on", cfg.HTTPServer.Address)

	// Create a channel to listen for OS signals for graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine to allow graceful shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Block until we receive a signal
	<-done
	slog.Info("Server shutting down...")

	// Create context with timeout for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Error during server shutdown", slog.String("error", err.Error()))
	} else {
		slog.Info("Server shutdown completed")
	}
}
