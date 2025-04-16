package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"radaroficial.app/internal/config"
	"radaroficial.app/internal/diarios"
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
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {

		today := time.Now()
		// yesterday := time.Now().AddDate(0, 0, -1)

		entries, err := diarios.FetchGovernoPiauiDiarios(context.Background(), today)
		if err != nil {
			log.Fatal("❌ Failed to fetch diarios:", err)
		}

		service := diarios.NewDiarioService(pool)
		for _, d := range entries {
			err := service.Insert(context.Background(), d)
			if err != nil {
				log.Printf("⚠️ Failed to insert diário for %s: %v", d.SourceURL, err)
			} else {
				log.Printf("✅ Inserted diário %s", d.SourceURL)
			}
		}
		fmt.Fprint(w, "pong")
	})

	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
