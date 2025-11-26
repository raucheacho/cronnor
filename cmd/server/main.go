package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rauche/cronnor/internal/config"
	"github.com/rauche/cronnor/internal/http"
	"github.com/rauche/cronnor/internal/jobs"
	"github.com/rauche/cronnor/internal/storage"
)

func main() {
	log.Println("🚀 Starting Cronnor HTTP Cron Server...")

	// Load configuration
	cfg := config.Load()
	log.Printf("Configuration loaded: Port=%s, DB=%s", cfg.Port, cfg.DBPath)

	// Initialize database
	repo, err := storage.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer repo.Close()

	// Run migrations
	log.Println("Running database migrations...")
	if err := repo.RunMigrations(cfg.MigrationPath); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("✅ Database migrations completed")

	// Initialize scheduler
	scheduler := jobs.NewScheduler(repo)
	if err := scheduler.Start(); err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}
	log.Println("✅ Job scheduler started")

	// Initialize HTTP server
	server, err := http.NewServer(repo, scheduler)
	if err != nil {
		log.Fatalf("Failed to create HTTP server: %v", err)
	}

	// Start server in a goroutine
	go func() {
		if err := server.Start(cfg.Port); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	log.Printf("✅ HTTP server listening on :%s", cfg.Port)
	log.Printf("🌐 Open http://localhost:%s in your browser", cfg.Port)

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down gracefully...")

	// Stop scheduler
	scheduler.Stop()

	// Give time for cleanup
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	<-ctx.Done()

	log.Println("👋 Server stopped")
}
