package api

import (
	"app-bookstore/lib"
	"app-bookstore/model"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type RatingModule struct {
	db   *sqlx.DB
	name string
	JWT  lib.Jwt
}

func NewRatingModule(db *sqlx.DB, jwt lib.Jwt) *RatingModule {
	return &RatingModule{
		db:   db,
		name: "ratings-module",
		JWT:  jwt,
	}
}

type RatingParam struct {
	BookID uuid.UUID `json:"book_id"`
	Rating int       `json:"rating"`
	Review string    `json:"review"`
}

type RatingUpdateParam struct {
	Rating int    `json:"rating"`
	Review string `json:"review"`
}

func (rt *RatingModule) List(ctx context.Context) ([]model.RatingResponse, error) {
	ratingRequest, err := model.GetAllRatings(ctx, rt.db)
	if err != nil {
		return nil, err
	}

	var response []model.RatingResponse
	for _, rating := range ratingRequest {
		response = append(response, rating.Response())
	}

	return response, nil
}

func (rt *RatingModule) Detail(ctx context.Context, id uuid.UUID) (model.RatingResponse, error) {
	ratingRequest, err := model.GetOneRating(ctx, rt.db, id)
	if err != nil {
		return model.RatingResponse{}, err
	}
	return ratingRequest.Response(), nil
}

func (rt *RatingModule) Create(ctx context.Context, token string, param RatingParam) (interface{}, error) {
	claims, err := rt.JWT.VerifyAccessToken(token)
	if err != nil {
		return nil, errors.New("invalid or expired access token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("failed to parse access token")
	}

	rating := model.RatingModel{
		ID:        uuid.New(),
		BookID:    param.BookID,
		UserID:    userID,
		Rating:    param.Rating,
		Review:    param.Review,
		CreatedAt: time.Now(),
		CreatedBy: userID,
	}

	err = rating.Insert(ctx, rt.db)
	if err != nil {
		return nil, err
	}

	return rating.Response(), nil
}

func (rt *RatingModule) Update(ctx context.Context, token string, param RatingUpdateParam, id uuid.UUID) (interface{}, error) {
	claims, err := rt.JWT.VerifyAccessToken(token)
	if err != nil {
		return nil, errors.New("invalid or expired access token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("failed to parse access token")
	}

	rating := model.RatingModel{
		ID:     id,
		Rating: param.Rating,
		Review: param.Review,
		UpdatedAt: pq.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedBy: uuid.NullUUID{
			UUID:  userID,
			Valid: true,
		},
	}

	err = rating.Update(ctx, rt.db, id)
	if err != nil {
		return nil, err
	}

	return rating.Response(), nil

}
