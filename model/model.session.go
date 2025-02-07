package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type SessionModel struct {
	ID         uuid.UUID     `db:"id"`
	UserID     uuid.UUID     `db:"user_id"`
	Token      string        `db:"token"`
	Expiration time.Time     `db:"expiration"`
	CreatedAt  time.Time     `db:"created_at"`
	CreatedBy  uuid.UUID     `db:"created_by"`
	UpdatedAt  pq.NullTime   `db:"updated_at"`
	UpdatedBy  uuid.NullUUID `db:"updated_by"`
}

type SessionResponse struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	Token      string    `json:"token"`
	Expiration time.Time `json:"expiration"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  uuid.UUID `json:"created_by"`
	UpdatedAt  time.Time `json:"updated_at"`
	UpdatedBy  uuid.UUID `json:"updated_by"`
}

func (s *SessionModel) Response() SessionResponse {
	return SessionResponse{
		ID:         s.ID,
		UserID:     s.UserID,
		Token:      s.Token,
		Expiration: s.Expiration,
		CreatedAt:  s.CreatedAt,
		CreatedBy:  s.CreatedBy,
		UpdatedAt:  s.UpdatedAt.Time,
		UpdatedBy:  s.UpdatedBy.UUID,
	}
}

func (s *SessionModel) Insert(ctx context.Context, db *sqlx.DB) error {
	query := `
		INSERT INTO sessions (user_id, token, expiration, created_at, created_by)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at
	`

	err := db.QueryRowxContext(ctx, query,
		s.UserID,
		s.Token,
		s.Expiration,
		s.CreatedAt,
		s.CreatedBy,
	).Scan(
		&s.ID,
		&s.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func CheckSessionExists(ctx context.Context, db *sqlx.DB, token string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(SELECT 1 FROM sessions WHERE token = $1 AND expiration > NOW())
	`
	err := db.QueryRowxContext(ctx, query, token).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func DeleteSessionByToken(ctx context.Context, db *sqlx.DB, token string) error {
	query := `DELETE FROM sessions WHERE token = $1`
	_, err := db.ExecContext(ctx, query, token)
	return err
}
