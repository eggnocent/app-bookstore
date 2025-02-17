package model

import (
	"app-bookstore/lib"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type PublisherModel struct {
	ID        uuid.UUID     `db:"id"`
	Name      string        `db:"name"`
	Address   string        `db:"address"`
	Phone     string        `db:"phone"`
	CreatedAt time.Time     `db:"created_at"`
	CreatedBy uuid.UUID     `db:"created_by"`
	UpdatedAt pq.NullTime   `db:"updated_at"`
	UpdatedBy uuid.NullUUID `db:"updated_by"`
}

type PublisherResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy uuid.UUID `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy uuid.UUID `json:"updated_by"`
}

func (p *PublisherModel) Response() PublisherResponse {
	return PublisherResponse{
		ID:        p.ID,
		Name:      p.Name,
		Address:   p.Address,
		Phone:     p.Phone,
		CreatedAt: p.CreatedAt,
		CreatedBy: p.CreatedBy,
		UpdatedAt: p.UpdatedAt.Time,
		UpdatedBy: p.UpdatedBy.UUID,
	}
}

func GetAllPublisher(ctx context.Context, db *sqlx.DB, filter lib.Filter, dateFilter DateFilter) ([]PublisherModel, error) {
	var filters []string
	if filter.Search != "" {
		filters = append(filters, fmt.Sprintf("name ILIKE '%%%s%%'", filter.Search))
	}

	if !dateFilter.StartDate.IsZero() && !dateFilter.EndDate.IsZero() {
		filters = append(filters, fmt.Sprintf(
			"created_at BETWEEN '%s' AND '%s'",
			dateFilter.StartDate.Format("2006-01-02"),
			dateFilter.EndDate.Format("2006-01-02"),
		))
	}

	query := fmt.Sprintf(
		`
		SELECT
			id, 
			name, 
			address, 
			phone, 
			created_at, 
			created_by, 
			updated_at, 
			updated_by
		FROM
			publishers 
		%s
		ORDER BY created_at %s
		LIMIT $1 OFFSET $2
	`, lib.SearchGenerate(ctx, "AND", filters), filter.Dir)

	rows, err := db.QueryxContext(ctx, query, filter.Limit, filter.Offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var publishers []PublisherModel
	for rows.Next() {
		var publisher PublisherModel
		err := rows.StructScan(&publisher)
		if err != nil {
			return nil, err
		}

		publishers = append(publishers, publisher)
	}

	return publishers, nil
}

func GetOnePublisher(ctx context.Context, db *sqlx.DB, id uuid.UUID) (PublisherModel, error) {
	query := `
		SELECT
			id, 
            name, 
            address, 
            phone, 
			created_at, 
			created_by, 
			updated_at,
			updated_by
		FROM
			publishers
		WHERE 
			id = $1
	`
	publisherRequest := PublisherModel{}
	err := db.QueryRowxContext(ctx, query, id).StructScan(&publisherRequest)
	if err != nil {
		return publisherRequest, err
	}

	return publisherRequest, nil
}

func (p *PublisherModel) Insert(ctx context.Context, db *sqlx.DB) error {
	query := `
		INSERT INTO publishers (
		    id, 
            name, 
            address, 
            phone,
			created_at, 
            created_by
		) VALUES (
			$1,$2,$3,$4,$5,$6 
		) RETURNING id, created_at
	`

	err := db.QueryRowxContext(ctx, query,
		p.ID,
		p.Name,
		p.Address,
		p.Phone,
		p.CreatedAt,
		p.CreatedBy,
	).Scan(
		&p.ID,
		&p.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (p *PublisherModel) Update(ctx context.Context, db *sqlx.DB) error {
	query := `
		UPDATE
			publishers
		SET
			name = $1,
			address = $2,
            phone = $3,
            updated_at = $4,
            updated_by = $5
		WHERE
			id = $6
		RETURNING id, name, address, phone, created_at, created_by, updated_at, updated_by
	`

	err := db.QueryRowxContext(ctx, query,
		p.Name,
		p.Address,
		p.Phone,
		p.UpdatedAt,
		p.UpdatedBy.UUID,
		p.ID,
	).Scan(
		&p.ID,
		&p.Name,
		&p.Address,
		&p.Phone,
		&p.CreatedAt,
		&p.CreatedBy,
		&p.UpdatedAt,
		&p.UpdatedBy,
	)

	if err != nil {
		return err
	}

	return nil
}
