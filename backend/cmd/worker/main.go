package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"radaroficial.app/internal/jobs"
)

var pool *pgxpool.Pool

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("⚠️ Warning: Error loading .env file: %v", err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPortStr := os.Getenv("DB_PORT")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil || dbPortStr == "" {
		dbPort = 5432 // Default PostgreSQL port
		log.Printf("⚠️ Warning: Invalid DB_PORT, using default: %d", dbPort)
	}

	sslMode := os.Getenv("DB_SSL_MODE")

	// e.g., postgres://user:pass@host:5432/dbname
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, sslMode)

	// Configure connection pool
	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("❌ Failed to parse database config: %v", err)
	}

	// Set reasonable pool limits
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = 1 * time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	// Create the connection pool
	pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
}

func main() {
	log.Printf("📋 Starting Radar Oficial worker process...")

	// Create context that listens for termination signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Make sure to close the pool when we exit
	defer pool.Close()

	// Verify database connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("❌ Database ping failed: %v", err)
	}
	log.Printf("✅ Connected to database")

	// Initialize and start River Queue client
	riverClient, err := jobs.NewRiverClient(ctx, pool)
	if err != nil {
		log.Fatalf("❌ Failed to initialize River Queue client: %v", err)
	}

	// Schedule initial jobs
	if err := riverClient.ScheduleInitialJobs(ctx); err != nil {
		log.Printf("⚠️ Failed to schedule initial jobs: %v", err)
		// Continue running even if scheduling fails
	} else {
		log.Printf("✅ Initial jobs scheduled successfully")
	}

	// Log active periodic jobs
	log.Printf("🔄 Worker configured with the following periodic jobs running every hour:")
	log.Printf("  • Diário dos Municípios do Piauí")
	log.Printf("  • Diários Oficiais do Governo do Piauí")

	// Wait for termination signal
	log.Printf("🔄 Worker is now running. Press Ctrl+C to exit...")
	<-ctx.Done()

	// Graceful shutdown
	log.Printf("⏹️ Shutting down worker...")

	// Create a new context with timeout for shutdown operations
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown River Queue client
	riverClient.Shutdown(shutdownCtx)

	log.Printf("👋 Worker shutdown complete")
}
