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
)

type BookModel struct {
	ID            uuid.UUID     `db:"id"`
	Title         string        `db:"title"`
	AuthorID      uuid.UUID     `db:"author_id"`
	PublisherID   uuid.NullUUID `db:"publisher_id"`
	CategoryID    uuid.NullUUID `db:"category_id"`
	PublishedYear int           `db:"published_year"`
	ISBN          string        `db:"isbn"`
	Status        string        `db:"status"`
	AccessLevel   string        `db:"access_level"`
	CreatedAt     time.Time     `db:"created_at"`
	CreatedBy     uuid.UUID     `db:"created_by"`
	UpdatedAt     pq.NullTime   `db:"updated_at"`
	UpdatedBy     uuid.NullUUID `db:"updated_by"`
}

type BookResponse struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	AuthorID      uuid.UUID `json:"author_id"`
	PublisherID   uuid.UUID `json:"publisher_id"`
	CategoryID    uuid.UUID `json:"category_id"`
	PublishedYear int       `json:"published_year"`
	ISBN          string    `json:"isbn"`
	Status        string    `json:"status"`
	AccessLevel   string    `json:"access_level"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     uuid.UUID `json:"created_by"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
	UpdatedBy     uuid.UUID `json:"updated_by,omitempty"`
}

func (b *BookModel) Response() BookResponse {
	return BookResponse{
		ID:            b.ID,
		Title:         b.Title,
		AuthorID:      b.AuthorID,
		PublisherID:   b.PublisherID.UUID,
		CategoryID:    b.CategoryID.UUID,
		PublishedYear: b.PublishedYear,
		ISBN:          b.ISBN,
		Status:        b.Status,
		AccessLevel:   b.AccessLevel,
		CreatedAt:     b.CreatedAt,
		CreatedBy:     b.CreatedBy,
		UpdatedAt:     b.UpdatedAt.Time,
		UpdatedBy:     b.UpdatedBy.UUID,
	}
}

