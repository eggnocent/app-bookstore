package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type AuthorsModel struct {
	ID        uuid.UUID     `db:"id"`
	Name      string        `db:"name"`
	Bio       string        `db:"bio"`
	CreatedAt time.Time     `db:"created_at"`
	CreatedBy uuid.UUID     `db:"created_by"`
	UpdatedAt pq.NullTime   `db:"updated_at"`
	UpdatedBy uuid.NullUUID `db:"updated_by"`
}

type AuthorsRespose struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy uuid.UUID `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy uuid.UUID `json:"updated_by"`
}

func (a *AuthorsModel) Response() AuthorsRespose {
	return AuthorsRespose{
		ID:        a.ID,
		Name:      a.Name,
		Bio:       a.Bio,
		CreatedAt: a.CreatedAt,
		CreatedBy: a.CreatedBy,
		UpdatedAt: a.UpdatedAt.Time,
		UpdatedBy: a.UpdatedBy.UUID,
	}
}

func GetAllAuthors(ctx context.Context, db *sqlx.DB) ([]AuthorsModel, error) {
	query := `
		SELECT
			id, 
			name, 
			bio, 
			created_at, 
			created_by, 
			updated_at, 
			updated_by
		FROM
			authors
	`

	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var authors []AuthorsModel
	for rows.Next() {
		var author AuthorsModel
		err := rows.StructScan(&author)
		if err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}

	return authors, nil
}

func GetOneAuthors(ctx context.Context, db *sqlx.DB, id uuid.UUID) (AuthorsModel, error) {
	query := `
		SELECT 
			id,
			name,
			bio, 
            created_at, 
            created_by, 
            updated_at, 
            updated_by
		FROM
			authors
		WHERE 
			id = $1
	`

	authorRequest := AuthorsModel{}
	err := db.QueryRowxContext(ctx, query, id).StructScan(&authorRequest)
	if err != nil {
		return authorRequest, err
	}

	return authorRequest, nil
}

func (a *AuthorsModel) Insert(ctx context.Context, db *sqlx.DB) error {
	query := `
		INSERT INTO authors (
			id, name, bio, created_at, created_by
		)  VALUES (
		  $1, $2, $3, $4, $5
		) RETURNING id, created_at
	`

	err := db.QueryRowxContext(ctx, query,
		a.ID,
		a.Name,
		a.Bio,
		a.CreatedAt,
		a.CreatedBy,
	).Scan(
		&a.ID,
		&a.CreatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}

func (a *AuthorsModel) Update(ctx context.Context, db *sqlx.DB, id uuid.UUID) error {
	query := `
		UPDATE
			authors
		SET
			name = $1,
			bio = $2,
            updated_at = $3,
            updated_by = $4
		WHERE
			id = $5
		RETURNING id, name, bio, created_at, created_by, updated_at, updated_by
	`

	err := db.QueryRowxContext(ctx, query,
		a.Name,
		a.Bio,
		a.UpdatedAt.Time,
		a.UpdatedBy.UUID,
		a.ID,
	).Scan(
		&a.ID,
		&a.Name,
		&a.Bio,
		&a.CreatedAt,
		&a.CreatedBy,
		&a.UpdatedAt,
		&a.UpdatedBy,
	)

	if err != nil {
		return err
	}

	return nil
}
