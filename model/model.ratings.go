package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type RatingModel struct {
	ID        uuid.UUID     `db:"id"`
	BookID    uuid.UUID     `db:"book_id"`
	UserID    uuid.UUID     `db:"user_id"`
	Rating    int           `db:"rating"`
	Review    string        `db:"review"`
	CreatedAt time.Time     `db:"created_at"`
	CreatedBy uuid.UUID     `db:"created_by"`
	UpdatedAt pq.NullTime   `db:"updated_at"`
	UpdatedBy uuid.NullUUID `db:"updated_by"`
}

type RatingResponse struct {
	ID        uuid.UUID `json:"id"`
	BookID    uuid.UUID `json:"book_id"`
	UserID    uuid.UUID `json:"user_id"`
	Rating    int       `json:"rating"`
	Review    string    `json:"review"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy uuid.UUID `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy uuid.UUID `json:"updated_by"`
}

func (rt *RatingModel) Response() RatingResponse {
	return RatingResponse{
		ID:        rt.ID,
		BookID:    rt.BookID,
		UserID:    rt.UserID,
		Rating:    rt.Rating,
		Review:    rt.Review,
		CreatedAt: rt.CreatedAt,
		CreatedBy: rt.CreatedBy,
		UpdatedAt: rt.UpdatedAt.Time,
		UpdatedBy: rt.UpdatedBy.UUID,
	}
}

func GetAllRatings(ctx context.Context, db *sqlx.DB) ([]RatingModel, error) {
	query := `
		SELECT
			id, book_id, user_id, rating, review, created_at, created_by, updated_at, updated_by
		FROM
			ratings
	`

	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var ratings []RatingModel
	for rows.Next() {
		var rating RatingModel
		err := rows.StructScan(&rating)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, rating)
	}

	return ratings, nil
}

func GetOneRating(ctx context.Context, db *sqlx.DB, id uuid.UUID) (RatingModel, error) {
	query := `
        SELECT
            id, book_id, user_id, rating, review, created_at, created_by, updated_at, updated_by
        FROM 
            ratings
        WHERE id = $1
    `

	var rating RatingModel
	err := db.QueryRowxContext(ctx, query, id).StructScan(&rating)
	if err != nil {
		return RatingModel{}, err
	}

	return rating, nil
}

func (rt *RatingModel) Insert(ctx context.Context, db *sqlx.DB) error {
	query := `
		INSERT INTO ratings (
		    id, book_id, user_id, rating, review, created_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7 
		) RETURNING id, created_at
	`

	err := db.QueryRowxContext(ctx, query,
		rt.ID,
		rt.BookID,
		rt.UserID,
		rt.Rating,
		rt.Review,
		rt.CreatedAt,
		rt.CreatedBy,
	).Scan(
		&rt.ID,
		&rt.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (rt *RatingModel) Update(ctx context.Context, db *sqlx.DB, id uuid.UUID) error {
	query := `
		UPDATE 
			ratings
		SET
			rating = $1,
			review = $2,
            updated_at = $3,
            updated_by = $4
		WHERE
			id = $5
		RETURNING id, book_id, user_id, rating, review, created_at, created_by, updated_at, updated_by
	`

	err := db.QueryRowxContext(ctx, query,
		rt.Rating,
		rt.Review,
		rt.UpdatedAt.Time,
		rt.UpdatedBy.UUID,
		rt.ID,
	).Scan(
		&rt.ID,
		&rt.BookID,
		&rt.UserID,
		&rt.Rating,
		&rt.Review,
		&rt.CreatedAt,
		&rt.CreatedBy,
		&rt.UpdatedAt,
		&rt.UpdatedBy,
	)

	if err != nil {
		return err
	}

	return nil
}
