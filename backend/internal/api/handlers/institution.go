package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"radaroficial.app/internal/institutions"
)

type InstitutionHandler struct{ DB *pgxpool.Pool }

func NewInstitutionHandler(db *pgxpool.Pool) *InstitutionHandler {
	return &InstitutionHandler{DB: db}
}

type Result struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func (h *InstitutionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" && r.URL.Path == "/institutions" {

		srv := institutions.NewInstitutionService(h.DB)
		institutions, err := srv.GetAll(r.Context())

		if err != nil {
			log.Printf("❌ Failed to fetch institutions: %v", err)
			http.Error(w, "Failed to fetch institutions.", http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/json")
			var result []Result

			for _, institution := range institutions {
				result = append(result, Result{Id: institution.ID, Name: institution.Name, Slug: institution.Slug})
			}

			json.NewEncoder(w).Encode(map[string]any{
				"institutions": result,
			})
		}
	} else {
		http.NotFound(w, r)
	}

}
