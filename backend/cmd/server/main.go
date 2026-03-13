package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joaodias/golang-svelte-fullstack/internal/config"
	"github.com/joaodias/golang-svelte-fullstack/internal/database"
	"github.com/joaodias/golang-svelte-fullstack/internal/handler"
	"github.com/joaodias/golang-svelte-fullstack/internal/middleware"
	"github.com/joaodias/golang-svelte-fullstack/internal/repository"
	"github.com/joaodias/golang-svelte-fullstack/internal/router"
)

func main() {
	cfg := config.Load()

	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db, "migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	userRepo := repository.NewUserRepository(db.Pool)
	tokenRepo := repository.NewTokenRepository(db.Pool)

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	userHandler := handler.NewUserHandler(userRepo, tokenRepo, logger)
	authMiddleware := middleware.NewAuthMiddleware(userRepo, tokenRepo)

	mux := router.New(userHandler, authMiddleware)

	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server gracefully stopped")
}
