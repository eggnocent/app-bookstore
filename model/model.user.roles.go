package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserRoleModel struct {
	ID        uuid.UUID     `db:"id"`
	UserID    uuid.UUID     `db:"user_id"`
	RoleID    uuid.UUID     `db:"role_id"`
	CreatedAt time.Time     `db:"created_at"`
	CreatedBy uuid.UUID     `db:"created_by"`
	UpdatedAt pq.NullTime   `db:"updated_at"`
	UpdatedBy uuid.NullUUID `db:"updated_by"`
}

type UserRoleResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	RoleID    uuid.UUID `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy uuid.UUID `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy uuid.UUID `json:"updated_by"`
}

func (uro *UserRoleModel) Response() UserRoleResponse {
	return UserRoleResponse{
		ID:        uro.ID,
		UserID:    uro.UserID,
		RoleID:    uro.RoleID,
		CreatedAt: uro.CreatedAt,
		CreatedBy: uro.CreatedBy,
		UpdatedAt: uro.UpdatedAt.Time,
		UpdatedBy: uro.UpdatedBy.UUID,
	}
}

func GetAllUserRoles(ctx context.Context, db *sqlx.DB) ([]UserRoleModel, error) {
	query := `
		SELECT 
			id, 
			user_id, 
			role_id, 
			created_at, 
			created_by, 
			updated_at, 
			updated_by
		FROM 
			user_roles
	`

	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var userRoles []UserRoleModel
	for rows.Next() {
		var userRole UserRoleModel
		err := rows.StructScan(&userRole)
		if err != nil {
			return nil, err
		}
		userRoles = append(userRoles, userRole)
	}

	return userRoles, err
}

func GetOneUserRole(ctx context.Context, db *sqlx.DB, id uuid.UUID) (UserRoleModel, error) {
	query := `
		SELECT 
			id, 
			user_id, 
			role_id, 
			created_at, 
			created_by, 
			updated_at
		FROM 
			user_roles
		WHERE
		    id = $1
	`

	userRole := UserRoleModel{}
	err := db.QueryRowxContext(ctx, query, id).StructScan(&userRole)
	if err != nil {
		return userRole, err
	}
	return userRole, nil
}

func (uro *UserRoleModel) Insert(ctx context.Context, db *sqlx.DB) error {
	query := `
		INSERT INTO user_roles (
			id, 
			user_id, 
			role_id, 
			created_at, 
			created_by 
		) VALUES (
			$1, $2, $3, $4, $5
		) RETURNING id, created_at
	`

	err := db.QueryRowxContext(ctx, query,
		uro.ID,
		uro.UserID,
		uro.RoleID,
		uro.CreatedAt,
		uro.CreatedBy,
	).Scan(
		&uro.ID,
		&uro.CreatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}

func (uro *UserRoleModel) Update(ctx context.Context, db *sqlx.DB) error {
	query := `
		UPDATE
			user_roles
        SET
			user_id = $1,
			role_id = $2,
            updated_at = $3,
            updated_by = $4
		WHERE
			id = $5
		RETURNING id, user_id, role_id, created_at, created_by, updated_at, updated_by
	`

	err := db.QueryRowxContext(ctx, query,
		uro.UserID,
		uro.RoleID,
		uro.UpdatedAt.Time,
		uro.UpdatedBy.UUID,
		uro.ID,
	).Scan(
		&uro.ID,
		&uro.UserID,
		&uro.RoleID,
		&uro.CreatedAt,
		&uro.CreatedBy,
		&uro.UpdatedAt,
		&uro.UpdatedBy,
	)

	if err != nil {
		return err
	}

	return nil
}

func GetGuestRoleID(ctx context.Context, db *sqlx.DB) (uuid.UUID, error) {
	var roleID uuid.UUID
	query := `
		SELECT
			id
		FROM
			roles
		WHERE
			identifier = 'guest'
		LIMIT 1
	`

	err := db.QueryRowxContext(ctx, query).Scan(&roleID)
	if err != nil {
		return uuid.Nil, err
	}

	return roleID, nil
}

func (uro *UserRoleModel) AssignGuestRoles(ctx context.Context, db *sqlx.DB) error {
	query := `
		INSERT INTO user_roles(
			id, 
			user_id,
			role_id,
			created_at,
			created_by
		) VALUES (
			$1, $2, $3, $4, $5 
		) RETURNING id, created_at
	`

	err := db.QueryRowxContext(ctx, query,
		uro.ID,
		uro.UserID,
		uro.RoleID,
		uro.CreatedAt,
		uro.CreatedBy,
	).Scan(
		&uro.ID,
		&uro.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}
