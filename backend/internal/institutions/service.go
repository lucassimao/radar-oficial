package institutions

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InstitutionService struct {
	DB *pgxpool.Pool
}

func NewInstitutionService(db *pgxpool.Pool) *InstitutionService {
	return &InstitutionService{DB: db}
}

func (s *InstitutionService) GetStates(ctx context.Context) ([]string, error) {
	query := `SELECT distinct state FROM institutions WHERE active = true`

	rows, err := s.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var states []string

	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		states = append(states, name)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return states, nil
}
