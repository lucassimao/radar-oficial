package diarios

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"radaroficial.app/internal/model"
)

const KNOWLEDGE_BASE_PIAUI_UUID = "a4fc4135-1a22-11f0-bf8f-4e013e2ddde4"

type DiarioService struct {
	DB *pgxpool.Pool
}

func NewInstitutionService(db *pgxpool.Pool) *DiarioService {
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

// Exists checks if a diario already exists in the database
func (s *DiarioService) Exists(ctx context.Context, institutionID int, description string) (bool, error) {
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
