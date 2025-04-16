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
	"radaroficial.app/internal/storage"
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

	http.HandleFunc("/crawl", handleCrawl)

	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleCrawl(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slug := r.URL.Query().Get("slug")
	dateStr := r.URL.Query().Get("date")

	// If slug is not "governo-pi", return 204
	if slug != "governo-pi" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Parse date or fallback to today
	var fetchDate time.Time
	var err error
	if dateStr != "" {
		fetchDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			http.Error(w, "Invalid date format. Use YYYY-MM-DD.", http.StatusBadRequest)
			return
		}
	} else {
		fetchDate = time.Now()
	}

	uploader, err := storage.NewSpacesUploader("radar-oficial-diarios-piaui")
	if err != nil {
		http.Error(w, "Failed to initialize storage", http.StatusInternalServerError)
		return
	}

	// Fetch and save diarios
	entries, err := diarios.FetchGovernoPiauiDiarios(ctx, fetchDate, uploader)
	if err != nil {
		log.Printf("❌ Failed to fetch diarios: %v", err)
		http.Error(w, "Failed to fetch diarios", http.StatusInternalServerError)
		return
	}

	service := diarios.NewDiarioService(pool)
	for _, d := range entries {
		if err := service.Insert(ctx, d); err != nil {
			log.Printf("⚠️ Failed to insert diário for %s: %v", d.SourceURL, err)
		} else {
			log.Printf("✅ Inserted diário %s", d.SourceURL)
		}
	}

	fmt.Fprintf(w, "📥 Processed %d diário(s)\n", len(entries))
}
