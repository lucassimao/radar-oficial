package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"radaroficial.app/internal/diarios"
	"radaroficial.app/internal/model"
	"radaroficial.app/internal/storage"
)

type CrawlHandler struct {
	DB *pgxpool.Pool
}

func NewCrawlHandler(db *pgxpool.Pool) *CrawlHandler {
	return &CrawlHandler{DB: db}
}

func (h *CrawlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slug := r.URL.Query().Get("slug")
	dateStr := r.URL.Query().Get("date")

	// Create the diario service
	service := diarios.NewDiarioService(h.DB)

	// Initialize storage
	uploader, err := storage.NewSpacesUploader("radar-oficial-diarios-piaui")
	if err != nil {
		http.Error(w, "Failed to initialize storage", http.StatusInternalServerError)
		return
	}

	var entries []*model.Diario

	if slug == "governo-pi" {
		// Parse date or fallback to today
		var fetchDate time.Time
		if dateStr != "" {
			fetchDate, err = time.Parse("2006-01-02", dateStr)
			if err != nil {
				http.Error(w, "Invalid date format. Use YYYY-MM-DD.", http.StatusBadRequest)
				return
			}
		} else {
			fetchDate = time.Now()
		}

		// Fetch and save diarios from Governo do Piauí
		entries, err = diarios.FetchGovernoPiauiDiarios(ctx, fetchDate, uploader, service)
		if err != nil {
			log.Printf("❌ Failed to fetch diarios from Governo do Piauí: %v", err)
			http.Error(w, "Failed to fetch diarios", http.StatusInternalServerError)
			return
		}
		
		// Diarios are now inserted directly by the fetcher functions
		fmt.Fprintf(w, "📥 Processed %d diário(s) from Governo do Piauí\n", len(entries))
	} else if slug == "municipios-pi" {
		// Create a background context that won't be canceled when the request ends
		bgCtx := context.Background()
		
		// Start the fetching process in a goroutine
		go func() {
			// Create a copy of the DB connection for the goroutine
			asyncService := diarios.NewDiarioService(h.DB)
			
			log.Printf("🔄 Starting asynchronous fetch of Diário dos Municípios do Piauí")
			
			// Fetch and save diario from Municípios do Piauí in the background
			municipiosEntries, err := diarios.FetchDiarioDosMunicipiosPiaui(bgCtx, uploader, asyncService)
			if err != nil {
				log.Printf("❌ Async fetch failed for diario dos municípios: %v", err)
				return
			}
			
			log.Printf("✅ Async process completed successfully. Processed %d diário(s) from Municípios do Piauí", len(municipiosEntries))
		}()
		
		// Immediately respond to the user
		fmt.Fprintf(w, "🚀 Background processing started for Diário dos Municípios\n")
		fmt.Fprintf(w, "📋 The process typically takes a few minutes to complete\n")
		fmt.Fprintf(w, "📊 Results will be logged to the server console\n")
	} else {
		// If slug is not supported, return 204
		w.WriteHeader(http.StatusNoContent)
		return
	}
}
