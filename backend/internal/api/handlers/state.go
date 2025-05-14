package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"radaroficial.app/internal/institutions"
)

type StateHandler struct{ DB *pgxpool.Pool }

func NewStateHandler(db *pgxpool.Pool) *StateHandler {
	return &StateHandler{DB: db}
}

func (h *StateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.NotFound(w, r)
		return
	}

	srv := institutions.NewInstitutionService(h.DB)
	states, err := srv.GetStates(r.Context())

	if err != nil {
		http.Error(w, "Failed to fetch states.", http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(map[string]any{
			"states": states,
		})
	}

}
