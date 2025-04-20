package diarios

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"slices"

	"github.com/jackc/pgx/v5/pgxpool"
	"radaroficial.app/internal/model"
)

const KNOWLEDGE_BASE_PIAUI_UUID = "a4fc4135-1a22-11f0-bf8f-4e013e2ddde4"

type DiarioService struct {
	DB *pgxpool.Pool
}

func NewDiarioService(db *pgxpool.Pool) *DiarioService {
	return &DiarioService{DB: db}
}

func (s *DiarioService) Insert(ctx context.Context, d *model.Diario) error {
	query := `
		INSERT INTO diarios (
			institution_id,
			published_at,
			last_modified_at,
			source_url,
			description,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		ON CONFLICT (institution_id, description) DO NOTHING
		RETURNING id, created_at, updated_at;
	`

	err := s.DB.QueryRow(ctx, query,
		d.InstitutionID,
		d.PublishedAt,
		d.LastModifiedAt,
		d.SourceURL,
		d.Description,
	).Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt)

	// If no rows were returned (i.e., conflict triggered), skip Scan
	if err != nil && err.Error() == "no rows in result set" {
		return nil
	}

	return err
}

// DiarioExists checks if a diario already exists in the database
func (s *DiarioService) DiarioExists(ctx context.Context, institutionID int, description string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM diarios 
			WHERE institution_id = $1 AND description = $2
		);
	`

	var exists bool
	err := s.DB.QueryRow(ctx, query, institutionID, description).Scan(&exists)

	return exists, err
}

func (s *DiarioService) GetPendingIndexing(ctx context.Context) ([]*model.Diario, error) {
	query := `
		SELECT 
			id, institution_id, published_at, last_modified_at, 
			source_url, description, 
			created_at, updated_at, indexing_submitted_at
		FROM diarios
		WHERE indexing_submitted_at IS NULL
		ORDER BY id ASC;
	`

	rows, err := s.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var diarios []*model.Diario
	for rows.Next() {
		d := &model.Diario{}
		err := rows.Scan(
			&d.ID, &d.InstitutionID, &d.PublishedAt, &d.LastModifiedAt,
			&d.SourceURL, &d.Description,
			&d.CreatedAt, &d.UpdatedAt, &d.IndexingSubmittedAt,
		)
		if err != nil {
			return nil, err
		}
		diarios = append(diarios, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return diarios, nil
}

func (s *DiarioService) MarkAsIndexingSubmitted(ctx context.Context, institutionIds []int) error {
	query := `
		UPDATE diarios
		SET 
			indexing_submitted_at = NOW(),
			updated_at = NOW()
		WHERE institution_id = ANY($1) AND indexing_submitted_at IS NULL ;
	`

	_, err := s.DB.Exec(ctx, query, institutionIds)
	return err
}

func (s *DiarioService) ReIndexKnowledgeBases(ctx context.Context) error {
	diarios, err := s.GetPendingIndexing(ctx)

	if err != nil {
		return err
	}

	if len(diarios) == 0 {
		log.Printf("✅ No pending diarios to index")
		return nil
	}

	var pendingIndexingInstitutionIds []int

	for _, d := range diarios {
		if !slices.Contains(pendingIndexingInstitutionIds, d.InstitutionID) {
			pendingIndexingInstitutionIds = append(pendingIndexingInstitutionIds, d.InstitutionID)
		}
	}

	// Map of institution IDs to knowledge base UUIDs
	institutionKBMapping := map[string][]int{
		KNOWLEDGE_BASE_PIAUI_UUID: {InstitutionIDGovernoPiaui, InstitutionIDMunicipiosPiaui},
	}

	// which knowledge bases we have already triggered the reindexing
	var reindexTriggeredForKb []string

	for _, pendingIndexingInstitutionId := range pendingIndexingInstitutionIds {

		for kbUUID, institutionIds := range institutionKBMapping {

			if slices.Contains(reindexTriggeredForKb, kbUUID) {
				continue // avoid reindexing the same kb more than once
			}

			if slices.Contains(institutionIds, pendingIndexingInstitutionId) {

				// Trigger reindex for this institution's knowledge base
				err := triggerReindex(ctx, kbUUID)

				if err == nil {
					log.Printf("✅ Reindexing triggered for Knowledge base %s", kbUUID)
					reindexTriggeredForKb = append(reindexTriggeredForKb, kbUUID)

					// Mark all diarios from these institutions as indexing submitted
					err = s.MarkAsIndexingSubmitted(ctx, institutionIds)
					if err != nil {
						log.Printf("❌ Failed to mark diarios as indexing submitted: %v", err)
						return err
					}

				} else {
					log.Printf("❌ Failed to trigger reindex for knowledge base ID %s: %v", kbUUID, err)
				}
			}
		}
	}

	return nil
}

func triggerReindex(ctx context.Context, kbUUID string) error {
	token := os.Getenv("DO_API_KEY")

	body := map[string]interface{}{
		"knowledge_base_uuid": kbUUID,
	}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.digitalocean.com/v2/gen-ai/indexing_jobs", bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		log.Printf("❌ Reindex request failed with status: %s", resp.Status)
		return err
	}

	log.Println("✅ Reindex triggered successfully")
	return nil
}
