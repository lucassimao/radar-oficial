package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"radaroficial.app/internal/diarios"
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

	// Create the diario service
	service := diarios.NewDiarioService(h.DB)
	
	// Fetch and save diarios
	entries, err := diarios.FetchGovernoPiauiDiarios(ctx, fetchDate, uploader, service)
	if err != nil {
		log.Printf("❌ Failed to fetch diarios: %v", err)
		http.Error(w, "Failed to fetch diarios", http.StatusInternalServerError)
		return
	}
	for _, d := range entries {
		if err := service.Insert(ctx, d); err != nil {
			log.Printf("⚠️ Failed to insert diário for %s: %v", d.SourceURL, err)
		} else {
			log.Printf("✅ Inserted diário %s", d.SourceURL)
		}
	}

	fmt.Fprintf(w, "📥 Processed %d diário(s)\n", len(entries))
}
