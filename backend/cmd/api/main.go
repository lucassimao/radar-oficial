package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"radaroficial.app/internal/api"
	"radaroficial.app/internal/config"
)

var pool *pgxpool.Pool

func init() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("Invalid DB_PORT: %v", err)
	}

	sslMode := "disable"
	if config.Env() == "production" {
		sslMode = "require"
	}

	// e.g., postgres://user:pass@host:5432/dbname
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", dbUser, dbPassword,
		dbHost, dbPort, dbName, sslMode)

	pool, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
}

func main() {
	// Load env variables
	_ = godotenv.Load()

	// Setup signal catching
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize and start the server
	server := api.NewServer(pool)
	server.RegisterHandlers()

	// Start the server in a goroutine
	go func() {
		if err := server.Start(port); err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()

	// Shutdown the server
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	// Explicitly close the DB connection pool
	log.Println("Closing database connection pool...")
	pool.Close()

	log.Println("Application shutdown complete")
}
