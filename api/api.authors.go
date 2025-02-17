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

type AuthorsModule struct {
	db   *sqlx.DB
	name string
	JWT  lib.Jwt
}

func NewUserAuthorsModule(db *sqlx.DB, jwt lib.Jwt) *AuthorsModule {
	return &AuthorsModule{
		db:   db,
		name: "authors-module",
		JWT:  jwt,
	}
}

type AuthorsParam struct {
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

func (a *AuthorsModule) List(ctx context.Context, filter lib.Filter, dateFilter model.DateFilter) ([]model.AuthorsRespose, error) {
	authorRequest, err := model.GetAllAuthors(ctx, a.db, filter, dateFilter)
	if err != nil {
		return nil, err
	}

	var response []model.AuthorsRespose
	for _, author := range authorRequest {
		response = append(response, author.Response())
	}
	return response, err
}

func (a *AuthorsModule) Detail(ctx context.Context, id uuid.UUID) (model.AuthorsRespose, error) {
	authorRequest, err := model.GetOneAuthors(ctx, a.db, id)
	if err != nil {
		return model.AuthorsRespose{}, err
	}
	return authorRequest.Response(), nil
}

func (a *AuthorsModule) Create(ctx context.Context, token string, param AuthorsParam) (interface{}, error) {
	claims, err := a.JWT.VerifyAccessToken(token)
	if err != nil {
		return nil, errors.New("invalid or expired access token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("failed to parse access token")
	}

	author := model.AuthorsModel{
		ID:        uuid.New(),
		Name:      param.Name,
		Bio:       param.Bio,
		CreatedAt: time.Now(),
		CreatedBy: userID,
	}

	err = author.Insert(ctx, a.db)
	if err != nil {
		return nil, err
	}

	return author.Response(), nil
}

func (a *AuthorsModule) Update(ctx context.Context, token string, param AuthorsParam, id uuid.UUID) (interface{}, error) {
	claims, err := a.JWT.VerifyAccessToken(token)
	if err != nil {
		return nil, errors.New("invalid or expired access token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("failed to parse access token")
	}

	authors := model.AuthorsModel{
		ID:   id,
		Name: param.Name,
		Bio:  param.Bio,
		UpdatedAt: pq.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedBy: uuid.NullUUID{
			UUID:  userID,
			Valid: true,
		},
	}

	err = authors.Update(ctx, a.db, id)
	if err != nil {
		return nil, err
	}

	return authors.Response(), nil
}
