package main

import (
	"Subscribe-service/internal/config"
	"Subscribe-service/internal/handler"
	"Subscribe-service/internal/repository"
	"Subscribe-service/internal/server"
	"Subscribe-service/internal/service"
	"Subscribe-service/internal/storage"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Subscribe Service API
// @version 1.0
// @description REST API для управления подписками пользователей
// @host localhost:8080
// @BasePath /
func main() {
	log.Println("INIT starting service")

	cfg := config.MustLoad()

	db, err := storage.New(cfg)
	if err != nil {
		log.Fatalf("DB connection error %v", err)
	}
	defer db.Close()

	repo := repository.NewSubsRepository(db)
	serv := service.NewSubsService(repo)
	h := handler.New(serv)

	srv := server.New(cfg, db)

	go func() {
		if err := srv.Run(h); err != nil {
			log.Printf("SERVER error %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("SHUTDOWN signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Stop(ctx); err != nil {
		log.Printf("SHUTDOWN error %v", err)
	}

	log.Println("SHUTDOWN completed")
}
