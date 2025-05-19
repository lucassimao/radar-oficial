package handlers

import (
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"radaroficial.app/internal/jobs"
)

type JobsHandler struct{ DB *pgxpool.Pool }

func NewJobsHandler(db *pgxpool.Pool) *JobsHandler {
	return &JobsHandler{DB: db}
}

func (h *JobsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	queryValues := r.URL.Query()

	if !queryValues.Has("name") {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	jobName := queryValues.Get("name")

	riverClient, err := jobs.NewRiverClient(r.Context(), h.DB)
	if err != nil {
		log.Printf("❌ Failed to initialize River Queue client: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	switch jobName {
	case "doepi":
		// YYY-MM-DD
		customDate := queryValues.Get("customDate") // defaults to "" if no custom date set, which is interpreted as the current date
		jobs.ScheduleGovernoPiauiJob(r.Context(), riverClient.Client, customDate)
		w.WriteHeader(http.StatusNoContent)
	case "municipios_pi":
		jobs.ScheduleDiarioDosMunicipiosJob(r.Context(), riverClient.Client)
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "bad request", http.StatusBadRequest)
	}

}
