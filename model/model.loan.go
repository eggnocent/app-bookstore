package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type LoansModel struct {
	ID         uuid.UUID     `db:"id"`
	BookID     uuid.UUID     `db:"book_id"`
	MemberID   uuid.UUID     `db:"member_id"`
	LoanDate   time.Time     `db:"loan_date"`
	ReturnDate pq.NullTime   `db:"return_date"`
	Status     string        `db:"status"`
	CreatedAt  time.Time     `db:"created_at"`
	CreatedBy  uuid.UUID     `db:"created_by"`
	UpdatedAt  pq.NullTime   `db:"updated_at"`
	UpdatedBy  uuid.NullUUID `db:"updated_by"`
}

type LoansResponse struct {
	ID         uuid.UUID `db:"id"`
	BookID     uuid.UUID `db:"book_id"`
	MemberID   uuid.UUID `db:"member_id"`
	LoanDate   time.Time `db:"loan_date"`
	ReturnDate time.Time `db:"return_date"`
	Status     string    `db:"status"`
	CreatedAt  time.Time `db:"created_at"`
	CreatedBy  uuid.UUID `db:"created_by"`
	UpdatedAt  time.Time `db:"updated_at"`
	UpdatedBy  uuid.UUID `db:"updated_by"`
}

func (l *LoansModel) Response() LoansResponse {
	return LoansResponse{
		ID:         l.ID,
		BookID:     l.BookID,
		MemberID:   l.MemberID,
		LoanDate:   l.LoanDate,
		ReturnDate: l.ReturnDate.Time,
		Status:     l.Status,
		CreatedAt:  l.CreatedAt,
		CreatedBy:  l.CreatedBy,
		UpdatedAt:  l.UpdatedAt.Time,
		UpdatedBy:  l.UpdatedBy.UUID,
	}
}

func GetAllLoans(ctx context.Context, db *sqlx.DB) ([]LoansModel, error) {
	query := `
		SELECT 
			id, book_id, member_id, loan_date, return_date, status, created_at, created_by, updated_at, updated_by
		FROM
			loans
	`

	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var loans []LoansModel
	for rows.Next() {
		var loan LoansModel
		err := rows.StructScan(&loan)
		if err != nil {
			return nil, err
		}
		loans = append(loans, loan)
	}

	return loans, nil
}

func GetOneLoans(ctx context.Context, db *sqlx.DB, id uuid.UUID) (LoansModel, error) {
	query := `
		SELECT
			id, book_id, member_id, loan_date, return_date, status, created_at, created_by, updated_at, updated_by
		FROM 
			loans
		WHERE
			id = $1
	`

	loan := LoansModel{}
	err := db.QueryRowxContext(ctx, query, id).StructScan(&loan)
	if err != nil {
		return loan, err
	}
	return loan, nil
}

func (l *LoansModel) Insert(ctx context.Context, db *sqlx.DB) error {
	query := `
		INSERT INTO loans (
			id, book_id, member_id, loan_date, return_date, status, created_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8 
		) RETURNING id, created_at
	`

	err := db.QueryRowxContext(ctx, query,
		l.ID,
		l.BookID,
		l.MemberID,
		l.LoanDate,
		l.ReturnDate,
		l.Status,
		l.CreatedAt,
		l.CreatedBy,
	).Scan(
		&l.ID,
		&l.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (l *LoansModel) Update(ctx context.Context, db *sqlx.DB, id uuid.UUID) error {
	query := `
		UPDATE
			loans
		SET
			status = $1,
			return_date = $2,
			updated_at = $3,
			updated_by = $4
		WHERE
			id = $5
		RETURNING id, book_id, member_id, loan_date, return_date, status, created_at, created_by, updated_at, updated_by
	`

	err := db.QueryRowxContext(ctx, query,
		l.Status,
		l.ReturnDate,
		l.UpdatedAt.Time,
		l.UpdatedBy.UUID,
		l.ID,
	).Scan(
		&l.ID,
		&l.BookID,
		&l.MemberID,
		&l.LoanDate,
		&l.ReturnDate,
		&l.Status,
		&l.CreatedAt,
		&l.CreatedBy,
		&l.UpdatedAt,
		&l.UpdatedBy,
	)

	if err != nil {
		return err
	}
	return nil
}
