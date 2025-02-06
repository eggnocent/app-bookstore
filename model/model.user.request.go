package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserRequestModel struct {
	ID              uuid.UUID     `db:"id"`
	UserID          uuid.UUID     `db:"user_id"`
	RequestUserRole uuid.UUID     `db:"requested_role_id"`
	Status          string        `db:"status"`
	CreatedAt       time.Time     `db:"created_at"`
	CreatedBy       uuid.UUID     `db:"created_by"`
	UpdatedAt       pq.NullTime   `db:"updated_at"`
	UpdatedBy       uuid.NullUUID `db:"updated_by"`
}

type UserRequestResponse struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	RequestUserRole uuid.UUID `json:"requested_role_id"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedBy       uuid.UUID `json:"created_by"`
	UpdatedAt       time.Time `json:"updated_at"`
	UpdatedBy       uuid.UUID `json:"updated_by"`
}

func (ur *UserRequestModel) Response() UserRequestResponse {
	return UserRequestResponse{
		ID:              ur.ID,
		UserID:          ur.UserID,
		RequestUserRole: ur.RequestUserRole,
		Status:          ur.Status,
		CreatedAt:       ur.CreatedAt,
		CreatedBy:       ur.CreatedBy,
		UpdatedAt:       ur.UpdatedAt.Time,
		UpdatedBy:       ur.UpdatedBy.UUID,
	}
}

func CheckRoleExist(ctx context.Context, db *sqlx.DB, roleID uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(
						SELECT 
							1 
						FROM 
							roles 
						WHERE 
							id = $1
						)`

	err := db.QueryRowxContext(ctx, query, roleID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func GetAllUserRequest(ctx context.Context, db *sqlx.DB) ([]UserRequestModel, error) {
	query := `
		SELECT
		    id, 
            user_id, 
            requested_role_id, 
            status, 
            created_at,
			created_by,
			updated_at,
			updated_by
		FROM
		    user_requests
	`

	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var userRequests []UserRequestModel
	for rows.Next() {
		var userRequest UserRequestModel
		err := rows.StructScan(&userRequest)
		if err != nil {
			return nil, err
		}

		userRequests = append(userRequests, userRequest)
	}
	return userRequests, err
}

func GetOneUserRequest(ctx context.Context, db *sqlx.DB, id uuid.UUID) (UserRequestModel, error) {
	query := `
		SELECT
			id, user_id, requested_role_id, status, created_at, created_by, updated_at, updated_by
		FROM
			user_requests
		WHERE id = $1
	`

	userRequest := UserRequestModel{}
	err := db.QueryRowxContext(ctx, query, id).StructScan(&userRequest)
	if err != nil {
		return userRequest, err
	}

	return userRequest, nil
}

func (ur *UserRequestModel) Insert(ctx context.Context, db *sqlx.DB) error {
	query := `
		INSERT INTO user_requests (
			id, 
			user_id, 
			requested_role_id, 
			status, 
			created_at, 
			created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6 
		) RETURNING id, created_at
	`

	err := db.QueryRowxContext(ctx, query,
		ur.ID,
		ur.UserID,
		ur.RequestUserRole,
		ur.Status,
		ur.CreatedAt,
		ur.CreatedBy,
	).Scan(
		&ur.ID,
		&ur.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRequestModel) Update(ctx context.Context, db *sqlx.DB, id uuid.UUID) error {
	query := `
		UPDATE
			user_requests
		SET
			status = $1,
			updated_at = $2,
			updated_by = $3
		WHERE
			id = $4
		RETURNING id, user_id, requested_role_id, status, created_at, created_by, updated_at, updated_by
	`

	err := db.QueryRowxContext(ctx, query,
		ur.Status,
		ur.UpdatedAt.Time,
		ur.UpdatedBy.UUID,
		ur.ID,
	).Scan(
		&ur.ID,
		&ur.UserID,
		&ur.RequestUserRole,
		&ur.Status,
		&ur.CreatedAt,
		&ur.CreatedBy,
		&ur.UpdatedAt,
		&ur.UpdatedBy,
	)

	if err != nil {
		return err
	}

	return nil
}
