package main

import (
	"Subscribe-service/internal/config"
	"Subscribe-service/internal/handler"
	"Subscribe-service/internal/repository"
	"Subscribe-service/internal/server"
	"Subscribe-service/internal/service"
	"Subscribe-service/internal/storage"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println("Loaded config:", cfg)

	db, err := storage.New(cfg)
	if err != nil {
		log.Fatalf("storage.New error: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}()

	subRepo := repository.NewSubRepository(db)

	SubService := service.NewSubService(teamRepo)

	h := handler.New(subService)

	srv := server.New(cfg, db)
	go func() {
		if err := srv.Run(h); err != nil {
			log.Printf("server run error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Stop(ctx); err != nil {
		log.Printf("Error stopping server: %v", err)
	}

	log.Println("Shutdown complete")
}
