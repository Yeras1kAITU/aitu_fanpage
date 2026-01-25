package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"aitu_fanpage/internal/app"
	"aitu_fanpage/internal/config"
)

func main() {
	cfg := config.Load()

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	go func() {
		if err := application.Start(); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	application.Shutdown()
}
