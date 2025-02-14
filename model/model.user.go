package model

import (
	"app-bookstore/lib"
	"context"
	"errors"
	"fmt"
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

type DateFilter struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

func GetAllUser(ctx context.Context, db *sqlx.DB, filter lib.Filter, statusUser string, dateFilter DateFilter) ([]UserModel, error) {
	var filters []string

	if filter.Search != "" {
		filters = append(filters, fmt.Sprintf("u.username ILIKE '%%%s%%'", filter.Search))
	}

	if filter.RoleID != uuid.Nil {
		filters = append(filters, fmt.Sprintf("ur.role_id = '%s'", filter.RoleID))
	}

	if statusUser != "" {
		filters = append(filters, fmt.Sprintf("u.status = '%s'", statusUser))
	}

	if !dateFilter.StartDate.IsZero() && !dateFilter.EndDate.IsZero() {
		filters = append(filters, fmt.Sprintf(
			"u.created_at BETWEEN '%s' AND '%s'",
			dateFilter.StartDate.Format("2006-01-02"),
			dateFilter.EndDate.Format("2006-01-02"),
		))
	}

	query := fmt.Sprintf(
		`
		SELECT 
			u.id, 
			u.username, 
			u.created_at, 
			u.created_by, 
			u.updated_at, 
			u.updated_by 
		FROM 
			users u
		INNER JOIN 
			user_roles ur 
		ON 
			u.id = ur.user_id
		%s
		ORDER BY u.created_at %s
		LIMIT $1 OFFSET $2
	`, lib.SearchGenerate(ctx, "AND", filters), filter.Dir)
	rows, err := db.QueryxContext(ctx, query, filter.Limit, filter.Offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []UserModel
	for rows.Next() {
		var user UserModel
		err := rows.StructScan(&user)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	return users, err
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
