package model

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserModel struct {
	ID        uuid.UUID     `db:"id"`
	Username  string        `db:"username"`
	Password  string        `db:"password"`
	CreatedAt time.Time     `db:"created_at"`
	CreatedBy uuid.UUID     `db:"created_by"`
	UpdatedAt pq.NullTime   `db:"updated_at"`
	UpdatedBy uuid.NullUUID `db:"updated_by"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy uuid.UUID `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy uuid.UUID `json:"updated_by"`
}

func (u *UserModel) Response() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		CreatedBy: u.CreatedBy,
		UpdatedAt: u.UpdatedAt.Time,
		UpdatedBy: u.UpdatedBy.UUID,
	}
}

func GetUserByUsername(ctx context.Context, db *sqlx.DB, username string) (*UserModel, error) {
	var user UserModel
	query := `
		SELECT 
			id,
			username,
			password
		FROM
			users
		WHERE
			username = $1
	`

	err := db.GetContext(ctx, &user, query, username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(ctx context.Context, db *sqlx.DB, userID uuid.UUID) (*UserModel, error) {
	var user UserModel
	query := `
		SELECT
			id,
			password
		FROM
			users
		WHERE
			id = $1
	`

	err := db.GetContext(ctx, &user, query, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func (u *UserModel) Insert(ctx context.Context, db *sqlx.DB) error {
	query := `
		INSERT INTO users (
			username,
			password,
			created_by
		) VALUES (
			$1, $2, $3
		) RETURNING id, created_at
	`

	err := db.QueryRowxContext(ctx, query,
		u.Username,
		u.Password,
		u.CreatedBy,
	).Scan(
		&u.ID,
		&u.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func Update(ctx context.Context, db *sqlx.DB, userID uuid.UUID, hashPass string) error {
	query := `
		UPDATE
			users
		SET
			password = $1, updated_at = $2
		WHERE
			id = $3
	`

	_, err := db.ExecContext(ctx, query, hashPass, time.Now(), userID)
	if err != nil {
		return err
	}
	return nil
}