func GetAllBooks(ctx context.Context, db *sqlx.DB, filter lib.Filter, dateFilter DateFilter) ([]BookModel, error) {
	var filters []string
	var statuses []string

	if filter.Search != "" {
		filters = append(filters, fmt.Sprintf("b.title ILIKE '%%%s%%'", filter.Search))
	}

	if filter.AuthorID != uuid.Nil {
		filters = append(filters, fmt.Sprintf("b.author_id = '%s'", filter.AuthorID))
	}

	if filter.PublisherID != uuid.Nil {
		filters = append(filters, fmt.Sprintf("b.publisher_id = '%s'", filter.PublisherID))
	}

	if filter.CategoryID != uuid.Nil {
		filters = append(filters, fmt.Sprintf("b.category_id = '%s'", filter.CategoryID))
	}

	if filter.PublishedYear > 0 {
		filters = append(filters, fmt.Sprintf("b.published_year = %d", filter.PublishedYear))
	}

	if filter.Status != "" {
		filters = append(filters, fmt.Sprintf("b.status = '%s'", filter.Status))
	}

	if filter.AccessLevel != "" {
		filters = append(filters, fmt.Sprintf("b.access_level = '%s'", filter.AccessLevel))
	}

	if filter.Available {
		statuses = append(statuses, lib.Available)
	}

	if filter.Borrowed {
		statuses = append(statuses, lib.Borrowed)
	}

	if filter.Public {
		statuses = append(statuses, lib.Public)
	}

	if filter.MemberOnly {
		statuses = append(statuses, lib.MemberOnly)
	}

	if filter.AdminOnly {
		statuses = append(statuses, lib.AdminOnly)
	}

	if len(statuses) > 0 {
		placeHolders := make([]string, len(statuses))
		for i, status := range statuses {
			placeHolders[i] = fmt.Sprintf("'%s'", status)
		}

		filters = append(filters, fmt.Sprintf("b.status IN (%s)", strings.Join(placeHolders, ", ")))
	}

	if !dateFilter.StartDate.IsZero() && !dateFilter.EndDate.IsZero() {
		filters = append(filters, fmt.Sprintf(
			"b.created_at BETWEEN '%s' AND '%s'",
			dateFilter.StartDate.Format("2006-01-02"),
			dateFilter.EndDate.Format("2006-01-02"),
		))
	}

	query := fmt.Sprintf(
		`
		SELECT
			b.id, 
			b.title, 
			b.author_id, 
			b.publisher_id,
			b.category_id, 
			b.published_year, 
			b.isbn, 
			b.status, 
			b.access_level, 
			b.created_at, 
			b.created_by, 
			b.updated_at, 
			b.updated_by
		FROM
			books b
		INNER JOIN 
			authors a 
		ON 
			b.author_id = a.id
		INNER JOIN 
			publishers p 
		ON 
			b.publisher_id = p.id
		INNER JOIN 
			categories c 
		ON 
			b.category_id = c.id
		%s
		ORDER BY b.created_at %s
		LIMIT $1 OFFSET $2
	`, lib.SearchGenerate(ctx, "AND", filters), filter.Dir)

	rows, err := db.QueryxContext(ctx, query, filter.Limit, filter.Offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var books []BookModel
	for rows.Next() {
		var book BookModel
		err := rows.StructScan(&book)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, nil
}

func GetOneBooks(ctx context.Context, db *sqlx.DB, id uuid.UUID) (BookModel, error) {
	query := `
		SELECT
			id, title, author_id, publisher_id, category_id, published_year, isbn, status, access_level, created_at, created_by, updated_at, updated_by
		FROM
			books
		WHERE 
			id = $1
	`

	bookRequest := BookModel{}
	err := db.QueryRowxContext(ctx, query, id).StructScan(&bookRequest)
	if err != nil {
		return bookRequest, err
	}

	return bookRequest, nil
}

func (b *BookModel) Insert(ctx context.Context, db *sqlx.DB) error {
	query := `
		INSERT INTO books (
			id, title, author_id, publisher_id, category_id, published_year, isbn, status, access_level, created_at, created_by, updated_at, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		) RETURNING id, created_at
	`

	err := db.QueryRowxContext(ctx, query,
		b.ID,
		b.Title,
		b.AuthorID,
		b.PublisherID.UUID,
		b.CategoryID.UUID,
		b.PublishedYear,
		b.ISBN,
		b.Status,
		b.AccessLevel,
		b.CreatedAt,
		b.CreatedBy,
		b.UpdatedAt.Time,
		b.UpdatedBy.UUID,
	).Scan(
		&b.ID,
		&b.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (b *BookModel) Update(ctx context.Context, db *sqlx.DB, id uuid.UUID) error {
	query := `
		UPDATE
			books
		SET
			title = $1, 
			author_id = $2, 
			publisher_id = $3, 
			category_id = $4, 
			published_year = $5, 
			isbn = $6, 
			status = $7, 
			access_level = $8,
			updated_at = NOW(),
			updated_by = $9
		WHERE id = $10
		RETURNING id, title, author_id, publisher_id, category_id, published_year, isbn, status, access_level, updated_at, updated_by
	`

	err := db.QueryRowxContext(ctx, query,
		b.Title,
		b.AuthorID,
		b.PublisherID.UUID,
		b.CategoryID.UUID,
		b.PublishedYear,
		b.ISBN,
		b.Status,
		b.AccessLevel,
		b.UpdatedBy.UUID,
		b.ID,
	).Scan(
		&b.ID,
		&b.Title,
		&b.AuthorID,
		&b.PublisherID.UUID,
		&b.CategoryID.UUID,
		&b.PublishedYear,
		&b.ISBN,
		&b.Status,
		&b.AccessLevel,
		&b.UpdatedAt.Time,
		&b.UpdatedBy.UUID,
	)

	if err != nil {
		return err
	}

	return nil
}
