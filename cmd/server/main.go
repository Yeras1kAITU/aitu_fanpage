package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Yeras1kAITU/aitu_fanpage/internal/app"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/config"
)

func main() {
	cfg := config.Load()

	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}

	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Printf("Environment: %s", cfg.Server.Env)
	log.Printf("Database: %s", cfg.Database.Name)

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	server := &http.Server{
		Addr:         "0.0.0.0:" + cfg.Server.Port,
		Handler:      application.Router(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	serverErrors := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		log.Printf("Server listening on %s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("Server error: %v", err)

	case sig := <-shutdown:
		log.Printf("Received signal: %v. Starting graceful shutdown...", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Graceful shutdown failed: %v", err)
		}

		log.Println("Server stopped gracefully")
	}
}
