package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PasswordResetModel struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}

type PasswordResetResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (pr *PasswordResetModel) Response() PasswordResetResponse {
	return PasswordResetResponse{
		ID:        pr.ID,
		UserID:    pr.UserID,
		Token:     pr.Token,
		ExpiresAt: pr.ExpiresAt,
		CreatedAt: pr.CreatedAt,
	}
}

func (pr PasswordResetModel) Insert(ctx context.Context, db *sqlx.DB) error {
	query := `
		INSERT INTO password_reset (
			user_id, token, expires_at
		) VALUES (
			$1, $2, $3 
		) RETURNING id, created_at
	`

	err := db.QueryRowxContext(ctx, query,
		pr.UserID,
		pr.Token,
		pr.ExpiresAt,
	).Scan(
		&pr.ID,
		&pr.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func ValidateResetToken(ctx context.Context, db *sqlx.DB, token string) (uuid.UUID, error) {
	var userID uuid.UUID
	query := `
		SELECT user_id FROM password_reset
		WHERE token = $1 AND expires_at > NOW()
		LIMIT 1
	`

	err := db.QueryRowxContext(ctx, query, token).Scan(&userID)
	if err != nil {
		return uuid.Nil, err
	}
	return userID, nil
}

func DeleteResetToken(ctx context.Context, db *sqlx.DB, token string) error {
	query := `
		DELETE FROM password_reset
		WHERE token = $1
	`

	_, err := db.ExecContext(ctx, query, token)
	return err
}

func UpdatePassword(ctx context.Context, db *sqlx.DB, userID uuid.UUID, hashedPassword string) error {
	query := `
		UPDATE users
		SET password = $1,
		updated_at = NOW()
		WHERE id = $2
	`

	_, err := db.ExecContext(ctx, query, hashedPassword, userID)
	return err
}
