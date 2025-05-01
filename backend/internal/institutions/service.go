package institutions

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"radaroficial.app/internal/model"
)

type InstitutionService struct {
	DB *pgxpool.Pool
}

func NewInstitutionService(db *pgxpool.Pool) *InstitutionService {
	return &InstitutionService{DB: db}
}

func (s *InstitutionService) GetAll(ctx context.Context) ([]*model.Institution, error) {
	query := `
		SELECT id,name,slug,type,state,city,source_url,active,created_at,updated_at 
		FROM institutions WHERE active = true
	`

	rows, err := s.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var institutions []*model.Institution

	for rows.Next() {
		i := &model.Institution{}
		err := rows.Scan(&i.ID, &i.Name, &i.Slug, &i.Type, &i.State, &i.City, &i.SourceUrl, &i.Active, &i.CreatedAt, &i.UpdatedAt)
		if err != nil {
			return nil, err
		}
		institutions = append(institutions, i)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return institutions, nil
}
