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

// GovernoPiauiArgs contains arguments for the job
type GovernoPiauiArgs struct {
	Date string `json:"date"` // Optional date in YYYY-MM-DD format
}

// Kind returns the kind of job
func (GovernoPiauiArgs) Kind() string { return "fetch_governo_piaui" }

// GovernoPiauiWorker handles diario-related jobs for Governo do Piauí
type GovernoPiauiWorker struct {
	// Embed worker defaults
	river.WorkerDefaults[GovernoPiauiArgs]
	
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

// NewGovernoPiauiWorker creates a new GovernoPiauiWorker
func NewGovernoPiauiWorker(diarioService *diarios.DiarioService, uploader *storage.SpacesUploader) *GovernoPiauiWorker {
	return &GovernoPiauiWorker{
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

// Work processes the fetch job for Governo do Piauí diarios
func (w *GovernoPiauiWorker) Work(ctx context.Context, job *river.Job[GovernoPiauiArgs]) error {
	log.Printf("🔄 Starting job to fetch Diários from Governo do Piauí (ID: %d)", job.ID)
	
	// Parse date if provided, otherwise use current date
	var fetchDate time.Time
	var err error
	
	if job.Args.Date != "" {
		fetchDate, err = time.Parse("2006-01-02", job.Args.Date)
		if err != nil {
			return fmt.Errorf("invalid date format %s, expected YYYY-MM-DD: %w", job.Args.Date, err)
		}
	} else {
		fetchDate = time.Now()
	}
	
	// Fetch and process the diarios
	entries, err := diarios.FetchGovernoPiauiDiarios(ctx, fetchDate, w.Uploader, w.DiarioService)
	if err != nil {
		return fmt.Errorf("failed to fetch diarios from Governo do Piauí: %w", err)
	}
	
	log.Printf("✅ Job completed successfully. Processed %d diário(s) from Governo do Piauí", len(entries))
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

// MaxRetries defines max attempts for this job
func (w *GovernoPiauiWorker) MaxRetries(job *river.Job[GovernoPiauiArgs]) int {
	return 3 // Retry up to 3 times
}

// Timeout sets the maximum execution time for this job
func (w *GovernoPiauiWorker) Timeout(job *river.Job[GovernoPiauiArgs]) time.Duration {
	return 10 * time.Minute // Allow up to 10 minutes to fetch and process
}

// CreateDiarioDosMunicipiosPeriodicJob returns a periodic job for Diario dos Municipios
func CreateDiarioDosMunicipiosPeriodicJob() *river.PeriodicJob {
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

// CreateGovernoPiauiPeriodicJob returns a periodic job for Governo do Piauí
func CreateGovernoPiauiPeriodicJob() *river.PeriodicJob {
	return river.NewPeriodicJob(
		// Run every hour
		river.PeriodicInterval(1 * time.Hour),
		
		// Args constructor function
		func() (river.JobArgs, *river.InsertOpts) {
			return GovernoPiauiArgs{}, &river.InsertOpts{
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

// ScheduleGovernoPiauiJob schedules a job to fetch the Governo do Piauí diarios to run immediately
func ScheduleGovernoPiauiJob(ctx context.Context, client *river.Client[pgx.Tx], date string) error {
	// Insert a job to run immediately
	_, err := client.Insert(ctx, GovernoPiauiArgs{
		Date: date,
	}, &river.InsertOpts{
		Queue:    "default",
		Priority: 1, // Higher number = higher priority
	})
	
	if err != nil {
		return fmt.Errorf("failed to schedule immediate governo job: %w", err)
	}
	
	log.Printf("✅ Scheduled immediate job to fetch Governo do Piauí diarios")
	return nil
}