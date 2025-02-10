package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type CategoryModel struct {
	ID        uuid.UUID     `db:"id"`
	Name      string        `db:"name"`
	CreatedAt time.Time     `db:"created_at"`
	CreatedBy uuid.UUID     `db:"created_by"`
	UpdatedAt pq.NullTime   `db:"updated_at"`
	UpdatedBy uuid.NullUUID `db:"updated_by"`
}

type CategoryResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy uuid.UUID `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy uuid.UUID `json:"updated_by"`
}

func (c *CategoryModel) Response() CategoryResponse {
	return CategoryResponse{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		CreatedBy: c.CreatedBy,
		UpdatedAt: c.UpdatedAt.Time,
		UpdatedBy: c.UpdatedBy.UUID,
	}
}

func GetAllCategory(ctx context.Context, db *sqlx.DB) ([]CategoryModel, error) {
	query := `
		SELECT
			id, name, created_at, created_by, updated_at, updated_by
		FROM
			categories
	`

	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var categories []CategoryModel
	for rows.Next() {
		var category CategoryModel
		err := rows.StructScan(&category)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func GetOneCategories(ctx context.Context, db *sqlx.DB, id uuid.UUID) (CategoryModel, error) {
	query := `
		SELECT
			id, name, created_at, created_by, updated_at, updated_by
        FROM
			categories
		WHERE 
			id = $1
	`

	categoryRequest := CategoryModel{}
	err := db.QueryRowxContext(ctx, query, id).StructScan(&categoryRequest)
	if err != nil {
		return categoryRequest, err
	}

	return categoryRequest, nil
}

func (c *CategoryModel) Insert(ctx context.Context, db *sqlx.DB) error {
	query := `
        INSERT INTO categories (
            id, name, created_at, created_by
        ) VALUES (
            $1, $2, $3, $4
        ) RETURNING id
    `

	err := db.QueryRowxContext(ctx, query,
		c.ID,
		c.Name,
		c.CreatedAt,
		c.CreatedBy,
	).Scan(
		&c.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *CategoryModel) Update(ctx context.Context, db *sqlx.DB, id uuid.UUID) error {
	query := `
		UPDATE
			categories
		SET
			name = $1, 
			updated_at = $2, 
			updated_by = $3
        WHERE
			id = $4
		RETURNING id, name, created_at, created_by, updated_at, updated_by
	`

	err := db.QueryRowxContext(ctx, query,
		c.Name,
		c.UpdatedAt.Time,
		c.UpdatedBy.UUID,
		c.ID,
	).Scan(
		&c.ID,
		&c.Name,
		&c.CreatedAt,
		&c.CreatedBy,
		&c.UpdatedAt,
		&c.UpdatedBy,
	)

	if err != nil {
		return err
	}

	return nil
}
