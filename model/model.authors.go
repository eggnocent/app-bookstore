package model

import (
	"app-bookstore/lib"
	"context"
	"fmt"
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

func GetAllAuthors(ctx context.Context, db *sqlx.DB, filter lib.Filter, dateFilter DateFilter) ([]AuthorsModel, error) {
	var filters []string

	if filter.Search != "" {
		filters = append(filters, fmt.Sprintf("a.name ILIKE '%%%s%%'", filter.Search))
	}

	if filter.AuthorBook != "" {
		filters = append(filters, fmt.Sprintf("b.title ILIKE '%%%s%%'", filter.AuthorBook))
	}

	if !dateFilter.StartDate.IsZero() && !dateFilter.EndDate.IsZero() {
		filters = append(filters, fmt.Sprintf(
			"a.created_at BETWEEN '%s' AND '%s'",
			dateFilter.StartDate.Format("2006-01-02"),
			dateFilter.EndDate.Format("2006-01-02"),
		))
	}
	query := fmt.Sprintf(
		`
		SELECT
			a.id, 
			a.name, 
			a.bio, 
			a.created_at, 
			a.created_by, 
			a.updated_at, 
			a.updated_by
		FROM
			authors a
		INNER JOIN
			books b
		ON
			a.id = b.author_id
		%s
		ORDER BY a.created_at %s
		LIMIT $1 OFFSET $2
	`, lib.SearchGenerate(ctx, "AND", filters), filter.Dir)

	rows, err := db.QueryxContext(ctx, query, filter.Limit, filter.Offset)
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
