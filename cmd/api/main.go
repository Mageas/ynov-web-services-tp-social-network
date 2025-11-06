package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ynov-social-api/internal/api/handler"
	"ynov-social-api/internal/api/router"
	"ynov-social-api/internal/config"
	"ynov-social-api/internal/pkg/logger"
	"ynov-social-api/internal/repository/sqlite"
	"ynov-social-api/internal/service/auth"
	"ynov-social-api/internal/service/post"
	"ynov-social-api/internal/service/user"
)

func main() {
	// Initialize logger
	log := logger.New()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := sqlite.New(cfg.Database.Path)
	if err != nil {
		log.Fatal("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Info("Database connected successfully")

	// Initialize repositories
	userRepo := sqlite.NewUserRepository(db.GetConn())
	postRepo := sqlite.NewPostRepository(db.GetConn())

	// Initialize services
	passwordService := auth.NewPasswordService()
	jwtService := auth.NewJWTService(cfg.JWT.Secret, cfg.JWT.TTL)
	userService := user.NewService(userRepo, passwordService)
	postService := post.NewService(postRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(userService, jwtService, log)
	postHandler := handler.NewPostHandler(postService, log)

	// Initialize router
	r := router.New(authHandler, postHandler, jwtService)

	// Configure HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Info("Starting server on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown: %v", err)
	}

	log.Info("Server stopped")
}
