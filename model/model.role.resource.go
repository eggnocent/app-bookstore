package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type RoleResourceModel struct {
	ID         uuid.UUID     `db:"id"`
	RoleID     uuid.UUID     `db:"role_id"`
	ResourceID uuid.UUID     `db:"resource_id"`
	Method     string        `db:"method"`
	IsActive   bool          `db:"is_active"`
	CreatedAt  time.Time     `db:"created_at"`
	CreatedBy  uuid.UUID     `db:"created_by"`
	UpdatedAt  pq.NullTime   `db:"updated_at"`
	UpdatedBy  uuid.NullUUID `db:"updated_by"`
}

type RoleResourceResponse struct {
	ID         uuid.UUID `json:"id"`
	RoleID     uuid.UUID `json:"role_id"`
	ResourceID uuid.UUID `json:"resource_id"`
	Method     string    `json:"method"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  uuid.UUID `json:"created_by"`
	UpdatedAt  time.Time `json:"updated_at"`
	UpdatedBy  uuid.UUID `json:"updated_by"`
}

func (rr *RoleResourceModel) Response() RoleResourceResponse {
	return RoleResourceResponse{
		ID:         rr.ID,
		RoleID:     rr.RoleID,
		ResourceID: rr.ResourceID,
		Method:     rr.Method,
		IsActive:   rr.IsActive,
		CreatedAt:  rr.CreatedAt,
		CreatedBy:  rr.CreatedBy,
		UpdatedAt:  rr.UpdatedAt.Time,
		UpdatedBy:  rr.UpdatedBy.UUID,
	}
}

func GetAllRoleResource(ctx context.Context, db *sqlx.DB) ([]RoleResourceModel, error) {
	query := `
		SELECT id, role_id, resource_id, method, is_active, created_at, created_by, updated_at, updated_by
		FROM role_resources
	`

	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var roleResources []RoleResourceModel
	for rows.Next() {
		var roleResource RoleResourceModel
		err := rows.StructScan(&roleResource)
		if err != nil {
			return nil, err
		}

		roleResources = append(roleResources, roleResource)
	}

	return roleResources, nil
}

func GetOneRoleResource(ctx context.Context, db *sqlx.DB, id uuid.UUID) (RoleResourceModel, error) {
	query := `
		SELECT id, role_id, resource_id, method, is_active, created_at, created_by, updated_at, updated_by
		FROM role_resources
		WHERE id = $1
	`

	roleResource := RoleResourceModel{}
	err := db.QueryRowxContext(ctx, query, id).StructScan(&roleResource)
	if err != nil {
		return roleResource, err
	}

	return roleResource, nil
}

func (rr *RoleResourceModel) Insert(ctx context.Context, db *sqlx.DB) error {
	query := `
		INSERT INTO role_resources (
			role_id,
			resource_id,
			method,
			is_active, 
			created_by
		) VALUES (
			$1, $2, $3, $4, $5 
		) RETURNING id, created_at
	`

	err := db.QueryRowxContext(ctx, query,
		rr.RoleID,
		rr.ResourceID,
		rr.Method,
		rr.IsActive,
		rr.CreatedBy,
	).Scan(
		&rr.ID,
		&rr.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (rr *RoleResourceModel) Update(ctx context.Context, db *sqlx.DB) error {
	query := `
		UPDATE
			role_resources
		SET
			role_id = $1,
			resource_id = $2,
			method = $3,
			is_active = $4,
			updated_at = $5,
			updated_by = $6
		WHERE
			id = $7
		RETURNING id, role_id, resource_id, method, is_active, created_at, created_by, updated_at, updated_by
	`

	err := db.QueryRowxContext(ctx, query,
		rr.RoleID,
		rr.ResourceID,
		rr.Method,
		rr.IsActive,
		time.Now(),
		rr.UpdatedBy.UUID,
		rr.ID,
	).Scan(
		&rr.ID,
		&rr.RoleID,
		&rr.ResourceID,
		&rr.Method,
		&rr.IsActive,
		&rr.CreatedAt,
		&rr.CreatedBy,
		&rr.UpdatedAt,
		&rr.UpdatedBy,
	)

	if err != nil {
		return err
	}

	return nil
}
