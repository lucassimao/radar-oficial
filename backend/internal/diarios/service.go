package diarios

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"radaroficial.app/internal/model"
)

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
