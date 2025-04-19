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

	sslMode := os.Getenv("DB_SSL_MODE")

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
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer stop() // stop receiving signals
	defer pool.Close()

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

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("❌ Error during server shutdown: %v", err)
	} else {
		log.Println("✅ Server shutdown complete")
	}
}
