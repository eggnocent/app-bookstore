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

type CategoryModule struct {
	db   *sqlx.DB
	name string
	JWT  lib.Jwt
}

func NewCategoriesModule(db *sqlx.DB, jwt lib.Jwt) *CategoryModule {
	return &CategoryModule{
		db:   db,
		name: "categories-module",
		JWT:  jwt,
	}
}

type CategoryParam struct {
	Name string `json:"name"`
}

func (c *CategoryModule) List(ctx context.Context) ([]model.CategoryResponse, error) {
	categoriesRequest, err := model.GetAllCategory(ctx, c.db)
	if err != nil {
		return nil, err
	}

	var response []model.CategoryResponse
	for _, category := range categoriesRequest {
		response = append(response, category.Response())
	}

	return response, nil
}

func (c *CategoryModule) Detail(ctx context.Context, id uuid.UUID) (model.CategoryResponse, error) {
	categoryResponse, err := model.GetOneCategories(ctx, c.db, id)
	if err != nil {
		return model.CategoryResponse{}, err
	}

	return categoryResponse.Response(), nil
}

func (c *CategoryModule) Create(ctx context.Context, token string, param CategoryParam) (interface{}, error) {
	claims, err := c.JWT.VerifyAccessToken(token)
	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("invalid user id in token")
	}

	category := model.CategoryModel{
		ID:        uuid.New(),
		Name:      param.Name,
		CreatedAt: time.Now(),
		CreatedBy: userID,
	}

	err = category.Insert(ctx, c.db)
	if err != nil {
		return nil, err
	}

	return category.Response(), nil
}

func (c *CategoryModule) Update(ctx context.Context, token string, param CategoryParam, id uuid.UUID) (interface{}, error) {
	claims, err := c.JWT.VerifyAccessToken(token)
	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("invalid user id in token")
	}

	category := model.CategoryModel{
		ID:   id,
		Name: param.Name,
		UpdatedAt: pq.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedBy: uuid.NullUUID{
			UUID:  userID,
			Valid: true,
		},
	}

	err = category.Update(ctx, c.db, id)
	if err != nil {
		return nil, err
	}

	return category.Response(), nil
}
