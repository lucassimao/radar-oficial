package handlers

import (
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"radaroficial.app/internal/diarios"
)

type ReindexHandler struct{ DB *pgxpool.Pool }

func NewReindexHandler(db *pgxpool.Pool) *ReindexHandler {
	return &ReindexHandler{DB: db}
}

func (h *ReindexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	srv := diarios.NewInstitutionService(h.DB)

	err := srv.ReIndexKnowledgeBases(r.Context())
	if err != nil {
		log.Printf("❌ Failed to trigger reindex: %v", err)
		http.Error(w, "Failed to trigger reindex.", http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
