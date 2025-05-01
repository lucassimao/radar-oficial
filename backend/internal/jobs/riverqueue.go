package jobs

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"radaroficial.app/internal/diarios"
	"radaroficial.app/internal/storage"
)

// RiverClient wraps the River configuration and client
type RiverClient struct {
	Client   *river.Client[pgx.Tx]
	RiverDB  *riverpgxv5.Driver
	Workers  *river.Workers
	StopFunc func(context.Context) error

	// Worker instances
	DiarioWorker       *DiarioWorker
	GovernoPiauiWorker *GovernoPiauiWorker
}

// NewRiverClient creates and configures a new River client and workers
func NewRiverClient(ctx context.Context, db *pgxpool.Pool) (*RiverClient, error) {
	// Create the River database driver using pgx
	riverDB := riverpgxv5.New(db)

	// Create the job service dependencies
	diarioService := diarios.NewInstitutionService(db)

	// Initialize storage uploader
	uploader, err := storage.NewSpacesUploader("radar-oficial-diarios-piaui")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage uploader: %w", err)
	}

	// Create our job workers
	diarioWorker := NewDiarioWorker(diarioService, uploader)
	governoWorker := NewGovernoPiauiWorker(diarioService, uploader)

	// Create a workers registry
	workers := river.NewWorkers()

	// Register all workers
	river.AddWorker(workers, diarioWorker)
	river.AddWorker(workers, governoWorker)

	// Add periodic jobs
	periodicJobs := []*river.PeriodicJob{
		CreateDiarioDosMunicipiosPeriodicJob(),
		CreateGovernoPiauiPeriodicJob(),
	}

	// Create the River client config
	riverConfig := river.Config{
		Queues: map[string]river.QueueConfig{
			"default": {MaxWorkers: 5},
		},
		Workers:      workers,
		PeriodicJobs: periodicJobs,
	}

	// Create the River client
	client, err := river.NewClient(riverDB, &riverConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create River client: %w", err)
	}

	// Start the River client
	if err := client.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start River client: %w", err)
	}

	log.Printf("✅ River queue client started successfully")

	return &RiverClient{
		Client:             client,
		RiverDB:            riverDB,
		Workers:            workers,
		StopFunc:           client.Stop,
		DiarioWorker:       diarioWorker,
		GovernoPiauiWorker: governoWorker,
	}, nil
}

// ScheduleInitialJobs sets up the initial job schedules when the system starts
func (r *RiverClient) ScheduleInitialJobs(ctx context.Context) error {
	// Schedule immediate jobs for testing
	if err := ScheduleDiarioDosMunicipiosJob(ctx, r.Client); err != nil {
		return fmt.Errorf("failed to schedule immediate diario job: %w", err)
	}

	if err := ScheduleGovernoPiauiJob(ctx, r.Client, ""); err != nil {
		return fmt.Errorf("failed to schedule immediate governo job: %w", err)
	}

	log.Printf("✅ Initial jobs scheduled successfully")
	return nil
}

// Shutdown safely stops the River client
func (r *RiverClient) Shutdown(ctx context.Context) {
	// Shutdown the client
	if r.StopFunc != nil {
		if err := r.StopFunc(ctx); err != nil {
			log.Printf("⚠️ Error stopping River client: %v", err)
		} else {
			log.Printf("✅ River client stopped successfully")
		}
	}
}
