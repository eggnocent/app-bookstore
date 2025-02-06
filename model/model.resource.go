package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type ResourceModel struct {
	ID          uuid.UUID     `db:"id"`
	Name        string        `db:"name"`
	Endpoint    string        `db:"endpoint"`
	Method      string        `db:"method"`
	Description string        `db:"description"`
	CreatedAt   time.Time     `db:"created_at"`
	CreatedBy   uuid.UUID     `db:"created_by"`
	UpdatedAt   pq.NullTime   `db:"updated_at"`
	UpdatedBy   uuid.NullUUID `db:"updated_by"`
}

type ResourceResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Endpoint    string    `json:"endpoint"`
	Method      string    `json:"method"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   uuid.UUID `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   uuid.UUID `json:"updated_by"`
}

func (rs *ResourceModel) Response() ResourceResponse {
	return ResourceResponse{
		ID:          rs.ID,
		Name:        rs.Name,
		Endpoint:    rs.Endpoint,
		Method:      rs.Method,
		Description: rs.Description,
		CreatedAt:   rs.CreatedAt,
		CreatedBy:   rs.CreatedBy,
		UpdatedAt:   rs.UpdatedAt.Time,
		UpdatedBy:   rs.UpdatedBy.UUID,
	}
}

func GetAllResources(ctx context.Context, db *sqlx.DB) ([]ResourceModel, error) {
	query := `
		SELECT id, name, endpoint, method, description, created_at, created_by, updated_at, updated_by
		FROM resources
	`

	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var resources []ResourceModel
	for rows.Next() {
		var resource ResourceModel
		err := rows.StructScan(&resource)
		if err != nil {
			return nil, err
		}

		resources = append(resources, resource)
	}

	return resources, nil
}

func GetOneResource(ctx context.Context, db *sqlx.DB, id uuid.UUID) (ResourceModel, error) {
	query := `
		SELECT 
			id, name, endpoint, method, description, created_at, created_by, updated_at, updated_by
		FROM
			resources
		WHERE 
			id = $1
	`

	resource := ResourceModel{}
	err := db.QueryRowxContext(ctx, query, id).StructScan(&resource)
	if err != nil {
		return resource, err
	}

	return resource, nil
}

func (rs *ResourceModel) Insert(ctx context.Context, db *sqlx.DB) error {
	query := `
		INSERT INTO resources (
			name,
			endpoint,
            method,
            description,
            created_by
		) VALUES (
			$1, $2, $3, $4, $5 
		) RETURNING id, created_at
	`

	err := db.QueryRowxContext(ctx, query,
		rs.Name,
		rs.Endpoint,
		rs.Method,
		rs.Description,
		rs.CreatedBy,
	).Scan(
		&rs.ID,
		&rs.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (rs *ResourceModel) Update(ctx context.Context, db *sqlx.DB) error {
	query := `
		UPDATE 
			resources
		SET
			name = COALESCE(NULLIF($1, ''), name),
			endpoint = COALESCE(NULLIF($2, ''), endpoint),
            method = COALESCE(NULLIF($3, ''), method),
            description = COALESCE(NULLIF($4, ''), description),
            updated_at = $5,
            updated_by = $6
		WHERE 
			id = $7
		RETURNING id, name, endpoint, method, description, created_at, created_by, updated_at, updated_by
	`

	err := db.QueryRowxContext(ctx, query,
		rs.Name,
		rs.Endpoint,
		rs.Method,
		rs.Description,
		time.Now(),
		rs.UpdatedBy.UUID,
		rs.ID,
	).Scan(
		&rs.ID,
		&rs.Name,
		&rs.Endpoint,
		&rs.Method,
		&rs.Description,
		&rs.CreatedAt,
		&rs.CreatedBy,
		&rs.UpdatedAt,
		&rs.UpdatedBy,
	)

	if err != nil {
		return err
	}

	return nil
}
