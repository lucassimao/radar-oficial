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

	// Only support governo-pi slug now
	if slug != "governo-pi" {
		// 'municipios-pi' has been moved to the scheduled worker
		if slug == "municipios-pi" {
			http.Error(w, "The municipios-pi endpoint has been moved to an automated scheduled job", http.StatusGone)
			return
		}

		// If slug is not supported, return 204
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Create the diario service
	service := diarios.NewDiarioService(h.DB)

	// Initialize storage
	uploader, err := storage.NewSpacesUploader("radar-oficial-diarios-piaui")
	if err != nil {
		http.Error(w, "Failed to initialize storage", http.StatusInternalServerError)
		return
	}

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
	entries, err := diarios.FetchGovernoPiauiDiarios(ctx, fetchDate, uploader, service)
	if err != nil {
		log.Printf("❌ Failed to fetch diarios from Governo do Piauí: %v", err)
		http.Error(w, "Failed to fetch diarios", http.StatusInternalServerError)
		return
	}

	// Diarios are now inserted directly by the fetcher functions
	fmt.Fprintf(w, "📥 Processed %d diário(s) from Governo do Piauí\n", len(entries))
}
