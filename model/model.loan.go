package model

import (
	"app-bookstore/lib"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
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

func GetAllLoans(ctx context.Context, db *sqlx.DB, filter lib.Filter, dateFilter DateFilter) ([]LoansModel, error) {
	var filters []string
	var statuses []string

	if filter.Search != "" {
		filters = append(filters, fmt.Sprintf("b.title ILIKE '%%%s%%'", filter.Search))
	}

	if filter.BookID != uuid.Nil {
		filters = append(filters, fmt.Sprintf("l.book_id = '%s'", filter.BookID))
	}

	if filter.MemberID != uuid.Nil {
		filters = append(filters, fmt.Sprintf("l.member_id = '%s'", filter.MemberID))
	}

	log.Logger.Println(filter.MemberID)

	if filter.Borrowed {
		statuses = append(statuses, lib.Borrowed)
	}

	if filter.Returned {
		statuses = append(statuses, lib.Returned)
	}

	if len(statuses) > 0 {
		placeHOlders := make([]string, len(statuses))
		for i, status := range statuses {
			placeHOlders[i] = fmt.Sprintf("'%s'", status)
		}

		filters = append(filters, fmt.Sprintf("l.status IN (%s)", strings.Join(placeHOlders, ", ")))
	}

	if !dateFilter.StartDate.IsZero() && !dateFilter.EndDate.IsZero() {
		filters = append(filters, fmt.Sprintf(
			"l.loan_date BETWEEN '%s' AND '%s'",
			dateFilter.StartDate.Format("2006-01-02"),
			dateFilter.EndDate.Format("2006-01-02"),
		))
	}

	query := fmt.Sprintf(
		`
		SELECT 
			l.id, 
			l.book_id, 
			l.member_id, 
			l.loan_date, 
			l.return_date, 
			l.status, 
			l.created_at, 
			l.created_by, 
			l.updated_at, 
			l.updated_by
		FROM
			loans l
		INNER JOIN
			books b
		ON
			b.id = l.book_id
		INNER JOIN
			user_roles ur
		ON
			ur.user_id = l.member_id
		%s 
		ORDER BY l.created_at %s
		LIMIT $1 OFFSET $2
	`, lib.SearchGenerate(ctx, "AND", filters), filter.Dir)

	rows, err := db.QueryxContext(ctx, query, filter.Limit, filter.Offset)
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
