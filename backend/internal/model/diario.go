package model

import "time"

type Diario struct {
	ID                  int        `db:"id"`
	InstitutionID       int        `db:"institution_id"`
	PublishedAt         *time.Time `db:"published_at"`
	LastModifiedAt      *time.Time `db:"last_modified_at"`
	SourceURL           string     `db:"source_url"`
	Description         *string    `db:"description"`
	CreatedAt           time.Time  `db:"created_at"`
	UpdatedAt           time.Time  `db:"updated_at"`
	IndexingSubmittedAt *time.Time `db:"indexing_submitted_at"`
}
