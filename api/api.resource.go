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

type ResourceModule struct {
	db   *sqlx.DB
	name string
	JWT  lib.Jwt
}

func NewResourceModule(db *sqlx.DB, jwt lib.Jwt) *ResourceModule {
	return &ResourceModule{
		db:   db,
		name: "resource-module",
		JWT:  jwt,
	}
}

type ResourceParam struct {
	Name        string `json:"name"`
	Endpoint    string `json:"endpoint"`
	Method      string `json:"method"`
	Description string `json:"description"`
}

func (rs *ResourceModule) List(ctx context.Context) ([]model.ResourceResponse, error) {
	resources, err := model.GetAllResources(ctx, rs.db)
	if err != nil {
		return nil, err
	}

	var response []model.ResourceResponse
	for _, resource := range resources {
		response = append(response, resource.Response())
	}
	return response, nil
}

func (rs *ResourceModule) Detail(ctx context.Context, id uuid.UUID) (model.ResourceResponse, error) {
	resources, err := model.GetOneResource(ctx, rs.db, id)
	if err != nil {
		return model.ResourceResponse{}, err
	}

	return resources.Response(), nil
}

func (rs *ResourceModule) Create(ctx context.Context, token string, param ResourceParam) (interface{}, error) {
	claims, err := rs.JWT.VerifyAccessToken(token)
	if err != nil {
		return nil, errors.New("failed to verify access token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("failed to parse access token")
	}

	resourceID := uuid.New()
	resource := model.ResourceModel{
		ID:          resourceID,
		Name:        param.Name,
		Endpoint:    param.Endpoint,
		Method:      param.Method,
		Description: param.Description,
		CreatedAt:   time.Now(),
		CreatedBy:   userID,
	}

	err = resource.Insert(ctx, rs.db)
	if err != nil {
		return nil, err
	}

	return resource.Response(), nil
}

func (rs *ResourceModule) Update(ctx context.Context, token string, param ResourceParam, id uuid.UUID) (interface{}, error) {
	claims, err := rs.JWT.VerifyAccessToken(token)
	if err != nil {
		return nil, errors.New("invalid or expired access token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("failed to parse access token")
	}

	resource := model.ResourceModel{
		ID:          id,
		Name:        param.Name,
		Endpoint:    param.Endpoint,
		Method:      param.Method,
		Description: param.Description,
		UpdatedAt: pq.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedBy: uuid.NullUUID{
			UUID:  userID,
			Valid: true,
		},
	}

	err = resource.Update(ctx, rs.db)
	if err != nil {
		return nil, err
	}

	return resource.Response(), nil
}
