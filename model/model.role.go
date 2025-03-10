package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type RoleModel struct {
	ID          uuid.UUID     `db:"id"`
	Identifier  string        `db:"identifier"`
	Description string        `db:"description"`
	CreatedAt   time.Time     `db:"created_at"`
	CreatedBy   uuid.UUID     `db:"created_by"`
	UpdatedAt   pq.NullTime   `db:"updated_at"`
	UpdatedBy   uuid.NullUUID `db:"updated_by"`
}

type RoleResponse struct {
	ID          uuid.UUID `json:"id"`
	Identifier  string    `json:"identifier"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   uuid.UUID `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   uuid.UUID `json:"updated_by"`
}

func (r *RoleModel) Response() RoleResponse {
	return RoleResponse{
		ID:          r.ID,
		Identifier:  r.Identifier,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		CreatedBy:   r.CreatedBy,
		UpdatedAt:   r.UpdatedAt.Time,
		UpdatedBy:   r.UpdatedBy.UUID,
	}
}

func GetAllRoles(ctx context.Context, db *sqlx.DB) ([]RoleModel, error) {
	query := `
		SELECT id, identifier, description, created_at, created_by, updated_at, updated_by
		FROM roles
	`

	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var roles []RoleModel
	for rows.Next() {
		var role RoleModel
		err := rows.StructScan(&role)
		if err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}
	return roles, err
}

func GetOneRole(ctx context.Context, db *sqlx.DB, id uuid.UUID) (RoleModel, error) {
	query := `
		SELECT
			id, identifier, description, created_at, created_by, updated_at, updated_by
			FROM roles
			WHERE id = $1
	`

	role := RoleModel{}
	err := db.QueryRowxContext(ctx, query, id).StructScan(&role)
	if err != nil {
		return role, err
	}

	return role, nil
}

func GetUserRoleID(ctx context.Context, db *sqlx.DB, userID uuid.UUID) (uuid.UUID, error) {
	query := `
		SELECT 
			role_id
		FROM
			user_roles
		WHERE user_id = $1
		LIMIT 1
	`

	var roleID uuid.UUID
	err := db.QueryRowxContext(ctx, query,
		userID,
	).Scan(
		&roleID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Error: user %s tidak memiliki role di db\n", userID)
			return uuid.Nil, fmt.Errorf("user does not have an assigned role")
		}

		fmt.Printf("Error: Query gagal di GetUserRoleID() - %v\n", err)
		return uuid.Nil, err
	}

	fmt.Printf("Error: User %s memiliki role %s\n", userID, roleID)
	return roleID, nil
}

func (r *RoleModel) Insert(ctx context.Context, db *sqlx.DB) error {
	query := `
		INSERT INTO roles (
			identifier,
			description,
			created_by
		) VALUES (
			$1, $2, $3 
		) RETURNING id, created_at
	`

	err := db.QueryRowxContext(ctx, query,
		r.Identifier,
		r.Description,
		r.CreatedBy,
	).Scan(
		&r.ID,
		&r.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *RoleModel) Update(ctx context.Context, db *sqlx.DB) error {
	query := `
		UPDATE 
			roles
		SET
			identifier = $1,
			description = $2,
			updated_at = $3,
			updated_by = $4
		WHERE 
			id = $5
	`

	_, err := db.ExecContext(ctx, query,
		r.Identifier,
		r.Description,
		r.UpdatedAt.Time,
		r.UpdatedBy.UUID,
		r.ID,
	)

	if err != nil {
		return err
	}
	return nil
}
