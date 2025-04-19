package jobs

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// VerifyRiverTables checks if the River Queue tables exist
func VerifyRiverTables(ctx context.Context, db *pgxpool.Pool) error {
	log.Printf("🔄 Verifying River Queue database tables...")
	
	// Check if the main river_jobs table exists
	var tableExists bool
	err := db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'river_jobs'
		)
	`).Scan(&tableExists)
	
	if err != nil {
		return fmt.Errorf("failed to check if river_jobs table exists: %w", err)
	}
	
	if !tableExists {
		return fmt.Errorf("required river_jobs table doesn't exist - please run migrations")
	}
	
	log.Printf("✅ River Queue tables verified")
	return nil
}

// CleanRiverTables removes old completed and failed jobs
func CleanRiverTables(ctx context.Context, db *pgxpool.Pool) error {
	// Start a transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	
	// Delete finished jobs older than 7 days
	result, err := tx.Exec(ctx, `
		DELETE FROM river_jobs
		WHERE state = 'completed'
		AND finished_at < NOW() - INTERVAL '7 days'
	`)
	if err != nil {
		return fmt.Errorf("failed to clean up completed jobs: %w", err)
	}
	
	rowsAffected := result.RowsAffected()
	
	// Delete discarded jobs older than 30 days
	result, err = tx.Exec(ctx, `
		DELETE FROM river_jobs
		WHERE state = 'discarded'
		AND discarded_at < NOW() - INTERVAL '30 days'
	`)
	if err != nil {
		return fmt.Errorf("failed to clean up discarded jobs: %w", err)
	}
	
	rowsAffected += result.RowsAffected()
	
	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	log.Printf("✅ Cleaned up %d old River Queue jobs", rowsAffected)
	return nil
}