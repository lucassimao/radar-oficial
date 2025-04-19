package jobs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"radaroficial.app/internal/diarios"
	"radaroficial.app/internal/storage"
)

// DiarioDosMunicipiosArgs contains arguments for the job
type DiarioDosMunicipiosArgs struct{}

// Kind returns the kind of job
func (DiarioDosMunicipiosArgs) Kind() string { return "fetch_diario_dos_municipios" }

// DiarioWorker handles diario-related jobs
type DiarioWorker struct {
	// Embed worker defaults
	river.WorkerDefaults[DiarioDosMunicipiosArgs]
	
	// Add dependencies
	DiarioService *diarios.DiarioService
	Uploader      *storage.SpacesUploader
}

// NewDiarioWorker creates a new DiarioWorker
func NewDiarioWorker(diarioService *diarios.DiarioService, uploader *storage.SpacesUploader) *DiarioWorker {
	return &DiarioWorker{
		DiarioService: diarioService,
		Uploader:      uploader,
	}
}

// Work processes the fetch job for Diario dos Municipios
func (w *DiarioWorker) Work(ctx context.Context, job *river.Job[DiarioDosMunicipiosArgs]) error {
	log.Printf("🔄 Starting job to fetch Diário dos Municípios (ID: %d)", job.ID)
	
	// Fetch and process the diario
	entries, err := diarios.FetchDiarioDosMunicipiosPiaui(ctx, w.Uploader, w.DiarioService)
	if err != nil {
		return fmt.Errorf("failed to fetch diario dos municipios: %w", err)
	}
	
	log.Printf("✅ Job completed successfully. Processed %d diário(s) from Municípios do Piauí", len(entries))
	return nil
}

// MaxRetries defines max attempts for this job
func (w *DiarioWorker) MaxRetries(job *river.Job[DiarioDosMunicipiosArgs]) int {
	return 3 // Retry up to 3 times
}

// Timeout sets the maximum execution time for this job
func (w *DiarioWorker) Timeout(job *river.Job[DiarioDosMunicipiosArgs]) time.Duration {
	return 10 * time.Minute // Allow up to 10 minutes to fetch and process
}

// CreatePeriodicJob returns a periodic job configuration for Diario dos Municipios
func CreatePeriodicJob() *river.PeriodicJob {
	return river.NewPeriodicJob(
		// Run every hour
		river.PeriodicInterval(1 * time.Hour),
		
		// Args constructor function
		func() (river.JobArgs, *river.InsertOpts) {
			return DiarioDosMunicipiosArgs{}, &river.InsertOpts{
				Queue:    "default",
				Priority: 1, // Higher number = higher priority
			}
		},
		
		// Options - use nil for default options
		nil,
	)
}

// ScheduleDiarioDosMunicipiosJob schedules a job to fetch the Diario dos Municipios to run immediately
func ScheduleDiarioDosMunicipiosJob(ctx context.Context, client *river.Client[pgx.Tx]) error {
	// Insert a job to run immediately
	_, err := client.Insert(ctx, DiarioDosMunicipiosArgs{}, &river.InsertOpts{
		Queue:    "default",
		Priority: 1, // Higher number = higher priority
	})
	
	if err != nil {
		return fmt.Errorf("failed to schedule immediate diario job: %w", err)
	}
	
	log.Printf("✅ Scheduled immediate job to fetch Diario dos Municipios")
	return nil
}